package controller

import (
	"context"
	"time"

	"backend-sisteminformasi/config"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// StatisticsResponse represents the statistics data structure
type StatisticsResponse struct {
	TotalKegiatan    int64            `json:"totalKegiatan"`
	TotalAnggota     int64            `json:"totalAnggota"`
	TotalKehadiran   int64            `json:"totalKehadiran"`
	KegiatanByStatus map[string]int64 `json:"kegiatanByStatus"`
	KegiatanByUkm    []UkmStats       `json:"kegiatanByUkm"`
	MembersByUkm     []MemberStats    `json:"membersByUkm"`
	RecentActivities []ActivityStats  `json:"recentActivities"`
}

type UkmStats struct {
	UKM        string  `json:"ukm"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}

type MemberStats struct {
	UKM         string `json:"ukm"`
	AdminCount  int64  `json:"adminCount"`
	MemberCount int64  `json:"memberCount"`
	Total       int64  `json:"total"`
}

type ActivityStats struct {
	Title     string `json:"title"`
	UKM       string `json:"ukm"`
	Date      string `json:"date"`
	Attendees int64  `json:"attendees"`
}

// @Security BearerAuth
// GetStatistics godoc
// @Summary Get statistics data
// @Tags Statistics
// @Produce json
// @Success 200 {object} StatisticsResponse
// @Failure 500 {object} map[string]interface{}
// @Router /statistics [get]
func GetStatistics(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var stats StatisticsResponse

	// Get total counts
	totalKegiatan, err := config.DB.Collection("kegiatan").CountDocuments(ctx, bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to count kegiatan"})
	}
	stats.TotalKegiatan = totalKegiatan

	totalAnggota, err := config.DB.Collection("users").CountDocuments(ctx, bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to count users"})
	}
	stats.TotalAnggota = totalAnggota

	totalKehadiran, err := config.DB.Collection("kehadiran").CountDocuments(ctx, bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to count kehadiran"})
	}
	stats.TotalKehadiran = totalKehadiran

	// Get kegiatan by status
	pipeline := []bson.M{
		{
			"$group": bson.M{
				"_id":   "$status",
				"count": bson.M{"$sum": 1},
			},
		},
	}

	cursor, err := config.DB.Collection("kegiatan").Aggregate(ctx, pipeline)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to aggregate kegiatan by status"})
	}

	var statusResults []bson.M
	if err := cursor.All(ctx, &statusResults); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to decode status results"})
	}

	stats.KegiatanByStatus = make(map[string]int64)
	stats.KegiatanByStatus["upcoming"] = 0
	stats.KegiatanByStatus["ongoing"] = 0
	stats.KegiatanByStatus["completed"] = 0
	stats.KegiatanByStatus["cancelled"] = 0

	for _, result := range statusResults {
		if status, ok := result["_id"].(string); ok {
			if count, ok := result["count"].(int32); ok {
				stats.KegiatanByStatus[status] = int64(count)
			}
		}
	}

	// Get kegiatan by UKM
	pipeline = []bson.M{
		{
			"$group": bson.M{
				"_id":   "$kategori",
				"count": bson.M{"$sum": 1},
			},
		},
		{
			"$sort": bson.M{"count": -1},
		},
	}

	cursor, err = config.DB.Collection("kegiatan").Aggregate(ctx, pipeline)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to aggregate kegiatan by UKM"})
	}

	var ukmResults []bson.M
	if err := cursor.All(ctx, &ukmResults); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to decode UKM results"})
	}

	stats.KegiatanByUkm = make([]UkmStats, 0)
	for _, result := range ukmResults {
		if ukm, ok := result["_id"].(string); ok && ukm != "" {
			if count, ok := result["count"].(int32); ok {
				percentage := float64(0)
				if totalKegiatan > 0 {
					percentage = (float64(count) / float64(totalKegiatan)) * 100
				}
				stats.KegiatanByUkm = append(stats.KegiatanByUkm, UkmStats{
					UKM:        ukm,
					Count:      int64(count),
					Percentage: percentage,
				})
			}
		}
	}

	// Get members by UKM
	pipeline = []bson.M{
		{
			"$group": bson.M{
				"_id": bson.M{
					"ukm":  "$ukm",
					"role": "$role",
				},
				"count": bson.M{"$sum": 1},
			},
		},
	}

	cursor, err = config.DB.Collection("users").Aggregate(ctx, pipeline)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to aggregate members by UKM"})
	}

	var memberResults []bson.M
	if err := cursor.All(ctx, &memberResults); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to decode member results"})
	}

	memberMap := make(map[string]*MemberStats)
	for _, result := range memberResults {
		if idMap, ok := result["_id"].(bson.M); ok {
			if ukm, ok := idMap["ukm"].(string); ok && ukm != "" {
				if role, ok := idMap["role"].(string); ok {
					if count, ok := result["count"].(int32); ok {
						if _, exists := memberMap[ukm]; !exists {
							memberMap[ukm] = &MemberStats{UKM: ukm}
						}

						if role == "admin" {
							memberMap[ukm].AdminCount = int64(count)
						} else {
							memberMap[ukm].MemberCount = int64(count)
						}
						memberMap[ukm].Total = memberMap[ukm].AdminCount + memberMap[ukm].MemberCount
					}
				}
			}
		}
	}

	stats.MembersByUkm = make([]MemberStats, 0)
	for _, memberStat := range memberMap {
		stats.MembersByUkm = append(stats.MembersByUkm, *memberStat)
	}

	// Get recent activities
	pipeline = []bson.M{
		{
			"$sort": bson.M{"tanggal": -1},
		},
		{
			"$limit": 5,
		},
		{
			"$lookup": bson.M{
				"from": "kehadiran",
				"let":  bson.M{"kegiatan_id": bson.M{"$toString": "$_id"}},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$eq": []interface{}{"$kegiatan_id", "$$kegiatan_id"},
							},
						},
					},
					{
						"$count": "total",
					},
				},
				"as": "kehadiran_count",
			},
		},
		{
			"$project": bson.M{
				"judul":    1,
				"kategori": 1,
				"tanggal":  1,
				"attendees_count": bson.M{
					"$ifNull": []interface{}{
						bson.M{"$arrayElemAt": []interface{}{"$kehadiran_count.total", 0}},
						0,
					},
				},
			},
		},
	}

	cursor, err = config.DB.Collection("kegiatan").Aggregate(ctx, pipeline)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get recent activities"})
	}

	var activityResults []bson.M
	if err := cursor.All(ctx, &activityResults); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to decode activity results"})
	}

	stats.RecentActivities = make([]ActivityStats, 0)
	for _, result := range activityResults {
		activity := ActivityStats{}

		if title, ok := result["judul"].(string); ok {
			activity.Title = title
		}
		if ukm, ok := result["kategori"].(string); ok {
			activity.UKM = ukm
		}
		if date, ok := result["tanggal"].(string); ok {
			activity.Date = date
		}
		if attendees, ok := result["attendees_count"].(int32); ok {
			activity.Attendees = int64(attendees)
		}

		stats.RecentActivities = append(stats.RecentActivities, activity)
	}

	return c.JSON(stats)
}

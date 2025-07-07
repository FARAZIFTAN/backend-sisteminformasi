package model

type Kehadiran struct {
	ID         string `bson:"_id,omitempty" json:"id"`
	UserID     string `bson:"user_id" json:"user_id" validate:"required"`
	KegiatanID string `bson:"kegiatan_id" json:"kegiatan_id" validate:"required"`
	Status     string `bson:"status" json:"status" validate:"required,oneof=hadir tidak"`
	WaktuCek   string `bson:"waktu_cek" json:"waktu_cek"`
}

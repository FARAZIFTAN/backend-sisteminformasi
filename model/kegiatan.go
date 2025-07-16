package model

type Kegiatan struct {
	ID              string `bson:"_id,omitempty" json:"id"`
	Judul           string `bson:"judul" json:"judul" validate:"required"`
	Deskripsi       string `bson:"deskripsi" json:"deskripsi"`
	Tanggal         string `bson:"tanggal" json:"tanggal" validate:"required"`
	Lokasi          string `bson:"lokasi" json:"lokasi"`
	Kategori        string `bson:"kategori" json:"kategori"`
	MaxParticipants int    `bson:"maxParticipants" json:"maxParticipants"`
	DokumentasiURL  string `bson:"dokumentasi_url" json:"dokumentasi_url"`
	CreatedBy       string `bson:"created_by" json:"created_by"`
}

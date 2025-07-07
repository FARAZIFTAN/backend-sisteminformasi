package model

type Kegiatan struct {
	ID             string `bson:"_id,omitempty" json:"id"`
	Judul          string `bson:"judul" json:"judul" validate:"required"`
	Deskripsi      string `bson:"deskripsi" json:"deskripsi"`
	Tanggal        string `bson:"tanggal" json:"tanggal" validate:"required"`
	Lokasi         string `bson:"lokasi" json:"lokasi"`
	DokumentasiURL string `bson:"dokumentasi_url" json:"dokumentasi_url"`
	Kategori       string `bson:"kategori" json:"kategori"`
	CreatedBy      string `bson:"created_by" json:"created_by"`
}

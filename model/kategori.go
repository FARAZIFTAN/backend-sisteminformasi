package model

type Kategori struct {
	ID           string `bson:"_id,omitempty" json:"id"`
	NamaKategori string `bson:"nama_kategori" json:"nama_kategori" validate:"required"`
}

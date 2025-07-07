package model

type User struct {
	ID       string `bson:"_id,omitempty" json:"id"`
	Nama     string `bson:"nama" json:"nama" validate:"required,min=2,max=100"`
	Email    string `bson:"email" json:"email" validate:"required,email"`
	Password string `bson:"password" json:"password" validate:"required,min=6"`
	Role     string `bson:"role" json:"role" validate:"required,oneof=admin member"`
	UKM      string `bson:"ukm" json:"ukm" validate:"required"`
}

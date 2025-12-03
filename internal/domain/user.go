package domain

type User struct {
	ID       int64  `json:"id" gorm:"autoIncrement;primaryKey"`
	Username string `json:"username" gorm:"type:varchar(45);unique;not null"`
	Password string `json:"password" gorm:"type:varchar(45);not null"`
	Role     string `json:"role" gorm:"type:varchar(45)"`
}

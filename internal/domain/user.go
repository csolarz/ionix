package domain

type User struct {
	Username string `json:"username" gorm:"primaryKey"`
	Password string `json:"password" gorm:"type:varchar(45);not null"`
	Role     string `json:"role" gorm:"type:varchar(45)"`
}

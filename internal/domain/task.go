package domain

type Task struct {
	ID          int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string `json:"title" gorm:"type:varchar(45);not null"`
	Description string `json:"description" gorm:"type:varchar(255)"`
	DueDate     string `json:"due_date" binding:"required,datetime=2006-01-02T15:04:05Z07:00" gorm:"type:datetime"`
}

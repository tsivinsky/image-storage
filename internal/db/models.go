package db

import "time"

type Image struct {
	Id        uint      `json:"id" gorm:"primaryKey"`
	Filename  string    `json:"filename" gorm:"filename"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

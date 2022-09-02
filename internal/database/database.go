package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Message struct {
	ID       uint `gorm:"primaryKey"`
	Username string
	Content  string
	Time     string
	Color    int
	To       string
	RoomCode string `gorm:"column: room-code"`
}

func NewDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("termchat.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(Message{})
	return db
}

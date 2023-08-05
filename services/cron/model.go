package cron

import "time"

type MangaUpdate struct {
	ID                 string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name               string    `gorm:"not null"`
	RawURL             string    `gorm:"not null"`
	NegativeIdentifier string    `gorm:"null"`
	LastChapter        int       `gorm:"not null"`
	LastCheckedAt      time.Time `gorm:"null"`
	ChatID             int       `gorm:"not null"`
	UpdatedAt          time.Time `gorm:"not null"`
	CreatedAt          time.Time `gorm:"not null"`
}

type UpdateMangaPayload struct {
	ID            string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name          string    `gorm:"not null"`
	RawURL        string    `gorm:"not null"`
	LastChapter   int       `gorm:"not null"`
	LastCheckedAt time.Time `gorm:"null"`
	UpdatedAt     time.Time `gorm:"not null"`
	CreatedAt     time.Time `gorm:"not null"`
}

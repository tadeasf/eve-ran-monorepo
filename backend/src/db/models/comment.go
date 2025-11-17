package models

import (
	"time"
)

// KillComment represents a comment on a kill
type KillComment struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	KillID     uint      `gorm:"index:idx_kill_comments_kill_id" json:"kill_id"`
	KillmailID int64     `gorm:"index:idx_kill_comments_killmail_id" json:"killmail_id"`
	UserID     int64     `json:"user_id"` // For future user management
	Username   string    `json:"username"`
	Comment    string    `gorm:"type:text" json:"comment"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

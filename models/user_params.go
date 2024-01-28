package models

import "time"

type UserInfo struct {
	LastActivity time.Time
	UserID       int
	ChatID       int64
}

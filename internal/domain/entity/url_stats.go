package entity

import (
	"time"
)

type URLStat struct {
	ID        string
	URLID     string
	ClickedAt time.Time
	IP        string
	UserAgent string
	Referer   string
}

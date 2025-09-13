package dto

import "time"

type Resume struct {
	Clicks    int       `json:"clicks"`
	LastClick time.Time `json:"last_click"`
}

type Data struct {
	ID        string    `json:"id"`
	URLID     string    `json:"url_id"`
	ClickedAt time.Time `json:"clicked_at"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	Referer   string    `json:"referer"`
}

type URLStats struct {
	StatsResume Resume `json:"resume"`
	StatsData   []Data `json:"data"`
}

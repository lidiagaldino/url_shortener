package repository

import (
	"url-shortener/internal/domain/entity"
)

type URLStatsRepository interface {
	Save(stat *entity.URLStat) error
}

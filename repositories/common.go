package repositories

import (
	"github.com/hodukihugi/winglets-api/models"
	"gorm.io/gorm"
)

func paginate(tx *gorm.DB, pagination *models.Pagination) {
	if pagination != nil {
		tx.Limit(pagination.PerPage).Offset((pagination.CurrentPage - 1) * pagination.PerPage)
	}
}

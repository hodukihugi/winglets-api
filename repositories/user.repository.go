package repositories

import (
	"github.com/hodukihugi/winglets-api/core"
	"github.com/hodukihugi/winglets-api/models"
	"gorm.io/gorm"
	"strings"
)

type IUserRepository interface {
	Create(models.User) error
	First(models.OneUserFilter) (*models.User, error)
}

// UserRepository database structure
type UserRepository struct {
	*core.Database
	logger *core.Logger
}

func (r *UserRepository) First(filter models.OneUserFilter) (user *models.User, err error) {
	tx := r.Table("users")
	r.filterUser(filter, tx)
	return user, tx.First(&user).Error
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *core.Database, logger *core.Logger) IUserRepository {
	return &UserRepository{
		Database: db,
		logger:   logger,
	}
}

func (r *UserRepository) Create(user models.User) error {
	user.Email = strings.ToLower(user.Email)
	return r.DB.Create(&user).Error
}

// -------- Private functions ---------
func (r *UserRepository) filterUser(filter models.OneUserFilter, tx *gorm.DB) {
	if filter.Fields != nil && len(filter.Fields.Values()) > 0 {
		tx.Select(filter.Fields.Values())
	}
	for _, join := range filter.Joins.Values() {
		if join == "orgs" {
			tx.Joins("Orgs")
		}
		if join == "user" {
			tx.Joins("User")
		}
		if join == "referrer" {
			tx.Joins("Referrer")
		}
	}
	if filter.Email != "" {
		tx.Where("users.email = ?", filter.Email)
	}
	if filter.ID != 0 {
		tx.Where("users.id = ?", filter.ID)
	}
}

package repositories

import "github.com/hodukihugi/winglets-api/core"

type IMatchRepository interface {
	Temp()
}

type MatchRepository struct {
	*core.Database
	logger *core.Logger
}

func NewMatchRepository(db *core.Database, logger *core.Logger) IMatchRepository {
	return &MatchRepository{
		Database: db,
		logger:   logger,
	}
}

func (m *MatchRepository) Temp() {

}

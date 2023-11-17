package repositories

import (
	"database/sql"

	"github.com/go-redis/redis/v8"
)

type Repositories struct {
	DB                       *sql.DB
	RDB                      *redis.Client
	GameSettingsRepository   GameSettingsRepository
	UserRepository           UserRepository
	ItemRepository           ItemRepository
	CollectionItemRepository CollectionItemRepository
	UserScoresRepository     UserScoresRepository
}

func NewRepositories(db *sql.DB, rdb *redis.Client) *Repositories {
	return &Repositories{
		DB:                       db,
		RDB:                      rdb,
		GameSettingsRepository:   NewGameSettingsRepository(db, rdb),
		UserRepository:           NewUserRepository(db),
		ItemRepository:           NewItemRepository(db, rdb),
		CollectionItemRepository: NewCollectionItemRepository(db),
		UserScoresRepository:     NewUserScoresRepository(db),
	}
}

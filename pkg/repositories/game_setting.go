package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"42tokyo-road-to-dojo-go/pkg/server/entities"

	"github.com/go-redis/redis/v8"
)

type GameSettingsRepository interface {
	GetAllGameSettings() (*[]entities.GameSettings, error)
	GetGameSettingsByID(ID entities.GameSettingID) (*entities.GameSettings, error)
	GetActiveGameSettings() (*entities.GameSettings, error)
	CacheActiveGameSettings() error
	GetActiveGameSettingsFromCache() (*entities.GameSettings, error)
	ChangeActiveGameSettings(ID entities.GameSettingID) error
	AddGameSettings(settings entities.GameSettings) error
}

func NewGameSettingsRepository(db *sql.DB, rdb *redis.Client) GameSettingsRepository {
	return &gameSettingsRepository{db, rdb}
}

type gameSettingsRepository struct {
	db  *sql.DB
	rdb *redis.Client
}

func (r *gameSettingsRepository) GetAllGameSettings() (*[]entities.GameSettings, error) {
	query := "SELECT * FROM game_settings"
	rows, err := r.db.Query(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var settings []entities.GameSettings
	for rows.Next() {
		var setting entities.GameSettings
		var createdAt []byte
		if err := rows.Scan(&setting.ID, &setting.GachaCoinConsumption, &setting.RankingListLimit, &setting.NWeight, &setting.RWeight, &setting.SrWeight, &setting.MaxGachaTimes, &createdAt, &setting.IsActive); err != nil {
			log.Println(err)
			return nil, err
		}
		if setting.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt)); err != nil {
			log.Println(err)
			return nil, err
		}

		settings = append(settings, setting)
	}
	return &settings, nil
}

func (r *gameSettingsRepository) GetGameSettingsByID(ID entities.GameSettingID) (*entities.GameSettings, error) {
	query := "SELECT * FROM game_settings WHERE id = ? LIMIT 1"
	row := r.db.QueryRow(query, ID)

	var setting entities.GameSettings
	var createdAt []byte
	if err := row.Scan(&setting.ID, &setting.GachaCoinConsumption, &setting.RankingListLimit, &setting.NWeight, &setting.RWeight, &setting.SrWeight, &setting.MaxGachaTimes, &createdAt, &setting.IsActive); err != nil {
		log.Println(err)
		return nil, err
	}

	return &setting, nil
}

func (r *gameSettingsRepository) GetActiveGameSettings() (*entities.GameSettings, error) {
	query := "SELECT * FROM game_settings WHERE is_active = true LIMIT 1"
	row := r.db.QueryRow(query)

	var setting entities.GameSettings
	var createdAt []byte
	if err := row.Scan(&setting.ID, &setting.GachaCoinConsumption, &setting.RankingListLimit, &setting.NWeight, &setting.RWeight, &setting.SrWeight, &setting.MaxGachaTimes, &createdAt, &setting.IsActive); err != nil {
		log.Println(err)
		return nil, err
	}
	var err error
	if setting.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt)); err != nil {
		log.Println(err)
		return nil, err
	}

	return &setting, nil
}

func (r *gameSettingsRepository) CacheActiveGameSettings() error {
	settings, err := r.GetActiveGameSettings()
	if err != nil {
		log.Println(err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	settingsJson, err := json.Marshal(settings)
	if err != nil {
		log.Println(err)
		return err
	}

	// 古いキャッシュを削除
	if err := r.rdb.Del(ctx, "game_settings").Err(); err != nil {
		log.Println(err)
		return err
	}

	if err := r.rdb.Set(ctx, "game_settings", settingsJson, 0).Err(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (r *gameSettingsRepository) GetActiveGameSettingsFromCache() (*entities.GameSettings, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	settingsJon, err := r.rdb.Get(ctx, "game_settings").Result()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var settings entities.GameSettings
	if err := json.Unmarshal([]byte(settingsJon), &settings); err != nil {
		log.Println(err)
		return nil, err
	}

	return &settings, nil
}

func (r *gameSettingsRepository) ChangeActiveGameSettings(ID entities.GameSettingID) error {
	query := "UPDATE game_settings SET is_active = false"
	if _, err := r.db.Exec(query); err != nil {
		log.Println(err)
		return err
	}

	query = "UPDATE game_settings SET is_active = true WHERE id = ?"
	if _, err := r.db.Exec(query, ID); err != nil {
		log.Println(err)
		return err
	}

	if err := r.CacheActiveGameSettings(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (r *gameSettingsRepository) AddGameSettings(settings entities.GameSettings) error {
	query := "INSERT INTO game_settings (gacha_coin_consumption, ranking_list_limit, n_weight, s_weight, sr_weight, max_gacha_times) VALUES (?, ?, ?, ?, ?)"
	if _, err := r.db.Exec(query, settings.GachaCoinConsumption, settings.RankingListLimit, settings.NWeight, settings.RWeight, settings.SrWeight, settings.MaxGachaTimes); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

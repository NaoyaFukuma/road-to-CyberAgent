package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"42tokyo-road-to-dojo-go/pkg/server/entities"

	"github.com/go-redis/redis/v8"
)

type ItemRepository interface {
	GetItems() (*entities.Items, error)
	CacheItems() error
	GetItemsFromCache() (*entities.Items, error)
	AddItem(item *entities.Item) error
	DeleteItemByID(ID entities.ItemID) error
}

func NewItemRepository(db *sql.DB, rdb *redis.Client) ItemRepository {
	return &itemRepository{db, rdb}
}

type itemRepository struct {
	db  *sql.DB
	rdb *redis.Client
}

func (r *itemRepository) GetItems() (*entities.Items, error) {
	query := "SELECT id, name, rarity FROM item"
	rows, err := r.db.Query(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
		if err := rows.Close(); err != nil {
			log.Println(err)
		}
	}()

	var items entities.Items
	for rows.Next() {
		item := new(entities.Item)
		if err := rows.Scan(&item.ID, &item.Name, &item.Rarity); err != nil {
			log.Println(err)
			return nil, err
		}
		items = append(items, *item)
	}

	return &items, nil
}

func (r *itemRepository) CacheItems() error {
	items, err := r.GetItems()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	for _, item := range *items {
		itemJson, err := json.Marshal(item)
		if err != nil {
			log.Println(err)
			return err
		}
		if err := r.rdb.Set(ctx, fmt.Sprintf("item:%d", item.ID), itemJson, 0).Err(); err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func (r *itemRepository) GetItemsFromCache() (*entities.Items, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	keys, err := r.rdb.Keys(ctx, "item:*").Result()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var items entities.Items
	for _, key := range keys {
		itemJson, err := r.rdb.Get(ctx, key).Result()
		if err != nil {
			log.Println(err)
			return nil, err
		}
		var item entities.Item
		if err := json.Unmarshal([]byte(itemJson), &item); err != nil {
			log.Println(err)
			return nil, err
		}
		items = append(items, item)
	}

	return &items, nil
}

func (r *itemRepository) AddItem(item *entities.Item) error {
	query := "INSERT INTO item (name, rarity) VALUES (?, ?)"
	result, err := r.db.Exec(query, item.Name, item.Rarity)
	if err != nil {
		log.Println(err)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return err
	}

	item.ID = entities.ItemID(id)

	// redisにも追加
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	itemJson, err := json.Marshal(item)
	if err != nil {
		log.Println(err)
		return err
	}
	if err := r.rdb.Set(ctx, fmt.Sprint(int64(item.ID)), itemJson, 0).Err(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (r *itemRepository) DeleteItemByID(ID entities.ItemID) error {
	query := "DELETE FROM item WHERE id = ?"
	_, err := r.db.Exec(query, ID)
	if err != nil {
		log.Println(err)
		return err
	}

	// redisからも削除
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := r.rdb.Del(ctx, fmt.Sprintf("%d", int64(ID))).Err(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

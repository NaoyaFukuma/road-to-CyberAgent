package repositories

import (
	"database/sql"
	"log"

	"42tokyo-road-to-dojo-go/pkg/server/entities"
)

type CollectionItemRepository interface {
	GetCollectionItems(userID entities.UserID) (*[]entities.ItemID, error)
	AddCollectionItem(userID entities.UserID, itemID entities.ItemID) error
	AddCollectionItemTransaction(tx *sql.Tx, userID entities.UserID, itemID entities.ItemID) error
}

func NewCollectionItemRepository(db *sql.DB) CollectionItemRepository {
	return &collectionItemRepository{db}
}

type collectionItemRepository struct {
	db *sql.DB
}

// 自身の所持するアイテムのIDを取得する。
func (r *collectionItemRepository) GetCollectionItems(userID entities.UserID) (*[]entities.ItemID, error) {
	query := "SELECT item_id FROM user_items WHERE user_id = ?"

	rows, err := r.db.Query(query, userID)
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

	var itemIDs []entities.ItemID
	for rows.Next() {
		var itemID entities.ItemID
		if err := rows.Scan(&itemID); err != nil {
			log.Println(err)
			return nil, err
		}
		itemIDs = append(itemIDs, itemID)
	}

	return &itemIDs, nil
}

// 所持アイテムを追加する。重複追加は発生しない。
func (r *collectionItemRepository) AddCollectionItem(userID entities.UserID, itemID entities.ItemID) error {
	query := "INSERT INTO user_items (user_id, item_id) VALUES (?, ?)"

	_, err := r.db.Exec(query, userID, itemID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// 所持アイテムを追加する。重複追加は発生しない。
func (r *collectionItemRepository) AddCollectionItemTransaction(tx *sql.Tx, userID entities.UserID, itemID entities.ItemID) error {
	query := "INSERT INTO user_items (user_id, item_id) VALUES (?, ?)"

	_, err := execQueryAndReturnAffectedRows(tx, query, userID, itemID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

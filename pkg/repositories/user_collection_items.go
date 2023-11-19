package repositories

import (
	"database/sql"
	"log"

	"42tokyo-road-to-dojo-go/pkg/server/entities"
)

type CollectionItemRepository interface {
	GetCollectionItems(userID entities.UserID) (*[]entities.ItemID, error)
	AddCollectionItems(userID entities.UserID, itemIDs []entities.ItemID) error
	AddCollectionItemsTransaction(tx *sql.Tx, userID entities.UserID, itemIDs []entities.ItemID) error
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
func (r *collectionItemRepository) AddCollectionItems(userID entities.UserID, itemIDs []entities.ItemID) error {
	if len(itemIDs) == 0 {
		return nil // 追加するアイテムがない場合は、何もしない
	}

	// クエリのベースを作成
	query := "INSERT INTO user_items (user_id, item_id) VALUES "

	// パラメータのスライスを作成
	var params []interface{}

	// 各アイテムIDについて、クエリを構築
	for _, itemID := range itemIDs {
		query += "(?, ?),"
		params = append(params, userID, itemID)
	}

	// 最後のカンマを削除
	query = query[:len(query)-1]

	// クエリを実行
	_, err := r.db.Exec(query, params...)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// 複数の所持アイテムを追加する。重複追加は発生しない。
func (r *collectionItemRepository) AddCollectionItemsTransaction(tx *sql.Tx, userID entities.UserID, itemIDs []entities.ItemID) error {
	if len(itemIDs) == 0 {
		return nil // 追加するアイテムがない場合は、何もしない
	}

	// クエリのベースを作成
	query := "INSERT INTO user_items (user_id, item_id) VALUES "

	// パラメータのスライスを作成
	var params []interface{}

	// 各アイテムIDについて、クエリを構築
	for _, itemID := range itemIDs {
		query += "(?, ?),"
		params = append(params, userID, itemID)
	}

	// 最後のカンマを削除
	query = query[:len(query)-1]

	// トランザクション内でクエリを実行
	_, err := tx.Exec(query, params...)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

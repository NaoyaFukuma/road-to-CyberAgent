package handler

import (
	"log"
	"net/http"

	"42tokyo-road-to-dojo-go/pkg/http/response"
	"42tokyo-road-to-dojo-go/pkg/repositories"
	"42tokyo-road-to-dojo-go/pkg/server/entities"
)

// ユーザの所持アイテムを取得
// キャッシュからItemを取得し、DBから自身の所持アイテムのIDを取得する。
// その後、所持アイテムにはHasItemをtrueに設定する。
// このとき、キャッシュからの取得に失敗した場合はDBから取得する。
func HandleGetCollectionList(repos *repositories.Repositories) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		itemRepo := repos.ItemRepository
		collectionItemRepo := repos.CollectionItemRepository
		// ContextからUserIDを取得
		ctx := request.Context()
		userID, ok := ctx.Value("userID").(entities.UserID)
		if !ok {
			log.Println("userID is not found")
			response.SetStatusAndJson(writer, http.StatusInternalServerError, map[string]string{"error": "userID is not found"})
			return
		}
		// キャッシュからアイテムを取得
		var items *entities.Items
		items, err := itemRepo.GetItemsFromCache()
		if err != nil {
			log.Println(err)
			// キャッシュからの取得に失敗した場合はDBから取得する
			items, err = itemRepo.GetItems()
			if err != nil {
				log.Println(err)
				response.SetStatusAndJson(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}
		}

		// 自身の所持アイテムのIDを取得
		collectItemIDs, err := collectionItemRepo.GetCollectionItems(userID)
		if err != nil {
			log.Println(err)
			response.SetStatusAndJson(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		var collectionItems []entities.CollectionItem
		// HashMapによる高速化のため、IDｓのスライスをmapに変換する
		collectItemIDsMap := make(map[entities.ItemID]bool)
		for _, id := range *collectItemIDs {
			collectItemIDsMap[id] = true
		}

		// 所持アイテムはHasItemをtrueに設定する
		for _, item := range *items {
			hasItem := entities.HasItem(collectItemIDsMap[item.ID])
			collectionItems = append(collectionItems, entities.CollectionItem{
				ID:      item.ID,
				Name:    item.Name,
				Rarity:  item.Rarity,
				HasItem: hasItem,
			})
		}

		var collectionItemList entities.CollectionItemList
		collectionItemList.Items = collectionItems

		// レスポンスヘッダにステータスコードを設定
		response.SetStatusAndJson(writer, http.StatusOK, collectionItemList)
	}
}

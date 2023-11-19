package handler

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"

	"42tokyo-road-to-dojo-go/pkg/http/response"
	"42tokyo-road-to-dojo-go/pkg/repositories"
	"42tokyo-road-to-dojo-go/pkg/server/entities"
)

// ガチャを引く
// 引く回数をJSONで"times": 10のように指定
// ctxからユーザーIDを取得し、ユーザーの所持コイン数を確認
// トランザクションの期間を短くするために、先にガチャの結果を計算
// Response用にガチャの結果を整形し
// トランザクション開始し、
// 所持コインを引く処理を行う、
// 所持アイテムに加える処理を行う、
// その際に、すでに持っているアイテムの場合は、Repositoryを呼び出さない。新規の場合のみ呼び出し、isNewをtrueにする
func HandleGachaDraw(repos *repositories.Repositories) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		userID := request.Context().Value("userID").(entities.UserID)

		type times struct {
			Times int `json:"times"`
		}
		var timesJSON times
		err := json.NewDecoder(request.Body).Decode(&timesJSON)
		if err != nil {
			log.Println(err)
			response.SetStatusAndJson(writer, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		// validation
		timesInt, err := int64(timesJSON.Times), nil
		if err != nil {
			log.Println(err)
			response.SetStatusAndJson(writer, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		if timesInt < 1 {
			response.SetStatusAndJson(writer, http.StatusBadRequest, map[string]string{"error": "times must be greater than 0"})
			return
		}

		// game_settingsテーブルからガチャの設定を取得
		gameSettingsRepo := repos.GameSettingsRepository
		gameSettings, err := gameSettingsRepo.GetActiveGameSettingsFromCache()
		if err != nil {
			log.Println(err)
			response.SetStatusAndJson(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		if timesInt > int64(gameSettings.MaxGachaTimes) {
			response.SetStatusAndJson(writer, http.StatusBadRequest, map[string]string{"error": "times must be less than max_gacha_times"})
			return
		}

		// IDからユーザー情報を取得
		userRepo := repos.UserRepository
		user, err := userRepo.GetUserByID(userID)
		if err != nil {
			log.Println(err)
			response.SetStatusAndJson(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		// 所持コイン < ガチャのコスト * 引く回数 の場合はエラー
		if user.Coin < entities.Coin(gameSettings.GachaCoinConsumption)*entities.Coin(timesInt) {
			log.Println("not enough coin")
			response.SetStatusAndJson(writer, http.StatusBadRequest, map[string]string{"error": "not enough coin"})
			return
		}

		// ガチャの結果を計算するために、重み付けの合計を計算
		// settingのNWeight, RWeight, SrWeightを各itemのrarityに応じて採用して足し合わせる
		// まずはitemRepositoryからitemの一覧を取得
		itemRepo := repos.ItemRepository
		items, err := itemRepo.GetItemsFromCache()
		if err != nil {
			log.Println(err)
			response.SetStatusAndJson(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		itemsWithWeight := make([]entities.ItemWithWeight, 0, len(*items))

		// 重み付けの合計を計算
		var totalWeight entities.Weight = 0
		for _, item := range *items {
			itemsWithWeight = append(itemsWithWeight, entities.ItemWithWeight{
				Item:   item,
				Weight: 0,
			})

			switch item.Rarity {
			case entities.N:
				itemsWithWeight[len(itemsWithWeight)-1].Weight = gameSettings.NWeight
				totalWeight += gameSettings.NWeight
			case entities.R:
				itemsWithWeight[len(itemsWithWeight)-1].Weight = gameSettings.RWeight
				totalWeight += gameSettings.RWeight
			case entities.SR:
				itemsWithWeight[len(itemsWithWeight)-1].Weight = gameSettings.SrWeight
				totalWeight += gameSettings.SrWeight
			}
		}

		// ガチャの結果であるアイテムのIDを保存するための変数
		var gachaGetIDs []entities.ItemID

		// times回ガチャを引く
		for i := int64(0); i < timesInt; i++ {
			// 乱数を生成
			randomNum := rand.Int63n(int64(totalWeight))

			// 重み付けの合計を超えるまで、重み付けを引いていく
			var weightSum int64 = 0
			for _, itemWithWeight := range itemsWithWeight {
				weightSum += int64(itemWithWeight.Weight)
				if randomNum < weightSum {
					gachaGetIDs = append(gachaGetIDs, itemWithWeight.ID)
					break
				}
			}
		}

		// トランザクションの開始
		tx, err := repos.DB.Begin()
		if err != nil {
			log.Println(err)
			response.SetStatusAndJson(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		// 所持コインを引く処理
		err = userRepo.UpdateUserCoinsByIDTransaction(tx, userID, user.Coin-entities.Coin(gameSettings.GachaCoinConsumption)*entities.Coin(timesInt))
		if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Println(err)
				response.SetStatusAndJson(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}
			log.Println(err)
			response.SetStatusAndJson(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		// 既存のコレクションアイテムを取得
		collectionItemRepo := repos.CollectionItemRepository
		collectionItems, err := collectionItemRepo.GetCollectionItems(userID)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Println(err)
				response.SetStatusAndJson(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}
			log.Println(err)
			response.SetStatusAndJson(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		// GachResultListを作成
		var gachaResultList entities.GachaResultList

		// ガチャの結果をループで回しながら、isNewの判定と設定、newはコレクションアイテムに追加する
		// ハッシュマップを使って、O(n)で済むようにする
		var itemsMap = make(map[entities.ItemID]entities.Item, len(*items))
		for _, item := range *items {
			itemsMap[item.ID] = item
		}
		var collectionItemMap = make(map[entities.ItemID]bool, len(*collectionItems))
		for _, collectionItem := range *collectionItems {
			collectionItemMap[collectionItem] = true
		}

		// isNewのアイテムIDを格納するためのスライス
		var isNewItemIDs []entities.ItemID
		for _, gachaGetID := range gachaGetIDs {
			if _, ok := collectionItemMap[gachaGetID]; ok {
				gachaResultList.Items = append(gachaResultList.Items, entities.GachaResult{
					ID:     gachaGetID,
					Name:   itemsMap[gachaGetID].Name,
					Rarity: itemsMap[gachaGetID].Rarity,
					IsNew:  false,
				})
			} else {
				isNewItemIDs = append(isNewItemIDs, gachaGetID)
				collectionItemMap[gachaGetID] = true
				gachaResultList.Items = append(gachaResultList.Items, entities.GachaResult{
					ID:     gachaGetID,
					Name:   itemsMap[gachaGetID].Name,
					Rarity: itemsMap[gachaGetID].Rarity,
					IsNew:  true,
				})
			}
		}

		// 所持アイテムに加える処理
		err = collectionItemRepo.AddCollectionItemsTransaction(tx, userID, isNewItemIDs)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Println(err)
				response.SetStatusAndJson(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}
			log.Println(err)
			response.SetStatusAndJson(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		// commit
		if err := tx.Commit(); err != nil {
			if err := tx.Rollback(); err != nil {
				log.Println(err)
				response.SetStatusAndJson(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}
			log.Println(err)
			response.SetStatusAndJson(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		// レスポンスを返却
		response.SetStatusAndJson(writer, http.StatusOK, gachaResultList)
	}
}

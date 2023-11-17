package handler

import (
	"log"
	"net/http"
	"strconv"

	"42tokyo-road-to-dojo-go/pkg/http/response"
	"42tokyo-road-to-dojo-go/pkg/repositories"
	"42tokyo-road-to-dojo-go/pkg/server/entities"
)

// ランクキングリスト取得処理
// 指定した順位から一定数の順位までのランキング情報を取得します。
// 例えば「サーバ側での1回あたりのランキング取得件数設定」が10で、PathQuery「startパラメータ」の指定が1だった場合は1位〜10位を、「startパラメータ」の指定が5だった場合は5位〜14位を返却します。
// 本課題では同率順位は考慮せず、同じスコアだった場合はユーザーIDの昇順で順位を決定するものとします。
func HandleGetRankingList(repos *repositories.Repositories) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		userScoresRepo := repos.UserScoresRepository
		// 入力の受け取り
		start := request.URL.Query().Get("start")
		if start == "" {
			response.SetStatusAndJson(writer, http.StatusBadRequest, map[string]string{"error": "start is required"})
			return
		}

		// validation
		start64, err := strconv.ParseInt(start, 10, 64)
		if err != nil {
			log.Println(err)
			response.SetStatusAndJson(writer, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if start64 < 1 {
			response.SetStatusAndJson(writer, http.StatusBadRequest, map[string]string{"error": "start must be greater than 0"})
			return
		}

		// GameSettingsRepositoryからサーバ側での1回あたりのランキング取得件数設定を取得
		gameSettings, err := repos.GameSettingsRepository.GetActiveGameSettingsFromCache()
		if err != nil {
			log.Println(err)
			response.SetStatusAndJson(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		// ランキング情報をリポジトリから取得
		offset := start64 - 1
		scores, err := userScoresRepo.GetUserScoreWithUserName(offset, gameSettings.RankingListLimit)
		if err != nil {
			log.Println(err)
			response.SetStatusAndJson(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		// 順位を付与する
		rankings := make([]entities.RankInfo, 0, len(*scores))
		for i, score := range *scores {
			var rankInfo entities.RankInfo = entities.RankInfo{
				UserID:   score.UserID,
				UserName: score.UserName,
				Rank:     entities.Rank(offset + int64(i) + 1),
				Score:    score.Score,
			}
			rankings = append(rankings, rankInfo)
		}

		var rankingList entities.RankingListResponse = entities.RankingListResponse{
			RankInfoList: rankings,
		}
		response.SetStatusAndJson(writer, http.StatusOK, rankingList)
	}
}

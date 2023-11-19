package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"42tokyo-road-to-dojo-go/pkg/http/response"
	"42tokyo-road-to-dojo-go/pkg/repositories"
	"42tokyo-road-to-dojo-go/pkg/server/entities"
)

func HandleGameFinish(repos *repositories.Repositories) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// 入力の受け取りとvalidation
		userID := request.Context().Value("userID").(entities.UserID)

		type score = struct {
			Score int `json:"score"`
		}
		var scoreJSON score
		err := json.NewDecoder(request.Body).Decode(&scoreJSON)
		if err != nil {
			log.Println(err)
			response.SetStatusAndJson(writer, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		// validation
		scoreInt, err := int64(scoreJSON.Score), nil
		if err != nil {
			log.Println(err)
			response.SetStatusAndJson(writer, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if scoreInt < 0 {
			response.SetStatusAndJson(writer, http.StatusBadRequest, map[string]string{"error": "score must be greater than 0"})
			return
		}

		// DB処理にあたり、トランザクションを開始
		tx, err := repos.DB.Begin()
		if err != nil {
			log.Println(err)
			response.SetStatusAndJson(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		// user_scoresテーブルにスコアを登録
		userScoresRepo := repos.UserScoresRepository
		err = userScoresRepo.AddUserScoreTransaction(tx, userID, entities.Score(scoreInt))
		if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Println(err)
				response.SetStatusAndJson(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}
			log.Println(err)
			response.SetStatusAndJson(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		// user.highScoreと比較して、今回のスコアの方が高ければuser.highScoreを更新
		userRepo := repos.UserRepository
		user, err := userRepo.GetUserByID(userID)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Println(err)
				response.SetStatusAndJson(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}
			log.Println(err)
			response.SetStatusAndJson(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		if user.HighScore < entities.Score(scoreInt) {
			err = userRepo.UpdateUserHighScoreByIDTransaction(tx, userID, entities.Score(scoreInt))
			if err != nil {
				if err := tx.Rollback(); err != nil {
					log.Println(err)
					response.SetStatusAndJson(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				}
				log.Println(err)
				response.SetStatusAndJson(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}
		}

		//scoreからコインの数を計算し、user.coinsに加算し、Responseに増加したコインの数を返却
		coin := entities.Coin(scoreInt / 10)
		err = userRepo.UpdateUserCoinsByID(userID, user.Coin+coin)
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

		// レスポンスを返却 coinをJSONに変換
		type coinResponse struct {
			Coin string `json:"coin"`
		}
		coinResponseJSON := coinResponse{
			Coin: strconv.FormatInt(int64(coin), 10),
		}
		response.SetStatusAndJson(writer, http.StatusOK, coinResponseJSON)
	}
}

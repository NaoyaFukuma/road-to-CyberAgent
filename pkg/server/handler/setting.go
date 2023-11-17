package handler

import (
	"log"
	"net/http"
	"strconv"

	"42tokyo-road-to-dojo-go/pkg/http/response"
	"42tokyo-road-to-dojo-go/pkg/repositories"
	"42tokyo-road-to-dojo-go/pkg/server/entities"
)

// HandleSettingGet ゲーム設定情報取得処理
func HandleSettingGet(repos *repositories.Repositories) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		settingRepo := repos.GameSettingsRepository
		// ゲーム設定情報をリポジトリから取得
		settings, err := settingRepo.GetActiveGameSettingsFromCache()
		if err != nil {
			log.Println(err)
			response.SetStatusAndJson(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		response.SetStatusAndJson(writer, http.StatusOK, settings)
	}
}

// settingRepo.GetAllGameSettings()を使うハンドラ
func HandleAllSettingGet(repos *repositories.Repositories) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		settingRepo := repos.GameSettingsRepository
		// ゲーム設定情報をリポジトリから取得
		settings, err := settingRepo.GetAllGameSettings()
		if err != nil {
			log.Println(err)
			response.SetStatusAndJson(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		response.SetStatusAndJson(writer, http.StatusOK, settings)
	}
}

// GetGameSettingsByID を使うハンドラ id はクエリパラメータ
func HandleGetGameSettingsByID(repos *repositories.Repositories) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		settingRepo := repos.GameSettingsRepository
		// 入力の受け取り
		id := request.URL.Query().Get("id")
		if id == "" {
			response.SetStatusAndJson(writer, http.StatusBadRequest, map[string]string{"error": "id is required"})
			return
		}

		// validation
		id64, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			log.Println(err)
			response.SetStatusAndJson(writer, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if id64 < 1 {
			response.SetStatusAndJson(writer, http.StatusBadRequest, map[string]string{"error": "id must be positive"})
			return
		}

		// データ層の呼び出し
		settings, err := settingRepo.GetGameSettingsByID(entities.GameSettingID(id64))
		if err != nil {
			log.Println(err)
			response.SetStatusAndJson(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		// レスポンスの作成
		response.SetStatusAndJson(writer, http.StatusOK, settings)
	}
}

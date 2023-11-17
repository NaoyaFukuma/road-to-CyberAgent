package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"

	"42tokyo-road-to-dojo-go/pkg/http/response"
	"42tokyo-road-to-dojo-go/pkg/repositories"
	"42tokyo-road-to-dojo-go/pkg/server/entities"
)

// ユーザ作成
func HandleUserCreate(repos *repositories.Repositories) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		userRepo := repos.UserRepository
		var user entities.User
		err := json.NewDecoder(request.Body).Decode(&user)
		if err != nil {
			log.Println(err)
			response.SetStatusAndJson(writer, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		// validation
		err = validateUser(&user)
		if err != nil {
			log.Println(err)
			response.SetStatusAndJson(writer, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		// トークンを生成
		user.AuthToken = entities.AuthToken(uuid.New().String())

		// ユーザを作成
		err = userRepo.CreateUser(&user)
		if err != nil {
			log.Println(err)
			response.SetStatusAndJson(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		// レスポンスヘッダにステータスコードを設定
		response.SetStatusAndJson(writer, http.StatusOK, map[string]string{"token": string(user.AuthToken)})
	}
}

func validateUser(user *entities.User) error {
	if user.Name == "" {
		return fmt.Errorf("name is required")
	}
	// TODO:下限も考える。上限もDBは128文字だが、実際にはどうなのか
	if len(user.Name) > 128 {
		return fmt.Errorf("name is too long")
	}
	return nil
}

// ユーザ取得 ContextからユーザIDを取得してユーザを取得する
func HandleUserGet(repos *repositories.Repositories) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		userRepo := repos.UserRepository
		ctx := request.Context()
		userID, ok := ctx.Value("userID").(entities.UserID)
		if !ok {
			log.Println("userID is not found")
			response.SetStatusAndJson(writer, http.StatusInternalServerError, map[string]string{"error": "userID is not found"})
			return
		}

		// ユーザを取得
		user, err := userRepo.GetUserByID(userID)
		if err != nil {
			log.Println(err)
			response.SetStatusAndJson(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		// レスポンスヘッダにステータスコードを設定
		response.SetStatusAndJson(writer, http.StatusOK, user)
	}
}

// ユーザ更新 ContextからユーザIDを取得してユーザを更新する
func HandleUserUpdate(repos *repositories.Repositories) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		userRepo := repos.UserRepository
		ctx := request.Context()
		userID, ok := ctx.Value("userID").(entities.UserID)
		if !ok {
			log.Println("userID is not found")
			response.SetStatusAndJson(writer, http.StatusInternalServerError, map[string]string{"error": "userID is not found"})
			return
		}

		var user entities.User
		err := json.NewDecoder(request.Body).Decode(&user)
		if err != nil {
			log.Println(err)
			response.SetStatusAndJson(writer, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		// validation
		err = validateUser(&user)
		if err != nil {
			log.Println(err)
			response.SetStatusAndJson(writer, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		// ユーザを更新
		user.ID = userID
		err = userRepo.UpdateUserNameByID(userID, user.Name)
		if err != nil {
			log.Println(err)
			response.SetStatusAndJson(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		// レスポンスヘッダにステータスコードを設定
		writer.WriteHeader(http.StatusOK)
	}
}

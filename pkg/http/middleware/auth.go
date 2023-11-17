package middleware

import (
	"42tokyo-road-to-dojo-go/pkg/http/response"
	"42tokyo-road-to-dojo-go/pkg/repositories"
	"42tokyo-road-to-dojo-go/pkg/server/entities"
	"context"
	"net/http"
)

// Authenticate ユーザ認証を行ってContextへユーザID情報を保存する
func Authenticate(repos *repositories.Repositories, nextFunc http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		userRepo := repos.UserRepository
		ctx := request.Context()
		if ctx == nil {
			ctx = context.Background()
		}

		// x-tokenヘッダからトークンを取得
		var token string
		if token = request.Header.Get("x-token"); token == "" {
			response.SetStatusAndJson(writer, http.StatusUnauthorized, map[string]string{"error": "x-token header is required"})
			return
		}

		authToken := entities.AuthToken(token)
		userID, err := userRepo.GetUserIDAuthToken(authToken)
		if err != nil {
			response.SetStatusAndJson(writer, http.StatusUnauthorized, map[string]string{"error": err.Error()})
			return
		}

		// ユーザIDをContextへ保存
		ctx = context.WithValue(ctx, "userID", userID)

		// 次のハンドラを実行
		nextFunc(writer, request.WithContext(ctx))
	}
}

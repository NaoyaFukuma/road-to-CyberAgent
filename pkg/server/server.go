package server

import (
	"log"
	"net/http"

	"42tokyo-road-to-dojo-go/pkg/http/middleware"
	"42tokyo-road-to-dojo-go/pkg/repositories"
	"42tokyo-road-to-dojo-go/pkg/server/handler"
)

func Serve(addr string, repos *repositories.Repositories) {
	/* ルーティング設定 */
	// ゲーム設定関連
	http.HandleFunc("/setting/get", get(handler.HandleSettingGet(repos)))
	http.HandleFunc("/setting/all", get(handler.HandleAllSettingGet(repos)))
	http.HandleFunc("/setting/get_by_id", get(handler.HandleGetGameSettingsByID(repos)))

	// ユーザ関連
	http.HandleFunc("/user/create", post(handler.HandleUserCreate(repos)))
	http.HandleFunc("/user/get", get(middleware.Authenticate(repos, handler.HandleUserGet(repos))))
	http.HandleFunc("/user/update", post(middleware.Authenticate(repos, handler.HandleUserUpdate(repos))))

	// 所持アイテム関連
	http.HandleFunc("/collection/list", get(middleware.Authenticate(repos, handler.HandleGetCollectionList(repos))))

	// ランキング関連
	http.HandleFunc("/ranking/list", get(middleware.Authenticate(repos, handler.HandleGetRankingList(repos))))

	// ゲーム関連
	http.HandleFunc("/game/finish", post(middleware.Authenticate(repos, handler.HandleGameFinish(repos))))

	// ガチャ関連
	http.HandleFunc("/gacha/draw", post(middleware.Authenticate(repos, handler.HandleGachaDraw(repos))))

	/* ===== サーバの起動 ===== */
	log.Println("Server running...")
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Listen and serve failed. %+v", err)
	}
}

// get GETリクエストを処理する
func get(apiFunc http.HandlerFunc) http.HandlerFunc {
	return httpMethod(apiFunc, http.MethodGet)
}

// post POSTリクエストを処理する
func post(apiFunc http.HandlerFunc) http.HandlerFunc {
	return httpMethod(apiFunc, http.MethodPost)
}

// httpMethod 指定したHTTPメソッドでAPIの処理を実行する
func httpMethod(apiFunc http.HandlerFunc, method string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		// CORS対応
		writer.Header().Add("Access-Control-Allow-Origin", "*")
		writer.Header().Add("Access-Control-Allow-Headers", "Content-Type,Accept,Origin,x-token")

		// プリフライトリクエストは処理を通さない
		if request.Method == http.MethodOptions {
			return
		}
		// 指定のHTTPメソッドでない場合はエラー
		if request.Method != method {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			writer.Write([]byte("Method Not Allowed"))
			return
		}

		// 共通のレスポンスヘッダを設定
		writer.Header().Add("Content-Type", "application/json")
		apiFunc(writer, request)
	}
}

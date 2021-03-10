package server

import (
	"20dojo-online/pkg/db"
	"20dojo-online/pkg/http/middleware"
	"20dojo-online/pkg/http/response"
	"20dojo-online/pkg/server/service"
	"log"
	"math/rand"
	"net/http"
	"time"

	"20dojo-online/pkg/server/handler"
	"20dojo-online/pkg/server/model"
)

var (
	httpResponse = response.NewHttpResponse()

	userRepository = model.NewUserRepository(db.Conn)
	authMiddleware = middleware.NewMiddleware(httpResponse, userRepository)

	gachaProbabilityRepository   = model.NewGachaRepositoryRepository(db.Conn)
	userCollectionItemRepository = model.NewUserCollectionItemRepository(db.Conn)
	collectionItemRepository     = model.NewCollectionItemRepository(db.Conn)

	gameService       = service.NewGameService(userRepository)
	gachaService      = service.NewGachaService(userRepository, gachaProbabilityRepository, userCollectionItemRepository, collectionItemRepository)
	rankingService    = service.NewRankingService(userRepository)
	collectionService = service.NewCollectionService(userCollectionItemRepository, collectionItemRepository)

	userHandler       = handler.NewUserHandler(httpResponse, userRepository)
	settingHandler    = handler.NewSettingHandler(httpResponse)
	gameHandler       = handler.NewGameHandler(httpResponse, gameService)
	gachaHandler      = handler.NewGachaHandler(httpResponse, gachaService)
	rankingHandler    = handler.NewRankingHandler(httpResponse, rankingService)
	collectionHandler = handler.NewCollectionHandler(httpResponse, collectionService)
)

// Serve HTTPサーバを起動する
func Serve(addr string) {

	rand.Seed(time.Now().UnixNano())

	/* ===== URLマッピングを行う ===== */
	http.HandleFunc("/setting/get", get(middleware.AccessLogging(settingHandler.HandleSettingGet)))
	http.HandleFunc("/user/create", post(middleware.AccessLogging(userHandler.HandleUserCreate)))
	http.HandleFunc("/user/get",
		get(authMiddleware.Authenticate(middleware.AccessLogging(userHandler.HandleUserGet))))
	http.HandleFunc("/user/update",
		post(authMiddleware.Authenticate(middleware.AccessLogging(userHandler.HandleUserUpdate))))

	http.HandleFunc("/game/finish", post(middleware.AccessLogging(authMiddleware.Authenticate(gameHandler.HandleGameFinish))))

	http.HandleFunc("/gacha/draw", post(middleware.AccessLogging(authMiddleware.Authenticate(gachaHandler.HandleGachaDraw))))

	http.HandleFunc("/ranking/list", get(middleware.AccessLogging(authMiddleware.Authenticate(rankingHandler.HandleRankingList))))

	http.HandleFunc("/collection/list", get(middleware.AccessLogging(authMiddleware.Authenticate(collectionHandler.HandleUserCollectionList))))

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

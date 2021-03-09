package handler

import (
	"20dojo-online/pkg/dcontext"
	"20dojo-online/pkg/http/response"
	"20dojo-online/pkg/myerror"
	"20dojo-online/pkg/server/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type gameFinishRequest struct {
	Score int `json:"score"`
}

type gameFinishResponse struct {
	Coin int `json:"coin"`
}

type GameHandler struct {
	HttpResponse response.HttpResponseInterface
	GameService  service.GameServiceInterface
}

func NewGameHandler(httpResponse response.HttpResponseInterface, gameService service.GameServiceInterface) *GameHandler {
	return &GameHandler{
		HttpResponse: httpResponse,
		GameService:  gameService,
	}
}

// HandleGameFinish インゲーム終了
func (h *GameHandler) HandleGameFinish(writer http.ResponseWriter, request *http.Request) {

	// リクエストbodyからスコアを取得
	var requestBody gameFinishRequest
	if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
		err = myerror.ApplicationError{
			Message:       "failed to decode request body",
			OriginalError: err,
			Code:          http.StatusBadRequest,
		}
		log.Println(err)
		h.HttpResponse.Failed(writer, err)
		return
	}
	// scoreが負の数のときエラーを返す
	if requestBody.Score < 0 {
		err := myerror.ApplicationError{
			Message: fmt.Sprintf("score is minus. score=%d", requestBody.Score),
			Code:    http.StatusBadRequest,
		}
		log.Println(err)
		h.HttpResponse.Failed(writer, err)
		return
	}

	// ミドルウェアでコンテキストに格納したユーザidの取得
	ctx := request.Context()
	userID := dcontext.GetUserIDFromContext(ctx)
	if userID == "" {
		userIDEmptyErr := myerror.ApplicationError{
			Message: "userID from context is empty",
			Code:    http.StatusInternalServerError,
		}
		log.Println(userIDEmptyErr)
		h.HttpResponse.Failed(writer, userIDEmptyErr)
		return
	}

	// ゲーム終了時のロジック
	res, err := h.GameService.FinishGame(&service.FinishGameRequest{
		UserId: userID,
		Score:  requestBody.Score,
	})
	if err != nil {
		err = myerror.ApplicationError{
			Message: "failed to finish game correctly",
			Code:    http.StatusInternalServerError,
		}
		log.Println(err)
		h.HttpResponse.Failed(writer, err)
		return
	}

	// 獲得コインをレスポンスとして返す
	h.HttpResponse.Success(writer, &gameFinishResponse{
		Coin: res.Coin,
	})
}

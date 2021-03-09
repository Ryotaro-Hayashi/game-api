package handler

import (
	"20dojo-online/pkg/dcontext"
	"20dojo-online/pkg/http/response"
	"20dojo-online/pkg/myerror"
	"20dojo-online/pkg/server/service"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type gachaDrawRequest struct {
	Times int `json:"times"`
}

type gachaDrawResponse struct {
	Results []*result `json:"results"`
}

type result struct {
	CollectionID string `json:"collectionID"`
	Name         string `json:"name"`
	Rarity       int    `json:"rarity"`
	IsNew        bool   `json:"isNew"`
}

type GachaHandler struct {
	HttpResponse response.HttpResponseInterface
	GachaService service.GachaServiceInterface
}

func NewGachaHandler(httpResponse response.HttpResponseInterface, gachaService service.GachaServiceInterface) *GachaHandler {
	return &GachaHandler{
		HttpResponse: httpResponse,
		GachaService: gachaService,
	}
}

func (h *GachaHandler) HandleGachaDraw(writer http.ResponseWriter, request *http.Request) {

	// リクエストbodyからガチャ実行回数を取得
	var requestBody gachaDrawRequest
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

	// timesが0以下のときエラーを返す
	if requestBody.Times <= 0 {
		timesLessErr := myerror.ApplicationError{
			Message: fmt.Sprintf("gacha draw times is 0 or less. times=%d", requestBody.Times),
			Code:    http.StatusBadRequest,
		}
		log.Println(timesLessErr)
		h.HttpResponse.Failed(writer, timesLessErr)
		return
	}

	// Contextから認証済みのユーザIDを取得
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

	// ガチャ実行ロジック
	res, err := h.GachaService.DrawGacha(&service.DrawGachaRequest{
		Times:  requestBody.Times,
		UserID: userID,
	})
	if err != nil {
		var appErr myerror.ApplicationError
		if errors.As(err, &appErr) {
			log.Println(err)
			h.HttpResponse.Failed(writer, err)
			return
		} else {
			err = myerror.ApplicationError{
				Message:       "failed to draw gacha correctly",
				OriginalError: err,
				Code:          http.StatusInternalServerError,
			}
			log.Println(err)
			h.HttpResponse.Failed(writer, err)
			return
		}
	}

	// レスポンス
	var results []*result
	for _, gachaResult := range res.GachaResults {
		result := &result{
			CollectionID: gachaResult.CollectionID,
			Name:         gachaResult.Name,
			Rarity:       gachaResult.Rarity,
			IsNew:        gachaResult.IsNew,
		}
		results = append(results, result)
	}

	h.HttpResponse.Success(writer, gachaDrawResponse{Results: results})

}

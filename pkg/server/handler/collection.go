package handler

import (
	"20dojo-online/pkg/dcontext"
	"20dojo-online/pkg/http/response"
	"20dojo-online/pkg/myerror"
	"20dojo-online/pkg/server/service"
	"log"
	"net/http"
)

type collectionListResponse struct {
	Collections []*collection `json:"collections"`
}

type collection struct {
	CollectionID string `json:"collectionID"`
	Name         string `json:"name"`
	Rarity       int    `json:"rarity"`
	HasItem      bool   `json:"hasItem"`
}

type CollectionHandler struct {
	HttpResponse      response.HttpResponseInterface
	CollectionService service.CollectionServiceInterface
}

func NewCollectionHandler(httpResponse response.HttpResponseInterface, collectionService service.CollectionServiceInterface) *CollectionHandler {
	return &CollectionHandler{
		HttpResponse:      httpResponse,
		CollectionService: collectionService,
	}
}

// HandleCollectionList ユーザのコレクションアイテム一覧情報取得
func (h *CollectionHandler) HandleUserCollectionList(writer http.ResponseWriter, request *http.Request) {

	// コンテキストからユーザidを取得
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

	// ユーザのコレクションアイテム一覧情報取得のロジック
	res, err := h.CollectionService.GetUserCollectionList(&service.GetUserCollectionListRequest{UserID: userID})
	if err != nil {
		err := myerror.ApplicationError{
			Message:       "failed to get user collection item",
			OriginalError: err,
			Code:          http.StatusInternalServerError,
		}
		log.Println(err)
		h.HttpResponse.Failed(writer, err)
		return
	}

	// レスポンスの整形
	var collections []*collection
	for _, collectionItem := range res.CollectionItems {
		collection := &collection{
			CollectionID: collectionItem.CollectionID,
			Name:         collectionItem.Name,
			Rarity:       collectionItem.Rarity,
			HasItem:      collectionItem.HasItem,
		}
		collections = append(collections, collection)
	}

	h.HttpResponse.Success(writer, &collectionListResponse{Collections: collections})
}

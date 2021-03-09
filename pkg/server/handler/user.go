package handler

import (
	"20dojo-online/pkg/myerror"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"

	"20dojo-online/pkg/dcontext"
	"20dojo-online/pkg/http/response"
	"20dojo-online/pkg/server/model"
)

type UserHandler struct {
	HttpResponse   response.HttpResponseInterface
	UserRepository model.UserRepositoryInterface
}

func NewUserHandler(httpResponse response.HttpResponseInterface, userRepository model.UserRepositoryInterface) *UserHandler {
	return &UserHandler{
		HttpResponse:   httpResponse,
		UserRepository: userRepository,
	}
}

type userCreateRequest struct {
	Name string `json:"name"`
}

type userCreateResponse struct {
	Token string `json:"token"`
}

// HandleUserCreate ユーザ情報作成処理
func (h *UserHandler) HandleUserCreate(writer http.ResponseWriter, request *http.Request) {

	// リクエストBodyから更新後情報を取得
	var requestBody userCreateRequest
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

	// UUIDでユーザIDを生成する
	userID, err := uuid.NewRandom()
	if err != nil {
		err = myerror.ApplicationError{
			Message:       "failed to generate userID",
			OriginalError: err,
			Code:          http.StatusInternalServerError,
		}
		log.Println(err)
		h.HttpResponse.Failed(writer, err)
		return
	}

	// UUIDで認証トークンを生成する
	authToken, err := uuid.NewRandom()
	if err != nil {
		err = myerror.ApplicationError{
			Message:       "failed to generate token",
			OriginalError: err,
			Code:          http.StatusInternalServerError,
		}
		log.Println(err)
		h.HttpResponse.Failed(writer, err)
		return
	}

	// データベースにユーザデータを登録する
	if err = h.UserRepository.InsertUser(&model.User{
		ID:        userID.String(),
		AuthToken: authToken.String(),
		Name:      requestBody.Name,
		HighScore: 0,
		Coin:      0,
	}); err != nil {
		err = myerror.ApplicationError{
			Message:       "failed to insert user correctly",
			OriginalError: err,
			Code:          http.StatusInternalServerError,
		}
		log.Println(err)
		h.HttpResponse.Failed(writer, err)
		return
	}

	// 生成した認証トークンを返却
	h.HttpResponse.Success(writer, &userCreateResponse{Token: authToken.String()})
}

type userGetResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	HighScore int    `json:"highScore"`
	Coin      int    `json:"coin"`
}

// HandleUserGet ユーザ情報取得処理
func (h *UserHandler) HandleUserGet(writer http.ResponseWriter, request *http.Request) {

	// Contextから認証済みのユーザIDを取得
	ctx := request.Context()
	userID := dcontext.GetUserIDFromContext(ctx)
	if userID == "" {
		userIDEmptyErr := myerror.ApplicationError{
			Message: "userID from is empty",
			Code:    http.StatusInternalServerError,
		}
		log.Println(userIDEmptyErr)
		h.HttpResponse.Failed(writer, userIDEmptyErr)
		return
	}

	user, err := h.UserRepository.SelectUserByPrimaryKey(userID)
	if err != nil {
		err = myerror.ApplicationError{
			Message:       "failed to select user correctly",
			OriginalError: err,
			Code:          http.StatusInternalServerError,
		}
		log.Println(err)
		h.HttpResponse.Failed(writer, err)
		return
	}
	if user == nil {
		userNotFoundErr := myerror.ApplicationError{
			Message: "user not found",
			Code:    http.StatusInternalServerError,
		}
		log.Println(userNotFoundErr)
		h.HttpResponse.Failed(writer, userNotFoundErr)
		return
	}

	// レスポンスに必要な情報を詰めて返却
	h.HttpResponse.Success(writer, &userGetResponse{
		ID:        user.ID,
		Name:      user.Name,
		HighScore: user.HighScore,
		Coin:      user.Coin,
	})
}

type userUpdateRequest struct {
	Name string `json:"name"`
}

// HandleUserUpdate ユーザ情報更新処理
func (h *UserHandler) HandleUserUpdate(writer http.ResponseWriter, request *http.Request) {

	// リクエストBodyから更新後情報を取得
	var requestBody userUpdateRequest
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

	// Contextから認証済みのユーザIDを取得
	ctx := request.Context()
	userID := dcontext.GetUserIDFromContext(ctx)
	if userID == "" {
		userIDEmptyErr := myerror.ApplicationError{
			Message: "userID from is empty",
			Code:    http.StatusInternalServerError,
		}
		log.Println(userIDEmptyErr)
		h.HttpResponse.Failed(writer, userIDEmptyErr)
		return
	}

	user, err := h.UserRepository.SelectUserByPrimaryKey(userID)
	if err != nil {
		err = myerror.ApplicationError{
			Message:       "failed to select user correctly",
			OriginalError: err,
			Code:          http.StatusInternalServerError,
		}
		log.Println(err)
		h.HttpResponse.Failed(writer, err)
		return
	}
	if user == nil {
		userNotFoundErr := myerror.ApplicationError{
			Message: "user not found",
			Code:    http.StatusInternalServerError,
		}
		log.Println(userNotFoundErr)
		h.HttpResponse.Failed(writer, userNotFoundErr)
		return
	}

	user.Name = requestBody.Name
	if err = h.UserRepository.UpdateUserByPrimaryKey(user); err != nil {
		err = myerror.ApplicationError{
			Message:       "failed to update user correctly",
			OriginalError: err,
			Code:          http.StatusInternalServerError,
		}
		log.Println(err)
		h.HttpResponse.Failed(writer, err)
		return
	}

	h.HttpResponse.Success(writer, nil)
}

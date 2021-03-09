package middleware

import (
	"20dojo-online/pkg/myerror"
	"context"
	"log"
	"net/http"

	"20dojo-online/pkg/dcontext"
	"20dojo-online/pkg/http/response"
	"20dojo-online/pkg/server/model"
)

type Middleware struct {
	HttpResponse   response.HttpResponseInterface
	UserRepository model.UserRepositoryInterface
}

func NewMiddleware(httpResponse response.HttpResponseInterface, userRepository model.UserRepositoryInterface) *Middleware {
	return &Middleware{
		HttpResponse:   httpResponse,
		UserRepository: userRepository,
	}
}

// Authenticate ユーザ認証を行ってContextへユーザID情報を保存する
func (m *Middleware) Authenticate(nextFunc http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		ctx := request.Context()
		if ctx == nil {
			ctx = context.Background()
		}

		// リクエストヘッダからx-token(認証トークン)を取得
		token := request.Header.Get("x-token")
		if token == "" {
			log.Println("x-token is empty")
			return
		}

		user, err := m.UserRepository.SelectUserByAuthToken(token)
		if err != nil {
			err = myerror.ApplicationError{
				Message:       "failed to select user in middleware",
				OriginalError: err,
				Code:          http.StatusInternalServerError,
			}
			log.Println(err)
			m.HttpResponse.Failed(writer, err)
			return
		}
		if user == nil {
			userNotFoundErr := myerror.ApplicationError{
				Message: "user not found",
				Code:    http.StatusInternalServerError,
			}
			log.Println(userNotFoundErr)
			m.HttpResponse.Failed(writer, userNotFoundErr)
			return
		}

		// ユーザIDをContextへ保存して以降の処理に利用する
		ctx = dcontext.SetUserID(ctx, user.ID)

		// 次の処理
		nextFunc(writer, request.WithContext(ctx))
	}
}

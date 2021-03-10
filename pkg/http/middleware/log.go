package middleware

import (
	"20dojo-online/pkg/dcontext"
	"20dojo-online/pkg/logging"
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func AccessLogging(nextFunc http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		if ctx == nil {
			ctx = context.Background()
		}

		requestID, err := uuid.NewRandom()
		if err != nil {
			log.Println("failed uuid.NewRandom()")
			return
		}
		ctx = dcontext.SetRequestID(ctx, requestID.String())
		request = request.WithContext(ctx)
		logging.AccessLogging(request)

		nextFunc(writer, request)
	}

}

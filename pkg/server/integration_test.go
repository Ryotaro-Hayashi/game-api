package server

import (
	"20dojo-online/pkg/db"
	"20dojo-online/pkg/http/middleware"
	"20dojo-online/pkg/server/handler"
	"20dojo-online/pkg/server/model"
	"20dojo-online/pkg/server/service"
	"encoding/json"
	"reflect"
)

var (
	testUserRepository = model.NewUserRepository(db.Conn)
	testAuthMiddleware = middleware.NewMiddleware(httpResponse, testUserRepository)
	testRankingService = service.NewRankingService(testUserRepository)
	testRankingHandler = handler.NewRankingHandler(httpResponse, testRankingService)
)

// deepEqualString 文字列同士を比較する
func deepEqualString(str1, str2 string) (bool, error) {
	if str1 == str2 {
		return true, nil
	} else {
		var strInterface1 interface{}
		err := json.Unmarshal([]byte(str1), &strInterface1)
		if err != nil {
			return false, err
		}

		var strInterface2 interface{}
		err = json.Unmarshal([]byte(str2), &strInterface2)
		if err != nil {
			return false, err
		}

		if reflect.DeepEqual(strInterface1, strInterface2) {
			return true, nil
		} else {
			return false, nil
		}
	}
}

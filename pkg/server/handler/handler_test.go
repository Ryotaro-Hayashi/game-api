package handler

import (
	"20dojo-online/pkg/server/service/mock_service"
	"encoding/json"
	"reflect"

	"github.com/golang/mock/gomock"
)

type mock struct {
	gameService       *mock_service.MockGameServiceInterface
	gachaService      *mock_service.MockGachaServiceInterface
	rankingService    *mock_service.MockRankingServiceInterface
	collectionService *mock_service.MockCollectionServiceInterface
}

func newMock(ctrl *gomock.Controller) *mock {
	return &mock{
		gameService:       mock_service.NewMockGameServiceInterface(ctrl),
		gachaService:      mock_service.NewMockGachaServiceInterface(ctrl),
		rankingService:    mock_service.NewMockRankingServiceInterface(ctrl),
		collectionService: mock_service.NewMockCollectionServiceInterface(ctrl),
	}
}

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

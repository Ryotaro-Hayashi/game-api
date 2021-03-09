//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./mock_$GOPACKAGE/mock_$GOFILE

package service

import (
	"20dojo-online/pkg/constant"
	"20dojo-online/pkg/db"
	"20dojo-online/pkg/myerror"
	"20dojo-online/pkg/server/model"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

type DrawGachaRequest struct {
	Times  int
	UserID string
}

type DrawGachaResponse struct {
	GachaResults []*GachaResult
}

type GachaResult struct {
	CollectionID string
	Name         string
	Rarity       int
	IsNew        bool
}

type GachaService struct {
	UserRepository               model.UserRepositoryInterface
	GachaProbabilityRepository   model.GachaProbabilityRepositoryInterface
	UserCollectionItemRepository model.UserCollectionItemRepositoryInterface
	CollectionItemRepository     model.CollectionItemRepositoryInterface
}

func NewGachaService(userRepository model.UserRepositoryInterface,
	gachaProbabilityRepository model.GachaProbabilityRepositoryInterface,
	userCollectionItemRepository model.UserCollectionItemRepositoryInterface,
	collectionItemRepository model.CollectionItemRepositoryInterface) *GachaService {

	return &GachaService{
		UserRepository:               userRepository,
		GachaProbabilityRepository:   gachaProbabilityRepository,
		UserCollectionItemRepository: userCollectionItemRepository,
		CollectionItemRepository:     collectionItemRepository,
	}
}

type GachaServiceInterface interface {
	DrawGacha(serviceRequest *DrawGachaRequest) (*DrawGachaResponse, error)
}

var _ GachaServiceInterface = (*GachaService)(nil)

// DrawGacha ガチャ実行時のロジック
func (s *GachaService) DrawGacha(serviceRequest *DrawGachaRequest) (*DrawGachaResponse, error) {

	// ガチャ排出確率情報からratioの合計を計算
	gachaProbabilities, err := s.GachaProbabilityRepository.SelectGachaProbabilityAll()
	if err != nil {
		return nil, err
	}
	var gachaProbabilitySum int // ratioの合計
	for _, gachaProbability := range gachaProbabilities {
		gachaProbabilitySum += gachaProbability.Ratio
	}

	// 排出アイテムの決定
	gottenCollectionItemIDSlice := make([]string, 0, serviceRequest.Times) // 排出アイテムのidを入れるスライス
	for i := 0; i < serviceRequest.Times; i++ {
		randomNum := rand.Intn(gachaProbabilitySum) // 0からratioの合計までの整数で乱数を生成
		gachaProbabilityThreshold := 0              // 閾値
		var gottenCollectionItemID string
		for _, gachaProbability := range gachaProbabilities {
			gachaProbabilityThreshold += gachaProbability.Ratio
			if randomNum < gachaProbabilityThreshold {
				gottenCollectionItemID = gachaProbability.CollectionItemId
				break
			}
		}
		gottenCollectionItemIDSlice = append(gottenCollectionItemIDSlice, gottenCollectionItemID)
	}

	// ユーザの全所持アイテムを取得
	userCollectionItems, err := s.UserCollectionItemRepository.SelectUserCollectionItemsByUserID(serviceRequest.UserID)
	if err != nil {
		return nil, err
	}
	userCollectionItemIDMap := make(map[string]struct{}, len(userCollectionItems)+serviceRequest.Times) // ユーザの所持アイテムのidを入れるマップ
	for _, userCollectionItem := range userCollectionItems {
		userCollectionItemIDMap[userCollectionItem.CollectionItemID] = struct{}{}
	}

	newUserCollectionItemSlice := make([]*model.UserCollectionItem, 0, serviceRequest.Times) // Newアイテムを入れるスライス
	allCollectionItems, err := s.CollectionItemRepository.SelectCollectionItemAll()          // 全アイテムのスライス
	if err != nil {
		return nil, err
	}
	allCollectionItemMap := make(map[string]*model.CollectionItem, len(allCollectionItems)) // idをキーにした全アイテムのマップ
	for _, collectionItem := range allCollectionItems {
		allCollectionItemMap[collectionItem.ID] = collectionItem
	}
	var results []*GachaResult
	// 排出アイテムと所持アイテムを比較してNewアイテムをマップに格納
	for _, gottenCollectionItemID := range gottenCollectionItemIDSlice {
		isNew := false
		if _, ok := userCollectionItemIDMap[gottenCollectionItemID]; !ok { // 既出アイテムかを確認
			isNew = true
			userCollectionItemIDMap[gottenCollectionItemID] = struct{}{}
			newUserCollectionItem := &model.UserCollectionItem{
				UserID:           serviceRequest.UserID,
				CollectionItemID: gottenCollectionItemID,
			}
			newUserCollectionItemSlice = append(newUserCollectionItemSlice, newUserCollectionItem)
		}
		// レスポンスデータを整形
		collectionItem := allCollectionItemMap[gottenCollectionItemID]
		gachaResult := &GachaResult{
			CollectionID: collectionItem.ID,
			Name:         collectionItem.Name,
			Rarity:       collectionItem.Rarity,
			IsNew:        isNew,
		}
		results = append(results, gachaResult)
	}

	// トランザクション開始
	tx, err := db.Conn.Begin()
	if err != nil {
		return nil, err
	}

	// ユーザ情報を排他ロック
	user, err := s.UserRepository.SelectUserByPrimaryKeyForUpdate(tx, serviceRequest.UserID)
	if err != nil {
		return nil, err
	}
	// 消費コインの計算
	gachaCoinConsumptionSum := constant.GachaCoinConsumption * serviceRequest.Times
	coinResult := user.Coin - gachaCoinConsumptionSum // ガチャ実行後の所持コイン

	// 所持コインが足りない場合のバリデーション
	if user.Coin < gachaCoinConsumptionSum {
		coinShortageErr := myerror.ApplicationError{
			Message: fmt.Sprintf("your coin is not enought. your coin=%s", strconv.Itoa(user.Coin)),
			Code:    http.StatusBadRequest,
		}
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Println(fmt.Sprintf("Rollback Error in validating possessed coin: %s", rollbackErr))
		}
		return nil, coinShortageErr
	}

	// Newアイテムがある時のみ登録
	if len(newUserCollectionItemSlice) >= 1 {
		if err = s.UserCollectionItemRepository.BulkInsertUserCollectionItem(tx, newUserCollectionItemSlice); err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Println(fmt.Sprintf("Rollback Error in inserting to user_collection_item table: %s", rollbackErr))
			}
			return nil, err
		}
	}

	// コインの消費
	if err := s.UserRepository.UpdateUserCoinByPrimaryKey(tx, serviceRequest.UserID, coinResult); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Println(fmt.Sprintf("Rollback Error in updating user coin: %s", rollbackErr))
		}
		return nil, err
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return nil, commitErr
	}

	return &DrawGachaResponse{GachaResults: results}, err
}

//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./mock_$GOPACKAGE/mock_$GOFILE

package service

import "20dojo-online/pkg/server/model"

type GetUserCollectionListRequest struct {
	UserID string
}

type GetUserCollectionListResponse struct {
	CollectionItems []*CollectionItem
}

type CollectionItem struct {
	CollectionID string
	Name         string
	Rarity       int
	HasItem      bool
}

type CollectionService struct {
	UserCollectionItemRepository model.UserCollectionItemRepositoryInterface
	CollectionItemRepository     model.CollectionItemRepositoryInterface
}

func NewCollectionService(userCollectionItemRepository model.UserCollectionItemRepositoryInterface, collectionItemRepository model.CollectionItemRepositoryInterface) *CollectionService {
	return &CollectionService{
		UserCollectionItemRepository: userCollectionItemRepository,
		CollectionItemRepository:     collectionItemRepository,
	}
}

type CollectionServiceInterface interface {
	GetUserCollectionList(serviceRequest *GetUserCollectionListRequest) (*GetUserCollectionListResponse, error)
}

var _ CollectionServiceInterface = (*CollectionService)(nil)

// GetCollectionList ユーザのコレクションアイテム一覧情報取得のロジック
func (s *CollectionService) GetUserCollectionList(serviceRequest *GetUserCollectionListRequest) (*GetUserCollectionListResponse, error) {
	// コレクションアイテムを全件取得
	collectionItems, err := s.CollectionItemRepository.SelectCollectionItemAll()
	if err != nil {
		return nil, err
	}

	// ユーザの所持アイテムを取得
	userCollectionItems, err := s.UserCollectionItemRepository.SelectUserCollectionItemsByUserID(serviceRequest.UserID)
	if err != nil {
		return nil, err
	}

	// ユーザの所持アイテムをマップに変換
	userCollectionItemsMap := make(map[string]struct{}, len(userCollectionItems))
	for _, userCollectionItem := range userCollectionItems {
		userCollectionItemsMap[userCollectionItem.CollectionItemID] = struct{}{}
	}

	// ユーザの所持アイテムをチェック
	collectionItemList := make([]*CollectionItem, 0, len(collectionItems))
	for _, collectionItem := range collectionItems {
		userCollectionItem := &CollectionItem{
			CollectionID: collectionItem.ID,
			Name:         collectionItem.Name,
			Rarity:       collectionItem.Rarity,
		}
		if _, ok := userCollectionItemsMap[collectionItem.ID]; ok {
			userCollectionItem.HasItem = true
		}
		collectionItemList = append(collectionItemList, userCollectionItem)
	}

	return &GetUserCollectionListResponse{CollectionItems: collectionItemList}, err
}

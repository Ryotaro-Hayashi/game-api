//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./mock_$GOPACKAGE/mock_$GOFILE

package service

import "20dojo-online/pkg/server/model"

type GetRankInfoListRequest struct {
	Limit  int
	Offset int
}

type GetRankInfoListResponse struct {
	RankInfoList []*RankInfo
}

// RankInfo ランキング情報
type RankInfo struct {
	UserId   string
	UserName string
	Rank     int
	Score    int
}

type RankingService struct {
	UserRepository model.UserRepositoryInterface
}

var _ RankingServiceInterface = (*RankingService)(nil)

func NewRankingService(userRepository model.UserRepositoryInterface) *RankingService {
	return &RankingService{
		UserRepository: userRepository,
	}
}

type RankingServiceInterface interface {
	GetRankInfoList(serviceRequest *GetRankInfoListRequest) (*GetRankInfoListResponse, error)
}

var _ RankingServiceInterface = (*RankingService)(nil)

// GetRankInfoList ランキング情報取得時のロジック
func (s *RankingService) GetRankInfoList(serviceRequest *GetRankInfoListRequest) (*GetRankInfoListResponse, error) {

	// ハイスコア順に指定順位から指定件数を取得
	usersOrderByHighScoreDesc, err := s.UserRepository.SelectUsersOrderByHighScoreDesc(serviceRequest.Limit, serviceRequest.Offset)
	if err != nil {
		return nil, err
	}

	var rankInfoList []*RankInfo

	// ランク付け
	for index, userRankedIn := range usersOrderByHighScoreDesc {
		rankInfo := &RankInfo{
			UserId:   userRankedIn.ID,
			UserName: userRankedIn.Name,
			Rank:     serviceRequest.Offset + index,
			Score:    userRankedIn.HighScore,
		}
		rankInfoList = append(rankInfoList, rankInfo)
	}

	return &GetRankInfoListResponse{RankInfoList: rankInfoList}, nil
}

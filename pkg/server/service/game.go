//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./mock_$GOPACKAGE/mock_$GOFILE

package service

import (
	"20dojo-online/pkg/constant"
	"20dojo-online/pkg/server/model"
	"errors"
	"fmt"
)

type FinishGameRequest struct {
	UserId string
	Score  int
}

type FinishGameResponse struct {
	Coin int
}

type GameService struct {
	UserRepository model.UserRepositoryInterface
}

func NewGameService(userRepository model.UserRepositoryInterface) *GameService {
	return &GameService{
		UserRepository: userRepository,
	}
}

type GameServiceInterface interface {
	FinishGame(serviceRequest *FinishGameRequest) (*FinishGameResponse, error)
}

var _ GameServiceInterface = (*GameService)(nil)

// GameFinish ゲーム終了時のロジック
func (s *GameService) FinishGame(serviceRequest *FinishGameRequest) (*FinishGameResponse, error) {
	// 報酬の計算
	rewardCoin := int(float64(serviceRequest.Score) * constant.RewardCoinRate)

	// ゲーム終了前のユーザ情報の取得
	user, err := s.UserRepository.SelectUserByPrimaryKey(serviceRequest.UserId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		err = errors.New(fmt.Sprintf("user not found. userID=%s", serviceRequest.UserId))
		return nil, err
	}

	// ユーザのハイスコアとリクエストのスコアを比較
	if user.HighScore < serviceRequest.Score {
		user.HighScore = serviceRequest.Score
	}
	user.Coin += rewardCoin // 所持コイン

	// 所持コインとハイスコアを更新
	if err = s.UserRepository.UpdateUserCoinAndHighScoreByPrimaryKey(user.ID, user.Coin, user.HighScore); err != nil {
		return nil, err
	}

	return &FinishGameResponse{Coin: rewardCoin}, err
}

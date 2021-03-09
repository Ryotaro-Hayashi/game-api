package service

import (
	"20dojo-online/pkg/server/model"
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestRankingService_GetRankInfoList(t *testing.T) {

	type args struct {
		serviceRequest *GetRankInfoListRequest
	}

	tests := []struct {
		name    string
		args    args
		before  func(mock *mockRepository, args args)
		want    *GetRankInfoListResponse
		wantErr bool
	}{
		{
			name: "正常:ユーザ複数",
			before: func(mock *mockRepository, args args) {
				mock.userRepository.EXPECT().SelectUsersOrderByHighScoreDesc(
					args.serviceRequest.Limit, args.serviceRequest.Offset).Return([]*model.User{
					{
						ID:        "UserId1",
						AuthToken: "User1AuthToken",
						Name:      "User1",
						HighScore: 10000,
						Coin:      1000,
					},
					{
						ID:        "UserId2",
						AuthToken: "User2AuthToken",
						Name:      "User2",
						HighScore: 10,
						Coin:      1000,
					},
				}, nil)
			},
			args: args{
				serviceRequest: &GetRankInfoListRequest{
					Limit:  10,
					Offset: 1,
				},
			},
			want: &GetRankInfoListResponse{
				RankInfoList: []*RankInfo{
					{
						UserId:   "UserId1",
						UserName: "User1",
						Rank:     1,
						Score:    10000,
					},
					{
						UserId:   "UserId2",
						UserName: "User2",
						Rank:     2,
						Score:    10,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "異常:ユーザ取得エラー",
			args: args{
				serviceRequest: &GetRankInfoListRequest{
					Limit:  10,
					Offset: 1,
				},
			},
			before: func(mock *mockRepository, args args) {
				mock.userRepository.EXPECT().SelectUsersOrderByHighScoreDesc(
					args.serviceRequest.Limit, args.serviceRequest.Offset).Return([]*model.User{
					nil,
				}, errors.New("SelectUsersOrderByHighScoreDesc failed"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := newMockRepository(ctrl)
			tt.before(mock, tt.args)
			s := NewRankingService(mock.userRepository)
			got, err := s.GetRankInfoList(tt.args.serviceRequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRankInfoList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRankInfoList() got = %v, want %v", got, tt.want)
			}
		})
	}
}

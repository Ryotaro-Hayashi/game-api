package handler

import (
	"20dojo-online/pkg/constant"
	"20dojo-online/pkg/http/response"
	"20dojo-online/pkg/server/service"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestRankingHandler_HandleRankingList(t *testing.T) {
	type args struct {
		request *http.Request
	}
	type want struct {
		statusCode int
		body       string
	}
	tests := []struct {
		name   string
		args   args
		before func(mock *mock, args args)
		want   want
	}{
		{
			name: "正常:ランキング取得",
			args: args{
				request: httptest.NewRequest("GET", "http://localhost:8080/ranking/list?start=1", nil),
			},
			before: func(mock *mock, args args) {
				mock.rankingService.EXPECT().GetRankInfoList(&service.GetRankInfoListRequest{
					Offset: 1,
					Limit:  constant.RankingListLimit,
				}).Return(&service.GetRankInfoListResponse{
					RankInfoList: []*service.RankInfo{
						{
							UserId:   "UserId2",
							UserName: "User2",
							Rank:     1,
							Score:    10000,
						},
						{
							UserId:   "UserId1",
							UserName: "User1",
							Rank:     2,
							Score:    10,
						},
					},
				}, nil)
			},
			want: want{
				statusCode: http.StatusOK,
				body: `{
						  "ranks": [
							{
							  "userId": "UserId2",
							  "userName": "User2",
							  "rank": 1,
							  "score": 10000
							},
							{
							  "userId": "UserId1",
							  "userName": "User1",
							  "rank": 2,
							  "score": 10
							}
						  ]
						}`,
			},
		},
		{
			name: "異常:クエリパラメータエラー",
			args: args{
				request: httptest.NewRequest("GET", "http://localhost:8080/ranking/list?nostart=1", nil),
			},
			before: func(mock *mock, args args) {},
			want: want{
				statusCode: http.StatusBadRequest,
				body: `{
							"code": 400,
							"message": "Bad Request"
						}`,
			},
		},
		{
			name: "異常:開始順位エラー",
			args: args{
				request: httptest.NewRequest("GET", "http://localhost:8080/ranking/list?start=-1", nil),
			},
			before: func(mock *mock, args args) {},
			want: want{
				statusCode: http.StatusBadRequest,
				body: `{
							"code": 400,
							"message": "Bad Request"
						}`,
			},
		},
		{
			name: "異常:開始順位エラー",
			args: args{
				request: httptest.NewRequest("GET", "http://localhost:8080/ranking/list?start=1", nil),
			},
			before: func(mock *mock, args args) {
				mock.rankingService.EXPECT().GetRankInfoList(&service.GetRankInfoListRequest{
					Offset: 1,
					Limit:  constant.RankingListLimit,
				}).Return(nil, errors.New("GetRankInfoList"))
			},
			want: want{
				statusCode: http.StatusInternalServerError,
				body: `{
							"code": 500,
							"message": "Internal Server Error"
						}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := newMock(ctrl)
			tt.before(mock, tt.args)

			writer := httptest.NewRecorder()

			h := NewRankingHandler(response.NewHttpResponse(), mock.rankingService)
			h.HandleRankingList(writer, tt.args.request)

			res := writer.Result()
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("ioutil.ReadAll failed %s", err)
			}

			if res.StatusCode != tt.want.statusCode {
				t.Errorf("status code = %d, want %d", res.StatusCode, tt.want.statusCode)
			}

			boolean, err := deepEqualString(string(body), tt.want.body)
			if err != nil {
				t.Errorf("response.DeepEqualString() failed %s", err)
			}
			if !boolean {
				t.Errorf("response body = \n%s\n, want \n%s\n", string(body), tt.want.body)
			}
		})
	}
}

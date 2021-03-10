package handler

import (
	"20dojo-online/pkg/constant"
	"20dojo-online/pkg/http/response"
	"20dojo-online/pkg/myerror"
	"20dojo-online/pkg/server/service"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type rankingListResponse struct {
	Ranks []*rank `json:"ranks"`
}

// rank ランキング情報
type rank struct {
	UserId   string `json:"userId"`
	UserName string `json:"userName"`
	Rank     int    `json:"rank"`
	Score    int    `json:"score"`
}

type RankingHandler struct {
	HttpResponse   response.HttpResponseInterface
	RankingService service.RankingServiceInterface
}

func NewRankingHandler(httpResponse response.HttpResponseInterface, rankingService service.RankingServiceInterface) *RankingHandler {
	return &RankingHandler{
		HttpResponse:   httpResponse,
		RankingService: rankingService,
	}
}

// HandleRankingList ランキング情報取得
func (h *RankingHandler) HandleRankingList(writer http.ResponseWriter, request *http.Request) {
	// クエリストリングから開始順位の受け取り
	param := request.URL.Query().Get("start")
	start, err := strconv.Atoi(param)
	if err != nil {
		err = myerror.ApplicationError{
			Message: "failed to get query parameter",
			Code:    http.StatusBadRequest,
		}
		log.Println(err)
		h.HttpResponse.Failed(writer, err)
		return
	}
	// startが0以下のときエラーを返す
	if start <= 0 {
		err := myerror.ApplicationError{
			Message: fmt.Sprintf("start rank is 0 or less. start=%d", start),
			Code:    http.StatusBadRequest,
		}
		log.Println(err)
		h.HttpResponse.Failed(writer, err)
		return
	}

	// ランキング情報取得のロジック
	res, err := h.RankingService.GetRankInfoList(&service.GetRankInfoListRequest{
		Offset: start,
		Limit:  constant.RankingListLimit,
	})
	if err != nil {
		err = myerror.ApplicationError{
			Message:       "failed to finish game correctly",
			OriginalError: err,
			Code:          http.StatusInternalServerError,
		}
		log.Println(err)
		h.HttpResponse.Failed(writer, err)
		return
	}

	// レスポンスの整形
	var ranks []*rank
	for _, rankInfo := range res.RankInfoList {
		rank := &rank{
			UserId:   rankInfo.UserId,
			UserName: rankInfo.UserName,
			Rank:     rankInfo.Rank,
			Score:    rankInfo.Score,
		}
		ranks = append(ranks, rank)
	}

	h.HttpResponse.Success(writer, rankingListResponse{Ranks: ranks})
}

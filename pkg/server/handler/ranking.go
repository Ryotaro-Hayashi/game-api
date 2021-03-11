package handler

import (
	"20dojo-online/pkg/constant"
	"20dojo-online/pkg/dcontext"
	"20dojo-online/pkg/http/response"
	"20dojo-online/pkg/logging"
	"20dojo-online/pkg/myerror"
	"20dojo-online/pkg/server/service"
	"fmt"
	"net/http"
	"strconv"

	"go.uber.org/zap"
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
	requestID := dcontext.GetRequestIDFromContext(request.Context())
	logging.ApplicationLogger.Info("start getting rank info list", zap.String("requestID", requestID))

	// クエリストリングから開始順位の受け取り
	param := request.URL.Query().Get("start")
	start, err := strconv.Atoi(param)
	if err != nil {
		err = myerror.ApplicationError{
			Message: "failed to strconv.Atoi()",
			Code:    http.StatusBadRequest,
		}
		logging.ApplicationLogger.Warn("failed to strconv.Atoi()", zap.String("requestID", requestID))
		h.HttpResponse.Failed(writer, err)
		return
	}
	// startが0以下のときエラーを返す
	if start <= 0 {
		err := myerror.ApplicationError{
			Message: fmt.Sprintf("start rank is 0 or less. start=%d", start),
			Code:    http.StatusBadRequest,
		}
		logging.ApplicationLogger.Warn("start rank is 0 or less", zap.String("requestID", requestID))
		h.HttpResponse.Failed(writer, err)
		return
	}
	logging.ApplicationLogger.Debug("success in getting query parameter", zap.String("requestID", requestID), zap.Int("query", start))

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
		logging.ApplicationLogger.Error("failed to finish game correctly", zap.String("requestID", requestID))
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

	logging.ApplicationLogger.Info("success in finishing game", zap.String("requestID", requestID))
	h.HttpResponse.Success(writer, rankingListResponse{Ranks: ranks})
}

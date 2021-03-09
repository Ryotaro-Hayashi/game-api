package handler

import (
	"net/http"

	"20dojo-online/pkg/constant"
	"20dojo-online/pkg/http/response"
)

type settingHandler struct {
	HttpResponse response.HttpResponseInterface
}

func NewSettingHandler(httpResponse response.HttpResponseInterface) *settingHandler {
	return &settingHandler{
		HttpResponse: httpResponse,
	}
}

// HandleSettingGet ゲーム設定情報取得処理
func (h *settingHandler) HandleSettingGet(writer http.ResponseWriter, request *http.Request) {
	h.HttpResponse.Success(writer, &settingGetResponse{
		GachaCoinConsumption: constant.GachaCoinConsumption,
	})
}

type settingGetResponse struct {
	GachaCoinConsumption int `json:"gachaCoinConsumption"`
}

package queue

import (
	"ecom_promotion_v2/internal"
	"ecom_promotion_v2/internal/models"
	"ecom_promotion_v2/internal/repositories"
	"ecom_promotion_v2/internal/utils"
	"encoding/json"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() *zap.Logger {
	loggerCore := zapcore.NewCore(logEndcoder(), zapcore.AddSync(os.Stdout), zap.DebugLevel)
	return zap.New(loggerCore, zap.AddCaller())
}
func logEndcoder() zapcore.Encoder {
	encodeConfig := zap.NewProductionEncoderConfig()
	encodeConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	return zapcore.NewJSONEncoder(encodeConfig)
}

type RespLocal struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}
type buttonActionRedirect struct {
	ActionType string      `json:"actionType"`
	DataAction string      `json:"dataAction"`
	Data       interface{} `json:"data"`
	Title      string      `json:"title"`
}
type Notify struct {
	CustomerPhone         string                 `json:"customerPhone"`
	ContractNo            string                 `json:"contractNo"`
	CustomerId            string                 `json:"customerId"`
	TitleVi               string                 `json:"titleVi"`
	TitleEn               string                 `json:"titleEn"`
	MessageVi             string                 `json:"messageVi"`
	MessageEn             string                 `json:"messageEn"`
	ButtonType            string                 `json:"buttonType"`
	ButtonCloseContentVi  string                 `json:"buttonCloseContentVi"`
	ButtonCloseContentEn  string                 `json:"buttonCloseContentEn"`
	ButtonCloseType       string                 `json:"buttonCloseType"`
	ButtonCloseKey        string                 `json:"buttonCloseKey"`
	ButtonCloseRedirect   string                 `json:"buttonCloseRedirect"`
	ButtonActionContentVi string                 `json:"buttonActionContentVi"`
	ButtonActionContentEn string                 `json:"buttonActionContentEn"`
	ButtonActionType      string                 `json:"buttonActionType"`
	ButtonActionKey       string                 `json:"buttonActionKey"`
	ButtonActionRedirect  buttonActionRedirect   `json:"buttonActionRedirect"`
	IsPinMessage          int                    `json:"isPinMessage"`
	IsDetail              int                    `json:"isDetail"`
	IsSave                int                    `json:"isSave"`
	IsReloadHome          int                    `json:"isReloadHome"`
	NotifyTypeId          string                 `json:"notifyTypeId"`
	ActionType            string                 `json:"actionType"`
	TypeFormat            string                 `json:"typeFormat"`
	Data                  map[string]interface{} `json:"data"`
}
type QueueNotify struct {
	Body   Notify
	URL    string
	Header map[string]string
	Repo   *repositories.Repositories
}

func (c *QueueNotify) Run() {
	param := map[string]string{}
	var input map[string]interface{}
	inrec, _ := json.Marshal(c.Body)
	json.Unmarshal(inrec, &input)
	log := NewLogger()
	log.Info("QueueNotify", zap.Any("input", input))
	responese, err := utils.Request(c.URL, false, c.Header, param, input, 10, false)
	res := RespLocal{}
	storeCall := &models.StoreCallApiTb{
		Url:   c.URL,
		Input: c.Body.ContractNo,
	}
	if err != nil {
		log.Error("Request QueueNotify Error", zap.Any("data", input), zap.Error(err))
		storeCall.Output = err.Error()
		err = c.Repo.StoreCallApi.Create(storeCall)
		if err != nil {
			log.Error("repo.StoreCallApi.Create", zap.Any("data", storeCall), zap.Error(err))
		}
		return
	}
	// YOUR CODE HERE: STORE CALL API
	storeCall.Output = responese.String()
	err = c.Repo.StoreCallApi.Create(storeCall)
	if err != nil {
		log.Error("repo.StoreCallApi.Create", zap.Any("data", storeCall), zap.Error(err))
	}
	json.Unmarshal([]byte(responese.String()), &res)
	log.Info("QueueNotify", zap.Any("input", input), zap.Any("responese", responese))
}

type QueueNotifyTemplate struct {
	TemplateId  int
	CustomerId  int
	ContractNo  string
	PhoneNb     string
	TypeNotify  string
	DataReplace map[string]interface{}
	Data        map[string]interface{}
	Repo        repositories.Repositories
}

func (c *QueueNotifyTemplate) Run() {
	url := internal.Domains.HiFPT
	headers := map[string]string{
		"Authorization": utils.CreateNotifyTemplateToken(),
		"ClientKey":     internal.Keys.NotifyTemplateClientKey,
	}

	param := map[string]string{}
	body := map[string]interface{}{
		"templateId":  c.TemplateId,
		"dataReplace": c.DataReplace,
	}
	if c.Data != nil {
		body["data"] = c.Data
	}
	if c.TypeNotify == "CONTRACT" {
		body["contractNo"] = c.ContractNo
		url += internal.Eps.HiFPTApi.NotifyTemplateByContractNo
	} else if c.TypeNotify == "CUSTOMER_ID" {
		body["customerId"] = c.CustomerId
		url += internal.Eps.HiFPTApi.NotifyTemplateByCustomerId
	} else if c.TypeNotify == "PHONE" {
		body["phone"] = c.PhoneNb
		url += internal.Eps.HiFPTApi.NotifyTemplateByPhone
	}
	byteInput, _ := json.Marshal(body)
	output := ""
	dt := 0.0
	now := utils.GetTimeUTC7()
	defer func() {
		inputStoreCall := &models.StoreCallApiTb{
			Url:     url,
			Input:   string(byteInput),
			Output:  output,
			TCreate: now,
		}
		err := c.Repo.StoreCallApi.Create(inputStoreCall)
		if err != nil {
			internal.Log.Error("SwapWifi6StoreCallApi.Create", zap.Any("input", inputStoreCall))
		}
	}()
	internal.Log.Info("Call", zap.Any("url", url), zap.Any("headers", headers), zap.Any("input", body))
	resp, err := utils.Request(url, false, headers, param, body, 10, false)
	if err != nil {
		internal.Log.Error("Call fail ", zap.Any("url", url), zap.Any("headers", headers), zap.Any("input", body), zap.Error(err))
		output = err.Error()
		return
	}
	internal.Log.Info("Response", zap.Any("url", url), zap.Any("headers", headers), zap.Any("input", body), zap.Any("status_http", resp.StatusCode()), zap.Any("output", resp.String()), zap.Any("dt", resp.Time().Seconds()), zap.Any("http_status", resp.StatusCode()))
	dt = utils.GetTimeUTC7().Sub(now).Seconds()
	output = fmt.Sprintf("status:%v response:%v dt:%v", resp.StatusCode(), resp.String(), dt)
}

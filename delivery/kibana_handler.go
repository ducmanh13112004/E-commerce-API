package delivery

import (
	"ecom_promotion_v2/internal"
	"ecom_promotion_v2/internal/models"
	"ecom_promotion_v2/internal/utils"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func WebkitSendKibana(ctx *fiber.Ctx, funcName string, uri string, tokenAuth string, body interface{}, resultFe models.RespWeb, startTime time.Time, userWeb models.UserInfo) {
	dt := time.Since(startTime).Seconds()
	internal.Log.Info("Response", zap.Any("dt", dt), zap.Any("url", uri), zap.Any("body", body), zap.Any("result", resultFe))
	errPool := internal.PoolLog.Execute(&utils.SendLogToKibanaTask{
		ServiceName:     internal.ServiceName,
		BootstrapServer: internal.Brokers,
		TopicName:       internal.KafkaTopicName,
		Message: utils.KibanaMessage{
			ServiceName:  internal.ServiceName,
			UserAgent:    string(ctx.Context().UserAgent()),
			FuncName:     funcName,
			ActionName:   uri,
			Token:        tokenAuth,
			Input:        body,
			Output:       resultFe,
			ExecutedTime: dt,
			Version:      utils.GetTimeUTC7().String(),
			Url:          ctx.Request().URI().String(),
		},
	})
	if errPool != nil {
		internal.Log.Error("Error Execute Pool Kibana", zap.Any("userWeb", userWeb), zap.Any("body", body), zap.Error(errPool))
	}

	stringHeaders, _ := json.Marshal(ctx.GetReqHeaders())
	stringBody, _ := json.Marshal(body)
	stringResp, _ := json.Marshal(resultFe)
	// TODO: Send to response to kibana all server
	errPoolAll := internal.PoolLog.Execute(&utils.SendLogToKibanaAllTask{
		BootstrapServer: internal.Brokers,
		Log:             internal.Logger,
		TopicName:       internal.KafkaTopicNameAll,
		Message: utils.KibanaMessageAll{
			AppVersion:   userWeb.AppVersion,
			ScreenId:     "",
			IpAddress:    userWeb.CustomerIp,
			Phone:        userWeb.PhoneNb,
			CustomerId:   userWeb.CustomerId,
			Status:       resultFe.Status,
			Input:        string(stringBody),
			Output:       string(stringResp),
			Headers:      string(stringHeaders),
			FunctionName: ctx.Path(),
			ActionName:   uri,
			DateAction:   utils.GetTimeUTC7().Format("2006-01-02 15:04:05"),
			Url:          ctx.Context().URI().String(),
			Note:         fmt.Sprintf("status:%v,msg:%v", resultFe.Status, resultFe.Msg),
			TypeLog:      "Webkit",
			ProcessTime:  dt,
			Topic_name:   internal.KafkaTopicNameAll,
			ServiceName:  internal.ServiceName,
		},
	})
	if errPoolAll != nil {
		internal.Log.Error("Error Execute Pool Kibana All", zap.Any("userWeb", userWeb), zap.Any("body", body), zap.Error(errPoolAll))
	}
}
func AppSendKibana(ctx *fiber.Ctx, funcName string, uri string, tokenAuth string, body interface{}, resultFe interface{}, startTime time.Time, userWeb models.UserInfo) {
	resultTmp := models.RespLocal{}
	byteData, _ := json.Marshal(resultFe)
	err := json.Unmarshal(byteData, &resultTmp)
	if err != nil {
		internal.Log.Error("Error Execute Pool Kibana", zap.Error(err))
	}
	dt := time.Since(startTime).Seconds()
	internal.Log.Info("Response", zap.Any("dt", dt), zap.Any("url", uri), zap.Any("body", body), zap.Any("result", resultFe))
	errPool := internal.PoolLog.Execute(&utils.SendLogToKibanaTask{
		ServiceName:     internal.ServiceName,
		BootstrapServer: internal.Brokers,
		TopicName:       internal.KafkaTopicName,
		Message: utils.KibanaMessage{
			ServiceName:  internal.ServiceName,
			UserAgent:    string(ctx.Context().UserAgent()),
			FuncName:     funcName + "-" + userWeb.CustomerId,
			ActionName:   uri,
			Token:        userWeb.AccessToken,
			Input:        body,
			Output:       resultFe,
			ExecutedTime: dt,
			Version:      userWeb.AppVersion,
			Url:          ctx.Request().URI().String(),
		},
	})
	if errPool != nil {
		internal.Log.Error("Error Execute Pool Kibana", zap.Any("userWeb", userWeb), zap.Any("body", body), zap.Error(errPool))
	}

	stringHeaders, _ := json.Marshal(ctx.GetReqHeaders())
	stringBody, _ := json.Marshal(body)
	stringResp, _ := json.Marshal(resultFe)
	// TODO: Send to response to kibana all server
	errPoolAll := internal.PoolLog.Execute(&utils.SendLogToKibanaAllTask{
		BootstrapServer: internal.Brokers,
		Log:             internal.Logger,
		TopicName:       internal.KafkaTopicNameAll,
		Message: utils.KibanaMessageAll{
			AppVersion:   userWeb.AppVersion,
			ScreenId:     "",
			IpAddress:    userWeb.CustomerIp,
			Phone:        userWeb.PhoneNb,
			CustomerId:   userWeb.CustomerId,
			Status:       resultTmp.StatusCode,
			Input:        string(stringBody),
			Output:       string(stringResp),
			Headers:      string(stringHeaders),
			FunctionName: ctx.Path(),
			ActionName:   uri,
			DateAction:   utils.GetTimeUTC7().Format("2006-01-02 15:04:05"),
			Url:          ctx.Context().URI().String(),
			Note:         fmt.Sprintf("status:%v,msg:%v", resultTmp.StatusCode, resultTmp.Message),
			TypeLog:      "App",
			ProcessTime:  dt,
			Topic_name:   internal.KafkaTopicNameAll,
			ServiceName:  internal.ServiceName,
		},
	})
	if errPoolAll != nil {
		internal.Log.Error("Error Execute Pool Kibana All", zap.Any("userWeb", userWeb), zap.Any("body", body), zap.Error(errPoolAll))
	}
}
func LocalSendKibana(ctx *fiber.Ctx, funcName string, uri string, tokenAuth string, body interface{}, resultFe models.RespLocal, startTime time.Time) {
	dt := time.Since(startTime).Seconds()
	internal.Log.Info("Response", zap.Any("dt", dt), zap.Any("url", uri), zap.Any("body", body), zap.Any("result", resultFe))
	errPool := internal.PoolLog.Execute(&utils.SendLogToKibanaTask{
		ServiceName:     internal.ServiceName,
		BootstrapServer: internal.Brokers,
		TopicName:       internal.KafkaTopicName,
		Message: utils.KibanaMessage{
			ServiceName:  internal.ServiceName,
			UserAgent:    string(ctx.Context().UserAgent()),
			FuncName:     funcName,
			ActionName:   uri,
			Token:        string(ctx.Request().Header.Peek("TOKEN")),
			Input:        body,
			Output:       resultFe,
			ExecutedTime: time.Since(startTime).Seconds(),
			Version:      utils.GetTimeUTC7().String(),
			Url:          ctx.Request().URI().String(),
		},
	})
	if errPool != nil {
		internal.Log.Error("Error Execute Pool Kibana", zap.Any("body", body), zap.Error(errPool))
	}

	stringHeaders, _ := json.Marshal(ctx.GetReqHeaders())
	stringBody, _ := json.Marshal(body)
	stringResp, _ := json.Marshal(resultFe)
	// TODO: Send to response to kibana all server
	errPoolAll := internal.PoolLog.Execute(&utils.SendLogToKibanaAllTask{
		BootstrapServer: internal.Brokers,
		Log:             internal.Logger,
		TopicName:       internal.KafkaTopicNameAll,
		Message: utils.KibanaMessageAll{
			AppVersion:   "",
			ScreenId:     "",
			Phone:        "",
			CustomerId:   "",
			IpAddress:    "",
			Status:       resultFe.StatusCode,
			Input:        string(stringBody),
			Output:       string(stringResp),
			Headers:      string(stringHeaders),
			FunctionName: ctx.Path(),
			ActionName:   uri,
			DateAction:   utils.GetStringTimeUTC7("Y-M-D H:M:S"),
			Url:          ctx.Context().URI().String(),
			Note:         fmt.Sprintf("statusCode:%v,message:%v", resultFe.StatusCode, resultFe.Message),
			TypeLog:      "Local",
			ProcessTime:  dt,
			Topic_name:   internal.KafkaTopicNameAll,
			ServiceName:  internal.ServiceName,
		},
	})
	if errPoolAll != nil {
		internal.Log.Error("Error Execute Pool Kibana All", zap.Any("body", body), zap.Error(errPoolAll))
	}
}

func PortalSendKibana(ctx *fiber.Ctx, funcName string, uri string, tokenAuth string, body interface{}, resultFe models.RespWeb, startTime time.Time, userPortal models.PortalPayload) {
	dt := time.Since(startTime).Seconds()
	internal.Log.Info("Response", zap.Any("dt", dt), zap.Any("url", uri), zap.Any("body", body), zap.Any("result", resultFe))
	errPool := internal.PoolLog.Execute(&utils.SendLogToKibanaTask{
		ServiceName:     internal.ServiceName,
		BootstrapServer: internal.Brokers,
		TopicName:       internal.KafkaTopicName,
		Message: utils.KibanaMessage{
			ServiceName:  internal.ServiceName,
			UserAgent:    string(ctx.Context().UserAgent()),
			FuncName:     funcName,
			ActionName:   uri,
			Token:        tokenAuth,
			Input:        body,
			Output:       resultFe,
			ExecutedTime: dt,
			Version:      utils.GetTimeUTC7().String(),
			Url:          ctx.Request().URI().String(),
		},
	})
	if errPool != nil {
		internal.Log.Error("Error Execute Pool Kibana", zap.Any("userPortal", userPortal), zap.Any("body", body), zap.Error(errPool))
	}

	stringHeaders, _ := json.Marshal(ctx.GetReqHeaders())
	stringBody, _ := json.Marshal(body)
	stringResp, _ := json.Marshal(resultFe)
	// TODO: Send to response to kibana all server
	errPoolAll := internal.PoolLog.Execute(&utils.SendLogToKibanaAllTask{
		BootstrapServer: internal.Brokers,
		Log:             internal.Logger,
		TopicName:       internal.KafkaTopicNameAll,
		Message: utils.KibanaMessageAll{
			Phone:        userPortal.PhoneNb,
			Status:       resultFe.Status,
			Input:        string(stringBody),
			Output:       string(stringResp),
			Headers:      string(stringHeaders),
			FunctionName: ctx.Path(),
			ActionName:   uri,
			DateAction:   utils.GetTimeUTC7().Format("2006-01-02 15:04:05"),
			Url:          ctx.Context().URI().String(),
			Note:         fmt.Sprintf("status:%v,msg:%v", resultFe.Status, resultFe.Msg),
			TypeLog:      "Portal",
			ProcessTime:  dt,
			Topic_name:   internal.KafkaTopicNameAll,
			ServiceName:  internal.ServiceName,
		},
	})
	if errPoolAll != nil {
		internal.Log.Error("Error Execute Pool Kibana All", zap.Any("userPortal", userPortal), zap.Any("body", body), zap.Error(errPoolAll))
	}
}

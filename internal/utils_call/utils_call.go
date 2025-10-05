package utils_call

import (
	"ecom_promotion_v2/internal"
	"ecom_promotion_v2/internal/models"
	"ecom_promotion_v2/internal/repositories"
	"ecom_promotion_v2/internal/utils"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

// func CallFboxGetScheduleTimeline(c *settings.AppSettings, input models.GetScheduleTimeline) (int, *models.GetScheduleTimelineCall, error) {
// 	headers := map[string]string{}
// 	url := c.AppUrls.Urls["PARTNERS_FBOX"] + "/hifpt"
// 	params := map[string]string{
// 		"method":    input.Method,
// 		"Date":      input.Date,
// 		"hour":      strconv.Itoa(input.Hour),
// 		"limithour": strconv.Itoa(input.Limithour),
// 		"Start":     strconv.Itoa(input.Start),
// 		"Limit":     strconv.Itoa(input.Limit),
// 	}
// 	body := map[string]interface{}{}

// 	c.Log.Info("Call "+url, zap.Any("headers", headers), zap.Any("params", params), zap.Any("body", body))
// 	resp, err := utils.Request(url, true, headers, params, body, 10, false)

// 	if err != nil || resp.StatusCode() != 200 {
// 		c.Log.Error("Error Call", zap.Any("resp.StatusCode", resp.StatusCode()), zap.Error(err), zap.Any("resp", resp))
// 		return c.ErrMsgs.CallFail.Code, nil, errors.New(c.ErrMsgs.CallFail.Msg)
// 	}
// 	c.Log.Info("Response", zap.Any("url", url), zap.Any("params", params), zap.Any("resp.StatusCode", resp.StatusCode()))

// 	res := &models.GetScheduleTimelineCall{}
// 	err = json.Unmarshal([]byte(resp.String()), &res)
// 	if err != nil {
// 		c.Log.Error("Error Unmarshal", zap.Error(err))
// 		return 0, nil, err
// 	}

// 	return 1, res, nil
// }

func CallGetCustomerTypeFromPhone(repo *repositories.Repositories, customer *models.UserInfo) (*models.RespGetContractsBOT, *internal.SystemStatus) {
	defer utils.ExecTime(utils.GetTimeUTC7(), "CallGetCustomerTypeFromPhone", nil)
	now := utils.GetTimeUTC7()

	url := internal.Domains.HiFPT + internal.Eps.HiFPTApi.GetContracts
	headers := map[string]string{
		"TOKEN": internal.Keys.TokenKeyHiChatBot,
	}
	input := map[string]interface{}{
		"phone": customer.PhoneNb,
	}
	byteBody, _ := json.Marshal(input)
	storeCall := &models.StoreCallApiTb{
		Url:     url,
		Input:   string(byteBody),
		TCreate: now,
	}
	defer func() {
		err := repo.StoreCallApi.Create(storeCall)
		if err != nil {
			internal.Log.Error("StoreCallApi.Create", zap.Any("input", storeCall), zap.Error(err))
		}
	}()
	internal.Log.Info("Call "+url, zap.Any("input", input))
	resp, err := utils.Request(url, true, headers, map[string]string{}, input, 15, true)
	if err != nil {
		internal.Log.Error("Call "+url, zap.Any("input", input), zap.Error(err))
		tReceive := utils.GetTimeUTC7()
		storeCall.Output = err.Error()
		storeCall.TResponse = tReceive
		storeCall.Dt = tReceive.Sub(now).Seconds()
		return nil, internal.SysStatus.SystemError
	}
	internal.Log.Info("Response", zap.Any("url", url), zap.Any("header", headers), zap.Any("input", input), zap.Any("response", resp.String()), zap.Any("dt", resp.Time().Seconds()), zap.Any("http_status", resp.StatusCode()))
	tReceive := utils.GetTimeUTC7()
	storeCall.TResponse = tReceive
	storeCall.Dt = tReceive.Sub(now).Seconds()
	storeCall.Output = resp.String()
	if resp.StatusCode() != 200 {
		return nil, &internal.SystemStatus{
			Status: internal.SysStatus.SystemBusy.Status,
			Msg:    internal.SysStatus.SystemBusy.Msg,
			Detail: fmt.Sprintf("http status code: %s", resp.Status()),
		}
	}
	res := &models.RespGetContractsBOT{}
	err = json.Unmarshal([]byte(resp.String()), &res)
	if err != nil {
		internal.Log.Error("Unmarshal", zap.Any("input", resp.String()), zap.Error(err))
		return nil, internal.SysStatus.SystemError
	}
	return res, nil
}

func CallGetGetVoucherByProgramId(customerId, programId string) (string, *internal.SystemStatus) {
	url := internal.Domains.HiFPT + internal.Eps.HiFPTApi.GetVoucherByProgramId
	header := map[string]string{
		"token": utils.CreateEcomTokenV2(),
	}
	body := map[string]interface{}{
		"customer_id": customerId,
		"program_id":  programId,
	}
	internal.Log.Info("Call", zap.Any("url", url), zap.Any("header", header), zap.Any("body", body))
	resp, err := utils.Request(url, false, header, nil, body, 10, false)
	if err != nil {
		internal.Log.Error("Call fail", zap.Any("url", url), zap.Any("header", header), zap.Any("body", body), zap.Error(err))
		return "", internal.SysStatus.SystemBusy
	}
	internal.Log.Info("Response", zap.Any("url", url), zap.Any("header", header), zap.Any("body", body), zap.Any("httpStatus", resp.Status()), zap.Any("dt", fmt.Sprintf("%v", resp.Time().Seconds())), zap.Any("resp", resp.String()))
	res := struct {
		StatusCode int    `json:"statusCode"`
		Message    string `json:"message"`
		Data       struct {
			Evoucher string `json:"evoucher"`
		} `json:"data"`
	}{}
	err = json.Unmarshal(resp.Body(), &res)
	if err != nil {
		internal.Log.Error("json.Unmarshal", zap.Any("resp", resp.String()), zap.Any("programId", programId), zap.Any("customerId", customerId), zap.Error(err))
		return "", internal.SysStatus.SystemError
	}
	if res.StatusCode != 0 {
		return "", internal.SysStatus.SystemBusy
	}
	return res.Data.Evoucher, nil
}
func CallSendNotiByTemplate(dataReplace map[string]interface{}, templateId int, sendBy interface{}, typeNoti string) (*models.RespLocal, *internal.SystemStatus) {
	url := internal.Domains.HiFPT
	body := map[string]interface{}{
		"templateId":  templateId,
		"dataReplace": dataReplace,
	}
	if typeNoti == "CONTRACT" {
		body["contractNo"] = sendBy
		url += internal.Eps.HiFPTApi.NotifyTemplateByContractNo
	} else if typeNoti == "CUSTOMER_ID" {
		body["customerId"] = sendBy
		url += internal.Eps.HiFPTApi.NotifyTemplateByCustomerId
	} else if typeNoti == "PHONE" {
		body["phone"] = sendBy
		url += internal.Eps.HiFPTApi.NotifyTemplateByPhone
	}
	headers := map[string]string{
		"Authorization": utils.CreateNotifyTemplateToken(),
		"ClientKey":     internal.Keys.NotifyTemplateClientKey,
	}
	resp, err := utils.Request(url, false, headers, nil, body, 10, false)
	if err != nil {
		internal.Log.Error("CallSendNotiByTemplate Request", zap.Any("url", url), zap.Any("headers", headers), zap.Any("body", body), zap.Error(err))
		return nil, internal.SysStatus.SystemBusy
	}
	internal.Log.Info("Response CallSendNotiByTemplate", zap.Any("url", url), zap.Any("headers", headers), zap.Any("body", body), zap.Any("httpStatus", resp.Status()), zap.Any("dt", fmt.Sprintf("%v", resp.Time().Seconds())), zap.Any("resp", resp.String()))
	if resp.StatusCode() != 200 {
		internal.Log.Error("CallSendNotiByTemplate resp", zap.Any("resp", resp.String()), zap.Any("url", url), zap.Any("body", body), zap.Any("headers", headers))
		return nil, internal.SysStatus.SystemBusy
	}
	respLocal := &models.RespLocal{}
	err = json.Unmarshal(resp.Body(), &respLocal)
	if err != nil {
		internal.Log.Error("CallSendNotiByTemplate Unmarshal", zap.Error(err), zap.Any("resp", resp.String()), zap.Any("url", url), zap.Any("body", body), zap.Any("headers", headers))
		return nil, internal.SysStatus.SystemError
	}
	if respLocal.StatusCode != 0 {
		return nil, &internal.SystemStatus{
			Status: respLocal.StatusCode,
			Msg:    respLocal.Message,
		}
	}
	return respLocal, nil
}
func CallSendMailByTemplate(dataReplace map[string]interface{}, templateId int, mailTo, ccTo []string) (*models.RespLocal, *internal.SystemStatus) {
	url := internal.Domains.HiFPT + internal.Eps.HiFPTApi.SendMailByTemplateId
	body := map[string]interface{}{
		"templateId": templateId,
		"data":       dataReplace,
		"mailTo":     mailTo,
		"cc":         ccTo,
	}
	headers := map[string]string{
		"Authorization": utils.CreateNotifyTemplateToken(),
		"ClientKey":     internal.Keys.NotifyTemplateClientKey,
	}
	resp, err := utils.Request(url, false, headers, nil, body, 10, false)
	if err != nil {
		internal.Log.Error("CallSendMailByTemplate Request", zap.Any("url", url), zap.Any("headers", headers), zap.Any("body", body), zap.Error(err))
		return nil, internal.SysStatus.SystemBusy
	}
	internal.Log.Info("Response CallSendMailByTemplate", zap.Any("url", url), zap.Any("headers", headers), zap.Any("body", body), zap.Any("httpStatus", resp.Status()), zap.Any("dt", fmt.Sprintf("%v", resp.Time().Seconds())), zap.Any("resp", resp.String()))
	if resp.StatusCode() != 200 {
		internal.Log.Error("CallSendMailByTemplate resp", zap.Any("resp", resp.String()), zap.Any("url", url), zap.Any("body", body), zap.Any("headers", headers))
		return nil, internal.SysStatus.SystemBusy
	}
	respLocal := &models.RespLocal{}
	err = json.Unmarshal(resp.Body(), &respLocal)
	if err != nil {
		internal.Log.Error("CallSendMailByTemplate Unmarshal", zap.Error(err), zap.Any("resp", resp.String()), zap.Any("url", url), zap.Any("body", body), zap.Any("headers", headers))
		return nil, internal.SysStatus.SystemError
	}
	if respLocal.StatusCode != 0 {
		return nil, &internal.SystemStatus{
			Status: respLocal.StatusCode,
			Msg:    respLocal.Message,
		}
	}
	return respLocal, nil
}
func CallGetCustomerFrPhone(phoneNb string, funcName string) (*models.APICustomerInfo, *internal.SystemStatus) {
	url := internal.Domains.HiFPT + internal.Eps.HiCustomerProvider.CustomerInfo
	headers := map[string]string{
		"clientKey":     internal.Keys.NotifyProviderClientKey,
		"Authorization": utils.CreateNotifyTemplateToken(),
	}
	body := map[string]interface{}{
		"phone": phoneNb,
	}
	internal.Log.Info("Call", zap.Any("funcName", funcName), zap.Any("url", url), zap.Any("headers", headers), zap.Any("body", body))
	resp, err := utils.Request(url, false, headers, nil, body, 10, false)
	if err != nil || resp.StatusCode() != 200 {
		internal.Log.Error("Call fail", zap.Any("funcName", funcName), zap.Any("url", url), zap.Any("headers", headers), zap.Any("body", body), zap.Error(err))
		return nil, internal.SysStatus.SystemBusy
	}
	internal.Log.Info("CallGetCustomerFrPhone Response", zap.Any("funcName", funcName), zap.Any("url", url), zap.Any("headers", headers), zap.Any("body", body), zap.Any("resp", string(resp.Body())))
	res := &struct {
		StatusCode int    `json:"statusCode"`
		Message    string `json:"message"`
		Data       struct {
			Info models.APICustomerInfo `json:"info"`
		} `json:"data"`
	}{}
	err = json.Unmarshal(resp.Body(), &res)
	if err != nil {
		internal.Log.Error("json.Unmarshal", zap.Any("funcName", funcName), zap.Any("input", string(resp.Body())), zap.Error(err))
		return nil, internal.SysStatus.SystemError
	}
	if res.StatusCode == 0 {
		return &res.Data.Info, nil
	} else {
		return nil, internal.SysStatus.SystemBusy
	}
}

func CallSendFormatCustomerSop(customer *models.CustomerSopRequest, funcName string) (*models.CustomerSopResponse, *internal.SystemStatus) {
	url := internal.Domains.CustomerSop + internal.Eps.CustomerSop.SyncCustomer
	headers := map[string]string{
		"Authorization": utils.CreateCustomerSopToken(),
	}

	body := map[string]interface{}{}
	byteInput, _ := json.Marshal(customer)
	err := json.Unmarshal(byteInput, &body)
	if err != nil {
		internal.Log.Error("json.Unmarshal", zap.Any("funcName", funcName), zap.Any("input", string(byteInput)), zap.Error(err))
		return nil, internal.SysStatus.SystemError
	}
	internal.Log.Info("Call", zap.Any("funcName", funcName), zap.Any("url", url), zap.Any("headers", headers), zap.Any("body", body))
	resp, err := utils.Request(url, false, headers, nil, body, 10, false)

	if err != nil {
		internal.Log.Error("Call fail", zap.Any("funcName", funcName), zap.Any("url", url), zap.Any("headers", headers), zap.Any("body", body), zap.Error(err))
		return nil, internal.SysStatus.SystemBusy
	}
	internal.Log.Info("Response", zap.Any("funcName", funcName), zap.Any("url", url), zap.Any("headers", headers), zap.Any("body", body), zap.Any("httpStatus", resp.Status()), zap.Any("dt", fmt.Sprintf("%v", resp.Time().Seconds())), zap.Any("resp", resp.String()))
	if resp.StatusCode() != 200 {
		internal.Log.Error("Error Response", zap.Any("resp", resp.String()), zap.Any("funcName", funcName), zap.Any("url", url), zap.Any("body", body), zap.Any("headers", headers))
		return nil, internal.SysStatus.SystemBusy
	}
	res := &models.CustomerSopResponse{}

	err = json.Unmarshal(resp.Body(), &res)
	if err != nil {
		internal.Log.Error("json.Unmarshal", zap.Any("funcName", funcName), zap.Any("input", string(resp.Body())), zap.Error(err))
		return nil, internal.SysStatus.SystemError
	}
	if res.Status != "success" {
		return res, internal.SysStatus.SystemError
	} else {
		return res, nil
	}
}

func CallGetCustomerAddress(tokenApp string, funcName string) (*models.CustomerLocationResponse, *internal.SystemStatus) {
	url := internal.Domains.HiFPT + internal.Eps.HiFPTApi.GetShopAddress
	headers := map[string]string{
		"Authorization": tokenApp,
	}
	body := map[string]interface{}{}

	internal.Log.Info("Call GetShopAddress", zap.Any("funcName", funcName), zap.Any("url", url), zap.Any("headers", headers), zap.Any("body", body))
	resp, err := utils.Request(url, false, headers, nil, body, 10, false)

	if err != nil {
		internal.Log.Error("Call fail", zap.Any("funcName", funcName), zap.Any("url", url), zap.Any("headers", headers), zap.Any("body", body), zap.Error(err))
		return nil, internal.SysStatus.SystemBusy
	}
	internal.Log.Info("Response Call GetShopAddress", zap.Any("funcName", funcName), zap.Any("url", url), zap.Any("headers", headers), zap.Any("body", body), zap.Any("httpStatus", resp.Status()), zap.Any("dt", fmt.Sprintf("%v", resp.Time().Seconds())), zap.Any("resp", resp.String()))
	if resp.StatusCode() != 200 {
		internal.Log.Error("Error GetShopAddress Response", zap.Any("resp", resp.String()), zap.Any("funcName", funcName), zap.Any("url", url), zap.Any("body", body), zap.Any("headers", headers))
		return nil, internal.SysStatus.SystemBusy
	}
	res := &models.CustomerLocationResponse{}
	err = json.Unmarshal(resp.Body(), &res)
	if err != nil {
		internal.Log.Error("json.Unmarshal", zap.Any("funcName", funcName), zap.Any("input", string(resp.Body())), zap.Error(err))
		return nil, internal.SysStatus.SystemError
	}
	if res.StatusCode != 0 {
		return nil, internal.SysStatus.SystemError
	}
	return res, nil
}

func CallCustomerWifi6SyncToSOP(bodyOrigin interface{}, funcName string) (*models.CustomerSopResponse, *internal.SystemStatus) {
	url := internal.Domains.CustomerSop + internal.Eps.CustomerSop.SyncCustomerWifi6
	headers := map[string]string{
		"Authorization": utils.CreateCustomerSopToken(),
	}

	body := map[string]interface{}{}
	byteInput, _ := json.Marshal(bodyOrigin)
	err := json.Unmarshal(byteInput, &body)
	if err != nil {
		internal.Log.Error("json.Unmarshal", zap.Any("funcName", funcName), zap.Any("input", string(byteInput)), zap.Error(err))
		return nil, internal.SysStatus.SystemError
	}
	internal.Log.Info("Call", zap.Any("funcName", funcName), zap.Any("url", url), zap.Any("headers", headers), zap.Any("body", body))
	resp, err := utils.Request(url, false, headers, nil, body, 10, false)

	if err != nil {
		internal.Log.Error("Call fail", zap.Any("funcName", funcName), zap.Any("url", url), zap.Any("headers", headers), zap.Any("body", body), zap.Error(err))
		return nil, internal.SysStatus.SystemBusy
	}
	internal.Log.Info("Response", zap.Any("funcName", funcName), zap.Any("url", url), zap.Any("headers", headers), zap.Any("body", body), zap.Any("httpStatus", resp.Status()), zap.Any("dt", fmt.Sprintf("%v", resp.Time().Seconds())), zap.Any("resp", resp.String()))
	if resp.StatusCode() != 200 {
		internal.Log.Error("Error Response", zap.Any("resp", resp.String()), zap.Any("funcName", funcName), zap.Any("url", url), zap.Any("body", body), zap.Any("headers", headers))
		return nil, internal.SysStatus.SystemBusy
	}
	res := &models.CustomerSopResponse{}
	err = json.Unmarshal(resp.Body(), &res)
	if err != nil {
		internal.Log.Error("json.Unmarshal", zap.Any("funcName", funcName), zap.Any("input", string(resp.Body())), zap.Error(err))
		return nil, internal.SysStatus.SystemError
	}
	return res, nil
}

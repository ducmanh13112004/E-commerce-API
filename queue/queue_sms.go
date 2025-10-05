package queue

// type QueueSMS struct {
// 	Body          map[string]interface{}
// 	URL           string
// 	URLLogin      string
// 	Header        map[string]string
// 	Repo          *repositories.Repositories
// 	Input         models.DBScheduleRemind
// 	CallBack      models.CallBackRemind
// 	UseProduction bool
// }

// func (c *QueueSMS) Run() {
// 	//check time call notify
// 	log := NewLogger()
// 	now := utils.GetTimeUTC7()
// 	dataCallback := map[string]interface{}{}
// 	dataCallback["noti_id"] = c.Input.NotiId
// 	dataCallback["order_id"] = c.Input.OrderId
// 	dataCallback["remind_type"] = c.Input.Type
// 	dataCallback["t_remind"] = utils.GetStringTime(now, "Y-M-D H:M:S")
// 	result := models.ResultCallBackRemind{}
// 	if now == c.Input.StartSchedule || now == c.Input.EndSchedule || (c.Input.StartSchedule.Before(now) && c.Input.EndSchedule.After(now)) {
// 		controlApi, err := c.Repo.ControlApi.GetInfoFrApiName("TEST_REMIND_SMS")
// 		if err != nil {
// 			log.Error("GetInfoFrApiName('TEST_REMIND_SMS')", zap.Error(err))
// 			controlApi = &models.ControlApiTb{Conditions: "0984990410;0968893364"}
// 		}
// 		if controlApi == nil {
// 			log.Error("TEST_REMIND_SMS nil")
// 			controlApi = &models.ControlApiTb{Conditions: "0984990410;0968893364"}
// 		}
// 		check := false
// 		phoneNb, ok := c.Body["PhoneNumber"].(string)
// 		if ok {
// 			listWhiteList := utils.ConvertStringToList(controlApi.Conditions, ";")
// 			for _, value := range listWhiteList {
// 				if value == phoneNb {
// 					check = true
// 				}
// 			}
// 		}
// 		if (c.UseProduction == false && check == true) || c.UseProduction {
// 			//login sms
// 			accessToken, err := c.Login()
// 			if err != nil {
// 				log.Error("Login failed", zap.Error(err))
// 				result.StatusCode = 300
// 				result.Message = "Login failed:" + err.Error()
// 				dataCallback["response"] = result
// 				c.CallBack.Body = map[string]interface{}{
// 					"type_request": "result_remind",
// 					"data": map[string]interface{}{
// 						"list_remind": []interface{}{dataCallback},
// 					},
// 				}
// 				c.APICallBack()
// 				return
// 			}
// 			strDataSend, _ := json.Marshal(c.Body)
// 			endDataSend := utils.EncodeBase64(string(strDataSend))
// 			secretKey := utils.EncodeBase64("kpdv")
// 			sendUrl := c.URL
// 			bodySendSMS := map[string]interface{}{
// 				"data": endDataSend + "@" + secretKey,
// 			}
// 			bodyByte, _ := json.Marshal(bodySendSMS)
// 			storeCall := &models.StoreCallApiTb{
// 				Url:   c.URL,
// 				Input: string(bodyByte),
// 			}
// 			respSend, errSend := utils.Request(sendUrl, false, map[string]string{
// 				"tokenSMS": accessToken,
// 			}, nil, map[string]interface{}{
// 				"data": endDataSend + "@" + secretKey,
// 			}, 3, false)
// 			if errSend != nil {
// 				log.Error("Send sms fail", zap.Any("data send", endDataSend), zap.Any("secretKey", secretKey), zap.Error(errSend))
// 				result.StatusCode = 300
// 				result.Message = "SMS failed:" + err.Error()
// 				storeCall.Output = err.Error()
// 				err = c.Repo.StoreCallApi.Create(storeCall)
// 				if err != nil {
// 					log.Error("repo.StoreCallApi.Create", zap.Any("data", storeCall), zap.Error(err))
// 				}
// 				dataCallback["response"] = result
// 				c.CallBack.Body = map[string]interface{}{
// 					"type_request": "result_remind",
// 					"data": map[string]interface{}{
// 						"list_remind": []interface{}{dataCallback},
// 					},
// 				}
// 				c.APICallBack()
// 				return
// 			}
// 			resp := &struct {
// 				ID      int         `json:"ID"`
// 				Message string      `json:"Message"`
// 				Detail  interface{} `json:"Detail"`
// 			}{}
// 			log.Info("Send SMS response", zap.Any("url", sendUrl), zap.Any("data send", endDataSend), zap.Any("secretKey", secretKey), zap.Any("http status code", respSend.StatusCode()), zap.Any("check nil", respSend != nil), zap.Any("Resp", respSend.String()))
// 			if respSend != nil && respSend.StatusCode() == 200 {
// 				storeCall.Output = respSend.String()
// 				err = c.Repo.StoreCallApi.Create(storeCall)
// 				if err != nil {
// 					log.Error("repo.StoreCallApi.Create", zap.Any("data", storeCall), zap.Error(err))
// 				}
// 				json.Unmarshal([]byte(respSend.String()), &resp)
// 				result.Data = resp.Detail
// 				result.Message = resp.Message
// 				result.StatusCode = resp.ID
// 				dataCallback["response"] = result
// 				c.CallBack.Body = map[string]interface{}{
// 					"type_request": "result_remind",
// 					"data": map[string]interface{}{
// 						"list_remind": []interface{}{dataCallback},
// 					},
// 				}
// 				c.APICallBack()
// 				return
// 			} else {
// 				log.Error("unknown error Send SMS")
// 				storeCall.Output = fmt.Sprintf("unknown error Send SMS:%v %v", respSend.StatusCode(), respSend.String())
// 				err = c.Repo.StoreCallApi.Create(storeCall)
// 				if err != nil {
// 					log.Error("repo.StoreCallApi.Create", zap.Any("data", storeCall), zap.Error(err))
// 				}
// 				result.Data = respSend.String()
// 				result.Message = "unknown error Send SMS"
// 				result.StatusCode = 300
// 				dataCallback["response"] = result
// 				c.CallBack.Body = map[string]interface{}{
// 					"type_request": "result_remind",
// 					"data": map[string]interface{}{
// 						"list_remind": []interface{}{dataCallback},
// 					},
// 				}
// 				c.APICallBack()
// 			}
// 		}
// 	} else {
// 		log.Info("Ngoài giờ", zap.Any("data", c.Input))
// 		result.Data = nil
// 		result.Message = "Ngoài giờ"
// 		result.StatusCode = 400
// 		dataCallback["response"] = result
// 		c.CallBack.Body = map[string]interface{}{
// 			"type_request": "result_remind",
// 			"data": map[string]interface{}{
// 				"list_remind": []interface{}{dataCallback},
// 			},
// 		}
// 		c.APICallBack()
// 	}
// }
// func (c *QueueSMS) Login() (string, error) {
// 	log := NewLogger()
// 	body := map[string]interface{}{
// 		"UserName": "kpdv",
// 		"Password": "xbGokOKN7dr6jPIdjfRsqd2XpbN8stedGoQ/9Ucann4=",
// 	}
// 	byteBody, _ := json.Marshal(body)
// 	storeCall := &models.StoreCallApiTb{
// 		Url:   c.URLLogin,
// 		Input: string(byteBody),
// 	}
// 	respLogin, errLogin := utils.Request(c.URLLogin, false, nil, nil, body, 3, false)
// 	if errLogin != nil {
// 		storeCall.Output = errLogin.Error()
// 		err := c.Repo.StoreCallApi.Create(storeCall)
// 		if err != nil {
// 			log.Error("repo.StoreCallApi.Create", zap.Any("data", storeCall), zap.Error(err))
// 		}
// 		return "", errLogin
// 	}
// 	log.Info("Resp SMS login", zap.Any("url", c.URL), zap.Any("body", body), zap.Any("respLogin", respLogin))
// 	smsLoginResp := &struct {
// 		Detail map[string]interface{}
// 	}{}
// 	if errLogin == nil && respLogin.StatusCode() == 200 {
// 		errMarshal := json.Unmarshal([]byte(respLogin.String()), &smsLoginResp)
// 		if errMarshal != nil {
// 			return "", errMarshal
// 		}
// 		accessToken := smsLoginResp.Detail["AccessToken"].(string)
// 		storeCall.Output = respLogin.String()
// 		err := c.Repo.StoreCallApi.Create(storeCall)
// 		if err != nil {
// 			log.Error("repo.StoreCallApi.Create", zap.Any("data", storeCall), zap.Error(err))
// 		}
// 		return accessToken, nil
// 	}
// 	return "", errLogin
// }
// func (c *QueueSMS) APICallBack() {
// 	log := NewLogger()
// 	param := map[string]string{}
// 	bodyByte, err := json.Marshal(c.CallBack.Body)
// 	if err != nil {
// 		log.Error("json.Marshal(c.CallBack.Body)", zap.Any("data", c.CallBack.Body), zap.Error(err))
// 		return
// 	}
// 	responese, err := utils.Request(c.CallBack.URL, false, c.CallBack.Header, param, c.CallBack.Body, 10, false)
// 	log.Info("Request APICallBack response", zap.Any("URL", c.URL), zap.Any("header", c.Header), zap.Any("body", c.Body), zap.Any("resp", responese.String()))
// 	storeCall := &models.StoreCallApiTb{
// 		Url:     c.CallBack.URL,
// 		Input:   string(bodyByte),
// 		TCreate: utils.GetTimeUTC7(),
// 	}
// 	if err != nil {
// 		log.Error("Request APICallBack fail", zap.Any("URL", c.URL), zap.Any("header", c.Header), zap.Any("body", c.Body), zap.Error(err))
// 		storeCall.Output = err.Error()
// 		err = c.Repo.StoreCallApi.Create(storeCall)
// 		if err != nil {
// 			log.Error("StoreCallApi.Create(storeCall)", zap.Any("data", storeCall), zap.Error(err))
// 		}
// 		return
// 	}
// 	storeCall.Output = responese.String()
// 	err = c.Repo.StoreCallApi.Create(storeCall)
// 	if err != nil {
// 		log.Error("StoreCallApi.Create(storeCall)", zap.Any("data", storeCall), zap.Error(err))
// 	}
// }

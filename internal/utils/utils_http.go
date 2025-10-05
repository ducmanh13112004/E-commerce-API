package utils

import (
	"ecom_promotion_v2/internal"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

func Request(url string, isGet bool, headers, param map[string]string, body map[string]interface{}, timeout int, proxy bool) (*resty.Response, error) {
	if headers == nil {
		headers = map[string]string{}
	}
	if param == nil {
		param = map[string]string{}
	}
	if body == nil {
		body = map[string]interface{}{}
	}
	_, ok := headers["Content-Type"]
	if !ok {
		headers["Content-Type"] = "application/json"
	}
	client := resty.New()
	client.SetTimeout(time.Second * time.Duration(timeout))
	if proxy {
		client.SetProxy("http://proxy.hcm.fpt.vn:80")
	}
	req := client.R().SetHeaders(headers).SetQueryParams(param).
		SetBody(body)
	if isGet {
		return req.Get(url)
	}
	return req.Post(url)
}

func RequestWithAuth(url string, isGet bool, headers, param map[string]string, body map[string]interface{}, timeout int, proxy bool, authToken string) (*resty.Response, error) {
	headers["Content-Type"] = "application/json"
	client := resty.New()
	client.SetTimeout(time.Second * time.Duration(timeout))
	if proxy {
		client.SetProxy("http://proxy.hcm.fpt.vn:80")
	}
	req := client.R().
		SetHeaders(headers).
		SetQueryParams(param).
		SetBody(body)
	if authToken != "" {
		req = req.SetAuthToken(authToken)
	}
	if isGet {
		return req.Get(url)
	}
	return req.Post(url)
}
func RequestV2(url string, isGet bool, headers, param map[string]string, body map[string]interface{}, formData map[string]string, timeout int, proxy bool) (*resty.Response, error) {
	now := GetTimeUTC7()
	defer func() { fmt.Printf("ExecTime url=%s, dt=%v \n", url, GetTimeUTC7().Sub(now).Milliseconds()) }()
	client := resty.New()
	client.SetTimeout(time.Second * time.Duration(timeout))
	if proxy {
		client.SetProxy("http://proxy.hcm.fpt.vn:80")
	}
	req := client.R().SetFormData(formData)
	if isGet {
		return req.Get(url)
	}
	return req.Post(url)
}

func RequestV3(url string, isGet bool, headers, param map[string]string, body map[string]interface{}, formData map[string]string, timeout int, proxy bool) (*resty.Response, error) {
	now := GetTimeUTC7()
	defer func() {
		fmt.Printf("ExecTime url=%s, dt=%v ms\n", url, GetTimeUTC7().Sub(now).Milliseconds())
	}()

	client := resty.New()
	client.SetTimeout(time.Second * time.Duration(timeout))

	if proxy {
		client.SetProxy("http://proxy.hcm.fpt.vn:80")
	}

	// === Log trước khi call ===
	fmt.Println("=== REQUEST INFO ===")
	fmt.Println("Method:", map[bool]string{true: "GET", false: "POST"}[isGet])
	fmt.Println("URL:", url)
	if len(headers) > 0 {
		fmt.Println("Headers:", headers)
	}
	if len(param) > 0 {
		fmt.Println("Params:", param)
	}
	if len(body) > 0 {
		fmt.Println("Body:", body)
	}
	if len(formData) > 0 {
		fmt.Println("FormData:", formData)
	}

	req := client.R()

	if len(headers) > 0 {
		req.SetHeaders(headers)
	}
	if len(param) > 0 {
		req.SetQueryParams(param)
	}
	if len(body) > 0 {
		req.SetBody(body)
	}
	if len(formData) > 0 {
		req.SetFormData(formData)
	}

	var (
		resp *resty.Response
		err  error
	)
	if isGet {
		resp, err = req.Get(url)
	} else {
		resp, err = req.Post(url)
	}

	// === Log sau khi call ===
	if err != nil {
		internal.Log.Error("ResponseError-RequestV3", zap.Any("url", url), zap.Any("headers", headers), zap.Any("body", body), zap.Any("param", param), zap.Any("formData", formData), zap.Error(err))
	} else {
		internal.Log.Info("ResponseInfo-RequestV3", zap.Any("url", url), zap.Any("headers", headers), zap.Any("body", body), zap.Any("param", param), zap.Any("formData", formData), zap.Any("response", resp.String()), zap.Any("dt", resp.Time().Seconds()), zap.Any("http_status", resp.StatusCode()))
	}

	return resp, err
}

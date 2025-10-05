package cmd

import (
	"ecom_promotion_v2/internal"
	"ecom_promotion_v2/internal/utils"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func TestAPI(ctx *fiber.Ctx) error {
	funcName := "TestAPI"
	type StructAPI struct {
		Url     string
		Header  map[string]string
		Body    map[string]interface{}
		IsGet   bool
		IsProxy bool
		Param   map[string]string

		TypeCall string
	}
	body := &StructAPI{}
	err := ctx.BodyParser(body)
	if err != nil {
		return ctx.JSON(err)
	}
	internal.Log.Info("Call", zap.Any("funcName", funcName), zap.Any("url", body.Url), zap.Any("header", body.Header), zap.Any("Body", body.Body), zap.Any("Param", body.Param))
	resp, err := utils.Request(body.Url, body.IsGet, body.Header, body.Param, body.Body, 30, body.IsProxy)
	if err != nil {
		internal.Log.Error("Response", zap.Any("funcName", funcName), zap.Any("url", body.Url), zap.Any("header", body.Header), zap.Any("Body", body.Body), zap.Any("Param", body.Param), zap.Error(err))
		return ctx.JSON(err)
	}
	internal.Log.Info("Response", zap.Any("funcName", funcName), zap.Any("url", body.Url), zap.Any("header", body.Header), zap.Any("Body", body.Body), zap.Any("Param", body.Param), zap.Any("resp", resp.String()), zap.Any("httpStatus", resp.Status()))
	res := map[string]interface{}{}
	err = json.Unmarshal(resp.Body(), &res)
	if err != nil {
		internal.Log.Error("Unmarshal", zap.Any("resp", resp.String()), zap.Error(err))
		return ctx.JSON(resp.String())
	}
	return ctx.JSON(res)

}

package delivery

import (
	"ecom_promotion_v2/internal"
	"ecom_promotion_v2/internal/models"
	"ecom_promotion_v2/internal/utils"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/golang-jwt/jwt"
	"github.com/techoner/gophp"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

var visitors = make(map[string]*rate.Limiter)
var mu sync.Mutex

func getVisitor(ip string, limit int) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(1, limit)
		visitors[ip] = limiter
	}

	return limiter
}
func (s *appHandlers) RateLimit(ctx *fiber.Ctx) error {
	// customer_id := ctx.Locals("customer_id").(string)
	// limiter := getVisitor(customer_id, 10)
	// if !limiter.Allow() {
	// 	internal.Log.Info("Rate limit", zap.Any("Customer_id", customer_id))
	// 	return ctx.Status(fiber.StatusOK).JSON(models.RespWeb{
	// 		Status: internal.SysStatus.RateLimit.Status,
	// 		Msg:    internal.SysStatus.RateLimit.Msg,
	// 	})
	// }
	return ctx.Next()
}

func (s *appHandlers) CustomRateLimit(ctx *fiber.Ctx) error {
	customer_id := ctx.Locals("customer_id").(string)
	limiter := getVisitor(customer_id, 30)
	if !limiter.Allow() {
		internal.Log.Info("Rate limit", zap.Any("Customer_id", customer_id))
		return ctx.Status(fiber.StatusOK).JSON(models.RespWeb{
			Status: internal.SysStatus.RateLimit.Status,
			Msg:    internal.SysStatus.RateLimit.Msg,
		})
	}
	return ctx.Next()
}
func (s *appHandlers) FiberRateLimit(maxRequests int, duration time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        maxRequests,
		Expiration: duration,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusOK).JSON(models.RespWeb{
				Status: internal.SysStatus.RateLimit.Status,
				Msg:    internal.SysStatus.RateLimit.Msg,
			})
		},
		KeyGenerator: func(c *fiber.Ctx) string {
			customerID := c.Locals("customer_id").(string)

			if customerID == "" {
				// Fallback to IP or reject
				customerID = c.IP()
			}
			return "rate_limit_customer_" + customerID
		},
	})
}
func (s *appHandlers) RequireTokenWeb(ctx *fiber.Ctx) error {
	token := string(ctx.Request().Header.Peek("TOKEN"))
	ip := string(ctx.Request().Header.Peek("X-Real-Ip"))
	clams := jwt.MapClaims{}
	resultFe := models.RespWeb{}
	defer func() {
		internal.Log.Info("RequireTokenWeb", zap.Any("ip", ip), zap.Any("url", ctx.Context().URI()), zap.Any("authen", token), zap.Any("jwt", clams), zap.Any("resultFe", resultFe))
	}()
	if utils.IsEmpty(token) {
		resultFe = models.RespWeb{
			Status: internal.SysStatus.TokenRequired.Status,
			Msg:    internal.SysStatus.TokenRequired.Msg,
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(resultFe)
	}
	decodeToken, err := jwt.ParseWithClaims(token, clams, func(t *jwt.Token) (interface{}, error) {
		return []byte(internal.Keys.TOKEN_SECRET_KEY_WEB), nil
	})
	if err != nil {
		jwtErr := err.(*jwt.ValidationError).Errors
		if jwtErr == jwt.ValidationErrorExpired {
			resultFe = models.RespWeb{
				Status: internal.SysStatus.TokenExpired.Status,
				Msg:    internal.SysStatus.TokenExpired.Msg,
			}
			return ctx.Status(fiber.StatusBadRequest).JSON(resultFe)
		}
	}

	if decodeToken == nil || !decodeToken.Valid {
		resultFe = models.RespWeb{
			Status: internal.SysStatus.InvalidToken.Status,
			Msg:    internal.SysStatus.InvalidToken.Msg,
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(resultFe)
	}

	ctx.Locals("token_webkit", token)
	ctx.Locals("customer_id", clams["customerId"])
	ctx.Locals("customer_phone", clams["phone"])
	ctx.Locals("app_version", clams["appVersion"])
	ctx.Locals("access_token", clams["accessToken"])
	ctx.Locals("customer_name", clams["fullname"])
	ctx.Locals("customer_ip", ip)
	// ctx.Locals("user_info", userInfo)

	return ctx.Next()
}
func (s *appHandlers) RequireTokenLocal(ctx *fiber.Ctx) error {
	token := string(ctx.Request().Header.Peek("TOKEN"))
	resultFe := models.RespLocal{}
	defer func() {
		internal.Log.Info("RequireTokenLocal", zap.Any("url", ctx.Context().URI()), zap.Any("authen", token), zap.Any("result", resultFe))
	}()
	if utils.IsEmpty(token) {
		resultFe = models.RespLocal{
			StatusCode: internal.SysStatus.TokenRequired.Status,
			Message:    internal.SysStatus.TokenRequired.Msg,
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(resultFe)
	}
	getToken := utils.CreateEcomToken(internal.Keys.HIFPT_ECOM_CLIENT_KEY, internal.Keys.HIFPT_ECOM_SECRET_KEY)

	if token != getToken {
		detail := map[string]interface{}{}
		if !internal.Envs.IsProduction {
			detail["token"] = getToken
		}
		resultFe = models.RespLocal{
			StatusCode: internal.SysStatus.InvalidToken.Status,
			Message:    internal.SysStatus.InvalidToken.Msg,
			Data:       detail,
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(resultFe)
	}
	return ctx.Next()
}

func (s *appHandlers) RequireTokenClient(ctx *fiber.Ctx) error {
	token := string(ctx.Request().Header.Peek("TOKEN"))
	clientID := string(ctx.Request().Header.Peek("CLIENT-ID"))
	var body interface{}
	ctx.BodyParser(&body)
	internal.Log.Info("RequireTokenClient", zap.Any("ip", ctx.IP()), zap.Any("url", ctx.Context().URI()), zap.Any("authen", token))
	if utils.IsEmpty(token) {
		return ctx.Status(fiber.StatusBadRequest).JSON(models.RespLocal{
			StatusCode: internal.SysStatus.TokenRequired.Status,
			Message:    internal.SysStatus.TokenRequired.Msg,
		})
	}
	//check clientId valid
	secretKey, ok := internal.Keys.ClientIdSercretKey[clientID]
	if !ok {
		internal.Log.Error("ClientId is not valid", zap.Any("clientId", clientID), zap.Any("clientData", internal.Keys.ClientIdSercretKey))
		return ctx.Status(fiber.StatusBadRequest).JSON(models.RespLocal{
			StatusCode: 400,
			Message:    "ClientId không hợp lệ.",
		})
	}
	getToken := utils.CreateEcomToken(clientID, secretKey)
	if token != getToken {
		detail := map[string]interface{}{}
		if !internal.Envs.IsProduction {
			detail["token"] = getToken
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(models.RespLocal{
			StatusCode: internal.SysStatus.InvalidToken.Status,
			Message:    internal.SysStatus.InvalidToken.Msg,
			Data:       detail,
		})
	}
	return ctx.Next()
}

func (s *appHandlers) RequireTokenPortal(ctx *fiber.Ctx) error {
	funcName := "RequireTokenPortal"
	token := string(ctx.Request().Header.Peek("TOKEN"))
	clams := jwt.MapClaims{}
	resultFe := models.RespWeb{
		Status: internal.CODE_TOKEN_REQUIRED,
		Msg:    internal.MSG_TOKEN_REQUIRED,
	}
	var body interface{}
	ctx.BodyParser(&body)
	internal.Log.Info(funcName, zap.Any("ip", ctx.IP()), zap.Any("url", ctx.Context().URI()), zap.Any("authen", token))
	if utils.IsEmpty(token) {
		return ctx.Status(fiber.StatusBadRequest).JSON(resultFe)
	}
	decodeToken, err := jwt.ParseWithClaims(token, clams, func(t *jwt.Token) (interface{}, error) {
		return []byte(internal.Keys.TOKEN_SECRET_KEY_PORTAL), nil
	})

	if err != nil {
		jwtErr := err.(*jwt.ValidationError).Errors
		if jwtErr == jwt.ValidationErrorExpired {
			resultFe = models.RespWeb{
				Status: internal.CODE_TOKEN_EXPIRED,
				Msg:    internal.MSG_TOKEN_EXPIRED,
			}
			return ctx.Status(fiber.StatusOK).JSON(resultFe)
		}
	}

	if decodeToken == nil || !decodeToken.Valid {
		resultFe = models.RespWeb{
			Status: internal.CODE_INVALID_TOKEN,
			Msg:    internal.MSG_INVALID_TOKEN,
		}
		return ctx.Status(fiber.StatusOK).JSON(resultFe)
	}
	payload := models.PortalPayload{}
	bytePortal, _ := json.Marshal(clams)
	err = json.Unmarshal(bytePortal, &payload)
	fmt.Println("payload:", payload)
	if err != nil {
		internal.Log.Error("Unmarshal", zap.Any("funcName", funcName), zap.Error(err))
	}

	ctx.Locals("user_info", payload)
	ctx.Locals("payload", payload)
	ctx.Locals("user_id", clams["userId"])
	ctx.Locals("user_code", clams["userCode"])
	ctx.Locals("phone_nb", clams["phoneNb"])
	ctx.Locals("fullname", clams["fullname"])
	ctx.Locals("email", clams["email"])
	ctx.Locals("update_by", clams["updateBy"])
	ctx.Locals("company", clams["company"])
	ctx.Locals("company_name", clams["companyName"])
	ctx.Locals("department", clams["department"])
	ctx.Locals("department_name", clams["departmentName"])
	ctx.Locals("emp_code", clams["empCode"])
	ctx.Locals("emp_id", clams["empId"])
	ctx.Locals("group_id", clams["groupId"])
	return ctx.Next()
}

func (s *appHandlers) RequireXKeyTelegramPortal(ctx *fiber.Ctx) error {
	funcName := "RequireXKeyTelegramPortal"
	xApiKey := string(ctx.Request().Header.Peek("X-API-KEY"))

	internal.Log.Info(funcName, zap.Any("ip", ctx.IP()), zap.Any("url", ctx.Context().URI()), zap.Any("authen", xApiKey))

	// Kiểm tra nếu API Key bị thiếu
	if utils.IsEmpty(xApiKey) {
		internal.Log.Error(funcName, zap.String("error", "Missing API Key"), zap.String("client_ip", ctx.IP()))
		resultFe := models.RespWeb{
			Status: internal.CODE_TOKEN_REQUIRED,
			Msg:    "API Key is required",
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(resultFe)
	}

	// Kiểm tra nếu API Key không hợp lệ
	if xApiKey != internal.Keys.X_API_KEY_TELEGRAM {
		internal.Log.Error(funcName, zap.String("error", "Invalid API Key"), zap.String("provided_key", xApiKey), zap.String("client_ip", ctx.IP()))
		resultFe := models.RespWeb{
			Status: internal.CODE_INVALID_TOKEN,
			Msg:    "Invalid API Key",
		}
		return ctx.Status(fiber.StatusUnauthorized).JSON(resultFe)
	}

	// internal.Log.Info(funcName, zap.String("message", "API Key validated successfully"), zap.String("client_ip", ctx.IP()))

	return ctx.Next()
}
func (s *appHandlers) RequireTokenApp(ctx *fiber.Ctx) error {
	token := string(ctx.Request().Header.Peek("Authorization"))
	var body interface{}
	ctx.BodyParser(&body)
	resultApp := models.RespLocal{
		StatusCode: internal.SysStatus.TokenAppError.Status,
		Message:    internal.SysStatus.TokenAppError.Msg,
		Data:       internal.SysStatus.TokenAppError.Detail,
	}
	internal.Log.Info("RequireTokenApp", zap.Any("ip", ctx.IP()), zap.Any("url", ctx.Context().URI()), zap.Any("authen", token), zap.Any("body", body))
	if utils.IsEmpty(token) {
		internal.Log.Error("Empty token")
		return ctx.Status(fiber.StatusOK).JSON(resultApp)
	}
	tokenStr := strings.Replace(token, "Bearer ", "", 1)
	_, err := s.rdb.Ping().Result()
	if err != nil {
		internal.Log.Error("Ping ", zap.Error(err), zap.Any("funcName", "RequireTokenApp"))
		return ctx.JSON(resultApp)
	}
	clientId, token, errCore := s.VerifyTokenApp(tokenStr, internal.Keys.TOKEN_SECRET_KEY_APP)
	if errCore != nil {
		internal.Log.Error("VerifyTokenApp Error", zap.Any("tokenStr", tokenStr), zap.Error(err), zap.Any("funcName", "RequireTokenApp"))
		return ctx.JSON(resultApp)
	}
	keyRedis := fmt.Sprintf("%s:%s", clientId, token)
	dataRedis, _, errCore := s.GetDataRedisFrKey(keyRedis)
	if errCore != nil {
		internal.Log.Error("GetDataRedisFrKey Error", zap.Any("funcName", "RequireTokenApp"))
		return ctx.JSON(resultApp)
	}
	// dataRedis := &models.RedisModel{
	// 	CustomerId: 123456,
	// 	Phone:      "0982777935",
	// 	AppVersion: "9.2",
	// }
	// tokenStr := "tokenApp"
	// keyRedis := tokenStr
	ctx.Locals("keyLimit", fmt.Sprintf("%v", dataRedis.CustomerId))
	ctx.Locals("customer_id", fmt.Sprintf("%v", dataRedis.CustomerId))
	ctx.Locals("customer_phone", fmt.Sprintf("%v", dataRedis.Phone))
	ctx.Locals("app_version", dataRedis.AppVersion)
	ctx.Locals("token_app", tokenStr)
	internal.Log.Info("RequireTokenApp", zap.Any("keyRedis", keyRedis), zap.Any("dataRedis", dataRedis))
	return ctx.Next()
}
func (a *appHandlers) VerifyTokenApp(tokenApp string, secretKey string) (string, string, *internal.SystemStatus) {
	claims := jwt.MapClaims{}
	decodeToken, err := jwt.ParseWithClaims(tokenApp, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		jwtErr := err.(*jwt.ValidationError).Errors
		if jwtErr == jwt.ValidationErrorExpired {
			internal.Log.Error("ValidationErrorExpired ", zap.Any("tokenParse", decodeToken), zap.Error(err))
			return "", "", internal.SysStatus.TokenExpired
		}
	}
	if decodeToken == nil {
		internal.Log.Error("Token Valid", zap.Any("tokenParse", decodeToken))
		return "", "", internal.SysStatus.InvalidToken
	}
	token, ok := claims["jti"].(string)
	if !ok {
		internal.Log.Error("claims: jti does not exit", zap.Any("data", claims))
		return "", "", internal.SysStatus.InvalidToken
	}
	clientId, ok := claims["clientId"].(string)
	if !ok {
		clientId = "session"
	}
	return clientId, token, nil
}
func (s *appHandlers) GetDataRedisFrKey(key string) (*models.RedisModel, string, *internal.SystemStatus) {
	if !internal.Envs.IsDev {
		result := &models.RedisModel{}
		val2, err := s.rdb.Get(key).Result()
		if err == redis.Nil {
			internal.Log.Info("Token does not exist", zap.Any("Token", key))
			return nil, "", &internal.SystemStatus{
				Status: 1004,
				Msg:    "Token không tồn tại",
			}
		} else if err != nil {
			internal.Log.Error("Redis error", zap.Error(err))
			return nil, "", internal.SysStatus.SystemError
		} else {
			dat := []byte(val2)
			out, err := gophp.Unserialize(dat)
			if err != nil {
				internal.Log.Error("gophp.Unserialize(dat)", zap.Any("data", dat), zap.Error(err))
				return nil, "", internal.SysStatus.SystemError
			}
			jsonbody, err := json.Marshal(out)
			if err != nil {
				internal.Log.Error("json.Marshal(out) ", zap.Any("data", out), zap.Error(err))
				return nil, "", internal.SysStatus.SystemError
			}
			if err := json.Unmarshal(jsonbody, &result); err != nil {
				internal.Log.Error("json.Unmarshal(jsonbody) ", zap.Any("data", jsonbody), zap.Error(err))
				return nil, "", internal.SysStatus.SystemError
			}
			return result, val2, nil
		}
	} else {
		// data test dev
		return &models.RedisModel{
			CustomerId: 2003344,
			Phone:      "0703970397",
			// ListContract: listContractTest,
		}, "", nil
	}
}

package healthcheckservice

import (
	"context"
	"ecom_promotion_v2/internal"
	"ecom_promotion_v2/internal/models"
	"ecom_promotion_v2/internal/repositories"
	"ecom_promotion_v2/internal/utils"
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var mapStatus = map[int]string{
	1: "normal",
	2: "warning",
	3: "critical",
	4: "no_use",
}
var statusDb string
var statusRedis string
var statusKafka string
var statusThird map[string]int
var funcName = "HealthCheck"
var uuidStr string

func (a *healthCheckService) HealthCheck(ctx *fiber.Ctx) (int, interface{}) {
	uuidStr = uuid.New().String()
	now := time.Now()
	internal.Log.Info(funcName, zap.String("url", ctx.Request().URI().String()), zap.Any("uuid", uuidStr))
	// TODO: Ping DB
	wg := new(sync.WaitGroup)
	wg.Add(3)
	statusDb = "normal"
	statusRedis = "normal"
	statusKafka = "normal"
	statusThird = map[string]int{}
	result := models.ResultHealthCheck{
		RequestID: uuidStr,
		Status:    "normal",
		Detail: models.DetailHealthCheck{
			CPUIndex:        "normal",
			DatabaseIndex:   "normal",
			KafkaIndex:      "normal",
			RAMIndex:        "normal",
			RedisIndex:      "normal",
			ThirdPartyIndex: "normal",
		},
	}
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	//check % RAM
	RamSys := 1500 // MB
	RamStat := float64(memStats.Sys / 1048576)
	ramPercent := RamStat / float64(RamSys)
	if ramPercent <= 0.75 {
		result.Detail.RAMIndex = "normal"
	} else if ramPercent <= 95 {
		result.Detail.RAMIndex = "warning"
	} else {
		result.Detail.RAMIndex = "critical"
	}

	// check CPU

	//Go Rountine check DB, Redis, Kafka, 3rd
	CheckThirdParty(wg)
	go CheckDb(wg, a.repo)
	go CheckRedis(wg)
	go CheckKafka(wg)
	wg.Wait()
	result.Detail.DatabaseIndex = statusDb
	result.Detail.RedisIndex = statusRedis
	result.Detail.KafkaIndex = statusKafka
	//status ThirdParty
	statusThirdInt := 1
	for _, value := range statusThird {
		if statusThirdInt < value {
			statusThirdInt = value
		}
	}
	result.Detail.ThirdPartyIndex = mapStatus[statusThirdInt]
	// Status final
	if result.Detail.RAMIndex == "critical" || result.Detail.CPUIndex == "critical" || result.Detail.DatabaseIndex == "critical" || result.Detail.KafkaIndex == "critical" || result.Detail.RedisIndex == "critical" || result.Detail.ThirdPartyIndex == "critical" {
		result.Status = "critical"
	} else if result.Detail.RAMIndex == "warning" || result.Detail.CPUIndex == "warning" || result.Detail.DatabaseIndex == "warning" || result.Detail.KafkaIndex == "warning" || result.Detail.RedisIndex == "warning" || result.Detail.ThirdPartyIndex == "warning" {
		result.Status = "warning"
	}
	internal.Log.Info(funcName, zap.Any("uuid", uuidStr), zap.Any("result", result), zap.Any("dt", fmt.Sprintf("%v", time.Since(now).Seconds())))
	return 200, result
}
func CheckDb(wg *sync.WaitGroup, repo *repositories.Repositories) {
	tStart := time.Now()
	status := 1 // normal
	timeLimit := 10.0
	defer func() {
		dt := time.Since(tStart).Seconds()
		if dt > timeLimit && status < 2 {
			status = 2
		}
		statusDb = mapStatus[status]
		internal.Log.Info("CheckDb", zap.Any("funcName", funcName), zap.Any("uuid", uuidStr), zap.Any("dt", fmt.Sprintf("%v", dt)), zap.Any("statusDb", statusDb))
		wg.Done() // This decreases counter by 1
	}()

	db, err := internal.CheckSQLDB()
	if err != nil {
		status = 3
		internal.Log.Error("CheckSQLDB", zap.Any("funcName", funcName), zap.Any("uuid", uuidStr), zap.Error(err))
		return
	}

	model := models.ControlApiTb{}
	err = db.Debug().Select("*").Table(model.TableName()).Where("api_name = ?", "HEALTH_CHECK").First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		internal.Log.Error("RecordNotFound", zap.Any("funcName", funcName), zap.Any("uuid", uuidStr), zap.Error(err))
	} else if err != nil {
		internal.Log.Error("Query FAIL", zap.Any("funcName", funcName), zap.Any("uuid", uuidStr), zap.Error(err))
		status = 2
	} else {
		status = model.Active
		timeInt, err := strconv.Atoi(model.Conditions)
		if err != nil {
			internal.Log.Error("strconv.Atoi", zap.Any("input", model.Conditions), zap.Error(err))
		} else {
			timeLimit = float64(timeInt)
		}
	}

	internal.Log.Info("CheckDb Lan 1", zap.Any("funcName", funcName), zap.Any("uuid", uuidStr), zap.Any("dt", fmt.Sprintf("%v", time.Since(tStart).Seconds())))
	tStart1 := time.Now()
	wgDb := new(sync.WaitGroup)
	wgDb.Add(1)
	go func(wgDb *sync.WaitGroup) {
		_, err := repo.ControlApi.GetInfoFrApiName("HEALTH_CHECK")
		if err != nil {
			internal.Log.Error("ControlApi.GetInfoFrApiName", zap.Any("funcName", funcName), zap.Any("uuid", uuidStr), zap.Any("apiName", "HEALTH_CHECK"), zap.Error(err))
		}
		wgDb.Done() // This decreases counter by 1
	}(wgDb)
	wgDb.Wait()
	internal.Log.Info("CheckDb Lan 2", zap.Any("funcName", funcName), zap.Any("uuid", uuidStr), zap.Any("dt", fmt.Sprintf("%v", time.Since(tStart1).Seconds())))
	sqlDB, err := db.DB()
	if err != nil {
		internal.Log.Error("Close db", zap.Any("funcName", funcName), zap.Any("uuid", uuidStr), zap.Error(err))
	}
	sqlDB.Close()
}
func CheckRedis(wg *sync.WaitGroup) {
	tStart := time.Now()
	status := 1 // normal
	defer func() {
		dt := time.Since(tStart).Seconds()
		if dt > 3 && status < 2 {
			status = 2
		}
		statusRedis = mapStatus[status]
		internal.Log.Info("CheckRedis", zap.Any("funcName", funcName), zap.Any("uuid", uuidStr), zap.Any("dt", fmt.Sprintf("%v", dt)), zap.Any("statusRedis", statusRedis))
		wg.Done() // This decreases counter by 1
	}()

	redisClient := internal.NewConnectRedis()
	if redisClient == nil {
		status = 3
		internal.Log.Error("NewConnectRedis nil", zap.Any("funcName", funcName), zap.Any("uuid", uuidStr))
		return
	}
	_, err := redisClient.Ping().Result()
	if err != nil {
		internal.Log.Error("redis ping fail", zap.Any("funcName", funcName), zap.Any("uuid", uuidStr), zap.Error(err))
		status = 3
		return
	}
	key := "test"
	_, err = redisClient.Get(key).Result()
	if err == redis.Nil || err == nil {
		return
	} else if err != nil {
		status = 2
		internal.Log.Error("Redis error", zap.Any("funcName", funcName), zap.Any("uuid", uuidStr), zap.Error(err))
		return
	}
}
func CheckKafka(wg *sync.WaitGroup) {
	tStart := time.Now()
	status := 1 // normal
	defer func() {
		dt := time.Since(tStart).Seconds()
		if dt > 3 && status < 2 {
			status = 2
		}
		statusKafka = mapStatus[status]
		internal.Log.Info("CheckKafka", zap.Any("funcName", funcName), zap.Any("uuid", uuidStr), zap.Any("dt", fmt.Sprintf("%v", dt)), zap.Any("statusKafka", statusKafka))
		wg.Done() // This decreases counter by 1
	}()

	writer := &kafka.Writer{
		Addr:     kafka.TCP(internal.Brokers...),
		Topic:    internal.KafkaTopicName,
		Balancer: &kafka.LeastBytes{},
		Async:    false,
	}
	ctxTimeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*3)
	defer cancelFunc()
	Message := map[string]interface{}{
		"ts":           tStart.Format("2006-01-02 15:04:05"),
		"service_name": internal.ServiceName,
		"title":        "healthcheck",
	}
	err := writer.WriteMessages(ctxTimeout,
		kafka.Message{
			Key:   []byte(internal.ServiceName),
			Value: internal.ToByte(Message),
		},
	)
	if err != nil {
		internal.Log.Error("CheckKafka", zap.Any("funcName", funcName), zap.Any("uuid", uuidStr), zap.Error(err))
		status = 3
	}
	if err := writer.Close(); err != nil {
		status = 2
		internal.Log.Error("failed to close writer", zap.Any("funcName", funcName), zap.Any("uuid", uuidStr), zap.Error(err))
	}
}
func CheckThirdParty(wg *sync.WaitGroup) {
	// input = số lượng các bên thứ 3
	// wg.Add(1)
	// go CheckApi(wg, "manapi.fpt.net", "url", false, nil, nil, nil, false)
}
func CheckApi(wg *sync.WaitGroup, key string, url string, isGet bool, headers, param map[string]string, body map[string]interface{}, proxy bool) {
	tStart := time.Now()
	status := 1 // normal
	defer func() {
		dt := time.Since(tStart).Seconds()
		if dt > 10 && status < 2 {
			status = 2
		}
		statusThird[key] = status
		wg.Done() // This decreases counter by 1
	}()
	resp, err := utils.Request(url, isGet, headers, param, body, 55, proxy)
	if err != nil {
		internal.Log.Error("CheckApi call fail", zap.Any("funcName", funcName), zap.Any("uuid", uuidStr), zap.Any("key", key), zap.Any("url", url), zap.Error(err))
		status = 3
		return
	}
	if resp.StatusCode() != 200 {
		internal.Log.Error("CheckApi call fail", zap.Any("funcName", funcName), zap.Any("uuid", uuidStr), zap.Any("key", key), zap.Any("url", url), zap.Any("statusCode", resp.StatusCode()))
		status = 3
		return
	}
}

package internal

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/segmentio/kafka-go"
	"github.com/shettyh/threadpool"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	Envs      = InitEnvVars()
	SysStatus = InitSystemStatus()
	Log       = TempLog{}
	Logger    = NewLogger()
	Db        = NewSQLDB()
)

var (
	Domains = InitAPIDomains(Envs.IsProduction)
	Eps     = InitAPIEndpoints()
	Keys    = InitAppKeys()
	Pool    = NewThreadPool()
	PoolLog = NewThreadPool()

	Brokers            = NewBrokers(Envs.IsProduction)
	KafkaTopicName     = NewKafkaTopicName(Envs.IsProduction)
	KafkaTopicPartner  = NewKafkaTopicPartner(Envs.IsProduction)
	KafkaTopicNameAll  = NewKafkaTopicNameAll(Envs.IsProduction)
	FROM_EMAIL         = "...."
	URL_SEND_MAIL_SMTP = "....."
	ServiceName        = "ecom-promotion-v2-api"
)

const (
	INF_COST                = 1e9
	NEG_INF_COST            = -1e9
	CODE_DB_FAILED          = 500
	CODE_WRONG_PARAMS       = 400
	CODE_RATE_LIMIT         = 999
	CODE_TOKEN_REQUIRED     = 1003
	CODE_TOKEN_EXPIRED      = 1001
	CODE_INVALID_TOKEN      = 1002
	CODE_SYSTEM_BUSY        = 800
	CODE_SYSTEM_ERROR       = 300
	CODE_TOKEN_APP_REQUIRED = 2003
	CODE_TOKEN_APP_EXPIRED  = 2001
	CODE_INVALID_TOKEN_APP  = 2002
	CODE_TOKEN_APP_ERROR    = 10403
	CODE_NOT_FOUND          = 404
	MSG_NOT_FOUND           = "Không tìm thấy thông tin."
	MSG_DB_FAILED           = "Chưa hiển thị được thông tin, vui lòng thử lại sau."
	MSG_WRONG_PARAMS        = "Chưa hiển thị được thông tin, vui lòng thử lại sau."
	MSG_RATE_LIMIT          = "Bạn truy cập quá nhanh."
	MSG_TOKEN_REQUIRED      = "Token không tồn tại"
	MSG_TOKEN_EXPIRED       = "Token hết hạn"
	MSG_INVALID_TOKEN       = "Token không hợp lệ"
	MSG_SYSTEM_BUSY         = "Chưa hiển thị được thông tin, vui lòng thử lại sau."
	MSG_SYSTEM_ERROR        = "Chưa hiển thị được thông tin, vui lòng thử lại sau."
	MSG_TOKEN_APP_REQUIRED  = "Token app không tồn tại"
	MSG_TOKEN_APP_EXPIRED   = "Token app hết hạn"
	MSG_INVALID_TOKEN_APP   = "Token app không hợp lệ"
	MSG_TOKEN_APP_ERROR     = "unauthorized"
	FuncNameKey             = contextKey("funcName")
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

type AppKeys struct {
	LOCAL_ECOM_CLIENT_KEY string
	LOCAL_ECOM_SECRET_KEY string
}
type EnvVars struct {
	SqlHost      string `mapstructure:"DB_HOST"`
	SqlUser      string `mapstructure:"DB_USERNAME"`
	SqlPassword  string `mapstructure:"DB_PASSWORD"`
	SqlDBName    string `mapstructure:"DB_DATABASE"`
	SqlPort      int    `mapstructure:"DB_PORT"`
	IsProduction bool   `mapstructure:"USE_PRODUCTION"`
	HostName     string `mapstructure:"HOSTNAME"`
	IsDev        bool   `mapstructure:"IS_DEV"`
	RedisHost    string `mapstructure:"REDIS_HOST"`
	RedisPort    int    `mapstructure:"REDIS_PORT"`
	RedisDb      int    `mapstructure:"REDIS_DB"`
}
type ApiDomains struct {
	ShoppingV1 string
	Billing    string
	Promotion  string

	OnlineReceipt     string
	FPTPayment        string
	LocalNotiProvider string
	HrAPI             string

	CustomerSop string
	Loyalty     string
}

type GroupCustomerSop struct {
	SyncCustomer      string
	SyncCustomerWifi6 string
}
type GroupLoyaltyApi struct {
	GetToken              string
	CustomerInfo          string
	VoucherAvailable      string
	VoucherExchanged      string
	VoucherExchange       string
	VoucherExchangeDetail string
}
type ApiEndpoints struct {
	CustomerSop GroupCustomerSop
	Loyalty     GroupLoyaltyApi
}

// Khởi tạo các API keys từ biến môi trường
func InitAppKeys() *AppKeys {
	fmt.Println("LOADING APP KEYS...")
	keys := &AppKeys{}

	// Đọc từ biến môi trường
	viper.BindEnv("LOCAL_ECOM_CLIENT_KEY")
	viper.BindEnv("LOCAL_ECOM_SECRET_KEY")
	// Thêm các keys khác nếu cần

	if err := viper.Unmarshal(keys); err != nil {
		fmt.Println("Error loading app keys:", err)
		fmt.Println("LOADING APP KEYS FAILED")
	}

	fmt.Println("LOADING APP KEYS SUCCESS")
	return keys
}

// Khởi tạo các trạng thái hệ thống với mã lỗi và message tương ứng
// Định nghĩa các response code chuẩn cho ứng dụng
func InitSystemStatus() *AllSystemStatus {
	return &AllSystemStatus{
		DbFailed: &SystemStatus{
			Status: CODE_DB_FAILED,
			Msg:    MSG_DB_FAILED,
		},
		WrongParams: &SystemStatus{
			Status: CODE_WRONG_PARAMS,
			Msg:    MSG_WRONG_PARAMS,
		}, RateLimit: &SystemStatus{
			Status: CODE_RATE_LIMIT,
			Msg:    MSG_RATE_LIMIT,
		}, TokenRequired: &SystemStatus{
			Status: CODE_TOKEN_REQUIRED,
			Msg:    MSG_TOKEN_REQUIRED,
		}, TokenExpired: &SystemStatus{
			Status: CODE_TOKEN_EXPIRED,
			Msg:    MSG_TOKEN_EXPIRED,
		}, InvalidToken: &SystemStatus{
			Status: CODE_INVALID_TOKEN,
			Msg:    MSG_INVALID_TOKEN,
		}, SystemBusy: &SystemStatus{
			Status: CODE_SYSTEM_BUSY,
			Msg:    MSG_SYSTEM_BUSY,
		}, SystemError: &SystemStatus{
			Status: CODE_SYSTEM_ERROR,
			Msg:    MSG_SYSTEM_ERROR,
		}, NotFound: &SystemStatus{
			Status: CODE_NOT_FOUND,
			Msg:    MSG_NOT_FOUND,
		},
		// Thêm Unauthorized
		Unauthorized: &SystemStatus{
			Status: 401,
			Msg:    "Unauthorized",
		},
		TokenAppRequired: &SystemStatus{
			Status: CODE_TOKEN_APP_REQUIRED,
			Msg:    MSG_TOKEN_APP_REQUIRED,
		},
		TokenAppExpired: &SystemStatus{
			Status: CODE_TOKEN_APP_EXPIRED,
			Msg:    MSG_TOKEN_APP_EXPIRED,
		},
		InvalidTokenApp: &SystemStatus{
			Status: CODE_INVALID_TOKEN_APP,
			Msg:    MSG_INVALID_TOKEN_APP,
		},
		TokenAppError: &SystemStatus{
			Status: CODE_TOKEN_APP_ERROR,
			Msg:    MSG_TOKEN_APP_ERROR,
		},
	}
}

type SystemStatus struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Detail interface{} `json:"detail"`
}
type AllSystemStatus struct {
	DbFailed         *SystemStatus
	WrongParams      *SystemStatus
	RateLimit        *SystemStatus
	TokenRequired    *SystemStatus
	TokenExpired     *SystemStatus
	InvalidToken     *SystemStatus
	NotFound         *SystemStatus
	SystemBusy       *SystemStatus
	SystemError      *SystemStatus
	Unauthorized     *SystemStatus // Thêm dòng này
	TokenAppRequired *SystemStatus
	TokenAppExpired  *SystemStatus
	InvalidTokenApp  *SystemStatus
	TokenAppError    *SystemStatus
}

func DB(funcName string) *gorm.DB {
	ctx := context.WithValue(context.Background(), FuncNameKey, funcName)
	return Db.WithContext(ctx)
}

// Tạo logger sử dụng zap, cấu hình output ra stdout với format JSON
func NewLogger() *zap.Logger {
	// Log vào stdOut
	writeSyncer := zapcore.AddSync(os.Stdout)
	// Thiết lập log
	loggerCore := zapcore.NewCore(logEndcoder(), writeSyncer, zap.DebugLevel)
	return zap.New(loggerCore)
}

func logEndcoder() zapcore.Encoder {
	encodeConfig := zap.NewProductionEncoderConfig()
	encodeConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	return zapcore.NewJSONEncoder(encodeConfig)
}

// Khởi tạo thread pool với 50 worker và queue size 100,000
func NewThreadPool() *threadpool.ThreadPool {
	fmt.Println("LOADING THREAD POOL ...")
	threadPool := threadpool.NewThreadPool(50, 100000)
	if threadPool == nil {
		fmt.Println("Failed")
	}
	fmt.Println("LOADING THREAD POOL SUCCESS...")
	return threadPool
}

// Khởi tạo và load các biến môi trường từ file app.env hoặc system env
// Sử dụng viper để đọc config, fallback sang system env nếu không tìm thấy file
func InitEnvVars() *EnvVars {
	fmt.Println("LOADING ENVS...")
	envs := &EnvVars{}
	viper.SetConfigFile("app.env")
	errEnvFile := viper.ReadInConfig()
	if errEnvFile != nil {
		viper.AutomaticEnv()
		viper.BindEnv("DB_HOST")
		viper.BindEnv("DB_USERNAME")
		viper.BindEnv("DB_PASSWORD")
		viper.BindEnv("DB_DATABASE")
		viper.BindEnv("DB_PORT")
		viper.BindEnv("USE_PRODUCTION")
		viper.BindEnv("HOSTNAME")
		viper.BindEnv("IS_DEV")
		viper.BindEnv("REDIS_HOST")
		viper.BindEnv("REDIS_PORT")
		viper.BindEnv("REDIS_DB")
	}
	if err := viper.Unmarshal(envs); err != nil {
		fmt.Println("Error viper.Unmarshal", err)
		fmt.Println("LOADING ENVS FAILED")
	}
	fmt.Println("LOADING ENVS SUCCESS")
	return envs
}

// Khởi tạo kết nối MySQL database sử dụng GORM
// Cấu hình connection pool với các thông số tối ưu
func NewSQLDB() *gorm.DB {
	fmt.Println("LOADING MYSQL DB ...")
	DBDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&interpolateParams=true",
		Envs.SqlUser, Envs.SqlPassword, Envs.SqlHost, Envs.SqlPort, Envs.SqlDBName)
	DBDSN = DBDSN + "&loc=Asia%2FHo_Chi_Minh"
	DB, err := sql.Open("mysql", DBDSN)
	if err != nil {
		// cm.Log.Debug("Open mysql connection failed", zap.Error(err))
		fmt.Println("NewSQLDB", err)
		return nil
	}
	DB.SetConnMaxLifetime(time.Minute * 10) //Connection pool trong Go(Mỗi connection trong pool sẽ bị đóng sau 10 phút tránh connection bị già )
	DB.SetMaxOpenConns(1000)                //Connection pool trong Go(Tối đa 1000 connection cùng lúc đến database.)
	DB.SetMaxIdleConns(1000)                //Connection pool trong Go(Giữ lại 1000 connection sẵn sàng, không đóng chúng.)
	errPing := DB.Ping()
	if errPing != nil {
		fmt.Println("DB PING ERR: ", errPing)
		return nil
	}
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: DB,
	}), &gorm.Config{
		Logger: &CustomLogger{},
	})
	if err != nil {
		fmt.Println("gorm DB connect failed", err.Error())
		return nil
	}
	fmt.Println("LOADING MYSQL DB SUCCESS...")
	return gormDB
}

// Lấy instance database với context chứa tên function để logging
func CheckSQLDB() (*gorm.DB, error) {
	fmt.Println("LOADING MYSQL DB ...")
	DBDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&interpolateParams=true",
		Envs.SqlUser, Envs.SqlPassword, Envs.SqlHost, Envs.SqlPort, Envs.SqlDBName)
	DBDSN = DBDSN + "&loc=Asia%2FHo_Chi_Minh"
	DB, err := sql.Open("mysql", DBDSN)
	if err != nil {
		// cm.Log.Debug("Open mysql connection failed", zap.Error(err))
		fmt.Println("NewSQLDB", err)
		return nil, err
	}
	DB.SetConnMaxLifetime(time.Minute * 10)
	DB.SetMaxOpenConns(1000)
	DB.SetMaxIdleConns(1000)
	errPing := DB.Ping()
	if errPing != nil {
		fmt.Println("DB PING ERR: ", errPing)
		return nil, err
	}
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: DB,
	}))
	if err != nil {
		fmt.Println("gorm DB connect failed", err.Error())
		return nil, err
	}
	fmt.Println("LOADING MYSQL DB SUCCESS...")
	return gormDB, nil
}

func NewBrokers(useProduction bool) []string {

	return []string{""}
}

func NewKafkaTopicName(useProduction bool) string {

	return ""
}

func NewKafkaTopicPartner(useProduction bool) string {

	return ""
}

func NewKafkaTopicNameAll(useProduction bool) string {

	return ""
}

func InitAPIDomains(isProduction bool) *ApiDomains {
	// productiongit puk
	if isProduction {
		return &ApiDomains{
			//....
		}
	}
	// staging
	return &ApiDomains{
		//...
	}
}

func InitAPIEndpoints() *ApiEndpoints {
	endpoints := &ApiEndpoints{}
	return endpoints
}

type TempLog struct{}

func (TempLog) InfoNoConsole(msg string, fields ...zap.Field) {
	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	listString := strings.Split(file, ServiceName+"/")
	file = shortCaller(listString[len(listString)-1])
	caller := fmt.Sprintf("%v %v:%v", shortFuncNameCaller(fn.Name()), file, line)
	SendLogToKibana("InfoNoConsole", msg, caller, fields...)
}

func (TempLog) Info(msg string, fields ...zap.Field) {
	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	listString := strings.Split(file, ServiceName+"/")
	file = shortCaller(listString[len(listString)-1])
	caller := fmt.Sprintf("%v %v:%v", shortFuncNameCaller(fn.Name()), file, line)
	SendLogToKibana("INFO", msg, caller, fields...)
	fields = append(fields, zap.Any("caller", caller))
	Logger.Info(msg, fields...)
	fmt.Println()
}

func (TempLog) InfoQuery(msg string, funcName string, query string, fields ...zap.Field) {
	n := 32
	listCaller := []string{}
	projectPath, _ := os.Getwd()
	var logQuery string
	isLogCaller := true
	now := time.Now()
	lengthCaller := 2
	for i := 3; i < n; i++ {
		pc, file, line, _ := runtime.Caller(i)
		fn := runtime.FuncForPC(pc)
		name := fn.Name()
		if strings.Contains(file, "preload") {
			logQuery = fmt.Sprintf("%v-%v--PRELOAD---%v", funcName, msg, query)
			isLogCaller = false
		} else if strings.Contains(file, "associations") {
			logQuery = fmt.Sprintf("%v-%v--ASSOCIATION---%v", funcName, msg, query)
			isLogCaller = false
		}
		if strings.HasPrefix(file, projectPath) && (strings.Contains(file, "/services") || strings.Contains(file, "/repositories")) {
			shortFile := file
			if idx := strings.LastIndex(file, "internal/"); idx != -1 {
				shortFile = file[idx+9:]
			}
			shortName := name
			if idx := strings.LastIndex(name, "("); idx != -1 {
				shortName = name[idx:]
			}
			listCaller = append(listCaller, fmt.Sprintf("%s:%d %s", shortFile, line, shortName))
			lengthCaller--
			if lengthCaller == 0 {
				break
			}
			continue
		}
	}
	// Nếu không tìm thấy caller, ghi log không có thông tin cụ thể
	if isLogCaller {
		logQuery = fmt.Sprintf("%v---%v---%v\n%v", funcName, msg, query, strings.Join(listCaller, "--"))
	}
	if len(fields) > 0 {
		fmt.Printf("%v:%v:%v\n", time.Since(now), logQuery, fields[0].Interface.(error).Error())
	} else {
		fmt.Printf("%v:%v\n", time.Since(now), logQuery)
	}
	// fmt.Printf("%v:%v\n", time.Since(now), logQuery)
	fields = append(fields, zap.Any("query", logQuery))
	SendLogToKibana("INFO", msg, listCaller, fields...)
}
func (TempLog) Error(msg string, fields ...zap.Field) {
	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	listString := strings.Split(file, ServiceName+"/")
	file = shortCaller(listString[len(listString)-1])
	caller := fmt.Sprintf("%v %v:%v", shortFuncNameCaller(fn.Name()), file, line)
	SendLogToKibana("ERROR", msg, caller, fields...)
	fields = append(fields, zap.Any("caller", caller))
	Logger.Error(msg, fields...)
	fmt.Println()
}

func ZapFieldToMap(logLevel, msg string, caller interface{}, fields ...zapcore.Field) map[string]interface{} {
	Message := map[string]interface{}{
		"a-title": msg,
		"z_log": map[string]interface{}{
			"caller":         caller,
			"log_level":      logLevel,
			"ts":             GetTimeUTC7().Format(time.DateTime),
			"service_name":   ServiceName,
			"container_name": Envs.HostName,
		},
	}
	for _, f := range fields {
		var value interface{}
		switch f.Type {
		case 15: // string
			value = f.String
		case 4: // bool
			value = f.Integer
		case 11, 12, 13, 14: // int64 int32 int16 int8
			value = f.Integer
		case 9: // float64
			value = math.Float64frombits(uint64(f.Integer))
		case 10: // float32
			value = math.Float32frombits(uint32(f.Integer))
		case 16: // TimeType
			if f.Interface != nil {
				value = time.Unix(0, f.Integer).In(f.Interface.(*time.Location))
			} else {
				// Fall back to UTC if location is nil.
				value = time.Unix(0, f.Integer)
			}
		case 17: // TimeFullType
			value = f.Interface.(time.Time)
		case 26: // ErrorType
			value = f.Interface.(error).Error()
		default:
			value = f.Interface
		}
		Message[f.Key] = value
	}
	return Message
}

func shortCaller(caller string) string {
	if strings.Contains(caller, "cmd/") {
		caller = strings.Split(caller, "cmd/")[1]
	}
	if strings.Contains(caller, "internal/") {
		caller = strings.Split(caller, "internal/")[1]
	}
	if strings.Contains(caller, "delivery/") {
		caller = strings.Split(caller, "delivery/")[1]
	}
	if strings.Contains(caller, "service/") {
		caller = strings.Split(caller, "service/")[1]
	}
	// if strings.Contains(caller, "repositories/") {
	// 	caller = strings.Split(caller, "repositories/")[1]
	// }
	return caller
}

func shortFuncNameCaller(funcName string) string {
	if strings.Contains(funcName, ServiceName+"/") {
		funcName = strings.Split(funcName, ServiceName+"/")[1]
	}
	if strings.Contains(funcName, "(") {
		funcName = "(" + strings.Split(funcName, "(")[1]
	}
	if strings.Contains(funcName, ".func1") {
		funcName = strings.Split(funcName, ".func1")[0]
	}
	return funcName
}

type CustomLogger struct {
	logger.Config
}

func (c *CustomLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *c
	newLogger.LogLevel = level
	return &newLogger
}

func (c *CustomLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if c.LogLevel >= logger.Info {
		log.Printf("[INFO] "+msg, data...)
	}
}

func (c *CustomLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if c.LogLevel >= logger.Warn {
		log.Printf("[WARN] "+msg, data...)
	}
}

func (c *CustomLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if c.LogLevel >= logger.Error {
		log.Printf("[ERROR] "+msg, data...)
	}
}

func (c *CustomLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if c.LogLevel <= 0 {
		return
	}
	funcName := "unknown"
	if ctx != nil {
		if val, ok := ctx.Value(FuncNameKey).(string); ok {
			funcName = val
		}
	}
	elapsed := time.Since(begin)
	sql, rows := fc()
	query := fmt.Sprintf("[%.3fms] [rows:%d] %s", float64(elapsed.Nanoseconds())/1e6, rows, sql)
	switch {
	case err != nil && c.LogLevel >= logger.Error:
		Log.InfoQuery("DEBUG_QUERY_ERROR", funcName, query, zap.Error(err))
	case elapsed > c.SlowThreshold && c.SlowThreshold != 0 && c.LogLevel >= logger.Warn:
		Log.InfoQuery("DEBUG_QUERY_SLOW", funcName, query)
	case c.LogLevel >= logger.Info:
		Log.InfoQuery("DEBUG_QUERY", funcName, query)
	}
}

// Convert interface thành byte array sử dụng JSON encoding
func ToByte(a interface{}) []byte {
	buffers := new(bytes.Buffer)
	json.NewEncoder(buffers).Encode(a)
	return buffers.Bytes()
}

// Gửi log đến hệ thống Kibana thông qua Kafka
// Chỉ hoạt động trong môi trường không phải dev
func SendLogToKibana(logLevel, msg string, caller interface{}, fields ...zapcore.Field) {
	if Envs.IsDev {
		return
	}
	writer := &kafka.Writer{
		MaxAttempts: 3,
		BatchBytes:  10485760,
		Addr:        kafka.TCP(Brokers...),
		Topic:       KafkaTopicName,
		Balancer:    &kafka.LeastBytes{},
		Async:       true,
	}
	ctxTimeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*3)
	defer cancelFunc()

	Message := ZapFieldToMap(logLevel, msg, caller, fields...)
	err := writer.WriteMessages(ctxTimeout,
		kafka.Message{
			Key:   []byte(ServiceName),
			Value: ToByte(Message),
		},
	)
	if err != nil {
		Log.Error("SendLogToKibana", zap.Error(err))
	}
	// Close the writer
	if err := writer.Close(); err != nil {
		Log.Error("failed to close writer", zap.Error(err))
	}
}

func GetTimeUTC7() time.Time {
	now := time.Now()
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	return now.In(loc)
}

// Log các request đến third-party APIs với đầy đủ thông tin:
// headers, body, response, timing, errors
func SendLogPartner(url string, headers map[string]string, body map[string]interface{}, resp *resty.Response, errCall error, status *int, startCall time.Time, endCall time.Time, funcName string) {
	if Envs.IsDev {
		return
	}
	writer := &kafka.Writer{
		MaxAttempts: 3,
		BatchBytes:  10485760,
		Addr:        kafka.TCP(Brokers...),
		Topic:       KafkaTopicPartner,
		Balancer:    &kafka.LeastBytes{},
		Async:       true,
	}
	ctxTimeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*3)
	defer cancelFunc()
	byteHeader, _ := json.Marshal(headers)
	byteBody, _ := json.Marshal(body)
	respStr := ""
	if resp != nil {
		respStr = string(resp.Body())
	}
	errStr := ""
	if errCall != nil {
		errStr = errCall.Error()
	}
	dt := endCall.Sub(startCall).Seconds()
	if resp != nil {
		dt = resp.Time().Seconds()
	}
	Message := map[string]interface{}{
		"container_name": Envs.HostName,
		"service":        ServiceName,
		"name":           "API",
		"function":       funcName,
		"date_created":   GetTimeUTC7().Format("2006-01-02 15:04:05"),
		"start_call":     startCall.Format("2006-01-02 15:04:05"),
		"end_call":       endCall.Format("2006-01-02 15:04:05"),
		"url_3rd":        url,
		"headers":        string(byteHeader),
		"input":          string(byteBody),
		"output":         respStr,
		"time_response":  dt,
		"status":         status,
		"error":          errStr,
	}
	err := writer.WriteMessages(ctxTimeout,
		kafka.Message{
			Key:   []byte(ServiceName),
			Value: ToByte(Message),
		},
	)
	if err != nil {
		Log.Error("SendLogToKibana", zap.Error(err))
	}
	// Close the writer
	if err := writer.Close(); err != nil {
		Log.Error("failed to close writer", zap.Error(err))
	}
}

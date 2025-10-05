package utils

import (
	"bytes"
	"context"
	"ecom_promotion_v2/internal"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitZapLogger(logType string) *zap.Logger {
	writer := getLogWriter(logType)
	encoder := getLogEncoder()
	core := zapcore.NewCore(encoder, writer, zapcore.DebugLevel)
	return zap.New(core, zap.AddCaller())
}
func MyCaller(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(filepath.Base(caller.FullPath()))
}
func getLogEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}
func getLogFileWriterSyncer() zapcore.WriteSyncer {
	loggerPath := initLoggerPath()
	layout := "2006-01-02"
	sep := getPathSeparator()
	fileNameSaved := fmt.Sprintf("%s%slogs%sinfo_debugs%s%s",
		loggerPath,
		sep,
		sep,
		sep,
		time.Now().Format(layout)+".log")
	file, err := os.OpenFile(fileNameSaved, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil
	}
	return zapcore.AddSync(file)
}
func getLogConsoleWriterSyncer() zapcore.WriteSyncer {
	return zapcore.AddSync(os.Stdout)
}
func getLogWriter(logType string) zapcore.WriteSyncer {
	switch logType {
	case "file":
		return getLogFileWriterSyncer()
	case "clg":
		return getLogConsoleWriterSyncer()
	case "*":
		return zap.CombineWriteSyncers(getLogFileWriterSyncer(), getLogConsoleWriterSyncer())
	default:
		return nil
	}
}
func initLoggerPath() string {
	GetAbsPath, err := os.Getwd()
	if err != nil {
		return ""
	}
	return GetAbsPath
}
func getPathSeparator() string {
	switch runtime.GOOS {
	case "linux":
		return "/"
	case "windows":
		return "\\"
	}
	return "/"
}

type KibanaMessage struct {
	Url          string      `json:"url"`
	ServiceName  string      `json:"mservice_name"`
	UserAgent    string      `json:"user_agent"`
	FuncName     string      `json:"fuc_name"`
	ActionName   string      `json:"action_name"`
	Token        string      `json:"token"`
	Input        interface{} `json:"input"`
	Output       interface{} `json:"output"`
	ExecutedTime float64     `json:"dt"`
	Version      string      `json:"version"`
}
type KibanaMessageAll struct {
	Phone          string  `json:"phone"`
	CustomerId     string  `json:"customerId"`
	IpAddress      string  `json:"ipAddress"`
	IsCanhTo       string  `json:"isCanhTo"`
	IsCustomer     string  `json:"isCustomer"`
	Provider       string  `json:"orovider"`
	DeviceId       string  `json:"deviceId"`
	DevicePlatform string  `json:"devicePlatform"`
	Lang           string  `json:"lang"`
	Input          string  `json:"input"`
	Output         string  `json:"output"`
	Headers        string  `json:"headers"`
	AppVersion     string  `json:"appVersion"`
	ContractNo     string  `json:"contractNo"`
	LocationZone   string  `json:"locationZone"`
	LocationCode   string  `json:"locationCode"`
	BranchName     string  `json:"branchName"`
	Status         int     `json:"status"`
	ServiceName    string  `json:"serviceName"`
	FunctionName   string  `json:"functionName"`
	ActionName     string  `json:"actionName"`
	Url            string  `json:"url"`
	DateAction     string  `json:"dateAction"`
	PositionIcon   string  `json:"positionIcon"`
	Referer        string  `json:"referer"`
	Note           string  `json:"note"`
	TypeLog        string  `json:"typeLog"`
	ProcessTime    float64 `json:"processTime"`
	Topic_name     string  `json:"topic_name"`
	ScreenId       string  `json:"screenId"`
}

type SendLogToKibanaTask struct {
	Message         KibanaMessage
	BootstrapServer []string
	TopicName       string
	ServiceName     string
}

func (m KibanaMessage) ToByte() []byte {
	buffers := new(bytes.Buffer)
	json.NewEncoder(buffers).Encode(m)
	return buffers.Bytes()
}
func (t *SendLogToKibanaTask) Run() {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(t.BootstrapServer...),
		Topic:    t.TopicName,
		Balancer: &kafka.LeastBytes{},
		Async:    true,
	}
	ctxTimeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*3)
	defer cancelFunc()
	err := writer.WriteMessages(ctxTimeout,
		kafka.Message{
			Key:   []byte(internal.ServiceName),
			Value: t.Message.ToByte(),
		},
	)
	if err != nil {
		internal.Log.Error("Error Send log to kibana", zap.Error(err))
	}
	// Close the writer
	if err := writer.Close(); err != nil {
		internal.Log.Error("failed to close writer", zap.Error(err))
	}
}

type SendLogToKibanaAllTask struct {
	Log             *zap.Logger
	Message         KibanaMessageAll
	BootstrapServer []string
	TopicName       string
}

func (m KibanaMessageAll) ToByte() []byte {
	buffers := new(bytes.Buffer)
	json.NewEncoder(buffers).Encode(m)
	return buffers.Bytes()
}
func (t *SendLogToKibanaAllTask) Run() {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(t.BootstrapServer...),
		Topic:    t.TopicName,
		Balancer: &kafka.LeastBytes{},
		Async:    true,
	}
	ctxTimeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*3)
	defer cancelFunc()
	err := writer.WriteMessages(ctxTimeout,
		kafka.Message{
			Key:   []byte(internal.ServiceName),
			Value: t.Message.ToByte(),
		},
	)
	if err != nil {
		internal.Log.Error("Error Send log to kibana all", zap.Error(err))
	}
	// Close the writer
	if err := writer.Close(); err != nil {
		internal.Log.Error("failed to close writer", zap.Error(err))
	}
}

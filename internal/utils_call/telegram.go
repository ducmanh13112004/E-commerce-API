package utils_call

import (
	"ecom_promotion_v2/internal"
	"ecom_promotion_v2/internal/utils"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// Map trạng thái tin nhắn
var mapStatus = map[int]string{
	1: "Thành công ✅",
	2: "Cảnh báo ⚠️",
	3: "Nghiêm trọng 🚨",
}

// BotTelegramService định nghĩa interface gửi tin nhắn
type BotTelegramService interface {
	SendMessage(message string) error
}

// telegramService quản lý việc gửi tin nhắn và retry
type telegramService struct {
	retryQueue chan string
	maxRetries int
	retryDelay time.Duration
}

// NewBotTelegramService khởi tạo và trả về đối tượng BotTelegramService
func NewBotTelegramService() (BotTelegramService, error) {
	service := &telegramService{
		retryQueue: make(chan string, 100),
		maxRetries: 3,
		retryDelay: 5 * time.Second,
	}
	go service.retryWorker()
	return service, nil
}

// retryWorker tự động retry tin nhắn nếu gửi thất bại
func (s *telegramService) retryWorker() {
	for message := range s.retryQueue {
		fmt.Println("Retrying message:", message)
		retries := 0
		for retries < s.maxRetries {
			if err := s.SendMessage(message); err != nil {
				retries++
				fmt.Printf("Retry #%d failed for message: %s\n", retries, message)
				time.Sleep(s.retryDelay)
			} else {
				break
			}
		}
		if retries >= s.maxRetries {
			internal.Log.Error("Message failed after retries", zap.String("message", message))
		}
	}
}

// SendMessage gửi tin nhắn qua bot Telegram
func (s *telegramService) SendMessage(message string) error {
	if !internal.Envs.IsProduction {
		return nil
	}
	if err := CallSendMessageTelegram(message); err != nil {
		internal.Log.Error("Failed to send message", zap.String("message", message))
		s.retryQueue <- message
		return fmt.Errorf("system error: %v", internal.SysStatus.SystemError)
	}
	return nil
}

// CallSendMessageTelegram gọi API Telegram để gửi tin nhắn
func CallSendMessageTelegram(message string) error {
	if !internal.Envs.IsProduction {
		return nil
	}
	defer utils.ExecTime(utils.GetTimeUTC7(), "CallSendMessageTelegram", nil)

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", internal.Keys.TokenKeyHiBotTelegram)
	input := map[string]interface{}{
		"chat_id": "-4643572148",
		"text":    message,
	}

	internal.Log.Info("Calling Telegram API", zap.String("url", url), zap.Any("input", input))

	resp, err := utils.Request(url, false, nil, nil, input, 15, true)
	if err != nil || resp.StatusCode() != 200 {
		internal.Log.Error("Telegram call failed", zap.Any("response", resp), zap.Error(err))
		return fmt.Errorf("telegram call failed: %v", err)
	}

	internal.Log.Info("Telegram response", zap.Any("response", resp.String()))
	return nil
}

// SendStatusMessage gửi tin nhắn với trạng thái cụ thể
func SendStatusMessage(content string, status int) {
	now := utils.GetStringTimeUTC7("Y-D-M H:M:S")
	message := fmt.Sprintf("%s\nService: %s\nTrạng thái: %s\nNội dung: %s",
		now, internal.ServiceName, mapStatus[status], content)
	CallSendMessageTelegram(message)
}

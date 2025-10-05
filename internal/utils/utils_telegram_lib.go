package utils

import (
	"context"
	"ecom_promotion_v2/internal"
	"ecom_promotion_v2/internal/models"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-telegram/bot"
	"go.uber.org/zap"
)

const chatID = -4643572148

// Map trạng thái tin nhắn
var mapStatus = map[int]string{
	1: "Thành công ✅",
	2: "Cảnh báo ⚠️",
	3: "Nghiêm trọng 🚨",
}

type MessageRetry struct {
	ChatID     int64
	Text       string
	RetryCount int
}

var (
	Bot        *bot.Bot = initBotTelegram()
	retryQueue chan MessageRetry
)

func initBotTelegram() *bot.Bot {
	// Sử dụng proxy
	proxyURL, _ := url.Parse("http://proxy.hcm.fpt.vn:80")
	httpTransport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	httpClient := &http.Client{
		Transport: httpTransport,
	}
	opts := []bot.Option{bot.WithHTTPClient(10*time.Second, httpClient)}
	bot, err := bot.New(internal.Keys.TOKEN_TELEGRAM, opts...)

	// Không sử dụng proxy
	// bot, err := bot.New(internal.Keys.TOKEN_TELEGRAM)
	if err != nil {
		internal.Log.Error("Error init telegram", zap.Any("funcName", "NewTelegramService"), zap.Error(err))
		return nil
	}

	return bot
}

// Dùng để tạo worker gửi lại tin nhắn, và bắt sự kiện từ người dùng telegram gửi tới khi cần
func StartTelegramService() {
	retryQueue = make(chan MessageRetry, 100)
	retryTelegramWorker()
	// go Bot.Start(context.Background())
	fmt.Println("Telegram bot started!")
}

// Worker chạy nền, retry tin nhắn sau 15s nếu lỗi, retry được 3 lần
func retryTelegramWorker() {
	if retryQueue == nil {
		internal.Log.Error("retryQueue is null", zap.Any("funcName", "sendMessageTelegram"))
		return
	}
	go func() {
		for msg := range retryQueue {
			if msg.RetryCount > 3 {
				internal.Log.Error("Retry 3 times", zap.Any("funcName", "sendMessageTelegram"), zap.Any("Message", msg))
				continue
			}
			time.Sleep(15 * time.Second)
			err := sendRetryMessage(msg.ChatID, msg.Text, msg.RetryCount)
			if err != nil {
				msg.RetryCount++
				retryQueue <- msg
			}
		}
	}()

}

func sendMessage(message string) error {
	if !internal.Envs.IsProduction {
		return nil
	}
	if Bot == nil {
		internal.Log.Error("Telegram not initialized", zap.Any("funcName", "sendMessage"), zap.Any("ChatID", chatID), zap.Any("Text", message))
		return nil
	}
	_, err := Bot.SendMessage(context.Background(), &bot.SendMessageParams{
		ChatID: chatID,
		Text:   message,
	})
	if err != nil {
		internal.Log.Error("Error sendTelegram", zap.Any("funcName", "sendMessageTelegram"), zap.Any("ChatID", chatID), zap.Any("Text", message), zap.Error(err))
		if retryQueue != nil {
			retryQueue <- MessageRetry{ChatID: chatID, Text: message, RetryCount: 1}
		}
	}
	return err
}

func ForwardMessage(message *models.ForwardMessage) error {
	if !internal.Envs.IsProduction {
		return nil
	}
	if Bot == nil {
		internal.Log.Error("Telegram not initialized", zap.Any("funcName", "ForwardMessage"), zap.Any("ChatID", chatID), zap.Any("Text", message))
		return nil
	}
	msg := fmt.Sprintf("%s\nTrạng thái: %s\nService: %s\nNội dung: %s", message.Time, mapStatus[message.Status], message.ServiceName, message.Message)
	_, err := Bot.SendMessage(context.Background(), &bot.SendMessageParams{
		ChatID: chatID,
		Text:   msg,
	})
	if err != nil {
		internal.Log.Error("Error ForwardMessage", zap.Any("funcName", "sendMessageTelegram"), zap.Any("ChatID", chatID), zap.Any("Text", message), zap.Error(err))
		return err
	}
	return nil
}

func SendTelegramMessage(content string, status int) {
	now := GetStringTimeUTC7("Y-D-M H:M:S")
	message := fmt.Sprintf("%s\nTrạng thái: %s\nService: %s\nNội dung: %s", now, mapStatus[status], internal.ServiceName, content)
	sendMessage(message)
}

func sendRetryMessage(chatID int64, text string, retryCount int) error {
	if Bot == nil {
		internal.Log.Error("Telegram not initialized", zap.Any("funcName", "sendRetryMessage"))
		return nil
	}
	_, err := Bot.SendMessage(context.Background(), &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	})
	if err != nil {
		internal.Log.Error("Error sendTelegram", zap.Any("funcName", "sendMessageTelegram"), zap.Any("ChatID", chatID), zap.Any("Text", text), zap.Error(err))
		return err
	}
	return nil
}

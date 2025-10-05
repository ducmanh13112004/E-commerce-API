package queue

import (
	"ecom_promotion_v2/internal"
	"ecom_promotion_v2/internal/utils"

	"go.uber.org/zap"
)

type QueueSendMail struct {
	MailReceive string
	MailCc      string
	MailBcc     string
	Subject     string
	Body        string
	AttachFile  string
	AttackUrl   string
}

func (c *QueueSendMail) Run() {
	url := internal.URL_SEND_MAIL_SMTP
	header := map[string]string{}
	body := map[string]interface{}{
		"FromEmail":        internal.FROM_EMAIL,
		"Recipients":       c.MailReceive,
		"CarbonCopys":      c.MailCc,
		"BlindCarbonCopys": c.MailBcc,
		"Subject":          c.Subject,
		"Body":             c.Body,
		"AttachFile":       c.AttachFile,
		"AttachUrl":        c.AttackUrl,
	}
	internal.Log.Info("QueueSendMail", zap.Any("url", url), zap.Any("body", body))
	responese, err := utils.Request(url, false, header, nil, body, 30, false)
	// res := struct {
	// 	Status      string
	// 	Description string
	// 	Data        interface{}
	// }{}
	if err != nil {
		internal.Log.Error("QueueSendMail fail", zap.Any("body", body), zap.Error(err))
		return
	}
	// json.Unmarshal([]byte(responese.String()), &res)
	internal.Log.Info("QueueSendMail", zap.Any("url", url), zap.Any("body", body), zap.Any("responese", responese))
}

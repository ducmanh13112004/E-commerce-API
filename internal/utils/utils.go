package utils

import (
	"crypto/aes"
	"crypto/hmac"
	"crypto/md5"
	cryptoRand "crypto/rand"
	"crypto/sha256"
	"ecom_promotion_v2/internal"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

func IsEmpty(value string) bool {
	return len(strings.TrimSpace(value)) == 0
}

// ///
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
func GetHmacSha256(data, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}
func GetSha256(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}

func CreateEcomToken(EcomClientKey, EcomSecretKey string) string {
	timestr := GetStringTimeUTC7("Y-D-M")
	return GetMD5Hash(EcomClientKey + "::" + EcomSecretKey + timestr)
}

///

func GetTimeUTC7() time.Time {
	now := time.Now()
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	return now.In(loc)
}

func GetTimeUTC7WithAddedDays(daysToAdd int) time.Time {
	return GetTimeUTC7().Add(time.Hour * time.Duration(24*daysToAdd))
}

func GetEndOfDayUTC7() time.Time {
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	timeWithOffset := GetTimeUTC7()
	year, month, day := timeWithOffset.Date()
	endOfDay := time.Date(year, month, day, 0, 0, 0, 0, loc).Add(time.Hour*24 - time.Second)
	return endOfDay
}

func GetTimeUTC7FrTime(input time.Time) time.Time {
	if input.Location().String() == "UTC" {
		input = input.Add(time.Hour * -7)
	}
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	return input.In(loc)
}
func GetStringTimeUTC7(stringtype string) string {
	t := GetTimeUTC7()
	switch stringtype {
	case "Y-D-M":
		return t.Format("2006-02-01")
	case "Y-M-D":
		return t.Format("2006-01-02")
	case "D-M-Y":
		return t.Format("02-01-2006")
	case "M-D-Y":
		return t.Format("01-02-2006")
	case "Y-D-M H:M:S":
		return t.Format("2006-02-01 15:04:05")
	case "Y-M-D H:M:S":
		return t.Format("2006-01-02 15:04:05")
	case "Y-M-D H:M:S -0700":
		return t.Format("2006-01-02 15:04:05 -0700")
	case "D/M/Y H:M:S":
		return t.Format("02/01/2006 15:04:05")
	}
	return ""
}
func GetStringTime(t time.Time, stringtype string) string {

	switch stringtype {
	case "Y-D-M":
		return t.Format("2006-02-01")
	case "Y-M-D":
		return t.Format("2006-01-02")
	case "D-M-Y":
		return t.Format("02-01-2006")
	case "D/M/Y":
		return t.Format("02/01/2006")
	case "D/M/Y H:M:S":
		return t.Format("02/01/2006 15:04:05")
	case "H:M:S D/M/Y":
		return t.Format("15:04:05 02/01/2006")
	case "M-D-Y":
		return t.Format("01-02-2006")
	case "Y-D-M H:M:S":
		return t.Format("2006-02-01 15:04:05")
	case "Y-M-D H:M:S":
		return t.Format("2006-01-02 15:04:05")
	case "Y-M-D H:M:S -0700":
		return t.Format("2006-01-02 15:04:05 -0700")
	}
	return ""
}
func GetBlackListContain(input string) string {
	blacklist := [...]string{"drop", "delete", "select", "update", "or", "and", "insert", "all",
		"=", "<>", "!=", ">", ">=", "<", "<=", "*", ";", "--"}
	for _, value := range blacklist {
		if strings.Contains(input, value) {
			return value
		}
	}
	return ""
}
func EncodeBase64(text string) string {
	return base64.StdEncoding.EncodeToString([]byte(text))
}
func DecodeBase64(textEncoded string) (string, error) {
	enc, err := base64.StdEncoding.DecodeString(textEncoded)
	return string(enc), err
}
func ParseTimeFrString(stringtype string, timeinput string) (time.Time, error) {
	layout := ""
	switch stringtype {
	case "Y-D-M":
		layout = "2006-02-01"
	case "Y-M-D":
		layout = "2006-01-02"
	case "D-M-Y":
		layout = "02-01-2006"
	case "D-M-Y H:M:S":
		layout = "02-01-2006 15:04:05"
	case "M-D-Y":
		layout = "01-02-2006"
	case "Y-D-M H:M:S":
		layout = "2006-02-01 15:04:05"
	case "Y-M-D H:M:S":
		layout = "2006-01-02 15:04:05"
	case "Y-M-D H:M:S -0700":
		layout = "2006-01-02 15:04:05 -0700"
	case "Y-M-DTH:M:S +0700":
		layout = "2006-01-02T15:04:05 +0700"
	case "Y-M-DTH:M:S+07:00":
		layout = "2006-01-02T15:04:05+07:00"
	case "Y-M-DTH:M:S.000":
		layout = "2006-01-02T15:04:05.999999999"
	case "D/M/Y":
		layout = "02/01/2006"
	case "D/M/Y H:M:S":
		layout = "02/01/2006 15:04:05"
	case "M/D/Y H:M:S":
		layout = "01/02/2006 15:04:05"
	case "H:M:S D/M/Y":
		layout = "15:04:05 02/01/2006"
	}
	t, err := time.Parse(layout, timeinput)
	if err != nil {
		return t, err
	}
	return GetTimeUTC7FrTime(t), nil
}
func FloatToTime(input float64) time.Time {
	integ, decim := math.Modf(input)
	return time.Unix(int64(integ), int64(decim*(1e9)))
}
func StringToFloat64(input string) (float64, error) {
	result, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0, err
	}
	return result, nil
}
func StringToInt(input string) (int, error) {
	result, err := strconv.Atoi(input)
	if err != nil {
		return 0, err
	}
	return result, nil
}

var Weekdays = map[string]string{
	"Sunday":    "Chủ nhật",
	"Monday":    "Thứ hai",
	"Tuesday":   "Thứ ba",
	"Wednesday": "Thứ tư",
	"Thursday":  "Thứ năm",
	"Friday":    "Thứ sáu",
	"Saturday":  "Thứ bảy",
}
var WeekdayV2 = map[int]string{
	0: "Chủ nhật",
	1: "Thứ 2",
	2: "Thứ 3",
	3: "Thứ 4",
	4: "Thứ 5",
	5: "Thứ 6",
	6: "Thứ 7",
}

func DayInWeek(day time.Weekday) string {
	// fmt.Printf("day.String(): %v\n", day.String())
	return Weekdays[day.String()]
}

func StringToFormatTime(stringtype string) string {
	switch stringtype {
	case "Y-D-M":
		return "2006-02-01"
	case "Y-M-D":
		return "2006-01-02"
	case "D-M-Y":
		return "02-01-2006"
	case "M-D-Y":
		return "01-02-2006"
	case "Y-D-M H:M:S":
		return "2006-02-01 15:04:05"
	case "Y-M-D H:M":
		return "2006-01-02 15:04"
	case "Y-M-D H:M:S":
		return "2006-01-02 15:04:05"
	case "Y-M-D H:M:S -0700":
		return "2006-01-02 15:04:05 -0700"
	case "D/M/Y":
		return "02/01/2006"
	case "D/M/Y H:M":
		return "02/01/2006 15:04"
	case "D-M-Y H:M":
		return "02-01-2006 15:04"
	case "H:M D/M/Y":
		return "15:04 02/01/2006"
	case "Y-M-DTH:M:S.000":
		return "2006-01-02T15:04:05.999999999"
	case "D/M/Y H:M:S.000":
		return "02-01-2006 15:04:05.999999999"
	case "D/M/Y H:M:S":
		return "02/01/2006 15:04:05"
	case "D/M/Y - H:M":
		return "02/01/2006 - 15:04"
	}
	return ""
}

func ParseJwt(token string, key string) (jwt.MapClaims, *internal.SystemStatus) {
	clams := jwt.MapClaims{}
	decodeToken, err := jwt.ParseWithClaims(token, clams, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		jwtErr := err.(*jwt.ValidationError).Errors
		if jwtErr == jwt.ValidationErrorExpired {
			internal.Log.Error("ValidationErrorExpired", zap.Any("input", token), zap.Any("tokenParse", decodeToken), zap.Error(err))
			return nil, internal.SysStatus.TokenExpired
		}
	}
	if decodeToken == nil || !decodeToken.Valid {
		internal.Log.Error("Token Valid", zap.Any("input", token), zap.Any("tokenParse", decodeToken))
		return nil, internal.SysStatus.InvalidToken
	}
	return clams, nil
}

func EncryptAES(key []byte, plaintext string) (string, error) {
	//which requires that the input data be a multiple of the block size (which is 16 bytes for AES).
	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	out := make([]byte, len(plaintext))

	c.Encrypt(out, []byte(plaintext))

	return hex.EncodeToString(out), nil
}

func DecryptAES(key []byte, ct string) (string, error) {
	ciphertext, _ := hex.DecodeString(ct)

	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	pt := make([]byte, len(ciphertext))
	c.Decrypt(pt, ciphertext)

	s := string(pt[:])
	return s, nil
}

func ExecTime(start time.Time, funcName string, log *zap.Logger) {
	logStr := fmt.Sprintf("ExecTime dt=%d ms", GetTimeUTC7().Sub(start).Milliseconds())
	internal.Log.Info("ExecTime of "+funcName, zap.Any("log", logStr), zap.Any("dt", GetTimeUTC7().Sub(start).Milliseconds()))
}

func CheckItemInListContains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

// Generate is a low-level function to change alphabet and ID size.
// Opensource: https://github.com/matoous/go-nanoid
func GenerateNanoId(alphabet string, size int) (string, error) {
	chars := []rune(alphabet)

	if len(alphabet) == 0 || len(alphabet) > 255 {
		return "", errors.New("alphabet must not be empty and contain no more than 255 chars")
	}
	if size <= 0 {
		return "", errors.New("size must be positive integer")
	}

	getMask := func(alphabetSize int) int {
		for i := 1; i <= 8; i++ {
			mask := (2 << uint(i)) - 1
			if mask >= alphabetSize-1 {
				return mask
			}
		}
		return 0
	}

	mask := getMask(len(chars))
	// estimate how many random bytes we will need for the ID, we might actually need more but this is tradeoff
	// between average case and worst case
	ceilArg := 1.6 * float64(mask*size) / float64(len(alphabet))
	step := int(math.Ceil(ceilArg))

	id := make([]rune, size)
	bytes := make([]byte, step)
	for j := 0; ; {
		_, err := cryptoRand.Read(bytes)
		if err != nil {
			return "", err
		}
		for i := 0; i < step; i++ {
			currByte := bytes[i] & byte(mask)
			if currByte < byte(len(chars)) {
				id[j] = chars[currByte]
				j++
				if j == size {
					return string(id[:size]), nil
				}
			}
		}
	}
}

func GetFunctionName() string {
	pc, _, _, _ := runtime.Caller(1) // 1 để lấy thông tin của hàm gọi GetFunctionName
	funcDetails := runtime.FuncForPC(pc)
	funcName := funcDetails.Name()
	// Chia tên hàm và lấy phần cuối cùng
	parts := strings.Split(funcName, ".")
	return parts[len(parts)-1]
}

func ResponseString(resp *resty.Response) interface{} {
	if resp == nil {
		return nil
	}
	result := map[string]interface{}{}
	err := json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return resp.String()
	}
	return result
}

package common

import (
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	ent "main/internal/v1/ticket/domain/entities"
	"math/rand"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"
)

var (
	Bundle *i18n.Bundle
)

func InitializeI18N() {
	Bundle = i18n.NewBundle(language.Turkish)
	Bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	Bundle.MustLoadMessageFile("i18n/tr.json")
	Bundle.MustLoadMessageFile("i18n/en.json")
}

func HTTPResponser(data interface{}, status bool, message string) fiber.Map {
	return fiber.Map{
		"error":   status,
		"message": message,
		"data":    data,
	}
}

func GenNum() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(999999-100000) + 100000
}

// Get config path for local or docker
func GetConfigPath(configPath string) string {
	if configPath == "docker" {
		return "./config/config-docker"
	}
	return "./config/config-local"
}

func ValueTrim(value string) string {
	result := strings.TrimSpace(value)
	return result
}

func Placeholder(data []string) string {
	var columns []string
	for i := 1; i <= len(data); i++ {
		columns = append(columns, fmt.Sprintf("$%d", i))
	}
	return strings.Join(columns, ",")
}

func Column(data []string) string {
	var columns []string
	for i := 0; i < len(data); i++ {
		columns = append(columns, data[i])
	}
	return strings.Join(columns, ",")
}

func ThrowError(defination string) error {
	return errors.New(defination)
}

func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}

// CheckStringIfContains check a string if contains given param
func CheckStringIfContains(input_text string, search_text string) bool {
	CheckContains := strings.Contains(input_text, search_text)
	return CheckContains
}

func GetAuthIdFromToken(c *fiber.Ctx) (int64, int8) {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	AuthID := int64(claims["user_id"].(float64))
	UserType := int8(claims["user_type"].(float64))

	return AuthID, UserType
}

func getStructType(struc interface{}) reflect.Type {
	sType := reflect.TypeOf(struc)
	if sType.Kind() == reflect.Ptr {
		sType = sType.Elem()
	}

	return sType
}

func GetAuthDataFromToken(c *fiber.Ctx) (dat ent.UserData) {

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	mapper := ent.UserData{
		UserId: int64(claims["user_id"].(float64)),
	}

	return mapper

}

func GetRawAccessToken(c *fiber.Ctx) (token string) {
	user := c.Locals("user").(*jwt.Token)
	rawToken := user.Raw
	return rawToken
}

func RemoveBasePath(path string) string {
	pathSlice := strings.Split(path, "/")
	pathSlice = pathSlice[2:]
	newPath := fmt.Sprintf("/%s", strings.Join(pathSlice, "/"))
	return newPath
}

func Translate(c *fiber.Ctx, id string) (translation string) {

	var (
		lang = "tr"
	)

	if strings.TrimSpace(c.Get("X-Lang-Code")) != "" {
		lang = c.Get("X-Lang-Code")
	}

	localized := i18n.NewLocalizer(Bundle, lang)
	translation = localized.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: id,
		},
	})
	return
}

package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

var config *Config
var once sync.Once

type Config struct {
	IsDebug    bool
	Web        WebConfig
	PostgresDB PostgresDBConfig
}

type WebConfig struct {
	BasePath                string
	PORT                    string
	CORSAllowOrigin         []string
	WebURL                  string
	CookieDomain            string
	CookieSessionKey        string
	UseRedisSession         bool
	RedisSessionIndex       int //index for session
	RedisUserInfoIndex      int //index for userinfo
	LoginSuccessRedirectUrl string
	LoginErrorRedirectUrl   string
	HTTPOnly                bool
	Secure                  bool // set true in production (HTTPS)
	SameSite                string
}

type PostgresDBConfig struct {
	ConnectionString   string
	MaxOpenConn        int
	MaxIdleConn        int
	ConnMaxLifeTimeTTL *time.Duration
}

func GetConfig() *Config {
	once.Do(func() {
		GetConfigWithFilename(".env")
		config = &Config{
			IsDebug: getEnvBool("DEBUG", true),
			Web: WebConfig{
				BasePath:                getEnvString("BASEPATH", ""),
				PORT:                    getEnvString("PORT", "8080"),
				CORSAllowOrigin:         getEnvStringArray("ACCESS_ORIGIN", []string{}),
				WebURL:                  getEnvString("WEB_URL", ""),
				CookieDomain:            getEnvString("COOKIE_DOMAIN", ""),
				CookieSessionKey:        getEnvString("COOKIE_SESSION_KEY", "session-id"),
				UseRedisSession:         getEnvBool("USE_REDIS_SESSION", false),
				RedisSessionIndex:       getEnvInt("REDIS_SESSION_INDEX", 13),
				RedisUserInfoIndex:      getEnvInt("REDIS_USERINFO_INDEX", 3),
				LoginSuccessRedirectUrl: joinPath(getEnvString("WEB_URL", ""), getEnvString("LOGIN_SUCCESS_REDIECT_PATH", "")),
				LoginErrorRedirectUrl:   joinPath(getEnvString("WEB_URL", ""), getEnvString("LOGIN_ERROR_REDIECT_PATH", "")),
				HTTPOnly:                getEnvBool("COOKIE_HTTPONLY", true),
				Secure:                  getEnvBool("COOKIE_SECURE", false),
				SameSite:                getEnvString("COOKIE_SAMESITE", "Lax"),
			},
			PostgresDB: PostgresDBConfig{
				ConnectionString:   getEnvString("POSTGRES_CONNECTION_STRING", ""),
				MaxOpenConn:        getEnvInt("POSTGRES_MAX_OPEN_CONN", 0),
				MaxIdleConn:        getEnvInt("POSTGRES_MAX_IDLE_CONN", 0),
				ConnMaxLifeTimeTTL: getEnvDurationFromSecondsNullable("POSTGRES_CONN_MAX_LIFE_TIME_SECONDS", 0),
			},
		}
	})
	return config
}

func GetConfigWithFilename(envFileName string) {

	doloadFile := false

	//will replace the existing env values
	if godotenv.Load(envFileName) == nil {
		doloadFile = true
	}
	if godotenv.Load(fmt.Sprintf("../%s", envFileName)) == nil {
		doloadFile = true
	}
	if godotenv.Load(fmt.Sprintf("../../%s", envFileName)) == nil {
		doloadFile = true
	}
	if godotenv.Load(fmt.Sprintf("../../../%s", envFileName)) == nil {
		doloadFile = true
	}

	if godotenv.Load(fmt.Sprintf("../../../../%s", envFileName)) == nil {
		doloadFile = true
	}
	if godotenv.Load(fmt.Sprintf("../../../../../%s", envFileName)) == nil {
		doloadFile = true
	}
	if godotenv.Load(fmt.Sprintf("../../../../../../%s", envFileName)) == nil {
		doloadFile = true
	}
	if godotenv.Load(fmt.Sprintf("../../../../../../%s", envFileName)) == nil {
		doloadFile = true
	}
	if godotenv.Load(fmt.Sprintf("../../../../../../../%s", envFileName)) == nil {
		doloadFile = true
	}

	if !doloadFile {
		log.Println("WARNING: failed to load .env file")
	}
}

//lint:ignore U1000 Ignore unused code, it may require in the future
func getEnvString(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

//lint:ignore U1000 Ignore unused code, it may require in the future
func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}

//lint:ignore U1000 Ignore unused code, it may require in the future
func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return boolValue
}

//lint:ignore U1000 Ignore unused code, it may require in the future
func getEnvStringArray(key string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	result := strings.Split(value, ",")
	for i := range result {
		result[i] = strings.TrimSpace(result[i])
	}

	return result
}

//lint:ignore U1000 Ignore unused code, it may require in the future
func joinPath(baseurl string, paths ...string) string {
	combined, err := url.JoinPath(baseurl, paths...)
	if err != nil {
		return ""
	}
	return combined
}

//lint:ignore U1000 Ignore unused code, it may require in the future
func getEnvDurationFromSecondsNullable(key string, defaultValue time.Duration) *time.Duration {
	value := os.Getenv(key)
	if value == "" {
		if defaultValue == 0 {
			return nil
		} else {
			return &defaultValue
		}
	}

	intValue, err := strconv.ParseInt(value, 10, 64)

	if err != nil {
		return &defaultValue
	}

	result := time.Duration(intValue) * time.Second
	return &result
}

package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type Config struct {
	// main
	Environment string
	LogLevel    string
	RPCPort     string

	// jwt
	SigningKey        []byte
	RefreshSigningKey []byte

	// database
	PostgresHost     string
	PostgresPort     int
	PostgresUser     string
	PostgresPassword string
	PostgresDatabase string

	// minio
	BucketName          string `env:"MINIO_BUCKET_NAME" default:"files"`
	MinioDomain         string `env:"MINIO_DOMAIN" default:"test-cdn.yerelektron.uz"` //"http://127.0.0.1:9199"
	MinioAccessKeyID    string `env:"MINIO_ACCESS_KEY" default:"8DbGdJfNjQmSqVsXv2x4z7C9EbHeKgNkRnTrWtYv3y5A7DaFcJfMhPmSpUrXuZw3z6B8EbGdJgNjQmTqVsXv2x4A7C"`
	MinioSecretAccesKey string `env:"MINIO_SECRET_KEY" default:"QmSpUsXuZx4z6B9EbGdKgNjQnTqVtYv2x5A7C9FcHeKhPkRpUrWtZw3y5B8DaFdJfMjQmSpVsXuZx4z6B9EbGeKgNj"`
	MinioFilesBucketURL string `env:"MINIO_FILES_BUCKET_URL" default:"https://test-cdn.yerelektron.uz/files/"`

	// services
	TemplateServiceHost string `env:"TEMPLATE_SERVICE_HOST"`
	TemplateServicePort int    `env:"TEMPLATE_SERVICE_PORT"`
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
	cfg := Config{}

	/*                MAIN                */
	cfg.Environment = cast.ToString(getOrReturnDefaultValue("ENVIRONMENT", "DEV"))
	cfg.LogLevel = cast.ToString(getOrReturnDefaultValue("LOG_LEVEL", "debug"))
	cfg.RPCPort = cast.ToString(getOrReturnDefaultValue("RPC_PORT", ":8109"))

	/*                 JWT                 */
	cfg.SigningKey = []byte(cast.ToString(getOrReturnDefaultValue("SIGNING_KEY", "ZWxlY3Ryb24ga2FkYXN0ciBzZXJ2aWNlIGZvciBnb3Zlcm1lbnQgc2VydmljZXMgLT4gdWRldnMgZGV2ZWxvcGVkCg==")))
	cfg.RefreshSigningKey = []byte(cast.ToString(getOrReturnDefaultValue("SIGNING_KEY", "ZWxlY3Ryb24ga2FkYXN0ciBzZXJ2aWNlIGZvciBnb3Zlcm1lbnQgc2VydmljZXMgcmVmcmVzaCB0b2tlbiAtPiB1ZGV2cyBkZXZlbG9wZWQK")))

	/*                MINIO                */
	cfg.MinioAccessKeyID = cast.ToString(getOrReturnDefaultValue("MINIO_ACCESS_KEY", "8DbGdJfNjQmSqVsXv2x4z7C9EbHeKgNkRnTrWtYv3y5A7DaFcJfMhPmSpUrXuZw3z6B8EbGdJgNjQmTqVsXv2x4A7C"))
	cfg.MinioSecretAccesKey = cast.ToString(getOrReturnDefaultValue("MINIO_SECRET_KEY", "QmSpUsXuZx4z6B9EbGdKgNjQnTqVtYv2x5A7C9FcHeKhPkRpUrWtZw3y5B8DaFdJfMjQmSpVsXuZx4z6B9EbGeKgNj"))
	cfg.MinioDomain = cast.ToString(getOrReturnDefaultValue("MINIO_DOMAIN", "test-cdn.yerelektron.uz"))
	cfg.MinioFilesBucketURL = cast.ToString(getOrReturnDefaultValue("MINIO_FILES_BUCKET_URL", "https://test-cdn.yerelektron.uz/files/"))
	cfg.BucketName = cast.ToString(getOrReturnDefaultValue("MINIO_BUCKET_NAME", "files"))

	/*           DATABASE (POSTGRES)        */
	// cfg.PostgresHost = cast.ToString(getOrReturnDefaultValue("POSTGRES_HOST", "localhost"))
	// cfg.PostgresPort = cast.ToInt(getOrReturnDefaultValue("POSTGRES_PORT", 5432))
	// cfg.PostgresUser = cast.ToString(getOrReturnDefaultValue("POSTGRES_USER", "falck"))
	// cfg.PostgresPassword = cast.ToString(getOrReturnDefaultValue("POSTGRES_PASSWORD", "falck"))
	// cfg.PostgresDatabase = cast.ToString(getOrReturnDefaultValue("POSTGRES_DATABASE", "postgres"))
	cfg.PostgresHost = cast.ToString(getOrReturnDefaultValue("POSTGRES_HOST", "192.168.112.15"))
	cfg.PostgresPort = cast.ToInt(getOrReturnDefaultValue("POSTGRES_PORT", 30931))
	cfg.PostgresDatabase = cast.ToString(getOrReturnDefaultValue("POSTGRES_DATABASE", "integration_service"))
	cfg.PostgresUser = cast.ToString(getOrReturnDefaultValue("POSTGRES_USER", "postgres"))
	cfg.PostgresPassword = cast.ToString(getOrReturnDefaultValue("POSTGRES_PASSWORD", "ATs7tCVrBA"))

	/*               SERVICES               */
	cfg.TemplateServiceHost = cast.ToString(getOrReturnDefaultValue("TEMPLATE_SERVICE_HOST", "localhost"))
	cfg.TemplateServicePort = cast.ToInt(getOrReturnDefaultValue("TEMPLATE_SERVICE_PORT", 6436))
	return cfg
}

func getOrReturnDefaultValue(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)

	if exists {
		return os.Getenv(key)
	}

	return defaultValue
}

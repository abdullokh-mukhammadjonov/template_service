package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type Config struct {
	Environment string

	PostgresHost     string
	PostgresPort     int
	PostgresUser     string
	PostgresPassword string
	PostgresDatabase string

	SigningKey        []byte
	RefreshSigningKey []byte
	LogLevel          string
	RPCPort           string

	AuctionUrl               string
	AuctionUsername          string
	AuctionUsernameTypeSix   string
	AuctionPassword          string
	AuctionPasswordTypeSix   string
	AuctionOrderGetURL       string
	AuctionOrderCreateURL    string
	AuctionDocumentCreateURL string
	AuctionImageCreateURL    string
	AuctionDocumentUpdateURL string
	AuctionImageUpdateURL    string
	AuctionSendOrderURL      string
	AuctionGetProtocolURL    string

	DiscussionLogicServiceHost string
	DiscussionLogicServicePort int

	EntityServiceHost string
	EntityServicePort int

	UserServiceHost string
	UserServicePort int

	BucketName          string `env:"MINIO_BUCKET_NAME" default:"files"`
	MinioDomain         string `env:"MINIO_DOMAIN" default:"test-cdn.yerelektron.uz"` //"http://127.0.0.1:9199"
	MinioAccessKeyID    string `env:"MINIO_ACCESS_KEY" default:"8DbGdJfNjQmSqVsXv2x4z7C9EbHeKgNkRnTrWtYv3y5A7DaFcJfMhPmSpUrXuZw3z6B8EbGdJgNjQmTqVsXv2x4A7C"`
	MinioSecretAccesKey string `env:"MINIO_SECRET_KEY" default:"QmSpUsXuZx4z6B9EbGdKgNjQnTqVtYv2x5A7C9FcHeKhPkRpUrWtZw3y5B8DaFdJfMjQmSpVsXuZx4z6B9EbGeKgNj"`
	MinioFilesBucketURL string `env:"MINIO_FILES_BUCKET_URL" default:"https://test-cdn.yerelektron.uz/files/"`

	HokimyatUrl             string
	HokimyatTokenUrl        string
	HokimyatUrlDownloadFile string
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
	cfg := Config{}

	cfg.Environment = cast.ToString(getOrReturnDefaultValue("ENVIRONMENT", "DEV"))
	cfg.LogLevel = cast.ToString(getOrReturnDefaultValue("LOG_LEVEL", "debug"))
	cfg.MinioAccessKeyID = cast.ToString(getOrReturnDefaultValue("MINIO_ACCESS_KEY", "8DbGdJfNjQmSqVsXv2x4z7C9EbHeKgNkRnTrWtYv3y5A7DaFcJfMhPmSpUrXuZw3z6B8EbGdJgNjQmTqVsXv2x4A7C"))
	cfg.MinioSecretAccesKey = cast.ToString(getOrReturnDefaultValue("MINIO_SECRET_KEY", "QmSpUsXuZx4z6B9EbGdKgNjQnTqVtYv2x5A7C9FcHeKhPkRpUrWtZw3y5B8DaFdJfMjQmSpVsXuZx4z6B9EbGeKgNj"))
	cfg.MinioDomain = cast.ToString(getOrReturnDefaultValue("MINIO_DOMAIN", "test-cdn.yerelektron.uz"))
	cfg.BucketName = cast.ToString(getOrReturnDefaultValue("MINIO_BUCKET_NAME", "files"))
	cfg.MinioFilesBucketURL = cast.ToString(getOrReturnDefaultValue("MINIO_FILES_BUCKET_URL", "https://test-cdn.yerelektron.uz/files/"))

	// cfg.PostgresHost = cast.ToString(getOrReturnDefaultValue("POSTGRES_HOST", "localhost"))
	// cfg.PostgresPort = cast.ToInt(getOrReturnDefaultValue("POSTGRES_PORT", 5432))
	// cfg.PostgresUser = cast.ToString(getOrReturnDefaultValue("POSTGRES_USER", "falck"))
	// cfg.PostgresPassword = cast.ToString(getOrReturnDefaultValue("POSTGRES_PASSWORD", "falck"))
	// cfg.PostgresDatabase = cast.ToString(getOrReturnDefaultValue("POSTGRES_DATABASE", "postgres"))

	cfg.PostgresHost = cast.ToString(getOrReturnDefaultValue("POSTGRES_HOST", "192.168.112.15"))
	cfg.PostgresPort = cast.ToInt(getOrReturnDefaultValue("POSTGRES_PORT", 30931))
	cfg.PostgresDatabase = cast.ToString(getOrReturnDefaultValue("POSTGRES_DATABASE", "ek_integration_service"))
	cfg.PostgresUser = cast.ToString(getOrReturnDefaultValue("POSTGRES_USER", "postgres"))
	cfg.PostgresPassword = cast.ToString(getOrReturnDefaultValue("POSTGRES_PASSWORD", "ATs7tCVrBA"))
	cfg.SigningKey = []byte(cast.ToString(getOrReturnDefaultValue("SIGNING_KEY", "ZWxlY3Ryb24ga2FkYXN0ciBzZXJ2aWNlIGZvciBnb3Zlcm1lbnQgc2VydmljZXMgLT4gdWRldnMgZGV2ZWxvcGVkCg==")))
	cfg.RefreshSigningKey = []byte(cast.ToString(getOrReturnDefaultValue("SIGNING_KEY", "ZWxlY3Ryb24ga2FkYXN0ciBzZXJ2aWNlIGZvciBnb3Zlcm1lbnQgc2VydmljZXMgcmVmcmVzaCB0b2tlbiAtPiB1ZGV2cyBkZXZlbG9wZWQK")))
	cfg.RPCPort = cast.ToString(getOrReturnDefaultValue("RPC_PORT", ":8009"))

	cfg.AuctionOrderGetURL = cast.ToString(getOrReturnDefaultValue("ORDER_GET_URL", "http://10.190.4.122:8390/ws/api/services/common/get-order"))
	cfg.AuctionOrderCreateURL = cast.ToString(getOrReturnDefaultValue("ORDER_CREATE_URL", "http://10.190.4.122:8390/ws/api/services/cadaster/order/creates"))
	cfg.AuctionDocumentCreateURL = cast.ToString(getOrReturnDefaultValue("DOCUMENT_CREATE_URL", "http://10.190.4.122:8390/ws/api/services/documents"))
	cfg.AuctionDocumentUpdateURL = cast.ToString(getOrReturnDefaultValue("DOCUMENT_CREATE_URL", "http://10.190.4.122:8390/ws/api/services/documents"))
	cfg.AuctionImageCreateURL = cast.ToString(getOrReturnDefaultValue("IMAGE_CREATE_URL", "http://10.190.4.122:8390/ws/api/services/images"))
	cfg.AuctionImageUpdateURL = cast.ToString(getOrReturnDefaultValue("IMAGE_CREATE_URL", "http://10.190.4.122:8390/ws/api/services/images"))
	cfg.AuctionSendOrderURL = cast.ToString(getOrReturnDefaultValue("SEND_ORDER_URL", "http://10.190.4.122:8390/ws/api/services/common/send-order"))
	cfg.AuctionGetProtocolURL = cast.ToString(getOrReturnDefaultValue("GET_PROTOCOL_URL", "http://10.190.4.122:8390/ws/api/services/common/get-protocol"))
	cfg.AuctionUsername = cast.ToString(getOrReturnDefaultValue("AUCTION_USERNAME", "yerelektron@ygk.uz"))
	cfg.AuctionUsernameTypeSix = cast.ToString(getOrReturnDefaultValue("AUCTION_USERNAME_TYPE_SIX", "auga_republic"))
	cfg.AuctionPassword = cast.ToString(getOrReturnDefaultValue("AUCTION_PASSWORD", "yerelektron@ygk.uz"))
	cfg.AuctionPasswordTypeSix = cast.ToString(getOrReturnDefaultValue("AUCTION_PASSWORD_TYPE_SIX", "G2unR7Ng"))
	cfg.AuctionUrl = cast.ToString(getOrReturnDefaultValue("AUCTION_HOST", "http://10.190.4.122:8390"))
	cfg.EntityServiceHost = cast.ToString(getOrReturnDefaultValue("ENTITY_SERVICE_HOST", "localhost"))
	cfg.EntityServicePort = cast.ToInt(getOrReturnDefaultValue("ENTITY_SERVICE_PORT", 8004))
	cfg.UserServiceHost = cast.ToString(getOrReturnDefaultValue("USER_SERVICE_HOST", "localhost"))
	cfg.UserServicePort = cast.ToInt(getOrReturnDefaultValue("USER_SERVICE_PORT", 8002))
	cfg.DiscussionLogicServiceHost = cast.ToString(getOrReturnDefaultValue("DISCUSSION_LOGIC_SERVICE_HOST", "localhost"))
	cfg.DiscussionLogicServicePort = cast.ToInt(getOrReturnDefaultValue("DISCUSSION_LOGIC_SERVICE_PORT", 8003))

	cfg.HokimyatUrl = cast.ToString(getOrReturnDefaultValue("RAQAMLI_HOKIMYAT_HOST", "https://apigateway.digitaltashkent.uz/api/v1/yerElectron"))
	cfg.HokimyatTokenUrl = cast.ToString(getOrReturnDefaultValue("RAQAMLI_HOKIMYAT_TOKEN_HOST", "https://apigateway.digitaltashkent.uz/api/token/"))
	cfg.HokimyatUrlDownloadFile = cast.ToString(getOrReturnDefaultValue("RAQAMLI_HOKIMYAT_URL_DOWNLOAD_FILE", "https://api.admin.yerelektron.uz/v1/raqamli-hokimyat/download/entity-file"))
	return cfg
}

func getOrReturnDefaultValue(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)

	if exists {
		return os.Getenv(key)
	}

	return defaultValue
}

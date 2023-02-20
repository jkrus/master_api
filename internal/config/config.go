package config

import (
	"log"
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/jkrus/master_api/pkg/zap-logger/v6/fields"
)

type Config struct {
	// -----------------------------------------------------------------------------
	LogLevel               string        `envconfig:"LOG_LEVEL" default:"info" validate:"oneof=debug info warn error dpanic panic fatal"`
	ServerHost             string        `envconfig:"SERVER_HOST" default:"127.0.0.1"`
	ServerPort             string        `envconfig:"SERVER_PORT" default:"8080"`
	HTTPServerReadTimeOut  time.Duration `envconfig:"HTTP_SERVER_READ_TIMEOUT" default:"10m"`
	HTTPServerWriteTimeOut time.Duration `envconfig:"HTTP_SERVER_WRITE_TIMEOUT" default:"13m"`
	PayloadSoftLimit       int           `envconfig:"PAYLOAD_SAFE_LIMIT" default:"5120"`    // 5 КиБ
	PayloadHardLimit       int           `envconfig:"PAYLOAD_HARD_LIMIT" default:"5242880"` // 5 МиБ
	PayloadQuantityLimit   int           `envconfig:"PAYLOAD_QUANTITY_LIMIT" default:"10000"`

	// CacheRequestTTL - время жизни кэшируемого запроса
	CacheRequestTTL time.Duration `envconfig:"CACHE_REQUEST_TTL" default:"30m"`
	// CacheResponseTTL - время жизни кэшируемого ответа
	CacheResponseTTL time.Duration `envconfig:"CACHE_RESPONSE_TTL" default:"30m"`

	/*----- MinIO -----*/
	MinioEndPoint  string `envconfig:"MINIO_END_POINT" default:"172.18.0.2:9000"`
	MinioAccessKey string `envconfig:"MINIO_ACCESS_KEY" default:"minio"`
	MinioSecretKey string `envconfig:"MINIO_SECRET_KEY" default:"minio123"`

	// DB

	DBUser           string `envconfig:"DB_USER" default:"postgres"`
	DBPass           string `envconfig:"DB_PASS" default:"postgres"`
	DBHost           string `envconfig:"DB_HOST" default:"localhost"`
	DBPort           string `envconfig:"DB_PORT" default:"5432"`
	DBName           string `envconfig:"DB_NAME" default:"postgres"`
	DBSSLMode        string `envconfig:"DB_SSL_MODE" default:"disable" validate:"oneof=disable enable"`
	SQLSlowThreshold int    `envconfig:"SQL_SLOW_THRESHOLD" default:"600"`
	TraceSQLCommands bool   `envconfig:"TRACE_SQL_COMMANDS" default:"false"`
	AutoMigrate      bool   `envconfig:"AUTO_MIGRATE" default:"true"`

	// Hyper Ledger

	HFMSPID         string `envconfig:"HF_MSP_ID" default:"Org1MSP"`
	HFCryptoPath    string `envconfig:"HF_CRYPTO_PATH" default:"internal/hf_certs/organizations/peerOrganizations/org1.example.com"`
	HFCertPath      string `envconfig:"HF_CERT_PATH" default:"internal/hf_certs/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/signcerts/User1@org1.example.com-cert.pem"`
	HFKeyPath       string `envconfig:"HF_KEY_PATH" default:"internal/hf_certs/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/keystore/"`
	HFTLSCertPath   string `envconfig:"HF_TLS_CERT_PATH" default:"internal/hf_certs/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt"`
	HFPeerEndpoint  string `envconfig:"HF_PEER_ENDPOINT" default:"localhost:7051"`
	HFGatewayPeer   string `envconfig:"HF_GATEWAY_PEER" default:"peer0.org1.example.com"`
	HFChannelName   string `envconfig:"HF_CHANNEL_NAME" default:"mychannel"`
	HFChaincodeName string `envconfig:"HF_CHAINCODE_NAME" default:"file"`
}

func (config Config) PayloadConfig() fields.PayloadConfig {
	payloadConfig := fields.DefaultPayloadConfig
	payloadConfig.SoftLimit = config.PayloadSoftLimit
	payloadConfig.HardLimit = config.PayloadHardLimit
	payloadConfig.QuantityLimit = config.PayloadQuantityLimit

	return payloadConfig
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}
		err := envconfig.Process("server", instance)
		if err != nil {
			log.Fatal(err)
		}
	})
	return instance
}

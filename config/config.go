package config

import "github.com/multiversx/mx-chain-core-go/core"

// Configs holds all configs
type Configs struct {
	MainConfig      MainConfig
	ApiRoutesConfig APIRoutesConfig
	Flags           FlagsConfig
}

// MainConfig defines the config setup based on main config file
type MainConfig struct {
	General                    GeneralConfig
	WebSocketConnector         WebSocketConfig
	ConnectorApi               ConnectorApiConfig
	Redis                      RedisConfig
	RabbitMQ                   RabbitMQConfig
	ExternalWebSocketConnector WebSocketConfig
}

// GeneralConfig maps the general config section
type GeneralConfig struct {
	ExternalMarshaller MarshallerConfig
	AddressConverter   AddressConverterConfig
	CheckDuplicates    bool
}

// MarshallerConfig maps the marshaller configuration
type MarshallerConfig struct {
	Type string
}

// AddressConverterConfig maps the address pubkey converter configuration
type AddressConverterConfig struct {
	Type   string
	Prefix string
	Length int
}

// ConnectorApiConfig maps the connector configuration
type ConnectorApiConfig struct {
	Enabled  bool
	Host     string
	Username string
	Password string
}

// APIRoutesConfig holds the configuration related to Rest API routes
type APIRoutesConfig struct {
	APIPackages map[string]APIPackageConfig
}

// APIPackageConfig holds the configuration for the routes of each package
type APIPackageConfig struct {
	Routes []RouteConfig
}

// RouteConfig holds the configuration for a single route
type RouteConfig struct {
	Name string
	Open bool
	Auth bool
}

// RedisConfig maps the redis configuration
type RedisConfig struct {
	Url            string
	MasterName     string
	SentinelUrl    string
	ConnectionType string
	TTL            uint32
}

// RabbitMQConfig maps the rabbitMQ configuration
type RabbitMQConfig struct {
	Url                     string
	EventsExchange          RabbitMQExchangeConfig
	RevertEventsExchange    RabbitMQExchangeConfig
	FinalizedEventsExchange RabbitMQExchangeConfig
	BlockTxsExchange        RabbitMQExchangeConfig
	BlockScrsExchange       RabbitMQExchangeConfig
	BlockEventsExchange     RabbitMQExchangeConfig
}

// RabbitMQExchangeConfig holds the configuration for a rabbitMQ exchange
type RabbitMQExchangeConfig struct {
	Name string
	Type string
}

// WebSocketConfig holds the configuration for websocket observer interaction config
type WebSocketConfig struct {
	Enabled                    bool
	URL                        string
	Mode                       string
	RetryDurationInSec         int
	AcknowledgeTimeoutInSec    int
	WithAcknowledge            bool
	BlockingAckOnError         bool
	DropMessagesIfNoConnection bool

	DataMarshallerType string
}

// FlagsConfig holds the values for CLI flags
type FlagsConfig struct {
	LogLevel          string
	SaveLogFile       bool
	GeneralConfigPath string
	APIConfigPath     string
	WorkingDir        string
	PublisherType     string
	RestApiInterface  string
}

// LoadMainConfig returns a MainConfig instance by reading the provided toml file
func LoadMainConfig(filePath string) (*MainConfig, error) {
	cfg := &MainConfig{}
	err := core.LoadTomlFile(cfg, filePath)
	if err != nil {
		return nil, err
	}
	return cfg, err
}

// LoadAPIConfig returns a APIRoutesConfig instance by reading the provided toml file
func LoadAPIConfig(filePath string) (*APIRoutesConfig, error) {
	cfg := &APIRoutesConfig{}
	err := core.LoadTomlFile(cfg, filePath)
	if err != nil {
		return nil, err
	}
	return cfg, err
}

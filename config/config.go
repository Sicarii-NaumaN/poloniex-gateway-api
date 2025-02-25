// Config package uses for getting env params

package config

type configKey string

const (
	LogLevel = configKey("")
	Port     = configKey("PORT")

	DBDSN = configKey("DB_DSN")

	Nodes            = configKey("NODES")
	ChainId          = configKey("CHAIN_ID")
	IsSwaggerCreated = configKey("IS_SWAGGER_CREATED")
)

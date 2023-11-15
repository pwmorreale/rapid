package config

// Constants for various configuration parameters.
const (

	// Viper configuration.
	ViperConfigFileName = "rapid"
	ViperConfigPath     = "/etc/rapid"
	ViperConfigFileType = "yaml"
	ViperConfigEnv      = "rapid_CONFIG_PATH"

	// REST Server configuration
	ServerPort    = "server.port"
	ServerAddress = "server.address"

	LogLevel = "log_level"
)

// DefaultFile builds a default path to the config file.
func DefaultFile() string {
	return ViperConfigPath + "/" + ViperConfigFileName + "." + ViperConfigFileType
}

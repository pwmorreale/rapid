//
//  Copyright © 2024 Peter W. Morreale. All Rights Reserved.
//

// Package config contains config variables.and utilities
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

	// Scenario configuration
	ScenarioName    = "name"
	ScenarioVersion = "version"
)

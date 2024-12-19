//
//  Copyright © 2024 Peter W. Morreale. All Rights Reserved.
//

// Package cmd defines the Cobra root
package cmd

import (
	"errors"
	"fmt"

	"github.com/pwmorreale/rapid/internal/config"
	"github.com/pwmorreale/rapid/internal/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "The server command starts the rapid web server.",
	Long:  `The server command starts the rapid web server.`,
	RunE:  serverRun,
}

func init() {
	defaultFile := config.ViperConfigPath + "/" + config.ViperConfigFileName + "." + config.ViperConfigFileType
	serverCmd.PersistentFlags().StringVarP(&cfgFile, "conf", "c", "", "config file (default is "+defaultFile+")")
	serverCmd.MarkFlagRequired("conf")
	serverCmd.MarkFlagFilename("conf")
}

// serverRun creates and runs a REST server instance
func serverRun(_ *cobra.Command, _ []string) error {

	fmt.Println("running in SERVER")

	addr := viper.GetString(config.ServerAddress) + ":" + viper.GetString(config.ServerPort)
	if addr == "" {
		return errors.New("missing server address/port")
	}

	sv := server.New(addr)
	return (sv.Start())
}

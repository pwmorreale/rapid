/*
Copyright © 2023 Peter W. Morreale

*/
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

// serverRun creates and runs a REST server instance
func serverRun(cmd *cobra.Command, args []string) error {

	fmt.Println("running in SERVER")

	addr := viper.GetString(config.ServerAddress) + ":" + viper.GetString(config.ServerPort)
	if addr == "" {
		return errors.New("Missing server address/port")
	}

	sv := server.New(addr)
	return (sv.Start())
}

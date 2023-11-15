/*
Copyright © 2023 Peter W. Morreale

*/

// Package cmd defines the Cobra root.
package cmd

import (
	"fmt"

	"github.com/pwmorreale/rapid/internal/config"
	"github.com/pwmorreale/rapid/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile  string
	tstCases []string
	rootCmd  = &cobra.Command{
		Use:   "rapid",
		Short: "rapid is an REST API diagnostic tool.",
		Long:  `rapid enables testing of REST APIs.`,
		RunE:  RunRoot,
	}
)

// rootCmd represents the base command when called without any subcommands

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

var log = logger.GetLogger("rapid")

func init() {
	cobra.OnInitialize(initConfig)

	defaultFile := config.DefaultFile()

	rootCmd.PersistentFlags().StringVar(&cfgFile, "rapid", "", "config file (default is "+defaultFile+")")
	rootCmd.PersistentFlags().StringArrayVarP(&tstCases, "testcase", "t", nil, "Path to testcase configs, can specify multiple times. All viper file extensions are supported.")

	rootCmd.AddCommand(serverCmd)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {

		// Search config in default directory with default name.
		viper.AddConfigPath(config.ViperConfigPath)
		viper.SetConfigType(config.ViperConfigFileType)
		viper.SetConfigName(config.ViperConfigFileName)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal().Err(err).Msg("Failed to read config")
	} else {
		log.Info().Str("Using config file:", viper.ConfigFileUsed())
	}
}

// RunRoot executes the CLI interface.
func RunRoot(cmd *cobra.Command, args []string) error {

	fmt.Println("test cases: ", tstCases)
	return nil
}

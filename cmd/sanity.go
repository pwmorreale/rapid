//
//  Copyright © 2024 Peter W. Morreale. All Rights Reserved.
//

// Package cmd defines the Cobra root
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// sanityCmd represents the server command
var sanityCmd = &cobra.Command{
	Use:   "sanity",
	Short: "Sanity check the specified scenario configuration",
	Long:  "Sanity check the specified scenario configuration",
	RunE:  runSanity,
}

// runSanity creates and runs a REST server instance
func runSanity(_ *cobra.Command, _ []string) error {

	fmt.Println("running in SANITY")

	return nil
}

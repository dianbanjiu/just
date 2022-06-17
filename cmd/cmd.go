package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"just/cmd/subscription"
	"just/cmd/traffic"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "just",
	Short: "just is a util for justmysocks",
	Long:  `just is a util for justmysocks, can read subscription, and get traffic usage`,
}

func Execute() {
	rootCmd.AddCommand(traffic.TrafficCmd)
	rootCmd.AddCommand(subscription.SubCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

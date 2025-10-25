/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/latocchi/gomailit/internal/providers"
	"github.com/spf13/cobra"
)

var provider string

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup email provider. Default is google",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Unsupported provider:", provider)
			fmt.Println("Switching to default provider 'google'")
			_, err := providers.SetupGoogle()
			if err != nil {
				fmt.Println("Error setting up Google provider:", err)
			}
			return
		}

		provider = args[0]

		switch provider {
		case "google", "gmail":
			_, err := providers.SetupGoogle()
			if err != nil {
				fmt.Println("Error setting up Google provider:", err)
			}
		// TODO: Add other providers here
		default:
			fmt.Println("Unsupported provider:", provider)
			fmt.Println("Switching to default provider 'google'")
			_, err := providers.SetupGoogle()
			if err != nil {
				fmt.Println("Error setting up Google provider:", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

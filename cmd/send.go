/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/latocchi/gomailit/internal/providers"
	"github.com/latocchi/gomailit/internal/utils"
	"github.com/spf13/cobra"
)

var (
	to      string
	body    string
	subject string
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if utils.IsFile(body) {
			data, err := os.ReadFile(body)
			if err != nil {
				panic(err)
			}
			body = string(data)
		}

		err := providers.SendEmailGMail(to, subject, body)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sendCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sendCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	sendCmd.Flags().StringVarP(&to, "to", "t", "", "Recipient of the email")
	sendCmd.Flags().StringVarP(&body, "body", "b", "No body", "Body of the email, can be '-' for stdin or a .txt file path")
	sendCmd.Flags().StringVarP(&subject, "subject", "s", "No subject", "Subject of the email")

	if err := sendCmd.MarkFlagRequired("to"); err != nil {
		panic(err)
	}
}

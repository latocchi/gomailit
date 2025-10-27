/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/latocchi/gomailit/internal/providers"
	"github.com/latocchi/gomailit/internal/utils"
	"github.com/spf13/cobra"
)

var (
	to          string
	body        string
	subject     string
	attachments []string
	recipients  []string
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Sends email using the configured provider",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !utils.IsFile(utils.TokenPath()) {
			fmt.Println("No token found, please run 'gomailit setup google' first to set up the Google provider.")
			os.Exit(1)
		}

		if cmd.Flags().Changed("attach") {
			attachments = args
		}

		// fmt.Println(attachments)
		if len(attachments) > 0 {
			// Requirements:
			// --attach file1 file2
			// --attach ~/directory/subfolder/*

			// Go through attachments and remove any that do not exist
			for index, path := range attachments {
				// if file does not exist, remove from attachments slice
				fmt.Println("Checking attachment:", path)
				if !utils.FileExists(path) {
					fmt.Printf("Attachment file not found: %s\n", path)
					attachments = append(attachments[:index], attachments[index+1:]...)
					break
				}

			}
		}

		if utils.FileExists(body) {
			data, err := os.ReadFile(body)
			if err != nil {
				panic(err)
			}
			body = string(data)
		}

		// Check if 'to' is a file with multiple recipients
		if utils.FileExists(to) {
			file, err := os.Open(to)
			if err != nil {
				panic(err)
			}

			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line != "" {
					recipients = append(recipients, line)
				}
			}

			if err := scanner.Err(); err != nil {
				log.Fatalf("failed to read file: %v", err)
			}

			for _, recipient := range recipients {
				// Should i use goroutines here?
				err := providers.SendEmailGMail(recipient, subject, body, attachments)
				if err != nil {
					panic(err)
				}
			}
		} else { // Single recipient
			err := providers.SendEmailGMail(to, subject, body, attachments)
			if err != nil {
				panic(err)
			}
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
	sendCmd.Flags().BoolP("attach", "a", false, "File paths to attach to the email")
	if err := sendCmd.MarkFlagRequired("to"); err != nil {
		panic(err)
	}

}

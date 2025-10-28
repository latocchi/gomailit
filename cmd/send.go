/*
Copyright © 2025 Jaycy Ivan Bañaga jaycybanaga@gmail.com
*/
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

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
	Long: `Usage:

Basic Send
gomailit send --from alice@example.com --to bob@example.com \
	--subject "Hello" --body "This is a test"

Examples:

Send with attachment
gomailit send --to bob@example.com --subject "Files" --body "See attached" \
	--attach report.pdf

Send multiple attachments
gomailit send --to bob@example.com --subject "Files" --body "See attached" \ 
	--attach report.pdf agenda.pdf

Attach all files from a directory
gomailit send --to bob@example.com --subject "Files" --body "See attached" \ 
	--attach ~/Documents/report/*

Use file for email body
gomailit send --to bob@example.com --subject "Files" --body ~/Documents/body.txt \
	--attach ~/Documents/report/*

Send to multiple recipients via .txt file
gomailit send --to ~/Documents/recipients.txt --subject "Files" \ 
	--body ~/Documents/body.txt --attach ~/Documents/report/*

Example contents of recipients.txt file:
recipient@example.com
recipient1@example.com
recipient2@example.com
recipient3@example.com
recipient4@example.com
recipient5@example.com
`,
	Run: func(cmd *cobra.Command, args []string) {
		if !utils.IsFile(utils.TokenPath()) {
			fmt.Println("No token found, please run 'gomailit setup google' first to set up the Google provider.")
			os.Exit(1)
		}

		srv, err := providers.GetGoogleService()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to get google mail service: %v", err)
		}

		profile, err := srv.Users.GetProfile("me").Do()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to get user profile: %v", err)
		}

		fmt.Printf("Sending email as %s\n", profile.EmailAddress)

		if cmd.Flags().Changed("attach") {
			attachments = args
		}

		// fmt.Println(attachments)
		if len(attachments) > 0 {
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

			sem := make(chan struct{}, 5) // limit to 5 concurrent goroutines
			var wg sync.WaitGroup

			for _, recipient := range recipients {
				// Should i use goroutines here?
				wg.Add(1)
				sem <- struct{}{}

				go func(recipient string) {
					defer wg.Done()
					defer func() { <-sem }()

					if err := providers.SendEmailGMail(recipient, subject, body, attachments); err != nil {
						fmt.Printf("Failed to send email to %s: %v\n", recipient, err)
					} else {
						fmt.Printf("Email sent to %s successfully.\n", recipient)
					}
				}(recipient)
			}
			wg.Wait()
			fmt.Println("All emails sent.")
		} else { // Single recipient
			if err := providers.SendEmailGMail(to, subject, body, attachments); err != nil {
				fmt.Printf("Failed to send email to %s: %v\n", to, err)
			} else {
				fmt.Printf("Email sent to %s successfully.\n", to)
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
	// sendCmd.Flags().BoolP("toggle", "p", false, "Help message for toggle")

	sendCmd.Flags().StringVarP(&to, "to", "t", "", "Recipient (single email or .txt file with one address per line)")
	sendCmd.Flags().StringVarP(&body, "body", "b", "No body", "Body of the email, can be '-' for stdin or a .txt file path (default \"No body\")")
	sendCmd.Flags().StringVarP(&subject, "subject", "s", "No subject", "Subject of the email (default \"No subject\")")
	sendCmd.Flags().BoolP("attach", "a", false, "One or more attachment files, or a directory path")
	if err := sendCmd.MarkFlagRequired("to"); err != nil {
		panic(err)
	}

}

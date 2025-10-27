/*
Copyright © 2025 Jaycy Ivan Bañaga jaycybanaga@gmail.com
*/
package providers

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"mime/quotedprintable"
	"net"
	"net/http"
	"net/textproto"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/latocchi/gomailit/internal/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type AttachmentPart struct {
	Header textproto.MIMEHeader
	Body   []byte
	Err    error
}

func buildMessageWithAttachments(to, subject, body string, attachments []string) *gmail.Message {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	boundary := writer.Boundary()

	headers := fmt.Sprintf(
		"To: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: multipart/mixed; boundary=%s\r\n\r\n",
		to, subject, boundary,
	)
	buf.WriteString(headers)

	bodyHeader := textproto.MIMEHeader{}
	bodyHeader.Set("Content-Type", "text/plain; charset=\"UTF-8\"")
	bodyHeader.Set("Content-Transfer-Encoding", "quoted-printable")

	bodyPart, _ := writer.CreatePart(bodyHeader)
	qp := quotedprintable.NewWriter(bodyPart)
	qp.Write([]byte(body))
	qp.Close()

	var wg sync.WaitGroup
	results := make(chan AttachmentPart, len(attachments))

	for _, path := range attachments {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			data, err := os.ReadFile(path)
			if err != nil {
				results <- AttachmentPart{Err: err}
				return
			}

			encoded := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
			base64.StdEncoding.Encode(encoded, data)

			filename := filepath.Base(path)

			attachmentHeader := textproto.MIMEHeader{}
			attachmentHeader.Set("Content-Type", fmt.Sprintf("application/octet-stream; name=\"%s\"", filename))
			attachmentHeader.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
			attachmentHeader.Set("Content-Transfer-Encoding", "base64")

			results <- AttachmentPart{Header: attachmentHeader, Body: encoded, Err: nil}
		}(path)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for part := range results {
		if part.Err != nil {
			panic(part.Err)
		}

		w, _ := writer.CreatePart(part.Header)

		for i := 0; i < len(part.Body); i += 76 {
			end := i + 76
			if end > len(part.Body) {
				end = len(part.Body)
			}
			w.Write(part.Body[i:end])
			w.Write([]byte("\r\n"))
		}
	}

	writer.Close()
	raw := utils.EncodeURLSafeBase64(buf.Bytes())
	return &gmail.Message{Raw: raw}

}

func buildMessage(to, subject, body string) *gmail.Message {
	message := []byte(
		fmt.Sprintf("To: %s\r\n", to) +
			fmt.Sprintf("Subject: %s\r\n", subject) +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=\"utf-8\"\r\n\r\n" +
			body,
	)
	raw := utils.EncodeURLSafeBase64(message)
	return &gmail.Message{Raw: raw}
}

func SendEmailGMail(to, subject, body string, attachments []string) error {
	srv, err := GetGoogleService()
	if err != nil {
		return fmt.Errorf("unable to get google mail service: %v", err)
	}

	if len(attachments) > 0 {
		mail := buildMessageWithAttachments(to, subject, body, attachments)
		err = send(srv, mail)
		if err != nil {
			return err
		}
		// fmt.Println("Email with attachments sent successfully to " + to + "!")
		return nil
	}

	mail := buildMessage(to, subject, body)

	err = send(srv, mail)
	if err != nil {
		return err
	}
	return nil
}

func send(srv *gmail.Service, mail *gmail.Message) error {
	_, err := srv.Users.Messages.Send("me", mail).Do()
	if err != nil {
		return fmt.Errorf("unable to send email: %v", err)
	}
	return nil
}

func GetGoogleService() (*gmail.Service, error) {
	client, err := SetupGoogle()
	if err != nil {
		return nil, fmt.Errorf("unable to setup google client: %v", err)
	}

	srv, err := gmail.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Gmail client: %v", err)
	}

	return srv, nil
}

func SetupGoogle() (*http.Client, error) {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailSendScope, gmail.GmailMetadataScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %v", err)
	}

	client := getClient(config)

	return client, nil
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := utils.TokenPath()
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	// Create local server to listen for redirect
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("Unable to start local server: %v", err)
	}
	defer listener.Close()

	config.RedirectURL = "http://localhost:8080"
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Printf("Your browser will open for authorization.\nIf it doesn't, open this link manually:\n%v\n", authURL)
	exec.Command("xdg-open", authURL).Start()

	codeChan := make(chan string)

	go func() {
		conn, _ := listener.Accept()
		defer conn.Close()

		req, _ := http.ReadRequest(bufio.NewReader(conn))
		code := req.URL.Query().Get("code")

		conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/html\r\n\r\nYou may now close this window."))
		codeChan <- code
	}()

	code := <-codeChan
	tok, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

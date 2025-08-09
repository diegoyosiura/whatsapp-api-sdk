package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp"
)

// envTokenProvider lê WA_ACCESS_TOKEN do ambiente.
type envTokenProvider struct{}

func (envTokenProvider) Token(ctx context.Context) (string, error) {
	t := os.Getenv("WA_ACCESS_TOKEN")
	if t == "" {
		return "", fmt.Errorf("WA_ACCESS_TOKEN not set")
	}
	return t, nil
}
func (envTokenProvider) Refresh(ctx context.Context) error { return nil }

func main() {
	to := flag.String("to", "", "Recipient in E.164 (e.g., +5511999999999)")
	body := flag.String("body", "Hello from Go SDK!", "Message text body")
	flag.Parse()

	if *to == "" {
		log.Fatal("-to is required")
	}

	opts := whatsapp.Options{
		Version:       getenvDefault("WA_GRAPH_VERSION", "v20.0"),
		WABAID:        mustGetenv("WA_WABA_ID"),
		PhoneNumberID: mustGetenv("WA_PHONE_NUMBER_ID"),
		TokenProvider: envTokenProvider{},
		Timeout:       15 * time.Second,
		RetryMax:      3,
		UserAgent:     "examples/send_text",
	}

	client, err := whatsapp.NewClient(opts)
	if err != nil {
		log.Fatalf("new client: %v", err)
	}

	ctx := context.Background()
	resp, err := client.Messages.SendText(ctx, *to, *body)
	if err != nil {
		log.Fatalf("send text: %v", err)
	}

	if len(resp.Messages) > 0 {
		fmt.Printf("sent ok — message_id=%s wa_id=%s\n",
			resp.Messages[0].ID,
			resp.Contacts[0].WaID,
		)
	} else {
		fmt.Printf("sent ok — response: %+v\n", resp)
	}
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("env %s not set", k)
	}
	return v
}

func getenvDefault(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

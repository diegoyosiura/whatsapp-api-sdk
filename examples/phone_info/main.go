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
	// Flag opcional para consultar um ID específico depois de listar
	getID := flag.String("get", "", "Fetch details for this Phone Number ID after listing")
	flag.Parse()

	opts := whatsapp.Options{
		Version:       getenvDefault("WA_GRAPH_VERSION", "v20.0"),
		WABAID:        mustGetenv("WA_WABA_ID"),
		PhoneNumberID: mustGetenv("WA_PHONE_NUMBER_ID"),
		TokenProvider: envTokenProvider{},
		Timeout:       10 * time.Second,
		RetryMax:      3,
		UserAgent:     "examples/phone_info",
	}

	client, err := whatsapp.NewClient(opts)
	if err != nil {
		log.Fatalf("new client: %v", err)
	}

	ctx := context.Background()

	// Lista números do WABA
	list, err := client.Phone.List(ctx)
	if err != nil {
		log.Fatalf("phone list: %v", err)
	}
	fmt.Println("ID\tDISPLAY\tQUALITY")
	for _, pn := range list.Data {
		fmt.Printf("%s\t%s\t%s\n", pn.ID, pn.DisplayPhoneNumber, pn.QualityRating)
	}

	// Consulta detalhada opcional
	if *getID != "" {
		p, err := client.Phone.Get(ctx, *getID)
		if err != nil {
			log.Fatalf("phone get %s: %v", *getID, err)
		}
		fmt.Printf("\nDetails for %s:\n", *getID)
		fmt.Printf("VerifiedName: %s\n", p.VerifiedName)
		fmt.Printf("Display:     %s\n", p.DisplayPhoneNumber)
		fmt.Printf("Quality:     %s\n", p.QualityRating)
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

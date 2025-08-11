package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/domain"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/ports"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/services"

	"log"
	"os"
	"time"
)

type whSecretProvider struct{}

func (w whSecretProvider) Get(ctx context.Context, key ports.SecretKey) (string, error) {
	//TODO implement me
	panic("implement me")
}

type envToken struct{}

func (envToken) Token(ctx context.Context) (string, error) { return os.Getenv("WA_ACCESS_TOKEN"), nil }
func (envToken) Refresh(ctx context.Context) error         { return nil }

func main() {
	cmd := flag.String("cmd", "", "send-text|phone-list|request-code|verify-code|register")
	to := flag.String("to", "", "E.164 number")
	body := flag.String("body", "Hello", "message body")
	locale := flag.String("locale", "en_US", "locale for code")
	code := flag.String("code", "", "verification code")
	pin := flag.String("pin", "", "two-step PIN")
	flag.Parse()

	_ = services.NewWebhookService(whSecretProvider{})
	c, err := whatsapp.NewClient(whatsapp.Options{
		Version:       os.Getenv("WA_GRAPH_VERSION"),
		WABAID:        os.Getenv("WA_WABA_ID"),
		PhoneNumberID: os.Getenv("WA_PHONE_NUMBER_ID"),
		TokenProvider: envToken{},
		Timeout:       10 * time.Second,
		RetryMax:      3,
		UserAgent:     "cli",
	})
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	switch *cmd {
	case "send-text":
		resp, err := c.Messages.SendText(ctx, *to, *body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ok id=%s\n", resp.Messages[0].ID)
	case "phone-list":
		out, err := c.Phone.List(ctx)
		if err != nil {
			log.Fatal(err)
		}
		for _, p := range out.Data {
			fmt.Println(p.ID, p.DisplayPhoneNumber, p.QualityRating)
		}
	case "request-code":
		out, err := c.Registration.RequestCode(ctx, domain.RequestCodeParams{CodeMethod: domain.CodeMethodSMS, Locale: *locale})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("success:", out.Success)
	case "verify-code":
		out, err := c.Registration.VerifyCode(ctx, domain.VerifyCodeParams{Code: *code})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("success:", out.Success)
	case "register":
		var p *string
		if *pin != "" {
			p = pin
		}
		out, err := c.Registration.Register(ctx, domain.RegisterParams{Pin: p})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("success:", out.Success)
	default:
		log.Fatal("missing or invalid -cmd")
	}
}

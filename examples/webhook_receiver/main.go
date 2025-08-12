package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/domain"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/ports"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/services"
	"github.com/diegoyosiura/whatsapp-sdk-go/pkg/whatsapp/transport/webhook"
)

type fileMan struct {
}

func (f fileMan) SetData(ctx context.Context, data []byte) error {
	//TODO implement me
	panic("implement me")
}

func (f fileMan) Save(ctx context.Context, fileName string) error {
	//TODO implement me
	panic("implement me")
}

func (f fileMan) Open(ctx context.Context, fileName string) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

// Minimal SecretsProvider backed by environment variables for the example only.
type envSecrets struct{}

func (s envSecrets) Get(ctx context.Context, key ports.SecretKey) (string, error) {
	v := os.Getenv(string(key))
	if v == "" {
		return "", fmt.Errorf("secret %s not set", key)
	}
	return v, nil
}

// logHandler logs webhook events; replace with your business logic.
type logHandler struct {
	c *whatsapp.Client
}

func (lh logHandler) OnMessage(m domain.InboundMessage) {
	log.Printf("incoming message id=%s", m.ID)
	fm, err := lh.c.Media.GetMedia(context.Background(), m, fileMan{})
	if err != nil {
		log.Println(err)
	}
	log.Printf("media info: %+v", fm)
}
func (lh logHandler) OnStatus(s domain.MessageStatus) {
	log.Printf("status id=%s status=%s", s.ID, s.Status)
}

type envToken struct{}

func (envToken) Token(ctx context.Context) (string, error) { return os.Getenv("WA_ACCESS_TOKEN"), nil }
func (envToken) Refresh(ctx context.Context) error         { return nil }

func main() {
	secrets := envSecrets{}
	svc := services.NewWebhookService(secrets)

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
		panic(err)

	}
	dispatcher := services.NewWebhookDispatcher(logHandler{c: c})
	http.Handle("/webhook", webhook.NewHandler(svc, dispatcher))

	addr := ":8080"
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

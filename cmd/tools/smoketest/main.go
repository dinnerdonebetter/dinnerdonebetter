package main

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"go.opentelemetry.io/otel/trace"

	"github.com/prixfixeco/api_server/pkg/client/httpclient"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

const (
	stagingAddress = "https://api.prixfixe.dev"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	example := fakes.BuildFakeUserRegistrationInput()
	parsedURI, err := url.Parse(stagingAddress)
	if err != nil {
		panic(err)
	}

	input := &types.UserRegistrationInput{
		Username:     example.Username,
		Password:     example.Password,
		EmailAddress: "verygoodsoftwarenotvirus@protonmail.com",
	}

	user, err := testutils.CreateServiceUser(ctx, stagingAddress, input)
	if err != nil {
		panic(err)
	}

	cookie, err := testutils.GetLoginCookie(ctx, stagingAddress, user)
	if err != nil {
		panic(err)
	}

	client, err := httpclient.NewClient(parsedURI, trace.NewNoopTracerProvider(), httpclient.UsingCookie(cookie))
	if err != nil {
		panic(err)
	}

	createdWebhook, err := client.CreateWebhook(ctx, fakes.BuildFakeWebhookCreationInput())
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second)

	webhook, err := client.GetWebhook(ctx, createdWebhook.ID)
	if err != nil {
		panic(err)
	}

	fmt.Println(webhook)
}

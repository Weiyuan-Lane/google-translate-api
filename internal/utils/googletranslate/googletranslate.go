package googletranslate

import (
	"context"

	"cloud.google.com/go/translate"
	translatev3 "cloud.google.com/go/translate/apiv3"
	"google.golang.org/api/option"
)

func InitTranslateV2Client(apiKey string) *translate.Client {
	ctx := context.Background()

	apiKeyOption := option.WithAPIKey(apiKey)

	client, err := translate.NewClient(ctx, apiKeyOption)
	if err != nil {
		panic(err)
	}

	return client
}

func CloseV2Client(client *translate.Client) {
	if err := client.Close(); err != nil {
		panic(err)
	}
}

func InitTranslateV3Client() *translatev3.TranslationClient {
	ctx := context.Background()

	client, err := translatev3.NewTranslationClient(ctx)
	if err != nil {
		panic(err)
	}

	return client
}

func CloseV3Client(client *translatev3.TranslationClient) {
	if err := client.Close(); err != nil {
		panic(err)
	}
}

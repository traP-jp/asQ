package clients

import (
	"context"
	"fmt"
	"os"
	"sync"

	traq "github.com/traPtitech/go-traq"
)

var (
	clientInstance *traq.APIClient
	clientOnce     sync.Once
	ctx            context.Context
	ctxOnce        sync.Once
)

func GetTraqClient() *traq.APIClient {
	clientOnce.Do(func() {
		fmt.Println("traQ client initialized")
		clientInstance = traq.NewAPIClient(traq.NewConfiguration())
	})
	return clientInstance
}

func GetTraqContext() context.Context {
	ctxOnce.Do(func() {
		fmt.Println("traQ context initialized")
		accessToken := os.Getenv("TRAQ_BOT_ACCESS_TOKEN")
		fmt.Println(len(accessToken))
		ctx = context.WithValue(context.Background(), traq.ContextAccessToken, accessToken)
	})
	return ctx
}

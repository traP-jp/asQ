package clients

import (
	"fmt"
	"sync"

	traq "github.com/traPtitech/go-traq"
)

var (
	clientInstance *traq.APIClient
	once           sync.Once
)

func GetTraqClient() *traq.APIClient {
	once.Do(func() {
		fmt.Println("traQ client initialized")
		clientInstance = traq.NewAPIClient(traq.NewConfiguration())
	})
	return clientInstance
}

package main

import (
	"github.com/traP-jp/h25s_05/trap_mcp/clients"
)

func main() {
	traq_client := clients.GetTraqClient()
	client2 := clients.GetTraqClient()
	if traq_client == client2 {

	}
}

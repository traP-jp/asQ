package llm

import "iter"

type Service interface {
	AskQuestion(question string, instruction string, previousResponseID string) Stream
}

type Stream interface {
	ID() string
	Next() iter.Seq[Response]
	Close() error
	Err() error
}

type Response struct {
	Text string
}

type MCP struct {
	ServerLabel string
	ServerURL   string
	Header      map[string]string // Optional headers for the request
}

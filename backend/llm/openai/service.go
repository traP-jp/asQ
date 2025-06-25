package openai

import (
	"context"
	"errors"
	"iter"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/packages/ssestream"
	"github.com/openai/openai-go/responses"
	"github.com/traP-jp/h25s_05/backend/llm"
)

type Service struct {
	client *openai.Client
	mcps   []llm.MCP // List of available MCPs
}

var (
	ErrNoDataFound = errors.New("no data found for the given id")
)

func NewService(mcps []llm.MCP) *Service {
	client := openai.NewClient(option.WithBaseURL("https://llm-proxy.trap.jp"))
	return &Service{
		client: &client,
		mcps:   mcps,
	}
}

// AskQuestion sends a question to the LLM and returns the response id.
func (s *Service) AskQuestion(question string, instruction string, previousResponseID string) llm.Stream {
	var previousID param.Opt[string]
	if previousResponseID != "" {
		previousID = param.NewOpt(previousResponseID)
	}

	tools := make([]responses.ToolUnionParam, 0, len(s.mcps))
	for _, mcp := range s.mcps {
		tools = append(tools, responses.ToolUnionParam{
			OfMcp: &responses.ToolMcpParam{
				ServerLabel: mcp.ServerLabel,
				ServerURL:   mcp.ServerURL,
				Headers:     mcp.Header,
				RequireApproval: responses.ToolMcpRequireApprovalUnionParam{
					OfMcpToolApprovalSetting: param.NewOpt("never"),
				},
			},
		})
	}
	stream := s.client.Responses.NewStreaming(context.TODO(), responses.ResponseNewParams{
		Input: responses.ResponseNewParamsInputUnion{
			OfString: param.NewOpt(question),
		},
		Instructions:       param.NewOpt(instruction),
		Store:              param.NewOpt(true),
		PreviousResponseID: previousID,
		Tools:              tools,
		Model:              "gpt-4o",
	})

	return &Stream{
		s: stream,
	}
}

type Stream struct {
	s   *ssestream.Stream[responses.ResponseStreamEventUnion]
	id  string
	err error
}

func (s *Stream) ID() string {
	if s.id != "" {
		return s.id
	}
	for s.s.Next() {
		if event, ok := s.s.Current().AsAny().(responses.ResponseCreatedEvent); ok {
			s.id = event.Response.ID
			return s.id
		}
	}
	return ""
}

func (s *Stream) Next() iter.Seq[llm.Response] {
	return func(yield func(llm.Response) bool) {
		for s.s.Next() {
			event := s.s.Current()

			switch e := event.AsAny().(type) {
			case responses.ResponseCreatedEvent:
				s.id = e.Response.ID // Store the response ID for later use
				continue
			case responses.ResponseTextDeltaEvent:
				if !yield(llm.Response{Text: e.Delta}) {
					return
				}
			default:
				// Ignore other event types
				continue
			}
		}

		if err := s.s.Err(); err != nil {
			s.err = err
			return
		}
	}
}

func (s *Stream) Err() error {
	if s.err != nil {
		return s.err
	}
	if s.s.Err() != nil {
		return s.s.Err()
	}
	if s.id == "" {
		return ErrNoDataFound
	}
	return nil
}

func (s *Stream) Close() error {
	if s.s != nil {
		return s.s.Close()
	}
	return nil
}

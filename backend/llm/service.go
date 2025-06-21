package llm

import (
	"context"
	"errors"
	"log/slog"
	"slices"
	"sync"

	"github.com/google/uuid"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/packages/ssestream"
	"github.com/openai/openai-go/responses"
)

type Service struct {
	client      *openai.Client
	data        map[uuid.UUID][]StreamData
	mu          sync.Mutex
	streamCh    chan StreamData // Channel to receive stream data
	subscribers []subscriber
}

type subscriber struct {
	id uuid.UUID
	ch chan StreamData
}

type StreamData struct {
	ID        uuid.UUID
	TextDelta string
	Err       error
}

type MCP struct {
	ServerLabel string
	ServerURL   string
	Header      map[string]string // Optional headers for the request
}

func NewService() *Service {
	client := openai.NewClient(option.WithBaseURL("https://llm-proxy.trap.jp"))
	return &Service{
		client:   &client,
		streamCh: make(chan StreamData, 100), // Buffered channel to handle stream data
		data:     make(map[uuid.UUID][]StreamData),
	}
}

func (s *Service) Run() {
	for data := range s.streamCh {
		s.mu.Lock()
		s.data[data.ID] = append(s.data[data.ID], data)
		s.mu.Unlock()

		s.mu.Lock()
		for _, sub := range s.subscribers {
			if sub.id == data.ID {
				sub.ch <- data
			}
		}
		s.mu.Unlock()
	}
}

func (s *Service) publishData(data StreamData) {
	s.streamCh <- data
}

func (s *Service) closeStream(id uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, id)
	for _, sub := range s.subscribers {
		if sub.id == id {
			close(sub.ch)
		}
	}
	s.subscribers = slices.DeleteFunc(s.subscribers, func(sub subscriber) bool {
		return sub.id == id
	})
}

func (s *Service) Subscribe(ctx context.Context, id uuid.UUID) (<-chan StreamData, error) {
	ch := make(chan StreamData)
	s.mu.Lock()
	if _, exists := s.data[id]; !exists {
		return nil, errors.New("no data found for the given id")
	}
	s.subscribers = append(s.subscribers, subscriber{id: id, ch: ch})
	s.mu.Unlock()
	go func() {
		// existing data for the subscriber
		s.mu.Lock()
		if data, ok := s.data[id]; ok {
			for _, d := range data {
				select {
				case ch <- d:
				default:
					// If the channel is full, we skip sending to avoid blocking
				}
			}
		}
		s.mu.Unlock()
	}()
	go func() {
		<-ctx.Done()
		s.mu.Lock()
		defer s.mu.Unlock()
		s.subscribers = slices.DeleteFunc(s.subscribers, func(sub subscriber) bool {
			return sub.id == id
		})
	}()
	return ch, nil
}

// AskQuestion sends a question to the LLM and returns the response id.
func (s *Service) AskQuestion(question string, instruction string, mcps ...MCP) uuid.UUID {
	tools := make([]responses.ToolUnionParam, 0, len(mcps))
	for _, mcp := range mcps {
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
		Instructions: param.NewOpt(instruction),
		Tools:        tools,
		Model:        "gpt-4o",
	})
	id := uuid.New()

	s.mu.Lock()
	s.data[id] = []StreamData{} // Initialize the data for this stream
	s.mu.Unlock()

	go func() {
		s.handleStream(id, stream)
	}()

	return id
}

func (s *Service) handleStream(id uuid.UUID, stream *ssestream.Stream[responses.ResponseStreamEventUnion]) {
	for stream.Next() {
		res := stream.Current()
		if text, ok := res.AsAny().(responses.ResponseTextDeltaEvent); ok {
			s.publishData(StreamData{
				ID:        id,
				TextDelta: text.Delta,
			})
		}
	}
	if err := stream.Err(); err != nil {
		slog.Error("Stream error", slog.String("id", id.String()), slog.Any("error", err))
		s.publishData(StreamData{
			ID:  id,
			Err: err,
		})
	}

	s.closeStream(id)
}

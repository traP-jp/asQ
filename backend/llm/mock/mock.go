package mock

import (
	"iter"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/traP-jp/h25s_05/backend/llm"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

type Stream struct {
	id string
}

func (s *Service) AskQuestion(question string, instruction string, previousResponseID string) llm.Stream {
	return &Stream{
		id: uuid.NewString(),
	}
}

func (s *Stream) ID() string {
	return s.id
}

func (s *Stream) Next() iter.Seq[llm.Response] {
	return func(yield func(llm.Response) bool) {
		res := strings.Repeat("This is a mock response for the question: "+s.id+"\n", 10)
		time.Sleep(100 * time.Millisecond)
		for _, c := range res {
			time.Sleep(10 * time.Millisecond)
			if !yield(llm.Response{Text: string(c)}) {
				return
			}
		}
	}
}

func (s *Stream) Close() error {
	return nil
}

func (s *Stream) Err() error {
	return nil
}

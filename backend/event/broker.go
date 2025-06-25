package event

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Message[T any] struct {
	value     T
	timestamp time.Time
}

type Partition[T any] struct {
	messages      []Message[T]
	mu            sync.RWMutex
	cond          *sync.Cond
	cleanupTicker *time.Ticker
}

const (
	WhenceStart = 0
	WhenceEnd   = 2
)

func NewMessage[T any](value T) Message[T] {
	return Message[T]{value: value, timestamp: time.Now()}
}

func (m Message[T]) Value() T {
	return m.value
}

func NewPartition[T any](retention time.Duration, cleanupInterval time.Duration) *Partition[T] {
	p := &Partition[T]{messages: make([]Message[T], 0)}
	p.cond = sync.NewCond(&p.mu)
	if retention > 0 {
		p.cleanupTicker = time.NewTicker(cleanupInterval)
		go p.startCleanupRoutine(retention)
	}
	return p
}

func (p *Partition[T]) startCleanupRoutine(retention time.Duration) {
	for range p.cleanupTicker.C {
		p.cleanupExpiredMessages(retention)
	}
}

func (p *Partition[T]) cleanupExpiredMessages(retention time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.messages) == 0 {
		return
	}
	cutoffTime := time.Now().Add(-retention)
	firstValidIndex := 0
	for i, msg := range p.messages {
		if !msg.timestamp.Before(cutoffTime) {
			firstValidIndex = i
			break
		}
		if i == len(p.messages)-1 {
			firstValidIndex = len(p.messages)
		}
	}
	if firstValidIndex > 0 {
		p.messages = p.messages[firstValidIndex:]
	}
}

func (p *Partition[T]) StopCleanup() {
	if p.cleanupTicker != nil {
		p.cleanupTicker.Stop()
	}
}

func (p *Partition[T]) Append(msg Message[T]) int64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.messages = append(p.messages, msg)
	offset := int64(len(p.messages) - 1)
	p.cond.Broadcast()
	return offset
}

func (p *Partition[T]) getMessagesInternal(offset int64, whence int) ([]Message[T], int64) {
	length := int64(len(p.messages))
	nextOffset := length
	var startOffset int64
	switch whence {
	case WhenceStart:
		startOffset = offset
	case WhenceEnd:
		startOffset = length + offset
	default:
		return nil, nextOffset
	}
	if startOffset < 0 || startOffset > nextOffset {
		return nil, nextOffset
	}
	return p.messages[startOffset:], nextOffset
}

type Topic[T any] struct {
	name         string
	partition    *Partition[T]
	lastAccessed time.Time
	mu           sync.Mutex
}

func NewTopic[T any](name string, retention time.Duration, cleanupInterval time.Duration) *Topic[T] {
	return &Topic[T]{
		name:         name,
		partition:    NewPartition[T](retention, cleanupInterval),
		lastAccessed: time.Now(),
	}
}

func (t *Topic[T]) touch() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lastAccessed = time.Now()
}
func (t *Topic[T]) Stop() {
	t.partition.StopCleanup()
}

type config struct {
	messageRetention       time.Duration
	messageCleanupInterval time.Duration
	topicInactivityTime    time.Duration
	topicCleanupInterval   time.Duration
}

type BrokerOption[T any] func(*config)

func defaultConfig() config {
	return config{
		messageRetention:       10 * time.Minute,
		messageCleanupInterval: 1 * time.Minute,
		topicInactivityTime:    30 * time.Minute,
		topicCleanupInterval:   1 * time.Minute,
	}
}

func WithMessageRetention[T any](messageRetention time.Duration, messageCleanupInterval time.Duration) BrokerOption[T] {
	return func(c *config) {
		c.messageRetention = messageRetention
		c.messageCleanupInterval = messageCleanupInterval
	}
}

func WithTopicInactivity[T any](inactivityTime, cleanupInterval time.Duration) BrokerOption[T] {
	return func(c *config) {
		c.topicInactivityTime = inactivityTime
		c.topicCleanupInterval = cleanupInterval
	}
}

type Broker[T any] struct {
	topics             map[string]*Topic[T]
	mu                 sync.RWMutex
	config             config
	topicCleanupTicker *time.Ticker
}

func NewBroker[T any](opts ...BrokerOption[T]) *Broker[T] {
	conf := defaultConfig()

	for _, opt := range opts {
		opt(&conf)
	}

	b := &Broker[T]{
		topics: make(map[string]*Topic[T]),
		config: conf,
	}

	if b.config.topicInactivityTime > 0 {
		b.topicCleanupTicker = time.NewTicker(b.config.topicCleanupInterval)
		go b.startTopicCleanupRoutine()
	}

	return b
}

func (b *Broker[T]) startTopicCleanupRoutine() {
	for range b.topicCleanupTicker.C {
		b.cleanupInactiveTopics()
	}
}

func (b *Broker[T]) cleanupInactiveTopics() {
	b.mu.Lock()
	defer b.mu.Unlock()
	cutoffTime := time.Now().Add(-b.config.topicInactivityTime)
	for name, topic := range b.topics {
		topic.mu.Lock()
		isInactive := topic.lastAccessed.Before(cutoffTime)
		topic.mu.Unlock()
		if isInactive {
			topic.Stop()
			delete(b.topics, name)
		}
	}
}

// CreateTopic creates a new topic if it does not already exist.
func (b *Broker[T]) CreateTopic(topicName string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.topics[topicName]; !ok {
		b.topics[topicName] = NewTopic[T](topicName, b.config.messageRetention, b.config.messageCleanupInterval)
	}
}

// Fetch retrieves messages from the specified topic starting from the given offset.
// The whence parameter determines the starting point for the offset:
//   - [WhenceStart]: Start from the given offset
//   - [WhenceEnd]: Start from the end of the topic
//
// If no messages are available, it waits until new messages are produced or the context is canceled.
// Returns the messages and the next offset to use for subsequent fetches.
func (b *Broker[T]) Fetch(ctx context.Context, topicName string, offset int64, whence int) ([]Message[T], int64, error) {
	partition, err := b.getPartition(topicName)
	if err != nil {
		return nil, -1, err
	}

	go func() {
		<-ctx.Done()
		partition.cond.Broadcast()
	}()

	partition.mu.Lock()
	defer partition.mu.Unlock()
	for {
		currentMessages, nextOffset := partition.getMessagesInternal(offset, whence)
		if len(currentMessages) > 0 {
			return currentMessages, nextOffset, nil
		}
		if err := ctx.Err(); err != nil {
			return nil, nextOffset, fmt.Errorf("context error: %w", err)
		}
		partition.cond.Wait()
	}
}

func (b *Broker[T]) getPartition(topicName string) (*Partition[T], error) {
	b.mu.RLock()
	topic, ok := b.topics[topicName]
	b.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("topic not found: %s", topicName)
	}
	topic.touch()
	return topic.partition, nil
}

// Produce appends a message to the specified topic and returns the offset.
func (b *Broker[T]) Produce(topicName string, msg Message[T]) (int64, error) {
	partition, err := b.getPartition(topicName)
	if err != nil {
		return -1, err
	}
	return partition.Append(msg), nil
}

package event_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/traP-jp/h25s_05/backend/event"
)

// Test 1: 基本的なメッセージの送受信テスト
func TestBroker_ProduceAndFetch_Simple(t *testing.T) {
	r := require.New(t)

	// --- セットアップ ---
	broker := event.NewBroker[string]()
	topic := "test-simple"
	broker.CreateTopic(topic)

	// --- 実行 ---
	msg := event.NewMessage("hello world")
	offset, err := broker.Produce(topic, msg)
	r.NoError(err)
	r.Equal(int64(0), offset) // 最初のメッセージのオフセットは0

	// --- 検証 ---
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	fetchedMessages, nextOffset, err := broker.Fetch(ctx, topic, 0, event.WhenceStart)
	r.NoError(err)
	r.Equal(int64(1), nextOffset)
	r.Len(fetchedMessages, 1) // 1件のメッセージが取得できるはず
	r.Equal("hello world", fetchedMessages[0].Value())
}

// Test 2: Fetchがcontextのタイムアウトに従うかのテスト
func TestBroker_Fetch_Timeout(t *testing.T) {
	r := require.New(t)

	// --- セットアップ ---
	broker := event.NewBroker[string]()
	topic := "test-timeout"
	broker.CreateTopic(topic)

	// --- 実行 & 検証 ---
	// 非常に短いタイムアウトを持つcontextを作成
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// まだメッセージはないので、Fetchはタイムアウトするはず
	_, _, err := broker.Fetch(ctx, topic, 0, event.WhenceStart)
	r.Error(err) // エラーが発生することを確認
	// エラーがcontext.DeadlineExceededであることを確認
	r.ErrorIs(err, context.DeadlineExceeded, "error should be context.DeadlineExceeded")
}

// Test 3: メッセージが来るまでFetchが待機する（ロングポーリング）テスト
func TestBroker_Fetch_BlocksUntilProduce(t *testing.T) {
	r := require.New(t)

	// --- セットアップ ---
	broker := event.NewBroker[string]()
	topic := "test-blocking"
	broker.CreateTopic(topic)

	resultCh := make(chan event.Message[string], 1)

	// --- 実行 ---
	// goroutineでFetchを呼び出す
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		messages, _, err := broker.Fetch(ctx, topic, 0, event.WhenceStart)
		if err == nil && len(messages) > 0 {
			resultCh <- messages[0]
		}
	}()

	// Fetchが待機状態に入るのを少し待つ
	time.Sleep(100 * time.Millisecond)

	// メッセージを送信して、待機中のFetchを再開させる
	msg := event.NewMessage("unblock me")
	_, err := broker.Produce(topic, msg)
	r.NoError(err)

	// --- 検証 ---
	select {
	case receivedMsg := <-resultCh:
		r.Equal("unblock me", receivedMsg.Value())
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for fetch to complete")
	}
}

// Test 4: メッセージが保持期間を過ぎると削除されるかのテスト
func TestBroker_MessageRetention(t *testing.T) {
	r := require.New(t)

	// --- セットアップ ---
	// 非常に短い保持期間でブローカーを作成
	broker := event.NewBroker(
		event.WithMessageRetention[string](50*time.Millisecond, 20*time.Millisecond),
	)
	topic := "test-retention"
	broker.CreateTopic(topic)

	// --- 実行 ---
	// 1. メッセージを送信し、即座に取得できることを確認
	_, err := broker.Produce(topic, event.NewMessage("i will disappear"))
	r.NoError(err)

	ctx1, cancel1 := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel1()
	messages1, _, err := broker.Fetch(ctx1, topic, 0, event.WhenceStart)
	r.NoError(err)
	r.Len(messages1, 1)

	// 2. 保持期間より長く待機
	time.Sleep(100 * time.Millisecond)

	// --- 検証 ---
	// 3. 再度Fetchを試みるが、メッセージは削除されているはず
	ctx2, cancel2 := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel2()
	// Fetchはメッセージがないためタイムアウトする
	_, _, err = broker.Fetch(ctx2, topic, 0, event.WhenceStart)
	r.Error(err)
	r.ErrorIs(err, context.DeadlineExceeded)
}

// Test 5: 非アクティブなトピックが削除されるかのテスト
func TestBroker_TopicRetention(t *testing.T) {
	r := require.New(t)

	// --- セットアップ ---
	// 非常に短いトピック非アクティブ期間でブローカーを作成
	broker := event.NewBroker(
		event.WithTopicInactivity[string](100*time.Millisecond, 20*time.Millisecond),
	)
	activeTopic := "active-topic"
	inactiveTopic := "inactive-topic"
	broker.CreateTopic(activeTopic)
	broker.CreateTopic(inactiveTopic)

	// --- 実行 ---
	// inactive-topicには一度だけアクセス
	_, err := broker.Produce(inactiveTopic, event.NewMessage("initial"))
	r.NoError(err)

	// active-topicには定期的にアクセスし続ける
	stopCh := make(chan struct{})
	go func() {
		ticker := time.NewTicker(40 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				broker.Produce(activeTopic, event.NewMessage("keepalive"))
			case <-stopCh:
				return
			}
		}
	}()

	// 2. 非アクティブ期間より長く待機
	time.Sleep(200 * time.Millisecond)
	close(stopCh) // keepalive goroutineを停止

	// --- 検証 ---
	// 3. active-topicは存在し、inactive-topicは削除されていることを確認
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, _, err = broker.Fetch(ctx, activeTopic, 0, event.WhenceStart)
	r.NoError(err) // active-topicはエラーにならない

	_, _, err = broker.Fetch(ctx, inactiveTopic, 0, event.WhenceStart)
	r.Error(err) // inactive-topicは "topic not found" エラーになる
	r.Contains(err.Error(), "topic not found")
}

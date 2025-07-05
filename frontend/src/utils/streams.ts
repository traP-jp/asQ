import { type Ref } from 'vue'
import { type TypewriterMessage, type TypewriterController } from '@/utils/typewriter'
import api from '@/utils/api'

// メッセージ型定義
export interface StreamMessage {
  id: string
  message: string
  time: Date
  isAi?: boolean
}

// ストリームイベントハンドラー型
export interface StreamEventHandlers {
  onUserMessage?: (message: StreamMessage) => void
  onAIMessage?: (message: StreamMessage) => void
  onStreamStart?: (characterId: string) => void
  onStreamData?: (data: string) => void
  onStreamEnd?: () => void
  onError?: (error: Error) => void
}

// ストリーム設定
export interface StreamConfig {
  chatId: string
  typewriter: TypewriterController
  chatMessages: Ref<StreamMessage[]>
  aiMessages: Ref<TypewriterMessage[]>
}

export class StreamManager {
  private aiResponseEventSource: EventSource | null = null
  private speakingStatusEventSource: EventSource | null = null
  private currentIndex = -1
  private config: StreamConfig
  private handlers: StreamEventHandlers

  constructor(config: StreamConfig, handlers: StreamEventHandlers = {}) {
    this.config = config
    this.handlers = handlers
  }

  // AIレスポンス用SSEストリームを開始
  openAIResponseStream(responseId: string): void {
    // 既存のストリームをクローズ
    this.aiResponseEventSource?.close()

    this.aiResponseEventSource = new EventSource(`/api/sse/ai/${responseId}`)
    let chatMessageIndex = -1

    // ストリーム開始イベント
    this.aiResponseEventSource.addEventListener('start', (e: MessageEvent) => {
      try {
        const { characterId } = JSON.parse(e.data)

        // aiMessagesに追加
        this.config.aiMessages.value.push({
          id: String(characterId),
          message: '',
          displayedMessage: '',
          currentIndex: 0,
        })
        this.currentIndex = this.config.aiMessages.value.length - 1

        // chatMessagesにもプレースホルダを追加
        chatMessageIndex = this.config.chatMessages.value.length
        const newMessage: StreamMessage = {
          id: characterId,
          message: '',
          time: new Date(),
          isAi: true,
        }
        this.config.chatMessages.value.push(newMessage)

        this.handlers.onStreamStart?.(characterId)
      } catch (err) {
        const error = new Error(`Failed to parse AI response start event: ${err}`)
        console.error(error.message)
        this.handlers.onError?.(error)
      }
    })

    // データ断片受信イベント
    this.aiResponseEventSource.addEventListener('data', (e: MessageEvent) => {
      if (this.currentIndex === -1) return

      try {
        const { message } = JSON.parse(e.data)
        this.config.aiMessages.value[this.currentIndex].message += message
        this.config.typewriter.updateMessage(this.config.aiMessages.value, this.currentIndex)

        // chatMessages側も更新
        if (chatMessageIndex !== -1) {
          this.config.chatMessages.value[chatMessageIndex].message += message
        }

        this.handlers.onStreamData?.(message)
      } catch (err) {
        const error = new Error(`Failed to parse AI response data event: ${err}`)
        console.error(error.message)
        this.handlers.onError?.(error)
      }
    })

    // 完了イベント
    this.aiResponseEventSource.addEventListener('close', () => {
      this.aiResponseEventSource?.close()
      this.currentIndex = -1
      this.handlers.onStreamEnd?.()
    })
  }

  // 発言状態監視SSEを開始
  async startSpeakingStatusStream(): Promise<void> {
    if (!this.config.chatId) {
      throw new Error('Chat ID is required to start speaking status stream')
    }

    this.speakingStatusEventSource = new EventSource(`/api/sse/events/${this.config.chatId}`)

    // ユーザーイベント
    this.speakingStatusEventSource.addEventListener('user', async (e: MessageEvent) => {
      try {
        console.log('Raw user event data:', e.data)
        const id = (e.data as string).trim()

        const { data } = await api.get(`/api/messages/${id}`)

        const userMessage: StreamMessage = {
          id: data.userId,
          message: data.message,
          time: new Date(data.createdAt),
          isAi: false,
        }

        this.config.chatMessages.value.push(userMessage)
        this.handlers.onUserMessage?.(userMessage)
      } catch (err) {
        const error = new Error(`Failed to parse user speaking event: ${err}`)
        console.error('Failed to parse user speaking event. Raw data:', e.data, 'Error:', err)
        this.handlers.onError?.(error)
      }
    })

    // AIイベント
    this.speakingStatusEventSource.addEventListener('ai', (e: MessageEvent) => {
      try {
        console.log('Raw AI event data:', e.data)
        const id = (e.data as string).trim()
        console.log('AI speaking event received. Response ID:', id)

        // AIレスポンス用のSSEストリームを開く
        this.openAIResponseStream(id)
      } catch (err) {
        const error = new Error(`Failed to parse AI speaking event: ${err}`)
        console.error('Failed to parse ai speaking event. Raw data:', e.data, 'Error:', err)
        this.handlers.onError?.(error)
      }
    })
  }

  // 既存のチャットログを取得
  async loadChatHistory(): Promise<StreamMessage[]> {
    try {
      const { data } = await api.get(`/api/chats/${this.config.chatId}/log`)
      const history: StreamMessage[] = []

      // ユーザーメッセージを追加
      if (data.messages) {
        for (const m of data.messages) {
          history.push({
            id: m.userId,
            message: m.message,
            time: new Date(m.createdAt),
            isAi: false,
          })
        }
      }

      // AIレスポンスを追加
      if (data.responses) {
        for (const r of data.responses) {
          history.push({
            id: r.characterId,
            message: r.message,
            time: new Date(r.createdAt),
            isAi: true,
          })
        }
      }

      // 時系列順に並べ替え
      history.sort((a, b) => a.time.getTime() - b.time.getTime())

      // chatMessagesに設定
      this.config.chatMessages.value = history

      return history
    } catch (err) {
      const error = new Error(`Failed to fetch chat log: ${err}`)
      console.error('Failed to fetch chat log', err)
      this.handlers.onError?.(error)
      throw error
    }
  }

  // メッセージ送信
  async sendMessage(message: string, characterId: string): Promise<void> {
    try {
      await api.post(`/api/chats/${this.config.chatId}/search`, {
        message,
        characterId,
      })
      console.log('message sent')
    } catch (err) {
      const error = new Error(`Failed to send message: ${err}`)
      console.error('Failed to send message:', err)
      this.handlers.onError?.(error)
      throw error
    }
  }

  // クリーンアップ
  cleanup(): void {
    this.aiResponseEventSource?.close()
    this.speakingStatusEventSource?.close()
    this.config.typewriter.cleanup()
  }
}

// ファクトリー関数
export function createStreamManager(
  config: StreamConfig,
  handlers?: StreamEventHandlers,
): StreamManager {
  return new StreamManager(config, handlers)
}

<template>
  <div class="container">
    <Header title="Chat Room"></Header>
    <div class="main-content">
      <div class="sidebar">
        <ChooseCharacters v-model="selectedCharacterId" />
      </div>
      <div class="chat-space">
        <div class="chat-container">
          <div v-for="(message, index) in chatMessages" :key="index">
            <div class="message">
              <UserMessage
                class="user"
                v-if="!message.isAi"
                :id="message.id"
                :message="message.message"
              />
              <AiMessage class="ai" v-else :imageUrl="message.id" :message="message.message" />
            </div>
          </div>
        </div>
        <div style="display: flex; justify-content: center; width: 100%; height: 15%; margin-top: 10px">
          <InputText class="input-text" @sendMessage="handleSendMessage" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import Header from '@/components/HeaderComponent.vue'
import AiMessage from '@/components/AiMessage.vue'
import api from '@/utils/api'
import { useRoute } from 'vue-router'
import { createTypewriter, type TypewriterMessage } from '@/utils/typewriter'
import ChooseCharacters from '@/components/ChooseCharacters.vue'

const aiMessages = ref<TypewriterMessage[]>([])
const selectedCharacterId = ref<string>('d7ec5756-4f20-11f0-ae2b-6ae4cc62b771')
let currentIndex = -1 // 現在構築中のメッセージのインデックス
let aiResponseEventSource: EventSource | null = null // AIの発言内容を受け取るSSE
let speakingStatusEventSource: EventSource | null = null // 誰が発言したかの状態を受け取るSSE

// URLパラメータからchatIdを取得（setup関数の直接スコープで実行）
const route = useRoute()
const chatId = route.params.id as string

// タイプライターコントローラーを初期化
const typewriter = createTypewriter(30) // 30ms間隔でより滑らかに

// AIレスポンス用SSEストリームを開き、メッセージを受信して更新するヘルパー
const openAIResponseStream = (responseId: string) => {
  // 既存のストリームをクローズ
  aiResponseEventSource?.close()

  aiResponseEventSource = new EventSource(`/api/sse/ai/${responseId}`)

  // chatMessages に追加した AI メッセージのインデックスを保持
  let chatMessageIndex = -1

  // ストリーム開始
  aiResponseEventSource.addEventListener('start', (e: MessageEvent) => {
    try {
      const { characterId } = JSON.parse(e.data)
      aiMessages.value.push({
        id: String(characterId),
        message: '',
        displayedMessage: '',
        currentIndex: 0,
      })
      currentIndex = aiMessages.value.length - 1

      // chatMessages にもプレースホルダを追加
      chatMessageIndex = chatMessages.value.length
      chatMessages.value.push({
        id: characterId,
        message: '',
        time: new Date(),
        isAi: true,
      })
    } catch (err) {
      console.error('Failed to parse AI response start event', err)
    }
  })

  // データ断片受信
  aiResponseEventSource.addEventListener('data', (e: MessageEvent) => {
    if (currentIndex === -1) return
    try {
      const { message } = JSON.parse(e.data)
      aiMessages.value[currentIndex].message += message
      typewriter.updateMessage(aiMessages.value, currentIndex)

      // chatMessages 側も更新
      if (chatMessageIndex !== -1) {
        chatMessages.value[chatMessageIndex].message += message
      }
    } catch (err) {
      console.error('Failed to parse AI response data event', err)
    }
  })

  // 完了
  aiResponseEventSource.addEventListener('close', () => {
    aiResponseEventSource?.close()
    currentIndex = -1

    // ソートせず順序を保持
  })
}

// メッセージ送信ハンドラー
const handleSendMessage = async (message: string) => {
  try {
    // メッセージをサーバーに送信
    await api.post(`/api/chats/${chatId}/search`, {
      message,
      characterId: selectedCharacterId.value,
    })

    console.log('message sent')
  } catch (err) {
    console.error('Failed to send message:', err)
  }
}

onMounted(async () => {
  // chatIdの存在確認
  if (!chatId) {
    console.error('Chat ID is not provided in route parameters')
    return
  }

  // ユーザーIDを取得
  const { data } = await api.get('/api/users/me')
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const userId = data.userId

  // 発言状態を監視するSSEエンドポイント（誰が発言したかを受け取る）
  speakingStatusEventSource = new EventSource(`/api/sse/events/${chatId}`)

  // エラーイベント
  speakingStatusEventSource.addEventListener('error', (e: Event) => {
    console.error('Error occurred in speakingStatusEventSource:', e)
    // Optionally, implement retry logic or notify the user here
  })
  // ユーザーイベント
  speakingStatusEventSource.addEventListener('user', async (e: MessageEvent) => {
    try {
      console.log('Raw user event data:', e.data)
      const id = (e.data as string).trim()

      const { data } = await api.get(`/api/messages/${id}`)

      chatMessages.value.push({
        id: data.userId, // アイコン表示のため userId を使用
        message: data.message,
        time: new Date(data.createdAt),
        isAi: false, // ユーザーメッセージなので false
      })
      // ソートせず順序を保持
    } catch (err) {
      console.error('Failed to parse user speaking event. Raw data:', e.data, 'Error:', err)
    }
  })

  // AIイベント
  speakingStatusEventSource.addEventListener('ai', (e: MessageEvent) => {
    try {
      console.log('Raw AI event data:', e.data)
      const id = (e.data as string).trim()
      console.log('AI speaking event received. Response ID:', id)

      // AIレスポンス用のSSEストリームを開く
      openAIResponseStream(id)
    } catch (err) {
      console.error('Failed to parse ai speaking event. Raw data:', e.data, 'Error:', err)
    }
  })

  // 既存のログを取得
  try {
    const { data } = await api.get(`/api/chats/${chatId}/log`)
    // data.messages (user), data.responses (ai)
    const history: Messages[] = []
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
    // 時系列順に並べて格納
    history.sort((a, b) => a.time.getTime() - b.time.getTime())
    chatMessages.value = history
  } catch (err) {
    console.error('Failed to fetch chat log', err)
  }
})

onBeforeUnmount(() => {
  aiResponseEventSource?.close()
  speakingStatusEventSource?.close()
  typewriter.cleanup()
})

import UserMessage from '@/components/UserMessage.vue'
import InputText from '@/components/InputText.vue'
type Messages = {
  id: string
  message: string
  time: Date
  isAi?: boolean
}

type ChatMessages = {
  userMessages: Messages[]
  aiMessages: Messages[]
}

// チャットメッセージ一覧
const chatMessages = ref<Messages[]>([])
</script>

<style scoped>
.container {
  height: 100vh;
  width: 100%;
  max-width: 100%;
  background: linear-gradient(135deg, #aad5f9 0%, #f5dcfe 100%);
  overflow: hidden;
  display: flex;
  flex-direction: column;
  /* height: calc(100vh - 64px); Header 分を差し引く。Header の高さに合わせて調整 */
  padding: 16px;
  overflow-y: auto;
}
.chat-space {
  justify-content: flex-end;
  width: 60%;
  margin: auto;
}

.main-content {
  display: flex;
  flex: 1;
  height: calc(100vh - 64px); /* Header分を差し引く */
  overflow: hidden;
}

.sidebar {
  flex-shrink: 0;
  width: 400px;
  padding: 20px;
  display: flex;
  align-items: flex-start;
  justify-content: center;
}

.chat-space {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 20px;
  height: 100%;
  max-width: calc(100% - 200px);
  min-width: 0; /* flexboxでの縮小を許可 */
}

.chat-container {
  background-color: #f3f6fb;
  overflow-y: auto;
  flex: 1;
  width: 100%;
  border-radius: 15px;
  padding: 20px;
  margin-bottom: 20px;
  box-shadow: 0px 10px 10px -6px rgba(0, 0, 0, 0.3);
  min-height: 0; /* flexboxでの縮小を許可 */
}

.message {
  display: flex;
  margin-bottom: 10px;
}

.user {
  margin-top: 5px;
  margin-bottom: 15px;
  flex: 1;
}

.ai {
  margin-bottom: 3px;
  flex: 2;
}

.input-text {
  flex-shrink: 0;
  z-index: 1000;
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>

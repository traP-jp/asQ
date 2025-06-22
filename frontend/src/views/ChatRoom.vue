
<template>
  <Header title="Chat Room"></Header>
  <div class="chat-room">
    <div class="chat-space">
      <div class="chat-container">
        <div v-for="(message, index) in chatMessages" :key="index">
          <div class="messag">
            <UserMessage v-if="!message.isAi" :id="message.id" :message="message.message" />
            <AiMessage v-else :id="message.id" :message="message.message" />
          </div>
        </div>
      </div>
      <InputText
        class="input-text"
        @sendMessage="
          (message: string) => {
            const newMessage: Messages = {
              id: `user${chatMessagesResponse.userMessages.length + 1}`,
              message,
              time: new Date(),
              isAi: false,
            }
            chatMessagesResponse.userMessages.push(newMessage)
            chatMessages = [...chatMessagesResponse.userMessages, ...chatMessagesResponse.aiMessages]
            chatMessages.sort((a, b) => a.time.getTime() - b.time.getTime())
          }
        "
      />
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

interface SpeakingStatus {
  id: string
  type: 'user' | 'ai'
}

const aiMessages = ref<TypewriterMessage[]>([])
const speakingStatus = ref<SpeakingStatus>({
  id: '',
  type: 'user',
})
let currentIndex = -1 // 現在構築中のメッセージのインデックス
let aiResponseEventSource: EventSource | null = null // AIの発言内容を受け取るSSE
let speakingStatusEventSource: EventSource | null = null // 誰が発言したかの状態を受け取るSSE

// URLパラメータからchatIdを取得（setup関数の直接スコープで実行）
const route = useRoute()
const chatId = route.params.id as string

// タイプライターコントローラーを初期化
const typewriter = createTypewriter(30) // 30ms間隔でより滑らかに

onMounted(async () => {
  // chatIdの存在確認
  if (!chatId) {
    console.error('Chat ID is not provided in route parameters')
    return
  }

  // ユーザーIDを取得
  // const { data } = await api.get('/api/users/me')
  // const userId = data.userId

  // 発言状態を監視するSSEエンドポイント（誰が発言したかを受け取る）
  speakingStatusEventSource = new EventSource(`/api/sse/events/${chatId}`)

  // エラーイベント
  speakingStatusEventSource.addEventListener('error', (e: Event) => {
    console.error('Error occurred in speakingStatusEventSource:', e)
    // Optionally, implement retry logic or notify the user here
  })
  // ユーザーイベント
  speakingStatusEventSource.addEventListener('user', (e: MessageEvent) => {
    try {
      const { id } = JSON.parse(e.data)
      console.log('User speaking event received:', id)

      // ユーザーの発言状態を追加
      speakingStatus.value = { id, type: 'user' }
    } catch (err) {
      console.error('Failed to parse user speaking event', err)
    }
  })

  // AIイベント
  speakingStatusEventSource.addEventListener('ai', (e: MessageEvent) => {
    try {
      const { id } = JSON.parse(e.data)
      console.log('AI speaking event received:', id)

      // AIの発言状態を追加
      speakingStatus.value = { id, type: 'ai' }
    } catch (err) {
      console.error('Failed to parse ai speaking event', err)
    }
  })

  // メッセージを送信  開発用！！！！！！！！！
  //　後で消す！！！！！！！！！！！！！！！！！！！！
  let responseId: string
  try {
    const { data } = await api.post(`/api/chats/${chatId}/search`, {
      message: 'ハッカソンについて詳しく教えて100字以上で',
      characterId: '12345',
    })
    responseId = data.id
  } catch (err) {
    console.error('Failed to send message:', err)
    return
  }

  // AIの発言内容を受け取るSSEを開く（レスポンスIDを使用）
  aiResponseEventSource = new EventSource(`/api/sse/ai/${responseId}`)

  // 会話開始：characterId が渡されるので新しいメッセージを作成
  aiResponseEventSource.addEventListener('start', (e: MessageEvent) => {
    try {
      const { characterId } = JSON.parse(e.data)
      aiMessages.value.push({
        id: characterId,
        message: '',
        displayedMessage: '',
        currentIndex: 0,
      })
      currentIndex = aiMessages.value.length - 1
    } catch (err) {
      console.error('Failed to parse AI response start event', err)
    }
  })

  // AIの発言内容の断片を受信
  aiResponseEventSource.addEventListener('data', (e: MessageEvent) => {
    if (currentIndex === -1) return // start がまだ来ていない
    try {
      const { message } = JSON.parse(e.data)
      aiMessages.value[currentIndex].message += message
      typewriter.updateMessage(aiMessages.value, currentIndex)
    } catch (err) {
      console.error('Failed to parse AI response data event', err)
    }
  })

  // AI発言完了
  aiResponseEventSource.addEventListener('close', () => {
    aiResponseEventSource?.close()
    currentIndex = -1
  })
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

const chatMessagesResponse = ref<ChatMessages>({
  userMessages: [
    {
      id: 'user1',
      message: 'Hello, how are you?',
      time: new Date(2022, 0, 1, 0, 0, 0),
      isAi: false,
    },
    {
      id: 'user2',
      message: 'What is the weather like today?',
      time: new Date(2022, 0, 1, 0, 10, 0),
      isAi: false,
    },
  ],
  aiMessages: [
    {
      id: 'ai1',
      message: 'I am fine, thank you!',
      time: new Date(2022, 0, 1, 0, 0, 5),
      isAi: true,
    },
    {
      id: 'ai2',
      message: 'The weather is sunny today.',
      time: new Date(2022, 0, 1, 0, 10, 5),
      isAi: true,
    },
  ],
})

const chatMessages = ref<Messages[]>([])

chatMessages.value = [
  ...chatMessagesResponse.value.userMessages,
  ...chatMessagesResponse.value.aiMessages,
]

chatMessages.value.sort((a, b) => a.time.getTime() - b.time.getTime())
</script>

<style scoped>
.chat-room-container {
  width: 100%;
  height: calc(100vh - 64px); /* Header 分を差し引く。Header の高さに合わせて調整 */
  padding: 16px;
  overflow-y: auto;
}

.message-wrapper {
  margin-bottom: 12px;
}

.chatspace{
  flex-direction: column;
  width: 80%;
}
.chat-container{
  background-color: aliceblue;
  overflow-y: scroll;
}

</style>




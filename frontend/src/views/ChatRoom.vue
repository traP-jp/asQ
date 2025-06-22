<template>
  <!-- <div class="w-100 h-100">
    <h1>Chat Room</h1>
    <v-btn> aaa </v-btn>
  </div> -->
  <Header title="Chat Room" />
  <div class="chat-room-container">
    <ChooseCharacters v-model="selectedCharacterId" />
    受信したAIメッセージを順次表示
    <div v-for="(msg, idx) in aiMessages" :key="idx" class="message-wrapper">
      <AiMessage :id="msg.id" :message="msg.displayedMessage" />
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
interface SpeakingStatus {
  id: string
  type: 'user' | 'ai'
}

const aiMessages = ref<TypewriterMessage[]>([])
const speakingStatus = ref<SpeakingStatus>({
  id: '',
  type: 'user',
})
const selectedCharacterId = ref<string>('ai1')
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

  // // 発言状態を監視するSSEエンドポイント（誰が発言したかを受け取る）
  // speakingStatusEventSource = new EventSource(`/api/sse/events/${chatId}`)

  // // エラーイベント
  // speakingStatusEventSource.addEventListener('error', (e: Event) => {
  //   console.error('Error occurred in speakingStatusEventSource:', e)
  //   // Optionally, implement retry logic or notify the user here
  // })
  // // ユーザーイベント
  // speakingStatusEventSource.addEventListener('user', (e: MessageEvent) => {
  //   // バックエンドからは id のみが渡されるため、そのまま使用する
  //   const id = e.data.trim()
  //   console.log('User speaking event received:', id)

  //   // ユーザーの発言状態を追加
  //   speakingStatus.value = { id, type: 'user' }
  // })

  // // AIイベント
  // speakingStatusEventSource.addEventListener('ai', (e: MessageEvent) => {
  //   const id = e.data.trim()
  //   console.log('AI speaking event received:', id)

  //   // AIの発言状態を追加
  //   speakingStatus.value = { id, type: 'ai' }
  // })

  // メッセージを送信  開発用！！！！！！！！！
  //　後で消す！！！！！！！！！！！！！！！！！！！！
  let responseId: string
  try {
    const { data } = await api.post(`/api/chats/${chatId}/search`, {
      message: 'ハッカソンについて詳しく教えて100字以上で',
      characterId: 'c72d594f-4f1b-11f0-b33d-629483e90542	',
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
        id: String(characterId),
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
  // speakingStatusEventSource?.close()
  typewriter.cleanup()
})
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
</style>

<template>
  <!-- <div class="w-100 h-100">
    <h1>Chat Room</h1>
    <v-btn> aaa </v-btn>
  </div> -->
  <Header title="Chat Room" />
  <div class="chat-room-container">
    <!-- 受信したAIメッセージを順次表示 -->
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

const aiMessages = ref<TypewriterMessage[]>([])
let currentIndex = -1 // 現在構築中のメッセージのインデックス
let eventSource: EventSource | null = null

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
  const { data } = await api.get('/api/users/me')
  const userId = data.userId

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

  // 3. SSE を開く（レスポンスIDを使用）
  eventSource = new EventSource(`/api/sse/ai/${responseId}`)

  // 会話開始：characterId が渡されるので新しいメッセージを作成
  eventSource.addEventListener('start', (e: MessageEvent) => {
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
      console.error('Failed to parse start event', err)
    }
  })

  // メッセージの断片を受信
  eventSource.addEventListener('data', (e: MessageEvent) => {
    if (currentIndex === -1) return // start がまだ来ていない
    try {
      const { message } = JSON.parse(e.data)
      aiMessages.value[currentIndex].message += message
      typewriter.updateMessage(aiMessages.value, currentIndex)
    } catch (err) {
      console.error('Failed to parse data event', err)
    }
  })

  // 会話終了
  eventSource.addEventListener('close', () => {
    eventSource?.close()
    currentIndex = -1
  })
})

onBeforeUnmount(() => {
  eventSource?.close()
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

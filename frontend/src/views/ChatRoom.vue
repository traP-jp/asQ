<template>
  <!-- <div class="w-100 h-100">
    <h1>Chat Room</h1>
    <v-btn> aaa </v-btn>
  </div> -->
  <Header title="Chat Room" />
  <div class="chat-room-container">
    <!-- 受信したAIメッセージを順次表示 -->
    <div v-for="(msg, idx) in aiMessages" :key="idx" class="message-wrapper">
      <AiMessage :id="msg.id" :message="msg.message" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, onMounted, onBeforeUnmount } from 'vue'
import { useRoute } from 'vue-router'
import Header from '@/components/HeaderComponent.vue'
import AiMessage from '@/components/AiMessage.vue'

// AIメッセージのリスト（characterId とメッセージ本文）
type AiMessage = { id: string; message: string }
const aiMessages = reactive<AiMessage[]>([])

let currentIndex = -1 // 現在構築中のメッセージのインデックス

onMounted(() => {
  const route = useRoute()
  const chatId = route.params.id as string
  // ChatID が無い場合は処理しない
  if (!chatId) return

  const es = new EventSource(`/api/sse/ai/${chatId}`)

  // 会話開始：characterId が渡されるので新しいメッセージを作成
  es.addEventListener('start', (e: MessageEvent) => {
    try {
      const { characterId } = JSON.parse(e.data)
      aiMessages.push({ id: characterId, message: '' })
      currentIndex = aiMessages.length - 1
    } catch (err) {
      console.error('Failed to parse start event', err)
    }
  })

  // メッセージの断片を受信
  es.addEventListener('data', (e: MessageEvent) => {
    if (currentIndex === -1) return // start がまだ来ていない
    try {
      const { message } = JSON.parse(e.data)
      aiMessages[currentIndex].message += message
    } catch (err) {
      console.error('Failed to parse data event', err)
    }
  })

  // 会話終了
  es.addEventListener('close', () => {
    es.close()
    currentIndex = -1
  })

  onBeforeUnmount(() => {
    es.close()
  })
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

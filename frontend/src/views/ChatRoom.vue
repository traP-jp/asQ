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
        <div
          style="display: flex; justify-content: center; width: 100%; height: 15%; margin-top: 10px"
        >
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
import { useRoute } from 'vue-router'
import { createTypewriter, type TypewriterMessage } from '@/utils/typewriter'
import { createStreamManager, type StreamMessage, type StreamEventHandlers } from '@/utils/streams'
import ChooseCharacters from '@/components/ChooseCharacters.vue'

const aiMessages = ref<TypewriterMessage[]>([])
const selectedCharacterId = ref<string>('d7ec5756-4f20-11f0-ae2b-6ae4cc62b771')

// URLパラメータからchatIdを取得（setup関数の直接スコープで実行）
const route = useRoute()
const chatId = route.params.id as string

// タイプライターコントローラーを初期化
const typewriter = createTypewriter(30) // 30ms間隔でより滑らかに

// チャットメッセージ一覧
const chatMessages = ref<StreamMessage[]>([])

// ストリームマネージャーを初期化
const streamConfig = {
  chatId,
  typewriter,
  chatMessages,
  aiMessages,
}


const streamManager = createStreamManager(streamConfig)

// メッセージ送信ハンドラー
const handleSendMessage = async (message: string) => {
  try {
    await streamManager.sendMessage(message, selectedCharacterId.value)
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

  try {
    // 既存のチャットログを読み込み
    await streamManager.loadChatHistory()

    // 発言状態監視SSEを開始
    await streamManager.startSpeakingStatusStream()
  } catch (err) {
    console.error('Failed to initialize chat:', err)
  }
})

onBeforeUnmount(() => {
  streamManager.cleanup()
})

import UserMessage from '@/components/UserMessage.vue'
import InputText from '@/components/InputText.vue'
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

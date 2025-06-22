<template>
  <div class="whole-page">
    <div class="header">
      <Header title="Chat List" />
    </div>
    <div class="explain">
      <div class="hero-content">
        <div class="icon-wrapper">
          <div class="chat-icon">💬</div>
          <div class="sparkle">✨</div>
        </div>
        <h2 class="hero-title">
          <span class="gradient-text">AI</span>に<span class="highlight">trap</span>のことを聞こう!!
        </h2>
        <p class="hero-description">
          <span class="pulse-dot">●</span>
          クリックして、chat roomを作成しよう!!
          <span class="pulse-dot">●</span>
        </p>
        <div class="decorative-line"></div>
      </div>
    </div>
    <div class="start-chat">
      <v-btn
        v-for="info in aiInfo"
        :key="info.id"
        height="90%"
        width="200"
        class="create-chat"
        @click="addNewRoom(info.id)"
      >
        <div class="contents">
          <AiIcon :imageUrl="info.iconUrl" style="width: 2.5rem; height: 2.5rem" />
          <div class="text">{{ info.description }}</div>
          <div class="text-start">チャットを始める</div>
        </div>
      </v-btn>
    </div>
    <div class="chat-history">
      <div class="chat-list">
        <RoomCard
          v-for="room in rooms"
          :key="room.id"
          :aiId="room.characterId"
          :message="room.title"
          :time="room.createdAt"
          :roomId="room.id"
          :userIcons="room.userIds"
          class="room-card"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import Header from '@/components/HeaderComponent.vue'
import RoomCard from '@/components/ChatRoomCard.vue'
import { useRouter } from 'vue-router'
import AiIcon from '@/components/AiIcon.vue'
import api from '@/utils/api'

interface AiEntry {
  id: string
  description: string
  iconUrl: string
}
const aiInfo = ref<AiEntry[]>([])

type Room = {
  id: string
  characterId: string
  title: string
  createdAt: string
  userIds: string[]
}

const rooms = ref<Room[]>([])

const router = useRouter()

const addNewRoom = async (aiId: string) => {
  const { data } = await api.post('/api/chats')

  // チャットルームページに遷移
  router.push(`/chat/${data.id}`)
}

onMounted(async () => {
  try {
    const response = await api.get('/api/chats')
    const chats: Room[] = response.data.chats
    if (!Array.isArray(chats)) {
      throw new Error('APIから配列が返ってきませんでした')
    }

    rooms.value = [
      ...chats.map((chat) => ({
        id: chat.id,
        characterId: chat.characterId ?? '',
        title: chat.title ?? '',
        createdAt: chat.createdAt
          ? new Date(chat.createdAt).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
          : '',
        userIds: Array.isArray(chat.userIds) ? chat.userIds : [],
      })),
      ...rooms.value,
    ]
  } catch (e) {
    console.error('チャット一覧の取得に失敗:', e)
  }

  try {
    const { data } = await api.get('/api/characters')
    aiInfo.value = data.characters
  } catch (e) {
    console.error('キャラクターの取得に失敗:', e)
  }
})
</script>

<style scoped>
.explain {
  text-align: center;
  margin: 30px auto;
  max-width: 800px;
  padding: 0 20px;
}

.hero-content {
  background: linear-gradient(145deg, rgba(255, 255, 255, 0.9), rgba(255, 255, 255, 0.7));
  backdrop-filter: blur(10px);
  border-radius: 20px;
  padding: 40px 30px;
  box-shadow:
    0 8px 32px rgba(0, 0, 0, 0.1),
    0 2px 8px rgba(0, 0, 0, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.3);
  position: relative;
  overflow: hidden;
}

.hero-content::before {
  content: '';
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: conic-gradient(from 0deg, transparent, rgba(170, 213, 249, 0.1), transparent);
  animation: rotate 20s linear infinite;
  z-index: -1;
}

.icon-wrapper {
  position: relative;
  display: inline-block;
  margin-bottom: 20px;
}

.chat-icon {
  font-size: 3rem;
  display: inline-block;
  animation: bounce 2s ease-in-out infinite;
  filter: drop-shadow(0 4px 8px rgba(0, 0, 0, 0.1));
}

.sparkle {
  position: absolute;
  top: -10px;
  right: -10px;
  font-size: 1.5rem;
  animation: sparkle 1.5s ease-in-out infinite alternate;
}

.hero-title {
  font-size: 2.2rem;
  font-weight: bold;
  margin: 20px 0;
  line-height: 1.3;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.gradient-text {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  font-weight: 900;
}

.highlight {
  color: #ff6b6b;
  font-weight: 900;
  text-shadow: 0 0 10px rgba(255, 107, 107, 0.3);
}

.hero-description {
  font-size: 1.2rem;
  color: #555;
  margin: 20px 0;
  font-weight: 500;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
}

.pulse-dot {
  color: #ff6b6b;
  animation: pulse 2s ease-in-out infinite;
}

.decorative-line {
  height: 3px;
  width: 80px;
  background: linear-gradient(90deg, #667eea, #764ba2, #ff6b6b);
  border-radius: 2px;
  margin: 20px auto 0;
  animation: shimmer 3s ease-in-out infinite;
}

@keyframes rotate {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

@keyframes bounce {
  0%,
  20%,
  50%,
  80%,
  100% {
    transform: translateY(0);
  }
  40% {
    transform: translateY(-10px);
  }
  60% {
    transform: translateY(-5px);
  }
}

@keyframes sparkle {
  0% {
    opacity: 0.5;
    transform: scale(0.8) rotate(0deg);
  }
  100% {
    opacity: 1;
    transform: scale(1.2) rotate(180deg);
  }
}

@keyframes pulse {
  0%,
  100% {
    opacity: 0.6;
    transform: scale(1);
  }
  50% {
    opacity: 1;
    transform: scale(1.1);
  }
}

@keyframes shimmer {
  0%,
  100% {
    opacity: 0.7;
  }
  50% {
    opacity: 1;
    transform: scaleX(1.1);
  }
}

/* レスポンシブ対応 */
@media (max-width: 768px) {
  .hero-content {
    padding: 30px 20px;
    margin: 20px;
  }

  .hero-title {
    font-size: 1.8rem;
  }

  .hero-description {
    font-size: 1rem;
  }

  .chat-icon {
    font-size: 2.5rem;
  }
}

.create-chat {
  margin: 1rem;
}
.start-chat {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 24vh;
  margin: 15px;
}
.chat-history {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.whole-page {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100vw;
  min-height: 100vh;
  background: linear-gradient(135deg, #aad5f9 0%, #f5dcfe 100%);
}

.chat-list {
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-start;
  margin: 5% auto;
  max-width: 1150px;
  width: 100%;
  /* デフォルト: PC表示のpaddingは36px */
  padding: 0 36px;
  gap: 0;
  box-sizing: border-box;
}

.chat-list > * {
  flex: 0 1 calc((100% - 72px) / 3); /* 3列・左右padding分を引く */
  max-width: calc((100% - 72px) / 3);
  min-width: 330px;
  box-sizing: border-box;
}

/* 画面幅が狭くなっても3列を優先し、paddingを減らす */
@media (max-width: 1150px) {
  .chat-list {
    padding: 0 16px;
  }
  .chat-list > * {
    flex: 0 1 calc((100% - 32px) / 3);
    max-width: 370px;
  }
}

@media (max-width: 1100px) {
  .chat-list {
    padding: 0 4px;
  }
  .chat-list > * {
    flex: 0 1 calc((100% - 8px) / 3);
    max-width: 370px; /* 3列の最大幅を設定 */
    min-width: 330px; /* 最小幅を設定 */
  }
}

/* 画面幅がカード3枚分＋paddingより小さくなったら2列に */
@media (max-width: 1020px) {
  .chat-list {
    padding: 0 4px;
    justify-content: center;
  }
  .chat-list > * {
    flex: 0 1 50%;
    max-width: 370px;
    min-width: 330px;
  }
}

/* スマホは1列 */
@media (max-width: 700px) {
  .chat-list {
    padding: 0 2px;
    justify-content: center;
  }
  .chat-list > * {
    flex: 1 1 100%;
    max-width: 370px;
    min-width: 330px;
  }
}

.create-chat {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  background-color: #f0f0f0;
  border-radius: 8px;
  box-shadow: 0 3px 6px rgba(0, 0, 0, 0.1);
  gap: 8px;
}
.text-start {
  margin: 10px;
  width: 100%;
  height: 90%;
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 16px;
  color: #333;
}
.contents {
  width: 100%;
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  align-items: center;
}
</style>

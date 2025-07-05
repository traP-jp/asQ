<template>
  <div class="whole-page">
    <div class="header">
      <Header title="asQ" />
    </div>
    <div class="text-center my-10">
      <h2 class="text-h4 font-weight-bold mb-3">
        <span class="text-primary">AI</span>に<span class="text-pink">traP</span>のことを聞こう!!
      </h2>
      <p class="text-subtitle-1">クリックして、chat roomを作成しよう!!</p>
    </div>
    <div class="d-flex justify-center align-center w-100" style="height: 24vh; margin: 15px">
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
    <div class="d-flex flex-column align-center w-100 h-100">
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
  // チャットルームの取得
  try {
    const response = await api.get('/api/chats')
    if (!Array.isArray(response.data.chats)) {
      throw new Error('APIから配列が返ってきませんでした')
    }

    rooms.value = [
      ...response.data.chats.map((chat: any) => ({
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

  // キャラクターの取得
  try {
    const { data } = await api.get('/api/characters')
    aiInfo.value = data.characters
  } catch (e) {
    console.error('キャラクターの取得に失敗:', e)
  }
})
</script>

<style scoped>
.create-chat {
  margin: 1rem;
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

/* 復元: 画面全体背景などのスタイル */
.whole-page {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100vw;
  min-height: 100vh;
  background: linear-gradient(135deg, #aad5f9 0%, #f5dcfe 100%);
}

/* チャットルームカード一覧 (CSS Grid でシンプルに) */
.chat-list {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(330px, 1fr));
  gap: 24px;
  max-width: 1150px;
  width: 100%;
  margin: 5% auto;
  padding: 0 16px;
  box-sizing: border-box;
}

.chat-list > * {
  width: 100%; /* Grid がサイズを制御するため明示 */
}

/* ───────── ❶ 3 列（PC） ───────── */
@media (min-width: 1030px) {
  /* 900px 以上なら横に 3 枚並べる */
  .chat-list {
    grid-template-columns: repeat(3, 1fr);
  }
}

/* ───────── ❷ 2 列（タブレット） ───────── */
@media (min-width: 768px) and (max-width: 1029.98px) {
  /* 600px〜899.98px なら 2 枚 */
  .chat-list {
    grid-template-columns: repeat(2, 1fr);
  }
}

/* ───────── ❸ 1 列（スマホ） ───────── */
@media (max-width: 767.98px) {
  /* 599.98px 以下は縦 1 列 */
  .chat-list {
    grid-template-columns: 1fr;
  }
}
</style>

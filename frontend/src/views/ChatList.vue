<template>
  <div class="whole-page">
    <div class="header">
      <Header title="Chat List" />
    </div>
    <div class="explain">
      <h2>Chat List</h2>
      <p>Click on a chat room to start chatting with AI!</p>
    </div>
    <div class="start-chat">
      <v-btn
        v-for="info in aiInfo"
        :key="info.aiId"
        height="90%"
        width="200"
        class="create-chat"
        @click="addNewRoom(info.aiId)"
      >
        <div class="contents">
          <AiIcon :imageUrl="info.imageUrl" style="width: 2.5rem; height: 2.5rem" />
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
  margin: 20px;
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
  height: calc(100vh - 9vh);
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

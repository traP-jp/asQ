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
          <AiIcon :id="info.imageUrl" />
          <div class="text">{{ info.description }}</div>
          <div class="text-start">Chat start!</div>
        </div>
      </v-btn>
    </div>

    <div class="chat-history">
      <div class="chat-list">
        <RoomCard
          v-for="room in rooms"
          :key="room.roomId"
          :aiId="room.aiId"
          :message="room.message"
          :time="room.time"
          :roomId="room.roomId"
          :userIcons="room.userIcons"
          class="room-card"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import Header from '@/components/HeaderComponent.vue'
import UserMessage from '@/components/UserMessage.vue'
import AiMessage from '@/components/AiMessage.vue'
import UserIcon from '@/components/UserIcon.vue'
import RoomCard from '@/components/ChatRoomCard.vue'
import { useRouter } from 'vue-router'
import AiIcon from '@/components/AiIcon.vue'

import { onMounted } from 'vue'

interface AiEntry {
  aiId: string
  description: string
  imageUrl: string
}
const aiInfo = ref<AiEntry[]>([
  {
    aiId: 'ai1',
    description: 'AI 1',
    imageUrl: 'https://q.trap.jp/api/v3/public/icon/ai1',
  },
  {
    aiId: 'ai2',
    description: 'AI 2',
    imageUrl: 'https://q.trap.jp/api/v3/public/icon/ai2',
  },
  // 他のAI情報も追加可能
])

// const aiInfo = ref<AiEntry[]>([])

// async () => {
//   const res = await fetch('https://your-api.com/ai-info')
//   aiInfo.value = await res.json()
// }

const rooms = ref([
  {
    roomId: '1',
    aiId: 'ai1',
    message: 'Hello from AI 1',
    time: '10:00',
    userIcons: ['user1', 'user2'],
  },
  {
    roomId: '2',
    aiId: 'ai2',
    message: 'Hello from AI 2',
    time: '11:00',
    userIcons: ['user3', 'user4'],
  },
])

const router = useRouter()
const nextRoomId = ref(3)

const addNewRoom = (aiId: string) => {
  const newRoom = {
    roomId: String(nextRoomId.value++),
    aiId,
    message: 'New chat started!',
    time: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
    userIcons: ['newUser'],
  }

  // rooms.value.unshift(newRoom)

  // チャットルームページに遷移
  router.push(`/chat/${newRoom.roomId}`)
}
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

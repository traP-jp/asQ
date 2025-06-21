<script setup lang="ts">
import { ref } from 'vue'
import UserIcon from './UserIcon.vue'
import AiMessage from './AiMessage.vue'

const props = defineProps<{
  aiId: string
  message: string
  time: string
  roomId: string
  userIcons: string[]
}>()

const text = ref<string>(props.message.substring(0, 20) + (props.message.length > 20 ? '...' : ''))
const userIcons = ref<string[]>(props.userIcons.slice(0, 3))
</script>

<template>
  <div class="room-card" @click="$router.push(`${$route.fullPath}chat/${roomId}`)">
    <div class="chat-info">
      <AiMessage :message="text" :id="aiId" class="chat" />
    </div>
    <div class="message-time">{{ time }}</div>

    <div class="members-info">
      <div class="users">
        <UserIcon v-for="id in userIcons" :id="id" :key:="id" class="user-icon" />
      </div>
    </div>
  </div>
</template>

<style>
.room-card {
  background-color: #fff;
  width: 350px;
  height: 150px;
  display: flex;
  flex-direction: column;
  /* justify-content: space-between; */
  padding: 10px;
  margin: 10px;
  border: 1px solid #cccccc;
  border-radius: 8px;
  box-shadow: 0 3px 6px rgba(0, 0, 0, 0.05);
  transition-duration: 0.2s;
}
.room-card:hover {
  transform: scale(1.07);
}
.chat-info {
  height: 50%;
  width: 90%;
  display: flex;
}
.chat .text-ai {
  max-width: 100%;
  background-color: #f0f0f0;
}
.chat .text-ai::after {
  border-right: 10px solid #f0f0f0; /* 吹き出しと同じ色で三角の左辺 */
}
.latest-member {
  padding: 5px;
  display: flex;
}
.latest-message {
  width: 80%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
.message-text {
  background-color: #dedede;
  border-radius: 8px;
  font-size: 1.2rem;
  padding: 5px;
}
.message-time {
  font-size: 0.8rem;
  color: #888888;
  margin: 5px;
  text-align: right;
}
.members-info {
  height: 35%;
  display: flex;
  justify-content: flex-end;
  align-items: center;
}
/* .users {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 5px;
  padding-right: 5px;
  margin: 0 0 0 auto;
} */

.users {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 5px;
  padding: 5px;
  /* margin: 0 0 0 auto; */
}
.user-icon img {
  width: 30px;
  height: 30px;
}
</style>

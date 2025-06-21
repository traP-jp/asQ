<script setup lang="ts">
import { ref } from 'vue'
import UserIcon from './UserIcon.vue'

const { aiId, message, time } = defineProps<{
  aiId: string
  message: string
  time: string
}>()

const text = ref<string>(message.substring(0, 17) + (message.length > 17 ? '...' : ''))

const userIcons = ref<string[]>(['mina', 'mumumu', 'kaomojikun'])
</script>

<template>
  <div class="room-card">
    <div class="chat-info">
      <div class="latest-member"><UserIcon :id="aiId" /></div>
      <div class="latest-message">
        <div class="message-text">{{ text }}</div>
        <div class="message-time">{{ time }}</div>
      </div>
    </div>
    <div class="members-info">
      <div class="users">
        <UserIcon v-for="id in userIcons" :id="id" />
      </div>
    </div>
  </div>
</template>

<style>
.room-card {
  background-color: #fff;
  width: 30%;
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
  height: 55%;
  display: flex;
}
.latest-member {
  padding: 5px;
  display: flex;
}
.latest-message {
  display: flex;
  flex-direction: column;
}
.message-text {
  font-size: 1.2rem;
  padding: 5px;
}
.message-time {
  font-size: 0.8rem;
  color: #888888;
  margin: 5px;
  padding-right: 5px;
  margin: 0 0 0 auto;
}
.members-info {
  height: 45%;
  display: flex;
  justify-content: flex-end;
  align-items: center;
}
.users {
  width: 20%;
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 5px;
  padding-right: 5px;
  /* margin: 0 0 0 auto; */
}
</style>

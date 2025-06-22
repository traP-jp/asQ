<script setup lang="ts">
import { ref, onMounted } from 'vue'
import UserIcon from './UserIcon.vue'
import api from '@/utils/api'
const { title } = defineProps<{
  title: string
}>()

const userId = ref<string>('')
const error = ref<string | null>(null)

// onMounted(async () => {
//   chatIdの存在確認
//   if (!chatId) {
//     console.error('Chat ID is not provided in route parameters')
//     return
//   }
// const { data } = await api.get('/api/users/me')
// const userId = data.userId

onMounted(async () => {
  try {
    const { data } = await api.get('/api/users/me')
    userId.value = data.userId
    console.log('User ID fetched:', userId.value)
  } catch (err) {
    console.error('Failed to fetch user data:', err)
    error.value = 'ユーザー情報の取得に失敗しました'
  } 
})
</script>

<template>
  <div class="header">
    <div class="home-button">
      <v-btn variant="text" to="/" :ripple="false">Home</v-btn>
    </div>
    <div class="header-text">{{ title }}</div>
    <div class="header-icon">
      <UserIcon :id="userId" :isHover="false" />
    </div>
  </div>
</template>

<style scoped>
.header {
  width: 100vw;
  height: 9vh;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 3px 6px rgba(0, 0, 0, 0.05);
}
.home-button {
  font-size: 1.5rem;
  margin: 1.5rem;
}
.header-text {
  font-size: 2rem;
  font-weight: bold;
  margin: 1rem;
}
.header-icon {
  margin: 1.5rem;
  width: 2.5rem;
  height: 2.5rem;
}
</style>

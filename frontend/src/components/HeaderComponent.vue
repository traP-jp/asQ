<script setup lang="ts">
import { ref } from 'vue'
import UserIcon from './UserIcon.vue'
const { title } = defineProps<{
  title: string
}>()
const userId = ref<string>('mina')
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

<style>
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
}
</style>

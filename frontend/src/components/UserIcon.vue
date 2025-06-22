<script setup lang="ts">
import { ref, watch } from 'vue'

const { id, isHover = true } = defineProps<{
  id: string
  isHover?: boolean
}>()

const imageUrl = ref(`https://q.trap.jp/api/v3/public/icon/${id}`)

// id が変化したら imageUrl を再構築
watch(() => id, (newId) => {
  imageUrl.value = `https://q.trap.jp/api/v3/public/icon/${newId}`
})
</script>

<template>
  <div>
    <img :src="imageUrl" class="usericon" />
    <v-tooltip v-if="isHover" activator="parent" location="bottom">{{ id }}</v-tooltip>
  </div>
</template>

<style scoped>
.usericon {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  box-shadow: 0 3px 6px rgba(0, 0, 0, 0.05);
}
</style>

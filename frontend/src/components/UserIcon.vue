<script setup lang="ts">
import { ref } from 'vue'

const imageUrl = ref('https://q.trap.jp/api/v3/public/icon/')

const { id, isHover = true } = defineProps<{
  id: string
  isHover?: boolean
}>()

const userId = ref<string>(id)
imageUrl.value += userId.value
</script>

<template>
  <div>
    <img :src="imageUrl" class="usericon" />
    <v-tooltip v-if="isHover" activator="parent" location="bottom">{{ id }}</v-tooltip>
  </div>
</template>

<style>
.usericon {
  width: 3rem;
  height: 3rem;
  border-radius: 50%;
  box-shadow: 0 3px 6px rgba(0, 0, 0, 0.05);
}
</style>

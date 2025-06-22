<script setup lang="ts">
defineProps<{
  imageUrl: string
}>()
</script>

<template>
  <div class="ai-icon">
    <img :src="imageUrl" class="ai-icon" />
  </div>
</template>

<style scoped>
.ai-icon {
  background-color: white;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  box-shadow: 0 3px 6px rgba(0, 0, 0, 0.05);
}
</style>

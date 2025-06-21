<script setup lang="ts">
const props = defineProps<{
  id: string
}>()

const images = import.meta.glob('@/assets/*.png', { eager: true, import: 'default' })
const logo : string = images[`/src/assets/${props.id}.png`] as string || images['/src/assets/default.png'] as string
</script>
<template>
  <div class="ai-icon">
    <img :src="logo" class="ai-icon" />
    <v-tooltip activator="parent" location="bottom">{{props.id}}</v-tooltip>
  </div>
</template>

<style scoped>
.ai-icon {
  background-color: white;
  width: 3rem;
  height: 3rem;
  border-radius: 50%;
  box-shadow: 0 3px 6px rgba(0, 0, 0, 0.05);
}
</style>

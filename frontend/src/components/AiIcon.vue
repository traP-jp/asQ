<script setup lang="ts">
const props = defineProps<{
  id: string
}>()
import ai1 from '@/assets/ai1.png'
import ai2 from '@/assets/ai2.png'
import ai3 from '@/assets/ai3.png'

let logo :string
if (props.id === 'ai1') {
  logo = ai1
} else if (props.id === 'ai2') {
  logo = ai2
} else if (props.id === 'ai3') {
  logo = ai3
}

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
  width: 100%;
  height: 100%;
  border-radius: 50%;
  box-shadow: 0 3px 6px rgba(0, 0, 0, 0.05);
}
</style>

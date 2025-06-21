<script setup lang="ts">
import { ref } from 'vue';

const imageUrl = ref('https://q.trap.jp/api/v3/public/icon/')

const props = defineProps<{
  id: string
}>()

const userId = ref<string>(props.id)
imageUrl.value += userId.value


</script>


<template>
    <div class="usericon">
        <img: :src=imageUrl class="usericon"/>
    </div>
</template>

<style>
.usericon {
  width: 3rem;
  height: 3rem;
  border-radius: 50%;
}

</style>
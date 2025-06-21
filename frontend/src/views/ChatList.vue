<template>
  <div class="d-flex flex-column align-center chat-list">
    <h1>Chat List</h1>

    <!-- 之だけでボタンがいい感じになる！凄い！ -->
    <v-btn class="my-4">aaaaaa</v-btn>

    <v-progress-circular model-value="20" indeterminate></v-progress-circular>
    <!-- <v-toolbar title="Toolbar"></v-toolbar> -->
    <div class="w-50 h-20">
    <v-combobox
      label="Combobox"
      variant="solo"
      :items="['California', 'Colorado', 'Florida', 'Georgia', 'Texas', 'Wyoming']"
    />
    </div>

    
    <UserIcon id="mina" />
  </div>
</template>

<script setup lang="ts">
import UserIcon from '@/components/UserIcon.vue';
</script>

<style scoped>
.chat-list {
  height: 100vh;
  width: 100vw;
  display: flex;
  flex-direction: column;
  align-items: center;
  background: linear-gradient(135deg, #42a5f5 0%, #b147d1 100%);
}

</style>

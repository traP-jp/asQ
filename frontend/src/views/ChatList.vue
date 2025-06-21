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

    <UserMessage :message="messageSample" id="mina" />
    <v-btn icon="mdi-send-variant" />
    <AiMessage :message="messageSample" id="ai1" />
    <UserIcon id="mina" />
    <RoomCard aiId="cp20" message="aaaaaaaaaaaaaaaaaaaaaaa" time="now" roomId="3" :userIcons="['mina', 'mumumu', 'kaomojikun']"/>
  </div>
</template>

<script setup lang="ts">
import UserMessage from '@/components/UserMessage.vue';
import AiMessage from '@/components/AiMessage.vue';
import UserIcon from '@/components/UserIcon.vue'
import RoomCard from '@/components/ChatRoomCard.vue';
import InputText from '@/components/InputText.vue';
import AiIcon from '@/components/AiIcon.vue';

const messageSample: string = `
# Hello

- list
- list
- list


ああああああああああああああああ

てすてすてす

# h1
## h2
### h3
#### h4
##### h5
###### h6


`
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

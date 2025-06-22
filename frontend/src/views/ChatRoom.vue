<script setup lang="ts">
import AiMessage from '@/components/AiMessage.vue'
import Header from '@/components/HeaderComponent.vue'
import UserMessage from '@/components/UserMessage.vue'
import InputText from '@/components/InputText.vue'
import { ref } from 'vue'
type Messages = {
  id: string
  message: string
  time: Date
  isAi?: boolean
}

type ChatMessages = {
  userMessages: Messages[]
  aiMessages: Messages[]
}

const chatMessagesResponse = ref<ChatMessages>({
  userMessages: [
    {
      id: 'user1',
      message: 'Hello, how are you?',
      time: new Date(2022, 0, 1, 0, 0, 0),
      isAi: false,
    },
    {
      id: 'user2',
      message: 'What is the weather like today?',
      time: new Date(2022, 0, 1, 0, 10, 0),
      isAi: false,
    },
  ],
  aiMessages: [
    {
      id: 'ai1',
      message: 'I am fine, thank you!',
      time: new Date(2022, 0, 1, 0, 0, 5),
      isAi: true,
    },
    {
      id: 'ai2',
      message: 'The weather is sunny today.',
      time: new Date(2022, 0, 1, 0, 10, 5),
      isAi: true,
    },
  ],
})

const chatMessages = ref<Messages[]>([])

chatMessages.value = [
  ...chatMessagesResponse.value.userMessages,
  ...chatMessagesResponse.value.aiMessages,
]

chatMessages.value.sort((a, b) => a.time.getTime() - b.time.getTime())
</script>

<template>
  <Header title="Chat Room"></Header>
  <div class="chat-room">
    <div class="chat-space">
      <div class="chat-container">
        <div v-for="(message, index) in chatMessages" :key="index">
          <div class="messag">
            <UserMessage v-if="!message.isAi" :id="message.id" :message="message.message" />
            <AiMessage v-else :id="message.id" :message="message.message" />
          </div>
        </div>
      </div>
      <InputText
        class="input-text"
        @sendMessage="
          (message: string) => {
            const newMessage: Messages = {
              id: `user${chatMessagesResponse.userMessages.length + 1}`,
              message,
              time: new Date(),
              isAi: false,
            }
            chatMessagesResponse.userMessages.push(newMessage)
            chatMessages = [...chatMessagesResponse.userMessages, ...chatMessagesResponse.aiMessages]
            chatMessages.sort((a, b) => a.time.getTime() - b.time.getTime())
          }
        "
      />
    </div>
  </div>
</template>
<style scoped>
.chatspace{
  flex-direction: column;
  width: 80%;
}
.chat-container{
  background-color: aliceblue;
  overflow-y: scroll;
}



</style>

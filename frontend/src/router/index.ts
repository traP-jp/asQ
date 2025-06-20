import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'chat-list',
      component: () => import('../views/ChatList.vue'),
    },
    {
      path: '/chat/:id',
      name: 'chat',
      component: () => import('../views/ChatRoom.vue'),
    },
  ],
})

export default router

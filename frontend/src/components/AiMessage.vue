<script setup lang="ts">
import { computed } from 'vue'
import AiIcon from '@/components/AiIcon.vue'
import Markdown from '@/components/Markdown.vue'

const props = defineProps<{
  message: string
  imageUrl: string
}>()
</script>

<template>
  <div class="ai-message">
    <AiIcon :imageUrl="props.imageUrl" style="width: 2.5rem; height: 2.5rem" />
    <div class="text-ai">
      <Markdown :text="props.message" />
    </div>
  </div>
</template>

<style scoped>
.ai-message {
  display: flex;
  width: 100%;
}
.text {
  display: flex;
  align-items: center;
}
.text-ai {
  display: flex;
  justify-content: center;
  /* 吹き出しの基本的なスタイル */
  flex: 1;
  background-color: #86e2e5; /* 明るい緑や水色など、ユーザー側を表す色 */
  color: #333;
  /* padding: 10px 15px; */
  border-radius: 10px; /* 角を丸くする */
  max-width: 70%; /* 吹き出しの最大幅。長文でも画面幅を占領しないように */
  word-wrap: break-word; /* 長い単語やURLがはみ出さないように */

  /* 右寄せ */
  margin-left: 20px; /* 左側のマージンを自動調整して右寄せ */
  margin-right: auto; /* 右側のマージン（任意） */
  margin-bottom: 10px; /* 他の吹き出しとの間隔 */

  /* 吹き出しの「しっぽ」（三角）の準備 */
  position: relative; /* 疑似要素の基準位置 */
  box-shadow: 0px 10px 10px -6px rgba(0, 0, 0, 0.3);
}

/* 吹き出しの「しっぽ」 (疑似要素で作成) */
.text-ai::after {
  content: '';
  position: absolute;
  left: -8px; /* 吹き出しの右端から外に出す位置 */
  top: 12px; /* 吹き出しの上端からの位置 */
  width: 0;
  height: 0;
  border-right: 10px solid #86e2e5; /* 吹き出しと同じ色で三角の左辺 */
  border-top: 10px solid transparent; /* 透明な上辺 */
  border-bottom: 10px solid transparent; /* 透明な下辺 */
}
</style>

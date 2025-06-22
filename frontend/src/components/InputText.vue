<script setup lang="ts">
import { ref } from 'vue'

const multilineText = ref('')

// イベントを定義
const emit = defineEmits<{
  sendMessage: [message: string]
}>()

// 送信ボタンのクリックハンドラー
const handleSend = () => {
  if (multilineText.value.trim()) {
    emit('sendMessage', multilineText.value.trim())
    multilineText.value = '' // 送信後にテキストフィールドをクリア
  }
}

// Enterキーでの送信（Shift+Enterは改行）
const handleKeydown = (event: KeyboardEvent) => {
  if (event.key === 'Enter' && !event.shiftKey) {
    event.preventDefault()
    handleSend()
  }
}
</script>

<template>
  <div class="input-text my-10">
    <textarea
      class="form"
      v-model="multilineText"
      placeholder="複数行のテキストを入力"
      rows="1"
      @keydown="handleKeydown"
    ></textarea>

    <v-btn class="send-button" color="primary" @click="handleSend">
      送信
      <v-icon class="button-icon" icon="mdi-send-variant"></v-icon>
    </v-btn>
  </div>
</template>

<style scoped>
.input-text {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
  width: 80%;
  background-color: #fbf4fe;
  padding: 10px;
  border-radius: 15px; /* 角を丸くする */
  box-shadow: 0px 10px 10px -6px rgba(0, 0, 0, 0.3);
}
.form {
  flex: 1;
  padding: 10px;
  border: 1px solid #fbf4fe; /* 入力フィールドの枠線 */
  border-radius: 15px; /* 角を丸くする */
  font-size: 16px; /* フォントサイズ */
  margin-right: 10px; /* ボタンとの間隔 */
  word-wrap: break-word;
  overflow-wrap: break-word;
  white-space: pre-wrap; /* 長い単語やURLがはみ出さないように */
  box-sizing: border-box; /* paddingとborderをwidth/heightに含める */
  font-family: inherit;
  font-size: 1em;
  line-height: 1.5em; /* 行の高さ (CSSで明確に指定推奨) */

  /* ここからが重要 */
  resize: none; /* ユーザーによるサイズ変更を無効にする */
  overflow-y: auto; /* コンテンツがはみ出したらスクロールバーを表示 */

  /* 最大高さを4行分に設定 */
  /* paddingとborderも考慮して正確な高さを計算する必要がある */
  /* 例: (line-height * 4行) + paddingTop + paddingBottom + borderTop + borderBottom */
  max-height: calc(
    1.5em * 4 + 8px * 2 + 1px * 2
  ); /* line-height(1.5em)*4行 + padding(8px*2) + border(1px*2) */
}
.send-button {
  display: flex;
  padding: 10px 20px; /* ボタンのパディング */
  border-radius: 15px; /* ボタンの角を丸くする */
  font-size: 16px;
  margin-top: 5px;
  margin-bottom: 5px; /* フォントサイズ */
}
/* フォーカス時の太枠を消す */
.button-icon {
  margin-left: 2px;
  flex: 1;
  justify-content: center;
}
input:focus,
textarea:focus {
  outline: none; /* フォーカス時の太枠を消す */
}

/* 特定の種類のinput要素に適用する場合 */
input[type='text']:focus,
input[type='email']:focus,
input[type='password']:focus {
  outline: none;
}

/* 特定のクラスを持つ要素に適用する場合 */
.my-custom-input:focus {
  outline: none;
}
</style>

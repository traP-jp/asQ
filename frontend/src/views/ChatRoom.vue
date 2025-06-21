<template>
  <div class="w-100 h-100" style="padding: 20px">
    <h1>Chat Room - Markdown Test</h1>

    <div style="margin-top: 20px; border: 1px solid #ccc; border-radius: 8px">
      <Markdown v-model:text="markdownText" :pad="16" style="height: 50%; width: 50%; overflow-y: auto" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import Markdown from '@/components/Markdown.vue'

const markdownText = ref(`# Markdownテストページ

## 見出し2

これは**太字**で、これは*斜体*です。

### コードブロック

\`\`\`javascript
function hello() {
  console.log("Hello, World!");
  return "success";
}
\`\`\`

### インラインコード

\`const message = "Hello"\` のようにインラインでコードを書けます。

### リスト

- 項目1
- 項目2
  - サブ項目1
  - サブ項目2
- 項目3

### 番号付きリスト

このテキストを編集して、リアルタイムでプレビューを確認できます！`)
</script>

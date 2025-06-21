<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { markedHighlight } from 'marked-highlight'
import { Marked } from 'marked'
import hljs from 'highlight.js'
import darkStyle from 'highlight.js/styles/github-dark.css?inline'
// ?inline をつけて読み込むことでCSSファイルもテキストデータになる

const text = defineModel<string>('text')

defineProps<{ pad?: number }>()

const html = computed(() => {
  const tempHtml = marked.parse(text.value || '') as string
  return tempHtml.replace('<pre><code>', '<pre><code class="hljs">')
  // <pre><code> のままでは表示がコードブロックの方式に従ってくれないので、末尾にクラス hljs を指定する
})

const marked = new Marked(
  markedHighlight({
    langPrefix: 'hljs language-',
    highlight(code, lang) {
      if (typeof lang === 'string' && lang.includes(':')) {
        lang = lang.substring(0, lang.indexOf(':'))
      }
      const language = hljs.getLanguage(lang) ? lang : 'plaintext'
      return hljs.highlight(code, { language }).value
    },
  }),
  {
    gfm: true, // GitHub Flavored Markdown を有効にする
    breaks: true, // 1段の改行を有効にする
  },
)

onMounted(async () => {
  const highlightStyleTag = document.createElement('style')
  highlightStyleTag.textContent = darkStyle
  document.head.appendChild(highlightStyleTag)
})
</script>

<template>
  <div :class="$style.container">
    <slot></slot>
    <div :class="$style.content">
      <div
        v-html="html"
        :class="$style.preview"
        :style="`padding: ${pad !== undefined ? pad : 12}px;`"
      ></div>
    </div>
  </div>
</template>

<style module>
.container {
  color: var(--color-text);
  position: relative;
  width: 100%;
  height: 100%;
  overflow-y: auto;
  z-index: 1;
}

.content {
  width: 100%;
  height: 100%;
  z-index: -1;
  overflow-y: auto;
}

.preview {
  width: 100%;
  height: fit-content;
  min-height: 100%;
  max-width: 1000px;
  margin: 0 auto;
}

/* 合宿のしおり、イベント告知文、個人メモのプレビューなどこのスタイルで統一する */
/* 使いやすくてクセのないスタイルであることが望ましい。基本的には traQ に寄せる */

.preview :global(p) {
  font-size: 1em;
  line-height: 1.4;
  margin-bottom: 0.3lh;
}

.preview :global(p):last-child {
  margin-bottom: 0 !important;
}

.preview :global(h1),
.preview :global(h2),
.preview :global(h3),
.preview :global(h4),
.preview :global(h5),
.preview :global(h6) {
  margin: 4px 0;
  font-weight: bold;
  letter-spacing: 0.05em;
}

.preview :global(h1):first-child,
.preview :global(h2):first-child,
.preview :global(h3):first-child,
.preview :global(h4):first-child,
.preview :global(h5):first-child,
.preview :global(h6):first-child {
  margin-top: 0 !important;
}

.preview :global(strong) {
  font-weight: bold;
}

.preview :global(code) {
  font-family: 'M PLUS Code Latin 60', 'M PLUS 1p';
  font-weight: 500;
  font-size: 11pt;
}

.preview :global(pre code) {
  margin: 0 4px 8px 4px !important;
  border-radius: 4px;
  max-height: 400px;
  overflow: scroll;
  background-color: var(--color-dark-gray) !important;
}

.preview :global(p code) {
  margin: 0 0px;
  background-color: var(--color-theme-pale);
  border-radius: 4px;
  border: 2px solid var(--color-theme-pale);
}

.preview :global(ul),
.preview :global(ol) {
  padding: 0px 0px 0px 20px;
  margin-top: 8px;
  margin-bottom: 8px;
}

.preview :global(li) {
  margin-bottom: 8px;
}

.preview :global(blockquote) {
  border-radius: 0px;
  color: var(--color-gray);
  padding: 0px 0px 0px 10px;
  border-left: 4px solid var(--color-gray);
  margin: 0 0 8px 4px;
}

.preview :global(tr:nth-child(2n + 1)) {
  background-color: #f0f0f0;
}

.preview :global(tr:nth-child(2n + 2)) {
  background-color: #e0e0e0;
}

.preview :global(th),
.preview :global(td) {
  padding: 4px 8px;
}

.preview :global(th) {
  font-weight: bold;
  background-color: var(--color-theme-pale) !important;
}

.preview :global(table) {
  border-collapse: collapse;
  margin: 0 4px 8px 4px !important;
  border-radius: 4px;
  overflow: hidden;
}

.preview :global(a) {
  color: #0066ff !important;
  text-decoration: none;
}

.preview :global(a):hover {
  text-decoration: underline;
}
</style>

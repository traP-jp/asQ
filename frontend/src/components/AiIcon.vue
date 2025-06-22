<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import api from '@/utils/api'

const props = defineProps<{
  /**
   * imageUrl には直接 URL が渡される場合と、キャラクター ID(UUID 等) が渡される場合の 2 パターンがある。
   *   - URL の場合   : そのまま <img> の src として利用
   *   - キャラクター ID の場合 : /api/characters から iconUrl を解決して利用
   */
  imageUrl?: string
}>()

const resolvedSrc = ref('')

// デフォルトアイコンURLを取得
const getDefaultIconUrl = () => {
  try {
    const url = new URL(`../assets/ai1.png`, import.meta.url).href
    console.log('Default icon URL generated:', url)
    return url
  } catch (error) {
    console.warn('Failed to generate default icon URL, using fallback:', error)
    return '/src/assets/ai1.png'
  }
}

// 最終的に表示するアイコンURL
const displaySrc = computed(() => {
  // resolvedSrcが有効な値でない場合はデフォルトアイコンを返す
  const src = resolvedSrc.value?.trim()
  const finalSrc = src && src.length > 0 ? src : getDefaultIconUrl()
  console.log('displaySrc computed:', { originalSrc: src, finalSrc, imageUrl: props.imageUrl })
  return finalSrc
})

/**
 * props.imageUrl を解決し、resolvedSrc に格納
 */
const resolveSrc = async () => {
  const original = (props.imageUrl ?? '').trim()

  // 空文字やnull/undefinedの場合は即座にデフォルトアイコンを使用
  if (!original || original.length === 0) {
    resolvedSrc.value = ''
    return
  }

  // URL( http(s):// または data: ) ならそのまま返す
  if (/^(https?:\/\/|data:)/.test(original)) {
    resolvedSrc.value = original
    return
  }

  // それ以外はキャラクター ID とみなしてAPI から検索
  try {
    const { data } = await api.get('/api/characters')

    if (data && data.characters && Array.isArray(data.characters)) {
      const character = data.characters.find(
        (c: { id: string; iconUrl: string }) => c.id === original,
      )
      if (character && character.iconUrl && character.iconUrl.trim().length > 0) {
        resolvedSrc.value = character.iconUrl.trim()
        return
      }
    }
  } catch (err) {
    console.warn('Failed to fetch character icon for ID:', original, err)
  }

  // API取得失敗またはキャラクターが見つからない場合、ローカルアセットを試行
  try {
    const assetUrl = new URL(`../assets/${original}.png`, import.meta.url).href
    // 実際にアセットが読み込めるかチェックはしないが、URLは生成する
    resolvedSrc.value = assetUrl
    return
  } catch (e) {
    console.warn('Local asset not found for:', original)
  }

  // 上記すべてに該当しない場合は、resolvedSrcを空文字にして
  // computed の displaySrc がデフォルトアイコンを返すようにする
  resolvedSrc.value = ''
}

onMounted(() => {
  console.log('AiIcon mounted with imageUrl:', props.imageUrl)
  // 初期状態でデフォルトアイコンを設定
  resolvedSrc.value = ''
  // 非同期でアイコンを解決
  resolveSrc()
})

// props.imageUrl が変更された時に再解決
watch(() => props.imageUrl, resolveSrc)

const handleImageError = () => {
  console.warn('Image loading failed, using default icon')
  resolvedSrc.value = getDefaultIconUrl()
}
</script>

<template>
  <div class="ai-icon-container">
    <!-- 常にアイコンを表示。displaySrcがデフォルトアイコンまたは解決されたアイコンを返す -->
    <img :src="displaySrc" class="ai-icon" alt="AI Icon" @error="handleImageError" />
  </div>
</template>

<style scoped>
.ai-icon-container {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.ai-icon {
  background-color: white;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  box-shadow: 0 3px 6px rgba(0, 0, 0, 0.05);
  object-fit: cover; /* アイコンのアスペクト比を保持 */
}
</style>

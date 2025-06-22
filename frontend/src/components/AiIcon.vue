<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
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

/**
 * props.imageUrl を解決し、resolvedSrc に格納
 */
const resolveSrc = async () => {
  const original = props.imageUrl ?? ''

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
      if (character && character.iconUrl) {
        resolvedSrc.value = character.iconUrl
        return
      }
    }
  } catch (err) {
    console.warn('Failed to fetch character icon for ID:', original, err)
  }

  // API取得失敗またはキャラクターが見つからない場合、ローカルアセットを試行
  if (original) {
    try {
      resolvedSrc.value = new URL(`../assets/${original}.png`, import.meta.url).href
      return
    } catch (e) {
      // ローカルアセットも見つからない場合
    }
  }

  // 最終的なフォールバック: デフォルト画像を表示
  try {
    resolvedSrc.value = new URL(`../assets/ai1.png`, import.meta.url).href
  } catch {
    resolvedSrc.value = ''
  }
}

onMounted(resolveSrc)

// props.imageUrl が変更された時に再解決
watch(() => props.imageUrl, resolveSrc)
</script>

<template>
  <div class="ai-icon">
    <!-- resolvedSrc が空の場合は何も表示しない -->
    <img v-if="resolvedSrc" :src="resolvedSrc" class="ai-icon" alt="AI Icon" />
  </div>
</template>

<style scoped>
.ai-icon {
  background-color: white;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  box-shadow: 0 3px 6px rgba(0, 0, 0, 0.05);
}
</style>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import api from '@/utils/api'

const props = defineProps<{
  /**
   * imageUrl には直接 URL が渡される場合と、キャラクター ID(UUID 等) が渡される場合の 2 パターンがある。
   *   - URL の場合   : そのまま <img> の src として利用
   *   - キャラクター ID の場合 : /api/characters から iconUrl を解決して利用
   */
  imageUrl?: string
}>()

// 簡易キャッシュ: モジュールスコープにキャラクター ID → iconUrl のマップを保持
let characterIconCache: Record<string, string> | null = null

const resolvedSrc = ref('')

/**
 * キャッシュが存在しない場合のみ /api/characters を叩いてキャッシュを構築
 */
const fetchCharacterIcons = async () => {
  if (characterIconCache) return // 既にキャッシュ済み
  try {
    const { data } = await api.get('/api/characters')

    // APIレスポンスの構造を確認してから処理
    if (data && data.characters && Array.isArray(data.characters)) {
      characterIconCache = Object.fromEntries(
        data.characters.map((c: { id: string; iconUrl: string }) => [c.id, c.iconUrl]),
      )
    } else {
      console.error('Invalid API response structure for character icons:', data)
      characterIconCache = {}
    }
  } catch (err) {
    console.error('Failed to fetch characters', err)
    characterIconCache = {}
  }
}

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

  // それ以外はキャラクター ID とみなしてキャッシュ(or API) から検索
  await fetchCharacterIcons()
  if (characterIconCache && characterIconCache[original]) {
    resolvedSrc.value = characterIconCache[original]
  } else {
    // ローカルアセット (assets ディレクトリ) に同名の画像があれば利用
    try {
      resolvedSrc.value = new URL(`../assets/${original}.png`, import.meta.url).href
    } catch (e) {
      // 見つからない場合はデフォルト画像
      try {
        resolvedSrc.value = new URL(`../assets/ai1.png`, import.meta.url).href
      } catch {
        resolvedSrc.value = ''
      }
    }
  }
}

onMounted(resolveSrc)
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

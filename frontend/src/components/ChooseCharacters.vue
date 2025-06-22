<script setup lang="ts">
import AiIcon from '@/components/AiIcon.vue'
import { api } from '@/utils/api'
import { onMounted, ref } from 'vue'

interface Character {
  id: string
  description: string
  iconUrl: string
  createdAt: string
}

const characters = ref<Character[]>([])

const modelValue = defineModel<string>({ default: 'ai1' })

onMounted(async () => {
  try {
    const { data } = await api.get('/api/characters')

    // APIレスポンスの構造を確認
    console.log('API response:', data)

    // data.charactersが配列であることを確認してから処理
    if (data && data.characters && Array.isArray(data.characters)) {
      characters.value = data.characters.map((char: Character) => ({
        id: char.id,
        description: char.description,
        iconUrl: char.iconUrl,
        createdAt: char.createdAt,
      }))
    } else {
      console.error('Invalid API response structure:', data)
      characters.value = []
    }

    console.log('Processed characters:', characters.value)
  } catch (error) {
    console.error('Failed to fetch characters:', error)
    characters.value = []
  }
})

// const cardColor = (active?: boolean) =>
//   active ? 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' : 'white'
</script>

<template>
  <v-card class="character-selector pa-4">
    <v-card-title class="text-h6 d-flex align-center mb-2">
      <v-icon icon="mdi-sparkles" size="20" class="mr-2" />
      相手を選択
    </v-card-title>
    <v-divider class="mb-2" />

    <!-- アイテムグループで単一選択 -->
    <v-item-group v-model="modelValue" mandatory>
      <v-container class="pa-0" fluid>
        <v-row no-gutters>
          <v-col cols="12" v-for="char in characters" :key="char.id" class="mb-2">
            <v-item :value="char.id" #="{ isSelected, toggle }">
              <v-card
                :elevation="isSelected ? 12 : 2"
                :class="[
                  'card-item d-flex align-center pa-3',
                  isSelected ? 'selected-card' : 'unselected-card',
                ]"
                @click="toggle"
              >
                <div class="flex-shrink-0 icon-wrapper">
                  <AiIcon :imageUrl="char.id" class="w-100 h-100" />
                </div>
                <div class="ml-4 text-left flex-grow-1">
                  <div
                    class="text-body-2"
                    :class="isSelected ? 'text-white' : 'text-grey-darken-1'"
                  >
                    {{ char.description }}
                  </div>
                </div>
              </v-card>
            </v-item>
          </v-col>
        </v-row>
      </v-container>
    </v-item-group>
  </v-card>
</template>

<style scoped>
.character-selector {
  width: 320px;
  /* カードの縁取りを柔らかく */
  border-radius: 16px;
}

.selected-card {
  transform: scale(1.02);
  box-shadow: 0 8px 25px rgba(102, 126, 234, 0.4) !important;
  border: 2px solid rgba(255, 255, 255, 0.3);
  transition:
    transform 0.4s ease,
    box-shadow 0.4s ease;
}

.unselected-card {
  transform: scale(1);
  border: 2px solid transparent;
  background: white;
  transition:
    transform 0.4s ease,
    box-shadow 0.4s ease;
}

.card-item {
  position: relative; /* 疑似要素配置のため */
  cursor: pointer;
  overflow: hidden; /* オーバーレイのはみ出し防止 */
  transition:
    transform 0.4s ease,
    box-shadow 0.4s ease;
}

/* グラデーションオーバーレイ */
.card-item::before {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: inherit;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  opacity: 0;
  transition: opacity 0.4s ease;
  pointer-events: none;
  z-index: 0;
}

.selected-card::before {
  opacity: 1;
}

/* コンテンツを前面に配置 */
.card-item > * {
  z-index: 1;
}

.selected-card:hover {
  transform: scale(1.03);
}

.unselected-card:hover {
  transform: scale(1.01);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15) !important;
}

.icon-wrapper {
  width: 2.5rem;
  height: 2.5rem;
}
</style>

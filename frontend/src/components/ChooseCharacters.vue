<script setup lang="ts">
import AiIcon from '@/components/AiIcon.vue'

interface Character {
  id: string
  description: string
}

const characters: Character[] = [
  {
    id: 'ai1',
    description: '親しみやすく、温かい語り口で回答します。',
  },
  {
    id: 'ai2',
    description: '専門知識を活かし、論理的かつ簡潔に回答します。',
  },
  {
    id: 'ai3',
    description: '友達感覚でフランクに回答します。',
  },
  {
    id: 'ai4',
    description: '元気いっぱいにモチベーションを高める回答をします。',
  },
]

const modelValue = defineModel<string>({ default: 'ai1' })

const cardColor = (active?: boolean) =>
  active ? 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' : 'white'
</script>

<template>
  <v-card class="character-selector pa-4">
    <v-card-title class="text-h6 d-flex align-center mb-2">
      <v-icon icon="mdi-sparkles" size="20" class="mr-2" />
      AIキャラクター
    </v-card-title>
    <v-divider class="mb-2" />

    <!-- アイテムグループで単一選択 -->
    <v-item-group v-model="modelValue" mandatory>
      <v-container class="pa-0" fluid>
        <v-row no-gutters>
          <v-col cols="12" v-for="char in characters" :key="char.id" class="mb-2">
            <v-item :value="char.id" #="{ isSelected, toggle }">
              <v-card
                :style="{ background: cardColor(isSelected) }"
                :elevation="isSelected ? 12 : 2"
                class="d-flex align-center pa-3 transition-all"
                :class="isSelected ? 'selected-card' : 'unselected-card'"
                style="cursor: pointer; display: flex; align-items: center"
                @click="toggle"
              >
                <div class="flex-shrink-0" style="width: 2.5rem; height: 2.5rem">
                  <AiIcon :id="char.id" style="width: 100%; height: 100%" />
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
}

.unselected-card {
  transform: scale(1);
  border: 2px solid transparent;
}

.transition-all {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.selected-card:hover {
  transform: scale(1.03);
}

.unselected-card:hover {
  transform: scale(1.01);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15) !important;
}
</style>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, nextTick } from 'vue'

interface AttributeValue {
  id: string
  name: string
  handle: string
}

const props = defineProps<{
  modelValue: string[]
  attributeValues?: AttributeValue[]
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', v: string[]): void
}>()

// ── Default colour palette ────────────────────────────────────────────────────
const DEFAULT_COLORS: { name: string; hex: string }[] = [
  { name: 'Beige',     hex: '#F5F5DC' },
  { name: 'Black',     hex: '#1C1C1C' },
  { name: 'Blue',      hex: '#2563EB' },
  { name: 'Bronze',    hex: '#CD7F32' },
  { name: 'Brown',     hex: '#92400E' },
  { name: 'Burgundy',  hex: '#800020' },
  { name: 'Coral',     hex: '#FF6B6B' },
  { name: 'Cream',     hex: '#FFFDD0' },
  { name: 'Cyan',      hex: '#06B6D4' },
  { name: 'Gold',      hex: '#F59E0B' },
  { name: 'Gray',      hex: '#6B7280' },
  { name: 'Green',     hex: '#16A34A' },
  { name: 'Indigo',    hex: '#4338CA' },
  { name: 'Ivory',     hex: '#FFFFF0' },
  { name: 'Khaki',     hex: '#C3B091' },
  { name: 'Lavender',  hex: '#E6E6FA' },
  { name: 'Maroon',    hex: '#800000' },
  { name: 'Mint',      hex: '#98D8C8' },
  { name: 'Navy',      hex: '#1E3A8A' },
  { name: 'Olive',     hex: '#808000' },
  { name: 'Orange',    hex: '#F97316' },
  { name: 'Pink',      hex: '#EC4899' },
  { name: 'Purple',    hex: '#9333EA' },
  { name: 'Red',       hex: '#DC2626' },
  { name: 'Rose',      hex: '#F43F5E' },
  { name: 'Silver',    hex: '#9CA3AF' },
  { name: 'Tan',       hex: '#D2B48C' },
  { name: 'Teal',      hex: '#0D9488' },
  { name: 'Turquoise', hex: '#40E0D0' },
  { name: 'Violet',    hex: '#7C3AED' },
  { name: 'White',     hex: '#F9FAFB' },
  { name: 'Yellow',    hex: '#FACC15' },
]

const hexMap = Object.fromEntries(DEFAULT_COLORS.map(c => [c.name.toLowerCase(), c.hex]))
const getHex = (name: string) => hexMap[name.toLowerCase()]

// ── State ─────────────────────────────────────────────────────────────────────
const isOpen = ref(false)
const search = ref('')
const newColorName = ref('')
const showNewInput = ref(false)
const containerRef = ref<HTMLElement | null>(null)
const searchInputRef = ref<HTMLInputElement | null>(null)

// ── Dropdown open/close ───────────────────────────────────────────────────────
const open = async () => {
  isOpen.value = true
  await nextTick()
  searchInputRef.value?.focus()
}

const close = () => {
  isOpen.value = false
  search.value = ''
  showNewInput.value = false
  newColorName.value = ''
}

const onClickOutside = (e: MouseEvent) => {
  if (containerRef.value && !containerRef.value.contains(e.target as Node)) close()
}

onMounted(() => document.addEventListener('mousedown', onClickOutside))
onBeforeUnmount(() => document.removeEventListener('mousedown', onClickOutside))

// ── Helpers ───────────────────────────────────────────────────────────────────
const isSelected = (name: string) => props.modelValue.includes(name)

const toggle = (name: string) => {
  if (isSelected(name)) {
    emit('update:modelValue', props.modelValue.filter(n => n !== name))
  } else {
    emit('update:modelValue', [...props.modelValue, name])
  }
}

const addCustomColor = () => {
  const name = newColorName.value.trim()
  if (!name) return
  if (!isSelected(name)) emit('update:modelValue', [...props.modelValue, name])
  newColorName.value = ''
  showNewInput.value = false
}

// ── Computed lists ────────────────────────────────────────────────────────────
const attrNames = computed(() => (props.attributeValues ?? []).map(v => v.name))

const filteredAttrValues = computed(() => {
  const q = search.value.toLowerCase()
  return (props.attributeValues ?? []).filter(v => v.name.toLowerCase().includes(q))
})

const filteredDefaults = computed(() => {
  const q = search.value.toLowerCase()
  const attrSet = new Set(attrNames.value.map(n => n.toLowerCase()))
  return DEFAULT_COLORS.filter(c =>
    !attrSet.has(c.name.toLowerCase()) && c.name.toLowerCase().includes(q),
  )
})

const showDefaultSection = computed(() => filteredDefaults.value.length > 0)
</script>

<template>
  <div ref="containerRef" class="relative">
    <!-- Trigger -->
    <div
      class="min-h-[38px] px-3 py-1.5 border rounded-lg bg-[#F6F6F7] flex flex-wrap gap-1.5 items-center cursor-text transition-colors"
      :class="isOpen ? 'border-[#1E212B]' : 'border-[#EDEDEE] hover:border-[#94BDD8]'"
      @click="open"
    >
      <!-- Selected colour chips -->
      <span
        v-for="name in modelValue"
        :key="name"
        class="flex items-center gap-1.5 px-2 py-0.5 bg-[#1E212B] text-white text-xs rounded-md"
      >
        <span
          class="w-3 h-3 rounded-full flex-shrink-0 border border-white/20"
          :style="{ background: getHex(name) || 'repeating-conic-gradient(#e5e7eb 0% 25%, #fff 0% 50%) 0 0 / 6px 6px' }"
        />
        {{ name }}
        <button type="button" class="opacity-70 hover:opacity-100" @click.stop="toggle(name)">
          <svg class="w-2.5 h-2.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
            <path d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </span>
      <span v-if="modelValue.length === 0" class="text-sm text-[#939499] pointer-events-none">Add colours…</span>
    </div>

    <!-- Dropdown -->
    <Transition
      enter-active-class="transition ease-out duration-100"
      enter-from-class="opacity-0 scale-95"
      enter-to-class="opacity-100 scale-100"
      leave-active-class="transition ease-in duration-75"
      leave-from-class="opacity-100 scale-100"
      leave-to-class="opacity-0 scale-95"
    >
      <div
        v-if="isOpen"
        class="absolute left-0 right-0 top-full mt-1 bg-white border border-[#EDEDEE] rounded-xl shadow-lg z-20 origin-top overflow-hidden"
      >
        <!-- Search -->
        <div class="px-3 pt-2.5 pb-2 border-b border-[#EDEDEE]">
          <input
            ref="searchInputRef"
            v-model="search"
            type="text"
            placeholder="Add color"
            class="w-full px-2 py-1.5 bg-[#F6F6F7] border border-[#EDEDEE] rounded-lg text-sm text-[#1E212B] placeholder-[#939499] outline-none focus:border-[#1E212B] transition-colors"
            @keydown.escape="close"
          >
        </div>

        <!-- List -->
        <div class="max-h-60 overflow-y-auto">
          <!-- Attribute values -->
          <template v-if="filteredAttrValues.length > 0">
            <button
              v-for="av in filteredAttrValues"
              :key="av.id"
              type="button"
              class="w-full flex items-center gap-2.5 px-3 py-2 hover:bg-[#F6F6F7] transition-colors text-left"
              @click="toggle(av.name)"
            >
              <span
                class="w-4 h-4 rounded border flex items-center justify-center flex-shrink-0 transition-colors"
                :class="isSelected(av.name) ? 'bg-[#1E212B] border-[#1E212B]' : 'border-[#DBDBDD]'"
              >
                <svg v-if="isSelected(av.name)" class="w-2.5 h-2.5 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                  <path d="M5 13l4 4L19 7" />
                </svg>
              </span>
              <span
                v-if="getHex(av.name)"
                class="w-5 h-5 rounded-md flex-shrink-0 border border-black/10"
                :style="{ background: getHex(av.name) }"
              />
              <span
                v-else
                class="w-5 h-5 rounded-md flex-shrink-0 border border-[#DBDBDD] overflow-hidden"
                style="background: repeating-conic-gradient(#e5e7eb 0% 25%, #fff 0% 50%) 0 0 / 8px 8px"
              />
              <span class="text-sm text-[#1E212B]">{{ av.name }}</span>
            </button>
          </template>

          <!-- Default entries -->
          <template v-if="showDefaultSection">
            <div v-if="filteredAttrValues.length > 0" class="px-3 py-1.5 text-xs font-medium text-[#939499]">
              Default entries
            </div>
            <button
              v-for="color in filteredDefaults"
              :key="color.name"
              type="button"
              class="w-full flex items-center gap-2.5 px-3 py-2 hover:bg-[#F6F6F7] transition-colors text-left"
              @click="toggle(color.name)"
            >
              <span
                class="w-4 h-4 rounded border flex items-center justify-center flex-shrink-0 transition-colors"
                :class="isSelected(color.name) ? 'bg-[#1E212B] border-[#1E212B]' : 'border-[#DBDBDD]'"
              >
                <svg v-if="isSelected(color.name)" class="w-2.5 h-2.5 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                  <path d="M5 13l4 4L19 7" />
                </svg>
              </span>
              <span
                class="w-5 h-5 rounded-md flex-shrink-0 border border-black/10"
                :style="{ background: color.hex }"
              />
              <span class="text-sm text-[#1E212B]">{{ color.name }}</span>
            </button>
          </template>

          <div v-if="filteredAttrValues.length === 0 && !showDefaultSection" class="px-3 py-3 text-sm text-[#939499] text-center">
            No colours found
          </div>
        </div>

        <!-- Add new entry -->
        <div class="border-t border-[#EDEDEE] px-3 py-2">
          <div v-if="showNewInput" class="flex items-center gap-2">
            <input
              v-model="newColorName"
              type="text"
              placeholder="Color name"
              class="flex-1 px-3 py-1.5 bg-[#F6F6F7] border border-[#EDEDEE] rounded-lg text-sm text-[#1E212B] placeholder-[#939499] outline-none focus:border-[#1E212B] transition-colors"
              @keydown.enter.prevent="addCustomColor"
              @keydown.escape="showNewInput = false"
            >
            <button type="button" class="text-xs font-medium text-[#2176AE] hover:underline" @click="addCustomColor">Add</button>
            <button type="button" class="text-xs text-[#939499] hover:text-[#1E212B]" @click="showNewInput = false">Cancel</button>
          </div>
          <button
            v-else
            type="button"
            class="flex items-center gap-1.5 text-sm text-[#1E212B] hover:text-[#2176AE] transition-colors"
            @click="showNewInput = true"
          >
            <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10" />
              <path d="M12 8v8M8 12h8" />
            </svg>
            Add new entry
          </button>
        </div>
      </div>
    </Transition>
  </div>
</template>

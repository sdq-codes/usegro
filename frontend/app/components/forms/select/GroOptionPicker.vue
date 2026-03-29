<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, nextTick } from 'vue'

const props = defineProps<{
  modelValue: string[]
  optionName: string
  attributeValues?: string[]
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', v: string[]): void
}>()

const isOpen = ref(false)
const search = ref('')
const newEntry = ref('')
const showNewInput = ref(false)
const containerRef = ref<HTMLElement | null>(null)
const searchInputRef = ref<HTMLInputElement | null>(null)

const open = async () => {
  isOpen.value = true
  await nextTick()
  searchInputRef.value?.focus()
}

const close = () => {
  isOpen.value = false
  search.value = ''
  showNewInput.value = false
  newEntry.value = ''
}

const onClickOutside = (e: MouseEvent) => {
  if (containerRef.value && !containerRef.value.contains(e.target as Node)) close()
}

onMounted(() => document.addEventListener('mousedown', onClickOutside))
onBeforeUnmount(() => document.removeEventListener('mousedown', onClickOutside))

const isSelected = (val: string) => props.modelValue.includes(val)

const toggle = (val: string) => {
  if (isSelected(val)) {
    emit('update:modelValue', props.modelValue.filter(v => v !== val))
  } else {
    emit('update:modelValue', [...props.modelValue, val])
  }
}

const addCustomEntry = () => {
  const val = newEntry.value.trim()
  if (!val) return
  if (!isSelected(val)) emit('update:modelValue', [...props.modelValue, val])
  newEntry.value = ''
  showNewInput.value = false
}

// All available options: attribute values + any selected custom entries not in attribute values
const allOptions = computed(() => {
  const attrSet = new Set((props.attributeValues ?? []).map(v => v.toLowerCase()))
  const customSelected = props.modelValue.filter(v => !attrSet.has(v.toLowerCase()))
  return [...(props.attributeValues ?? []), ...customSelected]
})

const filteredOptions = computed(() => {
  const q = search.value.toLowerCase()
  return allOptions.value.filter(v => v.toLowerCase().includes(q))
})

const placeholder = computed(() => `Add ${props.optionName.toLowerCase()}`)
</script>

<template>
  <div ref="containerRef" class="relative">
    <!-- Trigger -->
    <div
      class="min-h-[38px] px-3 py-1.5 border rounded-lg bg-[#F6F6F7] flex flex-wrap gap-1.5 items-center cursor-text transition-colors"
      :class="isOpen ? 'border-[#1E212B]' : 'border-[#EDEDEE] hover:border-[#94BDD8]'"
      @click="open"
    >
      <span
        v-for="val in modelValue"
        :key="val"
        class="flex items-center gap-1 px-2 py-0.5 bg-[#1E212B] text-white text-xs rounded-md"
      >
        {{ val }}
        <button type="button" class="opacity-70 hover:opacity-100 transition-opacity" @click.stop="toggle(val)">
          <svg class="w-2.5 h-2.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
            <path d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </span>
      <span v-if="modelValue.length === 0" class="text-sm text-[#939499] pointer-events-none">{{ placeholder }}</span>
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
            :placeholder="placeholder"
            class="w-full px-2 py-1.5 bg-[#F6F6F7] border border-[#EDEDEE] rounded-lg text-sm text-[#1E212B] placeholder-[#939499] outline-none focus:border-[#1E212B] transition-colors"
            @keydown.escape="close"
          >
        </div>

        <!-- List -->
        <div class="max-h-60 overflow-y-auto">
          <button
            v-for="val in filteredOptions"
            :key="val"
            type="button"
            class="w-full flex items-center gap-2.5 px-3 py-2 hover:bg-[#F6F6F7] transition-colors text-left"
            @click="toggle(val)"
          >
            <span
              class="w-4 h-4 rounded border flex items-center justify-center flex-shrink-0 transition-colors"
              :class="isSelected(val) ? 'bg-[#1E212B] border-[#1E212B]' : 'border-[#DBDBDD]'"
            >
              <svg v-if="isSelected(val)" class="w-2.5 h-2.5 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                <path d="M5 13l4 4L19 7" />
              </svg>
            </span>
            <span class="text-sm text-[#1E212B]">{{ val }}</span>
          </button>

          <div v-if="filteredOptions.length === 0" class="px-3 py-3 text-sm text-[#939499] text-center">
            No options found
          </div>
        </div>

        <!-- Add new entry -->
        <div class="border-t border-[#EDEDEE] px-3 py-2">
          <div v-if="showNewInput" class="flex items-center gap-2">
            <input
              v-model="newEntry"
              type="text"
              placeholder="Custom value"
              class="flex-1 px-3 py-1.5 bg-[#F6F6F7] border border-[#EDEDEE] rounded-lg text-sm text-[#1E212B] placeholder-[#939499] outline-none focus:border-[#1E212B] transition-colors"
              @keydown.enter.prevent="addCustomEntry"
              @keydown.escape="showNewInput = false"
            >
            <button type="button" class="text-xs font-medium text-[#2176AE] hover:underline" @click="addCustomEntry">Add</button>
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

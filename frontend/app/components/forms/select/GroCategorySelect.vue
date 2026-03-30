<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { useCategoryAPI } from '@/composables/api/catalog/category'

interface Category {
  id: string
  name: string
}

const model = defineModel<string[]>({ default: () => [] })

const categories = ref<Category[]>([])
const search = ref('')
const isOpen = ref(false)
const isCreating = ref(false)
const containerRef = ref<HTMLElement | null>(null)
const searchInputRef = ref<HTMLInputElement | null>(null)

const filtered = computed(() => {
  if (!search.value.trim()) return categories.value
  return categories.value.filter(c =>
    c.name.toLowerCase().includes(search.value.toLowerCase())
  )
})

const showCreate = computed(() =>
  search.value.trim() &&
  !categories.value.some(c => c.name.toLowerCase() === search.value.trim().toLowerCase())
)

const selectedCategories = computed(() =>
  categories.value.filter(c => model.value.includes(c.id))
)

const open = async () => {
  isOpen.value = true
  await nextTick()
  searchInputRef.value?.focus()
}

const toggle = (id: string) => {
  if (model.value.includes(id)) {
    model.value = model.value.filter(v => v !== id)
  } else {
    model.value = [...model.value, id]
    search.value = ''
  }
}

const createAndSelect = async () => {
  const name = search.value.trim()
  if (!name || isCreating.value) return
  isCreating.value = true
  const res = await useCategoryAPI().CreateCategory({ name })
  if (res.success && res.data?.data) {
    categories.value.push(res.data.data)
    model.value = [...model.value, res.data.data.id]
    search.value = ''
  }
  isCreating.value = false
}

const fetchCategories = async () => {
  const res = await useCategoryAPI().ListCategories()
  if (res.success && res.data?.data) categories.value = res.data.data
}

const onClickOutside = (e: MouseEvent) => {
  if (containerRef.value && !containerRef.value.contains(e.target as Node)) {
    isOpen.value = false
    search.value = ''
  }
}

onMounted(() => {
  fetchCategories()
  document.addEventListener('mousedown', onClickOutside)
})
onBeforeUnmount(() => document.removeEventListener('mousedown', onClickOutside))
</script>

<template>
  <div ref="containerRef" class="relative">
    <!-- Input area with selected chips -->
    <div
      class="min-h-[38px] px-3 py-1.5 border rounded-lg bg-[#F6F6F7] flex flex-wrap gap-1.5 items-center cursor-text transition-colors"
      :class="isOpen ? 'border-[#1E212B]' : 'border-[#EDEDEE] hover:border-[#94BDD8]'"
      @click="open"
    >
      <span
        v-for="cat in selectedCategories"
        :key="cat.id"
        class="flex items-center gap-1 px-2 py-0.5 bg-[#1E212B] text-white text-xs rounded-md"
      >
        {{ cat.name }}
        <button
          type="button"
          class="opacity-70 hover:opacity-100 transition-opacity"
          @click.stop="toggle(cat.id)"
        >
          <svg class="w-2.5 h-2.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
            <path d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </span>
      <input
        ref="searchInputRef"
        v-model="search"
        class="flex-1 min-w-[80px] bg-transparent text-sm outline-none text-[#1E212B] placeholder-[#939499]"
        placeholder="Search categories..."
        @focus="isOpen = true"
        @keydown.escape="isOpen = false"
      >
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
        class="absolute left-0 right-0 top-full mt-1 bg-white border border-[#EDEDEE] rounded-xl shadow-lg z-20 max-h-52 overflow-y-auto origin-top"
      >
        <button
          v-for="cat in filtered"
          :key="cat.id"
          type="button"
          class="w-full flex items-center gap-2.5 px-3 py-2 text-sm hover:bg-[#F6F6F7] transition-colors text-left"
          @click="toggle(cat.id)"
        >
          <span
            class="w-4 h-4 rounded border flex items-center justify-center flex-shrink-0 transition-colors"
            :class="model.includes(cat.id) ? 'bg-[#1E212B] border-[#1E212B]' : 'border-[#DBDBDD]'"
          >
            <svg v-if="model.includes(cat.id)" class="w-2.5 h-2.5 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
              <path d="M5 13l4 4L19 7" />
            </svg>
          </span>
          {{ cat.name }}
        </button>

        <button
          v-if="showCreate"
          type="button"
          class="w-full flex items-center gap-2.5 px-3 py-2 text-sm text-[#2176AE] hover:bg-[#EFF6FF] transition-colors"
          :disabled="isCreating"
          @click="createAndSelect"
        >
          <svg class="w-4 h-4 flex-shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <path d="M12 5v14M5 12h14" />
          </svg>
          Create "{{ search.trim() }}"
        </button>

        <div
          v-if="filtered.length === 0 && !showCreate"
          class="px-3 py-3 text-sm text-[#939499] text-center"
        >
          No categories found
        </div>
      </div>
    </Transition>
  </div>
</template>

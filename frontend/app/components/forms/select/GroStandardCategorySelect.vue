<script setup lang="ts">
import { ref, computed, watch, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { useStandardCategoryAPI } from '@/composables/api/catalog/category'
import type { StandardCategory, StandardCategorySearchResult, StandardAttribute } from '@/composables/dto/catalog/category'

interface BreadcrumbItem {
  id: string
  name: string
}

export interface StandardCategorySelection {
  id: string
  name: string
  attributes: StandardAttribute[]
}

// Selected leaf with its attributes — or null if nothing selected
const model = defineModel<StandardCategorySelection | null>({ default: null })

const api = useStandardCategoryAPI()

const isOpen = ref(false)
const containerRef = ref<HTMLElement | null>(null)
const searchInputRef = ref<HTMLInputElement | null>(null)

// Browse state
const levels = ref<StandardCategory[][]>([])
const breadcrumb = ref<BreadcrumbItem[]>([])
const loading = ref(false)

// Search state
const searchQuery = ref('')
const searchResults = ref<StandardCategorySearchResult[]>([])
const searchLoading = ref(false)
const searchDebounceTimer = ref<ReturnType<typeof setTimeout> | null>(null)

const isSearching = computed(() => searchQuery.value.trim().length > 0)

// ------------------------------------------------------------------
// Data loading
// ------------------------------------------------------------------
const loadRoots = async () => {
  loading.value = true
  const res = await api.ListRootCategories()
  if (res.success && res.data?.data) {
    levels.value = [res.data.data]
    breadcrumb.value = []
  }
  loading.value = false
}

const selectLeaf = async (cat: StandardCategory) => {
  // Close immediately so the UI feels responsive, then fetch attributes in the background
  model.value = { id: cat.id, name: fullLabel(cat), attributes: [] }
  close()
  const full = await api.GetCategory(cat.id)
  model.value = {
    id: cat.id,
    name: fullLabel(cat),
    attributes: full.data?.data?.attributes ?? [],
  }
}

const drillInto = async (cat: StandardCategory) => {
  if (cat.is_leaf) {
    await selectLeaf(cat)
    return
  }

  loading.value = true
  const res = await api.ListChildren(cat.id)
  loading.value = false

  if (!res.success) return

  breadcrumb.value = [...breadcrumb.value, { id: cat.id, name: cat.name }]
  levels.value = [...levels.value, res.data?.data ?? []]
}

const navigateTo = (index: number) => {
  levels.value = levels.value.slice(0, index + 2)
  breadcrumb.value = breadcrumb.value.slice(0, index + 1)
}

const selectSearchResult = async (result: StandardCategorySearchResult) => {
  model.value = { id: result.id, name: result.path, attributes: [] }
  close()
  const full = await api.GetCategory(result.id)
  model.value = { id: result.id, name: result.path, attributes: full.data?.data?.attributes ?? [] }
}

// ------------------------------------------------------------------
// Search
// ------------------------------------------------------------------
watch(searchQuery, (val) => {
  if (searchDebounceTimer.value) clearTimeout(searchDebounceTimer.value)
  if (!val.trim()) {
    searchResults.value = []
    return
  }
  searchDebounceTimer.value = setTimeout(async () => {
    searchLoading.value = true
    const res = await api.SearchCategories(val.trim())
    if (res.success && res.data?.data) {
      searchResults.value = res.data.data
    }
    searchLoading.value = false
  }, 300)
})

// ------------------------------------------------------------------
// Helpers
// ------------------------------------------------------------------
const currentLevel = computed(() => levels.value[levels.value.length - 1] ?? [])

const displayLabel = computed(() => {
  if (!model.value) return null
  const parts = model.value.name.split(' › ')
  return parts[parts.length - 1]
})

const fullLabel = (cat: StandardCategory) => {
  const parts = breadcrumb.value.map(b => b.name)
  parts.push(cat.name)
  return parts.join(' › ')
}

const clear = () => {
  model.value = null
}

// ------------------------------------------------------------------
// Open / close
// ------------------------------------------------------------------
const open = async () => {
  if (isOpen.value) return
  isOpen.value = true
  if (levels.value.length === 0) await loadRoots()
  await nextTick()
  searchInputRef.value?.focus()
}

const close = () => {
  isOpen.value = false
  searchQuery.value = ''
  searchResults.value = []
}

const onClickOutside = (e: MouseEvent) => {
  if (containerRef.value && !containerRef.value.contains(e.target as Node)) {
    close()
  }
}

onMounted(() => document.addEventListener('mousedown', onClickOutside))
onBeforeUnmount(() => {
  document.removeEventListener('mousedown', onClickOutside)
  if (searchDebounceTimer.value) clearTimeout(searchDebounceTimer.value)
})
</script>

<template>
  <div ref="containerRef" class="relative">
    <!-- Label row -->
    <div v-if="$slots.default" class="flex items-center space-x-4 mb-1">
      <slot />
      <div v-if="model && model.attributes.length > 0" class="ml-auto flex items-center gap-1.5">
        <span class="text-xs bg-[#EDEDEE] px-2 py-1 rounded-xl text-[#070707]">{{ model.attributes.length }}</span>
        <span class="text-xs text-[#939499]">Recommended variants</span>
      </div>
    </div>

    <!-- Trigger -->
    <div
      class="min-h-[38px] px-3 py-1.5 border rounded-lg bg-[#F6F6F7] flex items-center justify-between gap-2 cursor-pointer transition-colors select-none"
      :class="isOpen ? 'border-[#1E212B]' : 'border-[#EDEDEE] hover:border-[#94BDD8]'"
      @click="open"
    >
      <span v-if="model" class="text-sm text-[#1E212B] truncate">{{ displayLabel }}</span>
      <span v-else class="text-sm text-[#939499]">Select a category...</span>

      <div class="flex items-center gap-1 flex-shrink-0">
        <button
          v-if="model"
          type="button"
          class="text-[#939499] hover:text-[#1E212B] transition-colors"
          @click.stop="clear"
        >
          <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
            <path d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
        <svg
          class="w-4 h-4 text-[#939499] transition-transform"
          :class="isOpen ? 'rotate-180' : ''"
          viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
        >
          <path d="M6 9l6 6 6-6" />
        </svg>
      </div>
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
        class="absolute left-0 right-0 top-full mt-1 bg-white border border-[#EDEDEE] rounded-xl shadow-lg z-20 origin-top"
      >
        <!-- Search input -->
        <div class="px-3 pt-2.5 pb-2 border-b border-[#EDEDEE]">
          <div class="flex items-center gap-2 px-2.5 py-1.5 bg-[#F6F6F7] rounded-lg">
            <svg class="w-3.5 h-3.5 text-[#939499] flex-shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
              <circle cx="11" cy="11" r="8" />
              <path d="M21 21l-4.35-4.35" />
            </svg>
            <input
              ref="searchInputRef"
              v-model="searchQuery"
              type="text"
              placeholder="Search categories..."
              class="flex-1 bg-transparent text-sm text-[#1E212B] placeholder-[#939499] outline-none"
              @click.stop
            />
            <button
              v-if="searchQuery"
              type="button"
              class="text-[#939499] hover:text-[#1E212B] transition-colors"
              @click.stop="searchQuery = ''"
            >
              <svg class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                <path d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>

        <!-- Search results -->
        <div v-if="isSearching" class="max-h-52 overflow-y-auto">
          <div v-if="searchLoading" class="px-3 py-4 text-sm text-[#939499] text-center">
            Searching...
          </div>
          <template v-else-if="searchResults.length > 0">
            <button
              v-for="result in searchResults"
              :key="result.id"
              type="button"
              class="w-full flex flex-col px-3 py-2 text-left hover:bg-[#F6F6F7] transition-colors"
              :class="model?.id === result.id ? 'bg-[#F6F6F7]' : ''"
              @click="selectSearchResult(result)"
            >
              <span class="text-sm text-[#1E212B] truncate">{{ result.name }}</span>
              <span v-if="result.path !== result.name" class="text-xs text-[#939499] truncate mt-0.5">{{ result.path }}</span>
            </button>
          </template>
          <div v-else class="px-3 py-4 text-sm text-[#939499] text-center">
            No categories found
          </div>
        </div>

        <!-- Browse mode -->
        <template v-else>
          <!-- Breadcrumb -->
          <div v-if="breadcrumb.length > 0" class="flex items-center gap-1 px-3 py-2 border-b border-[#EDEDEE] overflow-x-auto">
            <button
              type="button"
              class="text-xs text-[#2176AE] hover:underline flex-shrink-0"
              @click="navigateTo(-1)"
            >
              All
            </button>
            <template v-for="(crumb, i) in breadcrumb" :key="crumb.id">
              <svg class="w-3 h-3 text-[#939499] flex-shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M9 6l6 6-6 6" />
              </svg>
              <button
                type="button"
                class="text-xs flex-shrink-0 transition-colors"
                :class="i === breadcrumb.length - 1 ? 'text-[#1E212B] font-medium' : 'text-[#2176AE] hover:underline'"
                @click="navigateTo(i - 1)"
              >
                {{ crumb.name }}
              </button>
            </template>
          </div>

          <!-- List -->
          <div class="max-h-52 overflow-y-auto">
            <div v-if="loading" class="px-3 py-4 text-sm text-[#939499] text-center">
              Loading...
            </div>

            <button
              v-for="cat in currentLevel"
              v-else
              :key="cat.id"
              type="button"
              class="w-full flex items-center justify-between px-3 py-2 text-sm hover:bg-[#F6F6F7] transition-colors text-left"
              :class="model?.id === cat.id ? 'bg-[#F6F6F7] font-medium' : ''"
              @click="drillInto(cat)"
            >
              <span class="truncate">{{ cat.name }}</span>
              <svg v-if="!cat.is_leaf" class="w-3.5 h-3.5 text-[#939499] flex-shrink-0 ml-2" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M9 6l6 6-6 6" />
              </svg>
            </button>

            <div v-if="!loading && currentLevel.length === 0" class="px-3 py-3 text-sm text-[#939499] text-center">
              No categories found
            </div>
          </div>
        </template>
      </div>
    </Transition>
  </div>
</template>

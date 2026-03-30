<script setup lang="ts">
import { ref, computed, watch, onUnmounted } from 'vue'
import { FilterIcon, Cancel01Icon } from '@hugeicons/core-free-icons'
import { HugeiconsIcon } from '@hugeicons/vue'

interface FilterOption {
  label: string
  value: string
}

export interface FilterDefinition {
  key: string
  label: string
  icon?: object
  priority?: boolean
  /** default: 'options' */
  type?: 'options' | 'date-range' | 'text' | 'number-range'
  options?: FilterOption[]
  /** For date-range / number-range: the row key whose raw value to compare against.
   *  Falls back to filter.key when omitted. */
  rawKey?: string
}

const props = defineProps<{
  filters: FilterDefinition[]
}>()

// model: Record<filter.key, string>
//   options      → '' | selected value
//   text         → '' | search string
//   date-range   → '' | JSON { preset, from, to }
//   number-range → '' | JSON { min, max }
const model = defineModel<Record<string, string>>({ default: () => ({}) })

const showPanel = ref(false)
const selectedCategory = ref(props.filters[0]?.key ?? '')
const triggerRef = ref<HTMLElement | null>(null)
const panelRef = ref<HTMLElement | null>(null)
const panelStyle = ref<Record<string, string>>({ top: '0px', left: '0px' })

// ── Date-range local state ────────────────────────────────────
const DATE_PRESETS = [
  { label: 'Custom', value: 'custom' },
  { label: 'Today', value: 'today' },
  { label: 'Yesterday', value: 'yesterday' },
  { label: 'This week', value: 'this_week' },
  { label: 'This month', value: 'this_month' },
  { label: 'Last month', value: 'last_month' },
]

const toIso = (d: Date) => d.toISOString().slice(0, 10)

const resolveDatePreset = (preset: string): { from: string; to: string } => {
  const today = new Date()
  if (preset === 'today') {
    const s = toIso(today)
    return { from: s, to: s }
  }
  if (preset === 'yesterday') {
    const d = new Date(today); d.setDate(d.getDate() - 1)
    const s = toIso(d); return { from: s, to: s }
  }
  if (preset === 'this_week') {
    const d = new Date(today); d.setDate(d.getDate() - d.getDay())
    return { from: toIso(d), to: toIso(today) }
  }
  if (preset === 'this_month') {
    const d = new Date(today.getFullYear(), today.getMonth(), 1)
    return { from: toIso(d), to: toIso(today) }
  }
  if (preset === 'last_month') {
    const first = new Date(today.getFullYear(), today.getMonth() - 1, 1)
    const last = new Date(today.getFullYear(), today.getMonth(), 0)
    return { from: toIso(first), to: toIso(last) }
  }
  return { from: '', to: '' }
}

// Local state per date-range / number-range filter
const parseDateRange = (key: string) => {
  const raw = model.value[key]
  if (!raw) return { preset: 'custom', from: '', to: '' }
  try { return JSON.parse(raw) } catch { return { preset: 'custom', from: '', to: '' } }
}

const parseNumberRange = (key: string) => {
  const raw = model.value[key]
  if (!raw) return { min: '', max: '' }
  try { return JSON.parse(raw) } catch { return { min: '', max: '' } }
}

const setDatePreset = (key: string, preset: string) => {
  if (preset === 'custom') {
    const cur = parseDateRange(key)
    model.value = { ...model.value, [key]: JSON.stringify({ preset: 'custom', from: cur.from, to: cur.to }) }
  } else {
    const { from, to } = resolveDatePreset(preset)
    model.value = { ...model.value, [key]: JSON.stringify({ preset, from, to }) }
  }
}

const setDateFrom = (key: string, from: string) => {
  const cur = parseDateRange(key)
  const val = JSON.stringify({ preset: 'custom', from, to: cur.to })
  model.value = { ...model.value, [key]: val }
}

const setDateTo = (key: string, to: string) => {
  const cur = parseDateRange(key)
  const val = JSON.stringify({ preset: 'custom', from: cur.from, to })
  model.value = { ...model.value, [key]: val }
}

const setNumberRange = (key: string, field: 'min' | 'max', value: string) => {
  const cur = parseNumberRange(key)
  cur[field] = value
  const serialized = cur.min === '' && cur.max === '' ? '' : JSON.stringify(cur)
  model.value = { ...model.value, [key]: serialized }
}

// ── Panel positioning ─────────────────────────────────────────
const priorityFilters = computed(() => props.filters.filter(f => f.priority).slice(0, 2))
const priorityKeys = computed(() => new Set(priorityFilters.value.map(f => f.key)))

const activeChips = computed(() =>
  props.filters
    .filter(f => priorityKeys.value.has(f.key) && model.value[f.key])
    .map(f => {
      const type = f.type ?? 'options'
      let label = ''
      if (type === 'options') {
        label = f.options?.find(o => o.value === model.value[f.key])?.label ?? model.value[f.key]
      } else if (type === 'text') {
        label = `"${model.value[f.key]}"`
      } else if (type === 'date-range') {
        try {
          const { preset, from, to } = JSON.parse(model.value[f.key])
          label = preset !== 'custom'
            ? DATE_PRESETS.find(p => p.value === preset)?.label ?? preset
            : [from, to].filter(Boolean).join(' → ') || f.label
        } catch { label = f.label }
      } else if (type === 'number-range') {
        try {
          const { min, max } = JSON.parse(model.value[f.key])
          label = [min && `$${min}`, max && `$${max}`].filter(Boolean).join('–')
        } catch { label = f.label }
      }
      return { key: f.key, label }
    })
)

const extraActiveCount = computed(() =>
  props.filters.filter(f => !priorityKeys.value.has(f.key) && model.value[f.key]).length
)

const activeCount = computed(() =>
  Object.values(model.value).filter(v => v !== '').length
)

const updatePanelPosition = () => {
  if (triggerRef.value) {
    const rect = triggerRef.value.getBoundingClientRect()
    const panelWidth = 560
    const overflowsRight = rect.left + panelWidth > window.innerWidth
    panelStyle.value = overflowsRight
      ? { top: `${rect.bottom + 6}px`, right: `${window.innerWidth - rect.right}px`, left: 'auto' }
      : { top: `${rect.bottom + 6}px`, left: `${rect.left}px`, right: 'auto' }
  }
}

const togglePanel = () => {
  if (showPanel.value) {
    showPanel.value = false
  } else {
    updatePanelPosition()
    showPanel.value = true
  }
}

const openCategory = (key: string) => {
  selectedCategory.value = key
  updatePanelPosition()
  showPanel.value = true
}

const clearAll = () => {
  const cleared: Record<string, string> = {}
  props.filters.forEach(f => { cleared[f.key] = '' })
  model.value = cleared
}

const clearOne = (key: string) => {
  model.value = { ...model.value, [key]: '' }
}

const setFilter = (key: string, value: string) => {
  model.value = { ...model.value, [key]: value }
}

const handleOutsideClick = (e: MouseEvent) => {
  const target = e.target as Node
  if (
    triggerRef.value && !triggerRef.value.contains(target) &&
    panelRef.value && !panelRef.value.contains(target)
  ) {
    showPanel.value = false
  }
}

watch(showPanel, (val) => {
  if (val) document.addEventListener('mousedown', handleOutsideClick)
  else document.removeEventListener('mousedown', handleOutsideClick)
})

onUnmounted(() => document.removeEventListener('mousedown', handleOutsideClick))

const selectedFilter = computed(() => props.filters.find(f => f.key === selectedCategory.value))
</script>

<template>
  <div ref="triggerRef" class="flex items-center gap-2 flex-wrap">
    <!-- Filters button -->
    <button
      class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm font-medium text-[#1E212B] hover:bg-[#F6F6F7] transition-colors shadow-[inset_0px_1px_0px_0px_#EBEBEB,inset_-1px_0px_0px_0px_#EBEBEB,inset_1px_0px_0px_0px_#EBEBEB,inset_0px_-1px_0px_0px_#CCCCCC]"
      :class="showPanel ? 'bg-[#EDEDEE]' : 'bg-white'"
      @click="togglePanel"
    >
      <HugeiconsIcon :icon="FilterIcon" class="w-4 h-4" color="#1E212B" />
      Filters
    </button>

    <!-- Priority quick-filter buttons -->
    <button
      v-for="filter in priorityFilters"
      :key="filter.key"
      class="flex items-center gap-1.5 px-3 py-1.5 border rounded-lg text-sm font-medium transition-colors"
      :class="model[filter.key]
        ? 'bg-[#1E212B] text-white border-[#1E212B]'
        : 'border-[#DBDBDD] text-[#1E212B] hover:bg-[#F6F6F7] shadow-[inset_0_1px_0_0_#E3E3E3]'"
      @click="openCategory(filter.key)"
    >
      {{ model[filter.key] || filter.label }}
      <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
        <path d="M6 9l6 6 6-6" />
      </svg>
    </button>

    <!-- Active filter chips -->
    <template v-if="activeChips.length > 0 || extraActiveCount > 0">
      <div class="w-px h-5 bg-[#DBDBDD]" />
      <span
        v-for="chip in activeChips"
        :key="chip.key"
        class="flex items-center gap-1 px-2.5 py-2 rounded-lg text-xs font-medium bg-white text-[#1E212B] shadow-[inset_0px_1px_0px_0px_#EBEBEB,inset_-1px_0px_0px_0px_#EBEBEB,inset_1px_0px_0px_0px_#EBEBEB,inset_0px_-1px_0px_0px_#CCCCCC]"
      >
        {{ chip.label }}
        <button class="ml-0.5 text-[#6F7177] hover:text-[#1E212B] transition-colors cursor-pointer" @click.stop="clearOne(chip.key)">
          <HugeiconsIcon :icon="Cancel01Icon" :size="11" />
        </button>
      </span>
      <button
        v-if="extraActiveCount > 0"
        class="flex items-center px-2.5 py-1.5 rounded-lg text-xs font-medium bg-[#1E212B] text-white cursor-pointer hover:bg-[#2d3142] transition-colors"
        @click="togglePanel"
      >
        +{{ extraActiveCount }}
      </button>
    </template>
  </div>

  <!-- Floating panel -->
  <Teleport to="body">
    <div
      v-if="showPanel"
      ref="panelRef"
      :style="panelStyle"
      class="fixed z-50 w-[560px] bg-white border border-[#EDEDEE] rounded-xl shadow-xl overflow-hidden"
    >
      <div class="flex">
        <!-- Left: category list -->
        <div class="px-4 w-56 border-r border-[#EDEDEE] shrink-0 bg-[#FCFCFC] py-3">
          <button
            v-for="filter in filters"
            :key="filter.key"
            class="w-full flex items-center cursor-pointer justify-between px-4 py-2.5 text-sm transition-colors"
            :class="selectedCategory === filter.key
              ? 'bg-[#EDEDEE] text-[#1E212B] font-medium rounded-lg'
              : 'text-[#4B4D55] hover:bg-[#EDEDEE] rounded-lg'"
            @click="selectedCategory = filter.key"
          >
            <span class="flex items-center gap-2">
              <HugeiconsIcon
                v-if="filter.icon"
                :icon="filter.icon"
                :size="15"
                :stroke-width="selectedCategory === filter.key ? 2 : 3"
                :color="selectedCategory === filter.key ? '#1E212B' : '#4B4D55'"
              />
              <span>{{ filter.label }}</span>
            </span>
            <span class="flex items-center gap-1.5">
              <span
                v-if="model[filter.key]"
                class="w-1.5 h-1.5 rounded-full bg-[#E87117]"
              />
              <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M9 18l6-6-6-6" />
              </svg>
            </span>
          </button>
        </div>

        <!-- Right: filter content -->
        <div class="flex-1 p-4 min-h-[240px]">
          <div class="flex items-center justify-between mb-4">
            <p class="text-xs font-semibold text-[#6F7177] uppercase tracking-wide">
              {{ selectedFilter?.label }}
            </p>
            <div class="flex items-center gap-3">
              <button
                v-if="activeCount > 0"
                class="text-xs font-medium text-[#D26B06] hover:text-[#b35a05] transition-colors cursor-pointer"
                @click="clearAll"
              >
                Clear all
              </button>
              <button class="text-[#6F7177] cursor-pointer hover:text-[#1E212B] transition-colors" @click="showPanel = false">
                <HugeiconsIcon :icon="Cancel01Icon" :size="15" />
              </button>
            </div>
          </div>

          <!-- OPTIONS -->
          <div v-if="!selectedFilter?.type || selectedFilter.type === 'options'" class="flex gap-2 flex-wrap">
            <button
              v-for="opt in selectedFilter?.options ?? []"
              :key="opt.value"
              class="px-3 py-1.5 cursor-pointer rounded-full text-xs font-medium border transition-colors"
              :class="model[selectedCategory] === opt.value
                ? 'bg-[#1E212B] text-white border-[#1E212B]'
                : 'bg-white text-[#4B4D55] border-gray-200 hover:border-gray-400'"
              @click="setFilter(selectedCategory, opt.value)"
            >
              {{ opt.label }}
            </button>
          </div>

          <!-- TEXT -->
          <div v-else-if="selectedFilter.type === 'text'">
            <input
              type="text"
              :value="model[selectedCategory]"
              :placeholder="`Search by ${selectedFilter.label.toLowerCase()}…`"
              class="w-full px-3 py-2 border border-[#DBDBDD] rounded-lg text-sm outline-none focus:border-[#1E212B] placeholder-[#9CA3AF]"
              @input="setFilter(selectedCategory, ($event.target as HTMLInputElement).value)"
            >
            <button
              v-if="model[selectedCategory]"
              class="mt-2 text-xs text-[#6F7177] hover:text-[#1E212B] cursor-pointer"
              @click="setFilter(selectedCategory, '')"
            >
              Clear
            </button>
          </div>

          <!-- DATE RANGE -->
          <div v-else-if="selectedFilter.type === 'date-range'" class="space-y-4">
            <div>
              <p class="text-xs text-[#6F7177] mb-1.5">Date range</p>
              <div class="relative">
                <select
                  :value="parseDateRange(selectedCategory).preset"
                  class="w-full px-3 py-2 border border-[#DBDBDD] rounded-lg text-sm outline-none focus:border-[#1E212B] bg-white appearance-none pr-8 cursor-pointer"
                  @change="setDatePreset(selectedCategory, ($event.target as HTMLSelectElement).value)"
                >
                  <option v-for="p in DATE_PRESETS" :key="p.value" :value="p.value">{{ p.label }}</option>
                </select>
                <svg class="pointer-events-none absolute right-2.5 top-1/2 -translate-y-1/2 w-4 h-4 text-[#6F7177]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path d="M6 9l6 6 6-6" />
                </svg>
              </div>
            </div>
            <div class="grid grid-cols-2 gap-3">
              <div>
                <p class="text-xs text-[#6F7177] mb-1.5">From</p>
                <input
                  type="date"
                  :value="parseDateRange(selectedCategory).from"
                  class="w-full px-3 py-2 border border-[#DBDBDD] rounded-lg text-sm outline-none focus:border-[#1E212B]"
                  @change="setDateFrom(selectedCategory, ($event.target as HTMLInputElement).value)"
                >
              </div>
              <div>
                <p class="text-xs text-[#6F7177] mb-1.5">To</p>
                <input
                  type="date"
                  :value="parseDateRange(selectedCategory).to"
                  class="w-full px-3 py-2 border border-[#DBDBDD] rounded-lg text-sm outline-none focus:border-[#1E212B]"
                  @change="setDateTo(selectedCategory, ($event.target as HTMLInputElement).value)"
                >
              </div>
            </div>
            <button
              v-if="model[selectedCategory]"
              class="text-xs text-[#6F7177] hover:text-[#1E212B] cursor-pointer"
              @click="setFilter(selectedCategory, '')"
            >
              Clear
            </button>
          </div>

          <!-- NUMBER RANGE -->
          <div v-else-if="selectedFilter.type === 'number-range'" class="space-y-3">
            <div class="grid grid-cols-2 gap-3">
              <div>
                <p class="text-xs text-[#6F7177] mb-1.5">Min amount</p>
                <div class="relative">
                  <span class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-[#6F7177]">$</span>
                  <input
                    type="number"
                    min="0"
                    :value="parseNumberRange(selectedCategory).min"
                    placeholder="0"
                    class="w-full pl-6 pr-3 py-2 border border-[#DBDBDD] rounded-lg text-sm outline-none focus:border-[#1E212B]"
                    @input="setNumberRange(selectedCategory, 'min', ($event.target as HTMLInputElement).value)"
                  >
                </div>
              </div>
              <div>
                <p class="text-xs text-[#6F7177] mb-1.5">Max amount</p>
                <div class="relative">
                  <span class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-[#6F7177]">$</span>
                  <input
                    type="number"
                    min="0"
                    :value="parseNumberRange(selectedCategory).max"
                    placeholder="Any"
                    class="w-full pl-6 pr-3 py-2 border border-[#DBDBDD] rounded-lg text-sm outline-none focus:border-[#1E212B]"
                    @input="setNumberRange(selectedCategory, 'max', ($event.target as HTMLInputElement).value)"
                  >
                </div>
              </div>
            </div>
            <button
              v-if="model[selectedCategory]"
              class="text-xs text-[#6F7177] hover:text-[#1E212B] cursor-pointer"
              @click="setFilter(selectedCategory, '')"
            >
              Clear
            </button>
          </div>
        </div>
      </div>
    </div>
  </Teleport>
</template>

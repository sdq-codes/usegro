<script setup lang="ts" generic="T extends Record<string, unknown>">
import { ref, computed } from 'vue'
import { Search01Icon } from '@hugeicons/core-free-icons'
import { HugeiconsIcon } from '@hugeicons/vue'
import GroFilterDropdown from '@/components/filters/GroFilterDropdown.vue'
import type { FilterDefinition } from '@/components/filters/GroFilterDropdown.vue'

export interface TableColumn {
  key: string
  label: string
  width?: string  // e.g. 'w-12', 'w-48'
}

export interface TableAction {
  key: string
  label: string
  class?: string  // e.g. 'text-red-600' for destructive actions
}

// Re-export so consumers only need one import
export type TableFilter = FilterDefinition

const props = defineProps<{
  columns: TableColumn[]
  rows: T[]
  rowKey: string
  actions?: TableAction[]
  filters?: TableFilter[]
  searchPlaceholder?: string
  page?: number
  totalPages?: number
  total?: number
  isLoading?: boolean
  emptyTitle?: string
  emptyMessage?: string
}>()

const emit = defineEmits<{
  (e: 'row-click', row: T): void
  (e: 'action', key: string, row: T): void
  (e: 'page-change', page: number): void
}>()

// Internal filter/search state
const filterValues = ref<Record<string, string>>(
  Object.fromEntries((props.filters ?? []).map(f => [f.key, '']))
)
const query = ref('')

const hasToolbar = computed(() => (props.filters && props.filters.length > 0) || props.searchPlaceholder !== undefined)

const filteredRows = computed<T[]>(() => {
  let list = props.rows

  if (props.filters) {
    for (const filter of props.filters) {
      const val = filterValues.value[filter.key]
      if (!val) continue

      const rowKey = filter.rawKey ?? filter.key
      const type = filter.type ?? 'options'

      if (type === 'options') {
        list = list.filter(row => row[rowKey] === val)
      } else if (type === 'text') {
        const term = val.toLowerCase()
        list = list.filter(row => String(row[rowKey] ?? '').toLowerCase().includes(term))
      } else if (type === 'date-range') {
        try {
          const { from, to } = JSON.parse(val)
          list = list.filter(row => {
            const raw = row[rowKey] as string | null
            if (!raw) return false
            const d = new Date(raw).getTime()
            if (from && d < new Date(from).getTime()) return false
            if (to && d > new Date(to + 'T23:59:59').getTime()) return false
            return true
          })
        } catch { /* invalid JSON, skip */ }
      } else if (type === 'number-range') {
        try {
          const { min, max } = JSON.parse(val)
          list = list.filter(row => {
            const num = Number(row[rowKey] ?? 0)
            if (min !== '' && !isNaN(Number(min)) && num < Number(min)) return false
            if (max !== '' && !isNaN(Number(max)) && num > Number(max)) return false
            return true
          })
        } catch { /* invalid JSON, skip */ }
      }
    }
  }

  if (query.value) {
    const term = query.value.toLowerCase()
    list = list.filter(row =>
      Object.values(row).some(v => {
        if (Array.isArray(v)) return v.some(x => String(x).toLowerCase().includes(term))
        if (typeof v === 'object' && v !== null) return JSON.stringify(v).toLowerCase().includes(term)
        return String(v).toLowerCase().includes(term)
      })
    )
  }

  return list
})

const openMenuRow = ref<unknown>(null)

const toggleMenu = (rowId: unknown) => {
  openMenuRow.value = openMenuRow.value === rowId ? null : rowId
}

const handleRowClick = (row: T) => {
  openMenuRow.value = null
  emit('row-click', row)
}

const handleAction = (actionKey: string, row: T) => {
  openMenuRow.value = null
  emit('action', actionKey, row)
}

const hasActions = computed(() => props.actions && props.actions.length > 0)
const currentPage = computed(() => props.page ?? 1)
const pages = computed(() => props.totalPages ?? 1)

const goToPage = (p: number) => {
  if (p < 1 || p > pages.value || p === currentPage.value) return
  emit('page-change', p)
}

const pageNumbers = computed<(number | '...')[]>(() => {
  const total = pages.value
  const current = currentPage.value
  if (total <= 7) return Array.from({ length: total }, (_, i) => i + 1)

  const list: (number | '...')[] = [1]
  if (current > 3) list.push('...')
  const start = Math.max(2, current - 1)
  const end = Math.min(total - 1, current + 1)
  for (let i = start; i <= end; i++) list.push(i)
  if (current < total - 2) list.push('...')
  list.push(total)
  return list
})
</script>

<template>
  <div>
    <!-- Toolbar -->
    <div v-if="hasToolbar" class="flex mb-4">
      <div class="flex items-center gap-2">
        <GroFilterDropdown
          v-if="filters && filters.length > 0"
          v-model="filterValues"
          :filters="filters"
        />
      </div>
      <div class="ml-auto relative">
        <HugeiconsIcon
          :icon="Search01Icon"
          class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4"
          color="#6F7177"
        />
        <input
          v-model="query"
          type="text"
          :placeholder="searchPlaceholder ?? 'Search'"
          class="pl-9 pr-4 py-2 bg-[#EDEDEE] border border-[#EDEDEE] rounded-3xl text-sm outline-none focus:border-[#1E212B] w-64 placeholder-[#939499] text-[#1E212B]"
        >
      </div>
    </div>

  <div class="bg-white rounded-2xl border border-[#EDEDEE] overflow-hidden">
    <table class="w-full" @click="openMenuRow = null">
      <thead>
        <tr class="bg-[#F6F6F7] border-b border-[#EDEDEE]">
          <th
            v-for="col in columns"
            :key="col.key"
            class="px-4 py-3 text-left text-xs font-semibold text-[#1E212B]"
            :class="col.width"
          >
            {{ col.label }}
          </th>
          <th v-if="hasActions" class="w-12 px-4 py-3" />
        </tr>
      </thead>

      <tbody v-if="!isLoading && filteredRows.length > 0">
        <tr
          v-for="row in filteredRows"
          :key="String(row[rowKey])"
          class="border-b border-[#EDEDEE] hover:bg-[#F6F6F7] cursor-pointer transition-colors last:border-0"
          @click="handleRowClick(row)"
        >
          <td
            v-for="col in columns"
            :key="col.key"
            class="px-4 py-4 text-sm text-[#4B4D55]"
            :class="col.width"
          >
            <!-- Named scoped slot per column key for custom rendering -->
            <slot :name="`cell-${col.key}`" :row="row" :value="row[col.key]">
              {{ row[col.key] ?? '—' }}
            </slot>
          </td>

          <!-- Row actions menu -->
          <td v-if="hasActions" class="px-4 py-4 relative" @click.stop>
            <button
              class="flex items-center justify-center w-8 h-8 rounded-full hover:bg-[#EDEDEE] text-[#4B4D55]"
              @click="toggleMenu(row[rowKey])"
            >
              <svg viewBox="0 0 24 24" fill="currentColor" class="w-4 h-4">
                <circle cx="12" cy="5" r="1.5" />
                <circle cx="12" cy="12" r="1.5" />
                <circle cx="12" cy="19" r="1.5" />
              </svg>
            </button>
            <div
              v-if="openMenuRow === row[rowKey]"
              class="absolute right-4 top-full mt-1 bg-white border border-[#EDEDEE] rounded-xl shadow-lg z-10 min-w-36 py-1"
            >
              <button
                v-for="action in actions"
                :key="action.key"
                class="w-full text-left px-4 py-2 text-sm hover:bg-[#F6F6F7] transition-colors"
                :class="action.class ?? 'text-[#1E212B]'"
                @click="handleAction(action.key, row)"
              >
                {{ action.label }}
              </button>
            </div>
          </td>
        </tr>
      </tbody>
    </table>

    <!-- Loading -->
    <div v-if="isLoading" class="py-12 text-center text-sm text-[#6F7177]">
      Loading...
    </div>

    <!-- Empty -->
    <div v-else-if="filteredRows.length === 0" class="py-16 text-center">
      <div class="flex items-center justify-center mb-4">
        <div class="w-16 h-16 rounded-full bg-[#EDEDEE]" />
      </div>
      <h5 class="text-base font-semibold text-[#1E212B]">
        {{ emptyTitle ?? 'No results found' }}
      </h5>
      <p v-if="emptyMessage" class="text-sm text-[#6F7177] mt-1">
        {{ emptyMessage }}
      </p>
    </div>

    <!-- Pagination -->
    <div
      v-if="!isLoading && pages > 1"
      class="flex items-center justify-between px-4 py-3 border-t border-[#EDEDEE]"
    >
      <p class="text-xs text-[#6F7177]">
        Page {{ currentPage }} of {{ pages }}
        <template v-if="total != null"> · {{ total }} records</template>
      </p>
      <div class="flex items-center gap-1">
        <button
          class="px-2.5 py-1.5 text-xs font-medium rounded-lg border border-[#DBDBDD] text-[#4B4D55] hover:bg-[#F6F6F7] disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
          :disabled="currentPage <= 1"
          @click="goToPage(currentPage - 1)"
        >
          Prev
        </button>
        <template v-for="p in pageNumbers" :key="String(p)">
          <span v-if="p === '...'" class="px-2 text-xs text-[#939499]">…</span>
          <button
            v-else
            class="w-7 h-7 text-xs font-medium rounded-lg transition-colors"
            :class="p === currentPage
              ? 'bg-[#1E212B] text-white'
              : 'border border-[#DBDBDD] text-[#4B4D55] hover:bg-[#F6F6F7]'"
            @click="goToPage(p as number)"
          >
            {{ p }}
          </button>
        </template>
        <button
          class="px-2.5 py-1.5 text-xs font-medium rounded-lg border border-[#DBDBDD] text-[#4B4D55] hover:bg-[#F6F6F7] disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
          :disabled="currentPage >= pages"
          @click="goToPage(currentPage + 1)"
        >
          Next
        </button>
      </div>
    </div>
  </div>
  </div>
</template>

<template>
  <div class="rounded-lg">
    <!-- Header -->
    <div class="flex justify-between items-center mb-5">
      <!-- Search -->
      <div class="relative flex-1 max-w-xs">
        <GroBasicSearch
          v-model="search"
          placeholder="Search"
          hint=""
          border-radius="8px"
        />
      </div>

      <!-- Filters Button -->
      <div
        class="relative flex items-center gap-1.5 py-2 px-3 text-[#1E212B] bg-white border border-gray-200 rounded-lg text-sm cursor-pointer hover:bg-gray-50 hover:border-gray-300 shadow-[inset_0_1px_0_0_#E3E3E3] transition-colors"
        @click="emit('filter-click')"
      >
        <HugeiconsIcon
          color="#1E212B"
          :icon="FilterIcon"
          class="h-4"
        />
        <span class="font-semibold text-xs">Filters</span>
        <span
          v-if="props.activeFiltersCount && props.activeFiltersCount > 0"
          class="flex items-center justify-center w-4 h-4 rounded-full bg-[#D26B06] text-white text-[10px] font-bold"
        >
          {{ props.activeFiltersCount }}
        </span>
      </div>
    </div>

    <!-- Table -->
    <div class="rounded-lg bg-white shadow overflow-hidden font-medium">
      <vue3-datatable
        :rows="props.rows"
        :columns="props.cols"
        :total-rows="total_rows"
        :is-server-mode="false"
        :page-size="params.pagesize"
        :has-checkbox="false"
        class="text-[#4B4D55] cursor-pointer"
        @row-click="(row) => props.onView?.(row)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import {computed, ref} from 'vue';
import Vue3Datatable from '@bhplugin/vue3-datatable';
import '@bhplugin/vue3-datatable/dist/style.css';
import { HugeiconsIcon } from '@hugeicons/vue'
import { FilterIcon } from '@hugeicons/core-free-icons'
import GroBasicSearch from "@/components/forms/input/GroBasicSearch.vue";
import { type ColsDefinition } from '@/composables/helpers/types/table'

const currentPage =  ref<number>(0);
const pageSize =  ref<number>(10);

const search = defineModel<string>();

const emit = defineEmits<{
  (e: 'filter-click'): void
}>()

const params = computed(() =>{
  return {
    current_page: currentPage,
    pagesize: pageSize
  }
})

const total_rows = computed(() =>{
  return props.rows?.length || 0;
})

const props = defineProps<{
  cols: ColsDefinition[]
  rows: [],
  onView?: (row) => void;
  activeFiltersCount?: number;
}>();
</script>

<style>
thead {
  background-color: #F6F6F7 !important;
  color: #1E212B !important;
}

tbody {
  background-color: #FFFFFF !important;
}

tbody tr {
  cursor: pointer;
  transition: background-color 0.15s ease;
}

tbody tr:hover {
  background-color: #F6F6F7 !important;
}

tbody tr:active {
  background-color: #EDEDEE !important;
  transform: scale(0.995);
  transition: background-color 0.05s ease, transform 0.05s ease;
}

.bh-pagination {
  padding-left: 16px !important;
  padding-right: 16px !important;
}

.bh-active {
  background-color: #1E212B !important;
  border: none !important;
  box-shadow: none !important;
}
</style>

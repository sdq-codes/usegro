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
        class="flex items-center py-2 px-3 text-[#1E212B] bg-white border border-gray-200 rounded-lg text-sm cursor-pointer hover:bg-gray-50 hover:border-gray-300 shadow-[inset_0_1px_0_0_#E3E3E3] transition-colors "
      >
        <HugeiconsIcon
          color="#1E212B"
          :icon="FilterIcon"
          class="h-4"
        />
        <span class="font-semibold text-xs">Filters</span>
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
        :has-checkbox="true"
        class="text-[#4B4D55]"
      >
        <template #actions="data">
          <div class="flex gap-4">
            <GroBasicButton
              color="primary"
              size="xs"
              shape="custom"
              class="w-max"
              @click="() => props.onView?.(data.value)"
            >
              View
            </GroBasicButton>
            <GroBasicButton
              color="tertiary"
              size="xs"
              shape="custom"
              class="w-max"
              @click="() => props.onDelete?.(data.value)"
            >
              Delete
            </GroBasicButton>
          </div>
        </template>
      </vue3-datatable>
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
import GroBasicButton from "@/components/buttons/GroBasicButton.vue";

const currentPage =  ref<number>(0);
const pageSize =  ref<number>(10);

const search = defineModel<string>();

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
  onDelete?: (row) => void;
}>();
</script>

<style >
thead {
  background-color: #F6F6F7 !important;
  color: #1E212B !important;
}

tbody {
  background-color: #FFFFFF !important;
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

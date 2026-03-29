<script setup lang="ts">
import { ref, computed } from 'vue'

interface Currency {
  code: string
  symbol: string
  label: string
}

interface Props {
  label?: string
  placeholder?: string
  currencies?: Currency[]
  defaultCurrency?: string
  hint?: string
}

const props = withDefaults(defineProps<Props>(), {
  label: '',
  placeholder: '0',
  hint: '',
  defaultCurrency: 'NGN',
  currencies: () => [
    { code: 'NGN', symbol: '₦', label: 'Nigerian Naira' },
    { code: 'USD', symbol: '$', label: 'US Dollar' },
    { code: 'EUR', symbol: '€', label: 'Euro' },
    { code: 'GBP', symbol: '£', label: 'British Pound' },
    { code: 'GHS', symbol: 'GH₵', label: 'Ghanaian Cedi' },
    { code: 'KES', symbol: 'KSh', label: 'Kenyan Shilling' },
    { code: 'ZAR', symbol: 'R', label: 'South African Rand' },
  ],
})

const amount = defineModel<string | number>('amount', { default: '' })
const currency = defineModel<string>('currency', { default: '' })

const selectedCode = ref(currency.value || props.defaultCurrency)
const showDropdown = ref(false)
const dropdownRef = ref<HTMLElement | null>(null)

const selectedCurrency = computed(() =>
  props.currencies.find(c => c.code === selectedCode.value) ?? props.currencies[0]
)

const selectCurrency = (code: string) => {
  selectedCode.value = code
  currency.value = code
  showDropdown.value = false
}

const onClickOutside = (e: MouseEvent) => {
  if (!dropdownRef.value?.contains(e.target as Node)) {
    showDropdown.value = false
  }
}

import { onMounted, onBeforeUnmount } from 'vue'
onMounted(() => document.addEventListener('mousedown', onClickOutside))
onBeforeUnmount(() => document.removeEventListener('mousedown', onClickOutside))
</script>

<template>
  <div>
    <label
      v-if="label"
      class="text-xs font-medium text-[#1E212B] mb-1 block"
    >
      {{ label }}
    </label>

    <div
      class="flex items-center bg-[#F6F6F7] border border-[#EDEDEE] rounded-lg hover:border-[#94BDD8] focus-within:border-[#1E212B] transition-colors overflow-visible"
    >
      <!-- Currency selector -->
      <div
        ref="dropdownRef"
        class="relative shrink-0"
      >
        <button
          type="button"
          class="flex bg-white items-center gap-1.5 px-3 py-2 h-full text-sm font-medium text-[#4B4D55] hover:text-[#1E212B] transition-colors border-r border-[#EDEDEE] rounded-l-lg hover:bg-[#EDEDEE]/60 cursor-pointer"
          @click="showDropdown = !showDropdown"
        >
          <span class="font-semibold">{{ selectedCurrency.symbol }}</span>
          <svg
            class="w-3 h-3 text-[#6F7177] transition-transform"
            :class="showDropdown ? 'rotate-180' : ''"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2.5"
          >
            <path d="M6 9l6 6 6-6" />
          </svg>
        </button>

        <!-- Dropdown -->
        <Transition
          enter-active-class="transition ease-out duration-150"
          enter-from-class="opacity-0 translate-y-1"
          enter-to-class="opacity-100 translate-y-0"
          leave-active-class="transition ease-in duration-100"
          leave-from-class="opacity-100"
          leave-to-class="opacity-0"
        >
          <div
            v-if="showDropdown"
            class="absolute left-0 top-full mt-1 bg-white border border-[#EDEDEE] rounded-xl shadow-lg z-50 min-w-48 py-1 overflow-hidden"
          >
            <button
              v-for="c in currencies"
              :key="c.code"
              type="button"
              class="w-full cursor-pointer flex items-center gap-3 px-4 py-2.5 text-sm hover:bg-[#F6F6F7] transition-colors text-left"
              :class="c.code === selectedCode ? 'bg-[#F6F6F7] font-semibold text-[#1E212B]' : 'text-[#4B4D55]'"
              @click="selectCurrency(c.code)"
            >
              <span class="w-6 text-center font-medium text-[#1E212B]">{{ c.symbol }}</span>
              <span class="flex-1">{{ c.label }}</span>
              <span class="text-xs text-[#939499]">{{ c.code }}</span>
              <svg
                v-if="c.code === selectedCode"
                class="w-3.5 h-3.5 text-[#AF513A]"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2.5"
              >
                <path d="M5 13l4 4L19 7" />
              </svg>
            </button>
          </div>
        </Transition>
      </div>

      <!-- Amount input -->
      <input
        v-model="amount"
        type="number"
        :placeholder="placeholder"
        min="0"
        class="flex-1 px-3 py-2 bg-transparent text-sm text-[#1E212B] outline-none placeholder-[#939499]"
      >
    </div>

    <p
      v-if="hint"
      class="text-xs text-[#6F7177] mt-1"
    >
      {{ hint }}
    </p>
  </div>
</template>

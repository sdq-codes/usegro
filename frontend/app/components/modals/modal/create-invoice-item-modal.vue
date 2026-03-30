<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useProductAPI } from '@/composables/api/catalog/product'
import { useServiceAPI } from '@/composables/api/catalog/service'
import type { CatalogEntry } from '@/pages/invoices/create-new-invoice.vue'

interface Props {
  type: 'product' | 'service'
}

const props = defineProps<Props>()

const model = defineModel<boolean>()

const emit = defineEmits<{
  (e: 'itemAdded', entry: CatalogEntry & { discountPercent?: number }): void
}>()

// ── Form state ────────────────────────────────────────────
const name = ref('')
const qty = ref(1)
const rate = ref<number | null>(null)
const currency = ref('USD')
const description = ref('')
const addToDiscount = ref(false)
const discountPercent = ref<number | null>(null)
const addToCatalogue = ref(true)
const showAdditionalInfo = ref(false)
const additionalInfo = ref('')
const isSaving = ref(false)
const error = ref('')

const currencies = ['USD', 'EUR', 'GBP', 'NGN', 'CAD', 'AUD']

const total = computed(() => {
  const r = rate.value ?? 0
  const q = qty.value ?? 1
  const base = r * q
  if (addToDiscount.value && discountPercent.value) {
    return base - (base * discountPercent.value) / 100
  }
  return base
})

const formatCurrency = (v: number) => {
  const sym: Record<string, string> = { USD: '$', EUR: '€', GBP: '£', NGN: '₦', CAD: 'CA$', AUD: 'A$' }
  return `${sym[currency.value] ?? ''}${v.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`
}

const reset = () => {
  name.value = ''
  qty.value = 1
  rate.value = null
  currency.value = 'USD'
  description.value = ''
  addToDiscount.value = false
  discountPercent.value = null
  addToCatalogue.value = true
  showAdditionalInfo.value = false
  additionalInfo.value = ''
  error.value = ''
}

watch(model, (v) => { if (!v) reset() })

const addItem = async () => {
  if (!name.value.trim()) { error.value = 'Name is required'; return }
  if (!rate.value || rate.value <= 0) { error.value = 'Rate must be greater than 0'; return }

  error.value = ''
  isSaving.value = true

  let catalogId = `custom-${Date.now()}`
  const defaultRate = rate.value

  if (addToCatalogue.value) {
    const payload = {
      name: name.value.trim(),
      description: description.value || undefined,
      price: rate.value,
      price_currency: currency.value,
      status: 'active',
      ...(addToDiscount.value && discountPercent.value ? { discount_percent: discountPercent.value } : {}),
    }

    const result = props.type === 'product'
      ? await useProductAPI().CreateProduct(payload)
      : await useServiceAPI().CreateService(payload)

    if (!result.success) {
      error.value = result.error ?? 'Failed to save to catalogue'
      isSaving.value = false
      return
    }

    const item = (result.data as any)?.data
    if (item?.id) catalogId = item.id
  }

  emit('itemAdded', {
    id: catalogId,
    name: name.value.trim(),
    type: props.type,
    defaultRate,
    discountPercent: (addToDiscount.value && discountPercent.value) ? discountPercent.value : undefined,
  })

  model.value = false
  isSaving.value = false
}
</script>

<template>
  <div
    v-if="model"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-[2px]"
    @mousedown.self="model = false"
  >
    <div class="bg-white rounded-2xl shadow-xl w-full max-w-lg mx-4 overflow-hidden">
      <!-- Header -->
      <div class="flex items-center justify-between px-6 py-5 border-b border-[#EDEDEE]">
        <h6 class="text-base font-semibold text-[#1E212B]">
          New {{ type === 'product' ? 'Product' : 'Service' }}
        </h6>
        <button
          class="w-7 h-7 flex items-center justify-center rounded-full hover:bg-[#F6F6F7] text-[#6F7177] hover:text-[#1E212B] transition-colors text-lg leading-none"
          @click="model = false"
        >
          ✕
        </button>
      </div>

      <!-- Body -->
      <div class="px-6 py-5 space-y-5 max-h-[70vh] overflow-y-auto">
        <!-- Name -->
        <div>
          <label class="block text-xs font-medium text-[#1E212B] mb-1.5">Name</label>
          <input
            v-model="name"
            type="text"
            placeholder="e.g. Product Design – Figma Prototype"
            class="w-full bg-[#F6F6F7] border border-[#EDEDEE] hover:border-[#94BDD8] focus:border-[#1E212B] rounded-xl px-3.5 py-2.5 text-sm text-[#1E212B] outline-none transition-colors placeholder-[#939499]"
          >
        </div>

        <!-- Qty + Rate + Currency -->
        <div class="grid grid-cols-3 gap-3">
          <div>
            <label class="block text-xs font-medium text-[#1E212B] mb-1.5">
              {{ type === 'service' ? 'Hours' : 'Qty' }}
            </label>
            <input
              v-model.number="qty"
              type="number"
              min="1"
              class="w-full bg-[#F6F6F7] border border-[#EDEDEE] hover:border-[#94BDD8] focus:border-[#1E212B] rounded-xl px-3.5 py-2.5 text-sm text-[#1E212B] outline-none transition-colors"
            >
          </div>
          <div>
            <label class="block text-xs font-medium text-[#1E212B] mb-1.5">Rate</label>
            <input
              v-model.number="rate"
              type="number"
              min="0"
              step="0.01"
              placeholder="0.00"
              class="w-full bg-[#F6F6F7] border border-[#EDEDEE] hover:border-[#94BDD8] focus:border-[#1E212B] rounded-xl px-3.5 py-2.5 text-sm text-[#1E212B] outline-none transition-colors placeholder-[#939499]"
            >
          </div>
          <div>
            <label class="block text-xs font-medium text-[#1E212B] mb-1.5">Currency</label>
            <select
              v-model="currency"
              class="w-full bg-[#F6F6F7] border border-[#EDEDEE] hover:border-[#94BDD8] focus:border-[#1E212B] rounded-xl px-3.5 py-2.5 text-sm text-[#1E212B] outline-none transition-colors cursor-pointer"
            >
              <option
                v-for="c in currencies"
                :key="c"
                :value="c"
              >
                {{ c }}
              </option>
            </select>
          </div>
        </div>

        <!-- Discount -->
        <div>
          <button
            class="flex items-center gap-1.5 text-sm font-semibold text-[#2176AE]"
            @click="addToDiscount = !addToDiscount"
          >
            <svg
              class="w-4 h-4"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2.5"
            >
              <path :d="addToDiscount ? 'M5 12h14' : 'M12 5v14M5 12h14'" />
            </svg>
            {{ addToDiscount ? 'Remove discount' : 'Add discount' }}
          </button>
          <div
            v-if="addToDiscount"
            class="mt-2 flex items-center gap-2"
          >
            <input
              v-model.number="discountPercent"
              type="number"
              min="0"
              max="100"
              placeholder="e.g. 10"
              class="w-24 bg-[#F6F6F7] border border-[#EDEDEE] focus:border-[#1E212B] rounded-xl px-3 py-2 text-sm text-[#1E212B] outline-none transition-colors"
            >
            <span class="text-sm text-[#6F7177]">% off</span>
          </div>
        </div>

        <!-- Total -->
        <div class="flex justify-end">
          <span class="text-sm font-semibold text-[#1E212B]">
            Total: <span class="text-[#E87117]">{{ formatCurrency(total) }}</span>
          </span>
        </div>

        <!-- Description -->
        <div>
          <label class="block text-xs font-medium text-[#1E212B] mb-1.5">Description</label>
          <textarea
            v-model="description"
            maxlength="100"
            rows="4"
            placeholder="Add a description…"
            class="w-full bg-[#F6F6F7] border border-[#EDEDEE] hover:border-[#94BDD8] focus:border-[#1E212B] rounded-xl px-3.5 py-2.5 text-sm text-[#1E212B] outline-none transition-colors resize-none placeholder-[#939499]"
          />
          <p class="text-right text-xs text-[#939499] mt-1">
            {{ description.length }}/100
          </p>
        </div>

        <!-- Additional info -->
        <div>
          <button
            class="flex items-center gap-1.5 text-sm font-semibold text-[#2176AE]"
            @click="showAdditionalInfo = !showAdditionalInfo"
          >
            <svg
              class="w-4 h-4"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2.5"
            >
              <path :d="showAdditionalInfo ? 'M5 12h14' : 'M12 5v14M5 12h14'" />
            </svg>
            Additional info
          </button>
          <textarea
            v-if="showAdditionalInfo"
            v-model="additionalInfo"
            rows="3"
            placeholder="Any extra notes…"
            class="mt-2 w-full bg-[#F6F6F7] border border-[#EDEDEE] hover:border-[#94BDD8] focus:border-[#1E212B] rounded-xl px-3.5 py-2.5 text-sm text-[#1E212B] outline-none transition-colors resize-none placeholder-[#939499]"
          />
        </div>

        <!-- Add to catalogue -->
        <label class="flex items-start gap-3 cursor-pointer select-none">
          <div class="relative mt-0.5 shrink-0">
            <input
              v-model="addToCatalogue"
              type="checkbox"
              class="sr-only"
            >
            <div
              class="w-5 h-5 rounded flex items-center justify-center border-2 transition-colors"
              :class="addToCatalogue ? 'bg-[#1E212B] border-[#1E212B]' : 'bg-white border-[#EDEDEE]'"
            >
              <svg
                v-if="addToCatalogue"
                class="w-3 h-3 text-white"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="3"
              >
                <path d="M5 13l4 4L19 7" />
              </svg>
            </div>
          </div>
          <div>
            <p class="text-sm font-medium text-[#1E212B]">
              Add item to {{ type === 'product' ? 'product' : 'service' }} catalogue
            </p>
            <p class="text-xs text-[#939499] mt-0.5">
              Adding this item to your catalogue will allow you to easily re-use it in future invoices
            </p>
          </div>
        </label>

        <!-- Error -->
        <p
          v-if="error"
          class="text-xs text-red-500"
        >
          {{ error }}
        </p>
      </div>

      <!-- Footer -->
      <div class="px-6 py-4 border-t border-[#EDEDEE] flex items-center justify-between">
        <a
          :href="`/catalog/${type === 'product' ? 'products' : 'services'}/new`"
          target="_blank"
          class="flex items-center gap-1.5 text-sm font-medium text-[#2176AE] hover:underline"
        >
          <svg
            class="w-3.5 h-3.5"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <path d="M7 17L17 7M17 7H7M17 7v10" />
          </svg>
          Edit in {{ type === 'product' ? 'Products' : 'Services' }} page
        </a>

        <div class="flex gap-3">
          <button
            class="px-5 py-2.5 text-sm font-medium text-[#1E212B] border border-[#EDEDEE] rounded-xl hover:bg-[#F6F6F7] transition-colors"
            @click="model = false"
          >
            Cancel
          </button>
          <button
            class="px-5 py-2.5 text-sm font-semibold text-white bg-[#1E212B] rounded-xl hover:bg-[#2d3140] transition-colors disabled:opacity-50"
            :disabled="isSaving"
            @click="addItem"
          >
            {{ isSaving ? 'Saving…' : 'Add item' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

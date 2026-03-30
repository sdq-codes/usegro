<script setup lang="ts">
import type { CatalogEntry, CatalogVariant, CatalogPlan, LineItem } from '@/pages/invoices/create-new-invoice.vue'

const props = defineProps<{ entry: CatalogEntry | null }>()

const model = defineModel<boolean>()

const emit = defineEmits<{
  (e: 'itemAdded', item: Omit<LineItem, 'id'>): void
}>()

const formatCurrency = (v: number, currency = 'USD') => {
  const sym: Record<string, string> = { USD: '$', EUR: '€', GBP: '£', NGN: '₦', CAD: 'CA$', AUD: 'A$' }
  return `${sym[currency] ?? '$'}${v.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`
}

const planCycleLabel = (plan: CatalogPlan) => {
  if (plan.plan_type === 'subscription' && plan.billing_cycle) {
    return `/ ${plan.billing_cycle}`
  }
  if (plan.plan_type === 'package') {
    const parts: string[] = []
    if (plan.session_count) parts.push(`${plan.session_count} sessions`)
    if (plan.validity_days) parts.push(`${plan.validity_days} days`)
    return parts.length ? `(${parts.join(', ')})` : ''
  }
  return ''
}

const selectVariant = (variant: CatalogVariant) => {
  emit('itemAdded', {
    catalogId: props.entry!.id,
    name: `${props.entry!.name} — ${variant.name}`,
    type: 'product',
    qty: 1,
    rate: variant.price,
    imageUrl: props.entry!.imageUrl,
    billingType: 'one-time',
  })
  model.value = false
}

const selectPlan = (plan: CatalogPlan) => {
  const isRecurring = plan.plan_type === 'subscription'
  emit('itemAdded', {
    catalogId: props.entry!.id,
    name: `${props.entry!.name} — ${plan.name}`,
    type: 'service',
    qty: 1,
    rate: plan.price,
    imageUrl: props.entry!.imageUrl,
    billingType: isRecurring ? 'recurring' : 'one-time',
    billingCycle: isRecurring ? (plan.billing_cycle ?? 'monthly') : undefined,
  })
  model.value = false
}

const addDefault = () => {
  emit('itemAdded', {
    catalogId: props.entry!.id,
    name: props.entry!.name,
    type: props.entry!.type,
    qty: 1,
    rate: props.entry!.defaultRate,
    imageUrl: props.entry!.imageUrl,
    billingType: 'one-time',
  })
  model.value = false
}
</script>

<template>
  <div
    v-if="model && entry"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-[2px]"
    @mousedown.self="model = false"
  >
    <div class="bg-white rounded-2xl shadow-xl w-full max-w-md mx-4 overflow-hidden">
      <!-- Header -->
      <div class="flex items-center justify-between px-6 py-5 border-b border-[#EDEDEE]">
        <div>
          <h6 class="text-base font-semibold text-[#1E212B]">
            {{ entry.name }}
          </h6>
          <p class="text-xs text-[#939499] mt-0.5">
            {{ entry.type === 'product' ? 'Select a variant' : 'Select a plan or package' }}
          </p>
        </div>
        <button
          class="w-7 h-7 flex items-center justify-center rounded-full hover:bg-[#F6F6F7] text-[#6F7177] hover:text-[#1E212B] transition-colors text-lg leading-none"
          @click="model = false"
        >
          ✕
        </button>
      </div>

      <!-- Product variants -->
      <div
        v-if="entry.type === 'product'"
        class="max-h-[60vh] overflow-y-auto"
      >
        <button
          v-for="variant in entry.variants"
          :key="variant.id"
          class="w-full flex items-center justify-between px-6 py-4 hover:bg-[#F6F6F7] transition-colors border-b border-[#EDEDEE] last:border-0 text-left"
          @click="selectVariant(variant)"
        >
          <div>
            <p class="text-sm font-medium text-[#1E212B]">
              {{ variant.name }}
            </p>
            <p
              v-if="variant.sku"
              class="text-xs text-[#939499] mt-0.5"
            >
              SKU: {{ variant.sku }}
            </p>
          </div>
          <span class="text-sm font-semibold text-[#1E212B] ml-4 shrink-0">
            {{ formatCurrency(variant.price) }}
          </span>
        </button>
      </div>

      <!-- Service plans -->
      <div
        v-else-if="entry.type === 'service'"
        class="max-h-[60vh] overflow-y-auto"
      >
        <button
          v-for="plan in entry.plans"
          :key="plan.id"
          class="w-full flex items-center justify-between px-6 py-4 hover:bg-[#F6F6F7] transition-colors border-b border-[#EDEDEE] last:border-0 text-left"
          @click="selectPlan(plan)"
        >
          <div>
            <div class="flex items-center gap-2">
              <p class="text-sm font-medium text-[#1E212B]">
                {{ plan.name }}
              </p>
              <span
                class="text-[10px] font-semibold uppercase px-1.5 py-0.5 rounded"
                :class="plan.plan_type === 'subscription'
                  ? 'bg-[#F0F7FF] text-[#2176AE]'
                  : 'bg-[#EEF8F5] text-[#00916E]'"
              >
                {{ plan.plan_type }}
              </span>
            </div>
            <p class="text-xs text-[#939499] mt-0.5">
              {{ planCycleLabel(plan) }}
            </p>
          </div>
          <span class="text-sm font-semibold text-[#1E212B] ml-4 shrink-0">
            {{ formatCurrency(plan.price, plan.price_currency) }}
          </span>
        </button>
      </div>

      <!-- Footer: add without selecting -->
      <div class="px-6 py-4 border-t border-[#EDEDEE] flex items-center justify-between">
        <p class="text-xs text-[#939499]">
          {{ entry.type === 'product' ? `Base price: ${formatCurrency(entry.defaultRate)}` : `Base price: ${formatCurrency(entry.defaultRate)}` }}
        </p>
        <button
          class="text-sm font-medium text-[#6F7177] hover:text-[#1E212B] transition-colors underline-offset-2 hover:underline"
          @click="addDefault"
        >
          Add without selecting
        </button>
      </div>
    </div>
  </div>
</template>

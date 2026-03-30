<script setup lang="ts">
import { computed } from 'vue'

interface LineItem {
  id: string
  name: string
  qty: number
  rate: number
  imageUrl?: string
  billingType: 'one-time' | 'recurring'
  billingCycle?: string
}

interface InvoiceData {
  businessName: string
  website: string
  email: string
  phone: string
  businessAddress: string
  city: string
  state: string
  country: string
  postalCode: string
  billedTo: { label: string; value: string }[]
  lineItems: LineItem[]
  invoiceNumber: string
  termsAndConditions: string
  subject: string
  taxRate: number
  dueDate?: string
  recurringStartDate?: string
  genDate?: string
}

const props = defineProps<{ invoiceData: InvoiceData }>()

const billedToDisplay = computed(() =>
  props.invoiceData.billedTo.map(c => c.label).join(', ')
)
const subtotal = computed(() =>
  (props.invoiceData.lineItems ?? []).reduce((sum, i) => sum + i.qty * i.rate, 0)
)
const taxAmount = computed(() => (subtotal.value * (props.invoiceData.taxRate || 0)) / 100)
const total = computed(() => subtotal.value + taxAmount.value)

const recurringItems = computed(() => (props.invoiceData.lineItems ?? []).filter(i => i.billingType === 'recurring'))
const oneTimeItems = computed(() => (props.invoiceData.lineItems ?? []).filter(i => i.billingType === 'one-time'))
const hasRecurring = computed(() => recurringItems.value.length > 0)
const hasOneTime = computed(() => oneTimeItems.value.length > 0)

const recurringByBillingCycle = computed(() => {
  const groups: Record<string, number> = {}
  for (const item of recurringItems.value) {
    const cycle = item.billingCycle ?? 'monthly'
    groups[cycle] = (groups[cycle] ?? 0) + item.qty * item.rate
  }
  return groups
})

const formatCurrency = (v: number) =>
  `$${v.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`

const hasItems = computed(() => (props.invoiceData.lineItems ?? []).length > 0)
const placeholder = computed(() => !hasItems.value)
</script>

<template>
  <div class="bg-white shadow-[0px_35.16px_70.32px_0px_#4E4D593D] rounded-lg overflow-hidden">
    <!-- Dark navy header -->
    <div class="bg-[#1E212B] px-6 py-5">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-3">
          <img
            src="https://res.cloudinary.com/sdq121/image/upload/v1764981178/ccrh2qinydoihgbswden.svg"
            class="w-10 h-10"
          >
          <div>
            <h3 class="text-white font-bold text-base leading-tight">{{ invoiceData.businessName }}</h3>
            <p class="text-gray-400 text-xs">{{ invoiceData.website }}</p>
          </div>
        </div>
        <div class="text-right">
          <p class="text-[#E87117] font-black text-2xl tracking-tight leading-none">INVOICE</p>
          <p v-if="invoiceData.invoiceNumber" class="text-gray-300 text-xs mt-1 font-mono">
            #{{ invoiceData.invoiceNumber }}
          </p>
          <div v-else class="w-20 h-3 bg-gray-600 rounded animate-pulse mt-1 ml-auto" />
        </div>
      </div>
    </div>

    <!-- Orange accent bar -->
    <div class="h-1 bg-[#E87117]" />

    <!-- Info section -->
    <div class="px-6 pt-5 pb-4">
      <div class="grid grid-cols-2 gap-6">
        <!-- Left: billed to + address -->
        <div>
          <p class="text-[10px] font-bold text-[#E87117] uppercase tracking-widest mb-2">Billed To</p>
          <div v-if="!billedToDisplay" class="w-3/4 h-3 bg-gray-200 rounded animate-pulse" />
          <p v-else class="text-sm font-semibold text-[#1E212B]">{{ billedToDisplay }}</p>
        </div>

        <!-- Right: dates + amount -->
        <div class="text-right space-y-3">
          <div>
            <p class="text-[10px] font-bold text-gray-400 uppercase tracking-widest mb-0.5">Amount Due</p>
            <p class="text-2xl font-black text-[#E87117]">{{ formatCurrency(total) }}</p>
          </div>
          <div class="flex justify-end gap-6 text-xs text-gray-500">
            <div>
              <span class="font-semibold text-[#1E212B]">Issued</span>
              <span class="ml-1">{{ invoiceData.genDate || '---' }}</span>
            </div>
            <div v-if="hasOneTime">
              <span class="font-semibold text-[#1E212B]">Due</span>
              <span class="ml-1">{{ invoiceData.dueDate || '---' }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Business contact row -->
      <div class="flex gap-4 mt-4 pt-4 border-t border-[#F0F0F0]">
        <p class="text-xs text-gray-400">{{ invoiceData.email }}</p>
        <p class="text-xs text-gray-400">{{ invoiceData.phone }}</p>
        <p class="text-xs text-gray-400">
          {{ invoiceData.businessAddress }}, {{ invoiceData.city }}, {{ invoiceData.state }} {{ invoiceData.postalCode }}
        </p>
      </div>
    </div>

    <!-- Recurring badge -->
    <div v-if="hasRecurring" class="mx-6 mb-4 bg-[#1E212B] rounded-xl px-4 py-3 flex items-center gap-2">
      <svg class="w-3.5 h-3.5 text-[#E87117] shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
      </svg>
      <p class="text-xs text-gray-300">
        <template v-for="(amount, cycle) in recurringByBillingCycle" :key="cycle">
          <span class="text-[#E87117] font-semibold">{{ formatCurrency(amount) }}</span> / {{ cycle }}
        </template>
        <template v-if="invoiceData.recurringStartDate"> · starts {{ invoiceData.recurringStartDate }}</template>
      </p>
    </div>

    <!-- Items table -->
    <div class="mx-6 mb-5 rounded-xl border border-[#EDEDEE] overflow-hidden">
      <table class="w-full">
        <thead>
          <tr class="bg-[#1E212B]">
            <th class="text-left text-[10px] font-bold text-gray-300 px-4 py-3 uppercase tracking-widest">Item</th>
            <th class="text-center text-[10px] font-bold text-gray-300 px-3 py-3 uppercase tracking-widest">Qty</th>
            <th class="text-center text-[10px] font-bold text-gray-300 px-3 py-3 uppercase tracking-widest">Rate</th>
            <th class="text-right text-[10px] font-bold text-gray-300 px-4 py-3 uppercase tracking-widest">Amount</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="placeholder" class="border-b border-[#EDEDEE]">
            <td class="px-4 py-3"><div class="w-3/4 h-3 bg-gray-200 rounded animate-pulse" /></td>
            <td class="px-3 py-3 text-center"><div class="w-4 h-3 bg-gray-200 rounded animate-pulse mx-auto" /></td>
            <td class="px-3 py-3 text-center"><div class="w-10 h-3 bg-gray-200 rounded animate-pulse mx-auto" /></td>
            <td class="px-4 py-3 text-right"><div class="w-12 h-3 bg-gray-200 rounded animate-pulse ml-auto" /></td>
          </tr>
          <tr
            v-for="(item, idx) in invoiceData.lineItems"
            v-else
            :key="item.id"
            :class="idx % 2 === 0 ? 'bg-white' : 'bg-[#FAFAFA]'"
            class="border-b border-[#F0F0F0] last:border-0"
          >
            <td class="px-4 py-3">
              <div class="flex items-center gap-2">
                <div class="w-7 h-7 rounded-md overflow-hidden bg-[#F6F6F7] border border-[#EDEDEE] shrink-0 flex items-center justify-center">
                  <img v-if="item.imageUrl" :src="item.imageUrl" :alt="item.name" class="w-full h-full object-cover">
                  <svg v-else class="w-3.5 h-3.5 text-[#D7DAE0]" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="3" y="3" width="18" height="18" rx="2" /><path d="M3 9l4-4 4 4 4-5 4 5" /></svg>
                </div>
                <div>
                  <p class="text-sm font-medium text-[#1E212B]">{{ item.name }}</p>
                  <p v-if="item.billingType === 'recurring'" class="text-[10px] text-[#2176AE]">
                    Recurring · {{ item.billingCycle ?? 'monthly' }}
                  </p>
                </div>
              </div>
            </td>
            <td class="px-3 py-3 text-sm text-center text-gray-600">{{ item.qty }}</td>
            <td class="px-3 py-3 text-sm text-center text-gray-600">{{ formatCurrency(item.rate) }}</td>
            <td class="px-4 py-3 text-sm text-right font-semibold text-[#1E212B]">{{ formatCurrency(item.qty * item.rate) }}</td>
          </tr>
        </tbody>
      </table>

      <!-- Totals -->
      <div class="bg-[#FAFAFA] border-t border-[#EDEDEE] px-4 py-3 space-y-1.5">
        <div class="flex justify-between text-xs text-gray-500">
          <span>Subtotal</span>
          <span class="font-medium text-[#1E212B]">{{ formatCurrency(subtotal) }}</span>
        </div>
        <div class="flex justify-between text-xs text-gray-500">
          <span>Tax ({{ invoiceData.taxRate }}%)</span>
          <span class="font-medium text-[#1E212B]">{{ invoiceData.taxRate ? formatCurrency(taxAmount) : '--' }}</span>
        </div>
        <div class="flex justify-between items-center pt-2 border-t border-[#E0E0E0]">
          <span class="text-sm font-bold text-[#1E212B]">Total</span>
          <span class="text-lg font-black text-[#E87117]">{{ formatCurrency(total) }}</span>
        </div>
      </div>
    </div>

    <!-- Notes & Terms -->
    <div v-if="invoiceData.subject || invoiceData.termsAndConditions" class="mx-6 mb-6 space-y-4">
      <div v-if="invoiceData.subject" class="bg-[#FFF8F2] border border-[#FFD8B0] rounded-xl px-4 py-3">
        <p class="text-[10px] font-bold text-[#E87117] uppercase tracking-widest mb-1">Notes</p>
        <p class="text-xs text-gray-600">{{ invoiceData.subject }}</p>
      </div>
      <div v-if="invoiceData.termsAndConditions">
        <p class="text-[10px] font-bold text-gray-400 uppercase tracking-widest mb-1">Terms & Conditions</p>
        <p class="text-xs text-gray-500 leading-relaxed">{{ invoiceData.termsAndConditions }}</p>
      </div>
    </div>

    <!-- Footer band -->
    <div class="bg-[#1E212B] px-6 py-3 flex justify-between items-center">
      <p class="text-[10px] text-gray-500">Powered by Gro</p>
      <p class="text-[10px] text-gray-500">{{ invoiceData.email }}</p>
    </div>
  </div>
</template>

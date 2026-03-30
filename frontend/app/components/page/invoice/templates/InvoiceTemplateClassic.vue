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

const oneTimeItems = computed(() => (props.invoiceData.lineItems ?? []).filter(i => i.billingType === 'one-time'))
const recurringItems = computed(() => (props.invoiceData.lineItems ?? []).filter(i => i.billingType === 'recurring'))
const oneTimeSubtotal = computed(() => oneTimeItems.value.reduce((sum, i) => sum + i.qty * i.rate, 0))
const oneTimeTax = computed(() => (oneTimeSubtotal.value * (props.invoiceData.taxRate || 0)) / 100)
const oneTimeTotal = computed(() => oneTimeSubtotal.value + oneTimeTax.value)
const hasOneTime = computed(() => oneTimeItems.value.length > 0)
const hasRecurring = computed(() => recurringItems.value.length > 0)

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
  <div class="bg-white shadow-[0px_35.16px_70.32px_0px_#4E4D593D] rounded-lg px-4 py-8">
    <!-- Header -->
    <div class="flex justify-between items-start mb-6 px-2">
      <div class="flex gap-x-4">
        <img
          src="https://res.cloudinary.com/sdq121/image/upload/v1764981178/ccrh2qinydoihgbswden.svg"
          class="w-14 h-14"
        >
        <div class="space-y-0.5">
          <h3 class="text-lg font-semibold text-[#E87117]">
            {{ invoiceData.businessName }}
          </h3>
          <p class="text-xs text-[#5E6470]">{{ invoiceData.website }}</p>
          <p class="text-xs text-[#5E6470]">{{ invoiceData.email }}</p>
          <p class="text-xs text-[#5E6470]">{{ invoiceData.phone }}</p>
        </div>
      </div>
      <div class="text-right space-y-0.5">
        <p class="text-xs text-[#5E6470]">{{ invoiceData.businessAddress }}</p>
        <p class="text-xs text-[#5E6470]">
          {{ invoiceData.city }}, {{ invoiceData.state }}, {{ invoiceData.country }} {{ invoiceData.postalCode }}
        </p>
      </div>
    </div>

    <div class="rounded-2xl border border-[#D7DAE0] bg-[#FAFAFA]">
      <!-- Meta grid -->
      <div class="p-5">
        <div class="grid grid-cols-3 gap-4 mb-5">
          <div>
            <p class="text-[10px] text-gray-400 uppercase tracking-wider mb-1">Billed to</p>
            <div v-if="!billedToDisplay" class="w-4/5 h-3 bg-gray-200 rounded animate-pulse" />
            <p v-else class="text-sm font-medium text-[#1E212B]">{{ billedToDisplay }}</p>
          </div>
          <div>
            <p class="text-[10px] text-gray-400 uppercase tracking-wider mb-1">Invoice #</p>
            <div v-if="!invoiceData.invoiceNumber" class="w-3/5 h-3 bg-gray-200 rounded animate-pulse" />
            <p v-else class="text-sm font-semibold text-[#1E212B]">{{ invoiceData.invoiceNumber }}</p>
          </div>
          <div class="text-right">
            <p class="text-[10px] text-gray-400 uppercase tracking-wider mb-1">Amount (USD)</p>
            <p class="text-xl font-bold text-[#E87117]">{{ formatCurrency(total) }}</p>
          </div>
        </div>
        <div class="grid grid-cols-3 gap-4">
          <div>
            <p class="text-[10px] text-gray-400 uppercase tracking-wider mb-1">Date issued</p>
            <p class="text-sm font-semibold text-[#1E212B]">{{ invoiceData.genDate || '---' }}</p>
          </div>
          <div>
            <p class="text-[10px] text-gray-400 uppercase tracking-wider mb-1">Due date</p>
            <p class="text-sm font-semibold text-[#1E212B]">{{ invoiceData.dueDate || '---' }}</p>
          </div>
          <div class="text-right">
            <p class="text-[10px] text-gray-400 uppercase tracking-wider mb-1">Payment</p>
            <p class="text-sm font-semibold text-[#1E212B]">
              {{ hasRecurring && hasOneTime ? 'Mixed' : hasRecurring ? 'Recurring' : 'One-time' }}
            </p>
          </div>
        </div>

        <!-- Recurring summary -->
        <div v-if="hasRecurring" class="mt-4 bg-[#F0F7FF] rounded-lg px-4 py-2.5 flex items-center gap-2">
          <svg class="w-3.5 h-3.5 text-[#2176AE] shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          <p class="text-xs text-[#2176AE]">
            <template v-for="(amount, cycle) in recurringByBillingCycle" :key="cycle">
              {{ formatCurrency(amount) }} / {{ cycle }}
            </template>
            <template v-if="invoiceData.recurringStartDate"> · starts {{ invoiceData.recurringStartDate }}</template>
          </p>
        </div>
      </div>

      <!-- Items table -->
      <div class="border-t border-[#D7DAE0] px-2 pb-2">
        <table class="w-full">
          <thead>
            <tr class="border-b border-[#D7DAE0]">
              <th class="text-left text-[10px] font-semibold text-gray-500 px-4 py-2.5 uppercase tracking-wider">Item</th>
              <th class="text-center text-[10px] font-semibold text-gray-500 px-3 py-2.5 uppercase tracking-wider">Qty</th>
              <th class="text-center text-[10px] font-semibold text-gray-500 px-3 py-2.5 uppercase tracking-wider">Rate</th>
              <th class="text-right text-[10px] font-semibold text-gray-500 px-4 py-2.5 uppercase tracking-wider">Amount</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="placeholder" class="border-b border-gray-100">
              <td class="px-4 py-3"><div class="w-3/4 h-3 bg-gray-200 rounded animate-pulse" /></td>
              <td class="px-3 py-3 text-center"><div class="w-4 h-3 bg-gray-200 rounded animate-pulse mx-auto" /></td>
              <td class="px-3 py-3 text-center"><div class="w-10 h-3 bg-gray-200 rounded animate-pulse mx-auto" /></td>
              <td class="px-4 py-3 text-right"><div class="w-12 h-3 bg-gray-200 rounded animate-pulse ml-auto" /></td>
            </tr>
            <tr
              v-for="item in invoiceData.lineItems"
              v-else
              :key="item.id"
              class="border-b border-gray-100 last:border-0"
            >
              <td class="px-4 py-3">
                <div class="flex items-center gap-2">
                  <div class="w-7 h-7 rounded-md overflow-hidden bg-[#F6F6F7] border border-[#EDEDEE] shrink-0 flex items-center justify-center">
                    <img v-if="item.imageUrl" :src="item.imageUrl" :alt="item.name" class="w-full h-full object-cover">
                    <svg v-else class="w-3.5 h-3.5 text-[#D7DAE0]" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="3" y="3" width="18" height="18" rx="2" /><path d="M3 9l4-4 4 4 4-5 4 5" /></svg>
                  </div>
                  <span class="text-sm font-medium text-[#1E212B]">{{ item.name }}</span>
                </div>
              </td>
              <td class="px-3 py-3 text-sm text-center text-gray-600">{{ item.qty }}</td>
              <td class="px-3 py-3 text-sm text-center text-gray-600">{{ formatCurrency(item.rate) }}</td>
              <td class="px-4 py-3 text-sm text-right font-medium text-[#1E212B]">{{ formatCurrency(item.qty * item.rate) }}</td>
            </tr>
          </tbody>
          <tfoot class="border-t border-[#EDEDEE]">
            <tr>
              <td class="px-4 py-2" /><td class="px-3 py-2 text-xs text-gray-500 text-center">Subtotal</td>
              <td class="px-3 py-2" /><td class="px-4 py-2 text-xs font-medium text-[#1E212B] text-right">{{ formatCurrency(subtotal) }}</td>
            </tr>
            <tr>
              <td class="px-4 py-2" /><td class="px-3 py-2 text-xs text-gray-500 text-center">Tax ({{ invoiceData.taxRate }}%)</td>
              <td class="px-3 py-2" />
              <td v-if="invoiceData.taxRate" class="px-4 py-2 text-xs font-medium text-[#1E212B] text-right">{{ formatCurrency(taxAmount) }}</td>
              <td v-else class="px-4 py-2 text-xs font-medium text-[#1E212B] text-right">--</td>
            </tr>
            <tr class="border-t border-gray-200">
              <td class="px-4 py-2.5" /><td class="px-3 py-2.5 text-sm font-bold text-[#1E212B] text-center">Total</td>
              <td class="px-3 py-2.5" /><td class="px-4 py-2.5 text-sm font-bold text-[#E87117] text-right">{{ formatCurrency(total) }}</td>
            </tr>
          </tfoot>
        </table>
      </div>
    </div>

    <!-- Notes -->
    <div v-if="invoiceData.subject" class="mt-5 px-2">
      <p class="text-xs font-semibold text-[#6F7177] uppercase tracking-wider mb-1">Notes</p>
      <p class="text-sm text-gray-600">{{ invoiceData.subject }}</p>
    </div>

    <!-- Terms -->
    <div v-if="invoiceData.termsAndConditions" class="mt-5 px-2">
      <p class="text-xs font-semibold text-[#6F7177] uppercase tracking-wider mb-1">Terms & Conditions</p>
      <p class="text-sm text-gray-600 leading-relaxed">{{ invoiceData.termsAndConditions }}</p>
    </div>
  </div>
</template>

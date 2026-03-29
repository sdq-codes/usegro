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
  <div class="bg-white shadow-[0px_35.16px_70.32px_0px_#4E4D593D] rounded-lg px-7 py-8">
    <!-- Top: Invoice identity + Business info side by side -->
    <div class="flex justify-between items-start mb-8">
      <!-- Left: INVOICE heading + number -->
      <div>
        <h1 class="text-3xl font-black text-[#1E212B] tracking-tight leading-none">INVOICE</h1>
        <div class="mt-2 space-y-0.5">
          <div v-if="invoiceData.invoiceNumber" class="flex items-center gap-1.5">
            <span class="text-xs text-gray-400">#</span>
            <span class="text-sm font-mono font-semibold text-[#2176AE]">{{ invoiceData.invoiceNumber }}</span>
          </div>
          <div v-else class="w-24 h-3 bg-gray-100 rounded animate-pulse" />
          <p class="text-xs text-gray-400">{{ invoiceData.genDate || '---' }}</p>
        </div>
      </div>

      <!-- Right: Business info -->
      <div class="text-right">
        <p class="text-sm font-bold text-[#1E212B]">{{ invoiceData.businessName }}</p>
        <p class="text-xs text-gray-400 mt-0.5">{{ invoiceData.website }}</p>
        <p class="text-xs text-gray-400">{{ invoiceData.email }}</p>
        <p class="text-xs text-gray-400">{{ invoiceData.phone }}</p>
      </div>
    </div>

    <!-- Thin divider -->
    <div class="h-px bg-[#1E212B]" />

    <!-- Billed to + Due date row -->
    <div class="flex justify-between items-start py-5">
      <div>
        <p class="text-[9px] font-bold text-gray-400 uppercase tracking-widest mb-1.5">Bill To</p>
        <div v-if="!billedToDisplay" class="w-32 h-3 bg-gray-100 rounded animate-pulse" />
        <p v-else class="text-sm font-semibold text-[#1E212B]">{{ billedToDisplay }}</p>
        <p class="text-xs text-gray-400 mt-0.5">
          {{ invoiceData.businessAddress }}, {{ invoiceData.city }}
        </p>
      </div>
      <div class="text-right space-y-2">
        <div v-if="hasOneTime && invoiceData.dueDate">
          <p class="text-[9px] font-bold text-gray-400 uppercase tracking-widest mb-0.5">Due Date</p>
          <p class="text-sm font-semibold text-[#1E212B]">{{ invoiceData.dueDate }}</p>
        </div>
        <div v-if="hasRecurring">
          <p class="text-[9px] font-bold text-gray-400 uppercase tracking-widest mb-0.5">Billing</p>
          <p class="text-xs font-semibold text-[#2176AE]">
            <template v-for="(amount, cycle) in recurringByBillingCycle" :key="cycle">
              {{ formatCurrency(amount) }} / {{ cycle }}
            </template>
          </p>
        </div>
      </div>
    </div>

    <!-- Thin divider -->
    <div class="h-px bg-[#EDEDEE]" />

    <!-- Items table -->
    <table class="w-full mt-4">
      <thead>
        <tr>
          <th class="text-left text-[9px] font-bold text-gray-400 pb-3 uppercase tracking-widest">Description</th>
          <th class="text-center text-[9px] font-bold text-gray-400 pb-3 uppercase tracking-widest w-12">Qty</th>
          <th class="text-center text-[9px] font-bold text-gray-400 pb-3 uppercase tracking-widest w-20">Rate</th>
          <th class="text-right text-[9px] font-bold text-gray-400 pb-3 uppercase tracking-widest w-24">Amount</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-[#F0F0F0]">
        <tr v-if="placeholder">
          <td class="py-3"><div class="w-3/4 h-3 bg-gray-100 rounded animate-pulse" /></td>
          <td class="py-3 text-center"><div class="w-4 h-3 bg-gray-100 rounded animate-pulse mx-auto" /></td>
          <td class="py-3 text-center"><div class="w-10 h-3 bg-gray-100 rounded animate-pulse mx-auto" /></td>
          <td class="py-3 text-right"><div class="w-14 h-3 bg-gray-100 rounded animate-pulse ml-auto" /></td>
        </tr>
        <tr
          v-for="item in invoiceData.lineItems"
          v-else
          :key="item.id"
        >
          <td class="py-3 pr-4">
            <p class="text-sm font-medium text-[#1E212B]">{{ item.name }}</p>
            <p v-if="item.billingType === 'recurring'" class="text-[10px] text-[#2176AE] mt-0.5">
              {{ item.billingCycle ?? 'monthly' }}
            </p>
          </td>
          <td class="py-3 text-sm text-center text-gray-500">{{ item.qty }}</td>
          <td class="py-3 text-sm text-center text-gray-500">{{ formatCurrency(item.rate) }}</td>
          <td class="py-3 text-sm text-right font-medium text-[#1E212B]">{{ formatCurrency(item.qty * item.rate) }}</td>
        </tr>
      </tbody>
    </table>

    <!-- Thin divider -->
    <div class="h-px bg-[#EDEDEE] mt-2" />

    <!-- Totals -->
    <div class="mt-3 space-y-1.5 w-64 ml-auto">
      <div class="flex justify-between text-xs text-gray-400">
        <span>Subtotal</span>
        <span class="text-[#1E212B] font-medium">{{ formatCurrency(subtotal) }}</span>
      </div>
      <div class="flex justify-between text-xs text-gray-400">
        <span>Tax ({{ invoiceData.taxRate }}%)</span>
        <span class="text-[#1E212B] font-medium">{{ invoiceData.taxRate ? formatCurrency(taxAmount) : '--' }}</span>
      </div>
      <div class="h-px bg-[#EDEDEE]" />
      <div class="flex justify-between items-center pt-0.5">
        <span class="text-sm font-bold text-[#1E212B]">Total</span>
        <span class="text-xl font-black text-[#2176AE]">{{ formatCurrency(total) }}</span>
      </div>
      <p v-if="hasRecurring && invoiceData.recurringStartDate" class="text-[10px] text-gray-400 text-right">
        Recurring starts {{ invoiceData.recurringStartDate }}
      </p>
    </div>

    <!-- Notes & Terms -->
    <div v-if="invoiceData.subject || invoiceData.termsAndConditions" class="mt-8 pt-5 border-t border-[#EDEDEE] space-y-4">
      <div v-if="invoiceData.subject">
        <p class="text-[9px] font-bold text-gray-400 uppercase tracking-widest mb-1">Note</p>
        <p class="text-xs text-gray-600">{{ invoiceData.subject }}</p>
      </div>
      <div v-if="invoiceData.termsAndConditions">
        <p class="text-[9px] font-bold text-gray-400 uppercase tracking-widest mb-1">Terms & Conditions</p>
        <p class="text-xs text-gray-500 leading-relaxed">{{ invoiceData.termsAndConditions }}</p>
      </div>
    </div>

    <!-- Footer -->
    <div class="mt-8 pt-4 border-t border-[#EDEDEE] flex justify-between items-center">
      <p class="text-[9px] text-gray-300 uppercase tracking-widest">Powered by Gro</p>
      <p class="text-[9px] text-gray-400">{{ invoiceData.email }}</p>
    </div>
  </div>
</template>

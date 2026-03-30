<script setup lang="ts">
import { ref } from 'vue'
import type { LineItem } from '@/pages/invoices/create-new-invoice.vue'
import type { InvoiceTemplate } from '@/constants/invoice/templates'
import InvoiceTemplateClassic from '@/components/page/invoice/templates/InvoiceTemplateClassic.vue'
import InvoiceTemplateModern from '@/components/page/invoice/templates/InvoiceTemplateModern.vue'
import InvoiceTemplateMinimal from '@/components/page/invoice/templates/InvoiceTemplateMinimal.vue'
import { computed } from 'vue'

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

const props = withDefaults(defineProps<{
  invoiceData: InvoiceData
  selectedTemplate?: InvoiceTemplate
}>(), {
  invoiceData: () => ({
    businessName: 'UseGro.inc',
    website: 'www.website.com',
    email: 'hello@email.com',
    phone: '+91 00000 00000',
    businessAddress: 'Business address',
    city: 'City',
    state: 'State',
    country: 'US',
    postalCode: '00000',
    billedTo: [],
    lineItems: [],
    invoiceNumber: '',
    termsAndConditions: '',
    subject: '',
    taxRate: 0,
  }),
  selectedTemplate: 'classic',
})

type TabKey = 'pdf' | 'payer' | 'email'
const activeTab = ref<TabKey>('pdf')

const tabs: { key: TabKey; label: string }[] = [
  { key: 'pdf', label: 'Invoice PDF' },
  { key: 'payer', label: 'Payer Page' },
  { key: 'email', label: 'Email Preview' },
]

const templateComponent = computed(() => {
  if (props.selectedTemplate === 'modern') return InvoiceTemplateModern
  if (props.selectedTemplate === 'minimal') return InvoiceTemplateMinimal
  return InvoiceTemplateClassic
})

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
</script>

<template>
  <div class="w-1/2 overflow-y-auto rounded-r-2xl border-y border-r bg-[#EDF4F9] border-[#EDEDEE]">
    <h2 class="text-xl text-[#1E212B] mb-6 pt-8 px-16">
      Preview
    </h2>

    <!-- Tabs -->
    <div class="flex border-b border-[#EDEDEE] mb-8 px-16 gap-x-6">
      <button
        v-for="tab in tabs"
        :key="tab.key"
        class="pb-3 font-medium text-sm transition-colors relative cursor-pointer whitespace-nowrap"
        :class="activeTab === tab.key ? 'text-[#E87117]' : 'text-gray-500 hover:text-gray-700'"
        @click="activeTab = tab.key"
      >
        {{ tab.label }}
        <span
          v-if="activeTab === tab.key"
          class="absolute bottom-0 left-0 right-0 h-0.5 bg-[#E87117]"
        />
      </button>
    </div>

    <!-- ── Invoice PDF ─────────────────────────────────────────── -->
    <div v-if="activeTab === 'pdf'" class="px-16 pb-10">
      <component :is="templateComponent" :invoice-data="invoiceData" />
    </div>

    <!-- ── Payer Page ──────────────────────────────────────────── -->
    <div
      v-else-if="activeTab === 'payer'"
      class="px-16 pb-10"
    >
      <div class="bg-white shadow-[0px_35.16px_70.32px_0px_#4E4D593D] rounded-lg overflow-hidden">
        <!-- Header band -->
        <div class="bg-[#1E212B] px-8 py-5 flex items-center justify-between">
          <div class="flex items-center gap-3">
            <img
              src="https://res.cloudinary.com/sdq121/image/upload/v1764981178/ccrh2qinydoihgbswden.svg"
              class="w-8 h-8"
            >
            <span class="text-white font-semibold text-sm">{{ invoiceData.businessName }}</span>
          </div>
          <div class="flex items-center gap-1 text-gray-400 text-xs">
            <svg class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="3" y="11" width="18" height="11" rx="2" /><path d="M7 11V7a5 5 0 0110 0v4" />
            </svg>
            Secure payment
          </div>
        </div>

        <div class="px-8 py-6">
          <!-- Invoice summary card -->
          <div class="bg-[#FAFAFA] rounded-xl p-5 mb-6 border border-[#EDEDEE]">
            <p class="text-[10px] font-semibold text-[#6F7177] uppercase tracking-wider mb-3">Invoice Summary</p>
            <div class="space-y-2">
              <div class="flex justify-between text-sm">
                <span class="text-gray-500">Billed to</span>
                <span class="font-medium text-[#1E212B]">{{ billedToDisplay || 'Customer' }}</span>
              </div>
              <div v-if="invoiceData.invoiceNumber" class="flex justify-between text-sm">
                <span class="text-gray-500">Invoice #</span>
                <span class="font-medium text-[#1E212B]">{{ invoiceData.invoiceNumber }}</span>
              </div>
              <div v-if="hasOneTime && invoiceData.dueDate" class="flex justify-between text-sm">
                <span class="text-gray-500">Due date</span>
                <span class="font-medium text-[#1E212B]">{{ invoiceData.dueDate }}</span>
              </div>
              <template v-if="hasRecurring">
                <div v-for="(amount, cycle) in recurringByBillingCycle" :key="cycle" class="flex justify-between text-sm">
                  <span class="text-gray-500">Every {{ cycle }}</span>
                  <span class="font-medium text-[#1E212B]">{{ formatCurrency(amount) }}</span>
                </div>
              </template>
            </div>
            <div class="flex justify-between items-center pt-3 border-t border-[#EDEDEE] mt-3">
              <span class="text-sm font-bold text-[#1E212B]">
                {{ hasOneTime ? 'Due today' : 'Per cycle' }}
              </span>
              <span class="text-2xl font-bold text-[#E87117]">
                {{ formatCurrency(hasOneTime ? oneTimeTotal : Object.values(recurringByBillingCycle)[0] ?? 0) }}
              </span>
            </div>
          </div>

          <!-- Items -->
          <div v-if="hasItems" class="mb-6">
            <p class="text-[10px] font-semibold text-[#6F7177] uppercase tracking-wider mb-3">Items</p>
            <div
              v-for="item in invoiceData.lineItems"
              :key="item.id"
              class="flex justify-between items-center py-2 border-b border-[#EDEDEE] last:border-0"
            >
              <div class="flex items-center gap-2.5">
                <div class="w-8 h-8 rounded-lg overflow-hidden bg-[#F6F6F7] border border-[#EDEDEE] shrink-0 flex items-center justify-center">
                  <img v-if="item.imageUrl" :src="item.imageUrl" :alt="item.name" class="w-full h-full object-cover">
                  <svg v-else class="w-3.5 h-3.5 text-[#D7DAE0]" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="3" y="3" width="18" height="18" rx="2" /><path d="M3 9l4-4 4 4 4-5 4 5" /></svg>
                </div>
                <div>
                  <p class="text-sm font-medium text-[#1E212B]">{{ item.name }}</p>
                  <p class="text-xs text-gray-500">{{ item.qty }} × {{ formatCurrency(item.rate) }}</p>
                </div>
              </div>
              <span class="text-sm font-semibold text-[#1E212B]">{{ formatCurrency(item.qty * item.rate) }}</span>
            </div>
          </div>

          <!-- Recurring schedule -->
          <div v-if="hasRecurring" class="mb-6 bg-[#F0F7FF] rounded-xl p-4">
            <p class="text-xs font-semibold text-[#2176AE] mb-2">Recurring Schedule</p>
            <div class="space-y-1">
              <p v-for="(amount, cycle) in recurringByBillingCycle" :key="cycle" class="text-xs text-[#2176AE]">
                {{ formatCurrency(amount) }} charged every {{ cycle }}<template v-if="invoiceData.recurringStartDate">, starting {{ invoiceData.recurringStartDate }}</template>
              </p>
            </div>
          </div>

          <!-- Pay button -->
          <button class="w-full bg-[#E87117] hover:bg-[#d06810] text-white font-semibold py-3.5 rounded-xl transition-colors text-sm">
            {{ hasRecurring && !hasOneTime
              ? `Subscribe — ${formatCurrency(Object.values(recurringByBillingCycle)[0] ?? 0)} / ${Object.keys(recurringByBillingCycle)[0] ?? 'month'}`
              : hasOneTime
                ? `Pay ${formatCurrency(oneTimeTotal)}`
                : `Pay ${formatCurrency(total)}`
            }}
          </button>

          <p class="text-center text-xs text-gray-400 mt-4 flex items-center justify-center gap-1">
            <svg class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="3" y="11" width="18" height="11" rx="2" /><path d="M7 11V7a5 5 0 0110 0v4" />
            </svg>
            Secured by Gro Payments
          </p>
        </div>
      </div>
    </div>

    <!-- ── Email Preview ───────────────────────────────────────── -->
    <div
      v-else-if="activeTab === 'email'"
      class="px-16 pb-10"
    >
      <div class="bg-white shadow-[0px_35.16px_70.32px_0px_#4E4D593D] rounded-lg overflow-hidden">
        <!-- Email client chrome -->
        <div class="bg-[#F6F6F7] px-5 py-3 border-b border-[#EDEDEE] space-y-0.5">
          <p class="text-xs text-gray-500">
            <span class="font-semibold text-[#1E212B]">From:</span>
            {{ invoiceData.businessName }} &lt;{{ invoiceData.email }}&gt;
          </p>
          <p class="text-xs text-gray-500">
            <span class="font-semibold text-[#1E212B]">To:</span>
            {{ billedToDisplay || 'customer@email.com' }}
          </p>
          <p class="text-xs text-gray-500">
            <span class="font-semibold text-[#1E212B]">Subject:</span>
            Invoice from {{ invoiceData.businessName }}{{ invoiceData.invoiceNumber ? ` — #${invoiceData.invoiceNumber}` : '' }}
          </p>
        </div>

        <!-- Email body -->
        <div class="bg-[#F6F8FA] px-5 py-6">
          <div class="max-w-sm mx-auto bg-white rounded-xl overflow-hidden shadow-sm">
            <!-- Brand band -->
            <div class="bg-[#E87117] px-6 py-5 text-center">
              <img
                src="https://res.cloudinary.com/sdq121/image/upload/v1764981178/ccrh2qinydoihgbswden.svg"
                class="w-10 h-10 mx-auto mb-2"
              >
              <h3 class="text-white font-bold text-base">{{ invoiceData.businessName }}</h3>
            </div>

            <!-- Content -->
            <div class="px-6 py-5">
              <p class="text-sm text-gray-700 mb-1">Hi {{ billedToDisplay || 'there' }},</p>
              <p class="text-sm text-gray-600 leading-relaxed mb-5">
                You have a new invoice from <span class="font-semibold text-[#1E212B]">{{ invoiceData.businessName }}</span>.
                <template v-if="hasRecurring">
                  <template v-for="(amount, cycle) in recurringByBillingCycle" :key="cycle">
                    Includes a recurring charge of <span class="font-semibold">{{ formatCurrency(amount) }} / {{ cycle }}</span>.
                  </template>
                </template>
              </p>

              <!-- Amount card -->
              <div class="bg-[#FAFAFA] rounded-xl p-4 mb-5 border border-[#EDEDEE] text-center">
                <p class="text-[10px] text-gray-400 uppercase tracking-wider mb-1">
                  {{ hasOneTime ? 'Amount due' : 'Per cycle' }}
                </p>
                <p class="text-2xl font-bold text-[#E87117]">
                  {{ formatCurrency(hasOneTime ? oneTimeTotal : Object.values(recurringByBillingCycle)[0] ?? 0) }}
                </p>
                <p v-if="hasOneTime && invoiceData.dueDate" class="text-xs text-gray-400 mt-1">
                  Due {{ invoiceData.dueDate }}
                </p>
                <p v-if="hasRecurring && invoiceData.recurringStartDate" class="text-xs text-gray-400 mt-1">
                  Starting {{ invoiceData.recurringStartDate }}
                </p>
              </div>

              <!-- Details -->
              <div class="space-y-2 mb-5 text-sm">
                <div v-if="invoiceData.invoiceNumber" class="flex justify-between">
                  <span class="text-gray-400">Invoice #</span>
                  <span class="font-medium text-[#1E212B]">{{ invoiceData.invoiceNumber }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-gray-400">Date issued</span>
                  <span class="font-medium text-[#1E212B]">{{ invoiceData.genDate || '---' }}</span>
                </div>
                <div v-if="hasItems" class="flex justify-between">
                  <span class="text-gray-400">Items</span>
                  <span class="font-medium text-[#1E212B]">{{ invoiceData.lineItems.length }}</span>
                </div>
              </div>

              <!-- Note to customer -->
              <div v-if="invoiceData.subject" class="bg-blue-50 rounded-lg px-4 py-3 mb-5">
                <p class="text-xs font-semibold text-[#2176AE] mb-1">Note</p>
                <p class="text-xs text-gray-600">{{ invoiceData.subject }}</p>
              </div>

              <!-- CTA -->
              <a href="#" class="block w-full bg-[#E87117] text-white text-sm font-semibold text-center py-3 rounded-xl hover:bg-[#d06810] transition-colors">
                {{ hasRecurring ? 'View Invoice & Subscribe' : 'View & Pay Invoice' }}
              </a>

              <p class="text-center text-[11px] text-gray-400 mt-4 leading-relaxed">
                Sent by {{ invoiceData.businessName }} via Gro.<br>
                Questions? Contact {{ invoiceData.email }}
              </p>
            </div>

            <!-- Footer -->
            <div class="bg-[#F6F6F7] px-6 py-3 text-center border-t border-[#EDEDEE]">
              <p class="text-[10px] text-gray-400">Powered by Gro · Unsubscribe</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import CreateInvoiceDashboard from '@/components/dashboard/page/invoice/create-invoice-dashboard.vue'
import GroBasicButton from '@/components/buttons/GroBasicButton.vue'
import CreateNewInvoiceForm from '@/components/page/invoice/create-new-invoice-form.vue'
import PreviewNewInvoice from '@/components/page/invoice/preview-new-invoice.vue'
import CreateCustomerModal from '@/components/modals/modal/create-customer-modal.vue'
import InvoiceTemplateSettings from '@/components/page/invoice/InvoiceTemplateSettings.vue'
import { useCustomerAPI } from '@/composables/api/customer/customer'
import { useProductAPI } from '@/composables/api/catalog/product'
import { useServiceAPI } from '@/composables/api/catalog/service'
import { useInvoiceAPI } from '@/composables/api/billing/invoice'
import type { InvoiceTemplate } from '@/constants/invoice/templates'

export interface CatalogVariant {
  id: string
  name: string
  sku: string
  price: number
}

export interface CatalogPlan {
  id: string
  name: string
  plan_type: 'subscription' | 'package'
  price: number
  price_currency: string
  billing_cycle?: string
  session_count?: number | null
  validity_days?: number | null
}

export interface CatalogEntry {
  id: string
  name: string
  type: 'product' | 'service'
  defaultRate: number
  imageUrl?: string
  variants?: CatalogVariant[]
  plans?: CatalogPlan[]
}

export interface LineItem {
  id: string
  catalogId: string
  name: string
  type: 'product' | 'service'
  qty: number
  rate: number
  imageUrl?: string
  discountPercent?: number
  billingType: 'one-time' | 'recurring'
  billingCycle?: string  // 'monthly' | 'yearly' | 'weekly' — only for recurring
}

// ── Template settings ─────────────────────────────────────────
const isTemplateSettingsOpen = ref(false)
const selectedTemplate = ref<InvoiceTemplate>('classic')

// ── Action state ──────────────────────────────────────────────
const isCreating = ref(false)
const isSending = ref(false)

// ── Data sources ─────────────────────────────────────────────
const isCustomerModalOpen = ref(false)
const crmCustomers = ref<{ label: string; value: string; email: string }[]>([])
const catalogProducts = ref<CatalogEntry[]>([])
const catalogServices = ref<CatalogEntry[]>([])

// ── Form state ────────────────────────────────────────────────
const selectedCustomers = ref<string[]>([])
const lineItems = ref<LineItem[]>([])
const dueDate = ref('')
const recurringStartDate = ref('')
const taxRate = ref(0)
const referenceNumber = ref('')
const termsAndCondition = ref('')
const notesToCustomer = ref('')
const memos = ref('')
const attachments = ref<File[]>([])

// ── Data fetching ─────────────────────────────────────────────
const loadCustomers = async () => {
  const result = await useCustomerAPI().FetchCustomers()
  if (!result.success) return
  // CommonResponse wraps paginated payload: result.data.data.data = FormSubmission[]
  const paginated = (result.data as any)?.data
  const submissions: any[] = paginated?.data ?? []
  crmCustomers.value = submissions.map((s: any) => {
    const answers = s.answers || {}
    const displayName = answers.customer_type === 'Individual'
      ? `${answers.first_name || ''} ${answers.last_name || ''}`.trim()
      : answers.company_name || answers.business_name || ''
    return {
      label: displayName || 'Unknown',
      value: s._id,
      email: answers.email || '',
    }
  })
}

const loadCatalog = async () => {
  const [productsRes, servicesRes] = await Promise.all([
    useProductAPI().ListProducts({ limit: 100 }),
    useServiceAPI().ListServices(),
  ])
  if (productsRes.success) {
    const data = (productsRes.data?.data as any)?.data ?? productsRes.data?.data ?? []
    catalogProducts.value = data.map((p: any) => ({
      id: p.id, name: p.name, type: 'product' as const, defaultRate: p.price || 0,
      imageUrl: p.media?.find((m: any) => m.display_image)?.url ?? p.media?.[0]?.url ?? undefined,
      variants: p.product_detail?.variants ?? [],
    }))
  }
  if (servicesRes.success) {
    const data = (servicesRes.data?.data as any) ?? []
    catalogServices.value = data.map((s: any) => ({
      id: s.id, name: s.name, type: 'service' as const, defaultRate: s.price || 0,
      imageUrl: s.media?.find((m: any) => m.display_image)?.url ?? s.media?.[0]?.url ?? undefined,
      plans: (s.plans ?? []).filter((pl: any) => pl.status === 'active'),
    }))
  }
}

// ── Invoice data for preview ──────────────────────────────────
const billedToDisplay = computed(() =>
  crmCustomers.value.filter(c => selectedCustomers.value.includes(c.value))
)

const invoiceData = computed(() => ({
  businessName: 'UseGro.inc',
  website: 'www.website.com',
  email: 'hello@email.com',
  phone: '+91 00000 00000',
  businessAddress: 'Business address',
  city: 'City',
  state: 'State',
  country: 'US',
  postalCode: '00000',
  billedTo: billedToDisplay.value,
  lineItems: lineItems.value,
  invoiceNumber: referenceNumber.value,
  termsAndConditions: termsAndCondition.value,
  subject: notesToCustomer.value,
  taxRate: taxRate.value,
  dueDate: dueDate.value,
  recurringStartDate: recurringStartDate.value,
  genDate: new Date().toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' }),
}))

// ── Invoice actions ───────────────────────────────────────────
const buildInvoicePayload = (status: 'draft' | 'sent') => {
  const selected = crmCustomers.value.filter(c => selectedCustomers.value.includes(c.value))
  return {
    customer_ids: selected.map(c => c.value),
    customer_names: selected.map(c => c.label),
    customer_emails: selected.map(c => c.email),
    line_items: lineItems.value.map(item => ({
      name: item.name,
      type: item.type,
      catalog_id: item.catalogId,
      qty: item.qty,
      rate: item.rate,
      amount: item.qty * item.rate,
      billing_type: item.billingType,
      billing_cycle: item.billingCycle || '',
      image_url: item.imageUrl || '',
    })),
    tax_rate: taxRate.value,
    subject: notesToCustomer.value,
    terms_and_conditions: termsAndCondition.value,
    reference_number: referenceNumber.value,
    memo: memos.value,
    due_date: dueDate.value || null,
    recurring_start_date: recurringStartDate.value || null,
    status,
  }
}

const saveDraft = async () => {
  if (isCreating.value) return
  isCreating.value = true
  try {
    const { CreateInvoice } = useInvoiceAPI()
    await CreateInvoice(buildInvoicePayload('draft'))
  } finally {
    isCreating.value = false
  }
}

const sendInvoice = async () => {
  if (isSending.value) return
  isSending.value = true
  try {
    const { CreateInvoice, SendInvoice } = useInvoiceAPI()
    const res = await CreateInvoice(buildInvoicePayload('draft'))
    const invoiceId = (res.data as any)?.data?.id
    if (invoiceId) {
      await SendInvoice(invoiceId)
    }
  } finally {
    isSending.value = false
  }
}

onMounted(async () => {
  await Promise.all([loadCustomers(), loadCatalog()])
})
</script>

<template>
  <div>
    <CreateInvoiceDashboard>
      <template #title>
        New Invoice
      </template>
      <template #menu>
        <div class="hidden md:flex items-center justify-start gap-x-8">
          <button
            class="text-[#2176AE] text-sm whitespace-nowrap hover:text-[#1a5f8a] transition-colors cursor-pointer"
            @click="isTemplateSettingsOpen = true"
          >
            Invoice Template Settings
          </button>
          <GroBasicButton
            class="whitespace-nowrap"
            color="tertiary"
          >
            Hide Preview
          </GroBasicButton>
          <GroBasicButton
            class="whitespace-nowrap"
            color="tertiary"
            :disabled="isCreating"
            @click="saveDraft"
          >
            {{ isCreating ? 'Saving...' : 'Save as Draft' }}
          </GroBasicButton>
          <GroBasicButton
            class="whitespace-nowrap"
            color="primary"
            :disabled="isSending"
            @click="sendInvoice"
          >
            {{ isSending ? 'Sending...' : 'Send invoice' }}
          </GroBasicButton>
        </div>
      </template>
      <template #body>
        <div class="flex h-screen">
          <CreateNewInvoiceForm
            v-model:selected-customers="selectedCustomers"
            v-model:line-items="lineItems"
            v-model:due-date="dueDate"
            v-model:recurring-start-date="recurringStartDate"
            v-model:tax-rate="taxRate"
            v-model:reference-number="referenceNumber"
            v-model:terms-and-condition="termsAndCondition"
            v-model:notes-to-customer="notesToCustomer"
            v-model:memos="memos"
            v-model:attachments="attachments"
            :crm-customers="crmCustomers"
            :catalog-products="catalogProducts"
            :catalog-services="catalogServices"
            @create-new-customer="isCustomerModalOpen = true"
          />
          <PreviewNewInvoice :invoice-data="invoiceData" :selected-template="selectedTemplate" />
          <InvoiceTemplateSettings v-model="isTemplateSettingsOpen" v-model:selected-template="selectedTemplate" />
          <CreateCustomerModal
            v-model="isCustomerModalOpen"
            @new-customer-added="loadCustomers"
          />
        </div>
      </template>
    </CreateInvoiceDashboard>
  </div>
</template>

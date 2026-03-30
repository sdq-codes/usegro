<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'nuxt/app'
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
import { notify } from '@/composables/helpers/notification/notification'
import type { InvoiceTemplate } from '@/constants/invoice/templates'
import type { CatalogEntry, LineItem } from '@/pages/invoices/create-new-invoice.vue'

const route = useRoute()
const router = useRouter()
const id = route.params.id as string

// ── Template settings ─────────────────────────────────────────
const isTemplateSettingsOpen = ref(false)
const selectedTemplate = ref<InvoiceTemplate>('classic')

// ── Action state ──────────────────────────────────────────────
const isSaving = ref(false)
const isLoading = ref(true)

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

const loadInvoice = async () => {
  const res = await useInvoiceAPI().GetInvoice(id)
  if (!res.success) {
    notify('Invoice not found', 'error')
    router.push('/invoices')
    return
  }
  const inv = (res.data as any)?.data
  if (!inv) {
    router.push('/invoices')
    return
  }
  if (inv.status !== 'draft') {
    notify('Only draft invoices can be edited', 'error')
    router.push(`/invoices/${id}`)
    return
  }

  // Pre-populate form state from existing invoice
  selectedCustomers.value = inv.customer_ids ?? []
  taxRate.value = inv.tax_rate ?? 0
  referenceNumber.value = inv.reference_number ?? ''
  termsAndCondition.value = inv.terms_and_conditions ?? ''
  notesToCustomer.value = inv.subject ?? ''
  memos.value = inv.memo ?? ''
  dueDate.value = inv.due_date ? inv.due_date.substring(0, 10) : ''
  recurringStartDate.value = inv.recurring_start_date ? inv.recurring_start_date.substring(0, 10) : ''

  lineItems.value = (inv.line_items ?? []).map((li: any) => ({
    id: li.id,
    catalogId: li.catalog_id ?? '',
    name: li.name,
    type: li.type,
    qty: li.qty,
    rate: li.rate,
    imageUrl: li.image_url || undefined,
    billingType: li.billing_type,
    billingCycle: li.billing_cycle || undefined,
  }))
}

// ── Preview data ──────────────────────────────────────────────
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

// ── Save ──────────────────────────────────────────────────────
const saveChanges = async () => {
  if (isSaving.value) return
  isSaving.value = true
  try {
    const selected = crmCustomers.value.filter(c => selectedCustomers.value.includes(c.value))
    const payload = {
      customer_ids: selected.map(c => c.value),
      customer_names: selected.map(c => c.label),
      customer_emails: selected.map(c => c.email),
      line_items: lineItems.value.map(item => ({
        name: item.name,
        type: item.type,
        catalog_id: item.catalogId,
        qty: item.qty,
        rate: item.rate,
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
    }

    const res = await useInvoiceAPI().UpdateInvoice(id, payload)
    if (res.success) {
      notify('Invoice saved', 'success')
      router.push(`/invoices/${id}`)
    } else {
      notify(res.error || 'Failed to save invoice', 'error')
    }
  } finally {
    isSaving.value = false
  }
}

onMounted(async () => {
  isLoading.value = true
  await Promise.all([loadCustomers(), loadCatalog(), loadInvoice()])
  isLoading.value = false
})
</script>

<template>
  <div>
    <CreateInvoiceDashboard>
      <template #title>
        Edit Invoice
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
            @click="router.push(`/invoices/${id}`)"
          >
            Cancel
          </GroBasicButton>
          <GroBasicButton
            class="whitespace-nowrap"
            color="primary"
            :disabled="isSaving || isLoading"
            @click="saveChanges"
          >
            {{ isSaving ? 'Saving...' : 'Save Changes' }}
          </GroBasicButton>
        </div>
      </template>
      <template #body>
        <div v-if="isLoading" class="flex items-center justify-center h-64">
          <div class="animate-pulse text-[#6F7177] text-sm">Loading invoice...</div>
        </div>
        <div v-else class="flex h-screen">
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

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRouter } from 'nuxt/app'
import MainDashboard from '@/components/dashboard/main-dashboard.vue'
import GroBasicButton from '@/components/buttons/GroBasicButton.vue'
import GroDataTable from '@/components/table/GroDataTable.vue'
import type { TableColumn, TableAction } from '@/components/table/GroDataTable.vue'
import GroFilterDropdown from '@/components/filters/GroFilterDropdown.vue'
import type { FilterDefinition } from '@/components/filters/GroFilterDropdown.vue'
import BasicModal from '@/components/modals/basic-modal.vue'
import { HugeiconsIcon } from '@hugeicons/vue'
import {
  PlusSignIcon, CheckmarkCircle02Icon, Invoice01Icon, InformationCircleIcon,
  Calendar01Icon, User03Icon, DollarCircleIcon, Invoice03Icon, Calendar03Icon, Tag01Icon,
} from '@hugeicons/core-free-icons'
import { useInvoiceAPI } from '@/composables/api/billing/invoice'
import type { Invoice } from '@/composables/dto/billing/invoice'
import { notify } from '@/composables/helpers/notification/notification'

// Display status is computed from backend status + due_date
type DisplayStatus = 'Overdue' | 'Not due yet' | 'Due today' | 'Paid' | 'Draft' | 'Cancelled' | 'Sent'

interface MappedInvoice extends Record<string, unknown> {
  id: string
  due_date_raw: string | null       // ISO — used by date-range filter
  created_at_raw: string            // ISO — used by date-range filter
  due_date: string                  // formatted for display
  display_status: DisplayStatus
  customer: string
  invoice_number: string
  invoice_number_display: string    // strip prefix — display only
  created_at: string
  billing_type: string
  billing_cycle: string
  total: number
}

const router = useRouter()
const isLoading = ref(false)
const isDeleting = ref(false)
const isDeleteModalOpen = ref(false)
const invoiceToDelete = ref<MappedInvoice | null>(null)

const PAGE_SIZE = 20
const currentPage = ref(1)
const totalPages = ref(1)
const totalCount = ref(0)
const hasAnyInvoices = ref<boolean | null>(null)
const invoices = ref<MappedInvoice[]>([])

const formatDate = (iso: string | null | undefined) => {
  if (!iso) return '—'
  return new Date(iso).toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' })
}

const formatCurrency = (amount: number) =>
  new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(amount)

// Derive a human-readable display status from backend status + due_date
const computeDisplayStatus = (inv: Invoice): DisplayStatus => {
  if (inv.status === 'paid') return 'Paid'
  if (inv.status === 'cancelled') return 'Cancelled'
  if (inv.status === 'draft') return 'Draft'
  // For sent / overdue: check due_date
  if (inv.due_date) {
    const due = new Date(inv.due_date)
    const today = new Date()
    due.setHours(0, 0, 0, 0)
    today.setHours(0, 0, 0, 0)
    if (due < today) return 'Overdue'
    if (due.getTime() === today.getTime()) return 'Due today'
    return 'Not due yet'
  }
  return inv.status === 'overdue' ? 'Overdue' : 'Sent'
}

// Get dominant billing type from line items
const deriveBillingType = (inv: Invoice): { label: string; cycle: string } => {
  if (!inv.line_items?.length) return { label: '—', cycle: '' }
  const recurring = inv.line_items.find(i => i.billing_type === 'recurring')
  if (recurring) {
    const cycleMap: Record<string, string> = {
      monthly: 'Monthly', weekly: 'Weekly', yearly: 'Yearly',
      daily: 'Daily', biweekly: 'Bi-Weekly',
    }
    return {
      label: cycleMap[recurring.billing_cycle?.toLowerCase()] ?? recurring.billing_cycle ?? 'Recurring',
      cycle: '',
    }
  }
  return { label: 'One-time', cycle: '' }
}

const mapInvoice = (inv: Invoice): MappedInvoice => {
  const { label, cycle } = deriveBillingType(inv)
  const numericSuffix = inv.invoice_number ? inv.invoice_number.split('-').pop() ?? inv.invoice_number : '—'
  return {
    id: inv.id,
    due_date_raw: inv.due_date,
    created_at_raw: inv.created_at,
    due_date: formatDate(inv.due_date),
    display_status: computeDisplayStatus(inv),
    customer: inv.customer_names?.join(', ') || '—',
    invoice_number: inv.invoice_number || '—',   // full value used for text-filter
    invoice_number_display: numericSuffix,        // short value shown in cell
    created_at: formatDate(inv.created_at),
    billing_type: label,
    billing_cycle: cycle,
    total: inv.total,
  }
}

const tableColumns: TableColumn[] = [
  { key: 'due_date', label: 'Due date' },
  { key: 'display_status', label: 'Status' },
  { key: 'customer', label: 'Customer' },
  { key: 'invoice_number_display', label: 'Invoice No.' },
  { key: 'created_at', label: 'Creation date' },
  { key: 'billing_type', label: 'Type' },
  { key: 'total', label: 'Amount' },
]

const tableActions: TableAction[] = [
  { key: 'view', label: 'View' },
  { key: 'delete', label: 'Delete', class: 'text-red-600' },
]

const tableFilters: FilterDefinition[] = [
  {
    key: 'due_date_range',
    label: 'Due Date',
    icon: Calendar01Icon,
    type: 'date-range',
    rawKey: 'due_date_raw',
  },
  {
    key: 'display_status',
    label: 'Status',
    icon: CheckmarkCircle02Icon,
    priority: true,
    type: 'options',
    options: [
      { label: 'All', value: '' },
      { label: 'Not due yet', value: 'Not due yet' },
      { label: 'Due today', value: 'Due today' },
      { label: 'Overdue', value: 'Overdue' },
      { label: 'Paid', value: 'Paid' },
      { label: 'Draft', value: 'Draft' },
      { label: 'Cancelled', value: 'Cancelled' },
    ],
  },
  {
    key: 'customer',
    label: 'Customers',
    icon: User03Icon,
    type: 'text',
    rawKey: 'customer',
  },
  {
    key: 'amount_range',
    label: 'Amount',
    icon: DollarCircleIcon,
    type: 'number-range',
    rawKey: 'total',
  },
  {
    key: 'invoice_number',
    label: 'Invoice Number',
    icon: Invoice03Icon,
    type: 'text',
    rawKey: 'invoice_number',
  },
  {
    key: 'created_at_range',
    label: 'Creation date',
    icon: Calendar03Icon,
    type: 'date-range',
    rawKey: 'created_at_raw',
  },
  {
    key: 'billing_type',
    label: 'Type',
    icon: Tag01Icon,
    type: 'options',
    options: [
      { label: 'All', value: '' },
      { label: 'One-time', value: 'One-time' },
      { label: 'Monthly', value: 'Monthly' },
      { label: 'Weekly', value: 'Weekly' },
      { label: 'Yearly', value: 'Yearly' },
      { label: 'Daily', value: 'Daily' },
      { label: 'Bi-Weekly', value: 'Bi-Weekly' },
    ],
  },
]

const STATUS_CLASSES: Record<DisplayStatus, string> = {
  'Overdue':     'bg-[#FFF3E0] text-[#E65100]',
  'Not due yet': 'bg-[#F3F4F6] text-[#6B7280]',
  'Due today':   'bg-[#F3F4F6] text-[#6B7280]',
  'Paid':        'bg-[#F0FDF4] text-[#15803D]',
  'Draft':       'bg-[#F9FAFB] text-[#9CA3AF]',
  'Sent':        'bg-[#EFF6FF] text-[#1D4ED8]',
  'Cancelled':   'bg-[#F9FAFB] text-[#9CA3AF]',
}

// ── Summary stats ─────────────────────────────────────────────
// We accumulate stats from all loaded invoices. For a full picture we keep
// a flat allInvoices list that is updated alongside the paged view.
const allInvoices = ref<MappedInvoice[]>([])

const stats = computed(() => {
  const outstanding = allInvoices.value.filter(
    i => i.display_status === 'Not due yet' || i.display_status === 'Due today' || i.display_status === 'Sent',
  )
  const overdue = allInvoices.value.filter(i => i.display_status === 'Overdue')
  const paid = allInvoices.value.filter(i => i.display_status === 'Paid')

  const sum = (list: MappedInvoice[]) =>
    list.reduce((acc, i) => acc + (Number(i.total) || 0), 0)

  return {
    outstanding: { amount: sum(outstanding), count: outstanding.length },
    overdue:     { amount: sum(overdue),     count: overdue.length },
    paid:        { amount: sum(paid),         count: paid.length },
  }
})

// ── Filter state (drives API params) ──────────────────────────
const filterValues = ref<Record<string, string>>(
  Object.fromEntries(tableFilters.map(f => [f.key, '']))
)

// Translate the filter model into API query params
const buildApiParams = (page: number) => {
  const p: Record<string, string | number> = { page, limit: PAGE_SIZE }

  // Status
  if (filterValues.value.display_status)
    p.status = filterValues.value.display_status

  // Customer name
  if (filterValues.value.customer)
    p.customer_name = filterValues.value.customer

  // Invoice number
  if (filterValues.value.invoice_number)
    p.invoice_number = filterValues.value.invoice_number

  // Billing type
  if (filterValues.value.billing_type)
    p.billing_type = filterValues.value.billing_type

  // Due date range
  if (filterValues.value.due_date_range) {
    try {
      const { from, to } = JSON.parse(filterValues.value.due_date_range)
      if (from) p.due_date_from = from
      if (to) p.due_date_to = to
    } catch {}
  }

  // Created-at range
  if (filterValues.value.created_at_range) {
    try {
      const { from, to } = JSON.parse(filterValues.value.created_at_range)
      if (from) p.created_from = from
      if (to) p.created_to = to
    } catch {}
  }

  // Amount range
  if (filterValues.value.amount_range) {
    try {
      const { min, max } = JSON.parse(filterValues.value.amount_range)
      if (min !== '') p.amount_min = Number(min)
      if (max !== '') p.amount_max = Number(max)
    } catch {}
  }

  return p
}

const fetchInvoices = async (page = 1) => {
  isLoading.value = true
  try {
    const params = buildApiParams(page)
    const result = await useInvoiceAPI().ListInvoices(params as any)
    if (!result.success) return
    const body = (result.data as any)?.data
    invoices.value = (body?.data ?? []).map(mapInvoice)
    totalCount.value = body?.total ?? 0
    totalPages.value = Math.ceil(totalCount.value / PAGE_SIZE) || 1
    currentPage.value = page

    // On initial load (no filters), fetch all for accurate stat cards
    if (hasAnyInvoices.value === null) {
      hasAnyInvoices.value = totalCount.value > 0
      if (totalCount.value > PAGE_SIZE) {
        const allRes = await useInvoiceAPI().ListInvoices({ page: 1, limit: totalCount.value })
        if (allRes.success) {
          allInvoices.value = ((allRes.data as any)?.data?.data ?? []).map(mapInvoice)
        }
      } else {
        allInvoices.value = invoices.value
      }
    }
  } finally {
    isLoading.value = false
  }
}

// Re-fetch from page 1 whenever any filter changes
watch(filterValues, () => {
  currentPage.value = 1
  fetchInvoices(1)
}, { deep: true })

onMounted(() => fetchInvoices(1))

const handleRowClick = (row: Record<string, unknown>) => {
  router.push(`/invoices/${row.id}`)
}

const handleAction = (key: string, row: Record<string, unknown>) => {
  const inv = row as MappedInvoice
  if (key === 'view') router.push(`/invoices/${inv.id}`)
  else if (key === 'delete') {
    invoiceToDelete.value = inv
    isDeleteModalOpen.value = true
  }
}

const confirmDelete = async () => {
  if (!invoiceToDelete.value) return
  isDeleting.value = true
  try {
    const res = await useInvoiceAPI().DeleteInvoice(invoiceToDelete.value.id)
    if (res.success) {
      notify('Invoice deleted', 'success')
      isDeleteModalOpen.value = false
      invoiceToDelete.value = null
      const targetPage = invoices.value.length === 1 && currentPage.value > 1
        ? currentPage.value - 1
        : currentPage.value
      await fetchInvoices(targetPage)
      hasAnyInvoices.value = totalCount.value > 0
    } else {
      notify('Failed to delete invoice', 'error')
    }
  } catch {
    notify('An error occurred', 'error')
  } finally {
    isDeleting.value = false
  }
}
</script>

<template>
  <div>
    <MainDashboard current="Invoices">
      <template #title>
        <div class="block md:flex justify-between">
          <div>
            <h6 class="text-2xl font-semibold">
              Invoices
              <span v-if="totalCount > 0" class="text-[#6F7177] font-normal">
                ({{ totalCount }})
              </span>
            </h6>
            <small v-if="hasAnyInvoices" class="text-[#6F7177]">
              Manage and track your invoices and payments.
            </small>
          </div>
          <div v-if="hasAnyInvoices" class="my-auto">
            <NuxtLink :to="{ name: 'invoices-create-new-invoice' }">
              <GroBasicButton color="primary" size="xs" shape="custom" class="w-max mt-2 md:mt-0">
                <template #frontIcon>
                  <HugeiconsIcon :icon="PlusSignIcon" :size="14" color="white" :stroke-width="2" />
                </template>
                New Invoice
              </GroBasicButton>
            </NuxtLink>
          </div>
        </div>
      </template>

      <template #body>
        <div>
          <!-- Empty state -->
          <div
            v-if="hasAnyInvoices === false"
            class="rounded-3xl bg-white px-8 py-12 mt-6"
          >
            <div class="flex items-center justify-center w-full">
              <div class="rounded-full bg-[#EDEDEE] h-48 w-48 flex items-center justify-center">
                <HugeiconsIcon :icon="Invoice01Icon" :size="64" color="#C0C2C8" />
              </div>
            </div>
            <h5 class="text-center mt-4 text-lg font-semibold">
              Collecting payments with Gro Invoices
            </h5>
            <div class="flex flex-col items-center gap-2 mt-2">
              <p class="text-[#6F7177] text-sm">Manage and track invoices all in one place</p>
              <p class="text-[#6F7177] text-sm">Get paid online right from the invoice in a click</p>
              <p class="text-[#6F7177] text-sm">Easily create and send invoices to your clients</p>
              <p class="text-[#6F7177] text-sm">Collect recurring payments for ongoing projects</p>
            </div>
            <div class="flex items-center justify-center w-full mt-6">
              <NuxtLink :to="{ name: 'invoices-create-new-invoice' }">
                <GroBasicButton color="primary" size="xs" shape="custom" class="w-max">
                  <template #frontIcon>
                    <HugeiconsIcon :icon="PlusSignIcon" :size="14" color="white" :stroke-width="2" />
                  </template>
                  Create Invoice
                </GroBasicButton>
              </NuxtLink>
            </div>
          </div>

          <!-- Invoice list -->
          <div v-else-if="hasAnyInvoices" class="mt-6">
            <!-- Summary stat cards -->
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
              <!-- Outstanding -->
              <div class="bg-white rounded-2xl px-6 py-5">
                <p class="text-xl font-semibold text-[#1E212B]">{{ formatCurrency(stats.outstanding.amount) }}</p>
                <div class="flex items-center gap-1.5 mt-1">
                  <span class="text-sm text-[#6F7177]">Outstanding</span>
                  <HugeiconsIcon :icon="InformationCircleIcon" :size="14" color="#9CA3AF" />
                </div>
                <p class="text-xs text-[#9CA3AF] mt-1">{{ stats.outstanding.count }} Invoice{{ stats.outstanding.count !== 1 ? 's' : '' }}</p>
              </div>

              <!-- Overdue -->
              <div class="bg-white rounded-2xl px-6 py-5">
                <p class="text-xl font-semibold text-[#E87117]">{{ formatCurrency(stats.overdue.amount) }}</p>
                <div class="flex items-center gap-1.5 mt-1">
                  <span class="text-sm text-[#6F7177]">Overdue Invoices</span>
                </div>
                <p class="text-xs text-[#9CA3AF] mt-1">{{ stats.overdue.count }} Invoice{{ stats.overdue.count !== 1 ? 's' : '' }}</p>
              </div>

              <!-- Paid -->
              <div class="bg-white rounded-2xl px-6 py-5">
                <p class="text-xl font-semibold text-[#1E212B]">{{ formatCurrency(stats.paid.amount) }}</p>
                <div class="flex items-center gap-1.5 mt-1">
                  <span class="text-sm text-[#6F7177]">Paid Invoices</span>
                  <HugeiconsIcon :icon="InformationCircleIcon" :size="14" color="#9CA3AF" />
                </div>
                <p class="text-xs text-[#9CA3AF] mt-1">{{ stats.paid.count }} Invoice{{ stats.paid.count !== 1 ? 's' : '' }}</p>
              </div>
            </div>
            <!-- Filter bar above the table — drives API params via watch -->
            <div class="mb-4">
              <GroFilterDropdown v-model="filterValues" :filters="tableFilters" />
            </div>

            <GroDataTable
              :columns="tableColumns"
              :rows="invoices"
              row-key="id"
              :actions="tableActions"
              :page="currentPage"
              :total-pages="totalPages"
              :total="totalCount"
              :is-loading="isLoading"
              empty-title="No invoices found"
              empty-message="Try adjusting your search or filters"
              @row-click="handleRowClick"
              @action="handleAction"
              @page-change="fetchInvoices"
            >
              <!-- Status badge -->
              <template #cell-display_status="{ row }">
                <span
                  class="inline-flex px-2.5 py-0.5 rounded-full text-xs font-medium"
                  :class="STATUS_CLASSES[row.display_status as DisplayStatus]"
                >
                  {{ row.display_status }}
                </span>
              </template>

              <!-- Invoice number — short numeric display -->
              <template #cell-invoice_number_display="{ row }">
                <span class="text-sm text-[#1E212B]">{{ row.invoice_number_display }}</span>
              </template>

              <!-- Type — label with cycle on second line -->
              <template #cell-billing_type="{ row }">
                <div class="leading-tight">
                  <p class="text-sm text-[#1E212B]">{{ row.billing_type }}</p>
                  <p v-if="row.billing_cycle" class="text-xs text-[#6F7177]">{{ row.billing_cycle }}</p>
                </div>
              </template>

              <!-- Amount -->
              <template #cell-total="{ row }">
                <span class="text-sm font-medium text-[#1E212B]">{{ formatCurrency(Number(row.total)) }}</span>
              </template>
            </GroDataTable>
          </div>

          <!-- Delete confirm modal -->
          <BasicModal v-model="isDeleteModalOpen" size="xs">
            <template #title>Delete Invoice</template>
            <template #default>
              <div class="py-4">
                <p class="text-gray-700">
                  Are you sure you want to delete invoice
                  <strong>{{ invoiceToDelete?.invoice_number }}</strong>?
                </p>
                <p class="text-red-500 font-bold text-sm mt-2">This action cannot be undone.</p>
              </div>
            </template>
            <template #footer>
              <div class="flex justify-end gap-x-4">
                <GroBasicButton
                  color="secondary" size="sm" shape="custom" class="w-max"
                  :disabled="isDeleting"
                  @click="isDeleteModalOpen = false"
                >
                  Cancel
                </GroBasicButton>
                <GroBasicButton
                  color="primary" size="sm" shape="custom"
                  class="w-max bg-red-600 hover:bg-red-700"
                  :disabled="isDeleting"
                  @click="confirmDelete"
                >
                  {{ isDeleting ? 'Deleting...' : 'Delete' }}
                </GroBasicButton>
              </div>
            </template>
          </BasicModal>
        </div>
      </template>
    </MainDashboard>
  </div>
</template>

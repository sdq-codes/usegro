<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'nuxt/app'
import MainDashboard from '@/components/dashboard/main-dashboard.vue'
import GroBasicButton from '@/components/buttons/GroBasicButton.vue'
import BasicModal from '@/components/modals/basic-modal.vue'
import { HugeiconsIcon } from '@hugeicons/vue'
import { ArrowLeft01Icon, Mail01Icon, Delete02Icon, CheckmarkCircle02Icon, PencilEdit01Icon } from '@hugeicons/core-free-icons'
import { useInvoiceAPI } from '@/composables/api/billing/invoice'
import type { Invoice, InvoiceStatus } from '@/composables/dto/billing/invoice'
import { notify } from '@/composables/helpers/notification/notification'

const route = useRoute()
const router = useRouter()
const id = route.params.id as string

const invoice = ref<Invoice | null>(null)
const isLoading = ref(true)
const isSending = ref(false)
const isDeleting = ref(false)
const isDeleteModalOpen = ref(false)

const STATUS_LABELS: Record<InvoiceStatus, string> = {
  draft: 'Draft',
  sent: 'Sent',
  paid: 'Paid',
  overdue: 'Overdue',
  cancelled: 'Cancelled',
}

const STATUS_CLASSES: Record<InvoiceStatus, string> = {
  draft: 'bg-[#F9FAFB] text-[#6F7177]',
  sent: 'bg-[#EFF6FF] text-[#1D4ED8]',
  paid: 'bg-[#F0FDF4] text-[#15803D]',
  overdue: 'bg-[#FEF2F2] text-[#DC2626]',
  cancelled: 'bg-[#F9FAFB] text-[#9CA3AF]',
}

const formatDate = (iso: string | null | undefined) => {
  if (!iso) return '—'
  return new Date(iso).toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' })
}

const formatCurrency = (amount: number) =>
  new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(amount)

const subtotal = computed(() =>
  invoice.value?.line_items?.reduce((sum, item) => sum + item.amount, 0) ?? 0
)
const taxAmount = computed(() =>
  subtotal.value * ((invoice.value?.tax_rate ?? 0) / 100)
)
const total = computed(() => subtotal.value + taxAmount.value)

const fetchInvoice = async () => {
  isLoading.value = true
  try {
    const res = await useInvoiceAPI().GetInvoice(id)
    if (res.success) {
      invoice.value = (res.data as any)?.data ?? null
    } else {
      notify('Invoice not found', 'error')
      router.push('/invoices')
    }
  } finally {
    isLoading.value = false
  }
}

const sendInvoice = async () => {
  if (!invoice.value || isSending.value) return
  isSending.value = true
  try {
    const res = await useInvoiceAPI().SendInvoice(id)
    if (res.success) {
      invoice.value = (res.data as any)?.data ?? invoice.value
      notify('Invoice sent successfully', 'success')
    } else {
      notify('Failed to send invoice', 'error')
    }
  } finally {
    isSending.value = false
  }
}

const confirmDelete = async () => {
  isDeleting.value = true
  try {
    const res = await useInvoiceAPI().DeleteInvoice(id)
    if (res.success) {
      notify('Invoice deleted', 'success')
      router.push('/invoices')
    } else {
      notify('Failed to delete invoice', 'error')
    }
  } finally {
    isDeleting.value = false
  }
}

onMounted(fetchInvoice)
</script>

<template>
  <div>
    <MainDashboard current="Invoices">
      <template #title>
        <div class="flex items-center gap-3">
          <button
            class="flex items-center justify-center w-8 h-8 rounded-full hover:bg-[#EDEDEE] transition-colors cursor-pointer"
            @click="router.push('/invoices')"
          >
            <HugeiconsIcon :icon="ArrowLeft01Icon" :size="18" color="#1E212B" />
          </button>
          <div>
            <h6 class="text-2xl font-semibold">
              {{ invoice?.invoice_number ?? '—' }}
            </h6>
            <small class="text-[#6F7177]">Invoice details</small>
          </div>
        </div>
      </template>

      <template #body>
        <!-- Loading skeleton -->
        <div v-if="isLoading" class="mt-6 space-y-4">
          <div class="rounded-2xl bg-white p-6 animate-pulse h-32" />
          <div class="rounded-2xl bg-white p-6 animate-pulse h-48" />
        </div>

        <div v-else-if="invoice" class="mt-6 space-y-4">
          <!-- Header card -->
          <div class="rounded-2xl bg-white p-6">
            <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
              <div class="flex items-center gap-4">
                <div>
                  <div class="flex items-center gap-3">
                    <h5 class="text-xl font-semibold text-[#1E212B]">{{ invoice.invoice_number }}</h5>
                    <span
                      class="inline-flex px-2.5 py-0.5 rounded-full text-xs font-medium"
                      :class="STATUS_CLASSES[invoice.status]"
                    >
                      {{ STATUS_LABELS[invoice.status] }}
                    </span>
                  </div>
                  <p class="text-sm text-[#6F7177] mt-1">
                    Created {{ formatDate(invoice.created_at) }}
                    <template v-if="invoice.sent_at"> · Sent {{ formatDate(invoice.sent_at) }}</template>
                  </p>
                </div>
              </div>
              <div class="flex items-center gap-3">
                <GroBasicButton
                  v-if="invoice.status === 'draft'"
                  color="secondary"
                  size="sm"
                  shape="custom"
                  class="w-max"
                  @click="router.push(`/invoices/${id}/edit`)"
                >
                  <template #frontIcon>
                    <HugeiconsIcon :icon="PencilEdit01Icon" :size="14" color="#6F7177" />
                  </template>
                  Edit
                </GroBasicButton>
                <GroBasicButton
                  v-if="invoice.status === 'draft'"
                  color="primary"
                  size="sm"
                  shape="custom"
                  class="w-max"
                  :disabled="isSending"
                  @click="sendInvoice"
                >
                  <template #frontIcon>
                    <HugeiconsIcon :icon="Mail01Icon" :size="14" color="white" />
                  </template>
                  {{ isSending ? 'Sending...' : 'Send Invoice' }}
                </GroBasicButton>
                <GroBasicButton
                  v-if="invoice.status === 'sent'"
                  color="primary"
                  size="sm"
                  shape="custom"
                  class="w-max"
                  :disabled="isSending"
                  @click="sendInvoice"
                >
                  <template #frontIcon>
                    <HugeiconsIcon :icon="Mail01Icon" :size="14" color="white" />
                  </template>
                  {{ isSending ? 'Resending...' : 'Resend Invoice' }}
                </GroBasicButton>
                <GroBasicButton
                  color="secondary"
                  size="sm"
                  shape="custom"
                  class="w-max"
                  @click="isDeleteModalOpen = true"
                >
                  <template #frontIcon>
                    <HugeiconsIcon :icon="Delete02Icon" :size="14" color="#6F7177" />
                  </template>
                  Delete
                </GroBasicButton>
              </div>
            </div>
          </div>

          <!-- Two-column layout -->
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <!-- Left: main invoice info -->
            <div class="md:col-span-2 space-y-4">
              <!-- Line items -->
              <div class="rounded-2xl bg-white p-6">
                <h6 class="text-sm font-semibold text-[#1E212B] mb-4">Line Items</h6>
                <div class="overflow-x-auto">
                  <table class="w-full text-sm">
                    <thead>
                      <tr class="border-b border-[#EDEDEE]">
                        <th class="text-left font-medium text-[#6F7177] pb-3 pr-4">Item</th>
                        <th class="text-left font-medium text-[#6F7177] pb-3 pr-4">Type</th>
                        <th class="text-right font-medium text-[#6F7177] pb-3 pr-4">Qty</th>
                        <th class="text-right font-medium text-[#6F7177] pb-3 pr-4">Rate</th>
                        <th class="text-right font-medium text-[#6F7177] pb-3">Amount</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr
                        v-for="item in invoice.line_items"
                        :key="item.id"
                        class="border-b border-[#EDEDEE] last:border-0"
                      >
                        <td class="py-3 pr-4 font-medium text-[#1E212B]">{{ item.name }}</td>
                        <td class="py-3 pr-4 text-[#6F7177] capitalize">{{ item.type }}</td>
                        <td class="py-3 pr-4 text-right text-[#1E212B]">{{ item.qty }}</td>
                        <td class="py-3 pr-4 text-right text-[#1E212B]">{{ formatCurrency(item.rate) }}</td>
                        <td class="py-3 text-right font-medium text-[#1E212B]">{{ formatCurrency(item.amount) }}</td>
                      </tr>
                    </tbody>
                  </table>
                </div>

                <!-- Totals -->
                <div class="mt-4 border-t border-[#EDEDEE] pt-4 space-y-2">
                  <div class="flex justify-between text-sm text-[#6F7177]">
                    <span>Subtotal</span>
                    <span>{{ formatCurrency(subtotal) }}</span>
                  </div>
                  <div v-if="invoice.tax_rate > 0" class="flex justify-between text-sm text-[#6F7177]">
                    <span>Tax ({{ invoice.tax_rate }}%)</span>
                    <span>{{ formatCurrency(taxAmount) }}</span>
                  </div>
                  <div class="flex justify-between text-base font-semibold text-[#1E212B] pt-2 border-t border-[#EDEDEE]">
                    <span>Total</span>
                    <span class="text-[#E87117]">{{ formatCurrency(total) }}</span>
                  </div>
                </div>
              </div>

              <!-- Notes / Terms -->
              <div v-if="invoice.subject || invoice.terms_and_conditions || invoice.memo" class="rounded-2xl bg-white p-6 space-y-4">
                <div v-if="invoice.subject">
                  <p class="text-xs font-semibold text-[#6F7177] uppercase tracking-wide mb-1">Subject / Notes</p>
                  <p class="text-sm text-[#1E212B]">{{ invoice.subject }}</p>
                </div>
                <div v-if="invoice.terms_and_conditions">
                  <p class="text-xs font-semibold text-[#6F7177] uppercase tracking-wide mb-1">Terms &amp; Conditions</p>
                  <p class="text-sm text-[#1E212B] whitespace-pre-line">{{ invoice.terms_and_conditions }}</p>
                </div>
                <div v-if="invoice.memo">
                  <p class="text-xs font-semibold text-[#6F7177] uppercase tracking-wide mb-1">Memo</p>
                  <p class="text-sm text-[#1E212B]">{{ invoice.memo }}</p>
                </div>
              </div>
            </div>

            <!-- Right: meta -->
            <div class="space-y-4">
              <!-- Billed to -->
              <div class="rounded-2xl bg-white p-6">
                <h6 class="text-xs font-semibold text-[#6F7177] uppercase tracking-wide mb-3">Billed To</h6>
                <div class="space-y-3">
                  <div
                    v-for="(name, i) in invoice.customer_names"
                    :key="i"
                    class="flex items-start gap-3"
                  >
                    <div class="w-8 h-8 rounded-full bg-[#EDEDEE] flex items-center justify-center text-xs font-semibold text-[#4B4D55] shrink-0">
                      {{ name.charAt(0).toUpperCase() }}
                    </div>
                    <div>
                      <p class="text-sm font-medium text-[#1E212B]">{{ name }}</p>
                      <p class="text-xs text-[#6F7177]">{{ invoice.customer_emails?.[i] || '—' }}</p>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Invoice details -->
              <div class="rounded-2xl bg-white p-6 space-y-3">
                <h6 class="text-xs font-semibold text-[#6F7177] uppercase tracking-wide mb-3">Details</h6>
                <div class="flex justify-between text-sm">
                  <span class="text-[#6F7177]">Invoice #</span>
                  <span class="font-medium text-[#1E212B]">{{ invoice.invoice_number }}</span>
                </div>
                <div v-if="invoice.reference_number" class="flex justify-between text-sm">
                  <span class="text-[#6F7177]">Reference</span>
                  <span class="font-medium text-[#1E212B]">{{ invoice.reference_number }}</span>
                </div>
                <div class="flex justify-between text-sm">
                  <span class="text-[#6F7177]">Status</span>
                  <span
                    class="inline-flex px-2 py-0.5 rounded-full text-xs font-medium"
                    :class="STATUS_CLASSES[invoice.status]"
                  >
                    {{ STATUS_LABELS[invoice.status] }}
                  </span>
                </div>
                <div v-if="invoice.due_date" class="flex justify-between text-sm">
                  <span class="text-[#6F7177]">Due Date</span>
                  <span class="font-medium text-[#1E212B]">{{ formatDate(invoice.due_date) }}</span>
                </div>
                <div v-if="invoice.recurring_start_date" class="flex justify-between text-sm">
                  <span class="text-[#6F7177]">Recurring Start</span>
                  <span class="font-medium text-[#1E212B]">{{ formatDate(invoice.recurring_start_date) }}</span>
                </div>
                <div class="flex justify-between text-sm">
                  <span class="text-[#6F7177]">Created</span>
                  <span class="font-medium text-[#1E212B]">{{ formatDate(invoice.created_at) }}</span>
                </div>
                <div v-if="invoice.sent_at" class="flex justify-between text-sm">
                  <span class="text-[#6F7177]">Sent</span>
                  <span class="font-medium text-[#1E212B]">{{ formatDate(invoice.sent_at) }}</span>
                </div>
              </div>

              <!-- Paid indicator -->
              <div v-if="invoice.status === 'paid'" class="rounded-2xl bg-[#F0FDF4] p-4 flex items-center gap-3">
                <HugeiconsIcon :icon="CheckmarkCircle02Icon" :size="20" color="#15803D" />
                <p class="text-sm font-medium text-[#15803D]">This invoice has been paid</p>
              </div>
            </div>
          </div>
        </div>
      </template>
    </MainDashboard>

    <!-- Delete confirm modal -->
    <BasicModal v-model="isDeleteModalOpen" size="xs">
      <template #title>Delete Invoice</template>
      <template #default>
        <div class="py-4">
          <p class="text-gray-700">
            Are you sure you want to delete invoice <strong>{{ invoice?.invoice_number }}</strong>?
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

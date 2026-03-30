<script setup lang="ts">
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue'
import {
  MinusSignIcon,
  PlusSignIcon,
  ShoppingCart01Icon,
  UserMultipleIcon,
  Search01Icon,
} from '@hugeicons/core-free-icons'
import GroBasicTextArea from '@/components/forms/input/GroBasicTextArea.vue'
import { HugeiconsIcon } from '@hugeicons/vue'
import GroBasicInput from '@/components/forms/input/GroBasicInput.vue'
import GroFileUpload from '@/components/forms/GroFileUpload.vue'
import type { CatalogEntry, LineItem } from '@/pages/invoices/create-new-invoice.vue'
import CreateInvoiceItemModal from '@/components/modals/modal/create-invoice-item-modal.vue'
import CatalogItemPickerModal from '@/components/modals/modal/catalog-item-picker-modal.vue'

// ── Models ────────────────────────────────────────────────────
const selectedCustomers = defineModel<string[]>('selectedCustomers', { default: () => [] })
const lineItems = defineModel<LineItem[]>('lineItems', { default: () => [] })
const dueDate = defineModel<string>('dueDate', { default: '' })
const recurringStartDate = defineModel<string>('recurringStartDate', { default: '' })
const taxRate = defineModel<number>('taxRate', { default: 0 })
const referenceNumber = defineModel<string>('referenceNumber', { default: '' })
const termsAndCondition = defineModel<string>('termsAndCondition', { default: '' })
const notesToCustomer = defineModel<string>('notesToCustomer', { default: '' })
const memos = defineModel<string>('memos', { default: '' })
const attachments = defineModel<File[]>('attachments', { default: () => [] })

// ── Props ─────────────────────────────────────────────────────
const props = defineProps<{
  crmCustomers: { label: string; value: string }[]
  catalogProducts: CatalogEntry[]
  catalogServices: CatalogEntry[]
}>()

const emit = defineEmits<{ (e: 'createNewCustomer'): void }>()

// ── Variant / plan picker ─────────────────────────────────────
const showPickerModal = ref(false)
const pickerEntry = ref<CatalogEntry | null>(null)

// ── New item modal ────────────────────────────────────────────
const showNewItemModal = ref(false)
const newItemType = ref<'product' | 'service'>('product')

const openNewItem = (type: 'product' | 'service') => {
  newItemType.value = type
  showItemsDropdown.value = false
  showNewItemModal.value = true
}

const onNewItemAdded = (entry: CatalogEntry & { discountPercent?: number }) => {
  pushLineItem({ catalogId: entry.id, name: entry.name, type: entry.type, qty: 1, rate: entry.defaultRate, imageUrl: entry.imageUrl, discountPercent: entry.discountPercent, billingType: 'one-time' })
}

// ── Customer picker ───────────────────────────────────────────
const showCustomerDropdown = ref(false)
const customerSearch = ref('')
const customerDropdownRef = ref<HTMLElement | null>(null)

const filteredCustomers = computed(() => {
  const q = customerSearch.value.toLowerCase()
  return props.crmCustomers.filter(c =>
    c.label.toLowerCase().includes(q) && !selectedCustomers.value.includes(c.value)
  )
})

const selectCustomer = (c: { label: string; value: string }) => {
  selectedCustomers.value = [...selectedCustomers.value, c.value]
  customerSearch.value = ''
  showCustomerDropdown.value = false
}

const removeCustomer = (value: string) => {
  selectedCustomers.value = selectedCustomers.value.filter(v => v !== value)
}

const selectedCustomerObjects = computed(() =>
  selectedCustomers.value
    .map(v => props.crmCustomers.find(c => c.value === v))
    .filter(Boolean) as { label: string; value: string }[]
)

const getInitials = (label: string) =>
  label.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2)

const onCustomerClickOutside = (e: MouseEvent) => {
  if (customerDropdownRef.value && !customerDropdownRef.value.contains(e.target as Node)) {
    showCustomerDropdown.value = false
  }
}

// ── Items picker ──────────────────────────────────────────────
const showItemsDropdown = ref(false)
const itemsSearch = ref('')
const productsExpanded = ref(false)
const servicesExpanded = ref(false)
const itemsDropdownRef = ref<HTMLElement | null>(null)

const filteredProducts = computed(() => {
  const q = itemsSearch.value.toLowerCase()
  return props.catalogProducts.filter(p => p.name.toLowerCase().includes(q))
})

const filteredServices = computed(() => {
  const q = itemsSearch.value.toLowerCase()
  return props.catalogServices.filter(s => s.name.toLowerCase().includes(q))
})

watch(itemsSearch, (q) => {
  const hasQuery = q.trim().length > 0
  productsExpanded.value = hasQuery
  servicesExpanded.value = hasQuery
})

const isItemSelected = (catalogId: string) =>
  lineItems.value.some(l => l.catalogId === catalogId)

const pushLineItem = (item: Omit<LineItem, 'id'>) => {
  lineItems.value = [
    ...lineItems.value,
    { id: `${Date.now()}-${Math.random().toString(36).slice(2)}`, ...item },
  ]
}

const addLineItem = (entry: CatalogEntry) => {
  if (isItemSelected(entry.id)) return
  const hasVariants = entry.type === 'product' && (entry.variants?.length ?? 0) > 0
  const hasPlans = entry.type === 'service' && (entry.plans?.length ?? 0) > 0
  if (hasVariants || hasPlans) {
    pickerEntry.value = entry
    showPickerModal.value = true
    showItemsDropdown.value = false
    itemsSearch.value = ''
    return
  }
  pushLineItem({ catalogId: entry.id, name: entry.name, type: entry.type, qty: 1, rate: entry.defaultRate, imageUrl: entry.imageUrl, billingType: 'one-time' })
  showItemsDropdown.value = false
  itemsSearch.value = ''
}

const removeLineItem = (id: string) => {
  if (editingRateId.value === id) editingRateId.value = null
  lineItems.value = lineItems.value.filter(l => l.id !== id)
}

const editingRateId = ref<string | null>(null)

const itemAmount = (item: LineItem) => item.qty * item.rate

const subtotal = computed(() => lineItems.value.reduce((sum, i) => sum + itemAmount(i), 0))
const taxAmount = computed(() => (subtotal.value * (taxRate.value || 0)) / 100)
const total = computed(() => subtotal.value + taxAmount.value)

const oneTimeItems = computed(() => lineItems.value.filter(i => i.billingType === 'one-time'))
const recurringItems = computed(() => lineItems.value.filter(i => i.billingType === 'recurring'))

const oneTimeSubtotal = computed(() => oneTimeItems.value.reduce((sum, i) => sum + i.qty * i.rate, 0))
const oneTimeTax = computed(() => (oneTimeSubtotal.value * (taxRate.value || 0)) / 100)
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

// Close dropdown on outside click
const onClickOutside = (e: MouseEvent) => {
  if (itemsDropdownRef.value && !itemsDropdownRef.value.contains(e.target as Node)) {
    showItemsDropdown.value = false
  }
}
onMounted(() => {
  document.addEventListener('mousedown', onClickOutside)
  document.addEventListener('mousedown', onCustomerClickOutside)
})
onBeforeUnmount(() => {
  document.removeEventListener('mousedown', onClickOutside)
  document.removeEventListener('mousedown', onCustomerClickOutside)
})

// ── Collapsible extras ────────────────────────────────────────
const showReferenceNumber = ref(false)
const showTermsAndCondition = ref(false)
const showMemos = ref(false)


</script>

<template>
  <div class="w-1/2 overflow-y-auto bg-white py-7 rounded-l-2xl border-y border-l border-[#EDEDEE]">

    <!-- ── Who are you billing? ─────────────────────────────── -->
    <div class="py-8 px-8 border-b border-[#EDEDEE]">
      <div class="flex items-start gap-x-2 mb-4">
        <HugeiconsIcon :icon="UserMultipleIcon" color="#4D91BE" class="mt-0.5 shrink-0" />
        <div>
          <h6 class="text-sm font-semibold">Who are you billing?</h6>
          <p class="text-xs text-[#939499] mt-0.5">You can bill multiple customers at once.</p>
        </div>
      </div>

      <!-- Search -->
      <div
        ref="customerDropdownRef"
        class="relative"
      >
        <div
          class="flex items-center gap-2 bg-[#F6F6F7] border border-[#EDEDEE] hover:border-[#94BDD8] rounded-lg px-3 py-2.5 cursor-text transition-colors"
          @click="showCustomerDropdown = true"
        >
          <HugeiconsIcon :icon="Search01Icon" :size="16" color="#939499" />
          <input
            v-model="customerSearch"
            class="flex-1 bg-transparent outline-none text-sm text-[#1E212B] placeholder-[#939499]"
            placeholder="Search customers…"
            @focus="showCustomerDropdown = true"
          >
        </div>

        <!-- Dropdown -->
        <div
          v-if="showCustomerDropdown"
          class="absolute top-full left-0 right-0 mt-1 bg-white border border-[#EDEDEE] rounded-xl shadow-lg z-30 overflow-hidden"
        >
          <div class="max-h-52 overflow-y-auto">
            <div
              v-if="filteredCustomers.length === 0"
              class="px-4 py-3 text-xs text-[#939499] italic"
            >
              No customers found
            </div>
            <button
              v-for="c in filteredCustomers"
              :key="c.value"
              class="w-full flex items-center gap-3 px-4 py-2.5 hover:bg-[#F6F6F7] transition-colors text-left"
              @click.prevent="selectCustomer(c)"
            >
              <div class="w-7 h-7 rounded-full bg-[#2176AE] flex items-center justify-center shrink-0">
                <span class="text-[10px] font-bold text-white">{{ getInitials(c.label) }}</span>
              </div>
              <span class="text-sm font-medium text-[#1E212B]">{{ c.label }}</span>
            </button>
          </div>
          <button
            class="w-full flex items-center gap-2 px-4 py-2.5 text-xs font-semibold text-[#2176AE] hover:bg-[#F0F7FF] transition-colors border-t border-[#EDEDEE]"
            @click.prevent="emit('createNewCustomer'); showCustomerDropdown = false"
          >
            <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M12 5v14M5 12h14" /></svg>
            New customer
          </button>
        </div>
      </div>

      <!-- Selected customer cards -->
      <div
        v-if="selectedCustomerObjects.length > 0"
        class="mt-3 border border-[#EDEDEE] rounded-xl overflow-hidden divide-y divide-[#EDEDEE]"
      >
        <div
          v-for="c in selectedCustomerObjects"
          :key="c.value"
          class="flex items-center gap-3 px-4 py-3"
        >
          <div class="w-9 h-9 rounded-full bg-[#2176AE] flex items-center justify-center shrink-0">
            <span class="text-xs font-bold text-white">{{ getInitials(c.label) }}</span>
          </div>
          <span class="flex-1 text-sm font-medium text-[#1E212B]">{{ c.label }}</span>
          <button
            class="w-6 h-6 flex items-center justify-center rounded text-[#939499] hover:text-[#AF513A] hover:bg-[#FEE2E2] transition-colors"
            @click="removeCustomer(c.value)"
          >
            <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M6 18L18 6M6 6l12 12" /></svg>
          </button>
        </div>
      </div>
    </div>

    <!-- ── What are they paying for? ────────────────────────── -->
    <div class="py-8 px-8 border-b border-[#EDEDEE]">
      <div class="flex items-start gap-x-2 mb-4">
        <HugeiconsIcon :icon="ShoppingCart01Icon" color="#D26B06" class="mt-0.5 shrink-0" />
        <div>
          <h6 class="text-sm font-semibold">What are they paying for?</h6>
          <p class="text-xs text-[#939499] mt-0.5">You can add as many products and services as you need.</p>
        </div>
      </div>

      <!-- Search -->
      <div
        ref="itemsDropdownRef"
        class="relative"
      >
        <div
          class="flex items-center gap-2 bg-[#F6F6F7] border border-[#EDEDEE] hover:border-[#94BDD8] rounded-lg px-3 py-2.5 cursor-text transition-colors"
          @click="showItemsDropdown = true"
        >
          <HugeiconsIcon
            :icon="Search01Icon"
            :size="16"
            color="#939499"
          />
          <input
            v-model="itemsSearch"
            class="flex-1 bg-transparent outline-none text-sm text-[#1E212B] placeholder-[#939499]"
            placeholder="Search products & services…"
            @focus="showItemsDropdown = true"
          >
        </div>

        <!-- Dropdown panel -->
        <div
          v-if="showItemsDropdown"
          class="absolute top-full left-0 right-0 mt-1 bg-white border border-[#EDEDEE] rounded-xl shadow-lg z-30 max-h-72 overflow-y-auto"
        >
          <!-- Products section -->
          <div>
            <button
              class="w-full flex items-center justify-between px-4 py-2.5 bg-[#F6F6F7] text-xs font-semibold text-[#6F7177] uppercase tracking-wider hover:bg-[#EDEDEE] transition-colors"
              @click.prevent="productsExpanded = !productsExpanded"
            >
              <span>Products ({{ filteredProducts.length }})</span>
              <svg
                class="w-3.5 h-3.5 transition-transform"
                :class="productsExpanded ? 'rotate-180' : ''"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2.5"
              >
                <path d="M6 9l6 6 6-6" />
              </svg>
            </button>
            <template v-if="productsExpanded">
              <div
                v-if="filteredProducts.length === 0"
                class="px-4 py-3 text-xs text-[#939499] italic"
              >
                No products found
              </div>
              <button
                v-for="p in filteredProducts"
                :key="p.id"
                class="w-full flex items-center gap-3 px-4 py-2 text-sm text-left transition-colors"
                :class="isItemSelected(p.id)
                  ? 'bg-[#F0F7FF] text-[#2176AE] cursor-default'
                  : 'hover:bg-[#F6F6F7] text-[#1E212B]'"
                @click.prevent="addLineItem(p)"
              >
                <div class="w-8 h-8 rounded-lg overflow-hidden bg-[#F6F6F7] border border-[#EDEDEE] shrink-0 flex items-center justify-center">
                  <img
                    v-if="p.imageUrl"
                    :src="p.imageUrl"
                    :alt="p.name"
                    class="w-full h-full object-cover"
                  >
                  <svg v-else class="w-4 h-4 text-[#EDEDEE]" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="3" y="3" width="18" height="18" rx="2" /><path d="M3 9l4-4 4 4 4-5 4 5" /><circle cx="8.5" cy="8.5" r="1.5" /></svg>
                </div>
                <span class="flex-1 font-medium truncate">{{ p.name }}</span>
                <span class="text-xs text-[#939499] shrink-0">{{ formatCurrency(p.defaultRate) }}</span>
              </button>
            </template>
          </div>

          <!-- Services section -->
          <div class="border-t border-[#EDEDEE]">
            <button
              class="w-full flex items-center justify-between px-4 py-2.5 bg-[#F6F6F7] text-xs font-semibold text-[#6F7177] uppercase tracking-wider hover:bg-[#EDEDEE] transition-colors"
              @click.prevent="servicesExpanded = !servicesExpanded"
            >
              <span>Services ({{ filteredServices.length }})</span>
              <svg
                class="w-3.5 h-3.5 transition-transform"
                :class="servicesExpanded ? 'rotate-180' : ''"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2.5"
              >
                <path d="M6 9l6 6 6-6" />
              </svg>
            </button>
            <template v-if="servicesExpanded">
              <div
                v-if="filteredServices.length === 0"
                class="px-4 py-3 text-xs text-[#939499] italic"
              >
                No services found
              </div>
              <button
                v-for="s in filteredServices"
                :key="s.id"
                class="w-full flex items-center gap-3 px-4 py-2 text-sm text-left transition-colors"
                :class="isItemSelected(s.id)
                  ? 'bg-[#F0F7FF] text-[#2176AE] cursor-default'
                  : 'hover:bg-[#F6F6F7] text-[#1E212B]'"
                @click.prevent="addLineItem(s)"
              >
                <div class="w-8 h-8 rounded-lg overflow-hidden bg-[#F6F6F7] border border-[#EDEDEE] shrink-0 flex items-center justify-center">
                  <img
                    v-if="s.imageUrl"
                    :src="s.imageUrl"
                    :alt="s.name"
                    class="w-full h-full object-cover"
                  >
                  <svg v-else class="w-4 h-4 text-[#EDEDEE]" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="3" y="3" width="18" height="18" rx="2" /><path d="M3 9l4-4 4 4 4-5 4 5" /><circle cx="8.5" cy="8.5" r="1.5" /></svg>
                </div>
                <span class="flex-1 font-medium truncate">{{ s.name }}</span>
                <span class="text-xs text-[#939499] shrink-0">{{ formatCurrency(s.defaultRate) }}</span>
              </button>
            </template>
          </div>

          <!-- New product / service — fixed footer inside dropdown -->
          <div class="flex items-center gap-1 px-4 py-2.5 border-t border-[#EDEDEE] bg-white sticky bottom-0">
            <button
              class="flex items-center gap-1.5 text-xs font-semibold text-[#D26B06] hover:text-[#b85c00] transition-colors"
              @click.prevent="openNewItem('product')"
            >
              <svg class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                <path d="M12 5v14M5 12h14" />
              </svg>
              New product
            </button>
            <span class="text-[#EDEDEE] mx-1">|</span>
            <button
              class="flex items-center gap-1.5 text-xs font-semibold text-[#00916E] hover:text-[#007a5c] transition-colors"
              @click.prevent="openNewItem('service')"
            >
              <svg class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                <path d="M12 5v14M5 12h14" />
              </svg>
              New service
            </button>
          </div>
        </div>
      </div>

      <!-- Line items table -->
      <div
        v-if="lineItems.length > 0"
        class="mt-4 border border-[#EDEDEE] rounded-xl overflow-hidden"
      >
        <!-- Line item cards -->
        <div class="divide-y divide-[#EDEDEE]">
          <div
            v-for="item in lineItems"
            :key="item.id"
            class="flex items-center gap-3 px-4 py-3"
          >
            <!-- Qty / Hrs stepper -->
            <div class="flex flex-col items-center w-10 shrink-0">
              <span class="text-[10px] font-semibold text-[#6F7177] uppercase mb-0.5">
                {{ item.type === 'service' ? 'Hrs' : 'Qty' }}
              </span>
              <button
                class="text-[#6F7177] hover:text-[#1E212B] transition-colors leading-none"
                @click="item.qty++"
              >
                <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M18 15l-6-6-6 6" /></svg>
              </button>
              <span class="text-sm font-semibold text-[#1E212B] leading-tight">{{ item.qty }}</span>
              <button
                class="text-[#6F7177] hover:text-[#1E212B] transition-colors leading-none"
                @click="item.qty > 1 && item.qty--"
              >
                <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M6 9l6 6 6-6" /></svg>
              </button>
            </div>

            <!-- Thumbnail -->
            <div class="w-9 h-9 rounded-lg overflow-hidden bg-[#F6F6F7] border border-[#EDEDEE] shrink-0 flex items-center justify-center">
              <img
                v-if="item.imageUrl"
                :src="item.imageUrl"
                :alt="item.name"
                class="w-full h-full object-cover"
              >
              <svg v-else class="w-4 h-4 text-[#D7DAE0]" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="3" y="3" width="18" height="18" rx="2" /><path d="M3 9l4-4 4 4 4-5 4 5" /></svg>
            </div>

            <!-- Name + discount badge -->
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-[#1E212B] truncate">
                {{ item.name }}
              </p>
              <span
                v-if="item.discountPercent"
                class="inline-block text-[10px] font-semibold text-[#2176AE] bg-[#EBF4FF] px-1.5 py-0.5 rounded mt-0.5"
              >
                {{ item.discountPercent }}% off
              </span>
            </div>

            <!-- Price + actions -->
            <div class="text-right shrink-0">
              <p class="text-sm font-semibold text-[#1E212B]">
                {{ formatCurrency(item.rate) }}
              </p>
              <div class="flex items-center justify-end gap-1.5 mt-1">
                <!-- Edit rate inline -->
                <button
                  class="w-6 h-6 flex items-center justify-center rounded text-[#939499] hover:text-[#2176AE] hover:bg-[#EBF4FF] transition-colors"
                  title="Edit rate"
                  @click="editingRateId = editingRateId === item.id ? null : item.id"
                >
                  <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7" /><path d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z" /></svg>
                </button>
                <button
                  class="w-6 h-6 flex items-center justify-center rounded text-[#939499] hover:text-[#AF513A] hover:bg-[#FEE2E2] transition-colors"
                  @click="removeLineItem(item.id)"
                >
                  <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M6 18L18 6M6 6l12 12" /></svg>
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Inline rate editor -->
        <div
          v-if="editingRateId"
          class="px-4 py-3 bg-[#F6F6F7] border-t border-[#EDEDEE] flex items-center gap-3"
        >
          <span class="text-xs font-medium text-[#6F7177] shrink-0">
            Edit rate for <span class="text-[#1E212B]">{{ lineItems.find(i => i.id === editingRateId)?.name }}</span>
          </span>
          <div class="flex items-center bg-white border border-[#EDEDEE] rounded-lg focus-within:border-[#1E212B] transition-colors">
            <span class="pl-2 text-xs text-[#939499]">$</span>
            <input
              :value="lineItems.find(i => i.id === editingRateId)?.rate"
              type="number"
              min="0"
              step="0.01"
              class="w-24 px-2 py-1.5 bg-transparent text-sm text-[#1E212B] outline-none"
              @change="(e) => { const item = lineItems.find(i => i.id === editingRateId); if (item) item.rate = Number((e.target as HTMLInputElement).value) }"
            >
          </div>
          <button
            class="text-xs font-semibold text-[#2176AE] hover:underline"
            @click="editingRateId = null"
          >
            Done
          </button>
        </div>

        <!-- Totals -->
        <div class="px-4 py-3 bg-[#FAFAFA] border-t border-[#EDEDEE] space-y-2">
          <div class="flex justify-between text-sm text-[#6F7177]">
            <span>Subtotal</span>
            <span class="font-medium text-[#1E212B]">{{ formatCurrency(subtotal) }}</span>
          </div>
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <span class="text-sm text-[#6F7177]">Tax (%)</span>
              <input
                v-model.number="taxRate"
                type="number"
                min="0"
                max="100"
                step="0.5"
                class="w-16 text-center text-sm bg-white border border-[#EDEDEE] rounded-lg px-2 py-1 outline-none focus:border-[#1E212B] transition-colors"
              >
            </div>
            <span class="text-sm font-medium text-[#1E212B]">{{ formatCurrency(taxAmount) }}</span>
          </div>
          <div class="flex justify-between pt-2 border-t border-[#EDEDEE]">
            <span class="text-sm font-bold text-[#1E212B]">Total</span>
            <span class="text-sm font-bold text-[#E87117]">{{ formatCurrency(total) }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- ── Payments Collection ───────────────────────────────── -->
    <div class="py-8 px-8 border-b border-[#EDEDEE]">
      <h6 class="text-sm font-semibold mb-4">
        Payments Collection
      </h6>

      <p
        v-if="lineItems.length === 0"
        class="text-xs text-[#939499] italic"
      >
        Add items above to configure payment collection.
      </p>

      <div
        v-else
        class="space-y-6"
      >
        <!-- One-time section -->
        <div v-if="hasOneTime">
          <div class="flex items-center gap-2 mb-3">
            <span class="w-2 h-2 rounded-full bg-[#1E212B] shrink-0" />
            <span class="text-xs font-semibold text-[#1E212B] uppercase tracking-wider">One-time</span>
          </div>
          <div class="pl-4 space-y-3">
            <div class="flex justify-between text-sm">
              <span class="text-[#6F7177]">Amount due</span>
              <span class="font-semibold text-[#1E212B]">{{ formatCurrency(oneTimeTotal) }}</span>
            </div>
            <div>
              <label class="block text-xs font-medium text-[#1E212B] mb-1">Due date</label>
              <input
                v-model="dueDate"
                type="date"
                class="w-full bg-[#F6F6F7] border border-[#EDEDEE] hover:border-[#94BDD8] focus:border-[#1E212B] rounded-lg px-3 py-2.5 text-sm text-[#1E212B] outline-none transition-colors"
              >
            </div>
          </div>
        </div>

        <!-- Recurring section -->
        <div v-if="hasRecurring">
          <div class="flex items-center gap-2 mb-3">
            <span class="w-2 h-2 rounded-full bg-[#2176AE] shrink-0" />
            <span class="text-xs font-semibold text-[#2176AE] uppercase tracking-wider">Recurring</span>
          </div>
          <div class="pl-4 space-y-3">
            <div
              v-for="(amount, cycle) in recurringByBillingCycle"
              :key="cycle"
              class="flex justify-between text-sm"
            >
              <span class="text-[#6F7177]">Every {{ cycle }}</span>
              <span class="font-semibold text-[#1E212B]">{{ formatCurrency(amount) }}</span>
            </div>
            <div>
              <label class="block text-xs font-medium text-[#1E212B] mb-1">Billing starts</label>
              <input
                v-model="recurringStartDate"
                type="date"
                class="w-full bg-[#F6F6F7] border border-[#EDEDEE] hover:border-[#94BDD8] focus:border-[#1E212B] rounded-lg px-3 py-2.5 text-sm text-[#1E212B] outline-none transition-colors"
              >
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- ── Notes & Attachments ───────────────────────────────── -->
    <div class="py-8 px-8">
      <h6 class="text-sm font-semibold mb-6">
        Notes & Attachments
      </h6>

      <GroBasicTextArea
        v-model="notesToCustomer"
        placeholder="Note to your customer"
        border-radius="8px"
        :max-length="100"
      >
        <div class="flex w-full">
          Notes
          <small class="text-xs ml-auto text-[#939499]">{{ notesToCustomer?.length || 0 }}/100</small>
        </div>
      </GroBasicTextArea>

      <div class="mt-8 space-y-6">
        <!-- Reference Number -->
        <div>
          <button
            class="flex items-center gap-x-2 text-left"
            @click="showReferenceNumber = !showReferenceNumber"
          >
            <HugeiconsIcon
              :icon="showReferenceNumber ? MinusSignIcon : PlusSignIcon"
              :size="20"
              color="#2176AE"
            />
            <span class="font-semibold text-sm text-[#2176AE]">Reference Number</span>
          </button>
          <GroBasicInput
            v-if="showReferenceNumber"
            v-model="referenceNumber"
            class="mt-3"
            placeholder="e.g. INV-2024-001"
          />
        </div>

        <!-- Terms & Conditions -->
        <div>
          <button
            class="flex items-center gap-x-2 text-left"
            @click="showTermsAndCondition = !showTermsAndCondition"
          >
            <HugeiconsIcon
              :icon="showTermsAndCondition ? MinusSignIcon : PlusSignIcon"
              :size="20"
              color="#2176AE"
            />
            <span class="font-semibold text-sm text-[#2176AE]">Terms and conditions</span>
          </button>
          <GroBasicTextArea
            v-if="showTermsAndCondition"
            v-model="termsAndCondition"
            placeholder="Enter your terms and conditions"
            border-radius="8px"
            class="mt-3"
          />
        </div>

        <!-- Memo to self -->
        <div>
          <button
            class="flex items-center gap-x-2 text-left"
            @click="showMemos = !showMemos"
          >
            <HugeiconsIcon
              :icon="showMemos ? MinusSignIcon : PlusSignIcon"
              :size="20"
              color="#2176AE"
            />
            <span class="font-semibold text-sm text-[#2176AE]">Memo to self</span>
          </button>
          <GroBasicTextArea
            v-if="showMemos"
            v-model="memos"
            placeholder="Private memo (not visible to customer)"
            border-radius="8px"
            class="mt-3"
          />
        </div>
      </div>

      <!-- Attachments -->
      <div class="mt-8">
        <h6 class="text-sm font-semibold text-[#1E212B] mb-4">
          Attachments
        </h6>
        <GroFileUpload
          accept="image/*,.pdf,.doc,.docx,.xls,.xlsx,.csv"
          :max-files="5"
          :max-size-mb="20"
          hint="Accepts PDF, images and documents"
          @update:files="attachments = $event"
        >
          <template #label>
            <span class="text-xs font-medium text-[#1E212B]">Attachments</span>
          </template>
        </GroFileUpload>
      </div>
    </div>
  </div>

  <!-- New product / service modal -->
  <CreateInvoiceItemModal
    v-model="showNewItemModal"
    :type="newItemType"
    @item-added="onNewItemAdded"
  />

  <!-- Variant / plan picker modal -->
  <CatalogItemPickerModal
    v-model="showPickerModal"
    :entry="pickerEntry"
    @item-added="pushLineItem"
  />
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import MainDashboard from "@/components/dashboard/main-dashboard.vue"
import GroBasicButton from "@/components/buttons/GroBasicButton.vue"
import { useRouter } from 'nuxt/app'
import { useProductAPI } from '@/composables/api/catalog/product'
import type { Product as APIProduct } from '@/composables/dto/catalog/product'
import GroDataTable from '@/components/table/GroDataTable.vue'
import type { TableColumn, TableAction, TableFilter } from '@/components/table/GroDataTable.vue'
import { Tag01Icon, WarehouseIcon, Store01Icon, PercentIcon, CheckmarkCircle02Icon, FilterIcon } from '@hugeicons/core-free-icons'

const router = useRouter()
const isLoading = ref(false)
const showMoreActions = ref(false)
const PAGE_SIZE = 10
const currentPage = ref(1)
const totalPages = ref(1)
const totalCount = ref(0)

interface Product extends Record<string, unknown> {
  id: string
  name: string
  category: string
  ribbon: string
  price: number
  price_currency: string
  has_discount: string
  has_variants: string
  status: string
  show_in_store: string
  track_inventory: string
  stock_status: string
  stock_count: number
  image_url: string | null
}

const products = ref<Product[]>([])

const mapProduct = (p: APIProduct): Product => {
  const detail = p.product_detail
  const rawStock = detail?.stock_status ?? ''
  const stockLabel = rawStock === 'in_stock' ? 'In stock'
    : rawStock === 'limited_stock' ? 'Limited stock'
    : rawStock === 'low_stock' ? 'Low in stock'
    : rawStock === 'out_of_stock' ? 'Out of stock'
    : '—'

  return {
    id: p.id,
    name: p.name,
    category: p.standard_category?.name || '—',
    ribbon: detail?.ribbon || '—',
    price: p.price,
    price_currency: p.price_currency,
    has_discount: String((p.discount_percent ?? 0) > 0),
    has_variants: String((detail?.variants?.length ?? 0) > 0),
    status: p.status || 'draft',
    show_in_store: String(p.show_in_store ?? false),
    track_inventory: String(detail?.track_inventory ?? false),
    stock_status: stockLabel,
    stock_count: detail?.quantity ?? 0,
    image_url: p.media?.length ? p.media.find(m => m.position === 0)?.url ?? p.media[0].url : null,
  }
}

const fetchProducts = async (page = 1) => {
  isLoading.value = true
  try {
    const result = await useProductAPI().ListProducts({ page, limit: PAGE_SIZE })
    if (result.success && result.data?.data) {
      const paginated = result.data.data
      products.value = paginated.data.map(mapProduct)
      totalCount.value = Number(paginated.total)
      totalPages.value = paginated.total_pages
      currentPage.value = paginated.page
    }
  } finally {
    isLoading.value = false
  }
}

onMounted(() => fetchProducts(1))

const tableColumns: TableColumn[] = [
  { key: 'name', label: 'Name' },
  { key: 'category', label: 'Category' },
  { key: 'price', label: 'Price' },
  { key: 'has_variants', label: 'Has Variants' },
  { key: 'status', label: 'Status' },
  { key: 'show_in_store', label: 'Show in Site' },
]

const tableActions: TableAction[] = [
  { key: 'edit', label: 'Edit' },
]

const productFilters: TableFilter[] = [
  {
    key: 'status',
    label: 'Status',
    icon: Tag01Icon,
    priority: true,
    options: [
      { label: 'All', value: '' },
      { label: 'Active', value: 'active' },
      { label: 'Draft', value: 'draft' },
      { label: 'Archived', value: 'archived' },
    ],
  },
  {
    key: 'show_in_store',
    label: 'Show in Site',
    icon: Store01Icon,
    priority: true,
    options: [
      { label: 'All', value: '' },
      { label: 'Visible', value: 'true' },
      { label: 'Hidden', value: 'false' },
    ],
  },
  {
    key: 'has_variants',
    label: 'Has Variants',
    icon: FilterIcon,
    priority: true,
    options: [
      { label: 'All', value: '' },
      { label: 'Yes', value: 'true' },
      { label: 'No', value: 'false' },
    ],
  },
  {
    key: 'has_discount',
    label: 'Has Discount',
    icon: PercentIcon,
    priority: true,
    options: [
      { label: 'All', value: '' },
      { label: 'Yes', value: 'true' },
      { label: 'No', value: 'false' },
    ],
  },
  {
    key: 'track_inventory',
    label: 'Track Inventory',
    icon: WarehouseIcon,
    priority: true,
    options: [
      { label: 'All', value: '' },
      { label: 'Yes', value: 'true' },
      { label: 'No', value: 'false' },
    ],
  },
  {
    key: 'stock_status',
    label: 'Stock Status',
    icon: CheckmarkCircle02Icon,
    priority: true,
    options: [
      { label: 'All', value: '' },
      { label: 'In stock', value: 'In stock' },
      { label: 'Low in stock', value: 'Low in stock' },
      { label: 'Limited stock', value: 'Limited stock' },
      { label: 'Out of stock', value: 'Out of stock' },
    ],
  },
]

const CURRENCY_SYMBOLS: Record<string, string> = {
  NGN: '₦', USD: '$', EUR: '€', GBP: '£', GHS: 'GH₵', KES: 'KSh', ZAR: 'R',
}

const formatPrice = (price: number, currency = 'USD') => {
  const sym = CURRENCY_SYMBOLS[currency] ?? currency
  return `${sym}${price.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`
}

const handleAction = (key: string, row: Record<string, unknown>) => {
  if (key === 'edit') router.push(`/catalog/products/${row.id}`)
}
</script>

<template>
  <MainDashboard current="Products">
    <template #title>
      <div class="flex items-center justify-between">
        <h6 class="text-2xl font-semibold">
          Products
          <span class="text-[#6F7177] font-normal">({{ isLoading ? '...' : totalCount }})</span>
        </h6>
        <div class="flex items-center gap-3">
          <div class="relative">
            <button
              class="flex items-center gap-2 px-4 py-2 cursor-pointer bg-white border border-[#DBDBDD] rounded-lg text-sm font-medium text-[#1E212B] hover:bg-[#F6F6F7] shadow-[inset_0_1px_0_0_#E3E3E3]"
              @click="showMoreActions = !showMoreActions"
            >
              More actions
              <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                <path d="M6 9l6 6 6-6" />
              </svg>
            </button>
            <div
              v-if="showMoreActions"
              class="absolute right-0 top-full mt-1 bg-white border border-[#EDEDEE] rounded-xl shadow-lg z-10 min-w-36 py-1"
            >
              <button class="w-full text-left px-4 py-2 text-sm text-[#1E212B] hover:bg-[#F6F6F7] cursor-pointer">Export</button>
              <button class="w-full text-left px-4 py-2 text-sm text-[#1E212B] hover:bg-[#F6F6F7] cursor-pointer">Import</button>
            </div>
          </div>
          <GroBasicButton
            color="primary"
            size="xs"
            shape="custom"
            class="w-max"
            @click="router.push('/catalog/products/new')"
          >
            + New Product
          </GroBasicButton>
        </div>
      </div>
    </template>

    <template #body>
      <div class="mt-6">
        <GroDataTable
          :columns="tableColumns"
          :rows="products"
          row-key="id"
          :actions="tableActions"
          :filters="productFilters"
          search-placeholder="Search products"
          :is-loading="isLoading"
          :page="currentPage"
          :total-pages="totalPages"
          :total="totalCount"
          empty-title="No products found"
          empty-message="Create your first product to get started"
          @row-click="(row) => router.push(`/catalog/products/${row.id}`)"
          @action="handleAction"
          @page-change="fetchProducts"
        >
          <!-- Image + name -->
          <template #cell-name="{ row }">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-lg bg-[#EDEDEE] overflow-hidden flex items-center justify-center shrink-0">
                <img v-if="row.image_url" :src="String(row.image_url)" :alt="String(row.name)" class="w-full h-full object-cover" />
              </div>
              <span class="text-sm font-medium text-[#1E212B]">{{ row.name }}</span>
            </div>
          </template>

          <!-- Formatted price -->
          <template #cell-price="{ row }">
            <span class="text-sm font-medium text-[#1E212B]">
              {{ formatPrice(Number(row.price), String(row.price_currency)) }}
            </span>
          </template>

          <!-- Has variants -->
          <template #cell-has_variants="{ row }">
            <span
              class="inline-flex px-2 py-0.5 rounded-full text-xs font-medium"
              :class="row.has_variants === 'true' ? 'bg-[#EFF6FF] text-[#1D4ED8]' : 'bg-[#F6F6F7] text-[#6F7177]'"
            >
              {{ row.has_variants === 'true' ? 'Yes' : 'No' }}
            </span>
          </template>

          <!-- Status badge -->
          <template #cell-status="{ row }">
            <span
              class="inline-flex px-2 py-0.5 rounded-full text-xs font-medium"
              :class="row.status === 'active' ? 'bg-[#F0FDF4] text-[#15803D]'
                : row.status === 'archived' ? 'bg-[#FFF1F2] text-[#BE123C]'
                : 'bg-[#F6F6F7] text-[#6F7177]'"
            >
              {{ String(row.status).charAt(0).toUpperCase() + String(row.status).slice(1) }}
            </span>
          </template>


          <!-- Show in site -->
          <template #cell-show_in_store="{ row }">
            <span
              class="inline-flex px-2 py-0.5 rounded-full text-xs font-medium"
              :class="row.show_in_store === 'true' ? 'bg-[#F0FDF4] text-[#15803D]' : 'bg-[#F6F6F7] text-[#6F7177]'"
            >
              {{ row.show_in_store === 'true' ? 'Visible' : 'Hidden' }}
            </span>
          </template>

        </GroDataTable>
      </div>
    </template>
  </MainDashboard>
</template>

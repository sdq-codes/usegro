<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import MainDashboard from '@/components/dashboard/main-dashboard.vue'
import GroBasicButton from '@/components/buttons/GroBasicButton.vue'
import { HugeiconsIcon } from '@hugeicons/vue'
import { Search01Icon } from '@hugeicons/core-free-icons'
import { useRouter } from 'nuxt/app'
import { useServiceAPI } from '@/composables/api/catalog/service'
import type { CatalogItem } from '@/composables/dto/catalog/product'

const router = useRouter()
const query = ref('')
const isLoading = ref(false)
const selectedServices = ref<string[]>([])

const services = ref<CatalogItem[]>([])

onMounted(async () => {
  isLoading.value = true
  const result = await useServiceAPI().ListServices()
  if (result.success && result.data?.data) {
    services.value = result.data.data
  }
  isLoading.value = false
})

const filteredServices = computed(() => {
  if (!query.value) return services.value
  const term = query.value.toLowerCase()
  return services.value.filter(s => s.name.toLowerCase().includes(term))
})

const CURRENCY_SYMBOLS: Record<string, string> = {
  NGN: '₦', USD: '$', EUR: '€', GBP: '£', GHS: 'GH₵', KES: 'KSh', ZAR: 'R',
}
const formatPrice = (price: number, currency = 'USD') => {
  if (price === 0) return 'Free'
  const sym = CURRENCY_SYMBOLS[currency] ?? currency
  return `${sym}${price.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`
}

const serviceTypeLabel = (item: CatalogItem) => {
  const t = item.service_detail?.service_type
  if (t === 'appointment') return 'Appointment'
  if (t === 'class') return 'Class'
  if (t === 'course') return 'Course'
  return 'Service'
}

const statusClass = (status: string) => {
  if (status === 'draft') return 'text-[#D26B06]'
  if (status === 'archived') return 'text-[#AF513A]'
  return 'text-[#1E212B]'
}

const isAllSelected = computed(() =>
  selectedServices.value.length === filteredServices.value.length && filteredServices.value.length > 0,
)
const toggleAll = () => {
  isAllSelected.value
    ? (selectedServices.value = [])
    : (selectedServices.value = filteredServices.value.map(s => s.id))
}
const toggleService = (id: string) => {
  const idx = selectedServices.value.indexOf(id)
  if (idx > -1) selectedServices.value.splice(idx, 1)
  else selectedServices.value.push(id)
}
</script>

<template>
  <MainDashboard current="Services">
    <template #title>
      <div class="flex items-center justify-between">
        <div>
          <h6 class="text-2xl font-semibold">
            Services
            <span class="text-[#6F7177] font-normal">({{ isLoading ? '...' : services.length }})</span>
          </h6>
          <p class="text-sm text-[#6F7177] mt-0.5">
            Manage your services and offerings.
          </p>
        </div>
        <GroBasicButton
          color="primary"
          size="xs"
          shape="custom"
          class="w-max"
          @click="router.push('/catalog/services/new')"
        >
          Add Service
        </GroBasicButton>
      </div>
    </template>

    <template #body>
      <!-- List view (once there are services) -->
      <div
        v-if="isLoading || services.length > 0"
        class="mt-6 bg-white rounded-2xl border border-[#EDEDEE] overflow-hidden"
      >
        <!-- Search bar -->
        <div class="flex items-center gap-3 px-4 py-3 border-b border-[#EDEDEE]">
          <div class="ml-auto relative">
            <HugeiconsIcon :icon="Search01Icon" class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4" color="#939499" />
            <input
              v-model="query"
              type="text"
              placeholder="Search"
              class="pl-9 pr-4 py-2 bg-[#F6F6F7] border border-[#EDEDEE] rounded-lg text-sm outline-none focus:border-[#1E212B] w-64 placeholder-[#939499] text-[#1E212B]"
            >
          </div>
        </div>

        <!-- Table -->
        <table class="w-full">
          <thead>
            <tr class="bg-[#F6F6F7] border-b border-[#EDEDEE]">
              <th class="w-12 px-4 py-3">
                <input type="checkbox" :checked="isAllSelected" class="rounded cursor-pointer accent-[#1E212B]" @change="toggleAll">
              </th>
              <th class="px-4 py-3 text-left text-xs font-semibold text-[#1E212B]">
                Name
              </th>
              <th class="px-4 py-3 text-left text-xs font-semibold text-[#1E212B]">
                Type
              </th>
              <th class="px-4 py-3 text-left text-xs font-semibold text-[#1E212B]">
                Price
              </th>
              <th class="px-4 py-3 text-left text-xs font-semibold text-[#1E212B]">
                Status
              </th>
              <th class="w-12 px-4 py-3" />
            </tr>
          </thead>
          <tbody v-if="!isLoading">
            <tr
              v-for="service in filteredServices"
              :key="service.id"
              class="border-b border-[#EDEDEE] hover:bg-[#F6F6F7] cursor-pointer transition-colors"
              @click="router.push(`/catalog/services/${service.id}`)"
            >
              <td class="px-4 py-4" @click.stop>
                <input type="checkbox" :checked="selectedServices.includes(service.id)" class="rounded cursor-pointer accent-[#1E212B]" @change="toggleService(service.id)">
              </td>
              <td class="px-4 py-4">
                <div class="text-sm font-medium text-[#1E212B]">
                  {{ service.name }}
                </div>
                <div v-if="service.service_detail?.tagline" class="text-xs text-[#6F7177] mt-0.5">
                  {{ service.service_detail.tagline }}
                </div>
              </td>
              <td class="px-4 py-4 text-sm text-[#4B4D55]">
                {{ serviceTypeLabel(service) }}
              </td>
              <td class="px-4 py-4 text-sm font-medium text-[#1E212B]">
                {{ formatPrice(service.price, service.price_currency) }}
              </td>
              <td class="px-4 py-4">
                <span class="text-sm font-medium capitalize" :class="statusClass(service.status)">{{ service.status }}</span>
              </td>
              <td class="px-4 py-4" @click.stop>
                <button class="flex items-center justify-center w-8 h-8 rounded-full hover:bg-[#EDEDEE] text-[#4B4D55]">
                  <svg viewBox="0 0 24 24" fill="currentColor" class="w-4 h-4">
                    <circle cx="12" cy="5" r="1.5" /><circle cx="12" cy="12" r="1.5" /><circle cx="12" cy="19" r="1.5" />
                  </svg>
                </button>
              </td>
            </tr>
          </tbody>
        </table>

        <div v-if="isLoading" class="py-16 text-center text-sm text-[#6F7177]">
          Loading services...
        </div>
        <div v-else-if="filteredServices.length === 0 && query" class="py-16 text-center">
          <h5 class="text-lg font-semibold text-[#1E212B]">
            No services found
          </h5>
          <p class="text-sm text-[#6F7177] mt-1">
            Try adjusting your search
          </p>
        </div>
      </div>

      <!-- Empty state (no services yet) -->
      <div
        v-else
        class="rounded-3xl bg-white px-8 py-12 mt-6"
      >
        <div class="flex items-center justify-center w-full">
          <div class="rounded-full bg-[#EDEDEE] h-48 w-48" />
        </div>
        <h5 class="text-center mt-4 text-lg font-semibold">
          Add your first service
        </h5>
        <h6 class="text-center mt-2 text-sm text-[#6F7177]">
          Define the services you offer, set pricing, <br> and manage availability.
        </h6>
        <div class="flex items-center justify-center w-full mt-6">
          <GroBasicButton color="primary" size="xs" shape="custom" class="w-max" @click="router.push('/catalog/services/new')">
            Add Service
          </GroBasicButton>
        </div>
      </div>
    </template>
  </MainDashboard>
</template>

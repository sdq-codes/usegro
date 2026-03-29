<script setup lang="ts">
import { onMounted, ref } from 'vue'
import MainDashboard from "@/components/dashboard/main-dashboard.vue"
import GroBasicButton from "@/components/buttons/GroBasicButton.vue"
import BasicModal from "@/components/modals/basic-modal.vue"
import { useCustomerAPI, type PaginatedCustomers } from "@/composables/api/customer/customer"
import { notify } from "@/composables/helpers/notification/notification"
import type { FormSubmission } from "@/composables/dto/customer/customer"
import { useRouter } from "nuxt/app"
import CreateCustomerModal from "@/components/modals/modal/create-customer-modal.vue"
import { PlusSignIcon, UserIcon, CheckmarkCircle02Icon } from "@hugeicons/core-free-icons"
import { HugeiconsIcon } from "@hugeicons/vue"
import GroDataTable from "@/components/table/GroDataTable.vue"
import type { TableColumn, TableAction, TableFilter } from "@/components/table/GroDataTable.vue"


// Mapped shape used inside this page
interface MappedCustomer extends Record<string, unknown> {
  customerId: string
  formId: string
  name: string
  email: string
  phone_number: string
  address: string
  customer_type: string
  status: string
}

const isOpen = ref(false)
const isDeleteModalOpen = ref(false)
const customerToDelete = ref<MappedCustomer | null>(null)
const isDeleting = ref(false)
const isLoading = ref(false)
const currentTab = ref('customers')

// Pagination
const currentPage = ref(1)
const totalPages = ref(1)
const totalCount = ref(0)
const PAGE_SIZE = 20

// null = not yet loaded
const hasAnyCustomers = ref<boolean | null>(null)
const rawCustomers = ref<MappedCustomer[]>([])

const customerFilters: TableFilter[] = [
  {
    key: 'customer_type',
    label: 'Type',
    icon: UserIcon,
    priority: true,
    options: [
      { label: 'All', value: '' },
      { label: 'Individual', value: 'Individual' },
      { label: 'Business', value: 'Business' },
    ],
  },
  {
    key: 'status',
    label: 'Status',
    icon: CheckmarkCircle02Icon,
    priority: true,
    options: [
      { label: 'All', value: '' },
      { label: 'Active', value: 'active' },
      { label: 'Archived', value: 'archived' },
    ],
  },
]

const tableColumns: TableColumn[] = [
  { key: 'name', label: 'Customer Name' },
  { key: 'email', label: 'Email' },
  { key: 'phone_number', label: 'Phone Number' },
  { key: 'address', label: 'Address' },
  { key: 'customer_type', label: 'Type' },
]

const tableActions: TableAction[] = [
  { key: 'view', label: 'View' },
  { key: 'delete', label: 'Delete', class: 'text-red-600' },
]

const router = useRouter()

function mapSubmission(submission: FormSubmission): MappedCustomer {
  const answers = submission.answers || {}
  const displayName =
    answers.customer_type === 'Individual'
      ? `${answers.first_name || ''} ${answers.last_name || ''}`.trim()
      : (answers.company_name as string) || ''

  const rawPhone = answers.phone_number
  let phone = ''
  if (typeof rawPhone === 'string') {
    phone = rawPhone
  } else if (rawPhone && typeof rawPhone === 'object') {
    const strVal = Object.values(rawPhone as Record<string, unknown>)
      .find(v => typeof v === 'string' && (v as string).length > 0)
    phone = typeof strVal === 'string' ? strVal : ''
  }

  return {
    ...(answers as Record<string, unknown>),
    phone_number: phone ? `(${phone})` : '',
    name: displayName,
    customerId: submission._id,
    formId: submission.formID,
    status: submission.status,
  } as MappedCustomer
}

const fetchCustomers = async (page = 1) => {
  isLoading.value = true
  try {
    const result = await useCustomerAPI().FetchCustomers(page, PAGE_SIZE)
    if (!result.success) return

    // result.data is the CommonResponse body; paginated payload lives at .data
    const body = result.data as unknown as { data: PaginatedCustomers }
    const paginated = body.data

    rawCustomers.value = paginated.data.map(mapSubmission)
    totalCount.value = paginated.total
    totalPages.value = paginated.total_pages
    currentPage.value = paginated.page

    if (hasAnyCustomers.value === null) {
      hasAnyCustomers.value = paginated.total > 0
    }
  } finally {
    isLoading.value = false
  }
}

onMounted(() => fetchCustomers(1))

const handlePageChange = (page: number) => fetchCustomers(page)

const viewCrmCustomer = (row: Record<string, unknown>) => {
  const customer = row as MappedCustomer
  router.push(`/customers/${customer.customerId}/${customer.formId}`)
}

const handleTableAction = (key: string, row: Record<string, unknown>) => {
  const customer = row as MappedCustomer
  if (key === 'view') viewCrmCustomer(row)
  else if (key === 'delete') {
    customerToDelete.value = customer
    isDeleteModalOpen.value = true
  }
}

const confirmDelete = async () => {
  if (!customerToDelete.value) return
  isDeleting.value = true
  try {
    const res = await useCustomerAPI().DeleteCustomer(
      customerToDelete.value.customerId,
      customerToDelete.value.formId,
    )
    if (res.success) {
      notify('Customer was successfully deleted', 'success')
      isDeleteModalOpen.value = false
      customerToDelete.value = null
      const targetPage = rawCustomers.value.length === 1 && currentPage.value > 1
        ? currentPage.value - 1
        : currentPage.value
      await fetchCustomers(targetPage)
      hasAnyCustomers.value = totalCount.value > 0
    } else {
      notify('Failed to delete customer', 'error')
    }
  } catch {
    notify('An error occurred while deleting the customer', 'error')
  } finally {
    isDeleting.value = false
  }
}

const cancelDelete = () => {
  isDeleteModalOpen.value = false
  customerToDelete.value = null
}
</script>

<template>
  <div>
    <MainDashboard current="Customers">
      <template #title>
        <div class="block md:flex justify-between">
          <div>
            <h6 class="text-2xl font-semibold">
              Customers
              <span v-if="totalCount > 0" class="text-[#6F7177] font-normal">
                ({{ totalCount }})
              </span>
            </h6>
            <small v-if="hasAnyCustomers" class="text-[#6F7177]">
              Manage and track your customers, leads and site members.
            </small>
          </div>
          <div v-if="hasAnyCustomers" class="my-auto">
            <div class="inline-flex space-x-4 mt-2 md:mt-0">
              <NuxtLink to="/customers/import-contacts" class="cursor-pointer">
                <GroBasicButton color="secondary" size="xs" shape="custom" class="w-max">
                  Import/export
                </GroBasicButton>
              </NuxtLink>
              <GroBasicButton
                color="primary"
                size="xs"
                shape="custom"
                class="w-max"
                @click="isOpen = true"
              >
                <template #frontIcon>
                  <HugeiconsIcon :icon="PlusSignIcon" :size="14" color="white" :stroke-width="2" />
                </template>
                Create New
              </GroBasicButton>
            </div>
          </div>
        </div>
      </template>

      <template #body>
        <div>
          <!-- Empty state -->
          <div
            v-if="hasAnyCustomers === false"
            class="rounded-3xl bg-white px-3 md:px-8 py-12 mt-3 md:mt-6"
          >
            <div class="flex items-center justify-center w-full">
              <div class="rounded-full bg-[#EDEDEE] h-48 w-48" />
            </div>
            <h5 class="text-center mt-4 text-lg font-semibold">
              Everything relating to customers here
            </h5>
            <h6 class="text-center mt-2 text-sm font-regular text-[#6F7177]">
              Manage customers details, see order history, emails,
              <br>
              and group customers into segments
            </h6>
            <div class="flex items-center justify-center w-full gap-x-6 mt-6">
              <NuxtLink to="/customers/import-contacts" class="cursor-pointer">
                <GroBasicButton color="secondary" size="xs" shape="custom">
                  Import Customers
                </GroBasicButton>
              </NuxtLink>
              <GroBasicButton
                color="primary"
                size="xs"
                class="w-max"
                shape="custom"
                @click="isOpen = true"
              >
                <template #frontIcon>
                  <HugeiconsIcon :icon="PlusSignIcon" :size="14" color="white" :stroke-width="2" />
                </template>
                Create New
              </GroBasicButton>
            </div>
          </div>

          <!-- Customer list -->
          <div v-else-if="hasAnyCustomers">
            <!-- Tabs -->
            <div class="flex mt-4 md:mt-10 space-x-6 border-b border-b-[#EDEDEE]">
              <div
                class="text-sm font-medium cursor-pointer pb-1"
                :class="currentTab === 'customers' ? 'text-[#D26B06] border-b-2 border-b-[#D26B06]' : 'text-[#1E212B]'"
                @click="currentTab = 'customers'"
              >
                <h6>Customers List</h6>
              </div>
              <div
                class="text-sm font-medium cursor-pointer pb-1"
                :class="currentTab === 'segments' ? 'text-[#D26B06] border-b-2 border-b-[#D26B06]' : 'text-[#1E212B]'"
                @click="currentTab = 'segments'"
              >
                <h6>Segments</h6>
              </div>
            </div>

            <div v-if="currentTab === 'customers'" class="mt-6">
              <GroDataTable
                :columns="tableColumns"
                :rows="rawCustomers"
                row-key="customerId"
                :actions="tableActions"
                :filters="customerFilters"
                search-placeholder="Search customers"
                :page="currentPage"
                :total-pages="totalPages"
                :total="totalCount"
                :is-loading="isLoading"
                empty-title="No customers found"
                empty-message="Try adjusting your search or filters"
                @row-click="viewCrmCustomer"
                @action="handleTableAction"
                @page-change="handlePageChange"
              >
                <!-- Avatar + name cell -->
                <template #cell-name="{ row }">
                  <div class="flex items-center gap-3">
                    <div class="w-8 h-8 rounded-full bg-[#EDEDEE] flex items-center justify-center text-xs font-semibold text-[#4B4D55] shrink-0">
                      {{ String(row.name || '?').charAt(0).toUpperCase() }}
                    </div>
                    <span class="text-sm font-medium text-[#1E212B]">{{ row.name || '—' }}</span>
                  </div>
                </template>

                <!-- Type badge cell -->
                <template #cell-customer_type="{ row }">
                  <span
                    class="inline-flex px-2 py-0.5 rounded-full text-xs font-medium"
                    :class="row.customer_type === 'Business'
                      ? 'bg-[#EFF6FF] text-[#1D4ED8]'
                      : 'bg-[#F0FDF4] text-[#15803D]'"
                  >
                    {{ row.customer_type || '—' }}
                  </span>
                </template>
              </GroDataTable>
            </div>
          </div>

          <CreateCustomerModal
            v-model="isOpen"
            @new-customer-added="fetchCustomers(1).then(() => { hasAnyCustomers = totalCount > 0 })"
          />

          <BasicModal v-model="isDeleteModalOpen" size="xs">
            <template #title>
              Delete Customer
            </template>
            <template #default>
              <div class="py-4">
                <h6 class="text-gray-700 text-md">
                  Are you sure you want to delete <strong>{{ customerToDelete?.name }}</strong>?
                </h6>
                <h6 class="text-red-500 font-bold text-sm mt-2">
                  This action cannot be undone.
                </h6>
              </div>
            </template>
            <template #footer>
              <div class="flex justify-end gap-x-4">
                <GroBasicButton
                  color="secondary"
                  size="sm"
                  shape="custom"
                  class="w-max"
                  :disabled="isDeleting"
                  @click="cancelDelete"
                >
                  Cancel
                </GroBasicButton>
                <GroBasicButton
                  color="primary"
                  size="sm"
                  shape="custom"
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

<style scoped>
</style>

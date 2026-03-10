<script setup lang="ts">
import {computed, onMounted, ref} from 'vue'
import MainDashboard from "@/components/dashboard/main-dashboard.vue"
import GroBasicButton from "@/components/buttons/GroBasicButton.vue"
import BasicModal from "@/components/modals/basic-modal.vue"
import GroBasicTable from "@/components/table/GroBasicTable.vue"
import { useCustomerAPI } from "@/composables/api/customer/customer"
import {notify} from "@/composables/helpers/notification/notification";
import type {FormSubmission} from "@/composables/dto/customer/customer";
import {LIST_CUSTOMER_COLUMNS} from "@/constants/tables/customers";
import {useRouter} from "nuxt/app";
import CreateCustomerModal from "@/components/modals/modal/create-customer-modal.vue";

const isOpen = ref(false)
const isDeleteModalOpen = ref(false)
const customerToDelete = ref<FormSubmission | null>(null)
const isDeleting = ref(false)
const query = ref<string>("")
const crmCustomersFullList = ref<FormSubmission[]>([])
const currentTab = ref<string>("customers")

const router = useRouter()


function searchCustomers(data: FormSubmission[], query: string): FormSubmission[] {
  const term = query.toLowerCase();

  return data.filter((item) =>
    Object.values(item).some((value) => {
      if (Array.isArray(value)) {
        // Search inside array elements
        return value.some((v) => String(v).toLowerCase().includes(term));
      }
      if (typeof value === 'object' && value !== null) {
        // Deep search inside nested objects
        return JSON.stringify(value).toLowerCase().includes(term);
      }
      // Regular scalar values
      return String(value).toLowerCase().includes(term);
    })
  );
}

const crmCustomers: FormSubmission[] = computed(() => {
  if(query.value === '') {
    return crmCustomersFullList.value;
  } else {
    return searchCustomers(crmCustomersFullList.value, query.value);
  }
})

const allCustomers = async () => {
  const createdCustomer = await useCustomerAPI().FetchCustomers()
  if (!createdCustomer.success) {
    return;
  }

  crmCustomersFullList.value = createdCustomer.data?.data.map(submission => {
    const answers = submission.Answers || {};

    // Apply your conditional display logic here
    const displayName =
      answers.customer_type === 'Individual'
        ? `${answers.first_name || ''} ${answers.last_name || ''}`.trim()
        : answers.company_name || '';

    return {
      ...answers,
      name: displayName,
      customerId: submission.SubmissionID,
      formId: submission.PK.split("#")[1].trim(),
    };
  });
}

const createCustomer = async () => {
  isOpen.value = true
}

onMounted(async () => {
  await allCustomers()
})


const viewCrmCustomer = async (crmCustomer: FormSubmission) => {
  await router.push(`/customers/${crmCustomer?.customerId}/${crmCustomer?.formId}`);
}

const deleteCrmCustomer = async (crmCustomer: FormSubmission) => {
  customerToDelete.value = crmCustomer
  isDeleteModalOpen.value = true
}

const confirmDelete = async () => {
  if (!customerToDelete.value) return

  isDeleting.value = true

  try {
    // Make API call to delete the customer
    const deleteResponse = await useCustomerAPI().DeleteCustomer(customerToDelete.value?.customerId, customerToDelete.value?.formId)

    if (deleteResponse.success) {
      // Remove the customer from the local list
      crmCustomersFullList.value = crmCustomersFullList.value.filter(
        customer => customer.customerId !== customerToDelete.value?.customerId
      )

      notify("Customer was successfully deleted", "success")
      isDeleteModalOpen.value = false
      customerToDelete.value = null
    } else {
      notify("Failed to delete customer", "error")
    }
  } catch (error) {
    notify("An error occurred while deleting the customer", "error")
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
            <h6 class="text-2xl font-semibold mb-1">
              Customers
            </h6>
            <small class="text-[#6F7177]">
              Manage and track your customers, leads and site members.
            </small>
          </div>
          <div
            v-if="crmCustomersFullList.length > 0"
            class="my-auto"
          >
            <div class="inline-flex space-x-4">
              <NuxtLink
                to="/customers/import-contacts"
                class="cursor-pointer"
              >
                <GroBasicButton
                  color="secondary"
                  size="xs"
                  shape="custom"
                  class="w-max"
                >
                  Import/export
                </GroBasicButton>
              </NuxtLink>
              <GroBasicButton
                color="primary"
                size="xs"
                shape="custom"
                class="w-max"
                @click="createCustomer"
              >
                Add Customer
              </GroBasicButton>
            </div>
          </div>
        </div>
      </template>
      <template #body>
        <div>
          <div
            v-if="crmCustomersFullList.length > 0"
            class="flex mt-10 space-x-6 border-b-1 border-b-[#EDEDEE]"
          >
            <div
              class="text-sm font-medium cursor-pointer pb-1"
              :class="[currentTab === 'customers' ? 'text-[#D26B06] border-b-2 border-b-[#D26B06]' : 'text-[#1E212B]']"
              @click="currentTab = 'customers'"
            >
              <h6 class="">
                Customers List
              </h6>
            </div>
            <div
              class="text-sm font-medium cursor-pointer pb-1"
              :class="[currentTab === 'segments' ? 'text-[#D26B06] border-b-2 border-b-[#D26B06]' : 'text-[#1E212B]']"
              @click="currentTab = 'segments'"
            >
              <h6>
                Segments
              </h6>
            </div>
          </div>
          <div
            v-if="crmCustomersFullList.length === 0"
            class="rounded-3xl bg-white px-8 py-12 mt-6"
          >
            <div class="flex items-center justify-center w-full">
              <div class="rounded-full bg-[#EDEDEE] h-48 w-48" />
            </div>
            <h5 class="text-center mt-4 text-lg font-semibold">
              Everything relating to customers here
            </h5>
            <h6 class="text-center mt-2 text-sm font-regular text-[#6F7177]">
              Manage customers details, see order history, emails, <br> and group customers into segments
            </h6>
            <div class="flex items-center justify-center w-full gap-x-6 mt-6">
              <NuxtLink
                to="/customers/import-contacts"
                class="cursor-pointer"
              >
                <GroBasicButton
                  color="secondary"
                  size="xs"
                  shape="custom"
                >
                  Import Customers
                </GroBasicButton>
              </NuxtLink>
              <GroBasicButton
                color="primary"
                size="xs"
                class="w-max"
                shape="custom"
                @click="createCustomer"
              >
                Add Customer
              </GroBasicButton>
            </div>
          </div>
          <div v-else>
            <div v-if="currentTab === 'customers'">
              <div class="mt-7">
                <GroBasicTable
                  v-model="query"
                  :cols="LIST_CUSTOMER_COLUMNS"
                  :rows="crmCustomers"
                  :on-view="viewCrmCustomer"
                  :on-delete="deleteCrmCustomer"
                />
              </div>
            </div>
          </div>

          <CreateCustomerModal
            v-model="isOpen"
            @new-customer-added="allCustomers"
          />

          <BasicModal
            v-model="isDeleteModalOpen"
            size="xs"
          >
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

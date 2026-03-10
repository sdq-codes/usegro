<script setup lang="ts">

import CreateInvoiceDashboard from "@/components/dashboard/page/invoice/create-invoice-dashboard.vue";
import GroBasicButton from "@/components/buttons/GroBasicButton.vue";
import {useCustomerAPI} from "@/composables/api/customer/customer";
import {onMounted, ref} from "vue";
import type {FormSubmission} from "@/composables/dto/customer/customer";
import CreateNewInvoiceForm from "@/components/page/invoice/create-new-invoice-form.vue";
import CreateCustomerModal from "@/components/modals/modal/create-customer-modal.vue";
import PreviewNewInvoice from "@/components/page/invoice/preview-new-invoice.vue";

const crmCustomersFullList = ref<FormSubmission[]>([])
const selectedCustomer = ref<Array<any>>([])
const selectedProducts = ref<Array<any>>([])
const remainingReOccuringInvoice = ref<number>(6)

const isOpen = ref(false)

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
      label: displayName,
      value: submission.SubmissionID,
    };
  });
}

onMounted(async () => {
  await allCustomers()
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
          <button class="text-[#2176AE] text-sm whitespace-nowrap">
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
          >
            Save as Draft
          </GroBasicButton>
          <GroBasicButton
            class="whitespace-nowrap"
            color="primary"
          >
            Send invoice
          </GroBasicButton>
        </div>
      </template>
      <template #body>
        <div class="flex h-screen">
          <CreateNewInvoiceForm
            v-model:selected-customer="selectedCustomer"
            v-model:selected-products="selectedProducts"
            v-model:remaining-re-occuring-invoice="remainingReOccuringInvoice"
            :crm-customers-full-list="crmCustomersFullList"
            @create-new-customer="isOpen = true"
          />
          <PreviewNewInvoice />
          <CreateCustomerModal
            v-model="isOpen"
            @new-customer-added="allCustomers"
          />
        </div>
      </template>
    </CreateInvoiceDashboard>
  </div>
</template>
cutt.ly/ownhcIV

<style scoped>

</style>

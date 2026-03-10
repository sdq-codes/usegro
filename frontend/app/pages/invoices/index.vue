<script setup lang="ts">
import {onMounted, ref} from 'vue'
import MainDashboard from "@/components/dashboard/main-dashboard.vue"
import GroBasicButton from "@/components/buttons/GroBasicButton.vue"
import { useCustomerAPI } from "@/composables/api/customer/customer"
import { HugeiconsIcon } from "@hugeicons/vue"
import { Tick01FreeIcons} from "@hugeicons/core-free-icons"
import type {FormSubmission} from "@/composables/dto/customer/customer";

const crmCustomersFullList = ref<FormSubmission[]>([])

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

onMounted(async () => {
  await allCustomers()
})
//
// const createCustomerApi = async () => {
//   const formData = {
//     type: "customer",
//     answers: answerMap.value,
//     versionSnap: formFields.value,
//   }
//   const createdCustomer = await useFormAPI().CreateCustomer(formVersion.value?.formID, formVersion.value?.SK.split("#")[1], formData)
//
//   return createdCustomer.success
// }

</script>

<template>
  <div>
    <MainDashboard current="Invoices">
      <template #title>
        <div class="block md:flex justify-between">
          <div>
            <h6 class="text-2xl font-semibold mb-1">
              Invoices
            </h6>
          </div>
        </div>
      </template>
      <template #body>
        <div>
          <div
            class="rounded-3xl bg-white px-8 py-12 mt-6"
          >
            <div class="flex items-center justify-center w-full">
              <div class="rounded-full bg-[#EDEDEE] h-48 w-48" />
            </div>
            <h5 class="text-center mt-4 text-lg font-semibold">
              Collecting payments with Gro Invoices
            </h5>
            <div class="flex w-full">
              <div class="flex mx-auto mt-2">
                <HugeiconsIcon
                  :icon="Tick01FreeIcons"
                  color="#6F7177"
                />
                <h6 class="text-[#1E212B] ml-2 text-sm">
                  Manage and track invoices all in one place
                </h6>
              </div>
            </div>
            <div class="flex w-full">
              <div class="flex mx-auto mt-2">
                <HugeiconsIcon
                  :icon="Tick01FreeIcons"
                  color="#6F7177"
                />
                <h6 class="text-[#1E212B] ml-2 text-sm">
                  Get paid online right from the invoice in a click
                </h6>
              </div>
            </div>
            <div class="flex w-full">
              <div class="flex mx-auto mt-2">
                <HugeiconsIcon
                  :icon="Tick01FreeIcons"
                  color="#6F7177"
                />
                <h6 class="text-[#1E212B] ml-2 text-sm">
                  Easily create and send invoices to your clients
                </h6>
              </div>
            </div>
            <div class="flex w-full">
              <div class="flex mx-auto mt-2">
                <HugeiconsIcon
                  :icon="Tick01FreeIcons"
                  color="#6F7177"
                />
                <h6 class="text-[#1E212B] ml-2 text-sm">
                  Collect recurring payments for ongoing projects
                </h6>
              </div>
            </div>
            <div class="flex items-center justify-center w-full gap-x-6 mt-6">
              <NuxtLink
                :to="{ name: 'invoices-create-new-invoice' }"
                class="mr-6"
              >
                <GroBasicButton
                  color="primary"
                  size="xs"
                  class="w-max"
                  shape="custom"
                >
                  Get Started
                </GroBasicButton>
              </NuxtLink>
            </div>

          </div>
        </div>
      </template>
    </MainDashboard>
  </div>
</template>

<style scoped>
</style>

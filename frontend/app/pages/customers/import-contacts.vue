<script setup lang="ts">
import {computed, defineModel, ref, watch} from "vue";
import MainDashboard from "@/components/dashboard/main-dashboard.vue";
import {CheckmarkCircle02Icon, LessThanIcon, Tick01Icon} from "@hugeicons/core-free-icons";
import {HugeiconsIcon} from "@hugeicons/vue";
import PrepareImport from "@/components/page/import-customers/PrepareImport.vue";
import UploadFile from "@/components/page/import-customers/UploadFile.vue";
import CSVMatch from "@/components/page/import-customers/CSVMatch.vue";
import type CSVData from "@/composables/helpers/types/csv";
import CSVCustomerTags from "@/components/page/import-customers/CSVCustomerTags.vue";
import type {FetchCrmTagResponse} from "@/composables/dto/tag/tag";
import {buildCustomerPayloads} from "@/composables/helpers/customers/buildCustomerImportPayloads";
import {useFormAPI} from "@/composables/api/customer/forms/create";
import type {FormVersionResponse} from "@/composables/dto/customer/form/form";
import BasicModal from "@/components/modals/basic-modal.vue";
import GroBasicButton from "@/components/buttons/GroBasicButton.vue";

const step = ref<number>(1);
const failureCount = ref<number>(0);
const showSuccessModal = ref<boolean>(false);
const csvData = ref<CSVData | null>(null);
const uploadedFile = ref<File | null>(null);
const tags = defineModel<FetchCrmTagResponse[]>('tags', {
  default: () => []
});
const formVersion = ref<FormVersionResponse>()

const openModal = computed(() => {
  return step.value === 5;
});

const createCustomer = async () => {

  try {
    const fetchUserApiResponse = await useFormAPI().FetchForm()

    if (fetchUserApiResponse.success) {
      formVersion.value = fetchUserApiResponse.data?.data
    }
  } catch (error) {
    console.error('Error fetching form:', error)
  }
}

// When UploadFile emits a file
const handleFileUploaded = (file: File) => {
  uploadedFile.value = file;
};

const csvDataUpdated = (csvDataModel: CSVData | null) => {
  csvData.value = csvDataModel;
};

const customFieldSlugs = ref<string[]>([]);

const handleCustomFieldsChanged = (slugs: string[]) => {
  customFieldSlugs.value = slugs;
};

const createAllCustomers = async (payloads: any[]) => {
  // Run all requests in parallel
  const results = await Promise.all(
    payloads.map(async (payload) => {
      const formData = {
        type: "customer",
        answers: payload.answers,
        versionSnap: payload.versionSnap,
      };

      try {
        const createdCustomer = await useFormAPI().CreateCustomer(
          formVersion.value?.version.formID,
          formVersion.value?.version.SK.split("#")[1],
          formData
        );
        if (createdCustomer.success) {
          successCount.value = successCount.value + 1
        }

        return {
          payload,
          success: createdCustomer?.success === true,
        };
      } catch (error) {
        console.error("Error creating customer:", error);
        failureCount.value = failureCount.value + 1
        return { payload, success: false };
      }
    })
  );

  // Summarize results
  const total = results.length;
  const successCount = results.filter((r) => r.success).length;
  const failedCount = total - successCount;

  showSuccessModal.value = true;

  return {
    total,
    successCount,
    failedCount,
    allSucceeded: failedCount === 0,
    results,
  };
};

watch(step, async (newVal) => {
  if (newVal === 5) {
    await createCustomer();
    const createCustomersPayload = buildCustomerPayloads(formVersion.value, csvData.value?.data, tags.value, customFieldSlugs.value)
    console.log(createCustomersPayload);
    await createAllCustomers(createCustomersPayload)
  }
});
</script>

<template>
  <MainDashboard current="Customers">
    <template #title>
      <div class="max-w-7xl mx-auto">
        <div class="flex flex-wrap gap-4">
          <div class="flex mb-3">
            <HugeiconsIcon
              :icon="LessThanIcon"
              size="13"
              class="my-auto mr-1"
              color="currentColor"
            />
            <NuxtLink
              to="/customers"
              class="text-[#2176AE] text-sm font-medium inline-flex items-center border-b border-[#2176AE]"
            >
              Customers
            </NuxtLink>
          </div>
        </div>
      </div>
    </template>

    <template #body>
      <!-- Step 1: Prepare -->
      <PrepareImport
        v-if="step === 1"
        v-model="step"
      />

      <!-- Step 2: Upload -->
      <UploadFile
        v-if="step === 2"
        v-model="step"
        @file-uploaded="handleFileUploaded"
      />

      <!-- Step 3: CSV Match -->
      <CSVMatch
        v-if="step === 3"
        v-model:step-model="step"
        :file="uploadedFile"
        @csv-data-changed="csvDataUpdated"
        @custom-fields-changed="handleCustomFieldsChanged"
      />

      <!-- Step 3: Tag selections -->
      <CSVCustomerTags
        v-if="step === 4"
        v-model:step-model="step"
        v-model:tags="tags"
      />

      <BasicModal
        v-model="showSuccessModal"
        size="xs"
      >
        <template #default>
          <div>
            <div class="py-4 flex items-center justify-center">
              <HugeiconsIcon
                size="100"
                :icon="CheckmarkCircle02Icon"
                color="white"
                fill="#00916E"
                class="ml-1 cursor-help"
              />
            </div>
            <h6 class="text-center w-full font-semibold text-[#4B4D55] text-xl">
              Contacts imported Successfully
            </h6>
            <div class="my-4">
              <div class="flex w-full items-center justify-center">
                <HugeiconsIcon
                  size="24"
                  :icon="Tick01Icon"
                  color="#4B4D55"
                  class="ml-1 cursor-help"
                />
                <h6 class="text-sm text-[#4B4D55]">
                  {{ csvData?.data.length }} new contacts were imported
                </h6>
              </div>
            </div>
            <NuxtLink
              to="/customers"
              class="cursor-pointer"
            >
              <GroBasicButton
                color="primary"
                size="xs"
                shape="custom"
              >
                Continue to Customers
              </GroBasicButton>
            </NuxtLink>
          </div>
        </template>
      </BasicModal>
    </template>
  </MainDashboard>
</template>

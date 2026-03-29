<script setup lang="ts">

import {GreaterThanIcon, Tick01FreeIcons} from "@hugeicons/core-free-icons";
import {HugeiconsIcon} from "@hugeicons/vue";
import ImportCustomerTemplate from "@/components/page/import-customers/ImportCustomerTemplate.vue";
import GroBasicButton from "@/components/buttons/GroBasicButton.vue";
import {defineModel} from "vue";

const model = defineModel<number>();
</script>

<template>
  <div>
    <ImportCustomerTemplate>
      <template #header>
        <div class="flex gap-x-4">
          <h6 class="text-[#1E212B] font-semibold text-sm my-auto">
            Prepare
          </h6>
          <HugeiconsIcon
            :icon="GreaterThanIcon"
            class="my-auto"
            :size="16"
            color="#939499"
            :stroke-width="3"
          />
          <h6 class="text-[#939499] font-semibold text-sm my-auto">
            Upload
          </h6>
          <HugeiconsIcon
            :icon="GreaterThanIcon"
            class="my-auto"
            :size="16"
            color="#939499"
            :stroke-width="3"
          />
          <h6 class="text-[#939499] font-semibold text-sm my-auto">
            Match
          </h6>
          <HugeiconsIcon
            :icon="GreaterThanIcon"
            class="my-auto"
            :size="16"
            color="#939499"
            :stroke-width="3"
          />
          <h6 class="text-[#939499] font-semibold text-sm my-auto">
            Label
          </h6>
        </div>
      </template>
      <template #body>
        <div class="p-4 md:p-8">
          <h4 class="font-semibold text-[#1E212B] text-lg">
            Prepare your file for import
          </h4>
          <div class="mt-6 space-y-3">
            <div class="flex ">
              <HugeiconsIcon
                :icon="Tick01FreeIcons"
                color="#6F7177"
              />
              <a
                href="https://res.cloudinary.com/sdq121/raw/upload/v1761901062/usegro-customer-import-template.csv"
                download
                class="text-[#2176AE] ml-2 text-sm"
              >
                Download a sample CSV file
              </a>
            </div>
            <div class="flex ">
              <HugeiconsIcon
                :icon="Tick01FreeIcons"
                color="#6F7177"
              />
              <h6 class="text-[#1E212B] ml-2 text-sm">
                Add field titles in the first row. e.g First name, Email. etc
              </h6>
            </div>
            <div class="flex ">
              <HugeiconsIcon
                :icon="Tick01FreeIcons"
                color="#6F7177"
              />
              <h6 class="text-[#1E212B] ml-2 text-sm">
                Make sure all contacts include a phone number or email address
              </h6>
            </div>
          </div>
          <div class="border border-[#EDEDEE] w-3/12 mt-10" />
          <div class="flex gap-x-2 mt-9 text-sm">
            <h5>
              Not sure how to create a CSV file?
            </h5>
            <h5 class="text-[#2176AE] cursor-pointer">
              Learn how
            </h5>
          </div>
          <div class="flex w-full mt-32">
            <GroBasicButton
              color="primary"
              size="sm"
              class="w-max ml-auto"
              shape="custom"
              @click="model = model + 1"
            >
              <template #default>
                Next
              </template>
            </GroBasicButton>
          </div>
        </div>
      </template>
    </ImportCustomerTemplate>
  </div>
</template>

<style scoped>

</style>

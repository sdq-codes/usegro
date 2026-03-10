<script setup lang="ts">

import {CheckmarkCircle01Icon, GreaterThanIcon} from "@hugeicons/core-free-icons";
import {HugeiconsIcon} from "@hugeicons/vue";
import ImportCustomerTemplate from "@/components/page/import-customers/ImportCustomerTemplate.vue";
import GroBasicButton from "@/components/buttons/GroBasicButton.vue";
import {defineModel} from "vue";
import GroBasicTagSelect from "@/components/forms/select/GroBasicTagSelect.vue";
import type {FetchCrmTagResponse} from "@/composables/dto/tag/tag";

const tags = defineModel<FetchCrmTagResponse[]>('tags', {
  default: () => []
});

const stepModel = defineModel<number>("stepModel");
</script>

<template>
  <div>
    <ImportCustomerTemplate>
      <template #header>
        <div class="flex gap-x-4">
          <div class="flex">
            <HugeiconsIcon
              :icon="CheckmarkCircle01Icon"
              class="my-auto"
              :size="21"
              fill="#007458"
              color="white"
              :stroke-width="3"
            />
            <h6 class="text-[#007458] font-semibold text-sm my-auto">
              Prepare
            </h6>
          </div>
          <HugeiconsIcon
            :icon="GreaterThanIcon"
            class="my-auto"
            :size="16"
            color="#939499"
            :stroke-width="3"
          />
          <div class="flex">
            <HugeiconsIcon
              :icon="CheckmarkCircle01Icon"
              class="my-auto"
              :size="21"
              fill="#007458"
              color="white"
              :stroke-width="3"
            />
            <h6 class="text-[#007458] font-semibold text-sm my-auto">
              Upload
            </h6>
          </div>
          <HugeiconsIcon
            :icon="GreaterThanIcon"
            class="my-auto"
            :size="16"
            color="#939499"
            :stroke-width="3"
          />
          <div class="flex">
            <HugeiconsIcon
              :icon="CheckmarkCircle01Icon"
              class="my-auto"
              :size="21"
              fill="#007458"
              color="white"
              :stroke-width="3"
            />
            <h6 class="text-[#007458] font-semibold text-sm my-auto">
              Match
            </h6>
          </div>
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
          <h4 class="font-semibold text-[#1E212B] text-lg mt-7">
            Label your contacts (Optional)
          </h4>
          <div class="mt-3 space-y-3">
            <p class="text-sm text-[#4B4D55]">
              Labels are a great way to organize your contacts into categories and mailing lists
            </p>
          </div>
          <div class="mt-10">
            <div>
              <GroBasicTagSelect v-model="tags" />
            </div>
          </div>
          <div class="flex w-full mt-32">
            <GroBasicButton
              color="tertiary"
              size="sm"
              class="w-max"
              shape="custom"
              @click="stepModel = stepModel - 1"
            >
              <template #frontIcon>
                <HugeiconsIcon
                  :icon="GreaterThanIcon"
                  size="12"
                  :stroke-width="3"
                  class="my-auto rotate-180"
                />
              </template>
              <template #default>
                Back
              </template>
            </GroBasicButton>
            <GroBasicButton
              color="primary"
              size="sm"
              class="w-max ml-auto"
              shape="custom"
              @click="stepModel = stepModel + 1"
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

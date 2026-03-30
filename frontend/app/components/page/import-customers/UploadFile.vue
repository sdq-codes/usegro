<script setup lang="ts">
import {defineEmits, defineModel, ref} from 'vue';
import {
  CheckmarkCircle01Icon, FileIconFreeIcons,
  GreaterThanIcon,
  LessThanIcon,
  PlusSignIcon
} from "@hugeicons/core-free-icons";
import {HugeiconsIcon} from "@hugeicons/vue";
import ImportCustomerTemplate from "@/components/page/import-customers/ImportCustomerTemplate.vue";
import GroBasicButton from "@/components/buttons/GroBasicButton.vue";

const fileInput = ref<HTMLInputElement | null>(null);
const selectedFile = ref<File | null>(null);
const isDragging = ref(false);

const model = defineModel<number>();

const emit = defineEmits<{
  (e: 'file-uploaded', value): void
}>()

const handleFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement;
  if (target.files && target.files.length > 0) {
    selectedFile.value = target.files[0];
    emit('file-uploaded', target.files[0])
  }
};

const handleDragOver = (event: DragEvent) => {
  event.preventDefault();
  isDragging.value = true;
};

const handleDragLeave = () => {
  isDragging.value = false;
};

const handleDrop = (event: DragEvent) => {
  event.preventDefault();
  isDragging.value = false;

  if (event.dataTransfer?.files && event.dataTransfer.files.length > 0) {
    selectedFile.value = event.dataTransfer.files[0];
  }
};

const triggerFileInput = () => {
  fileInput.value?.click();
};

const removeFile = () => {
  selectedFile.value = null;
  if (fileInput.value) {
    fileInput.value.value = '';
  }
};
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
          <h6 class="text-[#1E212B] font-semibold text-sm my-auto">
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
        <div class="p-4 md:p-6">
          <h4 class="font-semibold text-[#1E212B] text-lg mt-2">
            Upload your file
          </h4>
          <div class="mt-6 space-y-3">
            <input
              ref="fileInput"
              type="file"
              accept=".csv,.xlsx,.xls"
              class="hidden"
              @change="handleFileSelect"
            >

            <div
              v-if="!selectedFile"
              :class="[
                'border-[#4D91BE] p-4 md:px-12 md:py-24 border-2 border-dashed rounded-2xl cursor-pointer transition-colors',
                isDragging ? 'bg-blue-50' : 'hover:bg-gray-50'
              ]"
              @click="triggerFileInput"
              @dragover="handleDragOver"
              @dragleave="handleDragLeave"
              @drop="handleDrop"
            >
              <div class="flex justify-center">
                <div class="">
                  <HugeiconsIcon
                    :icon="PlusSignIcon"
                    class="mx-auto"
                    :size="80"
                    color="#4D91BE"
                    :stroke-width="1"
                  />
                  <h6 class="text-[#1E212B] font-semibold text-md my-auto text-center mt-3">
                    Drag your file or click to <span class="text-[#2176AE]">browse a file</span>
                  </h6>
                  <h6 class="text-[#6F7177] text-center mt-1 font-normal">
                    Note: Existing contacts will update to match your csv.
                  </h6>
                </div>
              </div>
            </div>
            <div v-else>
              <div
                class="border-[#4D91BE] border-dashed bg-[#EDF4F9] p-4 md:px-12 md:py-8 border-2 rounded-2xl"
              >
                <div class="flex items-center justify-between">
                  <div class="flex items-center gap-3">
                    <HugeiconsIcon
                      :icon="FileIconFreeIcons"
                      :size="32"
                      color="#1E212B"
                      :stroke-width="1"
                    />
                    <div>
                      <h6 class="text-[#1E212B] font-semibold text-md">
                        {{ selectedFile.name }}
                      </h6>
                      <p class="text-[#6F7177] text-sm">
                        {{ (selectedFile.size / 1024).toFixed(2) }} KB
                      </p>
                    </div>
                  </div>
                  <button
                    class="text-[#6F7177] hover:text-[#1E212B] font-medium text-sm underline"
                    @click="removeFile"
                  >
                    Remove
                  </button>
                </div>
              </div>
              <div class="mt-5">
                <h6 class="text-[#4B4D55] text-md">
                  Email subscribers are contacts who agreed to be added to your mailing list. Before marking them as subscribers and sending email campaigns, make sure they’ve agreed to recieve emails.
                  <span class="text-[#2176AE] cursor-pointer">Learn more</span>
                </h6>
              </div>
            </div>
          </div>

          <div class="flex w-full mt-32">
            <GroBasicButton
              color="tertiary"
              size="sm"
              class="w-max"
              shape="custom"
              @click="model = model - 1"
            >
              <template #frontIcon>
                <HugeiconsIcon
                  :icon="LessThanIcon"
                  size="12"
                  :stroke-width="3"
                  class="my-auto"
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
              :disabled="!selectedFile"
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

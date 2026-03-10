<script setup>
import { ref } from 'vue';
import {HugeiconsIcon} from "@hugeicons/vue";
import {Upload02FreeIcons} from "@hugeicons/core-free-icons";

const dragActive = ref(false);
const files = ref([]);
const fileInput = ref(null);

const handleDrag = (e) => {
  if (e.type === "dragenter" || e.type === "dragover") {
    dragActive.value = true;
  } else if (e.type === "dragleave") {
    dragActive.value = false;
  }
};

const handleDrop = (e) => {
  dragActive.value = false;

  if (e.dataTransfer.files && e.dataTransfer.files[0]) {
    handleFiles(e.dataTransfer.files);
  }
};

const handleChange = (e) => {
  if (e.target.files && e.target.files[0]) {
    handleFiles(e.target.files);
  }
};

const handleFiles = (fileList) => {
  const newFiles = Array.from(fileList);
  files.value = [...files.value, ...newFiles];
};

const onButtonClick = () => {
  fileInput.value.click();
};

const removeFile = (index) => {
  files.value = files.value.filter((_, i) => i !== index);
};

const formatFileSize = (bytes) => {
  return (bytes / 1024).toFixed(2) + ' KB';
};
</script>

<template>
  <div class="">
    <div
      @click="onButtonClick"
      :class="[
        dragActive ? 'border-blue-500 bg-blue-50' : 'border-gray-300 bg-white'
      ]"
      @dragenter.prevent="handleDrag"
      @dragleave.prevent="handleDrag"
      @dragover.prevent="handleDrag"
      @drop.prevent="handleDrop"
    >
      <input
        ref="fileInput"
        type="file"
        class="hidden"
        multiple
        accept=".jpg,.jpeg,.gif,.png,.pdf"
        @change="handleChange"
      >

      <div class="flex">
        <div class="my-auto">
          <HugeiconsIcon :icon="Upload02FreeIcons" color="#2176AE" />
        </div>
        <div
          class="text-sm font-semibold ml-3 my-auto text-[#2176AE]"
        >
          Upload files
        </div>
      </div>
      <div class="flex">
        <p class="text-black text-xs mt-2">
          JPG, GIF, PNG, PDF  Up to 3 foles, xMB per file
        </p>
      </div>
    </div>

    <div
      v-if="files.length > 0"
      class="mt-6 bg-white rounded-lg border border-gray-200 p-4"
    >
      <h3 class="font-semibold text-gray-900 mb-3">
        Uploaded Files ({{ files.length }})
      </h3>
      <ul class="space-y-2">
        <li
          v-for="(file, index) in files"
          :key="index"
          class="flex items-center justify-between p-3 bg-gray-50 rounded hover:bg-gray-100"
        >
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium text-gray-900 truncate">
              {{ file.name }}
            </p>
            <p class="text-xs text-gray-500">
              {{ formatFileSize(file.size) }}
            </p>
          </div>
          <button
            class="ml-4 text-red-600 hover:text-red-800 text-sm font-medium"
            @click="removeFile(index)"
          >
            Remove
          </button>
        </li>
      </ul>
    </div>
  </div>
</template>

<style scoped>
/* Add any additional custom styles here if needed */
</style>

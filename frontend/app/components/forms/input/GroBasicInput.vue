<script lang="ts" setup>
import {defineProps, computed, defineModel} from 'vue'
import { HugeiconsIcon } from '@hugeicons/vue'
import {CheckmarkSquare01Icon, InformationCircleIcon, AlertSquareIcon, CancelSquareIcon, Loading03Icon} from '@hugeicons/core-free-icons'

const props = withDefaults(defineProps<{
  placeholder?: string,
  hint?: string,
  disabled?: boolean,
  info?: string,
  loading?: boolean,
  color?: 'primary' | 'success' | 'error',
  readonly?: boolean,
  type?: string,
}>(), {
  type: 'text',
  color: 'primary',
  disabled: false,
  placeholder: "",
  hint: "",
  info: "",
})

const model = defineModel<string>();

// Base classes for input
const inputClass = computed(() => {
  return [
    'w-full rounded-lg border border-solid px-3 py-2 pr-10 transition-all duration-200 outline-none bg-[#F6F6F7] focus:bg-[#FFFFFF]',
    'disabled:bg-[#DBDBDD] disabled:text-[#AFB5B8] disabled:border-0 disabled:cursor-not-allowed',
    colorClasses[props.color || 'primary'],
  ].join(' ')
})

const colorClasses: Record<string, string> = {
  primary: 'border-[#EDEDEE] hover:border-[#94BDD8] focus:border-[#1E212B] text-[#4F5435] placeholder-[#939499]',
  success: 'border-[#007458] focus:bg-[#FFFFFF] text-[#4F5435] placeholder-[#939499]',
  error: 'border-[#703425] focus:bg-[#FFFFFF] text-[#4F5435] placeholder-[#939499]',
  custom: ''
}

const labelClasses: Record<string, string> = {
  primary: 'text-[#1E212B]',
  success: 'text-[#00916E]',
  error: 'text-[#AF513A]',
}

const hintClasses: Record<string, string> = {
  primary: 'text-[#1E212B]',
  success: 'text-[#00916E]',
  error: 'text-[#AF513A]',
}

const labelClass = computed(() => {
  return [
    'text-xs font-medium mb-1',
    labelClasses[props.color || 'primary'],
  ].join(' ')
})

const hintClass = computed(() => {
  return [
    hintClasses[props.color || 'primary'],
  ].join(' ')
})

// Wrapper for positioning trailing icon
const wrapperClass = 'relative flex flex-col gap-1'

</script>

<template>
  <div class="w-full">
    <div
      :class="labelClass"
      class="flex items-center relative"
    >
      <slot />

      <div
        v-if="props?.info"
        class="relative group"
      >
        <HugeiconsIcon
          size="18"
          :icon="InformationCircleIcon"
          color="#1E212B"
          class="ml-1 cursor-help"
        />
        <!-- Tooltip -->
        <div
          class="absolute mt-2 w-96
             bg-gray-800 text-white text-xs rounded-md px-2 py-1
             hidden group-hover:flex  group-active:flex transition-opacity"
        >
          <div class="absolute bottom-0 hidden w-full items-center mb-5 group-hover:flex group-active:flex">
            <span class="relative rounded-md z-10 p-4 text-xs leading-none text-white whitespace-no-wrap bg-black shadow-lg">
              {{ props.info }}
            </span>
          </div>
        </div>
      </div>
    </div>


    <!-- Input wrapper with trailing icon -->
    <div :class="wrapperClass">
      <input
        v-model="model"
        :type="type"
        :placeholder="placeholder"
        :disabled="props?.disabled"
        :class="inputClass"
        :readonly="props.readonly"
      >
      <!-- Trailing Icon -->
      <div class="absolute inset-y-0 right-3 flex items-center cursor-pointer">
        <svg
          v-if="!props.disabled && !props.loading && model?.length >0 && props.type !== 'date'"
          width="16"
          height="16"
          viewBox="0 0 12 12"
          class="my-auto cursor-pointer"
          @click="model = ''"
        >
          <path
            d="M6.03369 0.020813C7.31152 0.0208052 8.3155 0.0209263 9.09912 0.126282C9.90193 0.234254 10.539 0.459831 11.0396 0.960266C11.5402 1.46088 11.7656 2.09862 11.8735 2.90167C11.9788 3.68519 11.979 4.6886 11.979 5.96613V6.03351C11.979 7.31134 11.9789 8.31531 11.8735 9.09894C11.7655 9.90182 11.5401 10.5388 11.0396 11.0394C10.539 11.5399 9.902 11.7654 9.09912 11.8734C8.3155 11.9787 7.31152 11.9788 6.03369 11.9788H5.96631C4.68879 11.9788 3.68537 11.9787 2.90186 11.8734C2.09881 11.7654 1.46106 11.54 0.960449 11.0394C0.460014 10.5388 0.234437 9.90175 0.126465 9.09894C0.0211094 8.31531 0.0209883 7.31134 0.0209961 6.03351V5.9671C0.0209883 4.68937 0.021148 3.68529 0.126465 2.90167C0.234432 2.09862 0.459834 1.46088 0.960449 0.960266C1.46106 0.459651 2.09881 0.234249 2.90186 0.126282C3.68547 0.0209649 4.68955 0.0208052 5.96729 0.020813H6.03369ZM8.1626 3.83722C7.93482 3.60954 7.56519 3.60953 7.3374 3.83722L5.99951 5.17511L4.6626 3.8382C4.43478 3.6104 4.0652 3.61038 3.8374 3.8382C3.60995 4.06596 3.6099 4.43469 3.8374 4.66241L5.17529 6.00031L3.8374 7.33722C3.60963 7.56502 3.60962 7.93461 3.8374 8.16241C4.0652 8.39019 4.4348 8.39019 4.6626 8.16241L5.99951 6.82452L7.3374 8.16241C7.56522 8.39021 7.9348 8.39023 8.1626 8.16241C8.39032 7.93459 8.39039 7.56499 8.1626 7.33722L6.82471 6.00031L8.1626 4.66241C8.39031 4.43461 8.39036 4.065 8.1626 3.83722Z"
            fill="#939499"
          />
        </svg>
        <HugeiconsIcon
          v-if="props?.loading"
          :icon="Loading03Icon"
          color="#939499"
          class="animate-spin"
        />
      </div>
    </div>

    <!-- Hint with icon -->
    <div
      v-if="props.hint"
      class="flex items-center gap-1 mt-1 text-xs"
      :class="hintClass"
    >
      <HugeiconsIcon
        v-if="props.hint && props.color === 'success'"
        color="#FFFFFF"
        fill="#00916E"
        :icon="CheckmarkSquare01Icon"
      />
      <HugeiconsIcon
        v-else-if="props.hint && props.color === 'error'"
        color="#FFFFFF"
        fill="#AF513A"
        :icon="CancelSquareIcon"
      />
      <HugeiconsIcon
        v-else
        :icon="AlertSquareIcon"
        fill="#939499"
        color="#FFFFFF"
      />

      <span>{{ props.hint }}</span>
    </div>
  </div>
</template>

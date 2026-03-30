<script lang="ts" setup>
import {defineProps, computed, defineModel, defineEmits} from 'vue'
import { HugeiconsIcon } from '@hugeicons/vue'
import {
  CheckmarkSquare01Icon,
  AlertSquareIcon,
  CancelSquareIcon,
  Search01Icon,
  PlusSignIcon
} from '@hugeicons/core-free-icons'
import Multiselect from "vue-multiselect"
import "vue-multiselect/dist/vue-multiselect.css"

interface Props {
  placeholder?: string
  hint?: string
  addMoreText?: string
  disabled?: boolean
  loading?: boolean
  color?: 'primary' | 'success' | 'error'
  readonly?: boolean
  options: { value: string | number; label: string; description?: string }[]
  modelValue?: string | number | null,
  multipleSelect?: boolean,
  showInitials?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  multipleSelect: false,
  showInitials: false,
  multiple: false,
  placeholder: '',
  hint: '',
  addMoreText: '',
  color: 'primary',
  modelValue: ''
})

const getInitials = (label: string) => {
  return label
    .split(' ')
    .map((n) => n[0])
    .join('')
    .toUpperCase()
    .substring(0, 2);
}

const emit = defineEmits<{
  (e: 'add-new'): void
}>()


function onAddNew() {
  emit('add-new')
}

const model = defineModel<string | number | Array<any> | null>();

// Normalize between plain value ↔ full option object for vue-multiselect
const internalModel = computed({
  get() {
    if (props.multipleSelect) {
      const vals = model.value as (string | number)[] | undefined
      return vals?.map(v => props.options.find(o => o.value === v)).filter(Boolean) ?? []
    }
    return props.options.find(o => o.value === model.value) ?? null
  },
  set(val: any) {
    if (props.multipleSelect) {
      model.value = val?.map((o: any) => o.value) ?? []
    } else {
      model.value = val?.value ?? null
    }
  }
})

// Base classes for select
const selectClass = computed(() => {
  return [
    'w-full rounded-lg duration-200 outline-none appearance-none',
    'disabled:bg-[#DBDBDD] disabled:text-[#AFB5B8] disabled:cursor-not-allowed',
    colorClasses[props.color || 'primary'],
  ].join(' ')
})

const colorClasses: Record<string, string> = {
  primary: 'text-[#4F5435] placeholder-[#939499]',
  success: 'focus:bg-[#FFFFFF] text-[#4F5435] placeholder-[#939499]',
  error: 'focus:bg-[#FFFFFF] text-[#4F5435] placeholder-[#939499]',
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
    <div :class="labelClass">
      <slot />
    </div>

    <!-- Select wrapper with trailing icon -->
    <div :class="wrapperClass">
      <Multiselect
        v-model="internalModel"
        :options="props.options"
        :placeholder="placeholder"
        :searchable="true"
        :multiple="multipleSelect"
        track-by="value"
        label="label"
        select-label=""
        selected-label=""
        deselect-label=""
        :class="selectClass"
        class="rounded-lg hover:border-[#94BDD8]"
      >
        <template #option="{ option }">
          <div class="flex items-center gap-2">
            <div
              v-if="props.showInitials"
              class="flex items-center justify-center w-6 h-6 rounded-full bg-[#D26B06] text-white text-[10px] font-bold shrink-0"
            >
              {{ getInitials(option.label) }}
            </div>
            <div class="flex flex-col min-w-0">
              <span class="font-semibold text-sm whitespace-normal break-words">{{ option.label }}</span>
              <span v-if="option.description" class="text-xs text-[#939499] font-normal mt-0.5 whitespace-normal break-words">{{ option.description }}</span>
            </div>
          </div>
        </template>

        <template #singleLabel="{ option }">
          <div class="flex items-center gap-2">
            <div
              v-if="props.showInitials"
              class="flex items-center justify-center w-5 h-5 rounded-full bg-[#4F5435] text-white text-[9px] font-bold"
            >
              {{ getInitials(option.label) }}
            </div>
            <span>{{ option.label }}</span>
          </div>
        </template>

        <template #tag="{ option, remove }">
          <span class="multiselect__tag">
            <div class="flex items-center gap-1">
              <div
                v-if="props.showInitials"
                class="flex items-center justify-center w-4 h-4 rounded-full bg-white/20 text-white text-[8px] font-bold"
              >
                {{ getInitials(option.label) }}
              </div>
              <span>{{ option.label }}</span>
            </div>
            <i
              class="multiselect__tag-icon"
              @mousedown.prevent="remove(option)"
            />
          </span>
        </template>

        <template
          v-if="addMoreText != ''"
          #beforeList
        >
          <button
            class="w-full py-3 flex px-3 py-2 text-left text-sm font-semibold cursor-pointer text-[#2176AE] hover:bg-[#F6F6F7] border-t border-[#DBDBDD] transition-colors"
            @click.prevent="onAddNew"
          >
            <HugeiconsIcon
              :icon="PlusSignIcon"
              color="#2176AE"
              class="h-4"
              stroke-width="3"
            />
            {{ addMoreText }}
          </button>
        </template>
      </Multiselect>
    </div>

    <!-- Hint with icon -->
    <div
      v-if="props.hint"
      class="flex items-center gap-1 mt-1 text-xs"
      :class="hintClass"
    >
      <HugeiconsIcon
        v-if="props.color === 'success'"
        color="#FFFFFF"
        fill="#00916E"
        :icon="CheckmarkSquare01Icon"
      />
      <HugeiconsIcon
        v-else-if="props.color === 'error'"
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

<style>
.multiselect__tags:hover{
  border-color: #94BDD8 !important;
}

.multiselect__tags{
  border-radius: 8px !important;
  background-color: #F6F6F7 !important;
}

.multiselect__input {
  background-color: #F6F6F7 !important;
}

.multiselect__tags::placeholder {
  color: #939499 !important;
}

.multiselect__option--highlight {
  background-color: #DBDBDD !important;
  color: #000000 !important;
}

.multiselect__content-wrapper {
  box-shadow: 0px 8px 24px 0px rgba(78, 77, 89, 0.12) !important;
  border: 0 !important;
}

.multiselect__single {
  background-color: #F6F6F7 !important;
}

.multiselect__tag {
  background-color: #af513a !important;
  font-weight: 400;
}

.multiselect__tag-icon::after {
  color: #FFFFFF !important;
}
</style>

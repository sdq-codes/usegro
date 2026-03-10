<template>
  <!-- Background overlay -->
  <div class="" />
  <div
    v-if="model"
    class="fixed inset-0 bg-gradient-to-b from-black/10 to-black/0 backdrop-blur-[2px] z-50 flex items-center justify-center bg-black/50"
  >
    <!-- Modal container -->
    <div
      :class="modalClass"
      class="bg-white absolute rounded-xl shadow-lg w-[90vw] md:w-full"
    >
      <div
        v-if="slots.title"
        class="header h-[10%] p-6 "
      >
        <div class="flex w-full">
          <h6 class="text-xl font-semibold my-auto text-[#000000]">
            <slot name="title" />
          </h6>
          <button
            class="ml-auto cursor-pointer text-gray-500 hover:text-gray-800 text-2xl my-auto"
            @click="model = false"
          >
            ✕
          </button>
        </div>
      </div>
      <div class="body h-[80%] overflow-scroll relative px-6 py-3">
        <slot />
      </div>
      <div
        v-if="slots.footer"
        class="footer h-[10%] p-6"
      >
        <slot name="footer" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {computed, defineModel, defineProps, useSlots} from "vue";

const slots = useSlots();

const props = defineProps<{
  size: 'lg' | 'xl' | 'xs',
}>()

const model = defineModel<boolean>({ default: false });

const sizeClasses: Record<string, string> = {
  xs: 'max-w-sm h-max',
  xl: 'max-w-3xl h-[80vh]',
  lg: 'max-w-xl h-[80vh]'
}

const modalClass = computed(() => {
  return [
    sizeClasses[props.size || 'lg'],
  ].join(' ')
})
</script>

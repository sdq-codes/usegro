<script setup lang="ts">
const props = defineProps({
  url: {
    type: String,
    default: '',
  },
  current: {
    type: Boolean,
    default: false,
  },
  collapsed: {
    type: Boolean,
    default: false,
  }
});
</script>

<template>
  <NuxtLink
    :to="{ name: props.url }"
    class="flex cursor-pointer md:justify-center lg:justify-start w-full rounded-lg px-3 py-2 mb-1 hover:bg-gray-300"
  >
    <slot name="icon" />
    <div
      v-if="!collapsed"
      class=" ml-2 my-auto text-sm md:hidden lg:block mr-auto"
    >
      <slot />
    </div>
    <small
      v-if="$slots.count && !collapsed"
      class="px-2 py-[2px] my-auto text-xs text-[#6F7177] rounded-xl font-normal bg-[#DBDBDD] md:hidden lg:block"
    >
      <slot name="count" />
    </small>
    <small
      v-if="$slots.label && !collapsed"
      class="ml-2 px-2 my-auto text- text-[#3C3F28] rounded-lg bg-[#EAEBE3] md:hidden lg:block"
    >
      <slot name="label" />
    </small>
  </NuxtLink>
</template>

<style scoped>

</style>

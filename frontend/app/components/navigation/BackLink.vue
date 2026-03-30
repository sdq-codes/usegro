<script setup lang="ts">
import { LessThanIcon } from "@hugeicons/core-free-icons"
import { HugeiconsIcon } from "@hugeicons/vue"

defineProps<{
  to: string
  label: string
}>()
</script>

<template>
  <div class="flex mb-3">
    <HugeiconsIcon
      :icon="LessThanIcon"
      size="13"
      class="my-auto mr-1"
      color="currentColor"
    />
    <NuxtLink
      :to="to"
      class="text-[#2176AE] text-xs font-medium inline-flex items-center border-b border-[#2176AE]"
    >
      {{ label }}
    </NuxtLink>
  </div>
</template>

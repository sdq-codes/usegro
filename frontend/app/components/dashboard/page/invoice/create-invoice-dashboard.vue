<script setup lang="ts">
import CollapsedDashboard from "@/components/dashboard/collapsed-dashboard.vue";
import {Cancel01Icon} from "@hugeicons/core-free-icons";
import {HugeiconsIcon} from "@hugeicons/vue";
</script>

<template>
  <CollapsedDashboard current="Invoices">
    <template #body>
      <div class="bg-white md:fixed z-10 flex md:w-[96vw] py-4 px-7 gap-x-5 md:h-20">
        <HugeiconsIcon
          :icon="Cancel01Icon"
          :size="32"
          color="#1E212B"
          class="my-auto"
        />
        <div class="my-auto b-l border-l border-l-[#EDEDEE] font-semibold text-[#1E212B] pl-5">
          <slot name="title"/>
        </div>
        <div class="ml-auto  my-auto">
          <slot name="menu"/>
        </div>
      </div>
      <div class="p-5 md:p-10 mt-16 flex-1">
        <slot name="body"/>
      </div>
    </template>
  </CollapsedDashboard>
</template>

<style scoped>
</style>

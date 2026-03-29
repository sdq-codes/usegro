<script setup lang="ts">
import SideBar from "@/components/sidebar/side-bar.vue";

const props = defineProps({
  current: {
    type: String,
    default: "",
  },
});
</script>

<template>
  <div class="flex h-[100vh] max-w-screen">
    <side-bar
      class=""
      :current="props.current"
    />
      <div class="bg-[#fafafa] overflow-y-scroll w-full py-4 pt-20 pb-20 md:pt-4 md:pb-4">
      <div class="md:w-full lg:max-w-7xl lg:mx-auto px-4 md:px-12 md:py-8">
        <slot name="title" />
        <div class="block ">
          <slot name="body" />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>

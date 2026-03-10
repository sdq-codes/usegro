<script setup lang="ts">
import SideBar from "@/components/sidebar/side-bar.vue";

const props = defineProps({
  current: {
    type: String,
    default: "",
  }
});
</script>

<template>
  <div class="flex h-[100vh] max-w-screen">
    <side-bar
      class=""
      :current="props.current"
      :collapsed="true"
    />
    <div class="bg-[#F5F5F5] overflow-y-scroll w-full">
      <div class="md:w-full ">
        <div class="flex flex-col h-screen">
          <slot name="body" />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>

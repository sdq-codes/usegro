<template>
  <!-- Sidebar -->
  <div class="flex h-screen">
    <side-bar class="" />
    <div class="bg-[#F5F5F5] w-full">
    </div>
  </div>
</template>

<script>
import SideBar from "@/components/sidebar/side-bar.vue";

export default {
  components: {SideBar},
  methods: {
    toggle() {
      this.open = !this.open;
    },
  },
};
</script>

<style scoped>

</style>

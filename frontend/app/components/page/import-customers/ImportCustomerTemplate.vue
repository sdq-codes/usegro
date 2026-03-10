<script setup lang="ts">

</script>

<template>
  <div class="max-w-7xl mx-auto">
    <div class="bg-white rounded-tl-3xl rounded-tr-3xl mt-3 border border-[#EDEDEE]">
      <div class="flex p-4 md:p-8">
        <h2 class="font-black text-2xl text-[#1E212B]">
          Import Contacts
        </h2>
        <div class="hidden md:flex ml-auto gap-x-3">
          <slot name="header" />
        </div>
      </div>
    </div>
    <div class="bg-white rounded-bl-3xl rounded-br-3xl">
      <div class="">
        <slot name="body" />
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>

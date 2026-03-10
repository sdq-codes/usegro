<script lang="ts" setup>
</script>

<template>
  <section class="flex h-screen w-screen">
    <div class="w-[100%] md:w-[50%] bg-[#F6FAFC] px-6">
      <div class=" w-full px-4 py-6 md:py-12">
        <NuxtLink
          to="/"
          class="cursor-pointer"
        >
          <img
            src="https://res.cloudinary.com/sdq121/image/upload/v1755371764/ltumpn6lhocowaxdi68p.png"
            alt="Gro App Logo"
            class="h-8"
          >
        </NuxtLink>
      </div>
      <div class="w-full">
        <div class="w-full md:w-3/5 mx-auto mt-4 md:mt-20">
          <slot />
        </div>
      </div>
    </div>

    <!-- Right (40%) -->
    <div class="hidden md:block md:w-[50%] bg-white px-8 py-12">
      <div
        class="h-full w-full rounded-2xl bg-cover bg-center"
        style="background-image: url('https://res.cloudinary.com/sdq121/image/upload/v1755547137/yswl1zhrsn901swk3uyf.png')"
      />
    </div>
  </section>
</template>

<style>
</style>

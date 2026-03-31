<script lang="ts" setup>

</script>

<template>
  <section>
    <div class="container-fluid basic-page">
      <div class="px-10 pt-10">
        <div class="grid grid-cols-1 relative z-1 md:grid-cols-2">
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
      </div>
      <div class="px-20 absolute w-full top-12 right-0 pointer-events-none">
        <img
          src="https://res.cloudinary.com/sdq121/image/upload/v1755376170/idi6ql03xtnkyisewfac.png"
          alt="Cloud Image"
          class="w-full"
        >
      </div>
      <div class="flex mt-10 md:mt-28">
        <div class="mx-auto w-md px-4 md:px-0">
          <slot />
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.basic-page {
  background-color: #F6FAFC;
  width: 100vw;
  min-height: 100vh;
}
</style>

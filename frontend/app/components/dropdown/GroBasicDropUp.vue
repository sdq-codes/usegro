<script setup>
import { ref, onMounted, onBeforeUnmount } from "vue";

const open = ref(false);

function toggle() {
  open.value = !open.value;
}

function closeDropdown(e) {
  if (!e.target.closest(".dropdown-container")) {
    open.value = false;
  }
}

function handleClick() {
  open.value = !open.value;
}

onMounted(() => {
  document.addEventListener("click", closeDropdown);
});

onBeforeUnmount(() => {
  document.removeEventListener("click", closeDropdown);
});
</script>

<template>
  <div class="relative dropdown-container">
    <!-- Trigger -->
    <div @click="toggle">
      <slot name="button" />
    </div>

    <!-- Dropup Menu -->
    <Transition
      enter-active-class="transition ease-out duration-200 transform"
      enter-from-class="opacity-0 translate-y-2"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition ease-out duration-200"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0 translate-y-2"
    >
      <ul
        v-show="open"
        class="origin-bottom absolute bottom-full left-0 mb-2
               min-w-[180px] bg-white border border-slate-200 p-2
               rounded-xl shadow-xl
               before:content-[''] before:absolute before:top-full before:left-6
               before:border-[6px] before:border-transparent before:border-t-white
               before:-mt-px
               after:content-[''] after:absolute after:top-full after:left-6
               after:border-[7px] after:border-transparent after:border-t-slate-200
               text-xs z-50"
      >
        <div @click="handleClick">
          <slot name="menu-list" />
        </div>
      </ul>
    </Transition>
  </div>
</template>

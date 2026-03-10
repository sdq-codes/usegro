<template>
  <div>
    <div
      class="hidden md:flex  flex-col items-left md:2 h-full overflow-hidden text-[#4B4D55] font-medium bg-[#FFFFFF] border-1 border-[#EDEDEE] rounded"
      :class="[collapsed ? 'w-16' : 'md:w-16 lg:w-64']"
    >
      <a
        class="flex mt-8 h-12 md:px-3 lg:px-4"
        href="#"
      >
        <img
          src="https://res.cloudinary.com/sdq121/image/upload/v1755371764/ltumpn6lhocowaxdi68p.png"
          alt="Gro App Logo"
          class="h-6 lg:h-8 px-2 my-auto"
          :class="[collapsed ? 'hidden' : 'md:hidden lg:block']"
        >
        <img
          src="https://res.cloudinary.com/sdq121/image/upload/v1762074508/ouc2jezgplwvubsgsobv.svg"
          alt="Gro App Logo"
          width="100%"
          class="my-auto"
          :class="[collapsed ? '' : 'lg:hidden']"
        >
      </a>
      <div
        class="flex flex-col items-center mt-3 border-t border-[#EDEDEE] pt-5 pb-8"
        :class="[collapsed ? 'px-2' : 'px-2 lg:px-4']"
      >
        <SidebarMenu
          v-for="menu in SIDEBAR_MENU_SECTION_A"
          :key="menu.title"
          :collapsed="collapsed"
          :url="menu.url"
          :current="menu.title === props.current"
          :class="[
            menu.title === props.current ?
              'bg-[#EDEDEE] text-[#1E212B] font-semibold ' : 'font-normal']"
        >
          <template #default>
            {{ menu.title }}
          </template>
          <template #icon>
            <HugeiconsIcon
              :color=" menu.title === props.current ? '#EDEDEE' : '#FFFFFF'"
              class="w-6"
              :fill="menu.title === props.current ? '#AF513A' : '#B7B8BB'"
              :icon="menu.icon"
            />
          </template>
          <template
            v-if="menu.count"
            #count
          >
            {{ menu.count }}
          </template>
          <template
            v-if="menu.soon"
            #label
          >
            Soon
          </template>
        </SidebarMenu>
      </div>
      <h6
        v-if="!collapsed"
        class="text-left uppercase font-bold text-xs my-2 text-[#6F7177] px-3 lg:px-8 md:hidden lg:block"
      >
        Sales channels
      </h6>
      <div
        class="flex flex-col items-center pb-12"
        :class="[collapsed ? 'px-2' : 'px-2 lg:px-4']"
      >
        <SidebarMenu
          v-for="menu in SIDEBAR_MENU_SECTION_B"
          :key="menu.title"
          :url="menu.url"
          :collapsed="collapsed"
          :current="menu.title === props.current"
        >
          <template #icon>
            <HugeiconsIcon
              :color=" menu.title === props.current ? '#EDEDEE' : '#FFFFFF'"
              class="w-6"
              :fill="menu.title === props.current ? '#AF513A' : '#B7B8BB'"
              :icon="menu.icon"
            />
          </template>
          <template #default>
            {{ menu.title }}
          </template>
          <template
            v-if="menu.count"
            #count
          >
            {{ menu.count }}
          </template>
          <template
            v-if="menu.soon"
            #label
          >
            Soon
          </template>
        </SidebarMenu>
      </div>
      <div
        class="flex items-center justify-center h-16 mt-auto hover:bg-gray-300"
        :class="[collapsed ? 'px-2' : 'px-2 lg:px-4']"
      >
        <SidebarMenu
          v-for="menu in SIDEBAR_MENU_SECTION_C"
          :key="menu.title"
          :url="menu.url"
          :collapsed="collapsed"
          :current="menu.title === props.current"
        >
          <template #icon>
            <HugeiconsIcon
              :color=" menu.title === props.current ? '#EDEDEE' : '#FFFFFF'"
              class="w-6"
              :fill="menu.title === props.current ? '#AF513A' : '#B7B8BB'"
              :icon="menu.icon"
            />
          </template>
          <template #default>
            {{ menu.title }}
          </template>
          <template
            v-if="menu.count"
            #count
          >
            {{ menu.count }}
          </template>
          <template
            v-if="menu.soon"
            #label
          >
            Soon
          </template>
        </SidebarMenu>
      </div>
    </div>

    <nav
      class="fixed bottom-0 left-0 right-0 flex justify-around items-center h-16 bg-white border-t border-[#EDEDEE] md:hidden"
    >
      <SidebarMenu
        v-for="menu in [...SIDEBAR_MENU_SECTION_A, ...SIDEBAR_MENU_SECTION_B]"
        :key="menu.title"
        :icon="menu.icon"
        :current="menu.title === props.current"
      >
        <template #default>
          <span class="text-xs">{{ menu.title }}</span>
        </template>
        <template
          v-if="menu.count"
          #count
        >
          {{ menu.count }}
        </template>
      </SidebarMenu>
    </nav>
  </div>
</template>

<script setup lang="ts">
import {onMounted} from "vue";
import {useUserAPI} from "@/composables/api/user/user";
import {verifications} from "@/composables/helpers/format/verifications";
import {useRouter} from "nuxt/app";
import SidebarMenu from "@/components/sidebar/sidebar-menu.vue";

import {SIDEBAR_MENU_SECTION_A, SIDEBAR_MENU_SECTION_B, SIDEBAR_MENU_SECTION_C} from "@/constants/navigation/side-bar";
import {HugeiconsIcon} from "@hugeicons/vue";

const props = defineProps({
  current: {
    type: String,
    default: "",
  },
  collapsed: {
    type: Boolean,
    default: false,
  },
});

const router = useRouter()

onMounted(async () => {
  const fetchUserApiResponse = await useUserAPI().FetchLoggedInUser()
  if (fetchUserApiResponse.success) {
    // const userVerifications = verifications(fetchUserApiResponse.data?.data?.verifications);
    // if (userVerifications?.email !== "VERIFIED") {
    //   await router.push("/verification/email");
    // }
  } else {
    localStorage.clear()
    await router.push("/authentication/login");
  }
})
</script>

<style scoped>
</style>

<template>
  <div>
    <div
      class="hidden md:flex  flex-col items-left md:2 h-full overflow-hidden text-[#4B4D55] font-medium bg-[#FFFFFF] border-1 border-[#EDEDEE] rounded"
      :class="[collapsed ? 'w-16' : 'md:w-16 lg:w-64']"
    >
      <a
        class="flex mt-3 h-12 md:px-3 lg:px-4"
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
        class="flex flex-col items-center mt-1 border-t border-[#EDEDEE] pt-5 pb-8"
        :class="[collapsed ? 'px-2' : 'px-2 lg:px-4']"
      >
        <template
          v-for="menu in SIDEBAR_MENU_SECTION_A"
          :key="menu.title"
        >
          <!-- Item with children (expandable) -->
          <template v-if="menu.children">
            <div
              class="flex cursor-pointer md:justify-center lg:justify-start w-full rounded-lg px-3 py-2 mb-1 hover:bg-gray-300"
              :class="[isChildActive(menu.title) ? 'font-semibold' : 'font-normal']"
              @click="toggleExpanded(menu.title)"
            >
              <HugeiconsIcon
                :color="isChildActive(menu.title) ? '#EDEDEE' : '#FFFFFF'"
                class="w-6"
                :fill="isChildActive(menu.title) ? '#AF513A' : '#B7B8BB'"
                :icon="menu.icon"
              />
              <div
                v-if="!collapsed"
                class="ml-2 my-auto text-sm md:hidden lg:block mr-auto"
              >
                {{ menu.title }}
              </div>
            </div>
            <div
              v-if="isExpanded(menu.title) && !collapsed"
              class="w-full md:hidden lg:block"
            >
              <NuxtLink
                v-for="child in menu.children"
                :key="child.title"
                :to="{ name: child.url }"
                class="flex cursor-pointer w-full rounded-lg px-3 py-2 mb-1 hover:bg-gray-300"
                :class="[child.title === props.current ? 'bg-[#EDEDEE] text-[#1E212B] font-semibold' : 'font-normal text-[#4B4D55]']"
              >
                <span class="ml-8 text-sm">{{ child.title }}</span>
              </NuxtLink>
            </div>
          </template>

          <!-- Regular item -->
          <SidebarMenu
            v-else
            :collapsed="collapsed"
            :url="menu.url"
            :current="menu.title === props.current"
            :class="[menu.title === props.current ? 'bg-[#EDEDEE] text-[#1E212B] font-semibold' : 'font-normal']"
          >
            <template #default>
              {{ menu.title }}
            </template>
            <template #icon>
              <HugeiconsIcon
                :color="menu.title === props.current ? '#EDEDEE' : '#FFFFFF'"
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
        </template>
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

    <!-- Mobile Top Bar -->
    <div class="fixed top-0 left-0 right-0 z-40 flex items-center justify-between h-14 px-4 bg-white border-b border-[#EDEDEE] md:hidden">
      <a href="#">
        <img
          src="https://res.cloudinary.com/sdq121/image/upload/v1755371764/ltumpn6lhocowaxdi68p.png"
          alt="Gro App Logo"
          class="h-7"
        >
      </a>
    </div>

    <!-- Mobile Bottom Tab Bar -->
    <div class="fixed bottom-0 left-0 right-0 z-40 flex items-center bg-white border-t border-[#EDEDEE] md:hidden" style="padding-bottom: env(safe-area-inset-bottom)">
      <NuxtLink
        v-for="tab in MOBILE_BOTTOM_TABS"
        :key="tab.title"
        :to="{ name: tab.url }"
        class="flex flex-col items-center justify-center flex-1 py-2 gap-0.5"
        :class="[tab.title === props.current ? 'text-[#AF513A]' : 'text-[#6F7177]']"
      >
        <HugeiconsIcon
          :icon="tab.icon"
          :size="22"
          :color="tab.title === props.current ? '#AF513A' : '#6F7177'"
          :fill="tab.title === props.current ? '#AF513A' : '#6F7177'"
        />
        <span class="text-[10px] font-medium">{{ tab.shortTitle ?? tab.title }}</span>
      </NuxtLink>
      <button
        class="flex flex-col items-center justify-center flex-1 py-2 gap-0.5"
        :class="[mobileMoreOpen ? 'text-[#AF513A]' : 'text-[#6F7177]']"
        @click="mobileMoreOpen = true"
      >
        <svg viewBox="0 0 24 24" fill="currentColor" class="w-[22px] h-[22px]">
          <circle cx="5" cy="12" r="2" />
          <circle cx="12" cy="12" r="2" />
          <circle cx="19" cy="12" r="2" />
        </svg>
        <span class="text-[10px] font-medium">More</span>
      </button>
    </div>

    <!-- Mobile More Bottom Sheet -->
    <Transition name="sheet">
      <div
        v-if="mobileMoreOpen"
        class="fixed inset-0 z-50 md:hidden flex flex-col justify-end"
      >
        <!-- Backdrop -->
        <div
          class="absolute inset-0 bg-black/40"
          @click="mobileMoreOpen = false"
        />
        <!-- Sheet panel -->
        <div class="relative bg-white rounded-t-2xl max-h-[80vh] overflow-y-auto" style="padding-bottom: env(safe-area-inset-bottom)">
          <!-- Drag handle -->
          <div class="flex justify-center pt-3 pb-1">
            <div class="w-10 h-1 rounded-full bg-[#DBDBDD]" />
          </div>

          <div class="px-4 pt-3 pb-2">
            <!-- Section A overflow (items not in bottom tab bar) -->
            <p class="text-[11px] font-semibold text-[#6F7177] uppercase tracking-wider mb-3">
              Menu
            </p>
            <template
              v-for="menu in MOBILE_MORE_SECTION_A"
              :key="menu.title"
            >
              <template v-if="menu.children">
                <div
                  class="flex cursor-pointer w-full rounded-xl px-3 py-3 mb-1 hover:bg-[#F6F6F7]"
                  :class="[isChildActive(menu.title) ? 'bg-[#F6F6F7]' : '']"
                  @click="toggleExpanded(menu.title)"
                >
                  <div class="w-10 h-10 rounded-xl bg-[#F6F6F7] flex items-center justify-center mr-3 shrink-0">
                    <HugeiconsIcon
                      :icon="menu.icon"
                      :size="20"
                      :color="isChildActive(menu.title) ? '#AF513A' : '#4B4D55'"
                      :fill="isChildActive(menu.title) ? '#AF513A' : '#4B4D55'"
                    />
                  </div>
                  <div class="my-auto text-sm font-medium text-[#1E212B] mr-auto">
                    {{ menu.title }}
                  </div>
                  <svg
                    class="w-4 h-4 my-auto text-[#6F7177] transition-transform"
                    :class="[isExpanded(menu.title) ? 'rotate-90' : '']"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                  >
                    <path d="M9 18l6-6-6-6" />
                  </svg>
                </div>
                <div
                  v-if="isExpanded(menu.title)"
                  class="ml-4 mb-1"
                >
                  <NuxtLink
                    v-for="child in menu.children"
                    :key="child.title"
                    :to="{ name: child.url }"
                    class="flex cursor-pointer w-full rounded-xl px-3 py-2.5 mb-0.5 hover:bg-[#F6F6F7]"
                    :class="[child.title === props.current ? 'bg-[#F6F6F7] text-[#AF513A] font-semibold' : 'text-[#4B4D55]']"
                    @click="mobileMoreOpen = false"
                  >
                    <span class="text-sm ml-2">{{ child.title }}</span>
                  </NuxtLink>
                </div>
              </template>
              <NuxtLink
                v-else
                :to="{ name: menu.url }"
                class="flex cursor-pointer w-full rounded-xl px-3 py-3 mb-1 hover:bg-[#F6F6F7]"
                :class="[menu.title === props.current ? 'bg-[#F6F6F7]' : '']"
                @click="mobileMoreOpen = false"
              >
                <div class="w-10 h-10 rounded-xl bg-[#F6F6F7] flex items-center justify-center mr-3 shrink-0">
                  <HugeiconsIcon
                    :icon="menu.icon"
                    :size="20"
                    :color="menu.title === props.current ? '#AF513A' : '#4B4D55'"
                    :fill="menu.title === props.current ? '#AF513A' : '#4B4D55'"
                  />
                </div>
                <div class="my-auto text-sm font-medium text-[#1E212B] mr-auto">
                  {{ menu.title }}
                </div>
                <small
                  v-if="menu.soon"
                  class="ml-2 px-2 my-auto text-xs text-[#3C3F28] rounded-lg bg-[#EAEBE3]"
                >
                  Soon
                </small>
              </NuxtLink>
            </template>

            <!-- Sales Channels -->
            <p class="text-[11px] font-semibold text-[#6F7177] uppercase tracking-wider mt-4 mb-3">
              Sales channels
            </p>
            <NuxtLink
              v-for="menu in SIDEBAR_MENU_SECTION_B"
              :key="menu.title"
              :to="{ name: menu.url }"
              class="flex cursor-pointer w-full rounded-xl px-3 py-3 mb-1 hover:bg-[#F6F6F7]"
              :class="[menu.title === props.current ? 'bg-[#F6F6F7]' : '']"
              @click="mobileMoreOpen = false"
            >
              <div class="w-10 h-10 rounded-xl bg-[#F6F6F7] flex items-center justify-center mr-3 shrink-0">
                <HugeiconsIcon
                  :icon="menu.icon"
                  :size="20"
                  :color="menu.title === props.current ? '#AF513A' : '#4B4D55'"
                  :fill="menu.title === props.current ? '#AF513A' : '#4B4D55'"
                />
              </div>
              <div class="my-auto text-sm font-medium text-[#1E212B] mr-auto">
                {{ menu.title }}
              </div>
              <small
                v-if="menu.soon"
                class="ml-2 px-2 my-auto text-xs text-[#3C3F28] rounded-lg bg-[#EAEBE3]"
              >
                Soon
              </small>
            </NuxtLink>

            <!-- Settings -->
            <div class="border-t border-[#EDEDEE] mt-4 pt-3">
              <NuxtLink
                v-for="menu in SIDEBAR_MENU_SECTION_C"
                :key="menu.title"
                :to="{ name: menu.url }"
                class="flex cursor-pointer w-full rounded-xl px-3 py-3 mb-1 hover:bg-[#F6F6F7]"
                :class="[menu.title === props.current ? 'bg-[#F6F6F7]' : '']"
                @click="mobileMoreOpen = false"
              >
                <div class="w-10 h-10 rounded-xl bg-[#F6F6F7] flex items-center justify-center mr-3 shrink-0">
                  <HugeiconsIcon
                    :icon="menu.icon"
                    :size="20"
                    :color="menu.title === props.current ? '#AF513A' : '#4B4D55'"
                    :fill="menu.title === props.current ? '#AF513A' : '#4B4D55'"
                  />
                </div>
                <div class="my-auto text-sm font-medium text-[#1E212B]">
                  {{ menu.title }}
                </div>
              </NuxtLink>
            </div>
          </div>

          <!-- Close button -->
          <div class="px-4 pb-4">
            <button
              class="w-full py-3.5 rounded-2xl bg-[#F6F6F7] text-sm font-semibold text-[#1E212B] hover:bg-[#EDEDEE] transition-colors"
              @click="mobileMoreOpen = false"
            >
              Close
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import {onMounted, ref} from "vue";
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

const mobileMoreOpen = ref(false)
const expandedItems = ref<Set<string>>(new Set())

// First 4 items from Section A shown as bottom tabs
const MOBILE_BOTTOM_TABS = SIDEBAR_MENU_SECTION_A.filter(m => !m.children).slice(0, 4)
// Remaining Section A items (including children) shown in More sheet
const MOBILE_MORE_SECTION_A = SIDEBAR_MENU_SECTION_A.slice(4)

const isChildActive = (menuTitle: string) => {
  const item = SIDEBAR_MENU_SECTION_A.find(m => m.title === menuTitle)
  return item?.children?.some(child => child.title === props.current) ?? false
}

const toggleExpanded = (title: string) => {
  if (expandedItems.value.has(title)) {
    expandedItems.value.delete(title)
  } else {
    expandedItems.value.add(title)
  }
}

const isExpanded = (title: string) => {
  return expandedItems.value.has(title) || isChildActive(title)
}

onMounted(async () => {
  // Auto-expand parent items whose child is active
  for (const menu of SIDEBAR_MENU_SECTION_A) {
    if (menu.children && isChildActive(menu.title)) {
      expandedItems.value.add(menu.title)
    }
  }
  const fetchUserApiResponse = await useUserAPI().FetchLoggedInUser()
  if (fetchUserApiResponse.success) {
    const userVerifications = verifications(fetchUserApiResponse.data?.data?.verifications);
    if (userVerifications?.email !== "VERIFIED") {
      await router.push("/verification/email");
    }
  } else {
    localStorage.clear()
    await router.push("/authentication/login");
  }
})
</script>

<style scoped>
.sheet-enter-active,
.sheet-leave-active {
  transition: opacity 0.25s ease;
}
.sheet-enter-active .relative,
.sheet-leave-active .relative {
  transition: transform 0.3s cubic-bezier(0.32, 0.72, 0, 1);
}
.sheet-enter-from,
.sheet-leave-to {
  opacity: 0;
}
.sheet-enter-from .relative {
  transform: translateY(100%);
}
.sheet-leave-to .relative {
  transform: translateY(100%);
}
</style>

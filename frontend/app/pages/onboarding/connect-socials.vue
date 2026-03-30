<script lang="ts" setup>
definePageMeta({ middleware: ['auth'] })
import { useRouter } from 'nuxt/app'
import { HugeiconsIcon } from '@hugeicons/vue'
import { ArrowRight01Icon } from '@hugeicons/core-free-icons'
import GroBasicButton from '@/components/buttons/GroBasicButton.vue'
import { useAuthentication } from '@/composables/api/authentication/authentication'

const router = useRouter()

function finish() {
  router.push('/dashboard')
}

async function logout() {
  await useAuthentication().Logout()
  router.push('/authentication/login')
}
</script>

<template>
  <div
    class="min-h-screen w-full"
    style="background-color: #FAFAFA;"
  >
    <!-- Header -->
    <div class="flex items-center justify-between px-8 pt-7">
      <NuxtLink to="/">
        <img
          src="https://res.cloudinary.com/sdq121/image/upload/v1755371764/ltumpn6lhocowaxdi68p.png"
          alt="Gro"
          class="h-7"
        >
      </NuxtLink>
      <button
        class="text-sm text-[#1E212B] font-medium hover:opacity-70 transition-opacity"
        @click="logout"
      >
        Log out
      </button>
    </div>

    <!-- Content -->
    <div class="flex flex-col items-center mt-20 md:mt-[7%] px-4">
      <!-- Progress indicator (all 4 steps complete) -->
      <div class="flex flex-col items-center mb-8">
        <p class="font-dm-sans text-sm text-[#6F7177] mb-3 flex items-center gap-1.5">
          Step 4 of 4
          <HugeiconsIcon
            :icon="ArrowRight01Icon"
            :size="14"
            color="#6F7177"
          />
          <span class="font-medium text-[#1E212B]">Connect Socials</span>
        </p>
        <div class="flex items-center">
          <template
            v-for="dot in 4"
            :key="dot"
          >
            <div class="w-6 h-6 rounded-full flex items-center justify-center border-2 border-[#1E212B] bg-white">
              <div class="h-4 w-4 rounded-full bg-[#1E212B]" />
            </div>
            <div
              v-if="dot < 4"
              class="w-16 h-px bg-[#1E212B]"
            />
          </template>
        </div>
      </div>

      <!-- Connect Socials card -->
      <div class="bg-white border border-[#EDEDEE] rounded-2xl py-4 px-4 md:py-10 md:px-8 w-full max-w-lg">
        <h2 class="text-xl md:text-2xl font-bold text-[#1E212B] mb-1">
          Connect your social networks
        </h2>
        <p class="text-sm text-[#6F7177] mb-7">
          Connect your social networks to automatically track customer info and interests as they arise
        </p>

        <div class="flex flex-col gap-3 mb-8">
          <div class="flex items-center justify-between border border-[#EDEDEE] rounded-xl px-4 py-3.5">
            <div class="flex items-center gap-3">
              <img
                src="https://upload.wikimedia.org/wikipedia/commons/6/6b/WhatsApp.svg"
                alt="WhatsApp"
                class="w-8 h-8"
              >
              <span class="text-sm font-medium text-black">Whatsapp</span>
            </div>
            <div>
              <GroBasicButton
                color="tertiary"
                size="xs"
                shape="custom"
                class="px-2 py-2.5"
              >
                Connect
              </GroBasicButton>
            </div>
          </div>
          <div class="flex items-center justify-between border border-[#EDEDEE] rounded-xl px-4 py-3.5">
            <div class="flex items-center gap-3">
              <img
                src="https://upload.wikimedia.org/wikipedia/commons/a/a5/Instagram_icon.png"
                alt="Instagram"
                class="w-8 h-8 rounded-lg"
              >
              <span class="text-sm font-medium text-black">Instagram</span>
            </div>
            <div>
              <GroBasicButton
                color="tertiary"
                size="xs"
                shape="custom"
                class="px-2 py-2.5"
              >
                Connect
              </GroBasicButton>
            </div>
          </div>
          <div class="flex items-center justify-between border border-[#EDEDEE] rounded-xl px-4 py-3.5">
            <div class="flex items-center gap-3">
              <img
                src="https://upload.wikimedia.org/wikipedia/en/a/a9/TikTok_logo.svg"
                alt="TikTok"
                class="w-8 h-8"
              >
              <span class="text-sm font-medium text-black">Tiktok</span>
            </div>
            <div>
              <GroBasicButton
                color="tertiary"
                size="xs"
                shape="custom"
                class="px-2 py-2.5"
              >
                Connect
              </GroBasicButton>
            </div>
          </div>
        </div>

        <div class="flex justify-end">
          <GroBasicButton
            color="primary"
            size="md"
            shape="custom"
            class="px-6 py-2.5"
            @click="finish"
          >
            Get started
            <HugeiconsIcon
              :icon="ArrowRight01Icon"
              :size="16"
              color="white"
            />
          </GroBasicButton>
        </div>
      </div>

      <p
        class="mt-5 text-sm text-[#6F7177] cursor-pointer hover:text-[#1E212B] transition-colors"
        @click="finish"
      >
        I'll connect my social networks later
      </p>
    </div>
  </div>
</template>

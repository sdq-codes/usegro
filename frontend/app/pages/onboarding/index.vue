<script lang="ts" setup>
definePageMeta({ middleware: ['auth'] })
import { ref, computed } from 'vue'
import { useRouter } from 'nuxt/app'
import { HugeiconsIcon } from '@hugeicons/vue'
import {
  ArrowLeft01Icon,
  ArrowRight01Icon,
  CheckmarkCircle02Icon,
  InformationCircleIcon,
  MessageMultiple01Icon,
  GlobalIcon,
  Tag01Icon,
  Store04Icon, Trolley01Icon, Folder02Icon, Calendar02Icon,
} from '@hugeicons/core-free-icons'
import GroBasicInput from '@/components/forms/input/GroBasicInput.vue'
import GroBasicButton from '@/components/buttons/GroBasicButton.vue'
import { useCRMAPI } from '@/composables/api/crm/crm'
import { useAuthentication } from '@/composables/api/authentication/authentication'

const router = useRouter()

const screen = ref(1)
const loading = ref(false)
const apiError = ref('')
const step1Error = ref('')
const checkingName = ref(false)

// Step 1
const fullName = ref('')
const businessName = ref('')

// Step 2 — business_info: 'starter_package' | 'existing_package'
const businessInfo = ref('')

// Step 3 — sales_channel_type values from backend enum
const salesChannels = ref<string[]>([])

// Step 4 — product_type values from backend enum
const productTypes = ref<string[]>([])

const progressDot = computed(() => {
  if (screen.value === 1) return 1
  if (screen.value === 2) return 2
  if (screen.value <= 4) return 3
  return 4
})

const stepLabel = computed(() => ({
  1: 'Personal info',
  2: 'Business info',
  3: 'Sales channels',
  4: 'Sales channels',
}[screen.value]))

const stepNumber = computed(() => {
  if (screen.value === 1) return 1
  if (screen.value === 2) return 2
  if (screen.value <= 4) return 3
  return 4
})

const isDotComplete = (dot: number) => dot < progressDot.value
const isDotActive = (dot: number) => dot === progressDot.value

function toggleChannel(value: string) {
  const idx = salesChannels.value.indexOf(value)
  if (idx === -1) salesChannels.value.push(value)
  else salesChannels.value.splice(idx, 1)
}

function toggleProductType(value: string) {
  const idx = productTypes.value.indexOf(value)
  if (idx === -1) productTypes.value.push(value)
  else productTypes.value.splice(idx, 1)
}

// Submit all steps at once after step 4
async function submitAll() {
  loading.value = true
  apiError.value = ''

  // 1. Create CRM with all collected data
  const createRes = await useCRMAPI().CreateCRM({
    full_name: fullName.value.trim(),
    business_name: businessName.value.trim(),
    business_info: businessInfo.value,
  })
  if (!createRes.success) {
    apiError.value = createRes.error
    loading.value = false
    return
  }

  const crmId = createRes.data?.data?.id
  if (crmId) localStorage.setItem('crm-id', crmId)
  if (!crmId) {
    apiError.value = 'Failed to retrieve CRM ID'
    loading.value = false
    return
  }

  // 2. Save sales channels
  if (salesChannels.value.length > 0) {
    const channelsRes = await useCRMAPI().CreateSalesChannels(crmId, salesChannels.value)
    if (!channelsRes.success) {
      apiError.value = channelsRes.error
      loading.value = false
      return
    }
  }

  // 3. Save product types
  if (productTypes.value.length > 0) {
    const productsRes = await useCRMAPI().CreateStockProducts(crmId, productTypes.value)
    if (!productsRes.success) {
      apiError.value = productsRes.error
      loading.value = false
      return
    }
  }

  loading.value = false
  router.push('/onboarding/connect-socials')
}

async function checkBusinessNameAndNext() {
  step1Error.value = ''
  checkingName.value = true
  const res = await useCRMAPI().CheckBusinessNameExists(businessName.value.trim())
  checkingName.value = false
  if (res.error) { step1Error.value = res.error; return }
  if (res.exists) { step1Error.value = 'This business name is already taken'; return }
  screen.value = 2
}

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
      <!-- Progress indicator -->
      <div class="flex flex-col items-center mb-8">
        <p class="font-dm-sans text-sm text-[#6F7177] mb-3 flex items-center gap-1.5">
          Step {{ stepNumber }} of 4
          <HugeiconsIcon
            :icon="ArrowRight01Icon"
            :size="14"
            color="#6F7177"
          />
          <span class="font-medium text-[#1E212B]">{{ stepLabel }}</span>
        </p>
        <div class="flex items-center">
          <template
            v-for="dot in 4"
            :key="dot"
          >
            <div
              class="w-6 h-6 rounded-full flex items-center justify-center transition-all"
              :class="{
                'bg-white border-2 border-[#EDEDEE]': isDotActive(dot),
                'bg-white border-2 border-[#1E212B]': isDotComplete(dot),
                'bg-white border-2 border-[#DBDBDD]': !isDotActive(dot) && !isDotComplete(dot),
              }"
            >
              <div
                v-if="isDotComplete(dot)"
                :class="{
                  'bg-[#1E212B]': isDotActive(dot),
                }"
              >
                <HugeiconsIcon

                  :icon="CheckmarkCircle02Icon"
                  color="#1E212B"
                  class=" h-4 w-4 rounded-full transition-all"
                />
              </div>
              <div
                v-else
                class=" h-4 w-4 rounded-full transition-all"
                :class="{
                  'bg-[#1E212B]': isDotActive(dot),
                }"
              >
              </div>
            </div>
            <div
              v-if="dot < 4"
              class="w-16 h-px"
              :class="isDotComplete(dot) ? 'bg-[#1E212B]' : 'bg-[#DBDBDD]'"
            />
          </template>
        </div>
      </div>

      <!-- ── SCREEN 1: Personal info ── -->
      <div
        v-if="screen === 1"
        class="bg-white border border-[#EDEDEE] rounded-2xl py-10 px-8 w-full max-w-2xl"
      >
        <h2 class="text-2xl font-semibold text-black mb-1">
          Welcome to Gro.com!
        </h2>
        <p class="text-sm font-medium text-[#6F7177] mb-7">
          We just need a little more info to get started.<br>You'll be able to edit this later
        </p>

        <div class="mb-5">
          <div class="flex items-center gap-1.5 mb-1.5">
            <span class="text-sm font-medium text-[#1E212B]">Full Business Name</span>
            <HugeiconsIcon
              :icon="InformationCircleIcon"
              :size="15"
              color="#6F7177"
            />
          </div>
          <GroBasicInput
            v-model="fullName"
            placeholder="Text placeholder"
            color="primary"
          />
        </div>

        <div class="mb-6">
          <p class="text-sm font-medium text-[#1E212B] mb-1.5">
            Business identifier
          </p>
          <div class="flex rounded-lg border border-[#EDEDEE] bg-[#F6F6F7] overflow-hidden focus-within:border-[#1E212B] transition-colors">
            <span class="px-3 py-2.5 text-sm text-[#6F7177] border-r border-[#EDEDEE] bg-white shrink-0">
              Usegro.com/
            </span>
            <input
              v-model="businessName"
              placeholder="examplellc"
              class="flex-1 px-3 py-2.5 text-sm bg-[#F6F6F7] outline-none text-[#1E212B] placeholder-[#939499] focus:bg-white transition-colors"
              @keydown="(e) => e.key === ' ' && e.preventDefault()"
              @input="(e) => { const el = e.target as HTMLInputElement; el.value = el.value.replace(/ /g, ''); businessName = el.value }"
            >
          </div>
          <p class="text-xs text-[#939499] mt-1.5">
            Your business name may be displayed to customers
          </p>
        </div>

        <div
          v-if="step1Error"
          class="text-sm text-[#AF513A] mb-4"
        >
          {{ step1Error }}
        </div>


        <div class="flex justify-end">
          <div>
            <GroBasicButton
              color="primary"
              size="md"
              shape="custom"
              :loading="checkingName"
              :disabled="!fullName.trim() || !businessName.trim() || checkingName"
              class="px-6 py-2.5"
              @click="checkBusinessNameAndNext"
            >
              Next
            </GroBasicButton>
          </div>
        </div>
      </div>

      <!-- ── SCREEN 2: Business info ── -->
      <div
        v-else-if="screen === 2"
        class="bg-white border border-[#EDEDEE] rounded-2xl py-10 px-8 w-full max-w-3xl"
      >
        <h2 class="text-2xl font-bold text-black mb-1">
          Which of these best describes you.
        </h2>
        <p class="text-sm text-[#6F7177] mb-7">
          We'll help you get set up based on your business needs.
        </p>

        <div class="grid grid-cols-2 gap-3 mb-8">
          <label
            v-for="opt in [
              { value: 'starter_package', label: 'I\'m just starting out my business' },
              { value: 'existing_package', label: 'I\'m already selling online or in person' },
            ]"
            :key="opt.value"
            class="col-span-2 md:col-span-1 flex items-center gap-3 border rounded-xl px-4 py-4 cursor-pointer transition-colors"
            :class="businessInfo === opt.value ? 'border-[#1E212B] bg-[#F6F6F7]' : 'border-[#EDEDEE] hover:border-[#B7B8BB]'"
          >
            <div
              class="w-5 h-5 rounded-full border-2 flex items-center justify-center shrink-0"
              :class="businessInfo === opt.value ? 'border-[#1E212B]' : 'border-[#DBDBDD]'"
            >
              <div
                v-if="businessInfo === opt.value"
                class="w-2.5 h-2.5 rounded-full bg-[#1E212B]"
              />
            </div>
            <input
              v-model="businessInfo"
              type="radio"
              :value="opt.value"
              class="hidden"
            >
            <span class="text-sm text-[#1E212B]">{{ opt.label }}</span>
          </label>
        </div>

        <div class="flex items-center justify-between">
          <div>
            <button
              class="flex cursor-pointer items-center gap-1.5 text-sm text-[#6F7177] hover:text-[#1E212B] transition-colors"
              @click="screen--"
            >
              <HugeiconsIcon
                :icon="ArrowLeft01Icon"
                :size="16"
              />
              Back
            </button>
          </div>
          <div>
            <GroBasicButton
              color="primary"
              size="md"
              shape="custom"
              :disabled="!businessInfo"
              class="px-6 py-2.5 w-1/2"
              @click="screen = 3"
            >
              Next
              <HugeiconsIcon
                :icon="ArrowRight01Icon"
                :size="16"
                color="white"
              />
            </GroBasicButton>
          </div>
        </div>
      </div>

      <!-- ── SCREEN 3: Sales channels — where to sell ── -->
      <div
        v-else-if="screen === 3"
        class="bg-white border border-[#EDEDEE] rounded-2xl py-10 px-8 w-full max-w-3xl"
      >
        <h2 class="text-2xl font-bold text-black mb-1">
          Where would you like to sell?
        </h2>
        <p class="text-sm text-[#6F7177] mb-7">
          Pick as many as you like – You can always change these later. We'll make sure you're set up to sell in these places
        </p>

        <div class="grid grid-cols-2 gap-3 mb-3">
          <label
            v-for="opt in [
              { value: 'online_store', label: 'An online store', description: 'Create a fully customizable website', icon: Store04Icon },
              { value: 'social_media_store', label: 'Social media', description: 'Reach customers on Facebook, Instagram, tiktok and more', icon: MessageMultiple01Icon },
              { value: 'physical_store', label: 'In person', description: 'Sell at retail shops, pop-ups, or other physical locations', icon: Trolley01Icon },
              { value: 'existing_website_store', label: 'An existing website', description: 'Add a buy/subscribe button to your website', icon: GlobalIcon },
            ]"
            :key="opt.value"
            class="col-span-2 md:col-span-1 flex items-start gap-3 rounded-xl px-4 py-4 cursor-pointer transition-colors"
            :class="salesChannels.includes(opt.value) ? 'border-[#1E212B] bg-[#F6F6F7]' : 'border border-[#EDEDEE] hover:border-[#B7B8BB]'"
            @click="toggleChannel(opt.value)"
          >
            <HugeiconsIcon
              :icon="opt.icon"
              :size="30"
              color="#1E212B"
              class="my-auto shrink-0"
            />
            <div class="flex-1 my-auto min-w-0">
              <p class="text-sm font-semibold text-[#1E212B]">{{ opt.label }}</p>
              <p class="text-xs text-[#6F7177] mt-0.5 leading-relaxed">{{ opt.description }}</p>
            </div>
            <div
              class="w-5 h-5 my-auto rounded border-2 flex items-center justify-center shrink-0 transition-colors"
              :class="salesChannels.includes(opt.value) ? 'border-[#1E212B] bg-[#1E212B]' : 'border-[#DBDBDD]'"
            >
              <HugeiconsIcon
                v-if="salesChannels.includes(opt.value)"
                :icon="CheckmarkCircle02Icon"
                :size="12"
                color="white"
              />
            </div>
          </label>
        </div>

        <label
          class="flex items-center justify-between border rounded-xl px-4 py-3.5 cursor-pointer transition-colors mb-7"
          :class="salesChannels.includes('unknown_store') ? 'border-[#1E212B] bg-[#F6F6F7]' : 'border-[#EDEDEE] hover:border-[#B7B8BB]'"
          @click="toggleChannel('unknown_store')"
        >
          <span class="text-sm text-[#1E212B]">I'm not sure</span>
          <div
            class="w-5 h-5 rounded border-2 flex items-center justify-center transition-colors"
            :class="salesChannels.includes('unknown_store') ? 'border-[#1E212B] bg-[#1E212B]' : 'border-[#DBDBDD]'"
          >
            <HugeiconsIcon
              v-if="salesChannels.includes('unknown_store')"
              :icon="CheckmarkCircle02Icon"
              :size="12"
              color="white"
            />
          </div>
        </label>

        <div class="flex items-center justify-between">
          <button
            class="flex cursor-pointer items-center gap-1.5 text-sm text-[#6F7177] hover:text-[#1E212B] transition-colors"
            @click="screen--"
          >
            <HugeiconsIcon
              :icon="ArrowLeft01Icon"
              :size="16"
            />
            Back
          </button>
          <div>
            <GroBasicButton
              color="primary"
              size="md"
              shape="custom"
              class="px-6 py-2.5"
              @click="screen = 4"
            >
              Next
              <HugeiconsIcon
                :icon="ArrowRight01Icon"
                :size="16"
                color="white"
              />
            </GroBasicButton>
          </div>
        </div>
      </div>

      <!-- ── SCREEN 4: Stock products — what to sell ── -->
      <div
        v-else-if="screen === 4"
        class="bg-white border border-[#EDEDEE] rounded-2xl py-10 px-8 w-full max-w-3xl"
      >
        <h2 class="text-2xl font-bold text-black mb-1">
          What do you plan to sell first?
        </h2>
        <p class="text-sm text-[#6F7177] mb-7">
          Pick what you want to start with. We'll help you stock your store.
        </p>

        <div class="grid grid-cols-2 gap-3 mb-8">
          <label
            v-for="opt in [
              { value: 'physical_products', label: 'Products I buy or make myself', description: 'Shipped by me', icon: Tag01Icon },
              { value: 'digital_products', label: 'Digital products', description: 'Music, digital art or NFTs', icon: Folder02Icon },
              { value: 'services_products', label: 'Services', description: 'Coaching, housekeeping, consulting', icon: Calendar02Icon },
              { value: 'unknown_products', label: 'I\'ll decide later', description: '', icon: null },
            ]"
            :key="opt.value"
            class="col-span-2 md:col-span-1 flex items-start gap-3 rounded-xl px-4 py-4 cursor-pointer transition-colors"
            :class="productTypes.includes(opt.value) ? 'border-[#1E212B] bg-[#F6F6F7]' : 'border border-[#EDEDEE] hover:border-[#B7B8BB]'"
            @click="toggleProductType(opt.value)"
          >
            <HugeiconsIcon
              v-if="opt.icon"
              :icon="opt.icon"
              :size="30"
              color="#1E212B"
              class="my-auto shrink-0"
            />
            <div class="flex-1 my-auto min-w-0">
              <p class="text-sm font-semibold text-[#1E212B]">{{ opt.label }}</p>
              <p
                v-if="opt.description"
                class="text-xs text-[#6F7177] mt-0.5"
              >{{ opt.description }}</p>
            </div>
            <div
              class="w-5 h-5 rounded border-2 flex items-center justify-center shrink-0 my-auto transition-colors"
              :class="productTypes.includes(opt.value) ? 'border-[#1E212B] bg-[#1E212B]' : 'border-[#DBDBDD]'"
            >
              <HugeiconsIcon
                v-if="productTypes.includes(opt.value)"
                :icon="CheckmarkCircle02Icon"
                :size="12"
                color="white"
              />
            </div>
          </label>
        </div>

        <div
          v-if="apiError"
          class="text-sm text-[#AF513A] mb-4"
        >
          {{ apiError }}
        </div>

        <div class="flex items-center justify-between">
          <button
            class="flex cursor-pointer items-center gap-1.5 text-sm text-[#6F7177] hover:text-[#1E212B] transition-colors"
            @click="screen--"
          >
            <HugeiconsIcon
              :icon="ArrowLeft01Icon"
              :size="16"
            />
            Back
          </button>
          <div>
            <GroBasicButton
              color="primary"
              size="md"
              shape="custom"
              :loading="loading"
              :disabled="loading"
              class="px-6 py-2.5"
              @click="submitAll"
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
      </div>

    </div>
  </div>
</template>

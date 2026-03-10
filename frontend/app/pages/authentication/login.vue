<!-- eslint-disable vue/multi-word-component-names -->
<script lang="ts" setup>
import BasicPage from "@/layouts/Authentication/BasicPage.vue";
import GroBasicInput from "@/components/forms/input/GroBasicInput.vue";
import GroBasicButton from "@/components/buttons/GroBasicButton.vue";
import GroBasicPassword from "@/components/forms/input/GroBasicPassword.vue"
import GroBasicPin from "@/components/forms/input/GroBasicPin.vue";
import { HugeiconsIcon } from "@hugeicons/vue";
import { CancelSquareIcon, Mail01Icon } from "@hugeicons/core-free-icons";
import { ref } from "vue";
import { z } from "zod";
import { useAuthentication } from "@/composables/api/authentication/authentication";
import { setAccessToken } from "@/composables/api/useApi";
import { useRouter } from "nuxt/app";
import { notify } from "@/composables/helpers/notification/notification";

const router = useRouter()
const { CheckEmailExists, LoginUser, RequestEmailCode, VerifyEmailCode } = useAuthentication()

// --- state ---
type Step = 'email' | 'password' | 'email-code'
const step = ref<Step>('email')

const email = ref('')
const password = ref('')
const emailError = ref('')
const passwordError = ref('')
const emailCode = ref('')
const codeError = ref('')
const apiErrors = ref<string[]>([])
const loading = ref(false)

definePageMeta({ middleware: ['auth'] })

// --- email step ---
const continueWithEmail = async () => {
  apiErrors.value = []
  emailError.value = ''

  const parsed = z.string().email('Enter a valid email address').safeParse(email.value)
  if (!parsed.success) {
    emailError.value = parsed.error.issues[0].message
    return
  }

  loading.value = true
  const result = await CheckEmailExists(email.value)
  loading.value = false

  if (!result.exists && !result.errors?.length) {
    emailError.value = 'No account found with this email address'
    return
  }

  step.value = 'password'
}

// --- password step ---
const submitLogin = async () => {
  apiErrors.value = []
  passwordError.value = ''

  if (!password.value) {
    passwordError.value = 'Password is required'
    return
  }

  loading.value = true
  const result = await LoginUser({ email: email.value, password: password.value })
  loading.value = false

  if (!result.success) {
    apiErrors.value = [result.error]
    notify(result.error, 'error')
    return
  }

  setAccessToken(result.data?.data?.access_token)
  localStorage.setItem('refresh_token', result.data?.data?.refresh_token)
  notify(result.data?.response_message, 'success')
  router.push('/dashboard')
}

// --- email code step ---
const requestEmailCode = async () => {
  apiErrors.value = []
  loading.value = true
  const result = await RequestEmailCode(email.value)
  loading.value = false

  if (!result.success) {
    apiErrors.value = [result.error]
    notify(result.error, 'error')
    return
  }

  notify('A 6-digit code has been sent to your email.', 'success')
  step.value = 'email-code'
}

const submitEmailCode = async () => {
  apiErrors.value = []
  codeError.value = ''

  if (emailCode.value.length < 6) {
    codeError.value = 'Enter the 6-character code sent to your email'
    return
  }

  loading.value = true
  const result = await VerifyEmailCode(email.value, emailCode.value)
  loading.value = false

  if (!result.success) {
    codeError.value = result.error
    notify(result.error, 'error')
    return
  }

  setAccessToken(result.data?.data?.access_token)
  localStorage.setItem('refresh_token', result.data?.data?.refresh_token)
  notify(result.data?.response_message || 'Signed in successfully', 'success')
  router.push('/dashboard')
}

// --- google popup ---
let popupRef: Window | null = null

const loginWithGoogle = () => {
  const width = 500
  const height = 600
  const left = window.screenX + (window.outerWidth - width) / 2
  const top = window.screenY + (window.outerHeight - height) / 2

  popupRef = window.open(
    'http://localhost:8090/api/v1/authentication/google/login',
    'Google Login',
    `width=${width},height=${height},left=${left},top=${top},resizable=yes,scrollbars=yes`
  )
}

if (process.client) {
  window.addEventListener('message', (event: MessageEvent) => {
    if (event.origin !== window.location.origin) return

    if (event.data?.type === 'GOOGLE_AUTH_SUCCESS') {
      setAccessToken(event.data.accessToken)
      localStorage.setItem('refresh_token', event.data.refreshToken)
      popupRef = null
      router.push('/dashboard')
    }

    if (event.data?.type === 'GOOGLE_AUTH_ERROR') {
      apiErrors.value = ['Google login failed. Please try again.']
      notify('Google login failed. Please try again.', 'error')
      popupRef = null
    }
  })
}
</script>

<template>
  <section>
    <BasicPage>
      <img
        src="https://res.cloudinary.com/sdq121/image/upload/v1755378345/wnrloghcmniaatz3d9ie.png"
        alt="7 orange"
        class="mx-auto"
      >
      <h2 class="text-center text-3xl font-bold mt-6">
        Sign in to Gro
      </h2>

      <div class="mt-8 md:mt-16 bg-white border border-[#EDEDEE] py-10 px-7 rounded-xl">

        <!-- ── Step 1: Email ── -->
        <template v-if="step === 'email'">
          <GroBasicInput
            v-model="email"
            placeholder="Enter your email"
            :color="emailError ? 'error' : 'primary'"
            :hint="emailError"
            @keyup.enter="continueWithEmail"
          >
            <h6>Email</h6>
          </GroBasicInput>

          <div
            v-if="apiErrors.length"
            class="my-2 space-y-1"
          >
            <div
              v-for="(msg, i) in apiErrors"
              :key="i"
              class="flex gap-1 text-sm text-[#AF513A]"
            >
              <HugeiconsIcon
                color="#FFFFFF"
                fill="#AF513A"
                :icon="CancelSquareIcon"
              />
              <span class="my-auto">{{ msg }}</span>
            </div>
          </div>

          <GroBasicButton
            color="primary"
            size="md"
            shape="custom"
            class="w-full mt-4 py-3"
            :disabled="loading"
            @click="continueWithEmail"
          >
            {{ loading ? 'Checking...' : 'Continue' }}
          </GroBasicButton>

          <div class="inline-flex mt-5 w-full justify-center">
            <img
              class="h-[2px] w-full my-auto"
              src="https://res.cloudinary.com/sdq121/image/upload/v1755396695/vbqph8myxctpwezdseja.png"
              alt="single line"
            >
            <h6 class="mx-2 text-[#4B4D55]">
              Or
            </h6>
            <img
              class="h-[2px] w-full my-auto"
              src="https://res.cloudinary.com/sdq121/image/upload/v1755396695/vbqph8myxctpwezdseja.png"
              alt="single line"
            >
          </div>

          <GroBasicButton
            color="tertiary"
            size="md"
            shape="custom"
            class="w-full mt-5 py-3"
            @click="loginWithGoogle"
          >
            <template #frontIcon>
              <img
                src="https://res.cloudinary.com/sdq121/image/upload/v1755425958/xdefngg3jpctaghb7fnh.svg"
                alt="Google Login Icon"
              >
            </template>
            <template #default>
              Continue with Google
            </template>
          </GroBasicButton>

          <h6 class="w-full text-center text-sm mt-5">
            Don't have an account?
            <span class="text-[#D16B07] cursor-pointer">
              <NuxtLink :to="{ name: 'authentication-register' }">Sign Up</NuxtLink>
            </span>
          </h6>
        </template>

        <!-- ── Step 2: Password ── -->
        <template v-else-if="step === 'password'">
          <!-- Static email display -->
          <p class="text-sm text-[#1E212B] mb-2">
            Email
          </p>
          <p class="text-md text-[#6F7177] mb-4">
            {{ email }}
          </p>

          <GroBasicPassword
            v-model="password"
            placeholder="Enter your password"
            :color="passwordError ? 'error' : 'primary'"
            :hint="passwordError"
            @keyup.enter="submitLogin"
          >
            <h6>Password</h6>
          </GroBasicPassword>

          <div class="flex justify-end mt-2">
            <NuxtLink
              :to="{ name: 'authentication-forgot-password' }"
              class="text-sm text-[#2176AE] cursor-pointer"
            >
              Forgot Password
            </NuxtLink>
          </div>

          <GroBasicButton
            color="primary"
            size="md"
            shape="custom"
            class="w-full mt-4 py-3"
            :disabled="loading"
            @click="submitLogin"
          >
            {{ loading ? 'Signing in...' : 'Continue' }}
          </GroBasicButton>

          <div class="inline-flex mt-5 w-full justify-center">
            <img
              class="h-[2px] w-full my-auto"
              src="https://res.cloudinary.com/sdq121/image/upload/v1755396695/vbqph8myxctpwezdseja.png"
              alt="single line"
            >
            <h6 class="mx-2 text-[#4B4D55]">
              Or
            </h6>
            <img
              class="h-[2px] w-full my-auto"
              src="https://res.cloudinary.com/sdq121/image/upload/v1755396695/vbqph8myxctpwezdseja.png"
              alt="single line"
            >
          </div>

          <GroBasicButton
            color="tertiary"
            size="md"
            shape="custom"
            class="w-full mt-5 py-3"
            :disabled="loading"
            @click="requestEmailCode"
          >
            <template #frontIcon>
              <HugeiconsIcon :icon="Mail01Icon" />
            </template>
            <template #default>
              {{ loading ? 'Sending...' : 'Continue with email Code' }}
            </template>
          </GroBasicButton>
        </template>

        <!-- ── Step 3: Email Code ── -->
        <template v-else-if="step === 'email-code'">
          <p class="text-sm text-[#4B4D55] mb-1">
            Enter the 6-character code sent to
          </p>
          <p class="text-sm font-medium text-[#1E212B] mb-6">
            {{ email }}
          </p>

          <GroBasicPin
            v-model="emailCode"
            :length="6"
            :color="codeError ? 'error' : 'primary'"
            :hint="codeError"
          >
            <h6>Verification Code</h6>
          </GroBasicPin>

          <div
            v-if="apiErrors.length"
            class="my-2 space-y-1"
          >
            <div
              v-for="(msg, i) in apiErrors"
              :key="i"
              class="flex gap-1 text-sm text-[#AF513A]"
            >
              <HugeiconsIcon
                color="#FFFFFF"
                fill="#AF513A"
                :icon="CancelSquareIcon"
              />
              <span class="my-auto">{{ msg }}</span>
            </div>
          </div>

          <GroBasicButton
            color="primary"
            size="md"
            shape="custom"
            class="w-full mt-6 py-3"
            :disabled="loading"
            @click="submitEmailCode"
          >
            {{ loading ? 'Verifying...' : 'Verify & Sign in' }}
          </GroBasicButton>

          <h6 class="w-full text-center text-sm mt-5">
            Didn't receive a code?
            <span
              class="text-[#D16B07] cursor-pointer"
              @click="requestEmailCode"
            >
              Resend
            </span>
          </h6>
        </template>
      </div>
    </BasicPage>
  </section>
</template>

<style scoped></style>

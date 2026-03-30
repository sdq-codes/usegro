<!-- eslint-disable vue/multi-word-component-names -->
<script lang="ts" setup>
import RegistrationPage from "@/layouts/Authentication/RegistrationPage.vue";
import GroBasicPassword from "@/components/forms/input/GroBasicPassword.vue";
import GroBasicInput from "@/components/forms/input/GroBasicInput.vue";
import GroBasicButton from "@/components/buttons/GroBasicButton.vue";
import {ref} from "vue";
import {RegisterSchema} from "@/composables/helpers/validation/authentication/authentication";
import {FormatZodErrors} from "@/composables/helpers/validation/util";
import {useAuthentication} from "@/composables/api/authentication/authentication";
import {setAccessToken} from "@/composables/api/useApi";
import {useRouter} from "nuxt/app";
import {notify} from "@/composables/helpers/notification/notification";

const router = useRouter()

const email = ref<string>("");
const password = ref<string>("");
const validationErrors = ref<Record<string, string[]>>({});
const apiError = ref<string>("");

definePageMeta({
  middleware: ["auth"],
});

const validateLoginForm = () => {
  validationErrors.value = {};
  const validationResult = RegisterSchema.safeParse({
    email: email.value,
    password: password.value,
  })
  validationErrors.value = validationResult?.success ?  {} : FormatZodErrors(validationResult.error!.issues)
}

const submitLoginForm = async () => {
  const loginApi = await useAuthentication().RegisterUser({
    email: email.value,
    password: password.value,
  })
  if (!loginApi?.success) {
    notify(loginApi?.error, 'error');
    apiError.value = loginApi?.error
  } else {
    notify(loginApi.data?.response_message, 'success');
    setAccessToken(loginApi.data?.data?.access_token)
    localStorage.setItem("session", "1")
    router.push("/verification/email")
  }
}

const registerUser = async () => {
  apiError.value = "";
  validateLoginForm();
  if (Object.keys(validationErrors.value).length === 0) {
    await submitLoginForm()
  }
}

let popupRef: Window | null = null

const loginWithGoogle = () => {
  const width = 500
  const height = 600
  const left = window.screenX + (window.outerWidth - width) / 2
  const top = window.screenY + (window.outerHeight - height) / 2

  const base = (useRuntimeConfig().public.apiBase as string) || 'http://localhost/api/v1'
  const googleLoginUrl = `${base.replace(/\/$/, '')}/base/authentication/google/login`


  popupRef = window.open(
    googleLoginUrl,
    'Google Login',
    `width=${width},height=${height},left=${left},top=${top},resizable=yes,scrollbars=yes`
  )
}

if (process.client) {
  window.addEventListener('message', (event: MessageEvent) => {
    if (event.origin !== window.location.origin) return

    if (event.data?.type === 'GOOGLE_AUTH_SUCCESS') {
      setAccessToken(event.data.accessToken)
      localStorage.setItem('session', '1')
      popupRef = null
      router.push('/dashboard')
    }

    if (event.data?.type === 'GOOGLE_AUTH_ERROR') {
      apiError.value = 'Google login failed. Please try again.'
      notify('Google login failed. Please try again.', 'error')
      popupRef = null
    }
  })
}

</script>

<template>
  <div>
    <RegistrationPage>
      <slot>
        <img
          src="https://res.cloudinary.com/sdq121/image/upload/v1755378345/wnrloghcmniaatz3d9ie.png"
          alt="7 orange"
          class="mx-auto"
        >
        <h2 class="text-center text-3xl font-bold mt-8">
          Create your Gro Account
        </h2>
        <div class="mt-10 py-12 bg-white border border-[#EDEDEE]  px-7 rounded-xl">
          <GroBasicInput
            v-model="email"
            placeholder="Enter email"
            :color="validationErrors?.email?.[0] ? 'error' : 'primary'"
            :hint="validationErrors?.email?.[0] ?? ''"
          >
            <h6>Email</h6>
          </GroBasicInput>
          <GroBasicPassword
            v-model="password"
            placeholder="Enter password"
            class="mt-4"
            :color="validationErrors?.password?.[0] ? 'error' : 'primary'"
            :hint="validationErrors?.password?.[0] ?? ''"
          >
            <h6>Password</h6>
          </GroBasicPassword>
          <div class="inline-flex w-full">
            <NuxtLink
              :to="{name : 'authentication-forgot-password'}"
              class="w-full text-right mt-2 ml-auto text-sm"
            >
              <span class="text-[#2176AE] cursor-pointer">Forgot Password</span>
            </NuxtLink>
          </div>
          <GroBasicButton
            color="primary"
            size="md"
            shape="custom"
            class="w-full mt-4 py-3"
            @click="registerUser"
          >
            Continue
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
            Already have an account?
            <span class="text-[#D16B07] cursor-pointer">
              <NuxtLink :to="{ name: 'authentication-login' }">Sign in</NuxtLink>
            </span>
          </h6>
        </div>
      </slot>
    </RegistrationPage>
  </div>
</template>

<style scoped>

</style>


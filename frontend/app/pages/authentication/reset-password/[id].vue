<!-- eslint-disable vue/multi-word-component-names -->
<script lang="ts" setup>

import BasicPage from "@/layouts/Authentication/BasicPage.vue";
import GroBasicButton from "@/components/buttons/GroBasicButton.vue";
import {ref} from "vue";
import {ResetPasswordSchema} from "@/composables/helpers/validation/authentication/authentication";
import {FormatZodErrors} from "@/composables/helpers/validation/util";
import {useAuthentication} from "@/composables/api/authentication/authentication";
import {ArrowLeft01Icon, CancelSquareIcon} from "@hugeicons/core-free-icons";
import {HugeiconsIcon} from "@hugeicons/vue";
import GroBasicPassword from "@/components/forms/input/GroBasicPassword.vue";
import {useRoute} from "nuxt/app";

const password = ref<string>("");
const confirmPassword = ref<string>("");
const passwordReset = ref<boolean>(false);
const validationErrors = ref<Record<string, string[]>>({});
const apiError = ref<string>("");

const route = useRoute()

const validateForgetUserPasswordForm = () => {
  validationErrors.value = {};
  const validationResult = ResetPasswordSchema.safeParse({
    password: password.value,
    confirm_password: confirmPassword.value,
    token: route.params.id,
  })
  validationErrors.value = validationResult?.success ?  {} : FormatZodErrors(validationResult.error!.issues)
}

const submitForgetUserPasswordForm = async () => {
  const loginApi = await useAuthentication().ResetPassword({
    password: password.value,
    confirm_password: confirmPassword.value,
    token: route.params?.id,
  })
  if (!loginApi?.success) {
    apiError.value = loginApi?.error
  } else {
    passwordReset.value = true
  }
}

const forgetUserPassword = async () => {
  apiError.value = "";
  validateForgetUserPasswordForm();
  if (Object.keys(validationErrors.value).length === 0) {
    await submitForgetUserPasswordForm()
  }
}

</script>

<template>
  <section>
    <BasicPage>
      <div class="bg-white border border-[#EDEDEE] py-10 px-7 rounded-xl">
        <svg
          class="mx-auto"
          width="44"
          height="48"
          viewBox="0 0 44 48"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            d="M21.002 0.480469C28.306 0.480512 34.2314 6.40595 34.2314 13.71V18.1201H35.7021C38.9452 18.1203 41.582 20.7569 41.582 24V31.9424C40.0477 30.7272 38.1092 30 36 30C31.0294 30 27 34.0294 27 39C27 42.9551 29.552 46.3121 33.0986 47.5195H6.30176C3.05863 47.5195 0.421875 44.8828 0.421875 41.6396V24C0.421881 20.7569 3.05863 18.1202 6.30176 18.1201H7.77148V13.71C7.77151 6.40593 13.6979 0.480469 21.002 0.480469ZM21.002 6.36035C16.9411 6.36035 13.6524 9.64911 13.6523 13.71V18.1201H28.3516V13.71C28.3515 9.64914 25.0628 6.36039 21.002 6.36035Z"
            fill="#D16B07"
          />
          <path
            d="M37.7552 31.244C38.0721 31.0739 38.4664 31.1938 38.636 31.5117L39.0251 32.2413L39.0409 32.2709C39.1588 32.4917 39.2716 32.7033 39.34 32.8732C39.3756 32.9617 39.4197 33.0873 39.4326 33.2274C39.4465 33.3779 39.4283 33.607 39.2607 33.8105C39.0856 34.0232 38.8535 34.0778 38.7113 34.0902C38.5739 34.1022 38.4425 34.0829 38.3462 34.0645C38.1594 34.0288 37.918 33.9569 37.6587 33.8796L37.6587 33.8796L37.6273 33.8703C37.1141 33.7174 36.5683 33.6349 36.0013 33.6349C32.9276 33.6349 30.4696 36.0553 30.4696 39C30.4696 39.9744 30.7368 40.8871 31.2049 41.6747C31.3889 41.9844 31.2878 42.3851 30.9791 42.5697C30.6704 42.7543 30.2709 42.6529 30.0869 42.3432C29.5029 41.3605 29.168 40.2179 29.168 39C29.168 35.2976 32.2459 32.3294 36.0013 32.3294C36.3682 32.3294 36.7287 32.3576 37.0805 32.4122L37.4295 31.6667C37.5061 31.4907 37.5848 31.3354 37.7552 31.244Z"
            fill="#D16B07"
          />
          <path
            d="M41.0235 35.4304C41.3322 35.2458 41.7317 35.3472 41.9157 35.6568C42.4997 36.6395 42.8346 37.7822 42.8346 39C42.8346 42.7025 39.7567 45.6707 36.0013 45.6707C35.6619 45.6707 35.328 45.6465 35.0013 45.5997L34.5732 46.3334C34.5034 46.4884 34.4178 46.6646 34.2474 46.756C33.9305 46.9261 33.5362 46.8063 33.3666 46.4884L32.9617 45.7292C32.8438 45.5083 32.731 45.2967 32.6626 45.1268C32.627 45.0383 32.5829 44.9127 32.57 44.7727C32.5561 44.6222 32.5743 44.3931 32.7419 44.1895C32.917 43.9768 33.1491 43.9223 33.2913 43.9099C33.4287 43.8979 33.5601 43.9172 33.6564 43.9356C33.8432 43.9713 34.0847 44.0432 34.344 44.1204L34.3753 44.1298C34.8885 44.2826 35.4343 44.3651 36.0013 44.3651C39.075 44.3651 41.533 41.9448 41.533 39C41.533 38.0256 41.2658 37.1129 40.7977 36.3253C40.6137 36.0156 40.7148 35.615 41.0235 35.4304Z"
            fill="#D16B07"
          />
        </svg>

        <h2 class="text-center text-3xl font-bold mt-6">
          {{ passwordReset ? "Password Reset Successful" : "Reset your password" }}
        </h2>
        <h6 class="text-center text-sm font-light text-[#4B4D55] mt-6 mb-8">
          {{ passwordReset ? "New password set. Click the below button to login" : "Kindly enter your new password." }}
        </h6>
        <GroBasicPassword
          v-if="!passwordReset"
          v-model="password"
          placeholder="Enter password"
          class="mt-4"
          :color="validationErrors?.password?.[0] ? 'error' : 'primary'"
          :hint="validationErrors?.password?.[0] ?? ''"
        >
          <h6>New Password</h6>
        </GroBasicPassword>
        <GroBasicPassword
          v-if="!passwordReset"
          v-model="confirmPassword"
          placeholder="Enter password"
          class="mt-4"
          :color="validationErrors?.password?.[0] ? 'error' : 'primary'"
          :hint="validationErrors?.password?.[0] ?? ''"
        >
          <h6>Confirm New Password</h6>
        </GroBasicPassword>
        <div
          v-if="apiError && !passwordReset"
          class="md:inline-flex w-full justify-between"
        >
          <div

            class="flex gap-1 mt-1 text-sm text-[#AF513A]"
          >
            <HugeiconsIcon
              color="#FFFFFF"
              fill="#AF513A"
              :icon="CancelSquareIcon"
            />

            <span class="my-auto">{{ apiError }}</span>
          </div>
        </div>
        <GroBasicButton
          v-if="!passwordReset"
          color="primary"
          size="md"
          shape="custom"
          class="w-full mt-6 py-3"
          @click="forgetUserPassword"
        >
          Reset Password
        </GroBasicButton>
        <NuxtLink
          v-else
          :to="{ name: 'login' }"
        >
          <GroBasicButton
            color="primary"
            size="md"
            shape="custom"
            class="w-full mt-6 py-3"
            @click="forgetUserPassword"
          >
            Login
          </GroBasicButton>
        </NuxtLink>
      </div>
      <NuxtLink
        :to="{ name: 'login' }"
        class="inline-flex w-full mt-5 cursor-pointer"
      >
        <HugeiconsIcon
          color="#4B4D55"
          :icon="ArrowLeft01Icon"
          class="ml-auto"
        />
        <h6 class="text-[#4B4D55] text-sm pl-2 my-auto mr-auto">
          Back to Sign In
        </h6>
      </NuxtLink>
    </BasicPage>
  </section>
</template>

<style scoped>

</style>

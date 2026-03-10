<!-- eslint-disable vue/multi-word-component-names -->
<script lang="ts" setup>

import BasicPage from "@/layouts/Authentication/BasicPage.vue";
import GroBasicInput from "@/components/forms/input/GroBasicInput.vue";
import GroBasicButton from "@/components/buttons/GroBasicButton.vue";
import {ref} from "vue";
import {ForgotPasswordSchema} from "@/composables/helpers/validation/authentication/authentication";
import {FormatZodErrors} from "@/composables/helpers/validation/util";
import {useAuthentication} from "@/composables/api/authentication/authentication";
import {ArrowLeft01Icon, CancelSquareIcon} from "@hugeicons/core-free-icons";
import {HugeiconsIcon} from "@hugeicons/vue";
import {maskEmail} from "@/composables/helpers/format/email";

const email = ref<string>("");
const emailSent = ref<boolean>(false);
const validationErrors = ref<Record<string, string[]>>({});
const apiError = ref<string>("");

const validateForgetUserPasswordForm = () => {
  validationErrors.value = {};
  const validationResult = ForgotPasswordSchema.safeParse({
    email: email.value,
  })
  validationErrors.value = validationResult?.success ?  {} : FormatZodErrors(validationResult.error!.issues)
}

const submitForgetUserPasswordForm = async () => {
  const loginApi = await useAuthentication().ForgotPassword({
    email: email.value,
  })
  if (!loginApi?.success) {
    apiError.value = loginApi?.error
  } else {
    emailSent.value = true
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
          v-if="!emailSent"
          class="mx-auto"
          width="42"
          height="48"
          viewBox="0 0 42 48"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            d="M13.6519 13.71V18.12H28.3519V13.71C28.3519 9.64911 25.0627 6.35998 21.0019 6.35998C16.941 6.35998 13.6519 9.64911 13.6519 13.71ZM7.77187 18.12V13.71C7.77187 6.40592 13.6978 0.47998 21.0019 0.47998C28.3059 0.47998 34.2319 6.40592 34.2319 13.71V18.12H35.7019C38.9451 18.12 41.5819 20.7568 41.5819 24V41.64C41.5819 44.8832 38.9451 47.52 35.7019 47.52H6.30187C3.05869 47.52 0.421875 44.8832 0.421875 41.64V24C0.421875 20.7568 3.05869 18.12 6.30187 18.12H7.77187Z"
            fill="#D16B07"
          />
        </svg>
        <svg
          v-else
          class="mx-auto"
          width="48"
          height="48"
          viewBox="0 0 48 48"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            fill-rule="evenodd"
            clip-rule="evenodd"
            d="M40.5279 0.335111C41.5164 -0.0257491 42.5868 -0.0984197 43.6153 0.125914C44.6538 0.352494 45.6056 0.872533 46.3572 1.62421C47.1087 2.37588 47.6288 3.32758 47.8554 4.36617C48.0797 5.39447 48.007 6.46515 47.6463 7.45336L35.3926 44.1803C35.094 45.0833 34.5677 45.8973 33.8652 46.5391C33.1647 47.1792 32.312 47.6287 31.3884 47.8454C30.4647 48.0707 29.4982 48.0515 28.5839 47.7902C27.6704 47.5293 26.8406 47.0356 26.1755 46.3574L19.6018 39.8136L12.6964 43.3845C12.158 43.6629 11.5125 43.6361 10.9988 43.3142C10.4851 42.9926 10.1796 42.4231 10.1953 41.8173L10.4784 30.9015L34.6325 13.3561C35.5901 12.6606 35.8026 11.3206 35.107 10.3631C34.4113 9.40552 33.0714 9.19316 32.1139 9.88868L7.54994 27.7317L1.61953 21.8012C0.977277 21.1595 0.503839 20.3676 0.242936 19.4979C-0.0161112 18.6344 -0.0581981 17.7205 0.120257 16.837C0.298942 15.8718 0.730399 14.9712 1.37079 14.2269C2.01557 13.4776 2.84876 12.914 3.78424 12.5945L3.7958 12.5906L40.5279 0.335111Z"
            fill="#D16B07"
          />
        </svg>


        <h2 class="text-center text-3xl font-bold mt-6">
          {{ emailSent ? "Your email is on the way" : "Forgot your password" }}
        </h2>
        <h6 class="text-center text-sm font-light text-[#4B4D55] mt-6 mb-8">
          {{ emailSent ? `Check your email ${maskEmail(email)} and follow the instructions to reset your password.` : "Don’t worry. Enter your email and we’ll send you a link to reset your password." }}
        </h6>
        <GroBasicInput
          v-if="!emailSent"
          v-model="email"
          placeholder="Enter email"
          :color="validationErrors?.email?.[0] ? 'error' : 'primary'"
          :hint="validationErrors?.email?.[0] ?? ''"
        >
          <h6>Email</h6>
        </GroBasicInput>
        <div
          v-if="apiError && !emailSent"
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
          v-if="!emailSent"
          color="primary"
          size="md"
          shape="custom"
          class="w-full mt-6 py-3"
          @click="forgetUserPassword"
        >
          Reset Password
        </GroBasicButton>
      </div>
      <NuxtLink
        :to="{ name: 'authentication-login' }"
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

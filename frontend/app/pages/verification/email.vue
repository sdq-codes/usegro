<!-- eslint-disable vue/multi-word-component-names -->
<script lang="ts" setup>

import BasicPage from "@/layouts/Authentication/BasicPage.vue";
import GroBasicButton from "@/components/buttons/GroBasicButton.vue";
import {onMounted, ref} from "vue";
import {useUserAPI} from "@/composables/api/user/user";
import {useVerificationEmailAPI} from "@/composables/api/verification/email";
import {CancelSquareIcon, ArrowLeft01Icon} from "@hugeicons/core-free-icons";
import {HugeiconsIcon} from "@hugeicons/vue";
import {useRouter} from "nuxt/app";
import GroBasicPin from "@/components/forms/input/GroBasicPin.vue";
import {verifications} from "@/composables/helpers/format/verifications";
import {notify} from "@/composables/helpers/notification/notification";

const router = useRouter()

const email = ref<string>("");
const verificationCode = ref<string>("");
const apiError = ref<string>("");
const canResend = ref<boolean>(false);
const resendTimer = ref<number>(10);
let timerInterval: NodeJS.Timeout | null = null;

onMounted(async () => {
  const fetchUserApiResponse = await useUserAPI().FetchLoggedInUser()
  if (fetchUserApiResponse.success) {
    email.value = fetchUserApiResponse.data?.data?.email;
    const userVerifications = verifications(fetchUserApiResponse.data?.data?.verifications);
    if (userVerifications?.email === "VERIFIED") {
      router.push("/dashboard");
      return;
    }
    startCooldown();
  } else {
    localStorage.clear()
    router.push("/authentication/login");
  }
})

const clearLocalStorage = () => {
  localStorage.clear()
}

// start cooldown timer
const startCooldown = () => {
  canResend.value = false;
  resendTimer.value = 10;
  if (timerInterval) clearInterval(timerInterval);
  timerInterval = setInterval(() => {
    resendTimer.value--;
    if (resendTimer.value <= 0) {
      canResend.value = true;
      if (timerInterval) clearInterval(timerInterval);
    }
  }, 1000);
};

// resend API call
const resendVerification = async () => {
  if (!canResend.value) return;

  try {
    const response = await useVerificationEmailAPI().ResendVerificationEmail();
    if (response.success) {
      notify(response.data?.response_message, 'success');
    } else {
      notify("Failed to resend verification email.", 'error');
    }
  } catch (err) {
    console.error(err);
    notify("Failed to resend verification email.", 'error');
  }

  startCooldown();
};

// verify email API call
const verifyEmail = async () => {
  const response = await useVerificationEmailAPI().VerifyEmail({ code: verificationCode.value});
  if (!response?.success) {
    apiError.value = response?.error
    notify(response?.error, 'error');
  } else {
    notify("Email successfully verified", 'success');
    router.push("/dashboard");
  }
};
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
        Verify your email
      </h2>
      <div class="mt-8 md:mt-16 bg-white border border-[#EDEDEE] py-10 px-7 rounded-xl">
        <h6 class="text-left text-[#6F7177] text-sm">
          Enter the code sent to
        </h6>
        <h6 class="text-left text-[#1E212B] text-sm mt-3">
          {{ email }}
        </h6>
        <GroBasicPin
          v-model="verificationCode"
          :length="6"
          color="primary"
          class="mt-8"
        />
        <div class="md:inline-flex w-full mt-2 justify-between">
          <div
            v-if="apiError"
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
          color="tertiary"
          size="md"
          shape="custom"
          class="w-full mt-3 py-3"
          :disabled="!canResend"
          @click="resendVerification"
        >
          Resend Verification Code {{ resendTimer !== 0 ? `(${resendTimer})` : "" }}
        </GroBasicButton>
        <GroBasicButton
          color="primary"
          size="md"
          shape="custom"
          class="w-full mt-3 py-3"
          :disabled="verificationCode.length !== 6"
          @click="verifyEmail"
        >
          Continue
        </GroBasicButton>
      </div>
      <NuxtLink
        :to="{ name: 'authentication-register' }"
        class="inline-flex w-full mt-5 cursor-pointer"
        @click="clearLocalStorage"
      >
        <HugeiconsIcon
          color="#4B4D55"
          :icon="ArrowLeft01Icon"
          class="ml-auto"
        />
        <h6 class="text-[#4B4D55] text-sm pl-2 my-auto mr-auto">
          Back to Sign Up
        </h6>
      </NuxtLink>
    </BasicPage>
  </section>
</template>

<style scoped>

</style>

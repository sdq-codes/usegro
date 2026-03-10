<script setup lang="ts">
import {
  MinusSignIcon,
  PlusSignIcon,
  ShoppingCart01Icon,
  UserMultipleIcon
} from "@hugeicons/core-free-icons";
import { PAYMENT_TYPE } from "@/constants/invoice/paymentType";
import GroBasicAlertWell from "@/components/alerts/GroBasicAlertWell.vue";
import GroBasicTextArea from "@/components/forms/input/GroBasicTextArea.vue";
import GroBasicRadio from "@/components/forms/select/GroBasicRadio.vue";
import { HugeiconsIcon } from "@hugeicons/vue";
import GroBasicSelect from "@/components/forms/select/GroBasicSelect.vue";
import { defineEmits, defineModel, defineProps, ref } from "vue";
import type { FormSubmission } from "@/composables/dto/customer/customer";
import GroBasicInput from "@/components/forms/input/GroBasicInput.vue";
import GroBasicUpload from "@/components/forms/GroBasicUpload.vue";

// Models shared with parent
const selectedCustomer = defineModel<Array<any>>("selectedCustomer");
const selectedProducts = defineModel<Array<any>>("selectedProducts");
const referenceNumber = defineModel<string>("referenceNumber");
const termsAndCondition = defineModel<string>("termsAndCondition");
const notesToCustomer = defineModel<string>("notesToCustomer");
const memos = defineModel<string>("memos");
const remainingReOccuringInvoice = defineModel<number>("remainingReOccuringInvoice");

// Local UI toggles (these stay local to the form)
const showReferenceNumber = ref<boolean>(false)
const showTermsAndCondition = ref<boolean>(false);
const showMemos = ref<boolean>(false);

const props = withDefaults(defineProps<{
  crmCustomersFullList?: FormSubmission[] | null,
}>(), {
  crmCustomersFullList: null,
})

const emit = defineEmits<{
  (e: 'createNewCustomer'): void
}>()

const createNewCustomer = () => {
  emit('createNewCustomer')
}
</script>

<template>
  <div class="w-1/2 overflow-y-auto bg-white py-7 rounded-l-2xl border-y border-l border-[#EDEDEE]">
    <div class="py-10 px-8 border-b border-[#EDEDEE]">
      <div class="flex gap-x-2">
        <HugeiconsIcon
          :icon="UserMultipleIcon"
          color="#4D91BE"
        />
        <h6 class="text-sm font-semibold my-auto">
          Who are you billing?
        </h6>
      </div>
      <GroBasicSelect
        :show-initials="true"
        v-model="selectedCustomer"
        :options="props.crmCustomersFullList"
        add-more-text="New customer"
        placeholder="Choose a customer"
        color="primary"
        class="mt-4"
        :multiple-select="true"
        @add-new="createNewCustomer"
      >
        Customer(s)*
      </GroBasicSelect>
    </div>

    <div class="py-10 px-8 border-b border-[#EDEDEE]">
      <div class="flex gap-x-2">
        <HugeiconsIcon
          :icon="ShoppingCart01Icon"
          color="#D26B06"
        />
        <h6 class="text-sm font-semibold my-auto">
          What are they paying for
        </h6>
      </div>
      <GroBasicSelect
        v-model="selectedProducts"
        :multiple-select="true"
        :options="props.crmCustomersFullList"
        placeholder="Type, select or search for products & services"
        color="primary"
        class="mt-4"
        add-more-text="Create New Item"
      >
        Items
      </GroBasicSelect>
    </div>

    <div class="py-10 px-8 border-b border-[#EDEDEE]">
      <div class="flex gap-x-2">
        <h6 class="text-sm font-semibold my-auto">
          Payments Collection
        </h6>
      </div>
      <div class="mt-4">
        <GroBasicRadio
          :options="PAYMENT_TYPE"
          layout="vertical"
          border-class="pb-6 border-b border-[#EDEDEE] last:border-b-0"
        />
        <GroBasicAlertWell
          variant="success"
          class="mt-4"
        >
          <template #text>
            You have <span class="font-semibold"> {{ remainingReOccuringInvoice || 0 }} recurring payments plans</span> left.
            <a
              class="font-semibold underline text-blue-400"
              href="#"
            >Gro Pro</a>
          </template>
        </GroBasicAlertWell>
      </div>
    </div>

    <div class="py-10 px-8">
      <div class="flex gap-x-2">
        <h6 class="text-sm font-semibold my-auto">
          Notes & Attachments
        </h6>
      </div>
      <div class="mt-6">
        <GroBasicTextArea
          v-model="notesToCustomer"
          placeholder="Note to your customer"
          border-radius="8px"
        >
          <div class="flex w-full">
            Notes
            <small class="text-xs ml-auto text-gray-500">{{ notesToCustomer?.length || 0 }}/100</small>
          </div>
        </GroBasicTextArea>

        <div class="mt-10 space-y-10">
          <div
            class="flex cursor-pointer gap-x-2"
            @click="showReferenceNumber = !showReferenceNumber"
          >
            <HugeiconsIcon
              :icon="showReferenceNumber ? MinusSignIcon : PlusSignIcon"
              :size="20"
              color="#2176AE"
            />
            <h6 class="my-auto font-semibold text-[#2176AE]">
              Reference Number
            </h6>
          </div>
          <GroBasicInput
            v-if="showReferenceNumber"
            v-model="referenceNumber"
            class="-mt-8"
            placeholder="Reference Number"
          />

          <div
            class="flex cursor-pointer gap-x-2"
            @click="showTermsAndCondition = !showTermsAndCondition"
          >
            <HugeiconsIcon
              :icon="showTermsAndCondition ? MinusSignIcon : PlusSignIcon"
              :size="20"
              color="#2176AE"
            />
            <h6 class="my-auto font-semibold text-[#2176AE]">
              Terms and conditions
            </h6>
          </div>
          <GroBasicTextArea
            v-if="showTermsAndCondition"
            v-model="termsAndCondition"
            placeholder="Terms and conditions"
            border-radius="8px"
            class="-mt-8"
          />

          <div
            class="flex cursor-pointer gap-x-2"
            @click="showMemos = !showMemos"
          >
            <HugeiconsIcon
              :icon="showMemos ? MinusSignIcon : PlusSignIcon"
              :size="20"
              color="#2176AE"
            />
            <h6 class="my-auto font-semibold text-[#2176AE]">
              Memo to self
            </h6>
          </div>
          <GroBasicTextArea
            v-if="showMemos"
            v-model="memos"
            placeholder="Memo to self"
            border-radius="8px"
            class="-mt-8"
          />
        </div>

        <div class="mt-10 space-y-5">
          <h6 class="text-sm font-semibold text-[#1E212B]">
            Attachments
          </h6>
          <GroBasicUpload />
        </div>
      </div>
    </div>
  </div>
</template>

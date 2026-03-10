<!-- eslint-disable vue/multi-word-component-names -->
<template>
  <div class="p-6 space-y-6">
    <h1 class="text-2xl font-bold mb-4">
      Component Showcase
    </h1>

    <div
      class="border rounded-lg shadow p-4 gap-y-5"
    >
      <!-- Header -->
      <h2 class="text-lg font-semibold mb-5">
        Buttons
      </h2>
      <div class="inline-flex space-x-5 mb-5">
        <GroBasicButton
          color="primary"
          size="md"
          shape="custom"
        >
          Delete
        </GroBasicButton>
        <GroBasicButton
          color="primary"
          size="md"
          shape="custom"
        >
          <slot name="frontIcon">
            <HugeiconsIcon
              :icon="PlusSignIcon"
              :size="16"
              color="currentColor"
              :stroke-width="3"
            />
          </slot>
          Delete
          <slot name="backIcon">
            <HugeiconsIcon
              :icon="PlusSignIcon"
              :size="16"
              color="currentColor"
              :stroke-width="3"
            />
          </slot>
        </GroBasicButton>
        <GroBasicButton
          color="primary"
          size="md"
          shape="custom"
          :loading="true"
        >
          Delete
        </GroBasicButton>
        <GroBasicButton
          color="primary"
          size="md"
          shape="custom"
          :disabled="true"
        >
          <slot name="frontIcon">
            <HugeiconsIcon
              :icon="PlusSignIcon"
              :size="16"
              color="currentColor"
              :stroke-width="3"
            />
          </slot>
          Delete
          <slot name="backIcon">
            <HugeiconsIcon
              :icon="PlusSignIcon"
              :size="16"
              color="currentColor"
              :stroke-width="3"
            />
          </slot>
        </GroBasicButton>
        <GroBasicButton
          color="primary"
          size="md"
          shape="custom"
        >
          <slot name="backIcon">
            <HugeiconsIcon
              :icon="StarIcon"
              :size="16"
              color="currentColor"
              :stroke-width="3"
            />
          </slot>
        </GroBasicButton>
        <GroBasicButton
          color="primary"
          size="md"
          shape="round"
        >
          <slot name="backIcon">
            <HugeiconsIcon
              :icon="StarIcon"
              :size="16"
              color="currentColor"
              :stroke-width="3"
            />
          </slot>
        </GroBasicButton>
      </div>
      <div>
        <div class="inline-flex space-x-5 mb-5">
          <GroBasicButton
            color="secondary"
            size="md"
            shape="custom"
          >
            Delete
          </GroBasicButton>
          <GroBasicButton
            color="secondary"
            size="md"
            shape="custom"
          >
            <slot name="frontIcon">
              <HugeiconsIcon
                :icon="PlusSignIcon"
                :size="16"
                color="currentColor"
                fill="currentColor"
                :stroke-width="3"
              />
            </slot>
            Delete
            <slot name="backIcon">
              <HugeiconsIcon
                :icon="PlusSignIcon"
                :size="16"
                color="currentColor"
                fill="currentColor"
                :stroke-width="3"
              />
            </slot>
          </GroBasicButton>
          <GroBasicButton
            color="secondary"
            size="md"
            shape="custom"
            :loading="true"
          >
            Delete
          </GroBasicButton>
          <GroBasicButton
            color="secondary"
            size="md"
            shape="custom"
            :disabled="true"
          >
            <slot name="frontIcon">
              <HugeiconsIcon
                :icon="PlusSignIcon"
                :size="16"
                color="currentColor"
                fill="currentColor"
                :stroke-width="3"
              />
            </slot>
            Delete
            <slot name="backIcon">
              <HugeiconsIcon
                :icon="PlusSignIcon"
                :size="16"
                color="currentColor"
                fill="currentColor"
                :stroke-width="3"
              />
            </slot>
          </GroBasicButton>
          <GroBasicButton
            color="secondary"
            size="md"
            shape="custom"
          >
            <slot name="backIcon">
              <HugeiconsIcon
                :icon="StarIcon"
                :size="16"
                color="currentColor"
                fill="currentColor"
                :stroke-width="3"
              />
            </slot>
          </GroBasicButton>
          <GroBasicButton
            color="secondary"
            size="md"
            shape="round"
          >
            <slot name="backIcon">
              <HugeiconsIcon
                :icon="StarIcon"
                :size="16"
                color="currentColor"
                fill="currentColor"
                :stroke-width="3"
              />
            </slot>
          </GroBasicButton>
        </div>
      </div>
      <div>
        <div class="inline-flex space-x-5 mb-5">
          <GroBasicButton
            color="tertiary"
            size="md"
            shape="custom"
          >
            Delete
          </GroBasicButton>
          <GroBasicButton
            color="tertiary"
            size="md"
            shape="custom"
          >
            <slot name="frontIcon">
              <HugeiconsIcon
                :icon="PlusSignIcon"
                :size="16"
                color="currentColor"
                fill="currentColor"
                :stroke-width="3"
              />
            </slot>
            Delete
            <slot name="backIcon">
              <HugeiconsIcon
                :icon="PlusSignIcon"
                :size="16"
                color="currentColor"
                fill="currentColor"
                :stroke-width="3"
              />
            </slot>
          </GroBasicButton>
          <GroBasicButton
            color="tertiary"
            size="md"
            shape="custom"
            :loading="true"
          >
            Delete
          </GroBasicButton>
          <GroBasicButton
            color="tertiary"
            size="md"
            shape="custom"
            :disabled="true"
          >
            <slot name="frontIcon">
              <HugeiconsIcon
                :icon="PlusSignIcon"
                :size="16"
                color="currentColor"
                fill="currentColor"
                :stroke-width="3"
              />
            </slot>
            Delete
            <slot name="backIcon">
              <HugeiconsIcon
                :icon="PlusSignIcon"
                :size="16"
                color="currentColor"
                fill="currentColor"
                :stroke-width="3"
              />
            </slot>
          </GroBasicButton>
          <GroBasicButton
            color="tertiary"
            size="md"
            shape="custom"
          >
            <slot name="backIcon">
              <HugeiconsIcon
                :icon="StarIcon"
                :size="16"
                color="currentColor"
                fill="currentColor"
                :stroke-width="3"
              />
            </slot>
          </GroBasicButton>
          <GroBasicButton
            color="tertiary"
            size="md"
            shape="round"
          >
            <slot name="backIcon">
              <HugeiconsIcon
                :icon="StarIcon"
                :size="16"
                color="currentColor"
                fill="currentColor"
                :stroke-width="3"
              />
            </slot>
          </GroBasicButton>
        </div>
      </div>
    </div>
    <div
      class="border rounded-lg shadow p-4 gap-y-5"
    >
      <!-- Header -->
      <h2 class="text-lg font-semibold mb-5">
        Input Texts
      </h2>
      <div class="inline-flex space-x-5 gap-x-3 mb-5">
        <GroBasicInput
          placeholder="Type your query..."
          hint="We’ll search across all categories"
          border-radius="8px"
        >
          <slot>
            Search
          </slot>
        </GroBasicInput>
        <GroBasicInput
          placeholder="Type your query..."
          hint="We’ll search across all categories"
          border-radius="8px"
          color="success"
        >
          <slot>
            Search
          </slot>
        </GroBasicInput>
        <GroBasicInput
          placeholder="Type your query..."
          hint="We’ll search across all categories"
          border-radius="8px"
          color="error"
        >
          <slot>
            Search
          </slot>
        </GroBasicInput>
        <GroBasicInput
          placeholder="Type your query..."
          hint="We’ll search across all categories"
          border-radius="8px"
          color="success"
          :readonly="true"
        >
          <slot>
            Search
          </slot>
        </GroBasicInput>
        <GroBasicInput
          placeholder="Type your query..."
          hint="We’ll search across all categories"
          border-radius="8px"
          color="primary"
          :disabled="true"
        >
          <slot>
            Search
          </slot>
        </GroBasicInput>
        <GroBasicInput
          placeholder="Type your query..."
          hint="We’ll search across all categories"
          border-radius="8px"
          color="primary"
          :loading="true"
        >
          <slot>
            Search
          </slot>
        </GroBasicInput>
      </div>
    </div>
    <div
      class="border rounded-lg shadow p-4 gap-y-5"
    >
      <!-- Header -->
      <h2 class="text-lg font-semibold mb-5">
        PIN
      </h2>
      <div class="inline-flex space-x-5 mb-5">
        <GroBasicPin
          :length="4"
          color="primary"
          hint="This is a pin"
        >
          <slot>
            Search
          </slot>
        </GroBasicPin>
        <GroBasicPin
          :length="4"
          color="primary"
          disabled="disabled"
          hint="This is a pin"
        >
          <slot>
            Search
          </slot>
        </GroBasicPin>
        <GroBasicPin
          :length="4"
          color="success"
          hint="This is a pin"
        >
          <slot>
            Search
          </slot>
        </GroBasicPin>
        <GroBasicPin
          :length="4"
          color="error"
          hint="This is a pin"
        >
          <slot>
            Search
          </slot>
        </GroBasicPin>
      </div>
    </div>
    <div
      class="border rounded-lg shadow p-4 gap-y-5"
    >
      <!-- Header -->
      <h2 class="text-lg font-semibold mb-5">
        Select
      </h2>
      <div class="inline-flex space-x-5 gap-x-5 mb-5">
        <GroBasicSelect
          v-model="selectedCountryOptions"
          :options="countryOptions"
          placeholder="Select a country"
          hint="Choose your country of residence"
          color="primary"
        >
          Country
        </GroBasicSelect>
        <GroTelephoneSelect
          v-model="country"
          v-model:telephone="telephone"
          placeholder="Text placeholder"
          hint="Hint"
          :char-limit="100"
        >
          Label
        </GroTelephoneSelect>
      </div>
    </div>
    <div
      class="border rounded-lg shadow p-4 gap-y-5"
    >
      <!-- Header -->
      <h2 class="text-lg font-semibold mb-5">
        Password
      </h2>
      <div class="inline-flex space-x-5 gap-x-3 mb-5">
        <GroBasicPassword
          placeholder="Type your query..."
          hint="We’ll search across all categories"
          border-radius="8px"
        >
          <slot>
            Search
          </slot>
        </GroBasicPassword>
        <GroBasicPassword
          placeholder="Type your query..."
          hint="We’ll search across all categories"
          border-radius="8px"
          color="success"
        >
          <slot>
            Search
          </slot>
        </GroBasicPassword>
        <GroBasicPassword
          placeholder="Type your query..."
          hint="We’ll search across all categories"
          border-radius="8px"
          color="error"
        >
          <slot>
            Search
          </slot>
        </GroBasicPassword>
        <GroBasicPassword
          placeholder="Type your query..."
          hint="We’ll search across all categories"
          border-radius="8px"
          color="success"
          :readonly="true"
        >
          <slot>
            Search
          </slot>
        </GroBasicPassword>
        <GroBasicPassword
          placeholder="Type your query..."
          hint="We’ll search across all categories"
          border-radius="8px"
          color="primary"
          :disabled="true"
        >
          <slot>
            Search
          </slot>
        </GroBasicPassword>
        <GroBasicPassword
          placeholder="Type your query..."
          hint="We’ll search across all categories"
          border-radius="8px"
          color="primary"
          :loading="true"
        >
          <slot>
            Search
          </slot>
        </GroBasicPassword>
      </div>
    </div>
    <div
      class="border rounded-lg shadow p-4 gap-y-5"
    >
      <!-- Header -->
      <h2 class="text-lg font-semibold mb-5">
        Search
      </h2>
      <div class="inline-flex space-x-5 gap-x-3 mb-5">
        <GroBasicSearch
          placeholder="Type your query..."
          hint="We’ll search across all categories"
          border-radius="8px"
        >
          <slot>
            Search
          </slot>
        </GroBasicSearch>
        <GroBasicSearch
          placeholder="Type your query..."
          hint="We’ll search across all categories"
          border-radius="8px"
          color="success"
        >
          <slot>
            Search
          </slot>
        </GroBasicSearch>
        <GroBasicSearch
          placeholder="Type your query..."
          hint="We’ll search across all categories"
          border-radius="8px"
          color="error"
        >
          <slot>
            Search
          </slot>
        </GroBasicSearch>
        <GroBasicSearch
          placeholder="Type your query..."
          hint="We’ll search across all categories"
          border-radius="8px"
          color="success"
          :readonly="true"
        >
          <slot>
            Search
          </slot>
        </GroBasicSearch>
        <GroBasicSearch
          placeholder="Type your query..."
          hint="We’ll search across all categories"
          border-radius="8px"
          color="primary"
          :disabled="true"
        >
          <slot>
            Search
          </slot>
        </GroBasicSearch>
        <GroBasicSearch
          placeholder="Type your query..."
          hint="We’ll search across all categories"
          border-radius="8px"
          color="primary"
          :loading="true"
        >
          <slot>
            Search
          </slot>
        </GroBasicSearch>
      </div>
    </div>
    <div
      class="border rounded-lg shadow p-4 gap-y-5"
    >
      <!-- Header -->
      <h2 class="text-lg font-semibold mb-5">
        Textarea
      </h2>
      <div class="inline-flex space-x-5 gap-x-3 mb-5">
        <GroBasicTextArea
          placeholder="Type your query..."
          hint="We’ll search across all categories"
          border-radius="8px"
        >
          <slot>
            Search
          </slot>
        </GroBasicTextArea>
        <GroBasicTextArea
          placeholder="Type your query..."
          hint="We’ll search across all categories"
          border-radius="8px"
          color="success"
        >
          <slot>
            Search
          </slot>
        </GroBasicTextArea>
        <GroBasicTextArea
          placeholder="Type your query..."
          hint="We’ll search across all categories"
          border-radius="8px"
          color="error"
        >
          <slot>
            Search
          </slot>
        </GroBasicTextArea>
        <GroBasicTextArea
          placeholder="Type your query..."
          hint="We’ll search across all categories"
          border-radius="8px"
          color="success"
          :readonly="true"
        >
          <slot>
            Search
          </slot>
        </GroBasicTextArea>
        <GroBasicTextArea
          placeholder="Type your query..."
          hint="We’ll search across all categories"
          border-radius="8px"
          color="primary"
          :disabled="true"
        >
          <slot>
            Search
          </slot>
        </GroBasicTextArea>
        <GroBasicTextArea
          placeholder="Type your query..."
          hint="We’ll search across all categories"
          border-radius="8px"
          color="primary"
          :loading="true"
        >
          <slot>
            Search
          </slot>
        </GroBasicTextArea>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { HugeiconsIcon } from '@hugeicons/vue'
import { PlusSignIcon, StarIcon } from '@hugeicons/core-free-icons'

import GroBasicButton from "@/components/buttons/GroBasicButton.vue"
import GroBasicInput from "@/components/forms/input/GroBasicInput.vue";
import GroBasicPin from "@/components/forms/input/GroBasicPin.vue";
import GroBasicSelect from "@/components/forms/select/GroBasicSelect.vue";
import {ref} from "vue";
import GroTelephoneSelect from "@/components/forms/select/GroTelephoneSelect.vue";
import GroBasicPassword from "@/components/forms/input/GroBasicPassword.vue";
import GroBasicSearch from "@/components/forms/input/GroBasicSearch.vue";
import GroBasicTextArea from "@/components/forms/input/GroBasicTextArea.vue";

const countryOptions = [
  { label: 'Nigeria', value: 'NG' },
  { label: 'United States', value: 'US' },
  { label: 'Canada', value: 'CA' },
]

const selectedCountryOptions = ref<string>("")

const country = ref('us')
const telephone = ref('')

interface CountryOption {
  value: string
  label: string
  flag: string
}

const countryTelephoneOptions: CountryOption[] = [
  { value: 'us', label: '+1', flag: '🇺🇸' },
  { value: 'ng', label: '+234', flag: '🇳🇬' },
  { value: 'gb', label: '+44', flag: '🇬🇧' }
]
</script>

<script setup lang="ts">

import {computed, onMounted, ref, defineModel, defineEmits} from "vue";
import {PlusSignIcon} from "@hugeicons/core-free-icons";
import BasicModal from "@/components/modals/basic-modal.vue";
import DynamicForm from "@/components/forms/form/DynamicForm.vue";
import {HugeiconsIcon} from "@hugeicons/vue";
import GroBasicButton from "@/components/buttons/GroBasicButton.vue";
import GroBasicDropUp from "@/components/dropdown/GroBasicDropUp.vue";
import {useFormAPI} from "@/composables/api/customer/forms/create";
import type {FormField} from "@/composables/helpers/types/form";
import {notify} from "@/composables/helpers/notification/notification";
import {useFormFieldManager} from "@/composables/helpers/managers/forms/useFormFieldManager";
import {validateFormWithConditions} from "@/composables/helpers/validation/forms/createCustomer";

interface CustomerEditData {
  submissionId: string
  formId: string
  answers: Record<string, unknown>
}

interface Props {
  customer?: CustomerEditData
}

const props = withDefaults(defineProps<Props>(), {
  customer: undefined
})

const isEditMode = computed(() => !!props.customer)

const { getDefaultValue, addField } = useFormFieldManager()

const model = defineModel<boolean>();

const dynamicFormRef = ref<InstanceType<typeof DynamicForm> | null>(null)
const formVersion = ref({ title: '' })
const isLoading = ref(false)
const formFields = ref<FormField[]>([])
const answerMap = ref<Record<string, string | string[]>>({})
const validationErrors = ref<Record<string, string | string[]>>({});

const updateCustomerApi = async () => {
  if (!props.customer) return false
  const result = await useFormAPI().UpdateSubmission(
    props.customer.formId,
    props.customer.submissionId,
    answerMap.value as Record<string, unknown>
  )
  emit('customer-updated')
  return result.success
}

const createCustomerApi = async () => {
  const formData = {
    type: "customer",
    answers: answerMap.value,
    versionSnap: formFields.value,
  }
  const createdCustomer = await useFormAPI().CreateCustomer(formVersion.value?.formID, formVersion.value?.id, formData)
  emit('new-customer-added')
  return createdCustomer.success
}

const validateCustomer = () => {
  validationErrors.value = {};

  const validation = validateFormWithConditions(formFields.value, answerMap.value);

  if (!validation.success) {
    validationErrors.value = validation.errors;
  }

  if (!answerMap.value['phone_number'] && !answerMap.value['email']) {
    validationErrors.value['phone_number'] = "Phone number or email address must be filled";
    validationErrors.value['email'] = "Phone number or email address must be filled";
  }


  if (
    answerMap.value['customer_type'] === 'Individual' &&
    (!answerMap.value['first_name'] || !answerMap.value['last_name'])
  ) {
    validationErrors.value['first_name'] = "First name must be filled for individual customer types";
    validationErrors.value['last_name'] = "Last name must be filled for individual customer types";
  }

  if (answerMap.value['customer_type'] === 'Business' && !answerMap.value['company_name']) {
    validationErrors.value['company_name'] = "Company name must be filled for business customer types";
  }
};

const initForm = async () => {
  isLoading.value = true

  try {
    const fetchUserApiResponse = await useFormAPI().FetchForm()

    if (fetchUserApiResponse.success) {
      const data = fetchUserApiResponse.data?.data
      formVersion.value = data.version
      formFields.value = data.fields

      // Initialize answerMap with default values
      answerMap.value = data.fields.reduce((acc: Record<string, string | string[]>, field: FormField) => {
        acc[field.slug] = getDefaultValue(field)
        return acc
      }, {})

      // In edit mode, overlay with existing answers
      if (isEditMode.value && props.customer) {
        answerMap.value = {
          ...answerMap.value,
          ...props.customer.answers as Record<string, string | string[]>,
        }
      }
    }
  } catch (error) {
    console.error('Error fetching form:', error)
  } finally {
    isLoading.value = false
  }
}

const saveCustomer = async () => {
  validateCustomer()
  if (Object.keys(validationErrors.value).length !== 0) {
    return
  }

  if (isEditMode.value) {
    const success = await updateCustomerApi()
    if (success) {
      notify("Customer was successfully updated", "success")
      model.value = false
    } else {
      notify("Customer could not be updated", "error")
    }
  } else {
    const customerCreation = await createCustomerApi();
    if (customerCreation) {
      notify("Customer was successfully created", "success")
      model.value = false;
    } else {
      notify("Customer could not be created", "error")
    }
  }
}

const saveAndCreateAnother = async () => {
  validateCustomer()
  if (Object.keys(validationErrors.value).length !== 0) {
    return
  } else {
    const customerCreation = await createCustomerApi();
    if (customerCreation) {
      notify("Customer was successfully created", "success")
      await initForm()
    } else {
      notify("Customer could not be created", "error")
    }
  }
}

const handleAddField = (type: string) => {
  const newField = addField(type, formFields.value, answerMap.value)
  if (!newField) return
  if (Array.isArray(newField)) {
    formFields.value.push(...newField)
  } else {
    formFields.value.push(newField)
  }
}

const emit = defineEmits<{
  (e: 'new-customer-added'): void
  (e: 'customer-updated'): void
}>()

onMounted(async () => {
  await initForm()
})

</script>

<template>
  <div>
    <BasicModal
      v-model="model"
      size="xl"
    >
      <template #title>
        {{ isEditMode ? 'Edit Customer' : formVersion.title }}
      </template>
      <template #default>
        <DynamicForm
          ref="dynamicFormRef"
          v-model="answerMap"
          :fields="formFields"
          :title="formVersion.title"
          :is-loading="isLoading"
          layout="horizontal"
          :errors="validationErrors"
          @update:fields="formFields = $event"
        />
      </template>
      <template #footer>
        <div class="flex">
          <div class="flex w-full h-full">
            <div
              v-if="!isEditMode"
              class="relative"
            >
              <GroBasicDropUp>
                <template #button>
                  <div class="flex cursor-pointer my-1">
                    <HugeiconsIcon
                      :icon="PlusSignIcon"
                      :size="12"
                      color="#070707"
                      :stroke-width="3"
                      class="my-auto"
                    />
                    <h6 class="text-xs border-b text-blue-500">
                      Add New Fields
                    </h6>
                  </div>
                </template>
                <template #menu-list>
                  <li>
                    <a
                      class="text-slate-800 hover:bg-slate-50 flex items-center p-2 cursor-pointer"
                      @click.prevent="handleAddField('email')"
                    >
                      <span class="whitespace-nowrap">Email</span>
                    </a>
                  </li>
                  <li>
                    <a
                      class="text-slate-800 hover:bg-slate-50 flex items-center p-2 cursor-pointer"
                      @click.prevent="handleAddField('phone')"
                    >
                      <span class="whitespace-nowrap">Phone</span>
                    </a>
                  </li>
                  <li>
                    <a
                      class="text-slate-800 hover:bg-slate-50 flex items-center p-2 cursor-pointer"
                      @click.prevent="handleAddField('address')"
                    >
                      <span class="whitespace-nowrap">Address</span>
                    </a>
                  </li>
                  <li>
                    <a
                      class="text-slate-800 hover:bg-slate-50 flex items-center p-2 cursor-pointer"
                      @click.prevent="handleAddField('company')"
                    >
                      <span class="whitespace-nowrap">Company info</span>
                    </a>
                  </li>
                  <li>
                    <a
                      class="text-slate-800 hover:bg-slate-50 flex items-center p-2 cursor-pointer"
                      @click.prevent="handleAddField('position')"
                    >
                      <span class="whitespace-nowrap">Position</span>
                    </a>
                  </li>
                  <li>
                    <a
                      class="text-slate-800 hover:bg-slate-50 flex items-center p-2 cursor-pointer"
                      @click.prevent="handleAddField('birthdate')"
                    >
                      <span class="whitespace-nowrap">Date</span>
                    </a>
                  </li>
                  <li>
                    <a
                      class="text-slate-800 hover:bg-slate-50 flex items-center p-2 cursor-pointer"
                      @click.prevent="handleAddField('notes')"
                    >
                      <span class="whitespace-nowrap">Notes</span>
                    </a>
                  </li>
                  <li>
                    <a
                      class="text-slate-800 hover:bg-slate-50 flex items-center p-2 cursor-pointer"
                      @click.prevent="handleAddField('custom')"
                    >
                      <span class="whitespace-nowrap">Custom fields</span>
                    </a>
                  </li>
                </template>
              </GroBasicDropUp>
            </div>
          </div>
          <div class="ml-auto gap-x-4 flex">
            <div
              v-if="!isEditMode"
              class="hidden md:block"
            >
              <GroBasicButton
                color="secondary"
                size="sm"
                shape="custom"
                class="w-max h-max"
                @click="saveAndCreateAnother"
              >
                Save & create another
              </GroBasicButton>
            </div>
            <GroBasicButton
              color="primary"
              size="sm"
              shape="custom"
              class="w-max h-max"
              @click="saveCustomer"
            >
              {{ isEditMode ? 'Update' : 'Save' }}
            </GroBasicButton>
          </div>
        </div>
      </template>
    </BasicModal>
  </div>
</template>

<style scoped>

</style>

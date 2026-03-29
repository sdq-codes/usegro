<script setup lang="ts">
import {onMounted, ref, onBeforeUnmount} from 'vue'
import MainDashboard from "@/components/dashboard/main-dashboard.vue";
import {useCustomerAPI} from "@/composables/api/customer/customer";

interface CustomerActivity {
  id: string
  activityType: string
  description: string
  createdAt: string
}
import {useFormAPI} from "@/composables/api/customer/forms/create";
import {useRoute} from "nuxt/app";
import GroBasicButton from "@/components/buttons/GroBasicButton.vue";
import {CallIcon, Copy01Icon, GreaterThanIcon, Mail01Icon} from "@hugeicons/core-free-icons";
import {HugeiconsIcon} from "@hugeicons/vue";
import BackLink from "@/components/navigation/BackLink.vue";
import {notify} from "@/composables/helpers/notification/notification";
import GroReadOnlyTags from "@/components/tags/GroReadOnlyTags.vue";
import GroBasicTagSelect from "@/components/forms/select/GroBasicTagSelect.vue";
import CreateCustomerModal from "@/components/modals/modal/create-customer-modal.vue"
import CommentEditor from "@/components/editor/CommentEditor.vue";

const customerProfile = ref()
const rawSubmission = ref()
const extraFields = ref<{ label: string; slug: string; value: unknown }[]>([])

const comment = ref('')
const showTagsEdit = ref(false)
const showEditModal = ref(false)
const showInfoMenu = ref(false)

const closeInfoMenu = (e: MouseEvent) => {
  if (!(e.target as Element).closest('.info-menu-container')) {
    showInfoMenu.value = false
  }
}
onMounted(() => document.addEventListener('click', closeInfoMenu))
onBeforeUnmount(() => document.removeEventListener('click', closeInfoMenu))
const editTagIds = ref<string[]>([])
const isSavingTags = ref(false)
const activities = ref<CustomerActivity[]>([])

const route = useRoute()

onMounted(async () => {
  await fetchCustomer()
})

const fetchCustomer = async () => {
  const createdCustomer = await useCustomerAPI().FetchCustomer(String(route.params?.submission ?? ""), String(route.params?.id ?? ""))
  if (!createdCustomer.success) {
    return;
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const submission = (createdCustomer.data as any)?.data
  rawSubmission.value = submission

  const answers = submission?.answers ?? {}

  customerProfile.value = {
    ...answers,
    customerId: submission?._id,
    createdAt: submission?.createdAt,
    formId: submission?.formID,
  }

  editTagIds.value = [...(answers?.customer_tags ?? [])]

  const activityResponse = await useCustomerAPI().FetchCustomerActivity(submission?._id)
  if (activityResponse.success) {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    activities.value = ((activityResponse as any).data?.data ?? []) as CustomerActivity[]
  }

  // Parse versionSnap extra fields — each field is an array of {Key, Value} pairs
  const versionSnap: Array<Array<{ Key: string; Value: unknown }>> = submission?.versionSnap ?? []
  extraFields.value = versionSnap
    .map(fieldPairs => Object.fromEntries(fieldPairs.map(({ Key, Value }) => [Key, Value])))
    .filter(f => f.section === 'Extra fields' && answers[f.slug as string] !== undefined && answers[f.slug as string] !== '')
    .map(f => ({ label: f.label as string || f.slug as string, slug: f.slug as string, value: answers[f.slug as string] }))
}

const saveTagsUpdate = async () => {
  if (!rawSubmission.value) return
  isSavingTags.value = true
  try {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const { customerId: _cid, createdAt: _ca, formId: _fid, ...answerFields } = customerProfile.value
    const updatedAnswers = { ...answerFields, customer_tags: editTagIds.value }

    const result = await useFormAPI().UpdateSubmission(
      rawSubmission.value.formID,
      rawSubmission.value._id,
      updatedAnswers
    )
    if (result.success) {
      customerProfile.value.customer_tags = editTagIds.value
      notify("Tags updated successfully", "success")
      showTagsEdit.value = false
    } else {
      notify("Failed to update tags", "error")
    }
  } finally {
    isSavingTags.value = false
  }
}

function formatCustomerDuration(timestamp: string): string {
  const now = new Date();
  const createdAt = new Date(timestamp);
  const diffMs = now.getTime() - createdAt.getTime();

  // Convert to seconds/minutes/hours/days
  const seconds = Math.floor(diffMs / 1000);
  const minutes = Math.floor(seconds / 60);
  const hours = Math.floor(minutes / 60);
  const days = Math.floor(hours / 24);

  let duration: string;

  if (seconds < 5) duration = "less than 5 seconds";
  else if (seconds < 60) duration = `${seconds} seconds`;
  else if (minutes < 60) duration = `${minutes} minute${minutes > 1 ? "s" : ""}`;
  else if (hours < 24) duration = `${hours} hour${hours > 1 ? "s" : ""}`;
  else duration = `${days} day${days > 1 ? "s" : ""}`;

  return `${duration}`;
}

const handleWhatsapp = () => {
  console.log(customerProfile.value)
  window.open(`https://wa.me/${customerProfile.value.phone_number.replace(/\D/g, '')}`, '_blank')
}

const handleCall = () => {
  window.location.href = `tel:${customerProfile.value.phone_number}`
}

const handleEmail = () => {
  console.log(customerProfile.value.email)
  window.location.href = `mailto:${customerProfile.value.email}`
}

const copyToClipboard = (text: string, field: string) => {
  navigator.clipboard.writeText(text)
  notify(`${field} copied successfully`, "success");
}

const isPostingComment = ref(false)

const postComment = async () => {
  if (!comment.value.trim() || !rawSubmission.value) return
  isPostingComment.value = true
  try {
    const result = await useCustomerAPI().PostComment(rawSubmission.value._id, comment.value.trim())
    if (result.success) {
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const newActivity = (result as any).data?.data as CustomerActivity
      if (newActivity) {
        activities.value = [newActivity, ...activities.value]
      }
      comment.value = ''
    } else {
      notify("Failed to post comment", "error")
    }
  } finally {
    isPostingComment.value = false
  }
}
</script>

<template>
  <div>
    <!-- Header -->
    <MainDashboard current="Customers">
      <template #title>
        <div class="">
          <div class="max-w-7xl mx-auto">
            <div class="flex flex-wrap gap-4">
              <div class="">
                <BackLink
                  to="/customers"
                  label="Customers"
                />
                <h1
                  v-if="customerProfile"
                  class="text-xl sm:text-3xl font-bold text-[#1E212B]"
                >
                  {{ customerProfile?.customer_type === 'Business'
                    ? customerProfile?.company_name
                    : `${customerProfile?.first_name} ${customerProfile?.last_name}` }}
                </h1>
                <div class="flex flex-wrap items-center gap-2 mt-1 mb-4 text-sm text-gray-600">
                  <span>{{ customerProfile?.email
                    ? customerProfile?.email
                    : customerProfile?.address }}</span>
                  <span class="hidden sm:inline">•</span>
                  <span class="text-[#6F7177] text-xs">Customer for {{ formatCustomerDuration(customerProfile?.createdAt) }}</span>
                </div>
              </div>
              <div class="flex items-center gap-2 ml-auto">
                <GroBasicButton
                  color="secondary"
                  size="xs"
                  class="whitespace-nowrap"
                >
                  Create Payment
                </GroBasicButton>
                <GroBasicButton
                  color="secondary"
                  size="xs"
                  class="whitespace-nowrap"
                >
                  Create Invoice
                </GroBasicButton>
              </div>
            </div>
          </div>
        </div>
      </template>

      <template #body>
        <div class="max-w-7xl mx-auto">
          <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <!-- Left Sidebar -->
            <div class="lg:col-span-1 space-y-6">
              <!-- Quick Actions -->
              <div class="bg-white rounded-lg shadow-sm p-6">
                <div class="flex justify-around gap-4">
                  <button
                    v-if="customerProfile?.phone_number"
                    class="flex flex-col cursor-pointer items-center gap-2 hover:opacity-80 transition-opacity"
                    @click="handleWhatsapp"
                  >
                    <div class="w-12 h-12 rounded-2xl flex items-center justify-center contact-box">
                      <svg
                        class="w-6 h-6 text-green-600"
                        fill="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 005.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 00-3.48-8.413Z" />
                      </svg>
                    </div>
                    <span class="text-xs font-medium text-[#1E212B]">Whatsapp</span>
                  </button>
                  <button
                    v-if="customerProfile?.phone_number"
                    class="flex flex-col items-center cursor-pointer gap-2 hover:opacity-80 transition-opacity "
                    @click="handleCall"
                  >
                    <div class="w-12 h-12 rounded-full flex items-center justify-center contact-box">
                      <HugeiconsIcon
                        :icon="CallIcon"
                        class="my-auto"
                        :size="16"
                        fill="#1E212B"
                        :stroke-width="3"
                      />
                    </div>
                    <span class="text-xs font-medium text-[#1E212B]">Call</span>
                  </button>
                  <button
                    v-if="customerProfile?.email"
                    class="flex flex-col items-center cursor-pointer gap-2 hover:opacity-80 transition-opacity "
                    @click="handleEmail"
                  >
                    <div class="w-12 h-12 rounded-full  flex items-center justify-center contact-box">
                      <HugeiconsIcon
                        :icon="Mail01Icon"
                        class="my-auto"
                        :size="16"
                        color="#1E212B"
                        :stroke-width="3"
                      />
                    </div>
                    <span class="text-xs font-medium text-[#1E212B]">Send an Email</span>
                  </button>
                </div>
              </div>

              <!-- Customer Information -->
              <div
                v-if="customerProfile"
                class="bg-white rounded-lg shadow-sm p-6"
              >
                <div class="flex items-center justify-between mb-4">
                  <h2 class="text-lg font-bold text-[#1E212B]">
                    Customer information
                  </h2>
                  <div class="relative info-menu-container">
                    <button
                      class="p-1 hover:bg-gray-100 rounded transition-colors"
                      @click.stop="showInfoMenu = !showInfoMenu"
                    >
                      <svg
                        class="w-5 h-5 text-gray-500"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          :stroke-width="2"
                          d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"
                        />
                      </svg>
                    </button>
                    <Transition
                      enter-active-class="transition ease-out duration-150"
                      enter-from-class="opacity-0 scale-95"
                      enter-to-class="opacity-100 scale-100"
                      leave-active-class="transition ease-in duration-100"
                      leave-from-class="opacity-100 scale-100"
                      leave-to-class="opacity-0 scale-95"
                    >
                      <ul
                        v-show="showInfoMenu"
                        class="absolute right-0 mt-1 w-36 bg-white border border-slate-200 rounded-xl shadow-xl p-1 z-50 origin-top-right"
                      >
                        <li>
                          <button
                            class="w-full text-left px-3 py-2 text-sm text-slate-700 hover:bg-slate-50 rounded-lg"
                            @click="showEditModal = true; showInfoMenu = false"
                          >
                            Edit
                          </button>
                        </li>
                      </ul>
                    </Transition>
                  </div>
                </div>

                <div class="space-y-2">
                  <div
                    v-if="customerProfile?.email"
                    class="flex items-center justify-between group"
                  >
                    <a
                      href="#"
                      class="text-[#2176AE] font-medium text-sm"
                      @click.prevent="handleEmail"
                    >
                      {{ customerProfile?.email }}
                    </a>
                    <button
                      class=" rounded transition-all"
                      @click="copyToClipboard(customerProfile?.email, 'Email')"
                    >
                      <HugeiconsIcon
                        color="#B7B8BB"
                        fill="#B7B8BB"
                        :icon="Copy01Icon"
                        class="w-4 h-4 cursor-pointer"
                      />
                    </button>
                  </div>

                  <div
                    v-if="customerProfile?.phone_number"
                    class="flex items-center justify-between group"
                  >
                    <a
                      href="#"
                      class="text-[#2176AE] font-medium text-sm"
                      @click.prevent="handleCall"
                    >
                      {{ customerProfile?.phone_number }}
                    </a>
                    <button
                      class=" rounded transition-all"
                      @click="copyToClipboard(customerProfile?.phone_number, 'Phone number')"
                    >
                      <HugeiconsIcon
                        color="#B7B8BB"
                        fill="#B7B8BB"
                        :icon="Copy01Icon"
                        class="w-4 h-4 cursor-pointer"
                      />
                    </button>
                  </div>

                  <div
                    v-if="customerProfile"
                    class="pt-3"
                  >
                    <h2 class="text-md mb-2 font-bold text-[#1E212B]">
                      Address
                    </h2>
                    <p class="text-sm text-[#6F7177]">
                      {{ customerProfile.address }}
                    </p>
                    <p class="text-sm text-[#6F7177]">
                      {{ customerProfile.state }}
                    </p>
                    <p class="text-sm text-[#6F7177]">
                      {{ customerProfile.country }}
                    </p>
                  </div>

                  <div
                    v-if="customerProfile"
                    class="pt-3"
                  >
                    <h3 class="text-md mb-2 font-bold text-[#1E212B]">
                      Marketing
                    </h3>
                    <ul class="space-y-2 text-sm text-gray-600">
                      <li class="flex items-center gap-2">
                        <span class="w-1.5 h-1.5 rounded-full bg-gray-400" />
                        <span class="text-[#6F7177]">
                          {{ customerProfile?.subscribe_marketing_email?.includes('subscribe_marketing_email')
                            ? 'Email subscribed'
                            : 'Email not subscribed' }}
                        </span>
                      </li>
                      <li class="flex items-center gap-2">
                        <span class="w-1.5 h-1.5 rounded-full bg-gray-400" />
                        <span class="text-[#6F7177]">
                          {{ customerProfile?.subscribe_marketing_email?.includes('subscribe_sms')
                            ? 'SMS subscribed'
                            : 'SMS not subscribed' }}
                        </span>
                      </li>
                    </ul>
                  </div>

                  <div
                    v-if="extraFields.length"
                    class="pt-3"
                  >
                    <h3 class="text-md mb-2 font-bold text-[#1E212B]">
                      Extra fields
                    </h3>
                    <ul class="space-y-2">
                      <li
                        v-for="field in extraFields"
                        :key="field.slug"
                      >
                        <span class="text-xs font-bold text-[#1E212B]">{{ field.label }}</span>
                        <p class="text-sm text-gray-900">
                          {{ field.value }}
                        </p>
                      </li>
                    </ul>
                  </div>
                </div>
              </div>

              <!-- Tags -->
              <div class="bg-white rounded-lg shadow-sm p-6">
                <div class="flex items-center justify-between mb-4">
                  <h3 class="text-md mb-2 font-bold text-[#1E212B]">
                    Tags
                  </h3>
                  <button
                    class="p-1 hover:bg-gray-100 rounded transition-colors"
                    @click="showTagsEdit = !showTagsEdit; editTagIds = [...(customerProfile?.customer_tags ?? [])]"
                  >
                    <svg
                      class="w-5 h-5 text-gray-500"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        :stroke-width="2"
                        d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"
                      />
                    </svg>
                  </button>
                </div>

                <div v-if="showTagsEdit">
                  <GroBasicTagSelect v-model="editTagIds" />
                  <div class="flex gap-2 mt-3">
                    <GroBasicButton
                      color="primary"
                      size="xs"
                      :disabled="isSavingTags"
                      @click="saveTagsUpdate"
                    >
                      Save
                    </GroBasicButton>
                    <GroBasicButton
                      color="secondary"
                      size="xs"
                      @click="showTagsEdit = false"
                    >
                      Cancel
                    </GroBasicButton>
                  </div>
                </div>
                <div
                  v-else
                  class="flex flex-wrap gap-2"
                >
                  <GroReadOnlyTags :tagids="customerProfile?.customer_tags" />
                </div>
              </div>
            </div>

            <div class="lg:col-span-2 space-y-6">
              <div class="bg-white rounded-lg shadow-sm p-6">
                <h2 class="text-lg font-bold text-[#1E212B]">
                  Orders
                </h2>
                <p class="text-[#4B4D55] text-sm mt-2 mb-4">
                  This customer hasn't placed any orders yet
                </p>
                <div>
                  <NuxtLink
                    to="/invoices/create-new-invoice"
                  >
                    <GroBasicButton
                      color="tertiary"
                      size="xs"
                      class="whitespace-nowrap w-max"
                    >
                      Create Invoice
                    </GroBasicButton>
                  </NuxtLink>
                </div>
              </div>

              <!-- Timeline -->
              <div class="rounded-lg mt-10">
                <h2 class="text-lg font-bold text-[#1E212B] mb-4">
                  Timeline
                </h2>

                <!-- Full-width editor -->
                <CommentEditor v-model="comment">
                  <template #actions>
                    <button
                      type="button"
                      class="px-4 py-1.5 bg-[#1E212B] hover:bg-[#2d3142] text-white rounded-lg text-sm font-medium transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
                      :disabled="!comment.trim() || isPostingComment"
                      @click="postComment"
                    >
                      {{ isPostingComment ? 'Posting...' : 'Post' }}
                    </button>
                  </template>
                </CommentEditor>
                <p class="text-xs text-right text-[#6F7177] mt-1">
                  Only you can see this comment
                </p>

                <!-- Activity entries — rail starts immediately at editor bottom -->
                <div
                  v-if="activities.length > 0"
                  class="relative mt-1"
                >
                  <!-- Continuous vertical rail: left-3.5 = 14px = center of w-7 icon column -->
                  <div class="absolute left-3.5 top-0 bottom-0 w-px bg-gray-200 -translate-x-1/2" />

                  <div
                    v-for="activity in activities"
                    :key="activity.id"
                    class="flex gap-3"
                  >
                    <!-- Icon node sits on top of the rail via z-10; ring-white hides the rail behind it -->
                    <div class="w-7 shrink-0 flex justify-center pt-3">
                      <div
                        class="w-7 h-7 rounded-full flex items-center justify-center relative z-10 ring-2 ring-white shrink-0"
                        :class="activity.activityType === 'comment' ? 'bg-blue-100' : 'bg-gray-100'"
                      >
                        <svg
                          v-if="activity.activityType === 'comment'"
                          class="w-3.5 h-3.5 text-blue-500"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
                          />
                        </svg>
                        <svg
                          v-else
                          class="w-3.5 h-3.5 text-gray-500"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                          />
                        </svg>
                      </div>
                    </div>

                    <!-- Content -->
                    <div class="flex-1 pt-2 pb-5">
                      <div class="flex items-center gap-2 mb-0.5">
                        <span
                          class="text-xs font-medium"
                          :class="activity.activityType === 'comment' ? 'text-blue-600' : 'text-gray-500'"
                        >
                          {{ activity.activityType === 'comment' ? 'Comment' : 'Activity' }}
                        </span>
                        <span class="text-xs text-[#6F7177]">· {{ formatCustomerDuration(activity.createdAt) }} ago</span>
                      </div>
                      <!-- eslint-disable-next-line vue/no-v-html -->
                      <p
                        v-if="activity.activityType === 'comment'"
                        class="text-sm text-gray-800 bg-gray-50 rounded-lg px-3 py-2 comment-content"
                        v-html="activity.description"
                      />
                      <p
                        v-else
                        class="text-sm text-gray-700"
                      >
                        {{ activity.description }}
                      </p>
                    </div>
                  </div>
                </div>

                <p
                  v-else
                  class="text-sm text-gray-400 pl-10 py-2"
                >
                  No activity yet
                </p>
              </div>
            </div>
          </div>
        </div>
      </template>
    </MainDashboard>

    <!-- Edit Customer Modal -->
    <CreateCustomerModal
      v-if="rawSubmission"
      v-model="showEditModal"
      :customer="{
        submissionId: rawSubmission._id,
        formId: rawSubmission.formID,
        answers: customerProfile,
      }"
      @customer-updated="fetchCustomer"
    />
  </div>
</template>

<style scoped>
.contact-box {
  box-shadow:
    0 1px 0 0 #E3E3E3 inset,
    1px 0 0 0 #E3E3E3 inset,
    -1px 0 0 0 #E3E3E3 inset,
    0 -1px 0 0 #B5B5B5 inset;
}

/* Rich text comment rendering */
.comment-content :deep(p) { margin: 0; }
.comment-content :deep(p + p) { margin-top: 0.4em; }
.comment-content :deep(ul) { list-style: disc; padding-left: 1.25rem; }
.comment-content :deep(ol) { list-style: decimal; padding-left: 1.25rem; }
.comment-content :deep(strong) { font-weight: 600; }
.comment-content :deep(em) { font-style: italic; }
.comment-content :deep(s) { text-decoration: line-through; }
.comment-content :deep(a) { color: #2176AE; text-decoration: underline; }
.comment-content :deep(blockquote) {
  border-left: 3px solid #e5e7eb;
  padding-left: 0.75rem;
  color: #6b7280;
  font-style: italic;
  margin: 0.4rem 0;
}
</style>

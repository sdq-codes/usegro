<script setup lang="ts">
import {onMounted, ref} from 'vue'
import MainDashboard from "@/components/dashboard/main-dashboard.vue";
import {useCustomerAPI} from "@/composables/api/customer/customer";
import {useRoute} from "nuxt/app";
import GroBasicButton from "@/components/buttons/GroBasicButton.vue";
import {Copy01Icon, InformationCircleIcon, LessThanIcon} from "@hugeicons/core-free-icons";
import {HugeiconsIcon} from "@hugeicons/vue";
import {notify} from "@/composables/helpers/notification/notification";
import GroReadOnlyTags from "@/components/tags/GroReadOnlyTags.vue";

const customerProfile = ref()

const comment = ref('')
const showLabelsEdit = ref(false)

const route = useRoute()

onMounted(async () => {
  await fetchCustomer()
})

const fetchCustomer = async () => {
  const createdCustomer = await useCustomerAPI().FetchCustomer(String(route.params?.submission ?? ""), String(route.params?.id ?? ""))
  if (!createdCustomer.success) {
    return;
  }

  customerProfile.value = {
      ...createdCustomer.data?.data?.Answers,
      customerId: createdCustomer.data?.data?.SubmissionID,
      createdAt: createdCustomer.data?.data?.CreatedAt,
      formId: createdCustomer.data?.data?.PK?.split("#")[1].trim(),
    };

  console.log(customerProfile.value);
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
  window.open(`https://wa.me/${customer.value.phone.replace(/\D/g, '')}`, '_blank')
}

const handleCall = () => {
  window.location.href = `tel:${customer.value.phone}`
}

const handleEmail = () => {
  window.location.href = `mailto:${customer.value.email}`
}

const copyToClipboard = (text: string, field: string) => {
  navigator.clipboard.writeText(text)
  notify(`${field} copied successfully`, "success");
}

const postComment = () => {
  if (comment.value.trim()) {
    comment.value = ''
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
                <div class="flex mb-3">
                  <HugeiconsIcon
                    :icon="LessThanIcon"
                    size="13"
                    class="my-auto mr-1"
                    color="currentColor"
                  />
                  <NuxtLink
                    to="/customers"
                    class="text-[#2176AE] text-sm font-medium inline-flex items-center border-b border-[#2176AE]"
                  >
                    Customers
                  </NuxtLink>
                </div>
                <h1
                  v-if="customerProfile"
                  class="text-2xl sm:text-3xl font-bold text-[#1E212B]"
                >
                  {{ customerProfile?.customer_type === 'Business'
                    ? customerProfile?.company_name
                    : `${customerProfile?.first_name} ${customerProfile?.last_name}` }}
                </h1>
                <div class="flex flex-wrap items-center gap-2 mt-1 text-sm text-gray-600">
                  <span>{{ customerProfile?.email
                    ? customerProfile?.email
                    : customerProfile?.address }}</span>
                  <span class="hidden sm:inline">•</span>
                  <span class="text-[#6F7177]">Customer for {{ formatCustomerDuration(customerProfile?.createdAt) }}</span>
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
                    class="flex flex-col items-center gap-2 hover:opacity-80 transition-opacity"
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
                    class="flex flex-col items-center gap-2 hover:opacity-80 transition-opacity "
                    @click="handleCall"
                  >
                    <div class="w-12 h-12 rounded-full flex items-center justify-center contact-box">
                      <svg
                        class="w-6 h-6 text-gray-700"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          :stroke-width="2"
                          d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"
                        />
                      </svg>
                    </div>
                    <span class="text-xs font-medium text-[#1E212B]">Call</span>
                  </button>
                  <button
                    class="flex flex-col items-center gap-2 hover:opacity-80 transition-opacity "
                    @click="handleEmail"
                  >
                    <div class="w-12 h-12 rounded-full  flex items-center justify-center contact-box">
                      <svg
                        class="w-6 h-6 text-gray-700"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          :stroke-width="2"
                          d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"
                        />
                      </svg>
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
                  <button class="p-1 hover:bg-gray-100 rounded transition-colors">
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
                </div>

                <div class="space-y-3">
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
                    <p class="text-sm text-gray-600">
                      {{ customerProfile.address }}
                    </p>
                    <p class="text-sm text-gray-600">
                      {{ customerProfile.state }}
                    </p>
                    <p class="text-sm text-gray-600">
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
                        <span>
                          {{ customerProfile?.subscribe_marketing_email?.includes('subscribe_marketing_email')
                            ? 'Email subscribed'
                            : 'Email not subscribed' }}
                        </span>
                      </li>
                      <li class="flex items-center gap-2">
                        <span class="w-1.5 h-1.5 rounded-full bg-gray-400" />
                        <span>
                          {{ customerProfile?.subscribe_marketing_email?.includes('subscribe_sms')
                            ? 'SMS subscribed'
                            : 'SMS not subscribed' }}
                        </span>
                      </li>
                    </ul>
                  </div>
                </div>
              </div>

              <!-- Labels -->
              <div class="bg-white rounded-lg shadow-sm p-6">
                <div class="flex items-center justify-between mb-4">
                  <h3 class="text-md mb-2 font-bold text-[#1E212B]">
                    Labels
                  </h3>
                  <button
                    class="p-1 hover:bg-gray-100 rounded transition-colors"
                    @click="showLabelsEdit = !showLabelsEdit"
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
                <div class="flex flex-wrap gap-2">
                  <GroReadOnlyTags :tagids="customerProfile?.customer_tags" />
                </div>
              </div>
            </div>

            <div class="lg:col-span-2 space-y-6">
              <div class="bg-white rounded-lg shadow-sm p-6">
                <h2 class="text-lg font-bold text-[#1E212B]">
                  Orders
                </h2>
                <p class="text-[#4B4D55] text-sm mb-4">
                  This customer hasn’t placed any orders yet
                </p>
                <div>
                  <GroBasicButton
                    color="tertiary"
                    size="xs"
                    class="whitespace-nowrap w-max"
                  >
                    Create Invoice
                  </GroBasicButton>
                </div>
              </div>

              <!-- Timeline -->
              <div class=" rounded-lg mt-10">
                <h2 class="text-lg font-bold text-[#1E212B]">
                  Timeline
                </h2>
                <div class="timeline-shadow rounded-lg">
                  <div class="my-2 p-4 bg-white">
                    <textarea
                      v-model="comment"
                      placeholder="Leave a comment..."
                      class="w-full px-4 py-3 rounded-lg focus:ring-0 focus:border-transparent active:border-0 active:ring-0 resize-none text-sm"
                      rows="3"
                    />
                  </div>
                  <div>
                    <div class="flex bg-[f5f5f5] px-4 items-center justify-between">
                      <div class="flex items-center gap-2 gap-x-3 bg-[">
                        <button class="hover:bg-gray-100 rounded transition-colors">
                          <svg
                            class="w-6 h-6 text-gray-500"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                          >
                            <path
                              stroke-linecap="round"
                              stroke-linejoin="round"
                              :stroke-width="2"
                              d="M14.828 14.828a4 4 0 01-5.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                            />
                          </svg>
                        </button>
                        <button class="hover:bg-gray-100 rounded transition-colors">
                          <svg
                            class="w-6 h-6 text-gray-500"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                          >
                            <path
                              stroke-linecap="round"
                              stroke-linejoin="round"
                              :stroke-width="2"
                              d="M7 20l4-16m2 16l4-16M6 9h14M4 15h14"
                            />
                          </svg>
                        </button>
                        <button class="hover:bg-gray-100 rounded transition-colors">
                          <svg
                            class="w-6 h-6 text-gray-500"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                          >
                            <path
                              stroke-linecap="round"
                              stroke-linejoin="round"
                              :stroke-width="2"
                              d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"
                            />
                          </svg>
                        </button>
                      </div>
                      <button
                        class="px-4 py-2 my-3 bg-white border border-gray-300 hover:bg-gray-50 text-gray-700 rounded-lg text-sm font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                        :disabled="!comment.trim()"
                        @click="postComment"
                      >
                        Post
                      </button>
                    </div>
                  </div>
                </div>

                <div class="space-y-6">
                  <div class="flex gap-4 px-5">
                    <div class="flex flex-col items-center">
                      <div class="w-px h-full bg-gray-200" />
                      <div class="w-2 h-2 rounded-full bg-[#1E212B] mt-1" />
                    </div>
                    <div class="flex-1 pb-6 mt-7">
                      <p class="text-sm text-#6F7177] mb-1">
                        {{ formatCustomerDuration(customerProfile?.createdAt) }}
                      </p>
                      <p class="text-sm text-gray-900">
                        You created this customer
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>
    </MainDashboard>
  </div>
</template>

<style scoped>
/* Additional custom styles if needed */

.contact-box {
  box-shadow:
    0px 1px 0px 0px #E3E3E3 inset,
    1px 0px 0px 0px #E3E3E3 inset,
    -1px 0px 0px 0px #E3E3E3 inset,
    0px -1px 0px 0px #B5B5B5 inset;
}

.timeline-shadow {
  box-shadow:
    0px 1px 3px 0px #0000000A,
    0px 1px 0px 0px #0000000A;
}
</style>

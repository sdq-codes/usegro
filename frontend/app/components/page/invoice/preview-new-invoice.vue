<script setup lang="ts">
import { ref, computed } from 'vue';

interface InvoiceItem {
  name: string;
  quantity: number;
  rate: number;
  amount: number;
}

interface InvoiceData {
  businessName: string;
  website: string;
  email: string;
  phone: string;
  businessAddress: string;
  city: string;
  state: string;
  country: string;
  postalCode: string;

  // These fields now sync with the Form Models
  billedTo: any[] | null;
  selectedProducts: any[] | null;
  invoiceNumber: string;
  termsAndConditions: string;
  subject: string; // Used for "Notes"

  taxRate: number;
  genDate?: string;
  dueDate?: string;
}

const props = withDefaults(defineProps<{
  invoiceData: InvoiceData;
}>(), {
  invoiceData: () => ({
    businessName: 'UseGro.inc',
    website: 'www.website.com',
    email: 'hello@email.com',
    phone: '+91 00000 00000',
    businessAddress: 'Business address',
    city: 'City',
    state: 'State',
    country: 'US',
    postalCode: '00000',
    billedTo: [],
    selectedProducts: [],
    invoiceNumber: '',
    termsAndConditions: '',
    subject: '',
    taxRate: 0
  })
});

const activeTab = ref<'pdf' | 'mobile' | 'email'>('pdf');

// 1. Format Customer Names for the "Billed To" section
const billedToDisplay = computed(() => {
  const customers = props.invoiceData.billedTo;
  if (!customers || customers.length === 0) return '';
  // Adjust 'name' property based on your actual Customer object structure
  return customers.map((c: any) => c.name || c.label || 'New Customer').join(', ');
});

// 2. Map Selected Products to the Table Rows
const displayItems = computed<InvoiceItem[]>(() => {
  const products = props.invoiceData.selectedProducts;
  if (!products || products.length === 0) {
    return [{ name: 'Item Name', quantity: 0, rate: 0, amount: 0 }];
  }
  return products.map((p: any) => ({
    name: p.name || p.label || 'Product/Service',
    quantity: p.quantity || 1,
    rate: p.price || p.rate || 0,
    amount: (p.quantity || 1) * (p.price || p.rate || 0)
  }));
});

const subtotal = computed(() => {
  return displayItems.value.reduce((sum, item) => sum + item.amount, 0);
});

const taxAmount = computed(() => {
  return (subtotal.value * (props.invoiceData.taxRate || 0)) / 100;
});

const total = computed(() => {
  return subtotal.value + taxAmount.value;
});

const formatCurrency = (amount: number) => {
  return `$${amount.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;
};
</script>

<template>
  <div class="w-1/2 overflow-y-auto rounded-r-2xl border-y border-r bg-[#EDF4F9] border-[#EDEDEE]">
    <h2 class="text-xl text-[#1E212B] mb-6 pt-8 px-16">Preview</h2>

    <div class="grid grid-cols-3 mb-8 px-20">
      <button
        v-for="tab in (['pdf', 'mobile', 'email'] as const)"
        :key="tab"
        @click="activeTab = tab"
        :class="[
          'pb-3 font-medium text-sm transition-colors relative cursor-pointer capitalize',
          activeTab === tab ? 'text-[#E87117]' : 'text-gray-500 hover:text-gray-700 border-b-2 border-[#EDEDEE]'
        ]"
      >
        {{ tab }} preview
        <span v-if="activeTab === tab" class="absolute bottom-0 left-0 right-0 h-0.5 bg-[#E87117]" />
      </button>
    </div>

    <div id="preview-new-invoice-container" class="rounded-xl px-16 pb-10">
      <div class="bg-white shadow-[0px_35.16px_70.32px_0px_#4E4D593D] rounded-lg px-4 py-8">
        <div class="flex justify-between items-start mb-5">
          <div class="flex gap-x-4 px-6">
            <img src="https://res.cloudinary.com/sdq121/image/upload/v1764981178/ccrh2qinydoihgbswden.svg" class="w-16 h-16">
            <div class="space-y-0.5">
              <h3 class="text-xl font-semibold text-[#E87117]">{{ invoiceData.businessName }}</h3>
              <p class="text-xs font-light text-[#5E6470]">{{ invoiceData.website }}</p>
              <p class="text-xs font-light text-[#5E6470]">{{ invoiceData.email }}</p>
              <p class="text-xs font-light text-[#5E6470]">{{ invoiceData.phone }}</p>
            </div>
          </div>
          <div class="text-right ml-auto mt-auto space-y-0.5">
            <p class="text-xs font-light text-[#5E6470]">{{ invoiceData.businessAddress }}</p>
            <p class="text-xs font-light text-[#5E6470]">
              {{ invoiceData.city }}, {{ invoiceData.state }}, {{ invoiceData.country }} - {{ invoiceData.postalCode }}
            </p>
          </div>
        </div>

        <div class="rounded-2xl border-[0.44px] border-[#D7DAE0] bg-[#FAFAFA]">
          <div class="rounded-lg p-6 mb-2">
            <div class="grid grid-cols-3 gap-6 mb-6">
              <div>
                <p class="text-xs text-gray-400 mb-1">Billed to</p>
                <div v-if="!billedToDisplay" class="w-[80%] h-3 bg-gray-200 rounded animate-pulse" />
                <p class="text-sm font-medium text-[#1E212B]">{{ billedToDisplay }}</p>
              </div>
              <div>
                <p class="text-xs text-gray-400 mb-1">Invoice number</p>
                <div v-if="!invoiceData.invoiceNumber" class="w-[60%] h-3 bg-gray-200 rounded animate-pulse" />
                <p class="text-sm font-semibold text-[#1E212B]">{{ invoiceData.invoiceNumber }}</p>
              </div>
              <div class="text-right">
                <p class="text-xs text-gray-500 mb-1">Invoice of (USD)</p>
                <p class="text-xl font-bold text-[#E87117]">{{ formatCurrency(total) }}</p>
              </div>
            </div>

            <div class="grid grid-cols-3 gap-6">
              <div>
                <p class="text-xs text-gray-500 mb-1">Generated date</p>
                <p class="text-sm font-semibold text-[#1E212B]">{{ invoiceData.genDate || '---' }}</p>
              </div>
              <div>
                <p class="text-xs text-gray-500 mb-1">Due date</p>
                <p class="text-sm font-semibold text-[#1E212B]">{{ invoiceData.dueDate || '---' }}</p>
              </div>
              <div class="text-right">
                <p class="text-xs text-gray-500 mb-1">Notes</p>
                <p class="text-sm font-semibold text-[#1E212B] truncate">{{ invoiceData.subject || '---' }}</p>
              </div>
            </div>
          </div>

          <div class="rounded-lg overflow-hidden mb-6 px-2">
            <table class="w-full">
              <thead>
              <tr class="border-y border-[#D7DAE0]">
                <th class="text-left text-xs font-semibold text-gray-600 px-6 py-3 uppercase">Item Detail</th>
                <th class="text-center text-xs font-semibold text-gray-600 px-6 py-3 uppercase">QTY</th>
                <th class="text-center text-xs font-semibold text-gray-600 px-6 py-3 uppercase">Rate</th>
                <th class="text-right text-xs font-semibold text-gray-600 px-6 py-3 uppercase">Amount</th>
              </tr>
              </thead>
              <tbody>
              <tr v-for="(item, index) in displayItems" :key="index" class="border-b border-gray-100 last:border-b-0">
                <td class="px-6 py-4 text-sm font-medium text-[#1E212B]">{{ item.name }}</td>
                <td class="px-6 py-4 text-sm text-center text-gray-600">{{ item.quantity || '-' }}</td>
                <td class="px-6 py-4 text-sm text-center text-gray-600">{{ formatCurrency(item.rate) }}</td>
                <td class="px-6 py-4 text-sm text-right font-medium text-[#1E212B]">{{ formatCurrency(item.amount) }}</td>
              </tr>
              </tbody>
            </table>

            <div class="px-6 py-4 space-y-3 bg-white/50">
              <div class="flex justify-between items-center">
                <span class="text-sm font-medium text-gray-600">Subtotal</span>
                <span class="text-sm font-semibold text-[#1E212B]">{{ formatCurrency(subtotal) }}</span>
              </div>
              <div class="flex justify-between items-center">
                <span class="text-sm font-medium text-gray-600">Tax ({{ invoiceData.taxRate }}%)</span>
                <span class="text-sm font-semibold text-[#1E212B]">{{ formatCurrency(taxAmount) }}</span>
              </div>
              <div class="flex justify-between items-center pt-3 border-t border-gray-200">
                <span class="text-base font-bold text-[#1E212B]">Total</span>
                <span class="text-base font-bold text-[#1E212B]">{{ formatCurrency(total) }}</span>
              </div>
            </div>
          </div>
        </div>

        <div class="mt-8">
          <h4 class="text-sm font-bold text-[#1E212B] mb-2">Terms & Conditions</h4>
          <p class="text-sm text-gray-600 leading-relaxed">
            {{ invoiceData.termsAndConditions || 'No specific terms provided.' }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

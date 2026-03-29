<script setup lang="ts">
import { type InvoiceTemplate, INVOICE_TEMPLATES } from '@/constants/invoice/templates'

const model = defineModel<boolean>({ default: false })
const selectedTemplate = defineModel<InvoiceTemplate>('selectedTemplate', { default: 'classic' })
</script>

<template>
  <Teleport to="body">
    <div
      v-if="model"
      class="fixed inset-0 bg-black/50 backdrop-blur-[2px] z-50 flex items-center justify-center"
      @click.self="model = false"
    >
      <div class="bg-white rounded-2xl shadow-2xl w-[90vw] max-w-xl">
        <!-- Header -->
        <div class="flex items-center justify-between px-6 py-5 border-b border-[#EDEDEE]">
          <h2 class="text-base font-semibold text-[#1E212B]">Invoice Template Settings</h2>
          <button
            class="text-gray-400 hover:text-gray-700 transition-colors cursor-pointer"
            @click="model = false"
          >
            <svg class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M18 6L6 18M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- Template grid -->
        <div class="px-6 py-6">
          <p class="text-xs text-gray-400 mb-5">
            Choose a template to customise how your invoice looks to customers.
          </p>
          <div class="grid grid-cols-3 gap-4">
            <button
              v-for="tmpl in INVOICE_TEMPLATES"
              :key="tmpl.id"
              class="group relative rounded-xl border-2 transition-all cursor-pointer text-left focus:outline-none"
              :class="selectedTemplate === tmpl.id
                ? 'border-[#E87117] shadow-[0_0_0_3px_#E8711720]'
                : 'border-[#EDEDEE] hover:border-gray-300'"
              @click="selectedTemplate = tmpl.id"
            >
              <!-- Thumbnail -->
              <div class="rounded-t-[10px] overflow-hidden h-36 bg-[#F8F9FA]">
                <!-- Classic thumbnail -->
                <div v-if="tmpl.id === 'classic'" class="p-2 h-full flex flex-col gap-1.5">
                  <div class="flex items-center gap-1.5">
                    <div class="w-5 h-5 rounded bg-[#E87117] opacity-80" />
                    <div class="flex-1 space-y-0.5">
                      <div class="h-1.5 w-16 bg-[#E87117] rounded-full opacity-80" />
                      <div class="h-1 w-10 bg-gray-300 rounded-full" />
                    </div>
                  </div>
                  <div class="flex-1 rounded-lg border border-[#D7DAE0] bg-white p-1.5 space-y-1">
                    <div class="grid grid-cols-3 gap-1">
                      <div class="h-1 bg-gray-200 rounded-full" />
                      <div class="h-1 bg-gray-200 rounded-full" />
                      <div class="h-1.5 bg-[#E87117] rounded-full opacity-70" />
                    </div>
                    <div class="h-px bg-gray-200" />
                    <div class="space-y-1">
                      <div class="flex justify-between">
                        <div class="h-1 w-14 bg-gray-200 rounded-full" />
                        <div class="h-1 w-8 bg-gray-200 rounded-full" />
                      </div>
                      <div class="flex justify-between">
                        <div class="h-1 w-10 bg-gray-200 rounded-full" />
                        <div class="h-1 w-8 bg-gray-200 rounded-full" />
                      </div>
                    </div>
                    <div class="h-px bg-gray-200" />
                    <div class="flex justify-end">
                      <div class="h-1.5 w-12 bg-[#E87117] rounded-full opacity-70" />
                    </div>
                  </div>
                </div>

                <!-- Modern thumbnail -->
                <div v-else-if="tmpl.id === 'modern'" class="h-full flex flex-col">
                  <div class="bg-[#1E212B] px-2.5 py-2 flex items-center justify-between">
                    <div class="flex items-center gap-1.5">
                      <div class="w-4 h-4 rounded bg-gray-600" />
                      <div class="h-1.5 w-10 bg-gray-500 rounded-full" />
                    </div>
                    <div class="text-right">
                      <div class="h-2 w-10 bg-[#E87117] rounded-full" />
                      <div class="h-1 w-6 bg-gray-600 rounded-full mt-0.5 ml-auto" />
                    </div>
                  </div>
                  <div class="h-0.5 bg-[#E87117]" />
                  <div class="flex-1 bg-white p-2 space-y-1.5">
                    <div class="flex justify-between">
                      <div class="h-1 w-12 bg-gray-200 rounded-full" />
                      <div class="h-2 w-10 bg-[#E87117] rounded-full opacity-80" />
                    </div>
                    <div class="h-px bg-gray-100" />
                    <div class="bg-[#1E212B] rounded h-5 px-1.5 flex items-center">
                      <div class="h-1 w-8 bg-gray-500 rounded-full" />
                    </div>
                    <div class="space-y-1">
                      <div class="flex justify-between items-center bg-gray-50 px-1 py-0.5 rounded">
                        <div class="h-1 w-12 bg-gray-300 rounded-full" />
                        <div class="h-1 w-8 bg-gray-300 rounded-full" />
                      </div>
                      <div class="flex justify-between items-center px-1 py-0.5 rounded">
                        <div class="h-1 w-10 bg-gray-200 rounded-full" />
                        <div class="h-1 w-8 bg-gray-200 rounded-full" />
                      </div>
                    </div>
                    <div class="flex justify-between items-center border-t border-gray-100 pt-1">
                      <div class="h-1.5 w-6 bg-[#1E212B] rounded-full opacity-40" />
                      <div class="h-1.5 w-10 bg-[#E87117] rounded-full opacity-80" />
                    </div>
                  </div>
                  <div class="bg-[#1E212B] py-1 px-2">
                    <div class="h-1 w-10 bg-gray-600 rounded-full" />
                  </div>
                </div>

                <!-- Minimal thumbnail -->
                <div v-else-if="tmpl.id === 'minimal'" class="p-3 h-full flex flex-col">
                  <div class="flex justify-between mb-2">
                    <div>
                      <div class="h-2.5 w-14 bg-[#1E212B] rounded-full font-black" />
                      <div class="h-1 w-8 bg-[#2176AE] rounded-full mt-1 opacity-70" />
                      <div class="h-1 w-10 bg-gray-200 rounded-full mt-0.5" />
                    </div>
                    <div class="text-right">
                      <div class="h-1.5 w-12 bg-[#1E212B] rounded-full opacity-60" />
                      <div class="h-1 w-8 bg-gray-200 rounded-full mt-0.5 ml-auto" />
                    </div>
                  </div>
                  <div class="h-px bg-[#1E212B]" />
                  <div class="flex-1 pt-2 space-y-1">
                    <div class="flex justify-between pb-1">
                      <div class="h-1 w-12 bg-gray-200 rounded-full" />
                      <div class="h-1 w-8 bg-gray-200 rounded-full" />
                    </div>
                    <div class="h-px bg-gray-100" />
                    <div class="flex justify-between">
                      <div class="h-1 w-14 bg-gray-200 rounded-full" />
                      <div class="h-1 w-8 bg-gray-200 rounded-full" />
                    </div>
                    <div class="flex justify-between">
                      <div class="h-1 w-10 bg-gray-200 rounded-full" />
                      <div class="h-1 w-8 bg-gray-200 rounded-full" />
                    </div>
                  </div>
                  <div class="h-px bg-gray-100" />
                  <div class="flex justify-end mt-1.5">
                    <div class="h-2 w-12 bg-[#2176AE] rounded-full opacity-80" />
                  </div>
                </div>
              </div>

              <!-- Card footer -->
              <div class="px-3 py-2.5">
                <div class="flex items-center justify-between">
                  <p class="text-xs font-semibold text-[#1E212B]">{{ tmpl.label }}</p>
                  <div
                    class="w-3.5 h-3.5 rounded-full border-2 flex items-center justify-center transition-colors"
                    :class="selectedTemplate === tmpl.id ? 'border-[#E87117] bg-[#E87117]' : 'border-gray-300'"
                  >
                    <svg v-if="selectedTemplate === tmpl.id" class="w-2 h-2 text-white" viewBox="0 0 12 12" fill="currentColor">
                      <path d="M10 3L5 8.5 2 5.5" stroke="currentColor" stroke-width="1.5" fill="none" stroke-linecap="round" stroke-linejoin="round" />
                    </svg>
                  </div>
                </div>
                <p class="text-[10px] text-gray-400 mt-0.5 leading-relaxed">{{ tmpl.description }}</p>
              </div>
            </button>
          </div>
        </div>

        <!-- Footer -->
        <div class="flex justify-end gap-3 px-6 py-4 border-t border-[#EDEDEE]">
          <button
            class="px-4 py-2 text-sm text-gray-600 hover:text-gray-900 transition-colors cursor-pointer"
            @click="model = false"
          >
            Cancel
          </button>
          <button
            class="px-5 py-2 bg-[#E87117] hover:bg-[#d06810] text-white text-sm font-semibold rounded-xl transition-colors cursor-pointer"
            @click="model = false"
          >
            Apply Template
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

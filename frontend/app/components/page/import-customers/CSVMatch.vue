<script setup lang="ts">
import {defineModel, ref, watch} from 'vue';
import Papa from 'papaparse';
import {
  CheckmarkCircle01Icon,
  GreaterThanIcon,
  UnavailableFreeIcons, Tick01Icon
} from "@hugeicons/core-free-icons";
import { HugeiconsIcon } from "@hugeicons/vue";
import ImportCustomerTemplate from "@/components/page/import-customers/ImportCustomerTemplate.vue";
import GroBasicButton from "@/components/buttons/GroBasicButton.vue";
import type CSVData from "@/composables/helpers/types/csv";

interface Props {
  file: File | null;
}

interface FormField {
  slug: string;
  label: string;
  fieldTypeName: string;
  required: boolean;
}

interface ColumnValidation {
  valid: boolean;
  coverage: number;
}

const props = defineProps<Props>();
const stepModel = defineModel<number>("stepModel");
const csvData = defineModel<CSVData | null>("csvDataModel");

const emit = defineEmits<{
  (e: "csvDataChanged", value: CSVData | null): void;
}>();

watch(csvData, (newVal) => {
  emit("csvDataChanged", newVal);
});

// const csvData = ref<CSVData | null>(null);
const mappings = ref<Record<string, string>>({});
const confirmedColumns = ref<Set<string>>(new Set());
const ignoredColumns = ref<Set<string>>(new Set());


// Form fields from your JSON
const formFields: FormField[] = [
  { slug: 'customer_type', label: 'Customer Type', fieldTypeName: 'Radio Button', required: true },
  { slug: 'first_name', label: 'First name', fieldTypeName: 'Short Text', required: false },
  { slug: 'last_name', label: 'Last name', fieldTypeName: 'Short Text', required: false },
  { slug: 'company_name', label: 'Company Name', fieldTypeName: 'Short Text', required: false },
  { slug: 'email', label: 'Email', fieldTypeName: 'Email', required: false },
  { slug: 'phone_number', label: 'Phone number', fieldTypeName: 'Phone Number', required: false },
  { slug: 'subscribe_marketing_email', label: 'Subscribe Marketing Email, SMS', fieldTypeName: 'Checkbox', required: false },
  { slug: 'country', label: 'Country', fieldTypeName: 'Country', required: false },
  { slug: 'state', label: 'State/Province', fieldTypeName: 'State', required: false },
  { slug: 'address', label: 'Street Address', fieldTypeName: 'Short Text', required: false },
  { slug: 'city', label: 'City', fieldTypeName: 'City', required: false },
  { slug: 'postal_code', label: 'Zip/Postal Code', fieldTypeName: 'Short Text', required: false },
];

const validateColumn = (csvColumn: string): ColumnValidation => {
  const mapping = mappings.value[csvColumn];
  if (!mapping || !csvData.value) return { valid: true, coverage: 0 };

  const values = csvData.value.data.map(row => row[csvColumn]);
  const nonEmpty = values.filter(v => v !== null && v !== undefined && v !== '');
  const coverage = Math.round((nonEmpty.length / values.length) * 100);

  return { valid: true, coverage };
};

// Check if all mapped columns are confirmed or ignored
const canProceed = () => {
  if (!csvData.value) return false;

  return csvData.value.meta.fields.every(col => {
    // Column must be either confirmed or ignored
    return isColumnConfirmed(col) || isColumnIgnored(col);
  });
};

const handleMappingChange = (csvColumn: string, formField: string) => {
  if (formField === 'none') {
    const updated = { ...mappings.value };
    delete updated[csvColumn];
    mappings.value = updated;
  } else {
    mappings.value[csvColumn] = formField;
    // Remove from ignored if mapping is changed
    ignoredColumns.value.delete(csvColumn);
  }
};

const confirmMapping = (csvColumn: string) => {
  confirmedColumns.value.add(csvColumn);
};

const ignoreColumn = (csvColumn: string) => {
  ignoredColumns.value.add(csvColumn);
  const updated = { ...mappings.value };
  delete updated[csvColumn];
  mappings.value = updated;
  confirmedColumns.value.delete(csvColumn);
};

const unignoreColumn = (csvColumn: string) => {
  ignoredColumns.value.delete(csvColumn);
};

const getMappedField = (csvColumn: string) => {
  return formFields.find(f => f.slug === mappings.value[csvColumn]);
};

const getSampleData = (csvColumn: string) => {
  if (!csvData.value) return [];
  return csvData.value.data.slice(0, 3).map(row => row[csvColumn]);
};

const getColumnLetter = (index: number) => {
  return String.fromCharCode(65 + index);
};

const getColumnStatus = (csvColumn: string): 'confirmed' | 'ignored' | 'unmatched' | 'matched' => {
  if (confirmedColumns.value.has(csvColumn)) return 'confirmed';
  if (ignoredColumns.value.has(csvColumn)) return 'ignored';
  if (mappings.value[csvColumn]) return 'matched';
  return 'unmatched';
};

const isColumnConfirmed = (csvColumn: string) => confirmedColumns.value.has(csvColumn);
const isColumnIgnored = (csvColumn: string) => ignoredColumns.value.has(csvColumn);

// Auto-map columns when CSV is loaded
const autoMapColumns = () => {
  if (!csvData.value) return;

  const autoMappings: Record<string, string> = {};
  const headers = csvData.value.meta.fields;

  headers.forEach(header => {
    const normalized = header.toLowerCase().trim().replace(/\s+/g, '_');
    const match = formFields.find(f =>
      f.slug === normalized ||
      f.label.toLowerCase() === header.toLowerCase().trim()
    );

    if (match) {
      autoMappings[header] = match.slug;
    }
  });

  mappings.value = autoMappings;
};

// Watch for file changes and parse CSV
watch(() => props.file, (newFile) => {
  if (newFile) {
    Papa.parse(newFile, {
      header: true,
      skipEmptyLines: true,
      dynamicTyping: false,
      complete: (results) => {
        // Original parsed data
        const rawData = results.data;
        const rawFields = results.meta.fields;

        // Create a mapping from CSV header -> slug (if found)
        const headerToSlug: Record<string, string> = {};
        rawFields.forEach((header: string) => {
          const normalized = header.toLowerCase().trim().replace(/\s+/g, '_');
          const match = formFields.find(f =>
            f.slug === normalized ||
            f.label.toLowerCase() === header.toLowerCase().trim()
          );

          if (match) {
            headerToSlug[header] = match.slug;
          } else {
            // Keep the original header if no match
            headerToSlug[header] = normalized;
          }
        });

        // Transform data keys to use slugs
        const transformedData = rawData.map((row) => {
          const newRow = {};
          Object.keys(row).forEach(key => {
            const slug = headerToSlug[key];
            newRow[slug] = row[key];
          });
          return newRow;
        });

        // Save transformed data with slugs
        csvData.value = {
          data: transformedData,
          meta: { fields: Object.values(headerToSlug) }
        };

        // Auto-map columns (now using slugs)
        autoMapColumns();
      },
      error: (error: Error) => {
        console.error('Error parsing CSV:', error);
      }
    });
  }
}, { immediate: true });
</script>

<template>
  <div>
    <ImportCustomerTemplate>
      <template #header>
        <div class="flex gap-x-4">
          <div class="flex">
            <HugeiconsIcon
              :icon="CheckmarkCircle01Icon"
              class="my-auto"
              :size="21"
              fill="#007458"
              color="white"
              :stroke-width="3"
            />
            <h6 class="text-[#007458] font-semibold text-sm my-auto">
              Prepare
            </h6>
          </div>
          <HugeiconsIcon
            :icon="GreaterThanIcon"
            class="my-auto"
            :size="16"
            color="#939499"
            :stroke-width="3"
          />
          <div class="flex">
            <HugeiconsIcon
              :icon="CheckmarkCircle01Icon"
              class="my-auto"
              :size="21"
              fill="#007458"
              color="white"
              :stroke-width="3"
            />
            <h6 class="text-[#007458] font-semibold text-sm my-auto">
              Upload
            </h6>
          </div>
          <HugeiconsIcon
            :icon="GreaterThanIcon"
            class="my-auto"
            :size="16"
            color="#939499"
            :stroke-width="3"
          />
          <h6 class="text-[#1E212B] font-semibold text-sm my-auto">
            Match
          </h6>
          <HugeiconsIcon
            :icon="GreaterThanIcon"
            class="my-auto"
            :size="16"
            color="#939499"
            :stroke-width="3"
          />
          <h6 class="text-[#939499] font-semibold text-sm my-auto">
            Label
          </h6>
        </div>
      </template>

      <template #body>
        <div v-if="csvData">
          <div class="flex overflow-x-scroll border-b border-[#EDEDEE] px-4 md:px-8 mb-6">
            <div
              v-for="(col, index) in csvData.meta.fields"
              :key="index"
              :class="[
                'px-5 py-2 text-sm font-medium flex items-center gap-1.5',
                getColumnStatus(col) === 'confirmed' ? 'bg-[#EBF6F3] border-[#007458]' : '',
                getColumnStatus(col) === 'ignored' ? 'bg-[#EDEDEE] border-[#E8EAED]' : '',
                getColumnStatus(col) === 'matched' ? 'bg-white border-[#E8EAED]' : '',
                getColumnStatus(col) === 'unmatched' ? 'bg-white border-[#E8EAED]' : ''
              ]"
            >
              <HugeiconsIcon
                v-if="getColumnStatus(col) === 'confirmed'"
                :icon="CheckmarkCircle01Icon"
                :size="16"
                fill="#007458"
                color="white"
                :stroke-width="2"
              />
              <HugeiconsIcon
                v-else-if="getColumnStatus(col) === 'ignored'"
                :icon="UnavailableFreeIcons"
                :size="16"
                color="#6F7177"
                :stroke-width="2"
              />
              <span
                :class="[
                  getColumnStatus(col) === 'confirmed' ? 'text-[#007458]' : 'text-[#6F7177]'
                ]"
              >
                {{ getColumnLetter(index) }}
              </span>
            </div>
          </div>

          <!-- Mapping Cards -->
          <div class="space-y-6 px-4 md:px-8">
            <div
              v-for="(csvColumn, colIndex) in csvData.meta.fields"
              :key="colIndex"
              class="bg-white rounded-2xl shadow-[0_3px_6px_0_#424A531F] p-6"
            >
              <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
                <!-- Left Side: CSV Column -->
                <div class="col-span-1 border border-[#EDEDEE] rounded-2xl">
                  <div class="flex w-full">
                    <div class="w-full">
                      <div class="flex bg-[#F6F6F7]  rounded-t-lg border border-[#EDEDEE]">
                        <div class="flex md:w-1/2 px-4 py-2 md:px-8">
                          <div class="bg-[#F5F6F7] flex items-center justify-center font-bold text-[#1E212B]">
                            {{ getColumnLetter(colIndex) }}
                          </div>
                          <div class="font-bold text-[#1E212B] my-auto pl-4 md:pl-8">
                            {{ csvColumn }}
                          </div>
                        </div>
                        <select
                          :value="mappings[csvColumn] || 'none'"
                          class="md:w-1/2 py-2 md:ml-auto bg-white px-3 focus:outline-none focus:ring-2 focus:ring-[#2176AE] text-sm text-[#1E212B]"
                          @change="(e) => handleMappingChange(csvColumn, (e.target as HTMLSelectElement).value)"
                        >
                          <option value="none">
                            Lookup matching fields...
                          </option>
                          <option
                            v-for="field in formFields"
                            :key="field.slug"
                            :value="field.slug"
                          >
                            {{ field.label }}
                          </option>
                        </select>
                      </div>
                    </div>
                  </div>
                  <!-- Sample Data Preview -->
                  <div
                    v-if="!isColumnConfirmed(csvColumn) && !isColumnIgnored(csvColumn)"
                    class="pb-4"
                  >
                    <div
                      v-for="(value, idx) in getSampleData(csvColumn)"
                      :key="idx"
                      class="flex w-full items-center"
                    >
                      <span
                        class="text-sm px-4 py-2 md:px-8 text-[#939499] w-6"
                        :class="[idx+2 === 2 ? 'border-r border-b border-[#EDEDEE] rounded-br-sm mb-2': 'border-b border-[#EDEDEE]']"
                      >{{ idx + 2 }}</span>
                      <span
                        class="text-sm w-max pl-4 flex-1 py-2 text-[#1E212B]"
                        :class="[idx+2 !== 2 ? 'border-l border-b border-t border-[#EDEDEE] rounded-br-sm': '']"
                      >{{ value || '—' }}</span>
                    </div>
                  </div>
                </div>

                <!-- Right Side: Mapping Status -->
                <div class="flex flex-col justify-center">
                  <!-- Confirmed Header -->
                  <div
                    v-if="isColumnConfirmed(csvColumn) "
                    class="flex items-center justify-between my-auto"
                  >
                    <div class="flex items-center gap-2">
                      <HugeiconsIcon
                        :icon="Tick01Icon"
                        :size="30"
                        color="#007458"
                      />
                      <span class="text-sm font-semibold text-[#007458]">Confirmed mapping</span>
                    </div>
                    <GroBasicButton
                      size="xs"
                      class="w-max"
                      color="secondary"
                      @click="confirmedColumns.delete(csvColumn)"
                    >
                      Edit
                    </GroBasicButton>
                  </div>

                  <!-- Ignored Header -->
                  <div
                    v-if="isColumnIgnored(csvColumn)"
                    class="flex items-center justify-between my-auto"
                  >
                    <div class="flex items-center gap-2">
                      <HugeiconsIcon
                        :icon="UnavailableFreeIcons"
                        :size="20"
                        color="#6F7177"
                        :stroke-width="2"
                      />
                      <span class="text-sm font-semibold text-[#6F7177]">Ignored</span>
                    </div>
                    <GroBasicButton
                      size="xs"
                      class="w-max"
                      color="secondary"
                      @click="unignoreColumn(csvColumn)"
                    >
                      Edit
                    </GroBasicButton>
                  </div>
                  <template v-if="getMappedField(csvColumn) && !isColumnIgnored(csvColumn)">
                    <div
                      v-if="!isColumnConfirmed(csvColumn)"
                      class="flex items-start gap-2 mb-3"
                    >
                      <HugeiconsIcon
                        :icon="Tick01Icon"
                        class=""
                        :size="25"
                        color="#007458"
                      />
                      <div class="text-sm text-[#1E212B] font-semibold mt-auto">
                        <span>Matched to the </span>
                        <span class="bg-[#F5F6F7] px-2 py-1 rounded font-medium">
                          {{ getMappedField(csvColumn)?.label }}
                        </span>
                        <span> field.</span>
                      </div>
                    </div>

                    <div
                      v-if="!isColumnConfirmed(csvColumn)"
                      class="flex items-start gap-2 mb-4"
                    >
                      <HugeiconsIcon
                        :icon="Tick01Icon"
                        class=""
                        :size="25"
                        color="#007458"
                      />
                      <span class="text-sm my-auto text-[#6F7177]">
                        All values pass validation for this field
                      </span>
                    </div>

                    <div
                      v-if="!isColumnConfirmed(csvColumn)"
                      class="flex items-start gap-2 mb-6"
                    >
                      <div class="w-5 h-5 shrink-0 mt-0.5 rounded-full bg-[#2A353D] flex items-center justify-center">
                        <span class="text-xs text-white">i</span>
                      </div>
                      <span class="text-sm text-[#6F7177] my-auto ml-1.5">
                        {{ validateColumn(csvColumn).coverage }}% of your rows have a value for this column
                      </span>
                    </div>

                    <div
                      v-if="!isColumnConfirmed(csvColumn)"
                      class="flex gap-3"
                    >
                      <GroBasicButton
                        color="primary"
                        size="sm"
                        class="w-max"
                        shape="custom"
                        @click="confirmMapping(csvColumn)"
                      >
                        <template #default>
                          Confirm mapping
                        </template>
                      </GroBasicButton>
                      <GroBasicButton
                        color="tertiary"
                        size="sm"
                        class="w-max"
                        shape="custom"
                        @click="ignoreColumn(csvColumn)"
                      >
                        <template #default>
                          Ignore this column
                        </template>
                      </GroBasicButton>
                    </div>
                  </template>
                  <div
                    v-else-if="!isColumnIgnored(csvColumn)"
                    class="space-y-4"
                  >
                    <div class="flex items-start gap-2">
                      <div class="w-5 h-5 shrink-0 mt-0.5 rounded-full bg-[#FFA726] flex items-center justify-center">
                        <span class="text-xs text-white font-bold">!</span>
                      </div>
                      <span class="text-sm text-[#1E212B] font-medium">
                        Unable to automatically match
                      </span>
                    </div>
                    <p class="text-sm text-[#6F7177]">
                      Please select a field from the dropdown above to map this column, or ignore it if not needed.
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Navigation Buttons -->
          <div class="flex w-full mt-8 pb-4 px-4 md:pb-8 md:px-8">
            <GroBasicButton
              color="tertiary"
              size="sm"
              class="w-max"
              shape="custom"
              @click="stepModel = stepModel - 1"
            >
              <template #frontIcon>
                <HugeiconsIcon
                  :icon="GreaterThanIcon"
                  size="12"
                  :stroke-width="3"
                  class="my-auto rotate-180"
                />
              </template>
              <template #default>
                Back
              </template>
            </GroBasicButton>
            <GroBasicButton
              color="primary"
              size="sm"
              class="w-max ml-auto"
              shape="custom"
              :disabled="!canProceed()"
              @click="stepModel = stepModel + 1"
            >
              <template #default>
                Next
              </template>
            </GroBasicButton>
          </div>
        </div>

        <div
          v-else
          class="text-center py-12 text-[#6F7177]"
        >
          No CSV data loaded
        </div>
      </template>
    </ImportCustomerTemplate>
  </div>
</template>

<style scoped>
</style>

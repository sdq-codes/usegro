<script setup lang="ts">
import { ref, onMounted, defineModel } from "vue";
import GroBasicInput from "@/components/forms/input/GroBasicInput.vue";
import GroBasicButton from "@/components/buttons/GroBasicButton.vue";
import {PlusSignIcon} from "@hugeicons/core-free-icons";
import {HugeiconsIcon} from "@hugeicons/vue";
import {useCrmTagsAPI} from "@/composables/api/tag/tag.js";
const model = defineModel<string[]>({ default: () => [] })


const crmId = ref(null);
const tags = ref([]);
const newTagName = ref("");
const inputColor = ref("primary");
const addTagInputHint = ref("");
const showAddTag = ref(false);

// Fetch tags
const fetchTags = async () => {
  const fetchTagApiResponse = await useCrmTagsAPI().FetchCRMTag()
  if (fetchTagApiResponse.success) {
    tags.value = fetchTagApiResponse.data?.data;
  }
};

// Create a new tag
const createTag = async () => {
  const createTagApiResponse = await useCrmTagsAPI().CreateTag({tag: newTagName.value})
  if (createTagApiResponse.success) {
    tags.value = createTagApiResponse.data?.data;
    newTagName.value = "";
    showAddTag.value = false;
  } else {
    addTagInputHint.value = "Error fetching tag";
    inputColor.value = "error";
    showAddTag.value = true;
  }
};

// Handle selection
const selectTag = (tagId) => {
  if (model.value.includes(tagId)) {
    model.value = model.value.filter(item => item !== tagId);
  } else {
    model.value = [...model.value, tagId];
  }
};

onMounted(() => {
  crmId.value = localStorage.getItem("crmId");
  fetchTags();
});
</script>

<template>
  <div class="rounded-2xl border border-[#E8EAED] px-4 py-3">
    <div
      v-if="!showAddTag"
      class="flex flex-wrap gap-4 mb-3"
    >
      <button
        v-for="tag in tags"
        :key="tag.id"
        class="px-4 py-1.5 rounded-full cursor-pointer text-xs capitalize font-semibold"
        :class="{
          'bg-[#1E3A5F] text-white': model?.includes(tag.id),
          'bg-[#DBEAFE] text-[#1E3A5F]': !model?.includes(tag.id),
        }"
        @click.prevent="selectTag(tag.id)"
      >
        {{ tag.tag }}
      </button>
    </div>

    <div
      v-if="!showAddTag"
      class="flex cursor-pointer"
      @click.prevent="showAddTag = !showAddTag"
    >
      <HugeiconsIcon
        :icon="PlusSignIcon"
        :size="12"
        color="#2176AE"
        :stroke-width="3"
        class="my-auto"
      />
      <h6 class="text-xs text-[#2176AE] ml-1">
        Add labels
      </h6>
    </div>

    <!-- Create new tag -->
    <div
      v-if="showAddTag"
      class="grid grid-cols-3 items-center gap-2"
    >
      <GroBasicInput
        v-model="newTagName"
        class="col-span-2"
        type="text"
        placeholder="New tag name"
        :color="inputColor"
        :hint="addTagInputHint"
      />
      <div class="flex gap-x-2 col-span-1">
        <GroBasicButton
          @click.prevent="createTag"
        >
          Add
        </GroBasicButton>
        <GroBasicButton
          color="secondary"
          @click.prevent="showAddTag = !showAddTag"
        >
          Cancel
        </GroBasicButton>
      </div>
    </div>
  </div>
</template>

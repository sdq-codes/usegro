<script setup lang="ts">
import { ref, onMounted, defineModel } from "vue";
import GroBasicInput from "@/components/forms/input/GroBasicInput.vue";
import GroBasicButton from "@/components/buttons/GroBasicButton.vue";
import {PlusSignIcon} from "@hugeicons/core-free-icons";
import {HugeiconsIcon} from "@hugeicons/vue";
import {useCrmTagsAPI} from "@/composables/api/tag/tag.js";
import type {FetchCrmTagResponse} from "@/composables/dto/tag/tag";

const model = defineModel<FetchCrmTagResponse[]>({ default: () => [] })


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
    model.value.push(tagId);
  }
};

onMounted(() => {
  crmId.value = localStorage.getItem("crmId");
  model.value = [];
  fetchTags();
});
</script>

<template>
  <div class="rounded-lg">
    <div
      v-if="!showAddTag"
      class="flex flex-wrap gap-4 mb-3"
    >
      <button
        v-for="tag in tags"

        :key="tag.SK"
        class="px-3 py-1 rounded-lg cursor-pointer text-xs capitalize font-semibold"
        :class="{
          'bg-[#1E212B] text-white': model?.includes(tag.SK),
          'bg-gray-100': !model?.includes(tag.SK),
        }"
        @click.prevent="selectTag(tag.SK)"
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
        color="#070707"
        :stroke-width="3"
        class="my-auto"
      />
      <h6 class="text-xs border-b text-blue-500">
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

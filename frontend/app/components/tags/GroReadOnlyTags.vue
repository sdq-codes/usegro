<script setup lang="ts">
import { ref, onMounted, defineProps, computed } from "vue";
import { useCrmTagsAPI } from "@/composables/api/tag/tag.js";
import type { FetchCrmTagResponse } from "@/composables/dto/tag/tag";

// Props — accept a list of tag IDs
const props = defineProps<{
  tagids: string[];
}>();

const tags = ref<FetchCrmTagResponse[]>([]);

// Fetch all tags
const fetchTags = async () => {
  const fetchTagApiResponse = await useCrmTagsAPI().FetchCRMTag();

  if (fetchTagApiResponse.success) {
    tags.value = fetchTagApiResponse.data?.data || [];
  }
};

const filteredTags = computed(() => {
    if (tags.value && props.tagids) {
      return tags.value.filter(tag => props.tagids.includes(tag.id))
    }
    return []
  }
);

onMounted(fetchTags);
</script>

<template>
  <div class="rounded-lg flex flex-wrap gap-4 mb-3">
    <template v-if="filteredTags.length">
      <span
        v-for="tag in filteredTags"
        :key="tag.id"
        class="px-3 py-1 rounded-lg text-xs capitalize font-semibold bg-[#1E212B] text-white cursor-default select-none"
      >
        {{ tag.tag }}
      </span>
    </template>

    <p
      v-else
      class="text-sm text-gray-500 italic"
    >
      No tags available
    </p>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import GroBasicInput from '@/components/forms/input/GroBasicInput.vue'
import GroBasicButton from '@/components/buttons/GroBasicButton.vue'
import { HugeiconsIcon } from '@hugeicons/vue'
import { PlusSignIcon } from '@hugeicons/core-free-icons'
import { useCatalogTagAPI } from '@/composables/api/catalog/tag'

const model = defineModel<string[]>({ default: () => [] })

const tags = ref<Array<{ id: string; name: string }>>([])
const newTagName = ref('')
const inputColor = ref<'primary' | 'error'>('primary')
const addTagInputHint = ref('')
const showAddTag = ref(false)

const fetchTags = async () => {
  const result = await useCatalogTagAPI().ListTags()
  if (result.success && result.data) {
    tags.value = (result.data as any).data ?? result.data
  }
}

const createTag = async () => {
  const result = await useCatalogTagAPI().CreateTag(newTagName.value)
  if (result.success && result.data) {
    await fetchTags()
    const created = (result.data as any).data ?? result.data
    model.value = [...model.value, created.id]
    newTagName.value = ''
    showAddTag.value = false
  } else {
    addTagInputHint.value = 'Error creating tag'
    inputColor.value = 'error'
  }
}

const selectTag = (tagId: string) => {
  if (model.value.includes(tagId)) {
    model.value = model.value.filter(id => id !== tagId)
  } else {
    model.value = [...model.value, tagId]
  }
}

onMounted(fetchTags)
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
        :class="model.includes(tag.id)
          ? 'bg-[#1E3A5F] text-white'
          : 'bg-[#DBEAFE] text-[#1E3A5F]'"
        @click.prevent="selectTag(tag.id)"
      >
        {{ tag.name }}
      </button>
    </div>

    <div
      v-if="!showAddTag"
      class="flex cursor-pointer"
      @click.prevent="showAddTag = true"
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
        <GroBasicButton @click.prevent="createTag">
          Add
        </GroBasicButton>
        <GroBasicButton
          color="secondary"
          @click.prevent="showAddTag = false"
        >
          Cancel
        </GroBasicButton>
      </div>
    </div>
  </div>
</template>

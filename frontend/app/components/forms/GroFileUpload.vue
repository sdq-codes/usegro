<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'

interface UploadedFile {
  id: string
  file: File
  objectUrl: string
  error: string | null
}

interface Props {
  accept?: string
  maxFiles?: number
  maxSizeMb?: number
  label?: string
  hint?: string
}

const props = withDefaults(defineProps<Props>(), {
  accept: 'image/*',
  maxFiles: 10,
  maxSizeMb: 10,
  label: 'Add Files',
  hint: 'Accepts images, videos or 3D models',
})

const emit = defineEmits<{
  (e: 'update:files', files: File[]): void
  (e: 'error', message: string): void
}>()

const uploadedFiles = ref<UploadedFile[]>([])
const isDragging = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)

// Lightbox state
const lightboxEntry = ref<UploadedFile | null>(null)

const isImage = (file: File) => file.type.startsWith('image/')
const isVideo = (file: File) => file.type.startsWith('video/')

const validateFile = (file: File): string | null => {
  if (props.maxSizeMb && file.size > props.maxSizeMb * 1024 * 1024) {
    return `File exceeds ${props.maxSizeMb}MB limit`
  }
  return null
}

const addFiles = (files: FileList | File[]) => {
  const incoming = Array.from(files)
  const remaining = props.maxFiles - uploadedFiles.value.length

  if (remaining <= 0) {
    emit('error', `Maximum ${props.maxFiles} files allowed`)
    return
  }

  const toAdd = incoming.slice(0, remaining)

  for (const file of toAdd) {
    const error = validateFile(file)
    const objectUrl = URL.createObjectURL(file)
    uploadedFiles.value.push({
      id: `${Date.now()}-${Math.random().toString(36).slice(2)}`,
      file,
      objectUrl,
      error,
    })
  }

  emitFiles()
}

const removeFile = (id: string) => {
  const entry = uploadedFiles.value.find(f => f.id === id)
  if (entry) URL.revokeObjectURL(entry.objectUrl)
  if (lightboxEntry.value?.id === id) lightboxEntry.value = null
  uploadedFiles.value = uploadedFiles.value.filter(f => f.id !== id)
  emitFiles()
}

const emitFiles = () => {
  emit('update:files', uploadedFiles.value.filter(f => !f.error).map(f => f.file))
}

const onInputChange = (e: Event) => {
  const input = e.target as HTMLInputElement
  if (input.files) addFiles(input.files)
  input.value = ''
}

const onDrop = (e: DragEvent) => {
  isDragging.value = false
  if (e.dataTransfer?.files) addFiles(e.dataTransfer.files)
}

const openLightbox = (entry: UploadedFile) => {
  if (!entry.error) lightboxEntry.value = entry
}

const closeLightbox = () => {
  lightboxEntry.value = null
}

const formatSize = (bytes: number) => {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}

const fileCountLabel = computed(() => `${uploadedFiles.value.length}/${props.maxFiles}`)
const atLimit = computed(() => uploadedFiles.value.length >= props.maxFiles)

// Keyboard close
const onKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Escape') closeLightbox()
}
onMounted(() => document.addEventListener('keydown', onKeydown))
onBeforeUnmount(() => {
  document.removeEventListener('keydown', onKeydown)
  uploadedFiles.value.forEach(f => URL.revokeObjectURL(f.objectUrl))
})
</script>

<template>
  <div>
    <!-- File count -->
    <div class="flex items-center space-x-4 mb-2">
      <slot name="label" />
      <span class="text-xs bg-[#EDEDEE] px-2 py-1 rounded-xl text-[#070707]">{{ fileCountLabel }}</span>
    </div>

    <!-- Drop zone -->
    <label
      v-if="!atLimit"
      :class="isDragging
        ? 'border-[#2176AE] bg-[#EBF4FF]'
        : 'border-[#94BDD8] hover:bg-[#F0F7FF]'"
      @dragenter.prevent="isDragging = true"
      @dragover.prevent="isDragging = true"
      @dragleave.prevent="isDragging = false"
      @drop.prevent="onDrop"
    >
      <div
        class="relative border-2 border-dashed rounded-xl py-10 flex flex-col items-center justify-center gap-2 cursor-pointer transition-colors"
      >
        <svg
          class="w-9 h-9 text-[#2176AE]"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="1.5"
        >
          <path d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
        </svg>
        <span class="text-sm font-semibold text-[#2176AE]">
          {{ isDragging ? 'Drop files here' : label }}
        </span>
        <span class="text-xs text-[#6F7177]">{{ hint }}</span>
        <span class="text-xs text-[#939499]">
          Max {{ maxSizeMb }}MB per file · up to {{ maxFiles }} files
        </span>
        <input
          ref="fileInput"
          type="file"
          :accept="accept"
          :multiple="maxFiles > 1"
          class="hidden"
          @change="onInputChange"
        >
      </div>
    </label>

    <!-- Add more -->
    <button
      v-if="uploadedFiles.length > 0 && !atLimit"
      type="button"
      class="mt-2 cursor-pointer text-xs font-medium text-[#2176AE] hover:underline flex items-center gap-1"
      @click="fileInput?.click()"
    >
      <svg
        class="w-3.5 h-3.5"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2.5"
      >
        <path d="M12 5v14M5 12h14" />
      </svg>
      Add more
    </button>

    <!-- Preview grid -->
    <div
      v-if="uploadedFiles.length > 0"
      class="mt-3 grid grid-cols-3 sm:grid-cols-4 gap-2"
    >
      <div
        v-for="entry in uploadedFiles"
        :key="entry.id"
        class="relative group rounded-xl overflow-hidden border bg-[#F6F6F7]"
        :class="entry.error ? 'border-red-300' : 'border-[#EDEDEE]'"
      >
        <!-- Image thumbnail -->
        <img
          v-if="isImage(entry.file)"
          :src="entry.objectUrl"
          :alt="entry.file.name"
          class="w-full aspect-square object-cover"
        >

        <!-- Video thumbnail with play overlay -->
        <div
          v-else-if="isVideo(entry.file)"
          class="relative w-full aspect-square bg-black"
        >
          <video
            :src="entry.objectUrl"
            class="w-full h-full object-cover opacity-80"
            preload="metadata"
            muted
          />
          <div class="absolute inset-0 flex items-center justify-center pointer-events-none">
            <div class="w-9 h-9 rounded-full bg-black/50 flex items-center justify-center">
              <svg
                class="w-4 h-4 text-white ml-0.5"
                viewBox="0 0 24 24"
                fill="currentColor"
              >
                <path d="M8 5v14l11-7z" />
              </svg>
            </div>
          </div>
        </div>

        <!-- Generic file -->
        <div
          v-else
          class="w-full aspect-square flex flex-col items-center justify-center gap-1 text-[#6F7177] p-2"
        >
          <svg
            class="w-8 h-8"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="1.5"
          >
            <path d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          <span class="text-[10px] font-medium px-1 text-center truncate w-full">{{ entry.file.name }}</span>
        </div>

        <!-- Error overlay -->
        <div
          v-if="entry.error"
          class="absolute inset-0 bg-red-50/90 flex items-center justify-center p-1"
        >
          <span class="text-[10px] text-red-600 font-medium text-center leading-tight">{{ entry.error }}</span>
        </div>

        <!-- Hover overlay -->
        <div
          v-if="!entry.error"
          class="absolute inset-0 bg-black/0 group-hover:bg-black/40 transition-colors opacity-0 group-hover:opacity-100 flex flex-col items-end justify-start p-1.5 gap-1"
        >
          <!-- View button -->
          <button
            type="button"
            class="w-7 h-7 rounded-full bg-white/90 flex items-center justify-center hover:bg-white transition-colors"
            title="View"
            @click.prevent="openLightbox(entry)"
          >
            <svg
              class="w-3.5 h-3.5 text-[#1E212B]"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <path d="M15 3h6v6M9 21H3v-6M21 3l-7 7M3 21l7-7" />
            </svg>
          </button>
          <!-- Remove button -->
          <button
            type="button"
            class="w-7 h-7 rounded-full bg-white/90 flex items-center justify-center hover:bg-white transition-colors"
            title="Remove"
            @click.prevent="removeFile(entry.id)"
          >
            <svg
              class="w-3.5 h-3.5 text-red-500"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2.5"
            >
              <path d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
          <!-- File size -->
          <span class="mt-auto text-[9px] text-white font-medium">{{ formatSize(entry.file.size) }}</span>
        </div>
      </div>
    </div>

    <!-- Lightbox -->
    <Teleport to="body">
      <Transition
        enter-active-class="transition ease-out duration-200"
        enter-from-class="opacity-0"
        enter-to-class="opacity-100"
        leave-active-class="transition ease-in duration-150"
        leave-from-class="opacity-100"
        leave-to-class="opacity-0"
      >
        <div
          v-if="lightboxEntry"
          class="fixed inset-0 z-[9999] flex items-center justify-center bg-black/80 p-4"
          @click.self="closeLightbox"
        >
          <!-- Close -->
          <button
            type="button"
            class="absolute top-4 right-4 w-9 h-9 rounded-full bg-white/10 hover:bg-white/20 flex items-center justify-center text-white transition-colors"
            @click="closeLightbox"
          >
            <svg
              class="w-5 h-5"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <path d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>

          <!-- Image viewer -->
          <img
            v-if="isImage(lightboxEntry.file)"
            :src="lightboxEntry.objectUrl"
            :alt="lightboxEntry.file.name"
            class="max-w-full max-h-full rounded-xl object-contain shadow-2xl"
            style="max-height: 90vh"
          >

          <!-- Video player -->
          <video
            v-else-if="isVideo(lightboxEntry.file)"
            :src="lightboxEntry.objectUrl"
            class="max-w-full rounded-xl shadow-2xl"
            style="max-height: 90vh"
            controls
            autoplay
          />

          <!-- File info bar -->
          <div class="absolute bottom-4 left-1/2 -translate-x-1/2 bg-black/50 text-white text-xs px-4 py-2 rounded-full backdrop-blur-sm flex items-center gap-3">
            <span class="font-medium truncate max-w-48">{{ lightboxEntry.file.name }}</span>
            <span class="text-white/60">{{ formatSize(lightboxEntry.file.size) }}</span>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

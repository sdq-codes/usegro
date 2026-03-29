<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import Link from '@tiptap/extension-link'
import Placeholder from '@tiptap/extension-placeholder'

interface Props {
  placeholder?: string
}

withDefaults(defineProps<Props>(), {
  placeholder: 'Leave a comment...',
})

const model = defineModel<string>({ default: '' })

// ── Editor ────────────────────────────────────────────────────────────────────
const editor = useEditor({
  content: model.value,
  extensions: [
    StarterKit,
    Link.configure({ openOnClick: false }),
    Placeholder.configure({ placeholder: 'Leave a comment...' }),
  ],
  editorProps: {
    attributes: {
      class: 'prose prose-sm max-w-none min-h-[72px] px-4 py-3 focus:outline-none text-sm text-gray-800',
    },
  },
  onUpdate: ({ editor: e }) => {
    model.value = e.isEmpty ? '' : e.getHTML()
  },
})

// Keep model in sync if parent resets to ''
watch(model, (val) => {
  if (!val && editor.value && !editor.value.isEmpty) {
    editor.value.commands.clearContent(true)
  }
})

onBeforeUnmount(() => editor.value?.destroy())

// ── Emoji picker ──────────────────────────────────────────────────────────────
const showEmojiPicker = ref(false)
const EMOJI_GROUPS: Record<string, string[]> = {
  'Smileys': ['😀','😂','😍','🥹','😊','😎','🤔','😅','🙄','😭','🥰','😘','🤗','😤','🤩'],
  'Gestures': ['👍','👎','👋','🙌','👏','🤝','🙏','✌️','🤞','💪','🫶','❤️','🔥','✅','🎉'],
  'People': ['👤','👩','👨','👶','🧑','👩‍💻','👨‍💻','🧑‍🎨','👩‍🏫','👨‍🏫'],
  'Objects': ['📎','📌','📝','💡','🔗','📅','📊','💬','📧','📱'],
}
const activeEmojiGroup = ref('Smileys')

const insertEmoji = (emoji: string) => {
  editor.value?.commands.insertContent(emoji)
  showEmojiPicker.value = false
}

const closeEmoji = (e: MouseEvent) => {
  if (!(e.target as Element).closest('.emoji-picker-container')) {
    showEmojiPicker.value = false
  }
}
onMounted(() => document.addEventListener('click', closeEmoji))
onBeforeUnmount(() => document.removeEventListener('click', closeEmoji))

// ── Link dialog ───────────────────────────────────────────────────────────────
const showLinkDialog = ref(false)
const linkUrl = ref('')

const openLinkDialog = () => {
  linkUrl.value = editor.value?.getAttributes('link').href ?? ''
  showLinkDialog.value = true
}

const applyLink = () => {
  if (!linkUrl.value.trim()) {
    editor.value?.chain().focus().unsetLink().run()
  } else {
    editor.value?.chain().focus().setLink({ href: linkUrl.value.trim() }).run()
  }
  showLinkDialog.value = false
  linkUrl.value = ''
}

const removeLink = () => {
  editor.value?.chain().focus().unsetLink().run()
  showLinkDialog.value = false
}
</script>

<template>
  <div class="rounded-lg bg-white overflow-hidden">
    <!-- Formatting toolbar -->
    <div class="flex items-center gap-0.5 px-2 pt-2 border-b border-gray-100">
      <button
        type="button"
        class="p-1.5 rounded text-gray-500 hover:bg-gray-100 transition-colors text-xs font-bold"
        :class="{ 'bg-gray-100 text-gray-900': editor?.isActive('bold') }"
        title="Bold"
        @click="editor?.chain().focus().toggleBold().run()"
      >
        B
      </button>
      <button
        type="button"
        class="p-1.5 rounded text-gray-500 hover:bg-gray-100 transition-colors text-xs italic"
        :class="{ 'bg-gray-100 text-gray-900': editor?.isActive('italic') }"
        title="Italic"
        @click="editor?.chain().focus().toggleItalic().run()"
      >
        I
      </button>
      <button
        type="button"
        class="p-1.5 rounded text-gray-500 hover:bg-gray-100 transition-colors text-xs font-mono line-through"
        :class="{ 'bg-gray-100 text-gray-900': editor?.isActive('strike') }"
        title="Strikethrough"
        @click="editor?.chain().focus().toggleStrike().run()"
      >
        S
      </button>
      <div class="w-px h-4 bg-gray-200 mx-1" />
      <button
        type="button"
        class="p-1.5 rounded text-gray-500 hover:bg-gray-100 transition-colors"
        :class="{ 'bg-gray-100 text-gray-900': editor?.isActive('bulletList') }"
        title="Bullet list"
        @click="editor?.chain().focus().toggleBulletList().run()"
      >
        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M4 6h16M4 12h16M4 18h16"
          />
        </svg>
      </button>
      <button
        type="button"
        class="p-1.5 rounded text-gray-500 hover:bg-gray-100 transition-colors"
        :class="{ 'bg-gray-100 text-gray-900': editor?.isActive('orderedList') }"
        title="Numbered list"
        @click="editor?.chain().focus().toggleOrderedList().run()"
      >
        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M7 20h14M7 12h14M7 4h14M3 20v-2M3 12v-2M3 4V2"
          />
        </svg>
      </button>
      <button
        type="button"
        class="p-1.5 rounded text-gray-500 hover:bg-gray-100 transition-colors"
        :class="{ 'bg-gray-100 text-gray-900': editor?.isActive('blockquote') }"
        title="Quote"
        @click="editor?.chain().focus().toggleBlockquote().run()"
      >
        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
          />
        </svg>
      </button>
    </div>

    <!-- Editor content -->
    <EditorContent :editor="editor" />

    <!-- Bottom toolbar -->
    <div class="flex items-center justify-between px-3 py-2 border-t bg-[#F6F6F7] border-gray-100">
      <div class="flex items-center gap-1">
        <!-- Emoji picker -->
        <div class="relative emoji-picker-container">
          <button
            type="button"
            class="p-1.5 rounded text-gray-400 hover:bg-gray-100 hover:text-gray-600 transition-colors"
            title="Emoji"
            @click.stop="showEmojiPicker = !showEmojiPicker"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M14.828 14.828a4 4 0 01-5.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
          </button>
          <Transition
            enter-active-class="transition ease-out duration-150"
            enter-from-class="opacity-0 scale-95"
            enter-to-class="opacity-100 scale-100"
            leave-active-class="transition ease-in duration-100"
            leave-from-class="opacity-100"
            leave-to-class="opacity-0 scale-95"
          >
            <div
              v-show="showEmojiPicker"
              class="absolute bottom-full left-0 mb-2 w-64 bg-white border border-gray-200 rounded-xl shadow-xl z-50 p-2 origin-bottom-left"
            >
              <!-- Group tabs -->
              <div class="flex gap-1 mb-2 flex-wrap">
                <button
                  v-for="group in Object.keys(EMOJI_GROUPS)"
                  :key="group"
                  type="button"
                  class="px-2 py-0.5 text-xs rounded-full transition-colors"
                  :class="activeEmojiGroup === group ? 'bg-[#1E212B] text-white' : 'text-gray-500 hover:bg-gray-100'"
                  @click.stop="activeEmojiGroup = group"
                >
                  {{ group }}
                </button>
              </div>
              <!-- Emoji grid -->
              <div class="grid grid-cols-8 gap-0.5">
                <button
                  v-for="emoji in EMOJI_GROUPS[activeEmojiGroup]"
                  :key="emoji"
                  type="button"
                  class="p-1 text-lg hover:bg-gray-100 rounded cursor-pointer"
                  @click.stop="insertEmoji(emoji)"
                >
                  {{ emoji }}
                </button>
              </div>
            </div>
          </Transition>
        </div>

        <!-- Link -->
        <button
          type="button"
          class="p-1.5 rounded text-gray-400 hover:bg-gray-100 hover:text-gray-600 transition-colors"
          :class="{ 'text-blue-500 bg-blue-50': editor?.isActive('link') }"
          title="Insert link"
          @click="openLinkDialog"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"
            />
          </svg>
        </button>
      </div>

      <slot name="actions" />
    </div>

    <!-- Link dialog -->
    <Transition
      enter-active-class="transition ease-out duration-150"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition ease-in duration-100"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="showLinkDialog"
        class="border-t border-gray-100 px-3 py-2 bg-gray-50 flex items-center gap-2"
      >
        <svg class="w-4 h-4 text-gray-400 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"
          />
        </svg>
        <input
          v-model="linkUrl"
          type="url"
          placeholder="https://..."
          class="flex-1 text-sm bg-transparent border-b border-gray-300 focus:border-blue-500 focus:outline-none py-0.5"
          @keydown.enter="applyLink"
          @keydown.escape="showLinkDialog = false"
        />
        <button
          type="button"
          class="text-xs text-blue-600 font-medium hover:underline"
          @click="applyLink"
        >
          Apply
        </button>
        <button
          v-if="editor?.isActive('link')"
          type="button"
          class="text-xs text-red-500 hover:underline"
          @click="removeLink"
        >
          Remove
        </button>
        <button
          type="button"
          class="text-xs text-gray-400 hover:text-gray-600"
          @click="showLinkDialog = false"
        >
          ✕
        </button>
      </div>
    </Transition>
  </div>
</template>

<style>
/* Tiptap prose styles */
.tiptap p { margin: 0; }
.tiptap p + p { margin-top: 0.5em; }
.tiptap ul { list-style: disc; padding-left: 1.25rem; }
.tiptap ol { list-style: decimal; padding-left: 1.25rem; }
.tiptap blockquote {
  border-left: 3px solid #e5e7eb;
  padding-left: 0.75rem;
  color: #6b7280;
  font-style: italic;
  margin: 0.5rem 0;
}
.tiptap a { color: #2176AE; text-decoration: underline; }
.tiptap strong { font-weight: 600; }
.tiptap em { font-style: italic; }
.tiptap s { text-decoration: line-through; }
.tiptap p.is-editor-empty:first-child::before {
  color: #9ca3af;
  content: attr(data-placeholder);
  float: left;
  height: 0;
  pointer-events: none;
}
</style>

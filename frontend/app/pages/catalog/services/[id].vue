<script setup lang="ts">
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue'
import MainDashboard from '@/components/dashboard/main-dashboard.vue'
import GroBasicInput from '@/components/forms/input/GroBasicInput.vue'
import GroBasicTextArea from '@/components/forms/input/GroBasicTextArea.vue'
import GroBasicSelect from '@/components/forms/select/GroBasicSelect.vue'
import GroCurrencyInput from '@/components/forms/input/GroCurrencyInput.vue'
import GroBasicRadio from '@/components/forms/select/GroBasicRadio.vue'
import GroBasicButton from '@/components/buttons/GroBasicButton.vue'
import GroToggle from '@/components/forms/select/GroToggle.vue'
import CatalogTagSelect from '@/components/forms/select/CatalogTagSelect.vue'
import BackLink from '@/components/navigation/BackLink.vue'
import GroFormSection from '@/components/forms/GroFormSection.vue'
import CommentEditor from '@/components/editor/CommentEditor.vue'
import { useRoute, useRouter } from 'nuxt/app'
import { useServiceAPI } from '@/composables/api/catalog/service'
import { useCatalogMediaAPI } from '@/composables/api/catalog/media'
import { notify } from '@/composables/helpers/notification/notification'

const route = useRoute()
const router = useRouter()
const serviceId = route.params.id as string
const isNew = serviceId === 'new'

const isLoading = ref(!isNew)

// ── Step 1: type picker / Step 2: form ─────────────────────────
const step = ref<1 | 2>(isNew ? 1 : 2)

const serviceTypes = [
  {
    key: 'appointment' as const,
    label: 'Appointment',
    desc: 'Private session that can be booked according to availability',
  },
  {
    key: 'class' as const,
    label: 'Class',
    desc: 'A group session that can recur. Clients book any session they want to join.',
  },
  {
    key: 'course' as const,
    label: 'Course',
    desc: 'A set of group sessions. Clients book them all up front',
  },
]

const selectType = (type: 'appointment' | 'class' | 'course') => {
  serviceType.value = type
  step.value = 2
}

// ── Service type ───────────────────────────────────────────────
const serviceType = ref<'appointment' | 'class' | 'course'>('appointment')

// ── Accordion open/closed state ────────────────────────────────
const open = ref({
  serviceType: true,
  details: true,
  price: false,
  schedule: false,
  location: false,
  images: false,
  booking: false,
})
const toggle = (key: keyof typeof open.value) => { open.value[key] = !open.value[key] }

// ── Service details ────────────────────────────────────────────
const form = ref({
  name: '',
  tagline: '',
  description: '',
  duration: '1 Hour',
  bufferTime: '10 minutes',
  status: 'Active',
  showInSite: false,
})

const durationOptions = [
  '15 minutes', '30 minutes', '45 minutes',
  '1 Hour', '1.5 Hours', '2 Hours', '2.5 Hours', '3 Hours', '3.5 Hours',
  '4 Hours', '4.5 Hours', '5 Hours', '5.5 Hours', '6 Hours', '6.5 Hours',
  '7 Hours', '7.5 Hours', '8 Hours', '9 Hours', '10 Hours', '11 Hours',
  '12 Hours', '13 Hours', '14 Hours', '15 Hours', '16 Hours', '17 Hours',
  '18 Hours', '19 Hours', '20 Hours', '21 Hours', '22 Hours', '23 Hours',
  '24 Hours', '36 Hours', '48 Hours',
].map(v => ({ value: v, label: v }))

const bufferOptions = [
  'No buffer', '5 minutes', '10 minutes', '15 minutes', '30 minutes', '1 Hour', '1.5 Hours', '2 Hours', '3 Hours',
].map(v => ({ value: v, label: v }))

const priceTypeOptions = [
  { value: 'fixed', label: 'Fixed Price', description: 'Specify an amount (e.g., ₦2,000)' },
  { value: 'free', label: 'Free', description: 'The price will not be shown on your website' },
  { value: 'variable', label: 'Varied Pricing', description: 'Provide various pricing options (e.g., Day session ₦2,000, Night session ₦3,000)' },
  { value: 'custom', label: 'Custom Price', description: 'Explain your pricing structure (e.g., negotiable pricing)' },
]

// ── Price & Payments ───────────────────────────────────────────
type PaymentMode = 'per-session' | 'with-plan' | 'per-session-or-plan'
type PriceType = 'fixed' | 'variable' | 'free' | 'custom'

const paymentMode = ref<PaymentMode>('per-session')
const priceType = ref<PriceType>('fixed')
const priceAmount = ref('')
const priceCurrency = ref((process.client && localStorage.getItem('currency')) || 'USD')
const customPriceLabel = ref('')

interface Variation { id: string; name: string; price: string }
const variations = ref<Variation[]>([])
const showVariationsModal = ref(false)
const variationsTitle = ref('')
const editVariations = ref<Variation[]>([])

const openVariationsModal = () => {
  variationsTitle.value = ''
  editVariations.value = variations.value.map(v => ({ ...v }))
  showVariationsModal.value = true
}
const addEditVariation = () => {
  editVariations.value.push({ id: String(Date.now()), name: '', price: '' })
}
const removeEditVariation = (id: string) => {
  editVariations.value = editVariations.value.filter(v => v.id !== id)
}
const saveVariations = () => {
  variations.value = editVariations.value.filter(v => v.name.trim())
  showVariationsModal.value = false
}

// Packages & Subscriptions plans
interface PlanEntry {
  id: string
  name: string
  planType: 'subscription' | 'package'
  price: string
  billingCycle: 'monthly' | 'yearly'
  sessionCount: string  // '' = unlimited
  validityDays: string
}

const plans = ref<PlanEntry[]>([])
const showPlanModal = ref(false)
const planModalStep = ref<1 | 2>(1)
const editingPlanId = ref<string | null>(null)
const newPlan = ref({
  name: '',
  planType: null as 'subscription' | 'package' | null,
  price: '',
  billingCycle: 'monthly' as 'monthly' | 'yearly',
  sessionCount: '',
  validityDays: '',
})

const openCreatePlan = () => {
  editingPlanId.value = null
  newPlan.value = { name: '', planType: null, price: '', billingCycle: 'monthly', sessionCount: '', validityDays: '' }
  planModalStep.value = 1
  showPlanModal.value = true
}

const openEditPlan = (plan: PlanEntry) => {
  editingPlanId.value = plan.id
  newPlan.value = { name: plan.name, planType: plan.planType, price: plan.price, billingCycle: plan.billingCycle, sessionCount: plan.sessionCount, validityDays: plan.validityDays }
  planModalStep.value = 2
  showPlanModal.value = true
}

const selectPlanType = (type: 'subscription' | 'package') => {
  newPlan.value.planType = type
  planModalStep.value = 2
}

const savePlan = () => {
  if (!newPlan.value.name.trim() || !newPlan.value.planType) return
  if (editingPlanId.value) {
    const idx = plans.value.findIndex(p => p.id === editingPlanId.value)
    if (idx >= 0) plans.value[idx] = { id: editingPlanId.value, name: newPlan.value.name, planType: newPlan.value.planType!, price: newPlan.value.price, billingCycle: newPlan.value.billingCycle, sessionCount: newPlan.value.sessionCount, validityDays: newPlan.value.validityDays }
  } else {
    plans.value.push({ id: String(Date.now()), name: newPlan.value.name, planType: newPlan.value.planType!, price: newPlan.value.price, billingCycle: newPlan.value.billingCycle, sessionCount: newPlan.value.sessionCount, validityDays: newPlan.value.validityDays })
  }
  showPlanModal.value = false
}

const removePlan = (id: string) => {
  plans.value = plans.value.filter(p => p.id !== id)
}

// ── Location ───────────────────────────────────────────────────
type LocationType = 'zoom' | 'phone' | 'in-person' | 'google-meet' | 'ms-teams'
interface LocationEntry { id: string; type: LocationType; address?: string; phoneMethod?: 'require' | 'provide'; phone?: string }

const selectedLocations = ref<LocationEntry[]>([])
const showAllOptionsDropdown = ref(false)
const zoomConnected = ref(false)

const locationLabels: Record<LocationType, string> = {
  zoom: 'Zoom',
  phone: 'Phone Call',
  'in-person': 'In person',
  'google-meet': 'Google Meet',
  'ms-teams': 'Microsoft Teams',
}

const hasLocationType = (type: LocationType) => selectedLocations.value.some(l => l.type === type)

const addLocation = (type: LocationType) => {
  if (hasLocationType(type)) return
  showAllOptionsDropdown.value = false
  selectedLocations.value.push({
    id: String(Date.now()),
    type,
    address: '',
    phoneMethod: 'require',
    phone: '',
  })
}
const removeLocation = (id: string) => {
  selectedLocations.value = selectedLocations.value.filter(l => l.id !== id)
}

const locationDropdownRef = ref<HTMLElement | null>(null)
const handleLocationDropdownClickOutside = (e: MouseEvent) => {
  if (locationDropdownRef.value && !locationDropdownRef.value.contains(e.target as Node)) {
    showAllOptionsDropdown.value = false
  }
}

// ── Images ─────────────────────────────────────────────────────
const MAX_MEDIA = 10

interface MediaItem {
  kind: 'existing' | 'new'
  id: string
  url: string
  file?: File
  pendingDelete?: boolean
}

const allMedia = ref<MediaItem[]>([])
const isDraggingUpload = ref(false)
const fileInputRef = ref<HTMLInputElement | null>(null)
const draggingMediaId = ref<string | null>(null)

const orderedMedia = computed(() => allMedia.value.filter(m => !m.pendingDelete))
const totalMediaCount = computed(() => orderedMedia.value.length)

const addFiles = (files: FileList | File[]) => {
  const remaining = MAX_MEDIA - totalMediaCount.value
  Array.from(files).slice(0, remaining).forEach(file => {
    allMedia.value.push({
      kind: 'new',
      id: `new-${Date.now()}-${Math.random().toString(36).slice(2)}`,
      url: URL.createObjectURL(file),
      file,
    })
  })
}

const onFileInputChange = (e: Event) => {
  const input = e.target as HTMLInputElement
  if (input.files) addFiles(input.files)
  input.value = ''
}

const onDropZoneDrop = (e: DragEvent) => {
  isDraggingUpload.value = false
  if (e.dataTransfer?.files) addFiles(e.dataTransfer.files)
}

const markMediaForDelete = (item: MediaItem) => {
  if (item.kind === 'existing') {
    item.pendingDelete = true
  } else {
    URL.revokeObjectURL(item.url)
    allMedia.value = allMedia.value.filter(m => m.id !== item.id)
  }
}

const setAsFeatured = (item: MediaItem) => {
  const idx = allMedia.value.findIndex(m => m.id === item.id)
  if (idx <= 0) return
  const items = [...allMedia.value]
  const [moved] = items.splice(idx, 1)
  items.unshift(moved)
  allMedia.value = items
}

const onMediaDrop = (targetItem: MediaItem) => {
  if (!draggingMediaId.value || draggingMediaId.value === targetItem.id) {
    draggingMediaId.value = null
    return
  }
  const fromIdx = allMedia.value.findIndex(m => m.id === draggingMediaId.value)
  const toIdx = allMedia.value.findIndex(m => m.id === targetItem.id)
  if (fromIdx === -1 || toIdx === -1) { draggingMediaId.value = null; return }
  const items = [...allMedia.value]
  const [moved] = items.splice(fromIdx, 1)
  items.splice(toIdx, 0, moved)
  allMedia.value = items
  draggingMediaId.value = null
}

// ── Booking ────────────────────────────────────────────────────
const bookingMode = ref<'auto' | 'manual'>('auto')

// ── Tags ───────────────────────────────────────────────────────
const selectedTagIds = ref<string[]>([])

// ── Computed ───────────────────────────────────────────────────
const showSchedule = computed(() => serviceType.value === 'class' || serviceType.value === 'course')
const priceSectionLabel = computed(() => serviceType.value === 'course' ? 'Course price' : 'Price per session')
const paymentModes = [
  { key: 'per-session' as PaymentMode, label: 'Per session', desc: 'Clients pay for sessions based on the price you set' },
  { key: 'with-plan' as PaymentMode, label: 'With a plan', desc: 'Clients buy a membership or package to book sessions with' },
  { key: 'per-session-or-plan' as PaymentMode, label: 'Per session or with a plan', desc: 'Clients pay for sessions based on the price you set' },
]

// ── Sentinel / fixed footer ─────────────────────────────────────
const footerSentinel = ref<HTMLElement | null>(null)
const showFixedBar = ref(false)
const isSaving = ref(false)

onMounted(async () => {
  document.addEventListener('click', handleLocationDropdownClickOutside, true)

  const setupObserver = async () => {
    const observer = new IntersectionObserver(
      (entries) => { showFixedBar.value = !entries[0]?.isIntersecting },
      { threshold: 0 },
    )
    if (footerSentinel.value) observer.observe(footerSentinel.value)
    onBeforeUnmount(() => {
      observer.disconnect()
      document.removeEventListener('click', handleLocationDropdownClickOutside, true)
      allMedia.value.filter(m => m.kind === 'new').forEach(m => URL.revokeObjectURL(m.url))
    })
  }

  if (isNew) {
    await setupObserver()
    return
  }

  const result = await useServiceAPI().GetService(serviceId)
  if (!result.success || !result.data?.data) {
    notify('Service not found', 'error')
    router.push('/catalog/services')
    return
  }

  const s = result.data.data
  const detail = s.service_detail

  form.value = {
    name: s.name,
    tagline: detail?.tagline ?? '',
    description: s.description,
    duration: detail?.duration ?? '1 Hour',
    bufferTime: detail?.buffer_time ?? '10 minutes',
    status: s.status ? s.status.charAt(0).toUpperCase() + s.status.slice(1) : 'Active',
    showInSite: s.show_in_store,
  }

  if (detail?.service_type) {
    serviceType.value = detail.service_type as 'appointment' | 'class' | 'course'
  }
  paymentMode.value = (detail?.payment_mode as PaymentMode) || 'per-session'
  priceType.value = (detail?.price_type as PriceType) || 'fixed'
  priceAmount.value = s.price ? String(s.price) : ''
  priceCurrency.value = s.price_currency || 'USD'
  customPriceLabel.value = detail?.custom_price_label ?? ''
  bookingMode.value = (detail?.booking_mode as 'auto' | 'manual') || 'auto'

  selectedTagIds.value = (s.tags ?? []).map((t: any) => typeof t === 'string' ? t : t.id)

  if (detail?.variations?.length) {
    variations.value = detail.variations.map(v => ({
      id: v.id,
      name: v.name,
      price: v.price ? String(v.price) : '',
    }))
  }

  if (detail?.locations?.length) {
    selectedLocations.value = detail.locations.map(l => ({
      id: l.id,
      type: l.location_type as LocationType,
      address: l.address,
      phoneMethod: (l.phone_method as 'require' | 'provide') || 'require',
      phone: l.phone,
    }))
  }

  if (s.plans?.length) {
    plans.value = s.plans.map(p => ({
      id: p.id,
      name: p.name,
      planType: p.plan_type,
      price: p.price ? String(p.price) : '',
      billingCycle: (p.billing_cycle as 'monthly' | 'yearly') || 'monthly',
      sessionCount: p.session_count != null ? String(p.session_count) : '',
      validityDays: p.validity_days != null ? String(p.validity_days) : '',
    }))
  }

  const rawMedia = s.media ?? []
  const sorted = [...rawMedia].sort((a, b) => {
    if (a.display_image && !b.display_image) return -1
    if (!a.display_image && b.display_image) return 1
    return a.position - b.position
  })
  allMedia.value = sorted.map(m => ({
    kind: 'existing' as const,
    id: m.id,
    url: m.url,
    pendingDelete: false,
  }))

  isLoading.value = false
  await setupObserver()
})

watch(
  () => form.value.status,
  (status) => {
    if (status.toLowerCase() !== 'active') form.value.showInSite = false
  }
)

// ── Actions ────────────────────────────────────────────────────
const buildPayload = (status: string, mediaKeys: string[] = [], displayImageId?: string, displayImageKey?: string) => ({
  name: form.value.name,
  description: form.value.description,
  price: priceType.value === 'free' ? 0 : (parseFloat(priceAmount.value) || 0),
  price_currency: priceCurrency.value,
  status,
  show_in_store: form.value.showInSite,
  service_type: serviceType.value,
  tagline: form.value.tagline,
  duration: form.value.duration,
  buffer_time: form.value.bufferTime,
  price_type: priceType.value,
  payment_mode: paymentMode.value,
  custom_price_label: customPriceLabel.value,
  booking_mode: bookingMode.value,
  tag_ids: selectedTagIds.value,
  variations: priceType.value === 'variable'
    ? variations.value.map((v, i) => ({ name: v.name, price: parseFloat(v.price) || 0, position: i }))
    : [],
  plans: plans.value.map(p => ({
    name: p.name,
    plan_type: p.planType,
    price: parseFloat(p.price) || 0,
    price_currency: priceCurrency.value,
    billing_cycle: p.planType === 'subscription' ? p.billingCycle : undefined,
    session_count: p.sessionCount ? parseInt(p.sessionCount) : null,
    validity_days: p.planType === 'package' && p.validityDays ? parseInt(p.validityDays) : null,
  })),
  locations: selectedLocations.value.map(l => ({
    location_type: l.type,
    address: l.address || '',
    phone_method: l.phoneMethod || '',
    phone: l.phone || '',
  })),
  additional_fields: [],
  media_keys: mediaKeys,
  display_image_id: displayImageId,
  display_image_key: displayImageKey,
  deleted_media_ids: allMedia.value.filter(m => m.kind === 'existing' && m.pendingDelete).map(m => m.id),
})

const uploadAndGetMediaArgs = async () => {
  const newItems = orderedMedia.value.filter(m => m.kind === 'new')
  let mediaKeys: string[] = []
  if (newItems.length > 0) {
    const { keys, error } = await useCatalogMediaAPI().presignAndUpload(newItems.map(m => m.file!))
    if (error) { notify(error, 'error'); return null }
    mediaKeys = keys
  }
  const firstItem = orderedMedia.value[0]
  const displayImageId = firstItem?.kind === 'existing' ? firstItem.id : undefined
  const newFeaturedIdx = newItems.findIndex(m => m.id === firstItem?.id)
  const displayImageKey = newFeaturedIdx >= 0 ? mediaKeys[newFeaturedIdx] : undefined
  return { mediaKeys, displayImageId, displayImageKey }
}

const saveService = async () => {
  if (!form.value.name.trim()) {
    notify('Service name is required', 'error')
    return
  }
  isSaving.value = true
  try {
    const media = await uploadAndGetMediaArgs()
    if (!media) return
    const payload = buildPayload('active', media.mediaKeys, media.displayImageId, media.displayImageKey)
    const result = isNew
      ? await useServiceAPI().CreateService(payload)
      : await useServiceAPI().UpdateService(serviceId, payload)
    if (result.success) {
      notify(isNew ? 'Service created successfully' : 'Service updated successfully', 'success')
      router.push('/catalog/services')
    } else {
      notify(result.error || (isNew ? 'Failed to create service' : 'Failed to update service'), 'error')
    }
  } finally {
    isSaving.value = false
  }
}

const saveAsDraft = async () => {
  if (!form.value.name.trim()) {
    notify('Service name is required', 'error')
    return
  }
  isSaving.value = true
  try {
    const media = await uploadAndGetMediaArgs()
    if (!media) return
    const payload = { ...buildPayload('draft', media.mediaKeys, media.displayImageId, media.displayImageKey), show_in_store: false }
    const result = isNew
      ? await useServiceAPI().CreateService(payload)
      : await useServiceAPI().UpdateService(serviceId, payload)
    if (result.success) {
      notify('Draft saved', 'success')
      router.push('/catalog/services')
    } else {
      notify(result.error || 'Failed to save draft', 'error')
    }
  } finally {
    isSaving.value = false
  }
}

const discard = () => router.push('/catalog/services')
</script>

<template>
  <MainDashboard current="Services">
    <template #title>
      <BackLink
        to="/catalog/services"
        label="Services"
      />
    </template>

    <template #body>
      <div
        v-if="isLoading"
        class="flex items-center justify-center py-24"
      >
        <div class="w-8 h-8 border-4 border-[#2176AE] border-t-transparent rounded-full animate-spin" />
      </div>

      <div
        v-else
        class="rounded-3xl gap-6 mt-4"
        :class="{ 'bg-white': step !== 1 || !isNew }"
      >
        <!-- Header -->
        <div
          class="flex justify-between pt-8 items-center border-b border-[#F6F6F7]"
          :class="{'px-8' : step !== 1 || !isNew}"
        >
          <h2 class="text-xl font-semibold text-[#1E212B]">
            {{ !isNew ? (form.name || 'Edit service') : (step === 1 ? 'Add a New Service' : (form.name || 'Create new service')) }}
          </h2>
          <button
            v-if="step !== 1 || !isNew"
            class="w-8 h-8 flex items-center justify-center rounded-full hover:bg-[#EDEDEE] text-[#4B4D55] hover:text-[#1E212B] transition-colors"
            @click="isNew && step === 2 ? step = 1 : discard()"
          >
            <svg
              class="w-6 h-6"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <path d="M18 6L6 18M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- ── Step 1: Type picker (new only) ──────────────────────── -->
        <div
          v-if="isNew && step === 1"
          class="pb-8"
        >
          <p class="text-sm text-[#6F7177] mb-10">
            Edit a service below or start adding one from scratch.
          </p>
          <div class="grid grid-cols-3 gap-6">
            <button
              v-for="t in serviceTypes"
              :key="t.key"
              class="group col-span-3 md:col-span-1 bg-white rounded-2xl border border-[#EDEDEE] overflow-hidden text-left hover:border-[#94BDD8] hover:shadow-md transition-all"
              @click="selectType(t.key)"
            >
              <div class="aspect-[4/3] bg-[#EDEDEE] group-hover:bg-[#E0EDF5] transition-colors" />
              <div class="px-6 py-5">
                <p class="text-base font-semibold text-[#1E212B] text-center mb-1">
                  {{ t.label }}
                </p>
                <p class="text-sm text-[#6F7177] text-center leading-snug">
                  {{ t.desc }}
                </p>
              </div>
            </button>
          </div>
        </div>

        <!-- ── Step 2: Full form ────────────────────────────────── -->
        <div
          v-else
          class="grid grid-cols-5 gap-y-6 md:border-b md:border-[#EDEDEE]"
        >
          <!-- ── Left column ──────────────────────────────────── -->
          <div class="col-span-5 md:col-span-3 pt-3 md:border-r md:border-[#EDEDEE] md:border-dashed gap-4">
            <!-- Service Type -->
            <GroFormSection
              title="Service Type"
              :open="open.serviceType"
              @toggle="toggle('serviceType')"
            >
              <GroBasicRadio
                v-model="serviceType"
                name="serviceType"
                layout="vertical"
                :options="[
                  { value: 'appointment', label: 'Appointment', description: 'Private session that can be booked according to availability' },
                  { value: 'class', label: 'Class', description: 'A group session that can recur. Clients book any session they want to join.' },
                  { value: 'course', label: 'Course', description: 'A set of group sessions. Clients book them all up front' },
                ]"
              />
            </GroFormSection>

            <!-- Service Details -->
            <GroFormSection
              title="Service Details"
              :open="open.details"
              @toggle="toggle('details')"
            >
              <div class="flex flex-col gap-4">
                <!-- Name + Image -->
                <div class="grid grid-cols-2 gap-4 items-start">
                  <GroBasicInput
                    v-model="form.name"
                    placeholder="Service"
                  >
                    Name
                  </GroBasicInput>
                  <GroBasicInput
                    v-model="form.tagline"
                    placeholder="Add text"
                  >
                    Tagline (Optional)
                  </GroBasicInput>
                </div>

                <!-- Description -->
                <div>
                  <div class="flex justify-between items-center mt-2 mb-1">
                    <label class="text-xs font-medium text-[#1E212B]">Description</label>
                    <button class="text-xs text-[#2176AE] flex items-center gap-1 hover:underline">
                      <svg
                        class="w-3.5 h-3.5"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="1.5"
                      >
                        <path d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
                      </svg>
                      Generate AI Text
                    </button>
                  </div>
                  <div class="border rounded-2xl border-gray-200">
                    <CommentEditor
                      v-model="form.description"
                      placeholder="For SEO improvement and help customers find this item by adding a unique & detailed service description."
                    />
                  </div>
                </div>

                <!-- Timing (Appointment only) -->
                <div v-if="serviceType === 'appointment'">
                  <p class="text-[10px] font-semibold text-[#6F7177] uppercase tracking-widest mb-3">
                    Timing
                  </p>
                  <div class="grid grid-cols-2 gap-4">
                    <GroBasicSelect
                      v-model="form.duration"
                      :options="durationOptions"
                    >
                      Duration
                    </GroBasicSelect>
                    <GroBasicSelect
                      v-model="form.bufferTime"
                      :options="bufferOptions"
                    >
                      Buffer time
                    </GroBasicSelect>
                  </div>
                </div>
              </div>
            </GroFormSection>

            <!-- Price & Payments -->
            <GroFormSection
              title="Price &amp; Payments"
              :open="open.price"
              @toggle="toggle('price')"
            >
              <div class="flex flex-col gap-5">
                <p class="text-sm font-medium text-[#1E212B]">
                  Choose how your clients pay for this service
                </p>

                <!-- Payment mode cards -->
                <div class="grid grid-cols-3 gap-3">
                  <button
                    v-for="mode in paymentModes"
                    :key="mode.key"
                    class="relative p-4 rounded-xl border text-left transition-colors"
                    :class="paymentMode === mode.key
                      ? 'border-[#2176AE] bg-[#F0F7FF]'
                      : 'border-[#EDEDEE] hover:border-[#94BDD8]'"
                    @click="paymentMode = mode.key"
                  >
                    <div
                      v-if="paymentMode === mode.key"
                      class="absolute top-2 right-2 w-4 h-4 rounded-full bg-[#2176AE] flex items-center justify-center"
                    >
                      <svg
                        class="w-2.5 h-2.5 text-white"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="3"
                      >
                        <path d="M5 13l4 4L19 7" />
                      </svg>
                    </div>
                    <p class="text-sm font-semibold text-[#1E212B] mb-1">
                      {{ mode.label }}
                    </p>
                    <p class="text-xs text-[#6F7177] leading-snug">
                      {{ mode.desc }}
                    </p>
                  </button>
                </div>

                <!-- Price per session / Course price -->
                <div v-if="paymentMode === 'per-session' || paymentMode === 'per-session-or-plan'">
                  <p class="text-sm font-semibold text-[#1E212B] mb-3">
                    {{ priceSectionLabel }}
                  </p>
                  <div class="grid grid-cols-2 gap-4">
                    <GroBasicSelect
                      v-model="priceType"
                      :options="priceTypeOptions"
                    >
                      Price type
                    </GroBasicSelect>
                    <GroCurrencyInput
                      v-if="priceType === 'fixed' || priceType === 'variable'"
                      v-model:amount="priceAmount"
                      v-model:currency="priceCurrency"
                      label="Amount"
                      placeholder="0"
                    />
                    <div v-else-if="priceType === 'custom'">
                      <GroBasicInput
                        v-model="customPriceLabel"
                        placeholder="e.g., Starting from ₦500"
                      >
                        Custom Price
                      </GroBasicInput>
                    </div>
                  </div>

                  <!-- Variable Pricing: Variations table -->
                  <div
                    v-if="priceType === 'variable'"
                    class="mt-4"
                  >
                    <!-- Empty state -->
                    <div
                      v-if="variations.length === 0"
                      class="border border-[#EDEDEE] rounded-xl p-6 flex flex-col items-center gap-2 text-center"
                    >
                      <p class="text-sm font-semibold text-[#1E212B]">
                        Add variations to your service
                      </p>
                      <p class="text-sm text-[#6F7177]">
                        Variations can be based on service features, duration, staff, and more.
                      </p>
                      <GroBasicButton
                        color="tertiary"
                        size="sm"
                        class="mt-1 !w-auto"
                        @click="openVariationsModal"
                      >
                        <template #default>
                          + Add Variations
                        </template>
                      </GroBasicButton>
                    </div>

                    <!-- Populated table -->
                    <template v-else>
                      <div class="border border-[#EDEDEE] rounded-xl overflow-hidden">
                        <table class="w-full">
                          <thead>
                            <tr class="bg-[#F6F6F7] border-b border-[#EDEDEE]">
                              <th class="px-4 py-2.5 text-left text-xs font-semibold text-[#1E212B]">
                                Variation Name
                              </th>
                              <th class="px-4 py-2.5 text-left text-xs font-semibold text-[#1E212B]">
                                Price
                              </th>
                              <th class="w-10 px-4 py-2.5" />
                            </tr>
                          </thead>
                          <tbody>
                            <tr
                              v-for="v in variations"
                              :key="v.id"
                              class="border-b border-[#EDEDEE] last:border-0"
                            >
                              <td class="px-4 py-3 text-sm text-[#1E212B]">
                                {{ v.name || '--' }}
                              </td>
                              <td class="px-4 py-3 text-sm text-[#1E212B]">
                                ₦ {{ parseFloat(v.price).toLocaleString('en-NG', { minimumFractionDigits: 2 }) }}
                              </td>
                              <td class="px-4 py-3">
                                <button
                                  class="w-7 h-7 flex items-center justify-center rounded hover:bg-[#EDEDEE] text-[#6F7177] transition-colors"
                                  @click="openVariationsModal"
                                >
                                  <svg
                                    class="w-3.5 h-3.5"
                                    viewBox="0 0 24 24"
                                    fill="none"
                                    stroke="currentColor"
                                    stroke-width="2"
                                  >
                                    <path d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z" />
                                  </svg>
                                </button>
                              </td>
                            </tr>
                          </tbody>
                        </table>
                      </div>
                      <button
                        class="mt-2 text-xs font-medium text-[#2176AE] hover:underline"
                        @click="openVariationsModal"
                      >
                        Manage variations
                      </button>
                    </template>
                  </div>
                </div>

                <!-- Packages & Subscriptions -->
                <div
                  v-if="paymentMode === 'with-plan' || paymentMode === 'per-session-or-plan'"
                  class="border-t border-[#EDEDEE] pt-4"
                >
                  <p class="text-sm font-semibold text-[#1E212B] mb-3">
                    Packages &amp; Subscriptions
                  </p>

                  <!-- Plans list -->
                  <div
                    v-if="plans.length > 0"
                    class="flex flex-col gap-2 mb-3"
                  >
                    <div
                      v-for="plan in plans"
                      :key="plan.id"
                      class="flex items-center justify-between px-4 py-3 border border-[#EDEDEE] rounded-xl"
                    >
                      <div class="flex items-center gap-3">
                        <div
                          class="px-2 py-0.5 rounded-full text-[10px] font-semibold uppercase tracking-wide"
                          :class="plan.planType === 'subscription' ? 'bg-[#EEF2FF] text-[#4F46E5]' : 'bg-[#F0FDF4] text-[#16A34A]'"
                        >
                          {{ plan.planType }}
                        </div>
                        <div>
                          <p class="text-sm font-medium text-[#1E212B]">
                            {{ plan.name }}
                          </p>
                          <p class="text-xs text-[#6F7177]">
                            <template v-if="plan.planType === 'subscription'">
                              {{ plan.billingCycle === 'monthly' ? 'Monthly' : 'Yearly' }} ·
                              {{ plan.sessionCount ? `${plan.sessionCount} sessions/cycle` : 'Unlimited sessions' }}
                            </template>
                            <template v-else>
                              {{ plan.sessionCount ? `${plan.sessionCount} sessions` : 'Unlimited sessions' }}
                              <template v-if="plan.validityDays"> · {{ plan.validityDays }} days</template>
                            </template>
                          </p>
                        </div>
                      </div>
                      <div class="flex items-center gap-3">
                        <span class="text-sm font-semibold text-[#1E212B]">
                          {{ priceCurrency }} {{ parseFloat(plan.price || '0').toLocaleString() }}
                        </span>
                        <button
                          class="w-7 h-7 flex items-center justify-center rounded hover:bg-[#EDEDEE] text-[#6F7177] transition-colors"
                          @click="openEditPlan(plan)"
                        >
                          <svg
                            class="w-3.5 h-3.5"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2"
                          >
                            <path d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z" />
                          </svg>
                        </button>
                        <button
                          class="w-7 h-7 flex items-center justify-center rounded hover:bg-[#FEE2E2] text-[#6F7177] hover:text-[#AF513A] transition-colors"
                          @click="removePlan(plan.id)"
                        >
                          <svg
                            class="w-3.5 h-3.5"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="1.5"
                          >
                            <path d="M3 6h18M8 6V4h8v2M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6" />
                          </svg>
                        </button>
                      </div>
                    </div>
                  </div>

                  <!-- Empty state -->
                  <div
                    v-else
                    class="border border-dashed border-[#DBDBDD] rounded-xl p-6 flex flex-col items-center gap-3 text-center"
                  >
                    <p class="text-sm font-medium text-[#1E212B]">
                      Create your first plan
                    </p>
                    <p class="text-xs text-[#6F7177]">
                      Let your clients choose from multiple plans.
                    </p>
                  </div>

                  <button
                    class="flex items-center gap-1.5 text-sm font-medium text-[#2176AE] hover:underline mt-3"
                    @click="openCreatePlan"
                  >
                    <svg
                      class="w-4 h-4"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="2.5"
                    >
                      <path d="M12 5v14M5 12h14" />
                    </svg>
                    {{ plans.length > 0 ? 'Add another plan' : 'Create New Plan' }}
                  </button>
                </div>
              </div>
            </GroFormSection>

            <!-- Schedule (Class / Course only) -->
            <GroFormSection
              v-if="showSchedule"
              title="Schedule"
              :open="open.schedule"
              @toggle="toggle('schedule')"
            >
              <div class="flex flex-col items-center justify-center py-8 gap-3 text-center">
                <div class="w-10 h-10 rounded-full bg-[#F6F6F7] flex items-center justify-center">
                  <svg
                    class="w-5 h-5 text-[#939499]"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="1.5"
                  >
                    <path d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </div>
                <p class="text-sm font-semibold text-[#1E212B]">
                  Coming Soon
                </p>
                <p class="text-xs text-[#6F7177]">
                  Session scheduling will be available in a future update.
                </p>
              </div>
            </GroFormSection>

            <!-- Location -->
            <GroFormSection
              title="Location"
              :open="open.location"
              @toggle="toggle('location')"
            >
              <div class="flex flex-col gap-4">
                <!-- Type picker -->
                <div class="flex gap-2 flex-wrap">
                  <button
                    :disabled="hasLocationType('zoom')"
                    class="flex flex-col items-center gap-1.5 px-5 py-3 border rounded-xl text-sm transition-colors"
                    :class="hasLocationType('zoom') ? 'border-[#EDEDEE] text-[#C0C1C6] bg-[#F6F6F7] cursor-not-allowed' : 'border-[#EDEDEE] text-[#1E212B] hover:border-[#94BDD8] hover:bg-[#F0F7FF]'"
                    @click="addLocation('zoom')"
                  >
                    <svg
                      class="w-5 h-5"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="1.5"
                    >
                      <path d="M15 10l4.553-2.553A1 1 0 0121 8.382v7.236a1 1 0 01-1.447.894L15 14M5 8h10a2 2 0 012 2v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4a2 2 0 012-2z" />
                    </svg>
                    <span class="text-xs">Zoom</span>
                  </button>
                  <button
                    :disabled="hasLocationType('phone')"
                    class="flex flex-col items-center gap-1.5 px-5 py-3 border rounded-xl text-sm transition-colors"
                    :class="hasLocationType('phone') ? 'border-[#EDEDEE] text-[#C0C1C6] bg-[#F6F6F7] cursor-not-allowed' : 'border-[#EDEDEE] text-[#1E212B] hover:border-[#94BDD8] hover:bg-[#F0F7FF]'"
                    @click="addLocation('phone')"
                  >
                    <svg
                      class="w-5 h-5"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="1.5"
                    >
                      <path d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z" />
                    </svg>
                    <span class="text-xs">Phone</span>
                  </button>
                  <button
                    :disabled="hasLocationType('in-person')"
                    class="flex flex-col items-center gap-1.5 px-5 py-3 border rounded-xl text-sm transition-colors"
                    :class="hasLocationType('in-person') ? 'border-[#EDEDEE] text-[#C0C1C6] bg-[#F6F6F7] cursor-not-allowed' : 'border-[#EDEDEE] text-[#1E212B] hover:border-[#94BDD8] hover:bg-[#F0F7FF]'"
                    @click="addLocation('in-person')"
                  >
                    <svg
                      class="w-5 h-5"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="1.5"
                    >
                      <path d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                      <path d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                    </svg>
                    <span class="text-xs">In person</span>
                  </button>
                  <div
                    ref="locationDropdownRef"
                    class="relative"
                  >
                    <button
                      class="flex flex-col items-center gap-1.5 px-5 py-3 border border-[#EDEDEE] rounded-xl text-sm text-[#1E212B] hover:border-[#94BDD8] hover:bg-[#F0F7FF] transition-colors"
                      @click.stop="showAllOptionsDropdown = !showAllOptionsDropdown"
                    >
                      <svg
                        class="w-5 h-5"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="1.5"
                      >
                        <path d="M6 9l6 6 6-6" />
                      </svg>
                      <span class="text-xs">All options</span>
                    </button>
                    <div
                      v-if="showAllOptionsDropdown"
                      class="absolute top-full left-0 mt-1 bg-white border border-[#EDEDEE] rounded-xl shadow-lg z-20 min-w-44 py-1 overflow-hidden"
                    >
                      <button
                        :disabled="hasLocationType('google-meet')"
                        class="w-full text-left px-4 py-2.5 text-sm flex items-center gap-2.5 transition-colors"
                        :class="hasLocationType('google-meet') ? 'text-[#C0C1C6] cursor-not-allowed' : 'text-[#1E212B] hover:bg-[#F6F6F7]'"
                        @click="addLocation('google-meet')"
                      >
                        <svg
                          class="w-4 h-4 shrink-0"
                          viewBox="0 0 24 24"
                          fill="none"
                        >
                          <rect
                            width="24"
                            height="24"
                            rx="4"
                            fill="#00832D"
                          />
                          <path
                            d="M15 8.5V6l4 4-4 4v-2.5H9V8.5h6z"
                            fill="white"
                          />
                          <rect
                            x="5"
                            y="8.5"
                            width="10"
                            height="7"
                            rx="1"
                            fill="white"
                          />
                        </svg>
                        Google Meet
                      </button>
                      <button
                        :disabled="hasLocationType('ms-teams')"
                        class="w-full text-left px-4 py-2.5 text-sm flex items-center gap-2.5 transition-colors"
                        :class="hasLocationType('ms-teams') ? 'text-[#C0C1C6] cursor-not-allowed' : 'text-[#1E212B] hover:bg-[#F6F6F7]'"
                        @click="addLocation('ms-teams')"
                      >
                        <svg
                          class="w-4 h-4 shrink-0"
                          viewBox="0 0 24 24"
                          fill="none"
                        >
                          <rect
                            width="24"
                            height="24"
                            rx="4"
                            fill="#5059C9"
                          />
                          <path
                            d="M13 7h4a1 1 0 011 1v6a1 1 0 01-1 1h-1v-5h-3V7z"
                            fill="white"
                          />
                          <rect
                            x="6"
                            y="9"
                            width="8"
                            height="8"
                            rx="1.5"
                            fill="white"
                          />
                        </svg>
                        Microsoft Teams
                      </button>
                    </div>
                  </div>
                </div>

                <!-- Selected location entries -->
                <div
                  v-for="loc in selectedLocations"
                  :key="loc.id"
                  class="border border-[#EDEDEE] rounded-xl overflow-hidden"
                >
                  <div class="flex items-center justify-between px-4 py-3 bg-[#F6F6F7]">
                    <div class="flex items-center gap-2">
                      <div class="w-6 h-6 rounded-full bg-[#2176AE] flex items-center justify-center">
                        <svg
                          class="w-3 h-3 text-white"
                          viewBox="0 0 24 24"
                          fill="currentColor"
                        >
                          <circle
                            cx="12"
                            cy="12"
                            r="5"
                          />
                        </svg>
                      </div>
                      <span class="text-sm font-medium text-[#1E212B]">{{ locationLabels[loc.type] }}</span>
                    </div>
                    <button
                      class="text-[#939499] hover:text-[#1E212B] transition-colors"
                      @click="removeLocation(loc.id)"
                    >
                      <svg
                        class="w-4 h-4"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                      >
                        <path d="M18 6L6 18M6 6l12 12" />
                      </svg>
                    </button>
                  </div>
                  <div
                    v-if="loc.type === 'zoom'"
                    class="px-4 py-3"
                  >
                    <div
                      v-if="!zoomConnected"
                      class="flex items-center gap-3 px-3 py-2.5 bg-[#FFF4ED] border border-[#F59E0B]/30 rounded-lg"
                    >
                      <svg
                        class="w-4 h-4 text-[#D26B06] shrink-0"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                      >
                        <path d="M12 9v4M12 17h.01M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z" />
                      </svg>
                      <span class="text-xs text-[#D26B06] flex-1">Your Zoom account is not connected</span>
                      <button class="text-xs font-medium text-[#2176AE] flex items-center gap-1 hover:underline whitespace-nowrap">
                        Connect Zoom
                        <svg
                          class="w-3 h-3"
                          viewBox="0 0 24 24"
                          fill="none"
                          stroke="currentColor"
                          stroke-width="2"
                        >
                          <path d="M18 13v6a2 2 0 01-2 2H5a2 2 0 01-2-2V8a2 2 0 012-2h6M15 3h6v6M10 14L21 3" />
                        </svg>
                      </button>
                    </div>
                  </div>
                  <div
                    v-else-if="loc.type === 'phone'"
                    class="px-4 py-3 flex flex-col gap-2.5"
                  >
                    <p class="text-xs font-medium text-[#1E212B]">
                      How will you get in touch?
                    </p>
                    <label class="flex items-center gap-2 cursor-pointer">
                      <div
                        class="w-3.5 h-3.5 rounded-full border-2 flex items-center justify-center transition-colors"
                        :class="loc.phoneMethod === 'require' ? 'border-[#1E212B]' : 'border-[#DBDBDD]'"
                      >
                        <div
                          v-if="loc.phoneMethod === 'require'"
                          class="w-1.5 h-1.5 rounded-full bg-[#1E212B]"
                        />
                      </div>
                      <span
                        class="text-sm text-[#1E212B]"
                        @click="loc.phoneMethod = 'require'"
                      >Require invitee's phone number.</span>
                    </label>
                    <label class="flex items-start gap-2 cursor-pointer">
                      <div
                        class="w-3.5 h-3.5 rounded-full border-2 flex items-center justify-center mt-0.5 shrink-0 transition-colors"
                        :class="loc.phoneMethod === 'provide' ? 'border-[#1E212B]' : 'border-[#DBDBDD]'"
                        @click="loc.phoneMethod = 'provide'"
                      >
                        <div
                          v-if="loc.phoneMethod === 'provide'"
                          class="w-1.5 h-1.5 rounded-full bg-[#1E212B]"
                        />
                      </div>
                      <div @click="loc.phoneMethod = 'provide'">
                        <p class="text-sm text-[#1E212B]">Provide a phone number to invitees after they book.</p>
                        <div
                          v-if="loc.phoneMethod === 'provide'"
                          class="mt-2 flex items-center bg-[#F6F6F7] border border-[#EDEDEE] rounded-lg hover:border-[#94BDD8] focus-within:border-[#1E212B] transition-colors"
                        >
                          <span class="pl-3 text-sm text-[#4B4D55]">🇳🇬</span>
                          <svg
                            class="w-3 h-3 text-[#939499] ml-1"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2.5"
                          >
                            <path d="M6 9l6 6 6-6" />
                          </svg>
                          <input
                            v-model="loc.phone"
                            type="tel"
                            placeholder="Phone number"
                            class="flex-1 px-2 py-2 bg-transparent text-sm text-[#1E212B] outline-none placeholder-[#939499]"
                          >
                        </div>
                      </div>
                    </label>
                  </div>
                  <div
                    v-else-if="loc.type === 'in-person'"
                    class="px-4 py-3"
                  >
                    <GroBasicInput
                      v-model="loc.address"
                      placeholder="E.g. Earls Barton, 2301 Highland Ave, Lagos Island."
                    >
                      Location name/address
                    </GroBasicInput>
                  </div>
                  <div
                    v-else-if="loc.type === 'google-meet'"
                    class="px-4 py-3"
                  >
                    <div class="flex items-center gap-3 px-3 py-2.5 bg-[#F0FFF4] border border-[#34A853]/30 rounded-lg">
                      <svg
                        class="w-4 h-4 text-[#34A853] shrink-0"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                      >
                        <path d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                      </svg>
                      <span class="text-xs text-[#1E7D3A]">A unique Google Meet link will be generated for each booking.</span>
                    </div>
                  </div>
                  <div
                    v-else-if="loc.type === 'ms-teams'"
                    class="px-4 py-3"
                  >
                    <div class="flex items-center gap-3 px-3 py-2.5 bg-[#EEF2FF] border border-[#5059C9]/30 rounded-lg">
                      <svg
                        class="w-4 h-4 text-[#5059C9] shrink-0"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                      >
                        <path d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                      </svg>
                      <span class="text-xs text-[#3B3F9C]">A unique Microsoft Teams link will be generated for each booking.</span>
                    </div>
                  </div>
                </div>

                <button
                  v-if="selectedLocations.length > 0 && selectedLocations.length < 5"
                  class="flex items-center gap-1.5 text-sm font-medium text-[#2176AE] hover:underline"
                  @click.stop="showAllOptionsDropdown = true"
                >
                  <svg
                    class="w-4 h-4"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2.5"
                  >
                    <path d="M12 5v14M5 12h14" />
                  </svg>
                  Add more items
                </button>
                <p
                  v-if="selectedLocations.length > 0"
                  class="text-xs text-[#6F7177]"
                >
                  Let your invitee choose from multiple meeting locations.
                </p>
              </div>
            </GroFormSection>

            <!-- Images -->
            <GroFormSection
              title="Images"
              :open="open.images"
              @toggle="toggle('images')"
            >
              <div class="flex flex-col gap-4">
                <div class="flex items-center gap-2 px-4 py-3 bg-[#EEF8F5] rounded-xl border border-[#D1F0E8]">
                  <svg
                    class="w-4 h-4 text-[#00916E] shrink-0"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                  >
                    <circle
                      cx="12"
                      cy="12"
                      r="3"
                    /><path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42" />
                  </svg>
                  <p class="text-xs text-[#00916E] font-medium">
                    Services with images get booked more frequently.
                  </p>
                </div>

                <!-- Unified grid: existing + new files, all draggable -->
                <div
                  v-if="totalMediaCount > 0"
                  class="mb-2"
                >
                  <div class="flex items-center gap-3 mb-4">
                    <p class="text-xs font-bold text-[#6F7177] uppercase tracking-widest">
                      Images and Videos
                    </p>
                    <span class="text-xs bg-[#EDEDEE] px-2 py-1 rounded-xl text-[#070707]">{{ totalMediaCount }}/{{ MAX_MEDIA }}</span>
                  </div>
                  <div class="grid grid-cols-4 gap-2">
                    <!-- First image: large (2×2) -->
                    <div
                      class="col-span-2 row-span-2 rounded-2xl overflow-hidden bg-[#F6F6F7] relative group cursor-grab border-2 border-[#1E212B] transition-opacity"
                      :class="draggingMediaId === orderedMedia[0].id ? 'opacity-40' : ''"
                      draggable="true"
                      @dragstart="(e) => { e.dataTransfer?.setData('text/plain', orderedMedia[0].id); draggingMediaId = orderedMedia[0].id }"
                      @dragover.prevent
                      @drop.prevent="onMediaDrop(orderedMedia[0])"
                      @dragend="draggingMediaId = null"
                    >
                      <img
                        :src="orderedMedia[0].url"
                        draggable="false"
                        class="w-full h-full object-cover pointer-events-none"
                        alt=""
                      >
                      <div class="absolute top-2 left-2 bg-black/60 text-white text-[9px] font-semibold px-2 py-0.5 rounded-full flex items-center gap-1 pointer-events-none">
                        <svg
                          class="w-2.5 h-2.5"
                          viewBox="0 0 24 24"
                          fill="currentColor"
                        >
                          <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" />
                        </svg>
                        Featured
                      </div>
                      <div class="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex flex-col items-end justify-start p-1.5">
                        <button
                          type="button"
                          class="w-7 h-7 rounded-full bg-white/90 flex items-center justify-center hover:bg-red-50 transition-colors"
                          @click.stop="markMediaForDelete(orderedMedia[0])"
                        >
                          <svg
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="#AF513A"
                            stroke-width="2"
                            class="w-3.5 h-3.5"
                          >
                            <path d="M6 18L18 6M6 6l12 12" />
                          </svg>
                        </button>
                      </div>
                    </div>
                    <!-- Remaining thumbnails -->
                    <div
                      v-for="item in orderedMedia.slice(1)"
                      :key="item.id"
                      class="aspect-square rounded-xl overflow-hidden bg-[#F6F6F7] relative group cursor-grab border border-[#EDEDEE] transition-opacity"
                      :class="draggingMediaId === item.id ? 'opacity-40' : ''"
                      draggable="true"
                      @dragstart="(e) => { e.dataTransfer?.setData('text/plain', item.id); draggingMediaId = item.id }"
                      @dragover.prevent
                      @drop.prevent="onMediaDrop(item)"
                      @dragend="draggingMediaId = null"
                    >
                      <img
                        :src="item.url"
                        draggable="false"
                        class="w-full h-full object-cover pointer-events-none"
                        alt=""
                      >
                      <div class="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex flex-col items-end justify-start p-1.5 gap-1.5">
                        <button
                          type="button"
                          class="w-7 h-7 rounded-full bg-white/90 flex items-center justify-center hover:bg-red-50 transition-colors"
                          @click.stop="markMediaForDelete(item)"
                        >
                          <svg
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="#AF513A"
                            stroke-width="2"
                            class="w-3.5 h-3.5"
                          >
                            <path d="M6 18L18 6M6 6l12 12" />
                          </svg>
                        </button>
                        <button
                          type="button"
                          class="w-7 h-7 rounded-full bg-white/90 flex items-center justify-center hover:bg-yellow-50 transition-colors"
                          title="Set as featured"
                          @click.stop="setAsFeatured(item)"
                        >
                          <svg
                            class="w-3.5 h-3.5 text-[#1E212B]"
                            viewBox="0 0 24 24"
                            fill="currentColor"
                          >
                            <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" />
                          </svg>
                        </button>
                      </div>
                    </div>
                  </div>
                  <p class="text-xs text-[#939499] mt-3">
                    Drag to reorder · First image is the featured image
                  </p>
                </div>

                <!-- Upload drop zone -->
                <div v-if="totalMediaCount < MAX_MEDIA">
                  <div
                    v-if="totalMediaCount === 0"
                    class="flex items-center gap-3 mb-2"
                  >
                    <p class="text-xs font-bold text-[#6F7177] uppercase tracking-widest">
                      Images and Videos
                    </p>
                  </div>
                  <label
                    class="relative border-2 border-dashed rounded-xl py-10 flex flex-col items-center justify-center gap-2 cursor-pointer transition-colors"
                    :class="isDraggingUpload ? 'border-[#2176AE] bg-[#EBF4FF]' : 'border-[#94BDD8] hover:bg-[#F0F7FF]'"
                    @dragenter.prevent="isDraggingUpload = true"
                    @dragover.prevent="isDraggingUpload = true"
                    @dragleave.prevent="isDraggingUpload = false"
                    @drop.prevent="onDropZoneDrop"
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
                    <span class="text-sm font-semibold text-[#2176AE]">{{ isDraggingUpload ? 'Drop files here' : 'Add Files' }}</span>
                    <span class="text-xs text-[#6F7177]">Accepts images, videos or 3D models</span>
                    <span class="text-xs text-[#939499]">Max 20MB per file · up to {{ MAX_MEDIA }} files</span>
                    <input
                      ref="fileInputRef"
                      type="file"
                      multiple
                      accept="image/*,video/*,.glb,.gltf"
                      class="hidden"
                      @change="onFileInputChange"
                    >
                  </label>
                </div>
              </div>
            </GroFormSection>

            <!-- Booking Requests -->
            <GroFormSection
              title="Booking Requests"
              :open="open.booking"
              @toggle="toggle('booking')"
            >
              <GroBasicRadio
                v-model="bookingMode"
                name="bookingMode"
                layout="vertical"
                :options="[
                  { value: 'auto', label: 'Automatically accept all bookings when available' },
                  { value: 'manual', label: 'Manually approve or decline booking requests' },
                ]"
              />
            </GroFormSection>
          </div>

          <!-- ── Right column ──────────────────────────────────── -->
          <div class="col-span-5 md:col-span-2 pt-3 gap-4 md:pr-8">
            <!-- Service Status -->
            <div class="py-6 border-b border-[#EDEDEE] px-8">
              <h3 class="text-base font-semibold text-[#1E212B] mb-3">
                Service Status
              </h3>
              <select
                v-model="form.status"
                class="w-full px-3 py-2 bg-[#F6F6F7] border border-[#EDEDEE] rounded-lg text-sm text-[#1E212B] outline-none hover:border-[#94BDD8] focus:border-[#1E212B] mb-3 cursor-pointer"
              >
                <option value="Active">
                  Active
                </option>
                <option value="Draft">
                  Draft
                </option>
                <option value="Archived">
                  Archived
                </option>
              </select>
              <GroToggle
                v-model="form.showInSite"
                :disabled="form.status.toLowerCase() !== 'active'"
              >
                Show in your store
              </GroToggle>
              <button class="flex items-center gap-1.5 text-sm font-medium text-[#2176AE] hover:underline mt-3">
                <svg
                  class="w-4 h-4"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                >
                  <path d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  <path d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                </svg>
                Preview
              </button>
            </div>

            <!-- Service Tags -->
            <div class="py-6 border-b border-[#EDEDEE] px-8">
              <h3 class="text-base font-semibold text-[#1E212B] mb-3">
                Tags
              </h3>
              <CatalogTagSelect v-model="selectedTagIds" />
            </div>
          </div>
        </div>

        <!-- Sentinel for conditional fixed footer -->
        <div
          ref="footerSentinel"
          class="h-px"
        />

        <!-- Inline footer (visible when form is short) -->
        <div
          v-if="step == 2 || !isNew"
          class="flex items-center justify-between px-8 py-4 border-t border-[#EDEDEE]"
          :class="showFixedBar ? 'invisible' : ''"
        >
          <button
            class="text-sm font-medium text-[#2176AE] hover:underline disabled:opacity-50"
            :disabled="isSaving"
            @click="saveAsDraft"
          >
            {{ isSaving ? 'Saving...' : 'Save as Draft' }}
          </button>
          <div class="flex items-center gap-3">
            <GroBasicButton
              color="secondary"
              size="xs"
              shape="custom"
              class="w-max"
              @click="discard"
            >
              Discard
            </GroBasicButton>
            <GroBasicButton
              color="primary"
              size="xs"
              shape="custom"
              class="w-max"
              :disabled="isSaving"
              @click="saveService"
            >
              {{ isSaving ? 'Saving...' : 'Save' }}
            </GroBasicButton>
          </div>
        </div>
      </div>

      <!-- Fixed footer (visible when form overflows viewport) -->
      <Transition
        enter-active-class="transition ease-out duration-200"
        enter-from-class="translate-y-full opacity-0"
        enter-to-class="translate-y-0 opacity-100"
        leave-active-class="transition ease-in duration-150"
        leave-from-class="translate-y-0 opacity-100"
        leave-to-class="translate-y-full opacity-0"
      >
        <div
          v-if="showFixedBar && (step === 2 || !isNew)"
          class="fixed bottom-0 left-0 md:left-16 lg:left-64 right-0 bg-white border-t border-[#EDEDEE] px-8 py-4 flex items-center justify-between z-40"
        >
          <button
            class="text-sm font-medium text-[#2176AE] hover:underline disabled:opacity-50"
            :disabled="isSaving"
            @click="saveAsDraft"
          >
            {{ isSaving ? 'Saving...' : 'Save as Draft' }}
          </button>
          <div class="flex items-center gap-3">
            <GroBasicButton
              color="secondary"
              size="xs"
              shape="custom"
              class="w-max"
              @click="discard"
            >
              Discard
            </GroBasicButton>
            <GroBasicButton
              color="primary"
              size="xs"
              shape="custom"
              class="w-max"
              :disabled="isSaving"
              @click="saveService"
            >
              {{ isSaving ? 'Saving...' : 'Save' }}
            </GroBasicButton>
          </div>
        </div>
      </Transition>
    </template>
  </MainDashboard>

  <!-- ── Manage Variations modal ──────────────────────────────── -->
  <Teleport to="body">
    <Transition name="fade">
      <div
        v-if="showVariationsModal"
        class="fixed inset-0 bg-black/30 backdrop-blur-sm flex items-center justify-center z-50"
        @click.self="showVariationsModal = false"
      >
        <div class="bg-white rounded-2xl p-6 w-[480px] shadow-2xl mx-4">
          <div class="flex justify-between items-center mb-5">
            <h3 class="text-lg font-semibold text-[#1E212B]">
              Manage Variations
            </h3>
            <button
              class="w-8 h-8 flex items-center justify-center rounded-full hover:bg-[#EDEDEE] text-[#4B4D55] transition-colors"
              @click="showVariationsModal = false"
            >
              <svg
                class="w-4 h-4"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <path d="M18 6L6 18M6 6l12 12" />
              </svg>
            </button>
          </div>

          <div class="mb-4">
            <GroBasicInput
              v-model="variationsTitle"
              placeholder="e.g., Age group"
            >
              What is the title of your variations entry?
            </GroBasicInput>
          </div>

          <div class="border border-[#EDEDEE] rounded-xl overflow-hidden mb-3">
            <div class="grid grid-cols-[1fr_140px_32px] bg-[#F6F6F7] border-b border-[#EDEDEE] px-4 py-2">
              <span class="text-xs font-semibold text-[#1E212B]">Variation Name</span>
              <span class="text-xs font-semibold text-[#1E212B]">Price</span>
              <span />
            </div>
            <div
              v-for="v in editVariations"
              :key="v.id"
              class="grid grid-cols-[1fr_140px_32px] items-center border-b border-[#EDEDEE] last:border-0 px-4 py-2 gap-2"
            >
              <GroBasicInput
                v-model="v.name"
                placeholder="e.g., Adult"
              />
              <div class="flex items-center bg-[#F6F6F7] border border-[#EDEDEE] rounded-lg focus-within:border-[#1E212B] transition-colors">
                <span class="pl-2 text-sm text-[#4B4D55]">₦</span>
                <input
                  v-model="v.price"
                  type="number"
                  placeholder="0"
                  class="flex-1 px-1.5 py-1.5 bg-transparent text-sm text-[#1E212B] outline-none placeholder-[#939499]"
                >
              </div>
              <button
                class="w-7 h-7 flex items-center justify-center text-[#939499] hover:text-[#AF513A] transition-colors"
                @click="removeEditVariation(v.id)"
              >
                <svg
                  class="w-4 h-4"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="1.5"
                >
                  <path d="M3 6h18M8 6V4h8v2M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6M10 11v6M14 11v6" />
                </svg>
              </button>
            </div>
          </div>

          <button
            class="flex items-center gap-1.5 text-sm font-medium text-[#2176AE] hover:underline mb-6"
            @click="addEditVariation"
          >
            <svg
              class="w-4 h-4"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2.5"
            >
              <path d="M12 5v14M5 12h14" />
            </svg>
            Add variation
          </button>

          <div class="flex justify-end gap-3">
            <GroBasicButton
              color="secondary"
              size="sm"
              shape="custom"
              class="w-max"
              @click="showVariationsModal = false"
            >
              Cancel
            </GroBasicButton>
            <GroBasicButton
              color="primary"
              size="sm"
              shape="custom"
              class="w-max"
              @click="saveVariations"
            >
              Save
            </GroBasicButton>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>

  <!-- ── Plan modal ─────────────────────────────────────────────── -->
  <Teleport to="body">
    <Transition name="fade">
      <div
        v-if="showPlanModal"
        class="fixed inset-0 bg-black/30 backdrop-blur-sm flex items-center justify-center z-50"
        @click.self="showPlanModal = false"
      >
        <div class="bg-white rounded-2xl shadow-2xl mx-4 w-[520px]">
          <!-- Header -->
          <div class="flex items-center justify-between px-6 pt-6 pb-4 border-b border-[#F6F6F7]">
            <div class="flex items-center gap-2">
              <button
                v-if="planModalStep === 2"
                class="w-7 h-7 flex items-center justify-center rounded-full hover:bg-[#EDEDEE] text-[#4B4D55] transition-colors"
                @click="planModalStep = 1"
              >
                <svg
                  class="w-4 h-4"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                >
                  <path d="M19 12H5M12 19l-7-7 7-7" />
                </svg>
              </button>
              <h3 class="text-base font-semibold text-[#1E212B]">
                {{ planModalStep === 1 ? 'What kind of plan are you selling?' : (editingPlanId ? 'Edit plan' : 'New plan') }}
              </h3>
            </div>
            <button
              class="w-8 h-8 flex items-center justify-center rounded-full hover:bg-[#EDEDEE] text-[#4B4D55] transition-colors"
              @click="showPlanModal = false"
            >
              <svg
                class="w-4 h-4"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <path d="M18 6L6 18M6 6l12 12" />
              </svg>
            </button>
          </div>

          <!-- Step 1: Type picker -->
          <div
            v-if="planModalStep === 1"
            class="p-6"
          >
            <div class="grid grid-cols-2 gap-4">
              <button
                class="relative p-5 rounded-xl border-2 text-left transition-all hover:border-[#2176AE] hover:shadow-sm"
                :class="newPlan.planType === 'subscription' ? 'border-[#2176AE] bg-[#F0F7FF]' : 'border-[#EDEDEE]'"
                @click="selectPlanType('subscription')"
              >
                <div
                  v-if="newPlan.planType === 'subscription'"
                  class="absolute top-3 right-3 w-5 h-5 rounded-full bg-[#2176AE] flex items-center justify-center"
                >
                  <svg
                    class="w-3 h-3 text-white"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="3"
                  >
                    <path d="M5 13l4 4L19 7" />
                  </svg>
                </div>
                <div class="w-10 h-10 rounded-full bg-[#2176AE] flex items-center justify-center mb-4">
                  <svg
                    class="w-5 h-5 text-white"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="1.5"
                  >
                    <path d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                  </svg>
                </div>
                <p class="text-sm font-semibold text-[#1E212B] mb-1">
                  Subscription
                </p>
                <p class="text-xs text-[#6F7177] leading-relaxed">
                  Clients get recurring charges and redeem sessions until the plan expires.
                </p>
                <p class="text-[10px] text-[#939499] mt-2">
                  e.g., 10 sessions a month, Yearly unlimited
                </p>
              </button>

              <button
                class="relative p-5 rounded-xl border-2 text-left transition-all hover:border-[#2176AE] hover:shadow-sm"
                :class="newPlan.planType === 'package' ? 'border-[#2176AE] bg-[#F0F7FF]' : 'border-[#EDEDEE]'"
                @click="selectPlanType('package')"
              >
                <div
                  v-if="newPlan.planType === 'package'"
                  class="absolute top-3 right-3 w-5 h-5 rounded-full bg-[#2176AE] flex items-center justify-center"
                >
                  <svg
                    class="w-3 h-3 text-white"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="3"
                  >
                    <path d="M5 13l4 4L19 7" />
                  </svg>
                </div>
                <div class="w-10 h-10 rounded-full bg-[#2176AE] flex items-center justify-center mb-4">
                  <svg
                    class="w-5 h-5 text-white"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="1.5"
                  >
                    <path d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                  </svg>
                </div>
                <p class="text-sm font-semibold text-[#1E212B] mb-1">
                  Package
                </p>
                <p class="text-xs text-[#6F7177] leading-relaxed">
                  Clients get charged once and redeem a limited number of sessions.
                </p>
                <p class="text-[10px] text-[#939499] mt-2">
                  e.g., 10 sessions total, valid for 90 days
                </p>
              </button>
            </div>

            <div class="flex justify-end mt-5">
              <button
                class="text-sm text-[#6F7177] hover:text-[#1E212B] transition-colors"
                @click="showPlanModal = false"
              >
                Cancel
              </button>
            </div>
          </div>

          <!-- Step 2: Plan form -->
          <div
            v-else-if="planModalStep === 2"
            class="p-6 flex flex-col gap-4"
          >
            <!-- Plan name -->
            <GroBasicInput
              v-model="newPlan.name"
              placeholder="e.g., Monthly membership"
            >
              Plan name
            </GroBasicInput>

            <!-- Price -->
            <GroCurrencyInput
              v-model:amount="newPlan.price"
              v-model:currency="priceCurrency"
              label="Price"
              placeholder="0"
            />

            <!-- Subscription: billing cycle -->
            <div v-if="newPlan.planType === 'subscription'">
              <label class="block text-xs font-medium text-[#1E212B] mb-2">Billing cycle</label>
              <div class="grid grid-cols-2 gap-3">
                <button
                  class="py-2.5 rounded-lg border text-sm font-medium transition-colors"
                  :class="newPlan.billingCycle === 'monthly' ? 'border-[#2176AE] bg-[#F0F7FF] text-[#2176AE]' : 'border-[#EDEDEE] text-[#6F7177] hover:border-[#94BDD8]'"
                  @click="newPlan.billingCycle = 'monthly'"
                >
                  Monthly
                </button>
                <button
                  class="py-2.5 rounded-lg border text-sm font-medium transition-colors"
                  :class="newPlan.billingCycle === 'yearly' ? 'border-[#2176AE] bg-[#F0F7FF] text-[#2176AE]' : 'border-[#EDEDEE] text-[#6F7177] hover:border-[#94BDD8]'"
                  @click="newPlan.billingCycle = 'yearly'"
                >
                  Yearly
                </button>
              </div>
            </div>

            <!-- Sessions -->
            <div>
              <label class="block text-xs font-medium text-[#1E212B] mb-1">
                {{ newPlan.planType === 'subscription' ? 'Sessions per cycle' : 'Total sessions' }}
              </label>
              <div class="flex items-center gap-3">
                <div class="flex-1 flex items-center bg-[#F6F6F7] border border-[#EDEDEE] rounded-lg focus-within:border-[#1E212B] transition-colors">
                  <input
                    v-model="newPlan.sessionCount"
                    type="number"
                    min="1"
                    :placeholder="newPlan.sessionCount === '' ? 'Unlimited' : ''"
                    :disabled="newPlan.sessionCount === ''"
                    class="flex-1 px-3 py-2.5 bg-transparent text-sm text-[#1E212B] outline-none placeholder-[#939499] disabled:opacity-50"
                  >
                </div>
                <label class="flex items-center gap-2 text-sm text-[#4B4D55] whitespace-nowrap cursor-pointer">
                  <input
                    type="checkbox"
                    class="w-4 h-4 rounded accent-[#2176AE] cursor-pointer"
                    :checked="newPlan.sessionCount === ''"
                    @change="newPlan.sessionCount = (newPlan.sessionCount === '') ? '1' : ''"
                  >
                  Unlimited
                </label>
              </div>
            </div>

            <!-- Package: validity -->
            <div v-if="newPlan.planType === 'package'">
              <label class="block text-xs font-medium text-[#1E212B] mb-1">Valid for (days after purchase)</label>
              <div class="flex items-center bg-[#F6F6F7] border border-[#EDEDEE] rounded-lg focus-within:border-[#1E212B] transition-colors">
                <input
                  v-model="newPlan.validityDays"
                  type="number"
                  min="1"
                  placeholder="e.g., 90"
                  class="flex-1 px-3 py-2.5 bg-transparent text-sm text-[#1E212B] outline-none placeholder-[#939499]"
                >
                <span class="pr-3 text-sm text-[#939499]">days</span>
              </div>
            </div>

            <div class="flex justify-end gap-3 pt-1">
              <GroBasicButton
                color="secondary"
                size="sm"
                shape="custom"
                class="w-max"
                @click="showPlanModal = false"
              >
                Cancel
              </GroBasicButton>
              <GroBasicButton
                color="primary"
                size="sm"
                shape="custom"
                class="w-max"
                :disabled="!newPlan.name.trim()"
                @click="savePlan"
              >
                {{ editingPlanId ? 'Save changes' : 'Add plan' }}
              </GroBasicButton>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style>
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s; }
.fade-enter-from, .fade-leave-to { opacity: 0; }

button {
  cursor: pointer;
}
</style>

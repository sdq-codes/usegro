<script setup lang="ts">
import { ref, computed, watch, onMounted, onBeforeUnmount, nextTick } from 'vue'
import MainDashboard from "@/components/dashboard/main-dashboard.vue"
import GroBasicButton from "@/components/buttons/GroBasicButton.vue"
import GroBasicInput from "@/components/forms/input/GroBasicInput.vue"
import GroCurrencyInput from "@/components/forms/input/GroCurrencyInput.vue"
import GroToggle from "@/components/forms/select/GroToggle.vue"
import GroBasicSelect from "@/components/forms/select/GroBasicSelect.vue"
import GroStandardCategorySelect, { type StandardCategorySelection } from "@/components/forms/select/GroStandardCategorySelect.vue"
import GroColorPicker from "@/components/forms/select/GroColorPicker.vue"
import GroOptionPicker from "@/components/forms/select/GroOptionPicker.vue"
import CatalogTagSelect from "@/components/forms/select/CatalogTagSelect.vue"
import CommentEditor from "@/components/editor/CommentEditor.vue"
import CustomFormField from "@/components/forms/CustomForm/CustomFormField.vue"
import { HugeiconsIcon } from '@hugeicons/vue'
import { DragDropHorizontalIcon } from '@hugeicons/core-free-icons'
import { useRoute, useRouter } from 'nuxt/app'
import BackLink from "@/components/navigation/BackLink.vue"
import GroFormSection from "@/components/forms/GroFormSection.vue"
import { useFormFieldManager } from "@/composables/helpers/managers/forms/useFormFieldManager"
import type { FormField } from "@/composables/helpers/types/form"
import { useProductAPI } from "@/composables/api/catalog/product"
import { useCatalogMediaAPI } from "@/composables/api/catalog/media"
import { notify } from "@/composables/helpers/notification/notification"

const RIBBON_OPTIONS = [
  { value: '', label: 'None' },
  { value: 'new', label: 'New' },
  { value: 'sale', label: 'Sale' },
  { value: 'hot', label: 'Hot' },
]

const route = useRoute()
const router = useRouter()
const productId = route.params.id as string
const isNew = productId === 'new'

const isLoading = ref(!isNew)
const isSaving = ref(false)
const MAX_MEDIA = 10

// ── Media state ──────────────────────────────────────────────
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
const activeExistingCount = computed(() => orderedMedia.value.length)
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

// ── Form ────────────────────────────────────────────────────
const form = ref({
  name: '',
  brand: '',
  ribbon: '',
  description: '',
  status: 'Active',
  showInStore: false,
  price: '',
  priceCurrency: (process.client && localStorage.getItem('currency')) || 'USD',
  costPerItem: '',
  sku: '',
  barcode: '',
  chargeTax: true,
  trackInventory: true,
  continueSellingWhenOutOfStock: false,
  stockStatus: 'in_stock',
  quantity: '',
})

const loadedSku = ref('')
const selectedStandardCategory = ref<StandardCategorySelection | null>(null)
const tags = ref<string[]>([])

interface ProductOption {
  id: string
  name: string
  values: string[]
  newValue: string
  isEditing: boolean
}

interface Variant {
  id: string
  name: string
  price: string
  quantity: string
  imageKey: string
  imageUrl: string
  imagePreview: string  // local object URL or loaded URL
}

const productOptions = ref<ProductOption[]>([])
const variants = ref<Variant[]>([])
const editingVariants = ref(false) // kept for compat
const groupByOptionIndex = ref(0)
const expandedGroups = ref<Set<string>>(new Set())
const showGroupByDropdown = ref(false)
const groupByDropdownRef = ref<HTMLElement | null>(null)
const selectedVariants = ref<string[]>([])

const cartesian = (...arrays: string[][]): string[][] => {
  if (arrays.length === 0) return []
  return arrays.reduce<string[][]>(
    (acc, arr) => acc.flatMap(combo => arr.map(v => [...combo, v])),
    [[]]
  )
}

const MAX_VARIANTS = 2048

const generateVariants = (): Variant[] => {
  const active = productOptions.value.filter(o => o.values.filter(v => v.trim()).length > 0)
  if (active.length === 0) return []
  const combos = cartesian(...active.map(o => o.values.filter(v => v.trim())))
  return combos.slice(0, MAX_VARIANTS).map((combo) => ({
    id: combo.join(' / '),
    name: combo.join(' / '),
    price: form.value.price || '',
    quantity: '',
    imageKey: '',
    imageUrl: '',
    imagePreview: '',
  }))
}

watch(
  () => form.value.status,
  (status) => {
    if (status.toLowerCase() !== 'active') form.value.showInStore = false
  }
)

watch(
  productOptions,
  () => {
    const prev = variants.value
    variants.value = generateVariants().map(v => {
      const existing = prev.find(e => e.name === v.name)
      return existing ? { ...existing } : v
    })
  },
  { deep: true, immediate: true }
)

watch(
  [variants, selectedVariants],
  () => {
    if (selectedVariants.value.length === 0) {
      form.value.trackInventory = false
      return
    }
    const selected = variants.value.filter(v => selectedVariants.value.includes(v.id))
    const allInStock = selected.length > 0 && selected.every(v => !v.quantity || v.quantity.trim() === '')
    if (allInStock) form.value.trackInventory = false
  },
  { deep: true }
)

const hasVariants = computed(() => variants.value.length > 0)

const isAllVariantsSelected = computed(
  () => selectedVariants.value.length === variants.value.length && variants.value.length > 0
)

const toggleAllVariants = () => {
  if (isAllVariantsSelected.value) selectedVariants.value = []
  else selectedVariants.value = variants.value.map(v => v.id)
}

const toggleVariant = (id: string) => {
  const idx = selectedVariants.value.indexOf(id)
  if (idx > -1) selectedVariants.value.splice(idx, 1)
  else selectedVariants.value.push(id)
}

const CURRENCY_SYMBOLS: Record<string, string> = {
  NGN: '₦', USD: '$', EUR: '€', GBP: '£', GHS: 'GH₵', KES: 'KSh', ZAR: 'R',
}

const currencySymbol = computed(() => CURRENCY_SYMBOLS[form.value.priceCurrency] ?? form.value.priceCurrency)

const effectivePrice = computed(() => parseFloat(form.value.price) || 0)

const profit = computed(() => effectivePrice.value - (parseFloat(form.value.costPerItem) || 0))

const margin = computed(() => {
  if (effectivePrice.value === 0) return '0'
  return ((profit.value / effectivePrice.value) * 100).toFixed(1)
})

const selectedVariantsTotalQty = computed(() =>
  variants.value
    .filter(v => selectedVariants.value.includes(v.id))
    .reduce((sum, v) => sum + (parseInt(v.quantity) || 0), 0)
)

const groupByOption = computed(() => productOptions.value[groupByOptionIndex.value] ?? productOptions.value[0])

const groupedVariants = computed(() => {
  if (!hasVariants.value || productOptions.value.length === 0) return []
  const idx = groupByOptionIndex.value
  const groups = new Map<string, Variant[]>()
  for (const v of variants.value) {
    const parts = v.name.split(' / ')
    const key = parts[idx] ?? v.name
    if (!groups.has(key)) groups.set(key, [])
    groups.get(key)!.push(v)
  }
  return Array.from(groups.entries()).map(([value, items]) => ({ value, items }))
})

const allGroupsCollapsed = computed(() => expandedGroups.value.size === 0)

const toggleGroup = (value: string) => {
  const next = new Set(expandedGroups.value)
  if (next.has(value)) next.delete(value)
  else next.add(value)
  expandedGroups.value = next
}

const toggleAllGroups = () => {
  if (allGroupsCollapsed.value) {
    expandedGroups.value = new Set(groupedVariants.value.map(g => g.value))
  } else {
    expandedGroups.value = new Set()
  }
}

const isGroupSelected = (items: Variant[]) =>
  items.length > 0 && items.every(v => selectedVariants.value.includes(v.id))

const isGroupPartiallySelected = (items: Variant[]) =>
  items.some(v => selectedVariants.value.includes(v.id)) && !isGroupSelected(items)

const toggleGroupSelection = (items: Variant[]) => {
  const ids = items.map(v => v.id)
  if (isGroupSelected(items)) {
    selectedVariants.value = selectedVariants.value.filter(id => !ids.includes(id))
  } else {
    selectedVariants.value = [...new Set([...selectedVariants.value, ...ids])]
  }
}

const getGroupPrice = (items: Variant[]) => {
  const prices = new Set(items.map(v => v.price))
  return prices.size === 1 ? items[0].price : ''
}

const setGroupPrice = (items: Variant[], price: string) => {
  items.forEach(v => { v.price = price })
}

const getGroupQty = (items: Variant[]) =>
  items.reduce((sum, v) => sum + (parseInt(v.quantity) || 0), 0)

const onGroupByClickOutside = (e: MouseEvent) => {
  if (groupByDropdownRef.value && !groupByDropdownRef.value.contains(e.target as Node)) {
    showGroupByDropdown.value = false
  }
}

onMounted(() => document.addEventListener('mousedown', onGroupByClickOutside))
onBeforeUnmount(() => document.removeEventListener('mousedown', onGroupByClickOutside))

const colorHexMap: Record<string, string> = {
  beige: '#F5F5DC', black: '#1C1C1C', blue: '#2563EB', bronze: '#CD7F32',
  brown: '#92400E', burgundy: '#800020', coral: '#FF6B6B', cream: '#FFFDD0',
  cyan: '#06B6D4', gold: '#F59E0B', gray: '#6B7280', green: '#16A34A',
  indigo: '#4338CA', ivory: '#FFFFF0', khaki: '#C3B091', lavender: '#E6E6FA',
  maroon: '#800000', mint: '#98D8C8', navy: '#1E3A8A', olive: '#808000',
  orange: '#F97316', pink: '#EC4899', purple: '#9333EA', red: '#DC2626',
  rose: '#F43F5E', silver: '#9CA3AF', tan: '#D2B48C', teal: '#0D9488',
  turquoise: '#40E0D0', violet: '#7C3AED', white: '#F9FAFB', yellow: '#FACC15',
}

const categoryAttributes = computed(() => selectedStandardCategory.value?.attributes ?? [])
const getAttributeValues = (optionName: string) =>
  categoryAttributes.value.find(a => a.name.toLowerCase() === optionName.toLowerCase())?.values ?? []

const getAttributeValueStrings = (optionName: string) =>
  getAttributeValues(optionName).map(v => v.name)

const hasAttributeValues = (optionName: string) =>
  getAttributeValues(optionName).length > 0
const availableCategoryAttributes = computed(() => {
  const usedNames = new Set(productOptions.value.map(o => o.name.toLowerCase()))
  return categoryAttributes.value.filter(a => !usedNames.has(a.name.toLowerCase()))
})
const hasCategoryAttributes = computed(() => availableCategoryAttributes.value.length > 0)

const showAddOptionDropdown = ref(false)
const addOptionDropdownRef = ref<HTMLElement | null>(null)

const addOption = (name = '') => {
  if (productOptions.value.length >= 3) return
  productOptions.value.push({ id: String(Date.now()), name, values: [], newValue: '', isEditing: true })
  showAddOptionDropdown.value = false
}

const onAddOptionClickOutside = (e: MouseEvent) => {
  if (addOptionDropdownRef.value && !addOptionDropdownRef.value.contains(e.target as Node)) {
    showAddOptionDropdown.value = false
  }
}

onMounted(() => document.addEventListener('mousedown', onAddOptionClickOutside))
onBeforeUnmount(() => document.removeEventListener('mousedown', onAddOptionClickOutside))

const removeOption = (id: string) => {
  productOptions.value = productOptions.value.filter(o => o.id !== id)
}

const doneEditingOption = (option: ProductOption) => {
  if (option.newValue.trim()) { option.values.push(option.newValue.trim()); option.newValue = '' }
  option.values = option.values.filter(v => v.trim())
  option.isEditing = false
}

const addValueToOption = (option: ProductOption) => {
  if (option.newValue.trim()) { option.values.push(option.newValue.trim()); option.newValue = '' }
}

const removeValueFromOption = (option: ProductOption, idx: number) => {
  option.values.splice(idx, 1)
}

const draggingOptionId = ref<string | null>(null)
const onOptionDragStart = (id: string) => { draggingOptionId.value = id }
const onOptionDrop = (targetId: string) => {
  if (!draggingOptionId.value || draggingOptionId.value === targetId) return
  const from = productOptions.value.findIndex(o => o.id === draggingOptionId.value)
  const to = productOptions.value.findIndex(o => o.id === targetId)
  const items = [...productOptions.value]
  const [moved] = items.splice(from, 1)
  items.splice(to, 0, moved)
  productOptions.value = items
  draggingOptionId.value = null
}

const draggingValue = ref<{ optionId: string; idx: number } | null>(null)
const onValueDragStart = (optionId: string, idx: number) => { draggingValue.value = { optionId, idx } }
const onValueDrop = (optionId: string, toIdx: number) => {
  if (!draggingValue.value || draggingValue.value.optionId !== optionId) return
  const fromIdx = draggingValue.value.idx
  if (fromIdx === toIdx) return
  const option = productOptions.value.find(o => o.id === optionId)
  if (!option) return
  const vals = [...option.values]
  const [moved] = vals.splice(fromIdx, 1)
  vals.splice(toIdx, 0, moved)
  option.values = vals
  draggingValue.value = null
}

const saveVariants = () => { editingVariants.value = false }

const onVariantImagePick = async (variant: Variant, file: File) => {
  variant.imagePreview = URL.createObjectURL(file)
  const { keys, error } = await useCatalogMediaAPI().presignAndUpload([file])
  if (error || keys.length === 0) {
    variant.imagePreview = variant.imageUrl
    return
  }
  variant.imageKey = keys[0]
  variant.imageUrl = variant.imagePreview
}

interface VariantGroupImage {
  imageKey: string
  imageUrl: string
  imagePreview: string
}

// Keyed by "optionName:value"
const variantGroupImages = ref<Record<string, VariantGroupImage>>({})

const getGroupImageKey = (groupValue: string) => {
  const optName = groupByOption.value?.name ?? ''
  return `${optName}:${groupValue}`
}

const getGroupImage = (groupValue: string): VariantGroupImage => {
  return variantGroupImages.value[getGroupImageKey(groupValue)] ?? { imageKey: '', imageUrl: '', imagePreview: '' }
}

const onGroupImagePick = async (groupValue: string, file: File) => {
  const key = getGroupImageKey(groupValue)
  const preview = URL.createObjectURL(file)
  variantGroupImages.value[key] = { imageKey: '', imageUrl: preview, imagePreview: preview }
  const { keys, error } = await useCatalogMediaAPI().presignAndUpload([file])
  if (error || keys.length === 0) {
    variantGroupImages.value[key].imagePreview = ''
    variantGroupImages.value[key].imageUrl = ''
    return
  }
  variantGroupImages.value[key].imageKey = keys[0]
}

const formatVariantPrice = (price: string) => {
  const num = parseFloat(price)
  const sym = currencySymbol.value
  if (isNaN(num) || num === 0) return `${sym}0.00`
  return `${sym}${num.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`
}

const variantQtyDisplay = (v: Variant) => v.quantity ? `${v.quantity} Units` : 'In stock'

// ── Additional Info ──────────────────────────────────────────
const { addField } = useFormFieldManager()
const additionalFields = ref<FormField[]>([])
const additionalAnswers = ref<Record<string, string | string[]>>({})
const showFieldPicker = ref(false)

const ADDITIONAL_FIELD_TYPES = [
  { type: 'custom', label: 'Text' },
  { type: 'notes', label: 'Notes' },
  { type: 'birthdate', label: 'Date' },
]

const handleAddAdditionalField = (type: string) => {
  const newField = addField(type, additionalFields.value, additionalAnswers.value)
  if (!newField) return
  if (Array.isArray(newField)) additionalFields.value.push(...newField)
  else additionalFields.value.push(newField)
  showFieldPicker.value = false
}

const handleAdditionalLabelChanged = (payload: { slug: string; label: string }) => {
  const field = additionalFields.value.find(f => f.slug === payload.slug)
  if (field) field.label = payload.label
}

// ── Load product (edit mode) ─────────────────────────────────
const footerSentinel = ref<HTMLElement | null>(null)
const showFixedBar = ref(false)

onMounted(async () => {
  const setupObserver = async () => {
    await nextTick()
    const observer = new IntersectionObserver(
      (entries) => { showFixedBar.value = !entries[0]?.isIntersecting },
      { threshold: 0 }
    )
    if (footerSentinel.value) observer.observe(footerSentinel.value)
    onBeforeUnmount(() => {
      observer.disconnect()
      allMedia.value.filter(m => m.kind === 'new').forEach(m => URL.revokeObjectURL(m.url))
    })
  }

  if (isNew) {
    await setupObserver()
    return
  }

  const result = await useProductAPI().GetProduct(productId)
  if (!result.success || !result.data?.data) {
    notify('Product not found', 'error')
    router.push('/catalog/products')
    return
  }

  const p = result.data.data
  const detail = p.product_detail

  form.value = {
    name: p.name,
    brand: detail?.brand ?? '',
    ribbon: detail?.ribbon ?? '',
    description: p.description,
    status: p.status ? p.status.charAt(0).toUpperCase() + p.status.slice(1) : 'Active',
    showInStore: p.show_in_store,
    price: p.price ? String(p.price) : '',
    priceCurrency: p.price_currency || 'USD',
    costPerItem: p.cost_per_item ? String(p.cost_per_item) : '',
    sku: detail?.sku ?? '',
    barcode: p.barcode ?? '',
    chargeTax: p.charge_tax ?? true,
    trackInventory: detail?.track_inventory ?? false,
    continueSellingWhenOutOfStock: detail?.continue_selling_when_out_of_stock ?? false,
    stockStatus: detail?.stock_status ?? 'in_stock',
    quantity: detail?.quantity ? String(detail.quantity) : '',
  }

  loadedSku.value = detail?.sku ?? ''
  selectedStandardCategory.value = p.standard_category
    ? { id: p.standard_category.id, name: p.standard_category.name, attributes: [] }
    : null
  tags.value = (p.tags ?? []).map((t: any) => typeof t === 'string' ? t : t.id)

  // Load media — featured image (display_image=true) goes first, rest sorted by position
  const rawMedia = p.media ?? []
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

  if (p.additional_fields?.length) {
    p.additional_fields.forEach((f, i) => {
      const slug = `additional_${i}`
      additionalFields.value.push({
        SK: `FIELD#additional-${i}`,
        slug,
        section: 'Extra fields',
        order: i,
        fieldTypeName: f.field_type,
        label: f.label,
        required: false,
        options: [],
        configs: [],
        validations: [],
      } as FormField)
      additionalAnswers.value[slug] = f.value
    })
  }

  if (detail?.variants?.length) {
    variants.value = detail.variants.map(v => ({
      id: v.name,
      name: v.name,
      price: v.price ? String(v.price) : '',
      quantity: v.quantity ? String(v.quantity) : '',
      imageKey: v.image_key ?? '',
      imageUrl: v.image_url ?? '',
      imagePreview: v.image_url ?? '',
    }))
  }

  if (detail?.options?.length) {
    productOptions.value = detail.options
      .sort((a, b) => a.position - b.position)
      .map(o => ({
        id: o.id,
        name: o.name,
        values: o.values.sort((a, b) => a.position - b.position).map(v => v.value),
        newValue: '',
        isEditing: false,
      }))
  }

  if (detail?.variant_groups?.length) {
    detail.variant_groups.forEach((g: any) => {
      if (g.image_key) {
        const key = `${g.option_name}:${g.value}`
        variantGroupImages.value[key] = { imageKey: g.image_key, imageUrl: g.image_url, imagePreview: g.image_url }
      }
    })
  }

  isLoading.value = false
  await setupObserver()
})

// ── Build payload ────────────────────────────────────────────
const buildPayload = (status: string) => ({
  name: form.value.name,
  brand: form.value.brand,
  ribbon: (form.value.ribbon as any)?.value ?? form.value.ribbon ?? '',
  description: form.value.description,
  price: parseFloat(form.value.price) || 0,
  price_currency: form.value.priceCurrency,
  cost_per_item: parseFloat(form.value.costPerItem) || 0,
  status,
  barcode: form.value.barcode,
  charge_tax: form.value.chargeTax,
  continue_selling_when_out_of_stock: form.value.continueSellingWhenOutOfStock,
  show_in_store: form.value.showInStore,
  sku: hasVariants.value ? '' : form.value.sku,
  track_inventory: form.value.trackInventory,
  stock_status: (form.value.stockStatus as any)?.value ?? form.value.stockStatus ?? 'in_stock',
  quantity: parseInt(form.value.quantity) || 0,
  item_type: 'physical',
  tag_ids: tags.value,
  standard_category_id: selectedStandardCategory.value?.id ?? null,
  options: productOptions.value.map((o, i) => ({
    name: o.name,
    values: o.values.filter(v => v.trim()),
    position: i,
  })),
  variants: variants.value.map(v => ({
    name: v.name,
    price: parseFloat(v.price) || 0,
    quantity: parseInt(v.quantity) || 0,
    image_key: v.imageKey,
    image_url: v.imageUrl,
  })),
  additional_fields: additionalFields.value.map((f, i) => ({
    label: f.label,
    field_type: f.fieldTypeName,
    value: String(additionalAnswers.value[f.slug] ?? ''),
    position: i,
  })),
  variant_groups: Object.entries(variantGroupImages.value)
    .filter(([_, g]) => g.imageKey)
    .map(([key, g]) => {
      const colonIdx = key.indexOf(':')
      const optionName = key.slice(0, colonIdx)
      const value = key.slice(colonIdx + 1)
      return { option_name: optionName, value, image_key: g.imageKey, image_url: g.imageUrl }
    }),
})

// ── Save ─────────────────────────────────────────────────────
const saveProduct = async (status: string) => {
  if (!form.value.name.trim()) {
    notify('Product name is required', 'error')
    return
  }
  isSaving.value = true
  try {
    // 1. Delete pending removals
    const toDelete = allMedia.value.filter(m => m.kind === 'existing' && m.pendingDelete)
    if (toDelete.length > 0) {
      await Promise.all(toDelete.map(m => useCatalogMediaAPI().DeleteMedia(m.id)))
    }

    // 2. Upload new files in their current order
    const newItems = orderedMedia.value.filter(m => m.kind === 'new')
    let mediaKeys: string[] = []
    if (newItems.length > 0) {
      const { keys, error } = await useCatalogMediaAPI().presignAndUpload(newItems.map(m => m.file!))
      if (error) {
        notify(error, 'error')
        return
      }
      mediaKeys = keys
    }

    // 3. First item in orderedMedia is the featured/display image
    const firstItem = orderedMedia.value[0]
    const displayImageId = firstItem?.kind === 'existing' ? firstItem.id : undefined
    const newFeaturedIdx = newItems.findIndex(m => m.id === firstItem?.id)
    const displayImageKey = newFeaturedIdx >= 0 ? mediaKeys[newFeaturedIdx] : undefined

    // 5. Build and send payload
    const payload = {
      ...buildPayload(status),
      media_keys: mediaKeys,
      display_image_id: displayImageId,
      display_image_key: displayImageKey,
    }
    if (status === 'draft') payload.show_in_store = false

    const result = isNew
      ? await useProductAPI().CreateProduct(payload)
      : await useProductAPI().UpdateProduct(productId, payload)

    if (result.success) {
      notify(isNew ? 'Product created' : 'Product updated', 'success')
      router.push('/catalog/products')
    } else {
      notify((result as any).error || 'Failed to save product', 'error')
    }
  } finally {
    isSaving.value = false
  }
}

// ── Accordion ────────────────────────────────────────────────
const openSections = ref<Record<string, boolean>>({
  productInfo: true,
  pricing: false,
  options: false,
  inventory: false,
  additional: false,
})
const toggleSection = (key: string) => {
  const isOpen = openSections.value[key]
  for (const k in openSections.value) openSections.value[k] = false
  openSections.value[key] = !isOpen
}
</script>

<template>
  <MainDashboard current="Products">
    <template #title>
      <BackLink to="/catalog/products" label="Products" />
    </template>

    <template #body>
      <div
        v-if="isLoading"
        class="bg-white rounded-3xl mt-4 flex items-center justify-center py-32 text-sm text-[#6F7177]"
      >
        Loading product...
      </div>

      <div v-else class="bg-white rounded-3xl gap-6 mt-4">
        <!-- Header -->
        <div class="flex justify-between px-8 pt-8 items-center border-b border-[#F6F6F7] pb-4">
          <h2 class="text-xl font-semibold text-[#1E212B]">
            {{ isNew ? (form.name.trim() || 'New Product') : (form.name.trim() || 'Edit Product') }}
          </h2>
          <button
            class="w-8 h-8 flex items-center justify-center rounded-full hover:bg-[#EDEDEE] text-[#4B4D55] hover:text-[#1E212B] transition-colors"
            @click="router.push('/catalog/products')"
          >
            <svg class="w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M18 6L6 18M6 6l12 12" />
            </svg>
          </button>
        </div>

        <div class="grid grid-cols-5 gap-y-6 md:border-b md:border-[#EDEDEE]">
          <!-- Left column -->
          <div class="col-span-5 md:col-span-3 pt-3 md:border-r md:border-[#EDEDEE] md:border-dashed gap-4">

            <!-- Product Info -->
            <GroFormSection title="Product Info" :open="openSections.productInfo" @toggle="toggleSection('productInfo')">
              <p class="text-[10px] font-semibold text-[#6F7177] uppercase tracking-widest mb-4">Basic Info</p>

              <div class="flex flex-col gap-4">
                <GroBasicInput v-model="form.name" placeholder="Product Name">Name *</GroBasicInput>

                <div class="grid grid-cols-2 gap-4">
                  <GroBasicInput v-model="form.brand" placeholder="Add text" info="This is the brand for the product.">Brand</GroBasicInput>
                  <GroBasicSelect v-model="form.ribbon" :options="RIBBON_OPTIONS" placeholder="None">Ribbon</GroBasicSelect>
                </div>

                <div>
                  <GroStandardCategorySelect v-model="selectedStandardCategory">
                    <label class="text-xs font-medium text-[#1E212B]">Category</label>
                  </GroStandardCategorySelect>
                </div>

                <!-- Description -->
                <div>
                  <div class="flex justify-between items-center mb-1">
                    <label class="text-xs font-medium text-[#1E212B]">Description</label>
                    <button class="text-xs text-[#2176AE] flex items-center gap-1 hover:underline">
                      <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                        <path d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
                      </svg>
                      Generate AI Text
                    </button>
                  </div>
                  <div class="border rounded-2xl border-gray-200">
                    <CommentEditor
                      v-model="form.description"
                      placeholder="For SEO improvement and help customers find this item by adding a unique & detailed product description."
                    />
                  </div>
                </div>

                <!-- Media -->
                <div class="border-t border-[#EDEDEE] mt-5 pt-6">
                  <!-- Unified grid: existing + new files, all draggable -->
                  <div v-if="totalMediaCount > 0" class="mb-5">
                    <div class="flex items-center gap-3 mb-4">
                      <p class="text-xs font-bold text-[#6F7177] uppercase tracking-widest">Images and Videos</p>
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
                        <img :src="orderedMedia[0].url" draggable="false" class="w-full h-full object-cover pointer-events-none" alt="" />
                        <div class="absolute top-2 left-2 bg-black/60 text-white text-[9px] font-semibold px-2 py-0.5 rounded-full flex items-center gap-1 pointer-events-none">
                          <svg class="w-2.5 h-2.5" viewBox="0 0 24 24" fill="currentColor">
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
                            <svg viewBox="0 0 24 24" fill="none" stroke="#AF513A" stroke-width="2" class="w-3.5 h-3.5">
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
                        <img :src="item.url" draggable="false" class="w-full h-full object-cover pointer-events-none" alt="" />
                        <div class="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex flex-col items-end justify-start p-1.5 gap-1.5">
                          <button
                            type="button"
                            class="w-7 h-7 rounded-full bg-white/90 flex items-center justify-center hover:bg-red-50 transition-colors"
                            @click.stop="markMediaForDelete(item)"
                          >
                            <svg viewBox="0 0 24 24" fill="none" stroke="#AF513A" stroke-width="2" class="w-3.5 h-3.5">
                              <path d="M6 18L18 6M6 6l12 12" />
                            </svg>
                          </button>
                          <button
                            type="button"
                            class="w-7 h-7 rounded-full bg-white/90 flex items-center justify-center hover:bg-yellow-50 transition-colors"
                            title="Set as featured"
                            @click.stop="setAsFeatured(item)"
                          >
                            <svg class="w-3.5 h-3.5 text-[#1E212B]" viewBox="0 0 24 24" fill="currentColor">
                              <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" />
                            </svg>
                          </button>
                        </div>
                      </div>
                    </div>
                    <p class="text-xs text-[#939499] mt-3">Drag to reorder · First image is the featured image</p>
                  </div>

                  <!-- Upload drop zone (same visual as GroFileUpload) -->
                  <div v-if="totalMediaCount < MAX_MEDIA">
                    <div v-if="totalMediaCount === 0" class="flex items-center gap-3 mb-2">
                      <p class="text-xs font-bold text-[#6F7177] uppercase tracking-widest">Images and Videos</p>
                    </div>
                    <label
                      class="relative border-2 border-dashed rounded-xl py-10 flex flex-col items-center justify-center gap-2 cursor-pointer transition-colors"
                      :class="isDraggingUpload ? 'border-[#2176AE] bg-[#EBF4FF]' : 'border-[#94BDD8] hover:bg-[#F0F7FF]'"
                      @dragenter.prevent="isDraggingUpload = true"
                      @dragover.prevent="isDraggingUpload = true"
                      @dragleave.prevent="isDraggingUpload = false"
                      @drop.prevent="onDropZoneDrop"
                    >
                      <svg class="w-9 h-9 text-[#2176AE]" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
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
              </div>
            </GroFormSection>

            <!-- Pricing -->
            <GroFormSection title="Pricing" :open="openSections.pricing" @toggle="toggleSection('pricing')">
              <div class="flex flex-col gap-4">
                <div class="grid grid-cols-2 gap-4">
                  <GroCurrencyInput
                    v-model:amount="form.price"
                    v-model:currency="form.priceCurrency"
                    label="Price"
                    placeholder="0"
                    class="md:col-span-1"
                  />
                </div>

                <div class="grid grid-cols-4 gap-4">
                  <div class="md:col-span-2">
                    <label class="text-xs font-medium text-[#1E212B] mb-1 block">Cost</label>
                    <div class="flex items-center bg-[#F6F6F7] border border-[#EDEDEE] rounded-lg hover:border-[#94BDD8] focus-within:border-[#1E212B] transition-colors">
                      <span class="pl-3 text-sm font-medium text-[#4B4D55]">{{ currencySymbol }}</span>
                      <input
                        v-model="form.costPerItem"
                        type="number"
                        placeholder="0"
                        class="flex-1 px-2 py-2 bg-transparent text-sm text-[#1E212B] outline-none placeholder-[#939499]"
                      >
                    </div>
                  </div>
                  <div>
                    <label class="text-xs font-medium text-[#1E212B] mb-1 block">Profit</label>
                    <div class="pr-3 py-2 text-sm text-[#939499]">
                      {{ profit > 0 ? `${currencySymbol}${profit.toLocaleString()}` : '--' }}
                    </div>
                  </div>
                  <div>
                    <label class="text-xs font-medium text-[#1E212B] mb-1 block">Margin</label>
                    <div class="pr-3 py-2 text-sm text-[#939499]">
                      {{ Number(margin) > 0 ? `${margin}%` : '--' }}
                    </div>
                  </div>
                </div>
                <GroToggle v-model="form.chargeTax">Charge tax on this product</GroToggle>
              </div>
            </GroFormSection>

            <!-- Inventory -->
            <GroFormSection title="Inventory" :open="openSections.inventory" @toggle="toggleSection('inventory')">
              <div class="flex flex-col gap-4">
                <GroToggle v-model="form.trackInventory">Track inventory</GroToggle>

                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <GroBasicInput v-model="form.sku" placeholder="e.g. PROD-001">SKU (Stock Keeping Unit)</GroBasicInput>
                  <GroBasicInput v-model="form.barcode" placeholder="e.g. 012345678905">Barcode (ISBN, UPC, GTIN, etc.)</GroBasicInput>
                </div>

                <template v-if="form.trackInventory">
                  <GroBasicInput
                    :model-value="hasVariants ? String(selectedVariantsTotalQty) : form.quantity"
                    placeholder="--"
                    type="number"
                    :disabled="hasVariants"
                    @update:model-value="(v: string | undefined) => { if (!hasVariants) form.quantity = v ?? '' }"
                  >Qty</GroBasicInput>
                  <GroToggle v-model="form.continueSellingWhenOutOfStock">Continue selling when out of stock</GroToggle>
                </template>
              </div>
            </GroFormSection>

            <!-- Product Options -->
            <GroFormSection title="Product Options" :open="openSections.options" @toggle="toggleSection('options')">
              <div class="flex items-center justify-between mb-4">
                <p class="text-sm text-[#4B4D55] leading-relaxed">
                  <template v-if="productOptions.length === 0">To create variants of this product, first add options, like sizes and colors.</template>
                  <template v-else>You can add up to 3 options per product.</template>
                </p>
                <span class="text-xs text-[#939499] flex-shrink-0 ml-4">{{ productOptions.length }}/3</span>
              </div>

              <div class="flex flex-col gap-3">
                <div
                  v-for="option in productOptions"
                  :key="option.id"
                  draggable="true"
                  @dragstart="onOptionDragStart(option.id)"
                  @dragover.prevent
                  @drop.prevent="onOptionDrop(option.id)"
                >
                  <!-- Collapsed -->
                  <div v-if="!option.isEditing" class="flex items-center gap-3 py-1">
                    <div class="text-[#C0C0C0] cursor-grab shrink-0">
                      <HugeiconsIcon :icon="DragDropHorizontalIcon" :size="20" color="currentColor" />
                    </div>
                    <div class="flex-1">
                      <p class="text-sm font-semibold text-[#1E212B]">{{ option.name }}</p>
                      <!-- Color swatches for Color option -->
                      <div v-if="option.name.toLowerCase() === 'color'" class="flex flex-wrap gap-1.5 mt-1">
                        <span
                          v-for="val in option.values"
                          :key="val"
                          class="flex items-center gap-1 px-2 py-0.5 bg-[#EDEDEE] text-xs text-[#1E212B] rounded-full"
                        >
                          <span
                            class="w-3 h-3 rounded-full flex-shrink-0 border border-black/10"
                            :style="{ background: colorHexMap[val.toLowerCase()] || 'repeating-conic-gradient(#e5e7eb 0% 25%, #fff 0% 50%) 0 0 / 6px 6px' }"
                          />
                          {{ val }}
                        </span>
                      </div>
                      <!-- Plain pills for other options -->
                      <div v-else class="flex flex-wrap gap-1.5 mt-1">
                        <span
                          v-for="val in option.values"
                          :key="val"
                          class="px-2.5 py-0.5 bg-[#EDEDEE] text-xs text-[#1E212B] rounded-full"
                        >{{ val }}</span>
                      </div>
                    </div>
                    <button
                      class="text-sm font-medium text-[#1E212B] px-3 py-1.5 border border-[#DBDBDD] rounded-lg hover:bg-[#F6F6F7] transition-colors shrink-0"
                      @click="option.isEditing = true"
                    >Edit</button>
                  </div>

                  <!-- Expanded -->
                  <div v-else class="border border-[#EDEDEE] rounded-xl p-4 bg-[#FAFAFA]">
                    <p class="text-xs font-medium text-[#6F7177] mb-2 flex items-center gap-1">
                      Option name
                      <svg class="w-3.5 h-3.5 text-[#94BDD8]" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <circle cx="12" cy="12" r="10" /><path d="M12 16v-4M12 8h.01" />
                      </svg>
                    </p>
                    <div class="flex items-center gap-2 mb-5">
                      <div class="text-[#C0C0C0] cursor-grab shrink-0">
                        <HugeiconsIcon :icon="DragDropHorizontalIcon" :size="20" color="currentColor" />
                      </div>
                      <GroBasicInput v-model="option.name" placeholder="e.g., Color, Size" class="flex-1" />
                      <button
                        class="w-8 h-8 flex items-center justify-center text-[#939499] hover:text-[#AF513A] transition-colors shrink-0"
                        @click="removeOption(option.id)"
                      >
                        <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                          <path d="M3 6h18M8 6V4h8v2M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6M10 11v6M14 11v6" />
                        </svg>
                      </button>
                    </div>

                    <p class="text-xs font-medium text-[#6F7177] mb-2 flex items-center gap-1">
                      Choices *
                      <svg class="w-3.5 h-3.5 text-[#94BDD8]" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <circle cx="12" cy="12" r="10" /><path d="M12 16v-4M12 8h.01" />
                      </svg>
                    </p>

                    <!-- Color picker -->
                    <template v-if="option.name.toLowerCase() === 'color'">
                      <GroColorPicker
                        v-model="option.values"
                        :attribute-values="getAttributeValues(option.name)"
                      />
                    </template>

                    <!-- Option picker for attributes with predefined values -->
                    <template v-else-if="hasAttributeValues(option.name)">
                      <GroOptionPicker
                        v-model="option.values"
                        :option-name="option.name"
                        :attribute-values="getAttributeValueStrings(option.name)"
                      />
                    </template>

                    <!-- Standard value list -->
                    <template v-else>
                      <div class="flex flex-col gap-2">
                        <div
                          v-for="(val, idx) in option.values"
                          :key="idx"
                          class="flex items-center gap-2"
                          draggable="true"
                          @dragstart="onValueDragStart(option.id, idx)"
                          @dragover.prevent
                          @drop.prevent="onValueDrop(option.id, idx)"
                        >
                          <div class="text-[#C0C0C0] cursor-grab shrink-0">
                            <HugeiconsIcon :icon="DragDropHorizontalIcon" :size="20" color="currentColor" />
                          </div>
                          <GroBasicInput v-model="option.values[idx]" class="flex-1" />
                          <button
                            class="w-8 h-8 flex items-center justify-center text-[#939499] hover:text-[#AF513A] transition-colors shrink-0"
                            @click="removeValueFromOption(option, idx)"
                          >
                            <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                              <path d="M3 6h18M8 6V4h8v2M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6M10 11v6M14 11v6" />
                            </svg>
                          </button>
                        </div>

                        <div class="flex items-center gap-2">
                          <div class="text-[#D9D9D9] shrink-0">
                            <HugeiconsIcon :icon="DragDropHorizontalIcon" :size="20" color="currentColor" />
                          </div>
                          <input
                            v-model="option.newValue"
                            placeholder="Add another value"
                            class="flex-1 px-4 py-2.5 bg-white border border-dashed border-[#DBDBDD] rounded-xl text-sm text-[#1E212B] outline-none focus:border-solid focus:border-[#94BDD8] placeholder-[#939499]"
                            @keydown.enter.prevent="addValueToOption(option)"
                            @keydown="(e: KeyboardEvent) => e.key === ',' && (e.preventDefault(), addValueToOption(option))"
                          >
                          <div class="w-8 h-8 shrink-0" />
                        </div>
                      </div>

                      <p class="text-xs text-[#939499] mt-2 ml-7">Press Enter or add a comma after each choice</p>
                    </template>
                    <button
                      class="mt-3 ml-7 px-4 py-1.5 border border-[#DBDBDD] rounded-lg text-sm font-medium text-[#1E212B] hover:bg-white transition-colors"
                      @click="doneEditingOption(option)"
                    >Done</button>
                  </div>
                </div>
              </div>

              <!-- Add option: dropdown when category has attributes, plain button otherwise -->
              <div ref="addOptionDropdownRef" class="relative mt-4 inline-block">
                <button
                  class="flex items-center gap-1.5 cursor-pointer text-sm font-medium text-[#2176AE] hover:underline disabled:opacity-40 disabled:cursor-not-allowed"
                  :disabled="productOptions.length >= 3"
                  @click="hasCategoryAttributes ? (showAddOptionDropdown = !showAddOptionDropdown) : addOption('')"
                >
                  <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                    <path d="M12 5v14M5 12h14" />
                  </svg>
                  Add {{ productOptions.length > 0 ? 'another ' : '' }}option
                </button>

                <Transition
                  enter-active-class="transition ease-out duration-100"
                  enter-from-class="opacity-0 scale-95"
                  enter-to-class="opacity-100 scale-100"
                  leave-active-class="transition ease-in duration-75"
                  leave-from-class="opacity-100 scale-100"
                  leave-to-class="opacity-0 scale-95"
                >
                  <div
                    v-if="showAddOptionDropdown && hasCategoryAttributes"
                    class="absolute left-0 top-full mt-1.5 w-48 bg-white border border-[#EDEDEE] rounded-xl shadow-lg z-20 py-1 origin-top-left"
                  >
                    <button
                      v-for="attr in availableCategoryAttributes"
                      :key="attr.id"
                      type="button"
                      class="w-full text-left px-3 py-2 text-sm text-[#1E212B] hover:bg-[#F6F6F7] transition-colors"
                      @click="addOption(attr.name)"
                    >{{ attr.name }}</button>
                    <div class="border-t border-[#EDEDEE] my-1" />
                    <button
                      type="button"
                      class="w-full text-left px-3 py-2 text-sm text-[#2176AE] font-medium hover:bg-[#F6F6F7] transition-colors flex items-center gap-1.5"
                      @click="addOption('')"
                    >
                      <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                        <path d="M12 5v14M5 12h14" />
                      </svg>
                      Add custom option
                    </button>
                  </div>
                </Transition>
              </div>

              <!-- Variants table -->
              <div v-if="hasVariants" class="border-t border-[#EDEDEE] mt-6 pt-5">
                <!-- Header -->
                <div class="flex items-center justify-between mb-3">
                  <div class="flex items-center gap-2">
                    <h4 class="text-base font-semibold text-[#1E212B]">Variants</h4>
                    <span class="text-xs text-[#6F7177] bg-[#EDEDEE] px-2 py-0.5 rounded-full">{{ variants.length }}/2048</span>
                  </div>
                  <!-- Group by dropdown -->
                  <div ref="groupByDropdownRef" class="relative">
                    <button
                      type="button"
                      class="flex items-center gap-1 text-xs text-[#939499] hover:text-[#1E212B] transition-colors"
                      @click="showGroupByDropdown = !showGroupByDropdown"
                    >
                      Group by:
                      <span class="text-[#1E212B] font-medium">{{ groupByOption?.name ?? '—' }}</span>
                      <svg class="w-3 h-3 ml-0.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                        <path d="M6 9l6 6 6-6" />
                      </svg>
                    </button>
                    <div
                      v-if="showGroupByDropdown"
                      class="absolute right-0 top-full mt-1 w-36 bg-white border border-[#EDEDEE] rounded-xl shadow-lg z-20 py-1"
                    >
                      <button
                        v-for="(opt, i) in productOptions"
                        :key="opt.id"
                        type="button"
                        class="w-full text-left px-3 py-2 text-sm hover:bg-[#F6F6F7] transition-colors"
                        :class="i === groupByOptionIndex ? 'text-[#1E212B] font-semibold' : 'text-[#4B4D55]'"
                        @click="groupByOptionIndex = i; showGroupByDropdown = false"
                      >{{ opt.name }}</button>
                    </div>
                  </div>
                </div>

                <div class="border border-[#EDEDEE] rounded-xl overflow-hidden">
                  <!-- Column header -->
                  <div class="flex items-center gap-3 px-4 py-2.5 bg-[#F6F6F7] border-b border-[#EDEDEE]">
                    <input type="checkbox" :checked="isAllVariantsSelected" class="rounded cursor-pointer accent-[#1E212B] flex-shrink-0" @change="toggleAllVariants">
                    <div class="flex-1 flex items-center gap-1.5 min-w-0">
                      <span class="text-xs font-semibold text-[#1E212B]">Variant</span>
                      <span class="text-[#DBDBDD] text-xs">·</span>
                      <button type="button" class="text-xs text-[#2176AE] hover:underline" @click="toggleAllGroups">
                        {{ allGroupsCollapsed ? 'Expand all' : 'Collapse all' }}
                      </button>
                    </div>
                    <div class="w-32 text-xs font-semibold text-[#1E212B] flex-shrink-0">Price</div>
                    <div class="w-24 text-xs font-semibold text-[#1E212B] flex-shrink-0">Available</div>
                  </div>

                  <!-- Groups -->
                  <div v-for="group in groupedVariants" :key="group.value" class="border-b border-[#EDEDEE] last:border-0">
                    <!-- Group row -->
                    <div class="flex items-center gap-3 px-4 py-3">
                      <input
                        type="checkbox"
                        :checked="isGroupSelected(group.items)"
                        class="rounded cursor-pointer accent-[#1E212B] flex-shrink-0"
                        :indeterminate="isGroupPartiallySelected(group.items)"
                        @change="toggleGroupSelection(group.items)"
                      >
                      <!-- Group image upload -->
                      <label
                        class="w-8 h-8 rounded-lg border border-dashed border-[#DBDBDD] bg-[#F6F6F7] flex items-center justify-center flex-shrink-0 cursor-pointer hover:border-[#94BDD8] hover:bg-[#EFF6FF] transition-colors overflow-hidden group"
                        :title="getGroupImage(group.value).imagePreview ? 'Change group image' : 'Add group image'"
                      >
                        <img v-if="getGroupImage(group.value).imagePreview" :src="getGroupImage(group.value).imagePreview" class="w-full h-full object-cover" alt="">
                        <svg v-else class="w-4 h-4 text-[#939499] group-hover:text-[#2176AE] transition-colors" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                          <rect x="3" y="3" width="18" height="18" rx="2" /><circle cx="8.5" cy="8.5" r="1.5" /><path d="M21 15l-5-5L5 21" />
                        </svg>
                        <input type="file" accept="image/*" class="hidden" @change="(e) => { const f = (e.target as HTMLInputElement).files?.[0]; if (f) onGroupImagePick(group.value, f); (e.target as HTMLInputElement).value = '' }">
                      </label>
                      <div class="flex-1 min-w-0">
                        <p class="text-sm font-semibold text-[#1E212B]">{{ group.value }}</p>
                        <button
                          type="button"
                          class="flex items-center gap-1 text-xs text-[#939499] mt-0.5 hover:text-[#1E212B] transition-colors"
                          @click="toggleGroup(group.value)"
                        >
                          {{ group.items.length }} {{ group.items.length === 1 ? 'variant' : 'variants' }}
                          <svg
                            class="w-3 h-3 transition-transform"
                            :class="expandedGroups.has(group.value) ? 'rotate-180' : ''"
                            viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"
                          >
                            <path d="M6 9l6 6 6-6" />
                          </svg>
                        </button>
                      </div>
                      <!-- Group price (sets all variants) -->
                      <div class="w-32 flex items-center bg-[#F6F6F7] border border-[#EDEDEE] rounded-lg focus-within:border-[#1E212B] transition-colors flex-shrink-0">
                        <span class="pl-3 text-sm text-[#4B4D55]">{{ currencySymbol }}</span>
                        <input
                          :value="getGroupPrice(group.items)"
                          type="number"
                          placeholder="0.00"
                          class="flex-1 px-2 py-1.5 bg-transparent text-sm text-[#1E212B] outline-none placeholder-[#939499]"
                          @input="setGroupPrice(group.items, ($event.target as HTMLInputElement).value)"
                        >
                      </div>
                      <!-- Group qty (readonly sum) -->
                      <input
                        :value="getGroupQty(group.items) || undefined"
                        type="number"
                        disabled
                        placeholder="0"
                        class="w-24 px-3 py-1.5 bg-[#F6F6F7] border border-[#EDEDEE] rounded-lg text-sm text-[#939499] outline-none cursor-not-allowed flex-shrink-0"
                      >
                    </div>

                    <!-- Variant rows (accordion) -->
                    <template v-if="expandedGroups.has(group.value)">
                      <div
                        v-for="variant in group.items"
                        :key="variant.id"
                        class="flex items-center gap-3 px-4 py-2.5 border-t border-[#EDEDEE] bg-[#FAFAFA]"
                        style="padding-left: 3.25rem"
                      >
                        <input type="checkbox" :checked="selectedVariants.includes(variant.id)" class="rounded cursor-pointer accent-[#1E212B] flex-shrink-0" @change="toggleVariant(variant.id)">
                        <!-- Variant image upload (falls back to group image) -->
                        <label
                          class="w-8 h-8 rounded-lg border border-dashed border-[#DBDBDD] bg-[#F6F6F7] flex items-center justify-center flex-shrink-0 cursor-pointer hover:border-[#94BDD8] hover:bg-[#EFF6FF] transition-colors overflow-hidden group"
                          :title="variant.imagePreview ? 'Change image' : 'Add image'"
                        >
                          <img
                            v-if="variant.imagePreview || getGroupImage(group.value).imagePreview"
                            :src="variant.imagePreview || getGroupImage(group.value).imagePreview"
                            class="w-full h-full object-cover"
                            :class="!variant.imagePreview ? 'opacity-40' : ''"
                            alt=""
                          >
                          <svg v-else class="w-4 h-4 text-[#939499] group-hover:text-[#2176AE] transition-colors" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                            <rect x="3" y="3" width="18" height="18" rx="2" /><circle cx="8.5" cy="8.5" r="1.5" /><path d="M21 15l-5-5L5 21" />
                          </svg>
                          <input type="file" accept="image/*" class="hidden" @change="(e) => { const f = (e.target as HTMLInputElement).files?.[0]; if (f) onVariantImagePick(variant, f); (e.target as HTMLInputElement).value = '' }">
                        </label>
                        <span class="flex-1 text-sm text-[#1E212B] min-w-0 break-words">{{ variant.name }}</span>
                        <!-- Price -->
                        <div class="w-32 flex items-center bg-white border border-[#EDEDEE] rounded-lg focus-within:border-[#1E212B] transition-colors flex-shrink-0">
                          <span class="pl-3 text-sm text-[#4B4D55]">{{ currencySymbol }}</span>
                          <input v-model="variant.price" type="number" placeholder="0.00" class="flex-1 px-2 py-1.5 bg-transparent text-sm text-[#1E212B] outline-none placeholder-[#939499]">
                        </div>
                        <!-- Qty -->
                        <input
                          v-model="variant.quantity"
                          type="number"
                          placeholder="0"
                          class="w-24 px-3 py-1.5 bg-white border border-[#EDEDEE] rounded-lg text-sm text-[#1E212B] outline-none focus:border-[#1E212B] transition-colors placeholder-[#939499] flex-shrink-0"
                        >
                      </div>
                    </template>
                  </div>
                </div>
              </div>
            </GroFormSection>

            <!-- Additional Info -->
            <GroFormSection title="Additional Info" :open="openSections.additional" :bordered="false" @toggle="toggleSection('additional')">
              <p v-if="additionalFields.length === 0" class="text-sm text-[#6F7177] mb-3">
                Add custom fields that relate with your customer or product
              </p>

              <CustomFormField
                v-if="additionalFields.length > 0"
                v-model="additionalAnswers"
                :fields="additionalFields"
                :editable-labels="true"
                class="mb-3"
                @label-changed="handleAdditionalLabelChanged"
              />

              <div class="relative">
                <button
                  class="flex items-center gap-1.5 text-sm font-medium text-[#2176AE] hover:underline"
                  @click="showFieldPicker = !showFieldPicker"
                >
                  <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                    <path d="M12 5v14M5 12h14" />
                  </svg>
                  Add custom fields
                </button>
                <div
                  v-if="showFieldPicker"
                  class="absolute left-0 top-7 z-20 bg-white border border-[#EDEDEE] rounded-xl shadow-lg py-1 min-w-[160px]"
                >
                  <button
                    v-for="opt in ADDITIONAL_FIELD_TYPES"
                    :key="opt.type"
                    type="button"
                    class="w-full text-left px-4 py-2 text-sm text-[#1E212B] hover:bg-[#F6F6F7] transition-colors"
                    @click="handleAddAdditionalField(opt.type)"
                  >{{ opt.label }}</button>
                </div>
              </div>
            </GroFormSection>

          </div>

          <!-- Right column -->
          <div class="col-span-5 md:col-span-2 pt-3 gap-4">
            <div class="py-6 border-b border-[#EDEDEE] px-8">
              <h3 class="text-base font-semibold text-[#1E212B] mb-3">Product status</h3>
              <select
                v-model="form.status"
                class="w-full px-3 py-2 bg-[#F6F6F7] border border-[#EDEDEE] rounded-lg text-sm text-[#1E212B] outline-none hover:border-[#94BDD8] focus:border-[#1E212B] mb-2 cursor-pointer"
              >
                <option value="Active">Active</option>
                <option value="Draft">Draft</option>
                <option value="Archived">Archived</option>
              </select>
              <p class="text-xs text-[#6F7177]">This product will be available to your sales channels.</p>
              <div class="mt-3">
                <GroToggle v-model="form.showInStore" :disabled="form.status.toLowerCase() !== 'active'">
                  Show in your store
                </GroToggle>
                <button class="flex items-center gap-1.5 text-sm font-medium text-[#2176AE] hover:underline mt-3">
                  <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                    <path d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                  </svg>
                  Preview
                </button>
              </div>
            </div>

            <div class="px-8 py-6">
              <h3 class="text-base font-semibold text-[#1E212B] mb-1">Tags</h3>
              <p class="text-xs text-[#939499] mb-3">Used for filtering and search</p>
              <CatalogTagSelect v-model="tags" />
            </div>
          </div>
        </div>

        <!-- Inline footer -->
        <div
          ref="footerSentinel"
          class="px-8 py-5 border-t border-[#EDEDEE] flex items-center justify-between"
          :class="showFixedBar ? 'invisible' : ''"
        >
          <button
            class="text-sm font-medium text-[#2176AE] hover:underline transition-colors"
            @click="saveProduct('draft')"
          >
            Save as Draft
          </button>
          <div class="flex items-center gap-3">
            <GroBasicButton color="secondary" size="xs" shape="custom" class="w-max" @click="router.push('/catalog/products')">
              Discard
            </GroBasicButton>
            <GroBasicButton color="primary" size="xs" shape="custom" class="w-max" :disabled="isSaving" @click="saveProduct(form.status.toLowerCase())">
              {{ isSaving ? 'Saving...' : (isNew ? 'Save' : 'Update') }}
            </GroBasicButton>
          </div>
        </div>
      </div>

      <!-- Fixed bottom bar -->
      <Transition
        enter-active-class="transition ease-out duration-200"
        enter-from-class="opacity-0 translate-y-2"
        enter-to-class="opacity-100 translate-y-0"
        leave-active-class="transition ease-in duration-150"
        leave-from-class="opacity-100"
        leave-to-class="opacity-0 translate-y-2"
      >
        <div
          v-if="showFixedBar && !isLoading"
          class="fixed bottom-0 left-0 md:left-16 lg:left-64 right-0 z-50 bg-white border-t border-[#EDEDEE] px-8 py-4 flex items-center justify-between"
        >
          <button
            class="text-sm font-medium text-[#2176AE] hover:underline transition-colors"
            @click="saveProduct('draft')"
          >
            Save as Draft
          </button>
          <div class="flex items-center gap-3">
            <GroBasicButton color="secondary" size="xs" shape="custom" class="w-max" @click="router.push('/catalog/products')">
              Discard
            </GroBasicButton>
            <GroBasicButton color="primary" size="xs" shape="custom" class="w-max" :disabled="isSaving" @click="saveProduct(form.status.toLowerCase())">
              {{ isSaving ? 'Saving...' : (isNew ? 'Save' : 'Update') }}
            </GroBasicButton>
          </div>
        </div>
      </Transition>
    </template>
  </MainDashboard>
</template>

<style>
button { cursor: pointer; }
</style>

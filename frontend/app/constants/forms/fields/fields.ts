// Email Field Creator
export const CREATE_EMAIL_FIELDS = (fieldCount: number) => {
  return {
    SK: `FORM#${crypto.randomUUID()}`,
    fieldTypeID: 'email',
    fieldTypeName: 'Email',
    section: 'Extra fields',
    slug: `email_${fieldCount}`,
    label: `Email ${fieldCount}`,
    description: 'Enter a valid email address',
    required: false,
    order: fieldCount,
    options: [],
    configs: [],
    validations: []
  }
}

// Phone Field Creator
export const CREATE_PHONE_FIELDS = (fieldCount: number) => {
  return {
    SK: `FORM#${crypto.randomUUID()}`,
    fieldTypeID: 'phone',
    fieldTypeName: 'Phone Number',
    section: 'Extra fields',
    slug: `phone_${fieldCount}`,
    label: `Phone ${fieldCount}`,
    description: 'Enter a valid phone number',
    required: false,
    order: fieldCount,
    options: [],
    configs: [],
    validations: []
  }
}

// Address Field Creator — returns a group of 5 fields
export const CREATE_ADDRESS_FIELDS = (fieldCount: number) => {
  const base = fieldCount * 10
  return [
    {
      SK: `FORM#${crypto.randomUUID()}`,
      fieldTypeID: 'country',
      fieldTypeName: 'Country',
      section: 'Extra fields',
      slug: `country_${fieldCount}`,
      label: 'Country',
      description: 'Select a country',
      required: false,
      order: base + 1,
      options: [],
      configs: [],
      validations: []
    },
    {
      SK: `FORM#${crypto.randomUUID()}`,
      fieldTypeID: 'state',
      fieldTypeName: 'State',
      section: 'Extra fields',
      slug: `state_province_${fieldCount}`,
      label: 'State/Province',
      description: '',
      required: false,
      order: base + 2,
      options: [],
      configs: [{ key: 'countrySlug', fieldSlug: `country_${fieldCount}` }],
      validations: []
    },
    {
      SK: `FORM#${crypto.randomUUID()}`,
      fieldTypeID: 'street_address',
      fieldTypeName: 'Short Text',
      section: 'Extra fields',
      slug: `street_address_${fieldCount}`,
      label: 'Street Address',
      description: '',
      required: false,
      order: base + 3,
      options: [],
      configs: [],
      validations: []
    },
    {
      SK: `FORM#${crypto.randomUUID()}`,
      fieldTypeID: 'city',
      fieldTypeName: 'City',
      section: 'Extra fields',
      slug: `city_${fieldCount}`,
      label: 'City',
      description: '',
      required: false,
      order: base + 4,
      options: [],
      configs: [],
      validations: []
    },
    {
      SK: `FORM#${crypto.randomUUID()}`,
      fieldTypeID: 'zip',
      fieldTypeName: 'Short Text',
      section: 'Extra fields',
      slug: `zip_postal_${fieldCount}`,
      label: 'Zip/Postal Code',
      description: '',
      required: false,
      order: base + 5,
      options: [],
      configs: [],
      validations: []
    },
  ]
}

// Company Info Field Creator
export const CREATE_COMPANY_FIELDS = (fieldCount: number) => {
  return {
    SK: `FORM#${crypto.randomUUID()}`,
    fieldTypeID: 'company',
    fieldTypeName: 'Short Text',
    section: 'Extra fields',
    slug: `company_${fieldCount}`,
    label: `Company ${fieldCount}`,
    description: 'Enter company name',
    required: false,
    order: fieldCount,
    options: [],
    configs: [],
    validations: [
      {
        key: 'minLength',
        rule: '2',
        message: 'Company name must be at least 2 characters'
      }
    ]
  }
}

// Position Field Creator
export const CREATE_POSITION_FIELDS = (fieldCount: number) => {
  return {
    SK: `FORM#${crypto.randomUUID()}`,
    fieldTypeID: 'position',
    fieldTypeName: 'Short Text',
    section: 'Extra fields',
    slug: `position_${fieldCount}`,
    label: `Position ${fieldCount}`,
    description: 'Enter job position/title',
    required: false,
    order: fieldCount,
    options: [],
    configs: [],
    validations: [
      {
        key: 'minLength',
        rule: '2',
        message: 'Position must be at least 2 characters'
      },
      {
        key: 'maxLength',
        rule: '100',
        message: 'Position must not exceed 100 characters'
      }
    ]
  }
}

// Birthdate Field Creator
export const CREATE_BIRTHDATE_FIELDS = (fieldCount: number) => {
  return {
    SK: `FORM#${crypto.randomUUID()}`,
    fieldTypeID: 'birthdate',
    fieldTypeName: 'Date',
    section: 'Extra fields',
    slug: `birthdate_${fieldCount}`,
    label: `Date ${fieldCount}`,
    description: 'Select a date',
    required: false,
    order: fieldCount,
    options: [],
    configs: [],
    validations: [
      {
        key: 'pattern',
        rule: '^\\d{4}-\\d{2}-\\d{2}$',
        message: 'Please enter a valid date in YYYY-MM-DD format'
      }
    ]
  }
}

// Notes Field Creator
export const CREATE_NOTES_FIELDS = (fieldCount: number) => {
  return {
    SK: `FORM#${crypto.randomUUID()}`,
    fieldTypeID: 'notes',
    fieldTypeName: 'Long Text',
    section: 'Extra fields',
    slug: `notes_${fieldCount}`,
    label: `Notes ${fieldCount}`,
    description: 'Add any additional notes',
    required: false,
    order: fieldCount,
    options: [],
    configs: [],
    validations: [
      {
        key: 'maxLength',
        rule: '500',
        message: 'Notes must not exceed 500 characters'
      }
    ]
  }
}

// Custom Field Creator
export const CREATE_CUSTOM_FIELDS = (fieldCount: number) => {
  return {
    SK: `FORM#${crypto.randomUUID()}`,
    fieldTypeID: 'custom',
    fieldTypeName: 'Short Text',
    section: 'Extra fields',
    slug: `custom_field_${fieldCount}`,
    label: `Custom Field ${fieldCount}`,
    description: 'Enter custom information',
    required: false,
    order: fieldCount,
    options: [],
    configs: [],
    validations: [
      {
        key: 'maxLength',
        rule: '200',
        message: 'Custom field must not exceed 200 characters'
      }
    ]
  }
}

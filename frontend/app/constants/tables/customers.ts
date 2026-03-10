// Email Field Creator
import { type ColsDefinition } from "@/composables/helpers/types/table";

export const LIST_CUSTOMER_COLUMNS: ColsDefinition[] = [
  {
    field: 'name',
    title: 'Customer Name',
    type: 'string'
  },
  {
    field: 'customer_type',
    title: 'Customer Type',
    type: 'string',
    isUnique: false
  },
  {
    field: 'email',
    title: 'Email',
    type: 'string'
  },
  {
    field: 'phone_number',
    title: 'Phone Number',
    type: 'string'
  },
  { field: 'actions', title: 'Actions' },
];

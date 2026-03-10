import {
  AnalyticsUpIcon,
  BankIcon, GlobalIcon,
  Home05Icon,
  Invoice01Icon,
  MagnetIcon, Settings01Icon,
  Store01Icon, Store04Icon,
  User03Icon
} from "@hugeicons/core-free-icons";

type IconSvgObject = ([string, {
  [key: string]: string | number;
}])[] | readonly (readonly [string, {
  readonly [key: string]: string | number;
}])[];

export interface SidebarMenuItem {
  title: string;
  icon: IconSvgObject;
  url?: string;
  count: number | null;
  soon: boolean;
}

export const SIDEBAR_MENU_SECTION_A: Array<SidebarMenuItem> = [
  {
    title: 'Home',
    url: 'dashboard',
    icon: Home05Icon,
    count: null,
    soon: false
  },
  {
    title: 'Invoices',
    url: 'invoices',
    icon: Invoice01Icon,
    count: null,
    soon: false,
  },
  {
    title: 'Customers',
    url: 'customers',
    icon: User03Icon,
    count: null,
    soon: false,
  },
  {
    title: 'Leads Management',
    url: 'dashboard',
    icon: MagnetIcon,
    count: null,
    soon: false,
  },
  {
    title: 'Payments',
    url: 'dashboard',
    icon: BankIcon,
    count: null,
    soon: false,
  },
  {
    title: 'Analytics',
    url: 'dashboard',
    icon: AnalyticsUpIcon,
    count: null,
    soon: false,
  },
  {
    title: 'Product & Services',
    url: 'dashboard',
    icon: Store01Icon,
    count: null,
    soon: false,
  }
]

export const SIDEBAR_MENU_SECTION_B: Array<SidebarMenuItem> = [
  {
    title: 'My Store',
    url: 'dashboard',
    icon: GlobalIcon,
    count: null,
    soon: true
  },
  {
    title: 'POS',
    url: 'dashboard',
    icon: Store04Icon,
    count: null,
    soon: true
  },
]

export const SIDEBAR_MENU_SECTION_C: Array<SidebarMenuItem> = [
  {
    title: 'Settings',
    url: 'dashboard',
    icon: Settings01Icon,
    count: null,
    soon: false
  },
]

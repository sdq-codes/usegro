export interface ColsDefinition {
  field: string;
  title: string;
  isUnique?: boolean;
  type?: 'string' | 'number' | 'date' | 'bool' | 'array';
}

export type SubmissionStatus = 'active' | 'archived';

export type SubmissionType =
  | 'form'
  | 'customer'
  | 'invoice-template'
  | 'invoice';

export interface FormSubmission {
  _id: string;
  formID: string;
  crmID: string;
  formVersionID: string;
  status: SubmissionStatus;
  type: SubmissionType;
  answers: Record<string, unknown>;
  versionSnap: unknown;
  createdAt: string;
  archivedAt?: string | null;
}

export type SubmissionStatus = 'active' | 'archived';

export type SubmissionType =
  | 'form'
  | 'customer'
  | 'invoice-template'
  | 'invoice';

export interface FormSubmission {
  PK: string; // FORM#<formID>
  SK: string; // SUBMISSION#<submissionID>
  formID: string;
  crmID: string;
  formVersionID: string;
  submissionID: string;
  status: SubmissionStatus;
  type: SubmissionType;
  answers: Record<string, unknown>;
  versionSnap: unknown;
  createdAt: string; // ISO timestamp
  archivedAt?: string | null; // optional, may be null
}

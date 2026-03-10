export default interface CSVData {
  data: Record<string, unknown>[];
  meta: {
    fields: string[];
  };
}

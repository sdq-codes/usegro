import type {FormVersionResponse} from "@/composables/dto/customer/form/form";

type ImportedCustomer = Record<string, unknown>;

type CreateCustomerPayload = {
  type: string;
  answers: Record<string, unknown>;
  versionSnap: unknown[];
};

export function buildCustomerPayloads(
  formVersionData: FormVersionResponse,
  importedCustomers: ImportedCustomer[],
  sharedTags: string[]
): CreateCustomerPayload[] {
  const fields = formVersionData.fields;

  return importedCustomers.map((customer) => {
    const answers: Record<string, unknown> = {};

    for (const field of fields) {
      const slug = field.slug;

      // Skip if field doesn't exist in imported data
      if (!(slug in customer)) continue;

      // Handle showIf configs (conditional visibility)
      const showIfConfig = field.configs?.find(
        (c) => c.key === "showIf"
      );

      if (showIfConfig) {
        const dependsOn = showIfConfig.fieldSlug;
        const expectedValue = String(showIfConfig.fieldValue).toLowerCase();
        const actualValue = String(customer[dependsOn] || "").toLowerCase();

        if (expectedValue !== actualValue) {
          continue; // Skip this field
        }
      }

      let value = customer[slug];

      // Handle checkbox values like "[yes, no]" or "[ no, no]"
      if (field.fieldTypeName === "Checkbox" && typeof value === "string") {
        try {
          // Parse safely if it's like "[yes, no]"
          const parsed = value
            .replace(/[[]]/g, "")
            .split(",")
            .map((v) => v.trim().toLowerCase())
            .filter((v) => v === "yes" || v === "subscribe_marketing_email" || v === "subscribe_sms");

          // Map to real values from form options
          const options = field.options?.map((opt: unknown) => opt.value) || [];
          value = options.filter((opt: string) =>
            parsed.some((p) => opt.toLowerCase().includes(p))
          );
        } catch {
          value = [];
        }
      }

      answers[slug] = value;
    }

    // Add shared tags
    answers["customer_tags"] = sharedTags;

    return {
      type: "customer",
      answers,
      versionSnap: fields, // for your API
    };
  });
}

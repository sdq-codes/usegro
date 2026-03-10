import { describe, it, expect } from 'vitest';
import { buildCustomerPayloads } from './buildCustomerImportPayloads';
import type { FormVersionResponse } from '@/composables/dto/customer/form/form';

describe('buildCustomerPayloads', () => {
  const mockFormVersionData: FormVersionResponse = {
    fields: [
      {
        slug: 'first_name',
        fieldTypeName: 'Text',
        configs: [],
        options: [],
      },
      {
        slug: 'email',
        fieldTypeName: 'Email',
        configs: [],
        options: [],
      },
      {
        slug: 'subscribe_marketing_email',
        fieldTypeName: 'Checkbox',
        configs: [],
        options: [
          { value: 'subscribe_marketing_email', label: 'Yes' },
        ],
      },
    ],
    version: {} as unknown,
  } as FormVersionResponse;

  it('should build payloads for imported customers', () => {
    const importedCustomers = [
      {
        first_name: 'John',
        email: 'john@example.com',
      },
      {
        first_name: 'Jane',
        email: 'jane@example.com',
      },
    ];
    const sharedTags = ['import', 'batch-1'];

    const result = buildCustomerPayloads(
      mockFormVersionData,
      importedCustomers,
      sharedTags
    );

    expect(result).toHaveLength(2);
    expect(result[0]).toEqual({
      type: 'customer',
      answers: {
        first_name: 'John',
        email: 'john@example.com',
        customer_tags: ['import', 'batch-1'],
      },
      versionSnap: mockFormVersionData.fields,
    });
    expect(result[1].answers.first_name).toBe('Jane');
  });

  it('should skip fields not in imported data', () => {
    const importedCustomers = [
      {
        first_name: 'John',
        // email is missing
      },
    ];
    const sharedTags: string[] = [];

    const result = buildCustomerPayloads(
      mockFormVersionData,
      importedCustomers,
      sharedTags
    );

    expect(result[0].answers).toEqual({
      first_name: 'John',
      customer_tags: [],
    });
    expect(result[0].answers).not.toHaveProperty('email');
  });

  it('should handle conditional fields with showIf config', () => {
    const formData: FormVersionResponse = {
      fields: [
        {
          slug: 'has_business',
          fieldTypeName: 'Radio',
          configs: [],
          options: [],
        },
        {
          slug: 'business_name',
          fieldTypeName: 'Text',
          configs: [
            {
              key: 'showIf',
              fieldSlug: 'has_business',
              fieldValue: 'yes',
            },
          ],
          options: [],
        },
      ],
      version: {} as unknown,
    } as FormVersionResponse;

    const importedCustomers = [
      {
        has_business: 'yes',
        business_name: 'Acme Corp',
      },
      {
        has_business: 'no',
        business_name: 'Should be ignored',
      },
    ];

    const result = buildCustomerPayloads(formData, importedCustomers, []);

    expect(result[0].answers).toHaveProperty('business_name', 'Acme Corp');
    expect(result[1].answers).not.toHaveProperty('business_name');
  });

  it('should handle checkbox values as strings', () => {
    const importedCustomers = [
      {
        first_name: 'John',
        subscribe_marketing_email: '[yes, no]',
      },
    ];

    const result = buildCustomerPayloads(
      mockFormVersionData,
      importedCustomers,
      []
    );

    expect(result[0].answers.subscribe_marketing_email).toEqual([
      'subscribe_marketing_email',
    ]);
  });

  it('should handle checkbox values with subscribe_sms', () => {
    const formData: FormVersionResponse = {
      fields: [
        {
          slug: 'subscriptions',
          fieldTypeName: 'Checkbox',
          configs: [],
          options: [
            { value: 'subscribe_marketing_email', label: 'Email' },
            { value: 'subscribe_sms', label: 'SMS' },
          ],
        },
      ],
      version: {} as unknown,
    } as FormVersionResponse;

    const importedCustomers = [
      {
        subscriptions: '[subscribe_marketing_email, subscribe_sms]',
      },
    ];

    const result = buildCustomerPayloads(formData, importedCustomers, []);

    expect(result[0].answers.subscriptions).toEqual([
      'subscribe_marketing_email',
      'subscribe_sms',
    ]);
  });

  it('should handle empty checkbox values', () => {
    const importedCustomers = [
      {
        first_name: 'John',
        subscribe_marketing_email: '[no, no]',
      },
    ];

    const result = buildCustomerPayloads(
      mockFormVersionData,
      importedCustomers,
      []
    );

    expect(result[0].answers.subscribe_marketing_email).toEqual([]);
  });

  it('should handle invalid checkbox values gracefully', () => {
    const importedCustomers = [
      {
        first_name: 'John',
        subscribe_marketing_email: 'invalid-format',
      },
    ];

    const result = buildCustomerPayloads(
      mockFormVersionData,
      importedCustomers,
      []
    );

    expect(result[0].answers.subscribe_marketing_email).toEqual([]);
  });

  it('should add shared tags to all customers', () => {
    const importedCustomers = [
      { first_name: 'John' },
      { first_name: 'Jane' },
    ];
    const sharedTags = ['vip', 'early-adopter'];

    const result = buildCustomerPayloads(
      mockFormVersionData,
      importedCustomers,
      sharedTags
    );

    expect(result[0].answers.customer_tags).toEqual(['vip', 'early-adopter']);
    expect(result[1].answers.customer_tags).toEqual(['vip', 'early-adopter']);
  });

  it('should handle empty imported customers array', () => {
    const result = buildCustomerPayloads(mockFormVersionData, [], []);

    expect(result).toEqual([]);
  });

  it('should handle case-insensitive showIf conditions', () => {
    const formData: FormVersionResponse = {
      fields: [
        {
          slug: 'country',
          fieldTypeName: 'Text',
          configs: [],
          options: [],
        },
        {
          slug: 'state',
          fieldTypeName: 'Text',
          configs: [
            {
              key: 'showIf',
              fieldSlug: 'country',
              fieldValue: 'USA',
            },
          ],
          options: [],
        },
      ],
      version: {} as unknown,
    } as FormVersionResponse;

    const importedCustomers = [
      {
        country: 'usa', // lowercase
        state: 'California',
      },
    ];

    const result = buildCustomerPayloads(formData, importedCustomers, []);

    expect(result[0].answers).toHaveProperty('state', 'California');
  });

  it('should include versionSnap in payload', () => {
    const importedCustomers = [{ first_name: 'John' }];

    const result = buildCustomerPayloads(
      mockFormVersionData,
      importedCustomers,
      []
    );

    expect(result[0].versionSnap).toEqual(mockFormVersionData.fields);
  });
});

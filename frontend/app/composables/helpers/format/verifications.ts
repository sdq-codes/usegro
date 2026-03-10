type Verification = {
  id: string;
  user_id: string;
  type: string;
  status: string;
  created_at: string;
  updated_at: string;
};

export const verifications = (verifications: Verification[]): Record<string, string> => {
  return verifications.reduce((acc: Record<string, string>, v) => {
    acc[v.type.toLowerCase()] = v.status;
    return acc;
  }, {});
}

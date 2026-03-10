export function maskEmail(email: string): string {
  const [localPart, domain] = email.split("@");
  if (!localPart) return email;

  const visible = localPart.slice(-2); // keep last 2 characters
  const masked = "*".repeat(localPart.length - visible.length) + visible;

  return `${masked}@${domain}`;
}

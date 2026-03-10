export function FormatZodErrors(errors: unknown): Record<string, string[]> {
  return errors.reduce((acc: Record<string, string[]>, err: unknown) => {
    const key = err.path[0]
    if (key != undefined) {
      if (!acc?.key) acc[key] = []
      acc[key]?.push(err?.message)
    }
    return acc
  }, {})
}

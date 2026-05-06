export function logError(ctx: string, err: unknown) {
  if (import.meta.env.DEV) {
    console.error(`[${ctx}]`, err)
  }
}

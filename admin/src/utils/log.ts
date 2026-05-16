/* eslint-disable functional/no-expression-statements */
export function logError(ctx: string, err: unknown): never {
  if (import.meta.env.DEV) {
    console.error(`[${ctx}]`, err);
  }
  return undefined as never;
}

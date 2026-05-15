/** Case-insensitive substring match. */
export function fuzzyMatch(text: string, query: string): boolean {
  return text.toLowerCase().includes(query.toLowerCase());
}

/** Returns HTML string with matched portion wrapped in <mark> tags. */
export function highlightMatch(text: string, query: string): string {
  if (!query) return text;
  const idx = text.toLowerCase().indexOf(query.toLowerCase());
  if (idx === -1) return text;
  const before = text.slice(0, idx);
  const match = text.slice(idx, idx + query.length);
  const after = text.slice(idx + query.length);
  return `${before}<mark class="bg-primary/20 text-text rounded-sm">${match}</mark>${after}`;
}
export function monasteryCatalogLink(id: number | string) {
  return `#readings?${new URLSearchParams({ monastery: String(id) })}`
}

export function monasteryIdFromHash(hash: string): number | undefined {
  const [page, query] = hash.split('?')
  const value =
    page === '#readings' ? new URLSearchParams(query).get('monastery') : null
  const id = Number(value)
  return Number.isSafeInteger(id) && id > 0 ? id : undefined
}

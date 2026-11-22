export function monasteryCatalogLink(id: number | string) {
  return `#readings?${new URLSearchParams({ monastery: String(id) })}`
}

export function monasteryIdFromHash(hash: string): number | undefined {
  const [page, query] = hash.split('?')
  // Keep previously shared catalog links working as well as the new pages.
  const value = /^#monasteries\/\d+$/.test(page)
    ? page.slice('#monasteries/'.length)
    : page === '#readings'
      ? new URLSearchParams(query).get('monastery')
      : null
  const id = Number(value)
  return Number.isSafeInteger(id) && id > 0 ? id : undefined
}

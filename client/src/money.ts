const formatter = new Intl.NumberFormat('mn-MN', {
  style: 'currency',
  currency: 'MNT',
  currencyDisplay: 'narrowSymbol',
  minimumFractionDigits: 0,
  maximumFractionDigits: 0,
})
export function formatPrice(price: number) {
  return formatter.format(price)
}

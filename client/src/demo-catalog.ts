import type { Book, Category, Paginated } from './api'

const prayers: Record<
  number,
  Omit<Book, 'id' | 'price'> & { categoryId: number }
> = {
  1: {
    name: 'Авидын чого',
    description: 'Талийгаачийн хойдын буянд зориулан айлтгах залбирал.',
    categoryId: 1,
  },
  2: {
    name: 'Арван гурван сан',
    description: 'Эд хөрөнгө, хишиг буяныг дуудан оршоох залбирал.',
    categoryId: 2,
  },
  3: {
    name: 'Банзрагч',
    description: 'Доройтол, бузар түйтгэрийг ариусгах ном.',
    categoryId: 3,
  },
  4: {
    name: 'Цагаан шүхэрт',
    description: 'Саад бэрхшээл, гай барцадаас хамгаалах залбирал.',
    categoryId: 3,
  },
  5: {
    name: 'Алтангэрэл',
    description: 'Амар амгалан, сайн сайхныг ерөөх ном.',
    categoryId: 2,
  },
  6: {
    name: 'Ерөөл',
    description: 'Буян хишиг, сайн сайхныг хүсэж айлтгах ерөөл.',
    categoryId: 2,
  },
}

const categories: Category[] = [
  { id: 1, name: 'Remembrance' },
  { id: 2, name: 'Blessings' },
  { id: 3, name: 'Protection' },
]

export type DemoMonastery = {
  id: string
  name: string
  offers: { bookId: number; price: number }[]
}
export type DemoBook = Book & { monasteryId: string }

export const demoMonasteries: DemoMonastery[] = [
  {
    id: 'amber-valley',
    name: 'Amber Valley Monastery',
    offers: [
      { bookId: 1, price: 10000 },
      { bookId: 2, price: 2000 },
      { bookId: 3, price: 2500 },
    ],
  },
  {
    id: 'cedar-ridge',
    name: 'Cedar Ridge Monastery',
    offers: [
      { bookId: 1, price: 12000 },
      { bookId: 4, price: 4000 },
      { bookId: 5, price: 3500 },
    ],
  },
  {
    id: 'lotus-grove',
    name: 'Lotus Grove Monastery',
    offers: [
      { bookId: 2, price: 3000 },
      { bookId: 4, price: 5000 },
      { bookId: 6, price: 2500 },
    ],
  },
]

export function monasteryCatalogLink(id: string) {
  return `#readings?${new URLSearchParams({ monastery: id })}`
}

export function monasteryIdFromHash(hash: string) {
  const [page, query] = hash.split('?')
  return page === '#readings'
    ? new URLSearchParams(query).get('monastery') || undefined
    : undefined
}

export function getDemoCategories(monastery: DemoMonastery) {
  const categoryIds = new Set(
    monastery.offers.map(({ bookId }) => prayers[bookId].categoryId),
  )
  return categories.filter(({ id }) => categoryIds.has(id))
}

export async function getDemoBooks(
  monastery: DemoMonastery,
  { categoryId, query }: { categoryId?: number; query: string },
): Promise<Paginated<DemoBook>> {
  const search = query.trim().toLowerCase()
  const data = monastery.offers
    .filter(({ bookId }) => {
      const prayer = prayers[bookId]
      return (
        (categoryId === undefined || prayer.categoryId === categoryId) &&
        [prayer.name, prayer.description].some((value) =>
          value.toLowerCase().includes(search),
        )
      )
    })
    .map(({ bookId, price }) => {
      const { name, description } = prayers[bookId]
      return { id: bookId, name, description, price, monasteryId: monastery.id }
    })
  return { data, pagination: { nextCursor: null, hasMore: false } }
}

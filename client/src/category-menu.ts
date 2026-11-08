type TriggerBounds = {
  left: number
  top: number
  bottom: number
  width: number
}
export function categoryMenuPosition(
  rect: TriggerBounds,
  viewport: { width: number; height: number },
) {
  const margin = 8
  const gap = 6
  const desired = 6 * 44 + 10
  const below = Math.max(0, viewport.height - rect.bottom - margin - gap)
  const above = Math.max(0, rect.top - margin - gap)
  const opensAbove = below < desired && above > below
  const width = Math.min(
    Math.max(rect.width, 220),
    Math.max(0, viewport.width - margin * 2),
  )
  return {
    left: Math.max(
      margin,
      Math.min(rect.left, viewport.width - width - margin),
    ),
    top: opensAbove ? rect.top - gap : rect.bottom + gap,
    width,
    maxHeight: Math.min(desired, opensAbove ? above : below),
    opensAbove,
  }
}

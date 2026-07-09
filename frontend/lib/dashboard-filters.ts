import { trustLevel, type TrustLevel } from './sop'
import type { Opportunity, OpportunityStatus, Platform, SopStage, SourceType } from './types'

export type SortKey = 'confidence' | 'trust' | 'newest' | 'oldest'
export type TimeRange = 'all' | 'today' | '3d' | '7d' | 'custom'

export interface DashboardFilters {
  status: 'all' | OpportunityStatus
  platform: 'all' | Platform
  source: 'all' | SourceType
  timeRange: TimeRange
  customFrom: string
  customTo: string
  keywords: string[]
  trustLevels: TrustLevel[]
  stages: SopStage[]
  sort: SortKey
}

export const defaultFilters: DashboardFilters = {
  status: 'all',
  platform: 'all',
  source: 'all',
  timeRange: 'all',
  customFrom: '',
  customTo: '',
  keywords: [],
  trustLevels: [],
  stages: [],
  sort: 'newest',
}

// 与 Mock 数据保持一致的"当前时间"参考点
export const MOCK_NOW = new Date('2026-07-07T10:10:00+08:00')

export function countActiveAdvancedFilters(f: DashboardFilters) {
  let n = 0
  if (f.source !== 'all') n++
  if (f.timeRange !== 'all') n++
  if (f.keywords.length > 0) n++
  if (f.trustLevels.length > 0) n++
  if (f.stages.length > 0) n++
  return n
}

export function applyFilters(list: Opportunity[], f: DashboardFilters): Opportunity[] {
  const now = MOCK_NOW.getTime()
  const dayMs = 24 * 60 * 60 * 1000

  const filtered = list.filter((o) => {
    if (f.status !== 'all' && o.status !== f.status) return false
    if (f.platform !== 'all' && o.platform !== f.platform) return false
    if (f.source !== 'all' && o.sourceType !== f.source) return false

    const t = new Date(o.createdAt).getTime()
    if (f.timeRange === 'today') {
      const startOfDay = new Date(MOCK_NOW)
      startOfDay.setHours(0, 0, 0, 0)
      if (t < startOfDay.getTime()) return false
    } else if (f.timeRange === '3d') {
      if (t < now - 3 * dayMs) return false
    } else if (f.timeRange === '7d') {
      if (t < now - 7 * dayMs) return false
    } else if (f.timeRange === 'custom') {
      if (f.customFrom && t < new Date(`${f.customFrom}T00:00:00+08:00`).getTime()) return false
      if (f.customTo && t > new Date(`${f.customTo}T23:59:59+08:00`).getTime()) return false
    }

    if (f.keywords.length > 0 && !f.keywords.some((k) => o.matchedKeywords.includes(k))) return false
    if (f.trustLevels.length > 0 && !f.trustLevels.includes(trustLevel(o.trustScore))) return false
    if (f.stages.length > 0 && !f.stages.includes(o.sopStage)) return false
    return true
  })

  return filtered.sort((a, b) => {
    switch (f.sort) {
      case 'confidence':
        return b.confidenceScore - a.confidenceScore
      case 'trust':
        return b.trustScore - a.trustScore
      case 'oldest':
        return new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime()
      default:
        return new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
    }
  })
}

import type { Opportunity } from './types'

export interface MiraStats {
  active: Opportunity[]
  attention: Opportunity[]
  business: Opportunity[]
  digestCount: number
  focusCount: number
  jobs: Opportunity[]
  judgment: Opportunity[]
  latestAt: string | null
  pending: Opportunity[]
  quietCount: number
  totalProcessed: number
}

function isUnverifiedOrRisky(item: Opportunity) {
  const status = item.linkVerification.status
  return status === 'unverified' || status === 'verifying' || status === 'suspicious' || status === 'malicious'
}

export function getMiraStats(opportunities: Opportunity[]): MiraStats {
  const active = opportunities.filter((item) => !item.archivedAt)
  const pending = active.filter((item) => item.status === 'pending')
  const attention = pending.filter((item) => item.attentionRequired || item.priority === 'urgent' || item.priority === 'high')
  const jobs = active.filter((item) => item.opportunityType === 'job')
  const business = active.filter((item) => item.opportunityType !== 'job')
  const judgment = pending.filter((item) => item.rawMessageLinks.length > 0 && isUnverifiedOrRisky(item))
  const quietCount = opportunities.filter((item) => item.archivedAt || item.status === 'ignored').length
  const digestCount = active.filter((item) => item.status !== 'pending').length
  let latestAt: string | null = null
  for (const item of active) {
    if (!latestAt || Date.parse(item.createdAt) > Date.parse(latestAt)) latestAt = item.createdAt
  }

  return {
    active,
    attention,
    business,
    digestCount,
    focusCount: Math.max(attention.length, pending.length),
    jobs,
    judgment,
    latestAt,
    pending,
    quietCount,
    totalProcessed: opportunities.length,
  }
}

export function formatMiraClock(iso: string | null) {
  if (!iso) return '--:--'
  const date = new Date(iso)
  if (Number.isNaN(date.getTime())) return '--:--'
  return `${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
}

export function greetingForNow(date = new Date()) {
  const hour = date.getHours()
  if (hour < 5) return '夜深了'
  if (hour < 11) return '早上好'
  if (hour < 14) return '中午好'
  if (hour < 18) return '下午好'
  return '晚上好'
}

export function buildMiraSummary(stats: MiraStats) {
  const firstAttention = stats.attention[0]
  if (firstAttention) {
    const action = firstAttention.rawMessageLinks.length > 0 ? '先核验链接，再决定是否跟进' : '可以优先处理'
    return `${firstAttention.contactName} 的信息需要你先看，Mira 判断它和当前目标更相关，${action}。`
  }
  if (stats.pending.length > 0) {
    return `Mira 已把 ${stats.pending.length} 条消息放入待处理队列，其余消息会继续归入摘要或安静区。`
  }
  if (stats.totalProcessed > 0) {
    return `当前没有必须立刻处理的消息，Mira 已整理 ${stats.totalProcessed} 条记录，可在消息页继续查看。`
  }
  return 'Mira 还在认识你的信息胃口。接入几条真实消息后，正式偏好会在你确认后生效。'
}

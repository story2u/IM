import type { FriendRequestStatus, LinkVerificationStatus, SopStage } from './types'

export const sopStageOrder: SopStage[] = [
  'detected',
  'analyzing',
  'verified',
  'contact_extracted',
  'friend_requested',
  'ready_to_chat',
  'chatting',
  'closed',
]

export const sopStageConfig: Record<SopStage, { label: string; dotClass: string }> = {
  detected: { label: '已发现', dotClass: 'bg-muted-foreground' },
  analyzing: { label: 'AI 分析中', dotClass: 'bg-primary' },
  verified: { label: '已核验', dotClass: 'bg-primary' },
  contact_extracted: { label: '已提取联系方式', dotClass: 'bg-primary' },
  friend_requested: { label: '待加好友', dotClass: 'bg-warning' },
  ready_to_chat: { label: '可对话', dotClass: 'bg-success' },
  chatting: { label: '沟通中', dotClass: 'bg-success' },
  closed: { label: '已结束', dotClass: 'bg-muted-foreground' },
}

export type TrustLevel = 'trusted' | 'unverified' | 'suspicious' | 'risky'

export function trustLevel(score: number): TrustLevel {
  if (score >= 80) return 'trusted'
  if (score >= 60) return 'unverified'
  if (score >= 40) return 'suspicious'
  return 'risky'
}

export const trustLevelConfig: Record<TrustLevel, { label: string; className: string }> = {
  trusted: { label: '安全可信', className: 'bg-success/10 text-success border-success/30' },
  unverified: { label: '待核验', className: 'bg-muted text-muted-foreground border-border' },
  suspicious: { label: '可疑', className: 'bg-warning/10 text-warning border-warning/30' },
  risky: { label: '高风险', className: 'bg-destructive/10 text-destructive border-destructive/30' },
}

export const linkStatusConfig: Record<LinkVerificationStatus, { label: string; className: string }> = {
  unverified: { label: '未核验', className: 'bg-muted text-muted-foreground border-border' },
  verifying: { label: '分析中', className: 'bg-primary/10 text-primary border-primary/30' },
  safe: { label: '安全', className: 'bg-success/10 text-success border-success/30' },
  suspicious: { label: '可疑', className: 'bg-warning/10 text-warning border-warning/30' },
  malicious: { label: '高风险', className: 'bg-destructive/10 text-destructive border-destructive/30' },
}

export const friendRequestConfig: Record<FriendRequestStatus, { label: string; className: string }> = {
  not_sent: { label: '未发送', className: 'bg-muted text-muted-foreground border-border' },
  pending: { label: '待通过', className: 'bg-warning/10 text-warning border-warning/30' },
  accepted: { label: '已通过', className: 'bg-success/10 text-success border-success/30' },
  rejected: { label: '被拒绝', className: 'bg-destructive/10 text-destructive border-destructive/30' },
  'n/a': { label: '无需此步骤', className: 'bg-muted text-muted-foreground border-border' },
}

export function formatDateTime(iso: string | null) {
  if (!iso) return ''
  const d = new Date(iso)
  return `${d.getMonth() + 1}/${d.getDate()} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

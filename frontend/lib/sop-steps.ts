import type { Opportunity } from './types'

export type StepState = 'done' | 'active' | 'in_progress' | 'blocked' | 'failed' | 'skipped' | 'locked'

export interface SopStep {
  key: 'discovery' | 'link_verification' | 'contacts' | 'friend_request' | 'chat'
  index: number
  title: string
  state: StepState
  stateLabel: string
  completedAt: string | null
  blockReason?: string
}

export function hasAnyContact(o: Opportunity) {
  const c = o.extractedContacts
  return Boolean(c.phone || c.email || c.telegramHandle || c.wecomId)
}

export function deriveSteps(o: Opportunity): SopStep[] {
  const hasLinks = o.rawMessageLinks.length > 0
  const lv = o.linkVerification
  const linkFailed = hasLinks && (lv.status === 'suspicious' || lv.status === 'malicious')
  const linkOK = !hasLinks || lv.status === 'safe'
  const contactsOK = hasAnyContact(o)
  const friendNeeded = o.sourceType === 'group' && o.friendRequestStatus !== 'n/a'
  const friendOK = !friendNeeded || o.friendRequestStatus === 'accepted'
  const chatUnlocked = linkOK && contactsOK && friendOK

  const steps: SopStep[] = []

  // Step 1：商机发现（永远已完成）
  steps.push({
    key: 'discovery',
    index: 1,
    title: '商机发现',
    state: 'done',
    stateLabel: '已完成',
    completedAt: o.createdAt,
  })

  // Step 2：链接安全与真实性分析
  if (!hasLinks) {
    steps.push({
      key: 'link_verification',
      index: 2,
      title: '链接安全分析',
      state: 'skipped',
      stateLabel: '无链接，跳过核验',
      completedAt: o.createdAt,
    })
  } else if (lv.status === 'verifying') {
    steps.push({
      key: 'link_verification',
      index: 2,
      title: '链接安全分析',
      state: 'in_progress',
      stateLabel: '分析中',
      completedAt: null,
    })
  } else if (lv.status === 'safe') {
    steps.push({
      key: 'link_verification',
      index: 2,
      title: '链接安全分析',
      state: 'done',
      stateLabel: '已完成 · 安全',
      completedAt: lv.verifiedAt,
    })
  } else if (linkFailed) {
    steps.push({
      key: 'link_verification',
      index: 2,
      title: '链接安全分析',
      state: 'failed',
      stateLabel: lv.status === 'malicious' ? '高风险 · 流程中断' : '可疑 · 流程中断',
      completedAt: lv.verifiedAt,
    })
  } else {
    steps.push({
      key: 'link_verification',
      index: 2,
      title: '链接安全分析',
      state: 'active',
      stateLabel: '待分析',
      completedAt: null,
    })
  }

  // Step 3：联系方式提取
  if (!linkOK) {
    steps.push({
      key: 'contacts',
      index: 3,
      title: '联系方式提取',
      state: 'blocked',
      stateLabel: '等待前置步骤',
      completedAt: null,
      blockReason: linkFailed ? '链接被判定为有风险，流程已中断' : '需先完成链接安全分析',
    })
  } else if (contactsOK) {
    steps.push({
      key: 'contacts',
      index: 3,
      title: '联系方式提取',
      state: 'done',
      stateLabel: '已完成',
      completedAt: lv.verifiedAt ?? o.createdAt,
    })
  } else {
    steps.push({
      key: 'contacts',
      index: 3,
      title: '联系方式提取',
      state: 'active',
      stateLabel: '待补充',
      completedAt: null,
    })
  }

  // Step 4：建立联系
  if (!friendNeeded) {
    steps.push({
      key: 'friend_request',
      index: 4,
      title: '建立联系',
      state: 'skipped',
      stateLabel: '私聊来源，无需此步骤',
      completedAt: o.createdAt,
    })
  } else if (!linkOK || !contactsOK) {
    steps.push({
      key: 'friend_request',
      index: 4,
      title: '建立联系',
      state: 'blocked',
      stateLabel: '等待前置步骤',
      completedAt: null,
      blockReason: !linkOK ? '需先完成链接安全核验' : '需先获取至少一项联系方式',
    })
  } else if (o.friendRequestStatus === 'accepted') {
    steps.push({
      key: 'friend_request',
      index: 4,
      title: '建立联系',
      state: 'done',
      stateLabel: '已通过',
      completedAt: lv.verifiedAt ?? o.createdAt,
    })
  } else if (o.friendRequestStatus === 'pending') {
    steps.push({
      key: 'friend_request',
      index: 4,
      title: '建立联系',
      state: 'in_progress',
      stateLabel: '申请待通过',
      completedAt: null,
    })
  } else if (o.friendRequestStatus === 'rejected') {
    steps.push({
      key: 'friend_request',
      index: 4,
      title: '建立联系',
      state: 'failed',
      stateLabel: '申请被拒绝',
      completedAt: null,
    })
  } else {
    steps.push({
      key: 'friend_request',
      index: 4,
      title: '建立联系',
      state: 'active',
      stateLabel: '待发送申请',
      completedAt: null,
    })
  }

  // Step 5：聊天与回复
  if (chatUnlocked) {
    steps.push({
      key: 'chat',
      index: 5,
      title: '聊天与回复',
      state: o.sopStage === 'chatting' ? 'in_progress' : 'active',
      stateLabel: o.sopStage === 'chatting' ? '沟通中' : '已解锁',
      completedAt: null,
    })
  } else {
    const reasons: string[] = []
    if (!linkOK) reasons.push('链接安全核验')
    if (!contactsOK) reasons.push('联系方式提取')
    if (!friendOK) reasons.push('好友添加')
    steps.push({
      key: 'chat',
      index: 5,
      title: '聊天与回复',
      state: 'locked',
      stateLabel: '未解锁',
      completedAt: null,
      blockReason: `完成${reasons.join('、')}后，即可开始对话`,
    })
  }

  return steps
}

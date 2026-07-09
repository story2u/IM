import type {
  ExtractedContacts,
  LinkVerification,
  Opportunity,
  ReplyTemplate,
} from './types'

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL ?? ''

interface ApiOpportunity {
  id: string
  platform: Opportunity['platform']
  contactName: string
  contactAvatar?: string
  summary: string
  matchedKeywords: string[]
  confidenceScore: number
  status: Opportunity['status']
  priority: Opportunity['priority']
  lastMessagePreview: string
  createdAt: string
  sourceType?: Opportunity['sourceType']
  groupName?: string | null
  groupMemberRole?: Opportunity['groupMemberRole']
  rawMessageLinks?: string[]
  linkVerification?: LinkVerification
  extractedContacts?: ExtractedContacts
  friendRequestStatus?: Opportunity['friendRequestStatus']
  sopStage?: Opportunity['sopStage']
  trustScore?: number
}

const defaultLinkVerification: LinkVerification = {
  status: 'unverified',
  verifiedAt: null,
  riskReasons: [],
  resolvedInfo: null,
}

const defaultContacts: ExtractedContacts = {
  phone: null,
  email: null,
  telegramHandle: null,
  wecomId: null,
  extractionSource: null,
}

function apiUrl(path: string) {
  return `${API_BASE_URL}${path}`
}

async function fetchJson<T>(path: string): Promise<T> {
  const response = await fetch(apiUrl(path), {
    headers: { Accept: 'application/json' },
    cache: 'no-store',
  })
  if (!response.ok) {
    throw new Error(`API ${path} failed with ${response.status}`)
  }
  return response.json() as Promise<T>
}

export function toOpportunity(item: ApiOpportunity): Opportunity {
  return {
    id: item.id,
    platform: item.platform,
    contactName: item.contactName,
    contactAvatar: item.contactAvatar || '/placeholder-user.jpg',
    summary: item.summary,
    matchedKeywords: item.matchedKeywords ?? [],
    confidenceScore: item.confidenceScore,
    status: item.status,
    priority: item.priority,
    lastMessagePreview: item.lastMessagePreview,
    createdAt: item.createdAt,
    sourceType: item.sourceType ?? 'private',
    groupName: item.groupName ?? null,
    groupMemberRole: item.groupMemberRole ?? 'member',
    rawMessageLinks: item.rawMessageLinks ?? [],
    linkVerification: item.linkVerification ?? defaultLinkVerification,
    extractedContacts: item.extractedContacts ?? defaultContacts,
    friendRequestStatus: item.friendRequestStatus ?? (item.sourceType === 'group' ? 'not_sent' : 'n/a'),
    sopStage: item.sopStage ?? 'detected',
    trustScore: item.trustScore ?? 70,
  }
}

export async function fetchOpportunities(): Promise<Opportunity[]> {
  const items = await fetchJson<ApiOpportunity[]>('/api/v1/opportunities?limit=200')
  return items.map(toOpportunity)
}

export async function fetchReplyTemplates(): Promise<ReplyTemplate[]> {
  return fetchJson<ReplyTemplate[]>('/api/v1/templates')
}

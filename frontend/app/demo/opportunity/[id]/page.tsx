'use client'

import { ArrowLeft, CheckCircle2, ExternalLink, Mail, ShieldCheck, Sparkles, UserCheck } from 'lucide-react'
import Link from 'next/link'
import { useParams } from 'next/navigation'
import { useState } from 'react'
import { DemoShell } from '@/components/demo/demo-shell'
import { PlatformIcon } from '@/components/platform-icon'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { demoMessages, demoOpportunities } from '@/lib/demo/opportunity-demo-data'

export default function DemoOpportunityDetail() {
  const { id } = useParams<{ id: string }>()
  const opportunity = demoOpportunities.find(o => o.id === id) ?? demoOpportunities[0]
  const [draft, setDraft] = useState('')
  const generate = () => setDraft('您好，我们可以安排 50 套设备的方案演示。为了准备更有针对性的内容，想先确认使用场景和关键规格。请问下周二或周三哪天方便？')
  return <DemoShell><div className="mx-auto max-w-5xl px-4 py-6 md:px-8" data-testid="demo-opportunity-detail"><div className="flex items-center gap-2"><Button variant="ghost" size="icon" nativeButton={false} render={<Link href="/demo" />} aria-label="返回看板"><ArrowLeft className="size-4" /></Button><h1 className="text-xl font-semibold">商机详情</h1><Badge>待人工审核</Badge></div>
    <section className="mt-4 rounded-md border bg-card p-5"><div className="flex flex-wrap items-start gap-4"><span className="grid size-11 place-items-center rounded-md bg-muted"><PlatformIcon platform={opportunity.platform} /></span><div className="min-w-0 flex-1"><h2 className="font-semibold">{opportunity.contactName}</h2><p className="mt-1 text-xs text-muted-foreground">{opportunity.groupName ?? '私聊'} · {opportunity.platform === 'telegram' ? 'Telegram' : '企业微信'}</p></div><div className="flex gap-6 text-center"><Metric label="相关度" value={`${Math.round(opportunity.confidenceScore * 100)}%`} /><Metric label="可信度" value={`${opportunity.trustScore}`} /></div></div><p className="mt-5 text-base leading-7">{opportunity.summary}</p><div className="mt-4 flex gap-2">{opportunity.matchedKeywords.map(k => <Badge key={k} variant="secondary">{k}</Badge>)}</div></section>
    <div className="mt-4 grid gap-4 lg:grid-cols-[1fr_1.15fr]"><section className="rounded-md border bg-card p-5"><h2 className="text-sm font-semibold">原始上下文</h2><div className="mt-4 space-y-3">{(demoMessages[opportunity.id] ?? []).map(m => <div key={m.id} className="rounded-md bg-muted p-3 text-sm leading-6"><strong className="block text-xs">{m.senderName}</strong>{m.content}</div>)}</div><h2 className="mt-7 text-sm font-semibold">Pi Agent 结论</h2><dl className="mt-3 space-y-3 text-sm"><Row icon={ShieldCheck} label="风险核验" value="未包含外部链接，来源连接状态正常" /><Row icon={Mail} label="联系方式" value={opportunity.extractedContacts.email ?? 'procurement@example.com'} /><Row icon={UserCheck} label="建议动作" value="确认演示时间和设备规格，执行前需要人工批准" /><Row icon={CheckCircle2} label="SOP 阶段" value="联系方式已提取，等待人工跟进" /></dl></section>
    <section className="rounded-md border bg-card p-5"><div className="flex items-center justify-between"><h2 className="text-sm font-semibold">可编辑回复草稿</h2><Badge variant="outline">不会自动发送</Badge></div><Textarea className="mt-4 min-h-44" value={draft} onChange={e => setDraft(e.target.value)} placeholder="点击下方按钮生成演示草稿" data-testid="demo-draft" /><div className="mt-3 flex items-center justify-between"><Button variant="outline" onClick={generate} data-testid="generate-demo-draft"><Sparkles className="size-4" />AI 生成草稿</Button><Button disabled data-testid="disabled-send">发送（演示禁用）<ExternalLink className="size-3.5" /></Button></div><p className="mt-4 text-xs leading-5 text-muted-foreground">演示模式不会调用模型、不会写数据库，也不会向 Telegram 或企业微信发送消息。</p></section></div>
  </div></DemoShell>
}
function Metric({ label, value }: { label: string; value: string }) { return <div><strong className="text-xl text-primary">{value}</strong><span className="block text-[10px] text-muted-foreground">{label}</span></div> }
function Row({ icon: Icon, label, value }: { icon: typeof ShieldCheck; label: string; value: string }) { return <div className="flex gap-3 border-t pt-3"><Icon className="mt-0.5 size-4 shrink-0 text-primary" /><div><dt className="text-xs text-muted-foreground">{label}</dt><dd className="mt-1">{value}</dd></div></div> }

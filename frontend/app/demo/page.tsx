'use client'

import { AlertTriangle, SlidersHorizontal } from 'lucide-react'
import Link from 'next/link'
import { useMemo, useState } from 'react'
import { DemoShell } from '@/components/demo/demo-shell'
import { PlatformIcon } from '@/components/platform-icon'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { demoOpportunities } from '@/lib/demo/opportunity-demo-data'

export default function DemoDashboard() {
  const [platform, setPlatform] = useState<'all' | 'telegram' | 'wecom'>('all')
  const [advanced, setAdvanced] = useState(false)
  const filtered = useMemo(() => demoOpportunities.filter((o) => platform === 'all' || o.platform === platform), [platform])
  const attention = demoOpportunities.filter((o) => o.attentionRequired)
  return <DemoShell><div className="mx-auto max-w-6xl px-4 py-6 md:px-8" data-testid="demo-dashboard">
    <header className="flex flex-wrap items-end justify-between gap-3"><div><h1 className="text-2xl font-semibold">商机看板</h1><p className="mt-1 text-sm text-muted-foreground">自动识别 Telegram 与企业微信中的潜在商机</p></div><Badge variant="secondary">{demoOpportunities.filter(o => o.status === 'pending').length} 条待处理</Badge></header>
    <section className="mt-5 rounded-md border border-warning/40 bg-warning/10 p-4" data-testid="attention-alert"><div className="flex gap-3"><AlertTriangle className="size-5 text-warning" /><div><p className="text-sm font-semibold">Pi Agent 发现 {attention.length} 条重大商机</p><p className="mt-1 text-xs text-muted-foreground">请优先核对链接结论和后续行动建议，外部动作仍需人工批准。</p><Button className="mt-3" size="sm" nativeButton={false} render={<Link href="/demo/opportunity/demo-procurement-50" />} data-testid="open-attention">查看采购需求</Button></div></div></section>
    <div className="mt-5 flex flex-wrap gap-2" data-testid="dashboard-filters"><FilterButton active={platform === 'all'} onClick={() => setPlatform('all')}>全部平台</FilterButton><FilterButton active={platform === 'telegram'} onClick={() => setPlatform('telegram')} testId="filter-telegram">Telegram</FilterButton><FilterButton active={platform === 'wecom'} onClick={() => setPlatform('wecom')}>企业微信</FilterButton><Button variant="outline" size="sm" onClick={() => setAdvanced(v => !v)} data-testid="advanced-filter"><SlidersHorizontal className="size-3.5" />高级筛选</Button></div>
    {advanced && <div className="mt-3 flex flex-wrap gap-2 rounded-md border bg-card p-3 text-xs" data-testid="advanced-filter-panel"><Badge>群消息</Badge><Badge variant="secondary">安全可信</Badge><Badge variant="secondary">高相关度排序</Badge><span className="ml-auto text-muted-foreground">筛选仅作用于演示数据</span></div>}
    <p className="my-4 text-xs text-muted-foreground">当前共 {filtered.length} 条商机</p>
    <div className="grid gap-3 lg:grid-cols-2">{filtered.map((o) => <Link key={o.id} href={`/demo/opportunity/${o.id}`} className="rounded-md border bg-card p-4 shadow-sm transition-shadow hover:shadow-md" data-testid={`demo-card-${o.id}`}><div className="flex items-start gap-3"><span className="grid size-9 place-items-center rounded-md bg-muted"><PlatformIcon platform={o.platform} className="size-4" /></span><div className="min-w-0 flex-1"><div className="flex flex-wrap items-center gap-2"><strong className="text-sm">{o.contactName}</strong><Badge variant="outline" className="text-[10px]">{o.priority === 'urgent' ? '紧急' : o.priority === 'high' ? '高优先级' : '普通'}</Badge><span className="text-[10px] text-success">可信度 {o.trustScore}</span></div><p className="mt-2 text-sm leading-6">{o.summary}</p><div className="mt-3 flex flex-wrap gap-1.5">{o.matchedKeywords.map(k => <Badge key={k} variant="secondary" className="text-[10px]">{k}</Badge>)}</div><p className="mt-3 text-[10px] text-muted-foreground">{o.groupName ?? '私聊'} · SOP {o.sopStage}</p></div><strong className="text-lg text-primary">{Math.round(o.confidenceScore * 100)}%</strong></div></Link>)}</div>
  </div></DemoShell>
}

function FilterButton({ active, onClick, children, testId }: { active: boolean; onClick: () => void; children: React.ReactNode; testId?: string }) { return <Button variant={active ? 'default' : 'outline'} size="sm" onClick={onClick} data-testid={testId}>{children}</Button> }

import { Bell, CheckCircle2, Clock3, CreditCard, MessageSquare, Send, ShieldCheck } from 'lucide-react'
import Link from 'next/link'
import { DemoShell } from '@/components/demo/demo-shell'
import { Badge } from '@/components/ui/badge'

const panels = [
  ['套餐与用量', '/demo/settings/subscription', CreditCard, 'Pro · 本月 18 / 500 次分析'],
  ['Telegram 原生连接', '/demo/settings/telegram', Send, '2 个演示来源，全部只读'],
  ['工作时间', '/demo/settings/working-hours', Clock3, '周一至周五 09:00–18:00'],
  ['通知偏好', '/demo/settings#notifications', Bell, '重大商机与人工待办提醒'],
] as const

export function DemoSettings({ panel = 'overview' }: { panel?: 'overview' | 'subscription' | 'telegram' | 'working-hours' }) {
  return <DemoShell><div className="mx-auto max-w-4xl px-4 py-6 md:px-8" data-testid={`demo-settings-${panel}`}><h1 className="text-2xl font-semibold">设置中心</h1><p className="mt-1 text-sm text-muted-foreground">演示页面只展示当前已实现的配置形态，不调用外部平台。</p>
    {panel === 'overview' ? <div className="mt-6 grid gap-3 sm:grid-cols-2">{panels.map(([title, href, Icon, copy]) => <Link key={href} href={href} className="flex items-center gap-4 rounded-md border bg-card p-5 shadow-sm"><span className="grid size-10 place-items-center rounded-md bg-primary/10 text-primary"><Icon className="size-5" /></span><div className="min-w-0"><strong className="text-sm">{title}</strong><p className="mt-1 text-xs text-muted-foreground">{copy}</p></div></Link>)}</div> : null}
    {panel === 'subscription' && <section className="mt-6 rounded-md border bg-card p-6"><div className="flex items-start justify-between"><div><Badge>Pro</Badge><h2 className="mt-3 text-xl font-semibold">套餐与本月用量</h2></div><span className="text-sm text-success">权益有效</span></div><div className="mt-7 grid gap-3 sm:grid-cols-3"><Stat label="Telegram / 企微群" value="2 / 10" /><Stat label="Pi Agent 月分析" value="18 / 500" /><Stat label="使用周期" value="7 月 1–31 日" /></div><p className="mt-6 text-xs text-muted-foreground">演示模式不初始化 RevenueCat，也不提供购买按钮。真实页面价格来自 RevenueCat Offering。</p></section>}
    {panel === 'telegram' && <section className="mt-6 space-y-3"><Connection icon={Send} title="Telegram 普通账号" copy="已连接 · 只读监听 · Session 服务端加密" /><Connection icon={MessageSquare} title="企业采购需求站（演示）" copy="群组 · 已监听 · 未触发额度暂停" /><Connection icon={ShieldCheck} title="安全边界" copy="凭据不会进入浏览器、截图数据或 AI 请求" /><p className="pt-3 text-xs text-muted-foreground">二维码授权、Bot 连接和真实群选择在演示模式下全部禁用。</p></section>}
    {panel === 'working-hours' && <section className="mt-6 rounded-md border bg-card p-6"><h2 className="font-semibold">人工工作时间</h2><div className="mt-6 grid grid-cols-2 gap-3 sm:grid-cols-5">{['周一','周二','周三','周四','周五'].map(d => <div key={d} className="rounded-md border bg-primary/5 p-3 text-center text-sm"><CheckCircle2 className="mx-auto mb-2 size-4 text-success" />{d}</div>)}</div><div className="mt-5 flex flex-wrap gap-8 text-sm"><span><small className="block text-muted-foreground">开始</small>09:00</span><span><small className="block text-muted-foreground">结束</small>18:00</span><span><small className="block text-muted-foreground">时区</small>Asia/Shanghai</span></div><p className="mt-6 text-xs text-muted-foreground">语义模型新发现的商机始终进入人工审核，不会因非工作时间自动发送。</p></section>}
  </div></DemoShell>
}
function Stat({ label, value }: { label: string; value: string }) { return <div className="rounded-md bg-muted p-4"><span className="text-xs text-muted-foreground">{label}</span><strong className="mt-2 block text-lg">{value}</strong></div> }
function Connection({ icon: Icon, title, copy }: { icon: typeof Send; title: string; copy: string }) { return <div className="flex items-center gap-4 rounded-md border bg-card p-5"><span className="grid size-10 place-items-center rounded-md bg-primary/10 text-primary"><Icon className="size-5" /></span><div><strong className="text-sm">{title}</strong><p className="mt-1 text-xs text-muted-foreground">{copy}</p></div><Badge variant="secondary" className="ml-auto">演示</Badge></div> }

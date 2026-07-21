'use client'

import { Archive, Bot, ChevronRight, Clock3, MessageSquareText, Settings, ShieldCheck, Sparkles } from 'lucide-react'
import Link from 'next/link'
import { Button } from '@/components/ui/button'
import { useAppStore } from '@/lib/app-store'
import { buildMiraSummary, getMiraStats } from '@/lib/mira'

function ActionLink({
  copy,
  href,
  icon: Icon,
  title,
}: {
  copy: string
  href: string
  icon: typeof Sparkles
  title: string
}) {
  return (
    <Link href={href} className="flex items-center gap-4 rounded-md border border-white/10 bg-[#0c1d30] p-4 transition hover:border-cyan-300/45">
      <span className="grid size-11 shrink-0 place-items-center rounded-md bg-cyan-300/10 text-cyan-300">
        <Icon className="size-5" />
      </span>
      <span className="min-w-0 flex-1">
        <span className="block font-black text-white">{title}</span>
        <span className="mt-1 block text-sm leading-6 text-slate-400">{copy}</span>
      </span>
      <ChevronRight className="size-5 shrink-0 text-slate-500" />
    </Link>
  )
}

export default function MiraPage() {
  const { opportunities } = useAppStore()
  const stats = getMiraStats(opportunities)
  const summary = buildMiraSummary(stats)

  return (
    <div className="min-h-full bg-[#06111f] text-white" data-testid="mira-agent-page">
      <div className="mx-auto grid w-full max-w-6xl gap-5 px-5 py-6 md:px-8 lg:grid-cols-[0.95fr_1.05fr]">
        <header className="lg:col-span-2">
          <p className="text-xs font-black tracking-[0.35em] text-teal-300">OPENMIRA</p>
          <h1 className="mt-3 text-4xl font-black tracking-normal">Mira</h1>
          <p className="mt-3 max-w-3xl text-base leading-7 text-slate-300">
            Mira 负责整理信息、解释保留原因、学习偏好；Pi Agent 是底层运行时，外部动作仍然需要你批准。
          </p>
        </header>

        <section className="rounded-md border border-cyan-300/20 bg-[#0b2238] p-5">
          <div className="flex items-start justify-between gap-4">
            <div>
              <p className="text-xs font-black text-cyan-300">当前判断</p>
              <h2 className="mt-3 text-3xl font-black leading-tight">Mira 正在学习你的信息胃口</h2>
            </div>
            <span className="grid size-12 place-items-center rounded-md bg-cyan-300/10 text-cyan-300">
              <Sparkles className="size-6" />
            </span>
          </div>
          <p className="mt-5 text-lg font-semibold leading-8 text-white">{summary}</p>
          <div className="mt-6 grid grid-cols-2 gap-3">
            <div className="rounded-md bg-[#081a25] p-4">
              <p className="text-3xl font-black text-amber-300 tabular-nums">{stats.pending.length}</p>
              <p className="mt-2 text-sm text-slate-300">待你处理</p>
            </div>
            <div className="rounded-md bg-[#081a25] p-4">
              <p className="text-3xl font-black text-slate-300 tabular-nums">{stats.quietCount}</p>
              <p className="mt-2 text-sm text-slate-300">已安静收起</p>
            </div>
          </div>
          <Button size="lg" nativeButton={false} render={<Link href="/messages?category=pending" />} className="mt-6 w-full bg-cyan-300 text-slate-950 hover:bg-cyan-200">
            处理 Mira 放到前面的消息
          </Button>
        </section>

        <section className="grid gap-3">
          <ActionLink
            href="/messages?category=judgment"
            icon={MessageSquareText}
            title="需要你判断"
            copy={`${stats.judgment.length} 条消息的链接或上下文需要确认。`}
          />
          <ActionLink
            href="/messages?category=quiet"
            icon={Archive}
            title="安静区"
            copy="查看被忽略、归档或暂时收起的信息。"
          />
          <ActionLink
            href="/settings"
            icon={Settings}
            title="账户与连接"
            copy="管理 Telegram、企业微信、订阅和安全设置。"
          />
          <ActionLink
            href="/templates"
            icon={Bot}
            title="回复素材"
            copy="保留原来的人工审核与回复模板能力。"
          />
        </section>

        <section className="rounded-md bg-[#0c1d30] p-5 lg:col-span-2">
          <div className="grid gap-5 lg:grid-cols-3">
            <div>
              <Clock3 className="size-6 text-teal-300" />
              <h2 className="mt-4 font-black">按时间段总结</h2>
              <p className="mt-2 text-sm leading-6 text-slate-400">早间、午间和晚间摘要会聚合真正需要回看的内容。</p>
            </div>
            <div>
              <ShieldCheck className="size-6 text-teal-300" />
              <h2 className="mt-4 font-black">不绕过你行动</h2>
              <p className="mt-2 text-sm leading-6 text-slate-400">发送、加好友、外部联系等动作始终需要人工批准。</p>
            </div>
            <div>
              <Sparkles className="size-6 text-teal-300" />
              <h2 className="mt-4 font-black">偏好可撤销</h2>
              <p className="mt-2 text-sm leading-6 text-slate-400">Mira 的判断来自结构化样本，后续可继续纠正和回滚。</p>
            </div>
          </div>
        </section>
      </div>
    </div>
  )
}

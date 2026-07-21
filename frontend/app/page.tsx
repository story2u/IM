'use client'

import { AlertTriangle, ArrowRight, Loader2, RefreshCw, Settings, ShieldCheck, Sparkles } from 'lucide-react'
import Link from 'next/link'
import { ProductHome } from '@/components/landing/product-home'
import { Button } from '@/components/ui/button'
import { useAppStore } from '@/lib/app-store'
import { useAuth } from '@/lib/auth'
import { buildMiraSummary, formatMiraClock, getMiraStats, greetingForNow } from '@/lib/mira'
import { cn } from '@/lib/utils'

function MiraOrb() {
  return (
    <div className="relative grid size-10 place-items-center" aria-hidden="true">
      <span className="absolute size-10 animate-pulse rounded-full border border-cyan-300/50 bg-indigo-500/15" />
      <span className="size-4 rounded-full bg-cyan-300 shadow-[0_0_22px_rgba(103,232,249,0.42)]" />
    </div>
  )
}

function MetricTile({ label, tone, value }: { label: string; tone: string; value: number }) {
  return (
    <div className="min-h-24 rounded-md bg-[#081a25] p-4">
      <p className={cn('text-3xl font-black tabular-nums', tone)}>{value}</p>
      <p className="mt-2 text-sm leading-5 text-slate-300">{label}</p>
    </div>
  )
}

function BriefingRow({ label, state, time }: { label: string; state: string; time: string }) {
  return (
    <div className="flex items-center justify-between gap-4 border-b border-white/8 py-3 last:border-0">
      <div>
        <p className="text-sm font-semibold text-white">{label}</p>
        <p className="mt-1 text-xs text-slate-400">{state}</p>
      </div>
      <span className="rounded-full border border-white/10 px-3 py-1 text-xs text-slate-300">{time}</span>
    </div>
  )
}

export default function TodayPage() {
  const { user, loading } = useAuth()
  const { backendError, backendLoading, opportunities, reloadBackendData } = useAppStore()

  if (!loading && !user) return <ProductHome />
  if (loading && !user) return <ProductHome />
  if (!user) return null

  const stats = getMiraStats(opportunities)
  const latestClock = formatMiraClock(stats.latestAt)
  const displayName = user.displayName || user.email.split('@')[0] || 'Bruce'
  const summary = buildMiraSummary(stats)

  return (
    <div className="min-h-full bg-[#06111f] text-white" data-testid="mira-today-page">
      <div className="mx-auto flex w-full max-w-6xl flex-col gap-5 px-5 py-6 md:px-8">
        <header className="flex items-start justify-between gap-4 pt-2">
          <div className="min-w-0">
            <p className="text-sm font-black text-teal-300">{greetingForNow()}，{displayName}</p>
            <h1 className="mt-5 max-w-3xl text-4xl font-black leading-tight tracking-normal md:text-6xl">
              今天 Mira 会这样照顾你的注意力
            </h1>
            <p className="mt-4 text-base leading-7 text-slate-300">
              {latestClock === '--:--'
                ? 'Mira 还在认识你的信息胃口'
                : `Mira 已整理至 ${latestClock}`}
              {' · '}
              本地处理 {stats.totalProcessed} 条 · 深度分析 {stats.attention.length} 条
            </p>
          </div>
          <div className="flex shrink-0 items-center gap-3">
            <MiraOrb />
            <Button variant="outline" size="icon" nativeButton={false} render={<Link href="/settings" />} aria-label="设置中心">
              <Settings className="size-4" />
            </Button>
          </div>
        </header>

        {backendError ? (
          <section
            role="alert"
            className="flex flex-col gap-3 rounded-md border border-rose-300/25 bg-rose-950/45 p-4 text-rose-100 sm:flex-row sm:items-center sm:justify-between"
          >
            <div className="flex gap-3">
              <AlertTriangle className="mt-0.5 size-5 shrink-0" />
              <div>
                <p className="font-semibold">Mira 暂时读不到服务端数据</p>
                <p className="mt-1 text-sm text-rose-100/75">{backendError}</p>
              </div>
            </div>
            <Button onClick={reloadBackendData} className="gap-2 bg-rose-100 text-rose-950 hover:bg-white">
              <RefreshCw className="size-4" />
              重试
            </Button>
          </section>
        ) : null}

        <section className="rounded-md border border-cyan-300/20 bg-[#0b2238] p-5 shadow-2xl shadow-cyan-950/20">
          <p className="text-xs font-black uppercase text-cyan-300">OpenMIRA</p>
          <div className="mt-3 flex flex-col gap-5 lg:flex-row lg:items-end lg:justify-between">
            <div>
              <h2 className="max-w-2xl text-3xl font-black leading-tight md:text-5xl">
                今天只需要你看 {stats.focusCount} 条
              </h2>
              <p className="mt-3 max-w-2xl text-sm leading-6 text-slate-300">
                其他消息会先进入摘要、安静区或继续等待更多上下文，外部动作仍需你批准。
              </p>
            </div>
            <div className="flex flex-wrap gap-3">
              <Button size="lg" nativeButton={false} render={<Link href="/messages?category=pending" />} className="gap-2 bg-cyan-300 text-slate-950 hover:bg-cyan-200">
                开始处理
                <ArrowRight className="size-4" />
              </Button>
              <Button size="lg" variant="outline" nativeButton={false} render={<Link href="/mira" />}>
                为什么是这些
              </Button>
            </div>
          </div>
          <div className="mt-5 grid gap-3 md:grid-cols-3">
            <MetricTile value={stats.attention.length} label="条需要现在关注" tone="text-amber-300" />
            <MetricTile value={stats.business.length + stats.jobs.length} label="条可以稍后处理" tone="text-sky-300" />
            <MetricTile value={stats.judgment.length} label="条需要你帮助判断" tone="text-indigo-300" />
          </div>
        </section>

        <section className="grid gap-5 lg:grid-cols-[1fr_0.9fr]">
          <div className="rounded-md bg-[#0c1d30] p-5">
            <p className="text-xs font-black text-teal-300">Mira 的一句话</p>
            <p className="mt-3 text-xl font-bold leading-8 text-white">{summary}</p>
            <div className="mt-5 grid gap-3 sm:grid-cols-2">
              {stats.attention.slice(0, 2).map((item) => (
                <Link
                  key={item.id}
                  href={`/opportunity/${item.id}`}
                  className="rounded-md border border-white/10 bg-[#071725] p-4 transition hover:border-cyan-300/45"
                >
                  <p className="text-sm font-semibold text-white">{item.contactName}</p>
                  <p className="mt-2 line-clamp-3 text-sm leading-6 text-slate-300">{item.summary}</p>
                  <p className="mt-3 text-xs font-bold text-cyan-300">查看来源 ›</p>
                </Link>
              ))}
            </div>
          </div>

          <div className="rounded-md bg-[#0c1d30] p-5">
            <div className="flex items-center justify-between gap-4">
              <div>
                <p className="text-xs font-black text-teal-300">当前信息胃口</p>
                <h2 className="mt-2 text-xl font-black">Mira 还在认识你的偏好</h2>
              </div>
              {backendLoading ? <Loader2 className="size-5 animate-spin text-cyan-300" /> : <Sparkles className="size-5 text-cyan-300" />}
            </div>
            <div className="mt-5 rounded-md bg-[#071725] p-5 text-center">
              <Sparkles className="mx-auto size-7 text-teal-300" />
              <p className="mx-auto mt-4 max-w-sm text-sm leading-6 text-slate-300">
                用几条真实消息开始教学，正式偏好只会在你确认后生效。
              </p>
            </div>
            <Button nativeButton={false} render={<Link href="/mira" />} className="mt-4 w-full bg-[#155e75] hover:bg-[#0e7490]">
              教 Mira 几条
            </Button>
          </div>
        </section>

        <section className="grid gap-5 lg:grid-cols-[0.9fr_1.1fr]">
          <div className="rounded-md bg-[#0c1d30] p-5">
            <h2 className="text-xl font-black">三段简报</h2>
            <div className="mt-4">
              <BriefingRow label="早间简报" time="08:30" state={stats.totalProcessed > 0 ? '已基于当前消息更新' : '等待真实消息'} />
              <BriefingRow label="午间简报" time="12:00" state="聚合上午新增的需要关注项" />
              <BriefingRow label="晚间摘要" time="18:30" state={`${stats.digestCount} 条可进入晚间回看`} />
            </div>
          </div>
          <div className="rounded-md bg-[#0c1d30] p-5">
            <h2 className="text-xl font-black">今天的信息流</h2>
            <p className="mt-2 text-sm text-slate-400">Mira 已判断 {stats.totalProcessed} 条服务端消息</p>
            <div className="mt-5 grid grid-cols-2 gap-3">
              <MetricTile value={stats.attention.length} label="立即关注" tone="text-amber-300" />
              <MetricTile value={stats.pending.length} label="稍后处理" tone="text-sky-300" />
              <MetricTile value={stats.digestCount} label="摘要出现" tone="text-indigo-300" />
              <MetricTile value={stats.quietCount} label="安静收起" tone="text-slate-300" />
            </div>
          </div>
        </section>

        <footer className="flex items-center gap-2 pb-3 text-xs text-slate-500">
          <ShieldCheck className="size-4" />
          访问令牌只保存在设备安全存储中；Mira 不会绕过你执行外部动作。
        </footer>
      </div>
    </div>
  )
}

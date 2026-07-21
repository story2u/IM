'use client'

import { AlertTriangle, BriefcaseBusiness, ChevronRight, Inbox, Loader2, RefreshCw, SearchCheck, Sparkles } from 'lucide-react'
import Link from 'next/link'
import { useSearchParams } from 'next/navigation'
import { Button } from '@/components/ui/button'
import { useAppStore } from '@/lib/app-store'
import { getMiraStats } from '@/lib/mira'
import type { Opportunity } from '@/lib/types'
import { cn } from '@/lib/utils'

type Category = 'all' | 'pending' | 'business' | 'jobs' | 'judgment' | 'quiet'

const categoryLabels: Record<Category, string> = {
  all: '全部',
  pending: '待处理',
  business: '商机',
  jobs: '工作机会',
  judgment: '需要判断',
  quiet: '安静区',
}

const categoryDescriptions: Record<Category, string> = {
  all: '保留完整列表，便于搜索和检查。',
  pending: 'Mira 建议你优先看这些。',
  business: '采购、合作和业务线索。',
  jobs: '招聘、远程工作和岗位线索。',
  judgment: '链接或上下文需要你确认。',
  quiet: '已忽略或归档的信息。',
}

function categoryFromSearch(value: string | null): Category {
  return value === 'pending' || value === 'business' || value === 'jobs' || value === 'judgment' || value === 'quiet'
    ? value
    : 'all'
}

function SectionCard({
  active,
  category,
  count,
  tone,
}: {
  active: boolean
  category: Category
  count: number
  tone: string
}) {
  return (
    <Link
      href={`/messages?category=${category}`}
      className={cn(
        'min-h-32 rounded-md border bg-[#0c1d30] p-4 transition hover:border-cyan-300/45',
        active ? 'border-cyan-300/70' : 'border-white/10',
      )}
    >
      <div className="flex items-start justify-between gap-3">
        <p className="text-sm font-black text-white">{categoryLabels[category]}</p>
        <p className={cn('text-3xl font-black tabular-nums', tone)}>{count}</p>
      </div>
      <p className="mt-4 text-sm leading-6 text-slate-400">{categoryDescriptions[category]}</p>
      <p className="mt-3 text-xs font-bold text-cyan-300">打开 ›</p>
    </Link>
  )
}

function opportunityTone(item: Opportunity) {
  if (item.priority === 'urgent') return 'border-amber-300/55'
  if (item.attentionRequired) return 'border-cyan-300/45'
  return 'border-white/10'
}

function OpportunityRow({ item }: { item: Opportunity }) {
  const unverified = item.rawMessageLinks.length > 0 && (
    item.linkVerification.status === 'unverified' ||
    item.linkVerification.status === 'verifying' ||
    item.linkVerification.status === 'suspicious' ||
    item.linkVerification.status === 'malicious'
  )
  return (
    <Link
      href={`/opportunity/${item.id}`}
      className={cn('block rounded-md border bg-[#0c1d30] p-4 transition hover:border-cyan-300/55', opportunityTone(item))}
    >
      <div className="flex items-start justify-between gap-4">
        <div className="min-w-0">
          <div className="flex flex-wrap items-center gap-2">
            <p className="truncate text-base font-black text-white">{item.contactName}</p>
            <span className="rounded-full bg-white/8 px-2 py-1 text-xs text-slate-300">
              {item.opportunityType === 'job' ? '工作机会' : '商机'}
            </span>
            <span className="rounded-full bg-white/8 px-2 py-1 text-xs text-slate-300">
              {item.platform === 'telegram' ? 'Telegram' : '企业微信'}
            </span>
          </div>
          <p className="mt-3 line-clamp-3 text-sm leading-6 text-slate-300">{item.summary}</p>
        </div>
        <div className="shrink-0 text-right">
          <p className="text-2xl font-black tabular-nums text-cyan-300">{Math.round(item.confidenceScore * 100)}</p>
          <p className="text-xs text-slate-500">相关度</p>
        </div>
      </div>
      {unverified ? (
        <p className="mt-4 flex items-center gap-2 rounded-md border border-amber-300/30 bg-amber-300/10 px-3 py-2 text-sm text-amber-100">
          <AlertTriangle className="size-4 shrink-0" />
          含未核验链接，请先完成安全分析
        </p>
      ) : null}
      <div className="mt-4 flex flex-wrap items-center gap-2">
        {item.matchedKeywords.slice(0, 6).map((keyword) => (
          <span key={keyword} className="rounded-full bg-white/8 px-2 py-1 text-xs text-slate-300">{keyword}</span>
        ))}
        <span className="ml-auto flex items-center gap-1 text-xs font-bold text-cyan-300">
          查看详情
          <ChevronRight className="size-3.5" />
        </span>
      </div>
    </Link>
  )
}

export default function MessagesPage() {
  const searchParams = useSearchParams()
  const category = categoryFromSearch(searchParams.get('category'))
  const { backendError, backendLoading, opportunities, reloadBackendData } = useAppStore()
  const stats = getMiraStats(opportunities)
  const quietItems = opportunities.filter((item) => item.archivedAt || item.status === 'ignored')
  const categoryItems: Record<Category, Opportunity[]> = {
    all: stats.active,
    pending: stats.pending,
    business: stats.business,
    jobs: stats.jobs,
    judgment: stats.judgment,
    quiet: quietItems,
  }
  const items = categoryItems[category]

  return (
    <div className="min-h-full bg-[#06111f] text-white" data-testid="mira-messages-page">
      <div className="mx-auto flex w-full max-w-6xl flex-col gap-5 px-5 py-6 md:px-8">
        <header className="flex flex-col gap-4 pt-2 sm:flex-row sm:items-end sm:justify-between">
          <div>
            <p className="text-xs font-black tracking-[0.35em] text-teal-300">MIRA MESSAGE CENTER</p>
            <h1 className="mt-3 text-4xl font-black tracking-normal">消息</h1>
            <p className="mt-2 max-w-2xl text-base leading-7 text-slate-300">
              Mira 把 Telegram 与企业微信的消息先分层，再把真正需要你看的放到前面。
            </p>
          </div>
          <div className="rounded-md bg-[#0c1d30] px-5 py-4 text-center">
            <p className="text-3xl font-black text-amber-300 tabular-nums">{stats.pending.length}</p>
            <p className="text-xs text-slate-400">待处理</p>
          </div>
        </header>

        {backendError ? (
          <section role="alert" className="flex flex-col gap-3 rounded-md border border-rose-300/25 bg-rose-950/45 p-4 text-rose-100 sm:flex-row sm:items-center sm:justify-between">
            <p className="text-sm">{backendError}</p>
            <Button onClick={reloadBackendData} className="gap-2 bg-rose-100 text-rose-950 hover:bg-white">
              <RefreshCw className="size-4" />
              重试
            </Button>
          </section>
        ) : null}

        <section className="grid gap-3 md:grid-cols-2 xl:grid-cols-4">
          <SectionCard active={category === 'pending'} category="pending" count={stats.pending.length} tone="text-amber-300" />
          <SectionCard active={category === 'business'} category="business" count={stats.business.length} tone="text-cyan-300" />
          <SectionCard active={category === 'jobs'} category="jobs" count={stats.jobs.length} tone="text-indigo-300" />
          <SectionCard active={category === 'quiet'} category="quiet" count={stats.quietCount} tone="text-slate-300" />
        </section>

        <section className="rounded-md border border-teal-300/20 bg-teal-950/35 p-4">
          <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <div className="flex gap-3">
              <span className="grid size-12 shrink-0 place-items-center rounded-md bg-teal-300/10 text-teal-300">
                <Sparkles className="size-5" />
              </span>
              <div>
                <p className="font-black">教 Mira 认识你的偏好</p>
                <p className="mt-1 text-sm leading-6 text-slate-300">
                  每次保留或忽略都会形成可撤销样本，正式偏好只会在你确认后生效。
                </p>
              </div>
            </div>
            <Button nativeButton={false} render={<Link href="/mira" />} className="bg-[#155e75] hover:bg-[#0e7490]">
              教 Mira 几条
            </Button>
          </div>
        </section>

        <section className="flex flex-wrap items-center gap-2">
          {(Object.keys(categoryLabels) as Category[]).map((item) => (
            <Button
              key={item}
              nativeButton={false}
              render={<Link href={`/messages?category=${item}`} />}
              variant={category === item ? 'default' : 'outline'}
              className={category === item ? 'bg-cyan-300 text-slate-950 hover:bg-cyan-200' : ''}
            >
              {categoryLabels[item]}
            </Button>
          ))}
        </section>

        <section className="min-h-80">
          <div className="mb-4 flex items-center justify-between gap-4 text-sm text-slate-400">
            <span>{categoryLabels[category]} · {items.length} 条</span>
            {backendLoading ? <span className="flex items-center gap-2"><Loader2 className="size-4 animate-spin" />正在读取结果</span> : null}
          </div>
          {items.length > 0 ? (
            <div className="grid gap-3 lg:grid-cols-2">
              {items.map((item) => <OpportunityRow key={item.id} item={item} />)}
            </div>
          ) : (
            <div className="grid min-h-80 place-items-center rounded-md border border-dashed border-white/15 bg-[#081a25] p-8 text-center">
              <div>
                {category === 'jobs' ? <BriefcaseBusiness className="mx-auto size-10 text-slate-500" /> : category === 'judgment' ? <SearchCheck className="mx-auto size-10 text-slate-500" /> : <Inbox className="mx-auto size-10 text-slate-500" />}
                <p className="mt-4 text-xl font-black">暂无匹配消息</p>
                <p className="mt-2 max-w-md text-sm leading-6 text-slate-400">调整分类，或等待 Mira 继续整理新的聊天消息。</p>
              </div>
            </div>
          )}
        </section>
      </div>
    </div>
  )
}

'use client'

import { ArrowRight, BellOff, CheckCircle2, GitBranch, Layers3, LockKeyhole, MessageSquareText, MonitorSmartphone, ShieldCheck, Sparkles } from 'lucide-react'
import Link from 'next/link'
import { BrandLogo } from '@/components/brand-logo'
import { Button } from '@/components/ui/button'

const signalCards = [
  ['立即关注', '6', '采购、合作和需要你判断的消息'],
  ['晚间摘要', '18', '可以稍后回看的上下文'],
  ['安静收起', '392', '低价值、重复和已处理噪音'],
] as const

const capabilities = [
  [Layers3, '聚合聊天消息', 'Telegram、企业微信和后续渠道进入同一个信息流。'],
  [Sparkles, 'Mira 主动整理', '按立即关注、稍后处理、摘要和安静区组织消息。'],
  [MessageSquareText, '解释保留原因', '每条重要消息都能追溯到结构化证据和来源。'],
  [BellOff, '过滤噪音', '重复、低价值和已归档内容不会挤占第一屏。'],
  [ShieldCheck, '人工批准动作', '发送、加好友和外部联系不会绕过用户。'],
  [MonitorSmartphone, '三端一致体验', 'H5、iOS、Android 使用同一套 OpenMIRA 心智。'],
] as const

function Brand() {
  return (
    <span className="flex items-center gap-2.5">
      <BrandLogo size={36} priority />
      <span className="leading-tight">
        <strong className="block text-sm text-white">OpenMIRA</strong>
        <span className="block text-[10px] text-slate-400">智能消息管家</span>
      </span>
    </span>
  )
}

function SignalPreview() {
  return (
    <div className="mx-auto mt-12 w-full max-w-4xl rounded-md border border-cyan-300/20 bg-[#0b2238] p-4 shadow-2xl shadow-cyan-950/30">
      <div className="flex items-center justify-between gap-4">
        <div>
          <p className="text-xs font-black text-cyan-300">今天 Mira 会这样照顾你的注意力</p>
          <p className="mt-2 text-2xl font-black text-white">今天只需要你看 6 条</p>
        </div>
        <span className="grid size-12 place-items-center rounded-md bg-cyan-300/10 text-cyan-300">
          <Sparkles className="size-6" />
        </span>
      </div>
      <div className="mt-5 grid gap-3 md:grid-cols-3">
        {signalCards.map(([label, value, copy]) => (
          <div key={label} className="rounded-md bg-[#081a25] p-4">
            <p className="text-3xl font-black text-cyan-300 tabular-nums">{value}</p>
            <p className="mt-2 text-sm font-semibold text-white">{label}</p>
            <p className="mt-1 text-xs leading-5 text-slate-400">{copy}</p>
          </div>
        ))}
      </div>
    </div>
  )
}

export function ProductHome() {
  return (
    <div className="min-h-svh bg-[#06111f] text-white" data-testid="product-home">
      <header className="sticky top-0 z-40 border-b border-white/10 bg-[#06111f]/92 backdrop-blur">
        <div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-5 lg:px-8">
          <Link href="/"><Brand /></Link>
          <nav className="hidden items-center gap-6 text-sm text-slate-300 md:flex" aria-label="产品导航">
            <a href="#capabilities">能力</a>
            <a href="#safety">安全</a>
            <a href="https://github.com/story2u/IM" target="_blank" rel="noreferrer">GitHub</a>
          </nav>
          <div className="flex items-center gap-2">
            <Button nativeButton={false} render={<Link href="/login" />} variant="ghost" size="sm" className="text-white hover:bg-white/8 hover:text-white">登录</Button>
            <Button nativeButton={false} render={<Link href="/login" />} size="sm" className="bg-cyan-300 text-slate-950 hover:bg-cyan-200">开始体验</Button>
          </div>
        </div>
      </header>

      <main>
        <section className="border-b border-white/10 px-5 py-16 text-center md:py-24">
          <div className="mx-auto max-w-5xl">
            <p className="text-xs font-black tracking-[0.45em] text-teal-300">OPENMIRA</p>
            <h1 className="mt-5 text-5xl font-black leading-tight tracking-normal md:text-7xl">消息再多，也只看重要的</h1>
            <p className="mx-auto mt-6 max-w-2xl text-lg leading-8 text-slate-300">
              Mira 主动整理、过滤和汇总聊天消息，把真正需要你处理的信息放到前面，其余进入摘要或安静区。
            </p>
            <div className="mt-8 flex flex-wrap justify-center gap-3">
              <Button size="lg" nativeButton={false} render={<Link href="/login" />} className="gap-2 bg-cyan-300 text-slate-950 hover:bg-cyan-200">
                登录 OpenMIRA
                <ArrowRight className="size-4" />
              </Button>
              <Button size="lg" variant="outline" nativeButton={false} render={<a href="#capabilities" />} className="border-white/15 bg-transparent text-white hover:bg-white/8 hover:text-white">
                查看能力
              </Button>
            </div>
            <SignalPreview />
          </div>
        </section>

        <section id="capabilities" className="border-b border-white/10 px-5 py-16">
          <div className="mx-auto max-w-7xl">
            <p className="text-sm font-black text-teal-300">产品能力</p>
            <h2 className="mt-3 text-3xl font-black">从“自己翻消息”变成“让 Mira 先整理”</h2>
            <div className="mt-10 grid gap-5 md:grid-cols-2 lg:grid-cols-3">
              {capabilities.map(([Icon, title, copy]) => (
                <div key={title} className="rounded-md border border-white/10 bg-[#0c1d30] p-5">
                  <Icon className="size-6 text-cyan-300" />
                  <h3 className="mt-5 font-black">{title}</h3>
                  <p className="mt-2 text-sm leading-6 text-slate-400">{copy}</p>
                </div>
              ))}
            </div>
          </div>
        </section>

        <section id="safety" className="px-5 py-16">
          <div className="mx-auto grid max-w-7xl gap-6 lg:grid-cols-[0.85fr_1.15fr]">
            <div>
              <p className="text-sm font-black text-teal-300">安全边界</p>
              <h2 className="mt-3 text-3xl font-black">Mira 可以整理信息，但不会替你越权行动</h2>
            </div>
            <div className="grid gap-3 md:grid-cols-2">
              {[
                [LockKeyhole, '用户级数据隔离'],
                [ShieldCheck, '外部动作人工审批'],
                [CheckCircle2, '偏好变更可撤销'],
                [MonitorSmartphone, 'H5 / iOS / Android 同步体验'],
              ].map(([Icon, title]) => {
                const I = Icon as typeof LockKeyhole
                return (
                  <div key={title as string} className="flex items-center gap-3 rounded-md border border-white/10 bg-[#0c1d30] p-4">
                    <I className="size-5 text-teal-300" />
                    <span className="font-semibold">{title as string}</span>
                  </div>
                )
              })}
            </div>
          </div>
        </section>
      </main>

      <footer className="border-t border-white/10 px-5 py-8">
        <div className="mx-auto flex max-w-7xl flex-wrap items-center justify-between gap-4 text-xs text-slate-500">
          <Brand />
          <span className="flex items-center gap-2">
            <GitBranch className="size-4" />
            OpenMIRA · Web / iOS / Android
          </span>
        </div>
      </footer>
    </div>
  )
}

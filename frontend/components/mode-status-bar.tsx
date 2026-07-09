'use client'

import { Bot, UserCheck } from 'lucide-react'
import { useAppStore } from '@/lib/app-store'
import { cn } from '@/lib/utils'

export function ModeStatusBar() {
  const { workMode, toggleWorkMode } = useAppStore()
  const isWork = workMode === 'work'

  return (
    <button
      type="button"
      onClick={toggleWorkMode}
      className={cn(
        'flex w-full items-center justify-center gap-2 px-4 py-1.5 text-xs font-medium transition-colors',
        isWork
          ? 'bg-warning/10 text-warning'
          : 'bg-primary/10 text-primary',
      )}
      title="点击切换模式（原型演示）"
    >
      <span
        className={cn('size-2 rounded-full animate-pulse', isWork ? 'bg-warning' : 'bg-primary')}
        aria-hidden="true"
      />
      {isWork ? (
        <span className="inline-flex items-center gap-1.5">
          <UserCheck className="size-3.5" />
          工作时间 · 人工审核模式
        </span>
      ) : (
        <span className="inline-flex items-center gap-1.5">
          <Bot className="size-3.5" />
          非工作时间 · AI 自动回复模式
        </span>
      )}
      <span className="text-[10px] opacity-60">（点击切换）</span>
    </button>
  )
}

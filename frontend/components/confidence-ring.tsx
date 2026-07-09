import { cn } from '@/lib/utils'

export function ConfidenceRing({ score, className }: { score: number; className?: string }) {
  const pct = Math.round(score * 100)
  const radius = 15
  const circumference = 2 * Math.PI * radius
  const offset = circumference * (1 - score)
  const colorClass = score >= 0.8 ? 'text-success' : score >= 0.6 ? 'text-primary' : 'text-warning'

  return (
    <div className={cn('relative inline-flex size-10 items-center justify-center', className)}>
      <svg viewBox="0 0 36 36" className="size-10 -rotate-90" aria-hidden="true">
        <circle cx="18" cy="18" r={radius} fill="none" strokeWidth="3" className="stroke-muted" />
        <circle
          cx="18"
          cy="18"
          r={radius}
          fill="none"
          strokeWidth="3"
          strokeLinecap="round"
          strokeDasharray={circumference}
          strokeDashoffset={offset}
          className={cn('stroke-current transition-all', colorClass)}
        />
      </svg>
      <span className="absolute text-[9px] font-semibold tabular-nums">
        {pct}
        <span className="sr-only">% 置信度</span>
      </span>
    </div>
  )
}

import { ShieldCheck } from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { trustLevel, trustLevelConfig } from '@/lib/sop'
import { cn } from '@/lib/utils'

export function TrustBadge({ score, showScore = true, className }: { score: number; showScore?: boolean; className?: string }) {
  const level = trustLevelConfig[trustLevel(score)]
  return (
    <Badge
      variant="outline"
      className={cn('h-5 gap-1 px-1.5 text-[10px]', level.className, className)}
      title={`可信度评分 ${score}/100：该商机是否真实、安全`}
    >
      <ShieldCheck className="size-3" />
      {level.label}
      {showScore && <span className="opacity-80">{score}</span>}
    </Badge>
  )
}

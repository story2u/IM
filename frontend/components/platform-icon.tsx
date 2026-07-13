import { MessageSquare, Send } from 'lucide-react'
import { cn } from '@/lib/utils'
import type { Platform } from '@/lib/types'

export function PlatformIcon({ platform, className }: { platform: Platform; className?: string }) {
  if (platform === 'telegram') {
    return (
      <span
        className={cn(
          'inline-flex size-6 items-center justify-center rounded-full bg-sky-500/15 text-sky-600 dark:text-sky-400',
          className,
        )}
        title="Telegram"
        role="img"
        aria-label="Telegram"
      >
        <Send className="size-3.5" />
      </span>
    )
  }
  return (
    <span
      className={cn(
        'inline-flex size-6 items-center justify-center rounded-full bg-success/15 text-success',
        className,
      )}
      title="企业微信"
      role="img"
      aria-label="企业微信"
    >
      <MessageSquare className="size-3.5" />
    </span>
  )
}

export function platformLabel(platform: Platform) {
  return platform === 'telegram' ? 'Telegram' : '企业微信'
}

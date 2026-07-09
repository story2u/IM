import { Bot } from 'lucide-react'
import type { ChatMessage } from '@/lib/types'
import { cn } from '@/lib/utils'

function formatBubbleTime(iso: string) {
  const date = new Date(iso)
  return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
}

export function ChatBubble({ message }: { message: ChatMessage }) {
  const fromContact = message.isFromContact

  return (
    <div className={cn('flex w-full', fromContact ? 'justify-start' : 'justify-end')}>
      <div className={cn('flex max-w-[85%] flex-col gap-1 md:max-w-[70%]', !fromContact && 'items-end')}>
        <div className="relative">
          <div
            className={cn(
              'rounded-2xl px-3.5 py-2.5 text-sm leading-relaxed',
              fromContact
                ? 'rounded-tl-sm bg-muted text-foreground'
                : 'rounded-tr-sm bg-primary text-primary-foreground',
            )}
          >
            {message.content}
          </div>
          {message.source === 'ai' && (
            <span
              className="absolute -bottom-1.5 -left-1.5 flex size-5 items-center justify-center rounded-full border border-border bg-card text-primary shadow-sm"
              title="AI 自动回复"
            >
              <Bot className="size-3" />
              <span className="sr-only">AI 自动回复</span>
            </span>
          )}
        </div>
        <span className="px-1 text-[10px] text-muted-foreground">
          {formatBubbleTime(message.sentAt)}
          {message.source === 'ai' && ' · AI 代回复'}
          {message.source === 'human' && ' · 人工回复'}
        </span>
      </div>
    </div>
  )
}

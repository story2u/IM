'use client'

import { Clock3, CreditCard, LayoutGrid, Radar, Send, Settings, Smartphone } from 'lucide-react'
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { cn } from '@/lib/utils'

const nav = [
  ['看板', '/demo', LayoutGrid],
  ['设置', '/demo/settings', Settings],
  ['多端', '/#apps', Smartphone],
] as const

export function DemoShell({ children }: { children: React.ReactNode }) {
  const pathname = usePathname()
  return (
    <div className="flex min-h-svh bg-background" data-testid="demo-shell">
      <aside className="fixed inset-y-0 left-0 z-30 hidden w-56 border-r bg-sidebar md:flex md:flex-col">
        <Link href="/" className="flex items-center gap-2.5 px-5 py-5"><span className="grid size-8 place-items-center rounded-md bg-primary text-primary-foreground"><Radar className="size-4" /></span><span><strong className="block text-sm">商机雷达</strong><small className="text-muted-foreground">安全演示模式</small></span></Link>
        <nav className="space-y-1 px-3">{nav.map(([label, href, Icon]) => <Link key={href} href={href} className={cn('flex items-center gap-2.5 rounded-md px-3 py-2 text-sm', pathname === href || (href !== '/demo' && pathname.startsWith(href)) ? 'bg-sidebar-accent font-medium text-sidebar-accent-foreground' : 'text-muted-foreground')}><Icon className="size-4" />{label}</Link>)}</nav>
        <div className="mt-auto border-t p-4 text-xs text-muted-foreground"><p className="flex items-center gap-2"><Send className="size-3.5" />不连接真实 Telegram</p><p className="mt-2 flex items-center gap-2"><CreditCard className="size-3.5" />不执行真实支付</p></div>
      </aside>
      <div className="min-w-0 flex-1 md:pl-56"><div className="sticky top-0 z-20 flex h-11 items-center justify-between border-b bg-background/95 px-4 backdrop-blur md:px-7"><span className="flex items-center gap-2 text-xs font-medium text-success"><span className="size-1.5 rounded-full bg-success" />确定性演示数据</span><span className="flex items-center gap-1.5 text-xs text-muted-foreground"><Clock3 className="size-3.5" />2026-07-13 09:30</span></div>{children}</div>
    </div>
  )
}

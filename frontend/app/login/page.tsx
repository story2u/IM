'use client'

import { Apple, Loader2, Mail, Radar, ShieldCheck } from 'lucide-react'
import { useRouter } from 'next/navigation'
import { useEffect, useRef, useState } from 'react'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { fetchOAuthAuthorizeUrl } from '@/lib/api'
import { useAuth } from '@/lib/auth'
import type { OAuthProvider } from '@/lib/types'

const providerLabels: Record<OAuthProvider, string> = {
  google: '使用 Google 登录',
  apple: '使用 Apple 登录',
}

export default function LoginPage() {
  const router = useRouter()
  const { user, completeOAuth } = useAuth()
  const [loadingProvider, setLoadingProvider] = useState<OAuthProvider | null>(null)
  const [processingCallback, setProcessingCallback] = useState(false)
  const [error, setError] = useState('')
  const handledCallbackRef = useRef(false)

  useEffect(() => {
    if (user) {
      router.replace('/')
    }
  }, [router, user])

  useEffect(() => {
    if (handledCallbackRef.current || typeof window === 'undefined') return
    const hashParams = new URLSearchParams(window.location.hash.replace(/^#/, ''))
    const token = hashParams.get('token') ?? new URLSearchParams(window.location.search).get('token')
    if (!token) return

    handledCallbackRef.current = true
    setProcessingCallback(true)
    setError('')
    completeOAuth(token)
      .then(() => router.replace('/'))
      .catch((exc) => {
        setError(exc instanceof Error ? exc.message : '登录回调失败')
        window.history.replaceState(null, '', '/login')
      })
      .finally(() => setProcessingCallback(false))
  }, [completeOAuth, router])

  async function startOAuth(provider: OAuthProvider) {
    setError('')
    setLoadingProvider(provider)
    try {
      const authorizationUrl = await fetchOAuthAuthorizeUrl(provider)
      window.location.assign(authorizationUrl)
    } catch (exc) {
      setError(exc instanceof Error ? exc.message : '无法发起登录')
      setLoadingProvider(null)
    }
  }

  return (
    <div className="grid min-h-svh place-items-center px-4 py-10">
      <div className="w-full max-w-sm">
        <div className="mb-5 flex items-center gap-2">
          <span className="flex size-9 items-center justify-center rounded-lg bg-primary text-primary-foreground">
            <Radar className="size-5" />
          </span>
          <div>
            <h1 className="text-lg font-semibold tracking-tight">商机雷达</h1>
            <p className="text-xs text-muted-foreground">IM 商机助手</p>
          </div>
        </div>

        <Card className="gap-5 p-5 shadow-sm">
          <div className="flex items-center gap-2">
            <ShieldCheck className="size-4 text-muted-foreground" />
            <h2 className="text-base font-semibold">登录</h2>
          </div>

          <div className="grid gap-3">
            <Button
              type="button"
              variant="outline"
              className="justify-start"
              onClick={() => startOAuth('google')}
              disabled={Boolean(loadingProvider) || processingCallback}
            >
              {loadingProvider === 'google' ? <Loader2 className="size-4 animate-spin" /> : <Mail className="size-4" />}
              {providerLabels.google}
            </Button>
            <Button
              type="button"
              className="justify-start bg-zinc-950 text-white hover:bg-zinc-800 dark:bg-white dark:text-zinc-950 dark:hover:bg-zinc-200"
              onClick={() => startOAuth('apple')}
              disabled={Boolean(loadingProvider) || processingCallback}
            >
              {loadingProvider === 'apple' ? <Loader2 className="size-4 animate-spin" /> : <Apple className="size-4" />}
              {providerLabels.apple}
            </Button>
          </div>

          {processingCallback && (
            <p className="rounded-md bg-muted px-3 py-2 text-xs text-muted-foreground">正在完成登录</p>
          )}
          {error && <p className="rounded-md bg-destructive/10 px-3 py-2 text-xs text-destructive">{error}</p>}
        </Card>
      </div>
    </div>
  )
}

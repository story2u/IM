'use client'

import { ArrowLeft, Loader2 } from 'lucide-react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useEffect, useState } from 'react'
import { BrandLogo } from '@/components/brand-logo'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { confirmPasswordReset, requestPasswordReset } from '@/lib/api'

export default function ForgotPasswordPage() {
  const router = useRouter()
  const [email, setEmail] = useState('')
  const [code, setCode] = useState('')
  const [newPassword, setNewPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [requested, setRequested] = useState(false)
  const [loading, setLoading] = useState(false)
  const [message, setMessage] = useState('')
  const [error, setError] = useState('')

  useEffect(() => {
    const value = new URLSearchParams(window.location.search).get('email')
    if (value) setEmail(value)
    window.history.replaceState(null, '', '/forgot-password')
  }, [])

  async function requestCode(event: React.FormEvent) {
    event.preventDefault()
    setLoading(true)
    setError('')
    try {
      const result = await requestPasswordReset(email)
      setMessage(result.message)
      setRequested(true)
    } catch (exc) {
      setError(exc instanceof Error ? exc.message : '无法发送重置邮件')
    } finally {
      setLoading(false)
    }
  }

  async function resetWithCode(event: React.FormEvent) {
    event.preventDefault()
    if (newPassword !== confirmPassword) {
      setError('两次输入的新密码不一致')
      return
    }
    setLoading(true)
    setError('')
    try {
      await confirmPasswordReset({ email, code, newPassword })
      router.replace('/login?passwordReset=1')
    } catch (exc) {
      setError(exc instanceof Error ? exc.message : '密码重置失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <main className="flex min-h-svh items-center justify-center bg-[#202020] px-5 py-10 text-white">
      <section className="w-full max-w-md" aria-labelledby="forgot-password-heading">
        <BrandLogo size={56} className="mb-7" />
        <h1 id="forgot-password-heading" className="text-3xl font-semibold">重置密码</h1>
        <p className="mt-3 text-sm leading-6 text-zinc-400">
          输入登录邮箱。无论账户是否存在，系统都会返回相同提示。
        </p>

        <form className="mt-8 grid gap-5" onSubmit={requested ? resetWithCode : requestCode}>
          <div className="grid gap-2">
            <Label htmlFor="email">邮箱</Label>
            <Input id="email" type="email" autoComplete="email" required maxLength={320} value={email}
              onChange={(event) => setEmail(event.target.value)} disabled={requested || loading}
              className="h-12 border-white/20 bg-white/5 text-white" />
          </div>
          {requested ? (
            <>
              <div className="grid gap-2">
                <Label htmlFor="code">邮件验证码</Label>
                <Input id="code" autoComplete="one-time-code" required minLength={10} maxLength={10}
                  value={code} onChange={(event) => setCode(event.target.value.toUpperCase().replace(/\s/g, ''))}
                  className="h-12 border-white/20 bg-white/5 font-mono uppercase text-white" />
              </div>
              <div className="grid gap-2">
                <Label htmlFor="new-password">新密码</Label>
                <Input id="new-password" type="password" autoComplete="new-password" required minLength={10} maxLength={128}
                  value={newPassword} onChange={(event) => setNewPassword(event.target.value)}
                  className="h-12 border-white/20 bg-white/5 text-white" />
                <p className="text-xs text-zinc-500">至少 10 个字符</p>
              </div>
              <div className="grid gap-2">
                <Label htmlFor="confirm-password">确认新密码</Label>
                <Input id="confirm-password" type="password" autoComplete="new-password" required minLength={10} maxLength={128}
                  value={confirmPassword} onChange={(event) => setConfirmPassword(event.target.value)}
                  className="h-12 border-white/20 bg-white/5 text-white" />
              </div>
            </>
          ) : null}
          <Button type="submit" disabled={loading} className="h-12 rounded-full bg-white text-zinc-950 hover:bg-zinc-200">
            {loading ? <Loader2 className="size-5 animate-spin" /> : null}
            {requested ? '确认重置' : '发送重置邮件'}
          </Button>
        </form>

        <div className="mt-5 min-h-10" aria-live="polite">
          {message ? <p className="text-sm text-emerald-300">{message}</p> : null}
          {error ? <p role="alert" className="text-sm text-red-300">{error}</p> : null}
        </div>
        <Link href="/login" className="mt-5 inline-flex items-center gap-2 text-sm text-zinc-300 hover:text-white">
          <ArrowLeft className="size-4" /> 返回登录
        </Link>
      </section>
    </main>
  )
}

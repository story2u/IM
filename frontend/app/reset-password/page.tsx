'use client'

import { Loader2 } from 'lucide-react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useEffect, useState } from 'react'
import { BrandLogo } from '@/components/brand-logo'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { confirmPasswordReset } from '@/lib/api'

export default function ResetPasswordPage() {
  const router = useRouter()
  const [token, setToken] = useState('')
  const [newPassword, setNewPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  useEffect(() => {
    const value = new URLSearchParams(window.location.search).get('token') ?? ''
    setToken(value)
    window.history.replaceState(null, '', '/reset-password')
  }, [])

  async function submit(event: React.FormEvent) {
    event.preventDefault()
    if (!token) {
      setError('重置链接无效，请重新申请')
      return
    }
    if (newPassword !== confirmPassword) {
      setError('两次输入的新密码不一致')
      return
    }
    setLoading(true)
    setError('')
    try {
      await confirmPasswordReset({ token, newPassword })
      router.replace('/login?passwordReset=1')
    } catch (exc) {
      setError(exc instanceof Error ? exc.message : '密码重置失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <main className="flex min-h-svh items-center justify-center bg-[#202020] px-5 py-10 text-white">
      <section className="w-full max-w-md">
        <BrandLogo size={56} className="mb-7" />
        <h1 className="text-3xl font-semibold">设置新密码</h1>
        <form className="mt-8 grid gap-5" onSubmit={submit}>
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
          <Button type="submit" disabled={loading || !token} className="h-12 rounded-full bg-white text-zinc-950 hover:bg-zinc-200">
            {loading ? <Loader2 className="size-5 animate-spin" /> : null} 确认重置
          </Button>
        </form>
        {error ? <p role="alert" className="mt-5 text-sm text-red-300">{error}</p> : null}
        <Link href="/forgot-password" className="mt-7 inline-block text-sm text-zinc-300 hover:text-white">重新申请重置邮件</Link>
      </section>
    </main>
  )
}

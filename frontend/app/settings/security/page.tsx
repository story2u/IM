'use client'

import { KeyRound, Loader2 } from 'lucide-react'
import { useRouter } from 'next/navigation'
import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { changePassword } from '@/lib/api'
import { useAuth } from '@/lib/auth'

export default function SecuritySettingsPage() {
  const router = useRouter()
  const { user, logout } = useAuth()
  const [currentPassword, setCurrentPassword] = useState('')
  const [newPassword, setNewPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  async function submit(event: React.FormEvent) {
    event.preventDefault()
    if (newPassword !== confirmPassword) {
      setError('两次输入的新密码不一致')
      return
    }
    setLoading(true)
    setError('')
    try {
      await changePassword(currentPassword, newPassword)
      logout()
      router.replace('/login?passwordChanged=1')
    } catch (exc) {
      setError(exc instanceof Error ? exc.message : '密码修改失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="mx-auto w-full max-w-2xl px-4 py-6 md:px-8">
      <header className="mb-6">
        <h1 className="text-2xl font-semibold">账户安全</h1>
        <p className="mt-1 text-sm text-muted-foreground">修改密码后，所有设备需要重新登录。</p>
      </header>
      <Card className="p-5">
        <div className="mb-5 flex items-center gap-3">
          <KeyRound className="size-5 text-muted-foreground" />
          <div><p className="font-medium">登录密码</p><p className="text-sm text-muted-foreground">{user?.email}</p></div>
        </div>
        {user?.hasPassword ? (
          <form className="grid gap-4" onSubmit={submit}>
            <div className="grid gap-2"><Label htmlFor="current-password">当前密码</Label><Input id="current-password" type="password" autoComplete="current-password" required maxLength={128} value={currentPassword} onChange={(event) => setCurrentPassword(event.target.value)} /></div>
            <div className="grid gap-2"><Label htmlFor="new-password">新密码</Label><Input id="new-password" type="password" autoComplete="new-password" required minLength={10} maxLength={128} value={newPassword} onChange={(event) => setNewPassword(event.target.value)} /><p className="text-xs text-muted-foreground">至少 10 个字符</p></div>
            <div className="grid gap-2"><Label htmlFor="confirm-password">确认新密码</Label><Input id="confirm-password" type="password" autoComplete="new-password" required minLength={10} maxLength={128} value={confirmPassword} onChange={(event) => setConfirmPassword(event.target.value)} /></div>
            <Button type="submit" disabled={loading} className="w-fit">{loading ? <Loader2 className="size-4 animate-spin" /> : null} 修改密码</Button>
          </form>
        ) : (
          <div>
            <p className="text-sm text-muted-foreground">你当前通过 OAuth 登录。设置密码前需要验证登录邮箱。</p>
            <Button
              type="button"
              className="mt-4"
              onClick={() => router.push(`/forgot-password?email=${encodeURIComponent(user?.email ?? '')}`)}
            >
              验证邮箱并设置密码
            </Button>
          </div>
        )}
        {error ? <p role="alert" className="mt-4 text-sm text-destructive">{error}</p> : null}
      </Card>
    </div>
  )
}

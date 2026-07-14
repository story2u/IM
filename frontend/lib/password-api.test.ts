import { afterEach, describe, expect, it, vi } from 'vitest'
import { changePassword, confirmPasswordReset, requestPasswordReset } from './api'

afterEach(() => {
  vi.unstubAllGlobals()
})

function mockJsonResponse(body: unknown) {
  const fetchMock = vi.fn().mockResolvedValue(
    new Response(JSON.stringify(body), {
      status: 200,
      headers: { 'Content-Type': 'application/json' },
    }),
  )
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

describe('password API', () => {
  it('requests reset without exposing account state in the client payload', async () => {
    const fetchMock = mockJsonResponse({ message: '如果该邮箱已注册，重置邮件将在几分钟内送达' })

    await requestPasswordReset('member@example.com')

    expect(fetchMock).toHaveBeenCalledWith('/api/v1/auth/password/reset/request', expect.objectContaining({
      method: 'POST',
      body: JSON.stringify({ email: 'member@example.com' }),
    }))
  })

  it('submits only the chosen reset credential and new password', async () => {
    const fetchMock = mockJsonResponse({ message: '密码已重置，请使用新密码登录' })

    await confirmPasswordReset({
      email: 'member@example.com',
      code: 'ABCDEFGH23',
      newPassword: 'new-password-123',
    })

    expect(fetchMock).toHaveBeenCalledWith('/api/v1/auth/password/reset/confirm', expect.objectContaining({
      body: JSON.stringify({
        email: 'member@example.com',
        code: 'ABCDEFGH23',
        newPassword: 'new-password-123',
      }),
    }))
  })

  it('sends current and new password only to the authenticated change endpoint', async () => {
    const fetchMock = mockJsonResponse({ message: '密码已修改，请重新登录' })

    await changePassword('old-password', 'new-password-123')

    expect(fetchMock).toHaveBeenCalledWith('/api/v1/auth/password/change', expect.objectContaining({
      method: 'POST',
      body: JSON.stringify({ currentPassword: 'old-password', newPassword: 'new-password-123' }),
    }))
  })
})

import { describe, expect, it } from 'vitest'
import { isPublicAuthPath } from './auth-routes'

describe('public authentication routes', () => {
  it.each(['/login', '/forgot-password', '/reset-password'])(
    'allows unauthenticated access to %s',
    (pathname) => {
      expect(isPublicAuthPath(pathname)).toBe(true)
    },
  )

  it('keeps authenticated settings private', () => {
    expect(isPublicAuthPath('/settings/security')).toBe(false)
  })
})

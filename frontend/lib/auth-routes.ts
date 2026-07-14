const publicAuthPaths = new Set(['/login', '/forgot-password', '/reset-password'])

export function isPublicAuthPath(pathname: string): boolean {
  return publicAuthPaths.has(pathname)
}

import { NextResponse } from 'next/server'

export function proxy() {
  if (process.env.DEMO_MODE !== 'true') {
    return new NextResponse('Not Found', { status: 404 })
  }
  return NextResponse.next()
}

export const config = {
  matcher: ['/demo/:path*'],
}

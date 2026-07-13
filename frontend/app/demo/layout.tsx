import { notFound } from 'next/navigation'

export default function DemoLayout({ children }: { children: React.ReactNode }) {
  if (process.env.DEMO_MODE !== 'true') notFound()
  return children
}

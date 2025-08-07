import './styles.css'
import type { ReactNode } from 'react'

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="en">
      <body className="min-h-screen bg-slate-50 text-slate-900">
        <header className="p-4 border-b bg-white">
          <h1 className="text-xl font-semibold">RouteIQ Dashboard</h1>
        </header>
        <main className="p-4">{children}</main>
      </body>
    </html>
  )
}

import type { Metadata } from 'next'
import { Geist, Geist_Mono, Inter } from 'next/font/google'
import './globals.css'
import { QueryProvider } from '@/providers/QueryProvider'
import { ThemeProvider } from '@/components/layout/ThemeProvider'
import { ToastProvider } from '@/components/layout/ToastProvider'
import { NuqsAdapter } from 'nuqs/adapters/next/app'

const inter = Inter({ subsets: ['latin'], variable: '--font-sans' })

const geistSans = Geist({
  variable: '--font-geist-sans',
  subsets: ['latin'],
})

const geistMono = Geist_Mono({
  variable: '--font-geist-mono',
  subsets: ['latin'],
})

export const metadata: Metadata = {
  title: 'Organiza Aqui',
  description: 'Sistema de organização',
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang="pt-BR" className={inter.variable} suppressHydrationWarning>
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        <ThemeProvider
          attribute="class"
          defaultTheme="system"
          enableSystem
          disableTransitionOnChange
        >
          <NuqsAdapter>
            <QueryProvider>
              {children}
              <ToastProvider />
            </QueryProvider>
          </NuqsAdapter>
        </ThemeProvider>
      </body>
    </html>
  )
}

import AppLayoutContent from '@/components/AppLayoutContent'

export default function AppLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return <AppLayoutContent>{children}</AppLayoutContent>
}

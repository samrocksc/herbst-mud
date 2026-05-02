type StatGridProps = Readonly<{
  children: React.ReactNode
  cols?: number
}>

export function StatGrid({ children, cols = 5 }: StatGridProps) {
  const colClass = cols === 4
    ? 'grid-cols-[repeat(auto-fit,minmax(200px,1fr))]'
    : 'grid-cols-[repeat(auto-fit,minmax(200px,1fr))]'

  return (
    <div className={`grid ${colClass} gap-4`}>
      {children}
    </div>
  )
}

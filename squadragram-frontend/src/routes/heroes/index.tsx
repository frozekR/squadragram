import { HeroesList } from '#/components/HeroesList'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/heroes/')({
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <div className="p-8">
      <div className="flex flex-col items-center justify-center">
        <h1 className="text-3xl font-bold mb-4">Heroes</h1>
        </div>
      <section className="mt-8 p-2">
        <h2 className="text-2xl font-semibold mb-2">All Heroes</h2>
        <HeroesList />
      </section>
    </div>
  )
}

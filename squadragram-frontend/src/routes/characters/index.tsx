import { CharactersList } from '#/components/CharactersList'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/characters/')({
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <div className="p-8">
      <div className="flex flex-col items-center justify-center">
        <h1 className="text-3xl font-bold mb-4">Characters</h1>
        </div>
      <section className="mt-8 p-2">
        <h2 className="text-2xl font-semibold mb-2">All Characters</h2>
        <CharactersList />
      </section>
    </div>
  )
}

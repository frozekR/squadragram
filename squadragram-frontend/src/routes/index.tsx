import { CharactersList } from '#/components/CharactersList'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/')({
    head: () => ({
    meta: [
      {
        charSet: 'utf-8',
      },
      {
        name: 'viewport',
        content: 'width=device-width, initial-scale=1',
      },
      {
        title: 'Squadragram',
      },
    ],
    links: [
      {
        rel: 'stylesheet',
        href: '/styles.css',
      },
    ],
  }),
  component: Home
})

function Home() {
  return (
    <div className="p-8">
      <div className="flex flex-col items-center justify-center">
        <h1 className="text-3xl font-bold mb-4">Welcome to Squadragram!</h1>
      </div>
      <section className="mt-8 p-2">
        <h2 className="text-2xl font-semibold mb-2">Featured Characters</h2>
        <div className="flex flex-wrap gap-4">
        <CharactersList></CharactersList>
      </div>
      </section>
      <section className="mt-8 p-2">
        <h2 className="text-2xl font-semibold mb-2">Latest News</h2>
      </section>
    </div>
  )
}

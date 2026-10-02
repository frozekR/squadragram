Welcome to your new TanStack Start app!

## Админка Squadragram

Откройте `http://localhost:3000/admin/` или пункт **Admin** в меню. Для входа
используйте значение `ADMIN_API_TOKEN` из `.env` в корне репозитория. Backend должен
быть пересобран с этой переменной. Токен вводится вручную, хранится только в памяти
вкладки и исчезает при выходе или обновлении страницы.

1. Выберите героя слева или нажмите **Добавить**. Заполните имя, роль и описание.
2. После сохранения героя загрузите его иконку, рендер и демо-видео в блоке **Медиа**.
3. На вкладке **Скиллы** добавляйте и редактируйте умения. Поле **Герой** позволяет
   перенести скилл; его UUID и собственные медиа сохраняются.
4. В медиа сначала выберите файл и проверьте предпросмотр, затем нажмите
   **Загрузить** или **Заменить**. Удаление требует отдельного подтверждения.

На публичной вкладке Skills иконки сгруппированы по типам, а ниже показан только
выбранный скилл. Порядок категорий: Passive → Rush attack → Skill → Super attack →
Max super attack → Transformation. Категории без скиллов скрыты.
В редакторе скилла поле **Порядок в категории** управляет положением его иконки:
меньшее число — левее. Измените число и нажмите **Сохранить изменения**. При равных
значениях используется порядок создания. Поле **Тип** меняет категорию скилла.
Проверка группировки: `node --experimental-strip-types --test tests/skill-order.test.mjs`.

Для изображений доступны PNG/JPEG/GIF/WebP, для видео MP4/WebM, размер до 32 MiB.
При переходе между редакторами интерфейс предупреждает о несохранённых полях и
выбранных файлах. Сохранение обновляет данные в админке и кеше публичных страниц;
уже открытая в другой вкладке публичная страница обновляется после перезагрузки.
Операции записи защищены на backend, включая запросы вне админки.

API и интеграционные проверки описаны в `../squadragram-backend/README.md`.

Админка оформлена утилитами Tailwind CSS; общие стили компонентов находятся в
`src/components/admin/styles.ts`. На телефоне список героев расположен над
редактором, формы и медиа показаны в одну колонку. От 640 px парные поля идут
рядом, от 1024 px появляется боковая колонка героев. Кнопки имеют высоту не менее
44 px, поля на телефоне используют шрифт 16 px, диалоги помещаются в экран.

# Getting Started

To run this application:

```bash
npm install
npm run dev
```

# Building For Production

To build this application for production:

```bash
npm run build
```

## Styling

This project uses [Tailwind CSS](https://tailwindcss.com/) for styling.

### Removing Tailwind CSS

If you prefer not to use Tailwind CSS:

1. Remove the demo pages in `src/routes/demo/`
2. Replace the Tailwind import in `src/styles.css` with your own styles
3. Remove `tailwindcss()` from the plugins array in `vite.config.ts`
4. Remove `@tailwindcss/vite` and `tailwindcss` from `package.json`

## Linting & Formatting

This project uses [Biome](https://biomejs.dev/) for linting and formatting. The following scripts are available:


```bash
npm run lint
npm run format
npm run check
```



## Routing

This project uses [TanStack Router](https://tanstack.com/router) with file-based routing. Routes are managed as files in `src/routes`.

### Adding A Route

To add a new route to your application just add a new file in the `./src/routes` directory.

TanStack will automatically generate the content of the route file for you.

Now that you have two routes you can use a `Link` component to navigate between them.

### Adding Links

To use SPA (Single Page Application) navigation you will need to import the `Link` component from `@tanstack/react-router`.

```tsx
import { Link } from "@tanstack/react-router";
```

Then anywhere in your JSX you can use it like so:

```tsx
<Link to="/about">About</Link>
```

This will create a link that will navigate to the `/about` route.

More information on the `Link` component can be found in the [Link documentation](https://tanstack.com/router/v1/docs/framework/react/api/router/linkComponent).

### Using A Layout

In the File Based Routing setup the layout is located in `src/routes/__root.tsx`. Anything you add to the root route will appear in all the routes. The route content will appear in the JSX where you render `{children}` in the `shellComponent`.

Here is an example layout that includes a header:

```tsx
import { HeadContent, Scripts, createRootRoute } from '@tanstack/react-router'

export const Route = createRootRoute({
  head: () => ({
    meta: [
      { charSet: 'utf-8' },
      { name: 'viewport', content: 'width=device-width, initial-scale=1' },
      { title: 'My App' },
    ],
  }),
  shellComponent: ({ children }) => (
    <html lang="en">
      <head>
        <HeadContent />
      </head>
      <body>
        <header>
          <nav>
            <Link to="/">Home</Link>
            <Link to="/about">About</Link>
          </nav>
        </header>
        {children}
        <Scripts />
      </body>
    </html>
  ),
})
```

More information on layouts can be found in the [Layouts documentation](https://tanstack.com/router/latest/docs/framework/react/guide/routing-concepts#layouts).

## Server Functions

TanStack Start provides server functions that allow you to write server-side code that seamlessly integrates with your client components.

```tsx
import { createServerFn } from '@tanstack/react-start'

const getServerTime = createServerFn({
  method: 'GET',
}).handler(async () => {
  return new Date().toISOString()
})

// Use in a component
function MyComponent() {
  const [time, setTime] = useState('')
  
  useEffect(() => {
    getServerTime().then(setTime)
  }, [])
  
  return <div>Server time: {time}</div>
}
```

## API Routes

You can create API routes by using the `server` property in your route definitions:

```tsx
import { createFileRoute } from '@tanstack/react-router'
import { json } from '@tanstack/react-start'

export const Route = createFileRoute('/api/hello')({
  server: {
    handlers: {
      GET: () => json({ message: 'Hello, World!' }),
    },
  },
})
```

## Data Fetching

There are multiple ways to fetch data in your application. You can use TanStack Query to fetch data from a server. But you can also use the `loader` functionality built into TanStack Router to load the data for a route before it's rendered.

For example:

```tsx
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/people')({
  loader: async () => {
    const response = await fetch('https://swapi.dev/api/people')
    return response.json()
  },
  component: PeopleComponent,
})

function PeopleComponent() {
  const data = Route.useLoaderData()
  return (
    <ul>
      {data.results.map((person) => (
        <li key={person.name}>{person.name}</li>
      ))}
    </ul>
  )
}
```

Loaders simplify your data fetching logic dramatically. Check out more information in the [Loader documentation](https://tanstack.com/router/latest/docs/framework/react/guide/data-loading#loader-parameters).



# Learn More

You can learn more about all of the offerings from TanStack in the [TanStack documentation](https://tanstack.com).

For TanStack Start specific documentation, visit [TanStack Start](https://tanstack.com/start).
# Медиа персонажей и скиллов

Запуск фронтенда: `npm run dev` (порт 3000). Backend должен быть доступен на
`http://localhost:3001/api`. Для другого адреса задайте `VITE_API_BASE_URL` в
`.env.local` и перезапустите Vite.

Карточки открывают `/characters/{id}` с последовательным ID персонажа. Для прямой
ссылки фронтенд находит UUID по ID в списке персонажей, затем использует UUID
для запросов API. На странице персонажа вкладка **Skills**
загружает `/characters/{uuid}/skills` и показывает иконки, рендеры и DEMO-видео.
Файлы берутся через `/media-metadata/{id}/file`, а не напрямую из приватного MinIO.
В списке используется ICON (или RENDER); в профиле — RENDER (или ICON).

После загрузки нового медиа через Postman обновите страницу или нажмите **Refresh**
на вкладке Skills. Если изображение отсутствует или недоступно, показывается заглушка;
ошибки загрузки видео и данных допускают повторную попытку.

Проверки: `node node_modules/typescript/bin/tsc --noEmit` и `npm run build`.

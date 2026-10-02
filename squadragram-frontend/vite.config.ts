import { defineConfig } from 'vite'
import { devtools } from '@tanstack/devtools-vite'

import { tanstackStart } from '@tanstack/react-start/plugin/vite'

import viteReact from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

const config = defineConfig({
  resolve: { tsconfigPaths: true },
  // Console forwarding loops with SSR logs; keep the devtools UI without piping.
  plugins: [devtools({ consolePiping: { enabled: false } }), tailwindcss(), tanstackStart(), viteReact()],
})

export default config

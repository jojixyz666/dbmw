/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: 'var(--color-primary)',
          container: 'var(--color-primary-container)',
          fixed: '#d1fae5',
          'fixed-dim': '#a7f3d0',
        },
        'on-primary': 'var(--color-on-primary)',
        'on-primary-container': 'var(--color-on-primary-container)',
        secondary: {
          DEFAULT: '#006c49',
          container: '#005236',
        },
        surface: {
          DEFAULT: 'var(--bg-surface)',
          dim: 'var(--bg-surface-lowest)',
          variant: 'var(--bg-surface-container)',
          container: 'var(--bg-surface-container)',
          'container-low': 'var(--bg-surface-low)',
          'container-lowest': 'var(--bg-surface-lowest)',
          'container-high': 'var(--bg-surface-high)',
          'container-highest': 'var(--bg-surface-high)',
        },
        'on-surface': 'var(--text-on-surface)',
        'on-surface-variant': 'var(--text-on-surface-variant)',
        outline: {
          DEFAULT: 'var(--border-outline)',
          variant: 'var(--border-outline-variant)',
        },
        brand: {
          50: '#ecfdf5',
          100: '#d1fae5',
          200: '#a7f3d0',
          300: '#6ee7b7',
          400: '#34d399',
          500: '#10b981',
          600: '#059669',
          700: '#047857',
          800: '#065f46',
          900: '#064e3b',
        },
        dark: {
          950: '#070a11',
          900: '#0b0f19',
          850: '#101726',
          800: '#151e32',
          750: '#1c2840',
          700: '#24334f',
          600: '#33476b',
          500: '#475e85',
        }
      },
      fontFamily: {
        sans: ['Plus Jakarta Sans', 'Inter', 'sans-serif'],
        mono: ['JetBrains Mono', 'Fira Code', 'monospace'],
        'body-base': ['Plus Jakarta Sans', 'Inter', 'sans-serif'],
        'code-base': ['JetBrains Mono', 'monospace'],
        'headline-md': ['Plus Jakarta Sans', 'sans-serif'],
        'headline-lg': ['Plus Jakarta Sans', 'sans-serif'],
        'label-caps': ['Plus Jakarta Sans', 'sans-serif'],
      }
    },
  },
  plugins: [],
}
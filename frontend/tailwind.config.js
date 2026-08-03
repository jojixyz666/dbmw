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
        brand: {
          50: '#ecfdf5',
          100: '#d1fae5',
          500: '#10b981',
          600: '#059669',
          700: '#047857',
        },
        dark: {
          900: '#0b0f19',
          800: '#131b2e',
          700: '#1e293b',
          600: '#334155',
          500: '#475569',
        }
      }
    },
  },
  plugins: [],
}

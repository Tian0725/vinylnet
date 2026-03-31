/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        brand: {
          navy: {
            dark: '#071527',    /* Background Base */
            base: '#0a1a35',    /* Official Darker Navy */
            surface: '#162e55', /* Shirt Navy (Original) */
          },
          royal: '#2563eb',     /* Dynamic Royal Blue */
          accent: '#f97316',    /* Brand Accent */
        },
      },
    },
  },
  plugins: [],
}

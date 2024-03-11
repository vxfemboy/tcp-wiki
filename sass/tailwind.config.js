/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["../assets/**/*.{html,js}"],
  theme: {
    extend: {},
  },
  daisyui: {
    themes: ["forest"]
  },
  plugins: [require("@tailwindcss/typography"), require("daisyui")],
}


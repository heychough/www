/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/**/*.html"],
  theme: {
    extend: {
        fontFamily: {
            sans: ["Inter", "system-ui", "sans-serif"]
        }
    },
  },
  plugins: [
    require("@tailwindcss/typography"),
  ],
}


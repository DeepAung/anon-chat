// const { colors } = require("tailwindcss/defaultTheme");

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/**/*.{html,js,templ,go}"],
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: "rgb(94, 234, 212)",
          dark: "rgb(20, 184, 166)",
        },
      },
    },
    fontFamily: {
      mono: ["Fira Code", "ui-monospace", "SFMono-Regular"],
    },
  },
  plugins: [],
};

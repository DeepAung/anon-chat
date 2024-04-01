/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/**/*.{html,js,templ,go}"],
  theme: {
    extend: {},
    fontFamily: {
      mono: ["Fira Code", "ui-monospace", "SFMono-Regular"],
    },
  },
  plugins: [],
};

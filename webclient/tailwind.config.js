/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      colors: {
        primary: "#2C3639",
        primaryLight: "#3F4E4F",
        secondary: "#DCD7C9",
        secondaryDark: "#A27B5C",
      },
    },
  },
  plugins: [require("daisyui")],
};

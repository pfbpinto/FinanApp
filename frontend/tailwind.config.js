/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./src/**/*.{js,jsx,ts,tsx}", // Garante que Tailwind analise os arquivos do React
    "./public/index.html", // Adiciona o arquivo base do React
  ],
  theme: {
    extend: {},
  },
  plugins: [],
};

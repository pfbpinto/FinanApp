/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./views/**/*.html", // Corrigir para a pasta de views
    "./assets/**/*.html", // Remover se não houver arquivos HTML dentro de assets
  ],
  theme: {
    extend: {},
  },
  plugins: [],
};

// src/index.js
import React from "react";
import ReactDOM from "react-dom/client";
import "./index.css";
import App from "./App";
import reportWebVitals from "./reportWebVitals";

// Importa o Font Awesome e Toastr aqui para carregá-los globalmente
import "font-awesome/css/font-awesome.min.css"; // Font Awesome
import "toastr/build/toastr.min.css"; // Toastr CSS

import toastr from "toastr"; // Importando Toastr para configuração global

import { AuthProvider } from "./components/AuthContext"; // Importa o AuthProvider

// Configuração global para o Toastr
toastr.options = {
  positionClass: "toast-top-right", // Posição no canto superior direito
  timeOut: 3000, // Tempo de exibição
  extendedTimeOut: 1000, // Tempo de exibição ao passar o mouse
  preventDuplicates: true, // Evitar notificações duplicadas
  closeButton: true, // Exibir botão de fechar
  progressBar: true, // Barra de progresso
  toastClass: "custom-toast", // Adicionando uma classe personalizada
};

const root = ReactDOM.createRoot(document.getElementById("root"));
root.render(
  <AuthProvider>
    <App />
  </AuthProvider>
);

reportWebVitals();

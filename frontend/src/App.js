// src/App.js
import React from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import Header from "./layouts/Header";
import Footer from "./layouts/Footer";
import Home from "./pages/Home";
import Login from "./pages/Login";
import User from "./pages/User";
import UserEdit from "./pages/UserPage/Edit";
import Register from "./pages/Register";
import NotFoundPage from "./pages/NotFoundPage";

function App() {
  return (
    <Router>
      <div className="flex flex-col min-h-screen">
        <Header />
        <main className="flex-grow pt-20">
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/login" element={<Login />} />
            <Route path="/home" element={<Home />} />
            <Route path="/user" element={<User />} />
            <Route path="/UserPage/Edit/:userId" element={<UserEdit />} />
            <Route path="/register" element={<Register />} />
            {/* Rota para página 404 */}
            <Route path="*" element={<NotFoundPage />} />
          </Routes>
        </main>
        <Footer />
      </div>
    </Router>
  );
}

export default App;

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
import UserIncomeF from "./pages/UserPage/IncomeForecast";
import UserIncomeA from "./pages/UserPage/IncomeActuals";
import UserExpenseF from "./pages/UserPage/ExpenseForecast";
import UserExpenseA from "./pages/UserPage/ExpenseActuals";
import UserAssetF from "./pages/UserPage/AssetForecast";
import UserAssetA from "./pages/UserPage/AssetActuals";

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
            <Route path="/user-page/edit/:userId" element={<UserEdit />} />
            <Route path="/register" element={<Register />} />
            <Route path="/user-income-forecast" element={<UserIncomeF />} />
            <Route path="/user-income-actuals" element={<UserIncomeA />} />
            <Route path="/user-expense-forecast" element={<UserExpenseF />} />
            <Route path="/user-expense-actuals" element={<UserExpenseA />} />
            <Route path="/user-asset-forecast" element={<UserAssetF />} />
            <Route path="/user-asset-actuals" element={<UserAssetA />} />
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

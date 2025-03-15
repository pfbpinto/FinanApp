import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../components/AuthContext";
import toastr from "toastr";
import "toastr/build/toastr.min.css";

const Login = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [passwordVisible, setPasswordVisible] = useState(false);
  const { setIsLoggedIn, setUser } = useAuth();
  const navigate = useNavigate();

  const handleLogin = async (e) => {
    e.preventDefault();

    const loginData = {
      email: email,
      password: password,
    };

    try {
      const response = await fetch("/api/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(loginData),
        credentials: "include",
      });

      if (!response.ok) {
        const errorMessage = await response.text();
        toastr.error(errorMessage);
        return;
      }

      const data = await response.json();
      if (data.status === "success") {
        setIsLoggedIn(true);
        setUser(data.user);
        navigate("/user");
      }
    } catch (error) {
      console.error("Login failed:", error);
      toastr.error("An error occurred. Please try again.");
    }
  };

  const togglePasswordVisibility = () => {
    setPasswordVisible(!passwordVisible);
  };

  return (
    <div className="max-w-md mx-auto p-6 mt-10 bg-white rounded-lg shadow-lg">
      <h1 className="text-2xl font-bold text-center text-blue-600 mb-6">
        Login
      </h1>
      <form onSubmit={handleLogin}>
        {/* Email Input */}
        <div className="mb-4">
          <label
            htmlFor="email"
            className="block text-sm font-medium text-gray-700"
          >
            Email
          </label>
          <input
            type="email"
            id="email"
            name="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="mt-1 block w-full px-4 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
            required
          />
        </div>

        {/* Password Input with Toggle */}
        <div className="mb-6 relative">
          <label
            htmlFor="password"
            className="block text-sm font-medium text-gray-700"
          >
            Password
          </label>
          <div className="relative">
            <input
              type={passwordVisible ? "text" : "password"}
              id="password"
              name="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="mt-1 block w-full px-4 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 pr-12"
              required
            />
            <button
              type="button"
              onClick={togglePasswordVisibility}
              className="absolute inset-y-0 right-2 flex items-center text-gray-500 hover:text-gray-700 focus:outline-none"
            >
              <i
                className={passwordVisible ? "fa fa-eye-slash" : "fa fa-eye"}
              />
            </button>
          </div>
        </div>

        {/* Submit Button */}
        <button
          type="submit"
          className="w-full py-2 px-4 bg-blue-600 text-white font-semibold rounded-md hover:bg-blue-700"
        >
          Login
        </button>
      </form>

      {/* Links for Register and Forgot Password */}
      <div className="mt-4 text-center">
        <a href="/register" className="text-blue-600 hover:underline">
          Register
        </a>
        <span className="text-gray-400 mx-2">|</span>
        <a href="/forgot-password" className="text-blue-600 hover:underline">
          Forgot Password
        </a>
      </div>
    </div>
  );
};

export default Login;

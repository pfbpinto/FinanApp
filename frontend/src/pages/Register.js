import React, { useState } from "react";
import toastr from "toastr";
import "toastr/build/toastr.min.css";

const Register = () => {
  const [passwordVisible, setPasswordVisible] = useState(false);
  const [formData, setFormData] = useState({
    first_name: "",
    last_name: "",
    dob: "",
    email: "",
    password: "",
  });

  const togglePasswordVisibility = () => {
    setPasswordVisible(!passwordVisible);
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const response = await fetch("/api/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(formData),
      });

      if (response.ok) {
        // Success
        const data = await response.json();
        toastr.success(data.message); // Exibe a mensagem de sucesso com toastr
        // Redireciona para a página de login
        window.location.href = "/login"; // Redireciona para a página de login
      } else {
        // Error
        const errorData = await response.json();
        toastr.error(errorData.message); // Exibe os erros com toastr
      }
    } catch (error) {
      toastr.error("Error during registration: " + error.message); // Exibe o erro com toastr
    }
  };

  return (
    <div className="max-w-md mx-auto p-6 mt-10 bg-white rounded-lg shadow-lg">
      <h1 className="text-2xl font-bold text-center text-blue-600 mb-6">
        Register
      </h1>
      <form onSubmit={handleSubmit}>
        {/* FirstName Input */}
        <div className="mb-4">
          <label
            htmlFor="first_name"
            className="block text-sm font-medium text-gray-700"
          >
            First Name
          </label>
          <input
            type="text"
            id="first_name"
            name="first_name"
            className="mt-1 block w-full px-4 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
            value={formData.first_name}
            onChange={handleInputChange}
            required
          />
        </div>

        {/* LastName Input */}
        <div className="mb-4">
          <label
            htmlFor="last_name"
            className="block text-sm font-medium text-gray-700"
          >
            Last Name
          </label>
          <input
            type="text"
            id="last_name"
            name="last_name"
            className="mt-1 block w-full px-4 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
            value={formData.last_name}
            onChange={handleInputChange}
            required
          />
        </div>

        {/* DataOfBirth Input */}
        <div className="mb-4">
          <label
            htmlFor="dob"
            className="block text-sm font-medium text-gray-700"
          >
            Date of Birth
          </label>
          <input
            type="date"
            id="dob"
            name="dob"
            className="mt-1 block w-full px-4 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
            value={formData.dob}
            onChange={handleInputChange}
            required
          />
        </div>

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
            className="mt-1 block w-full px-4 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
            value={formData.email}
            onChange={handleInputChange}
            required
          />
        </div>

        {/* Password Input */}
        <div className="mb-6">
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
              className="mt-1 block w-full px-4 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 pr-12"
              value={formData.password}
              onChange={handleInputChange}
              required
            />
            <button
              type="button"
              onClick={togglePasswordVisibility}
              className="absolute inset-y-0 right-2 flex items-center text-gray-500 hover:text-gray-700 focus:outline-none"
            >
              <i
                className={`fa ${passwordVisible ? "fa-eye-slash" : "fa-eye"}`}
              ></i>
            </button>
          </div>
        </div>

        {/* Submit Button */}
        <button
          type="submit"
          className="w-full py-2 px-4 bg-blue-600 text-white font-semibold rounded-md hover:bg-blue-700"
        >
          Register
        </button>
      </form>

      {/* Links for Login */}
      <div className="mt-4 text-center">
        <a href="/login" className="text-blue-600 hover:underline">
          Already have an account? Login
        </a>
      </div>
    </div>
  );
};

export default Register;

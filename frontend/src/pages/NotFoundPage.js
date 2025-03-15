// src/pages/NotFoundPage.js
import React from "react";
import { Link } from "react-router-dom";

const NotFoundPage = () => {
  return (
    <div className="flex items-center justify-center min-h-screen bg-white-100">
      <div className="text-center p-6 bg-white shadow-lg rounded-lg">
        <h1 className="text-6xl font-bold text-red-500">404</h1>
        <p className="text-xl text-gray-700 mt-4">
          Oops! The page you're looking for doesn't exist.
        </p>
        <p className="text-gray-500 mt-2">
          It might have been moved or deleted.
        </p>
        <Link
          to="/"
          className="mt-6 inline-block px-6 py-2 text-white bg-blue-600 hover:bg-blue-700 rounded-lg"
        >
          Go Back to Home
        </Link>
      </div>
    </div>
  );
};

export default NotFoundPage;

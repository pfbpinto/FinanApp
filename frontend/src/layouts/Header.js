import React, { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../components/AuthContext";
import toastr from "toastr";
import "toastr/build/toastr.min.css";
import userAvatar from "../assets/images/user.svg";

const Header = () => {
  const { isLoggedIn, user, loading, setIsLoggedIn, setUser } = useAuth();
  const navigate = useNavigate();
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);

  const handleLogout = async () => {
    try {
      const response = await fetch("/api/logout", {
        method: "POST",
        credentials: "include",
      });

      if (response.ok) {
        setIsLoggedIn(false);
        setUser(null);
        navigate("/login");
        toastr.success("You are logged out!");
      } else {
        throw new Error("Error trying to logout");
      }
    } catch (error) {
      console.error("Error trying to logout:", error);
      toastr.error("Error trying to logout.");
    }
  };

  const toggleDropdown = () => {
    setIsDropdownOpen(!isDropdownOpen);
  };

  return (
    <header className="bg-gray-800 text-white fixed top-0 left-0 w-full z-10 shadow-md">
      <div className="container mx-auto flex justify-between items-center p-4">
        {/* Logo */}
        <div className="text-xl font-bold">
          <Link to="/">FinanAPP</Link>
        </div>

        {/* Navegação */}
        <nav className="flex space-x-4 ml-auto">
          <Link to="/" className="hover:text-blue-500 mt-1">
            Home
          </Link>

          {!loading && isLoggedIn ? (
            <>
              <Link to="/user" className="hover:text-blue-500 mt-1">
                Dashboard
              </Link>

              {/* Dropdown do usuário */}
              <div className="relative">
                <button
                  onClick={toggleDropdown}
                  className="hover:text-blue-500 focus:outline-none flex items-center"
                >
                  {/* Foto do usuário ou avatar padrão */}
                  <img
                    src={user?.profilePicture || userAvatar}
                    alt="Perfil do Usuário"
                    className="w-8 h-8 rounded-full mr-2"
                  />
                  {/* Ícone do dropdown */}
                  <i className="fa fa-chevron-down ml-2"></i>
                </button>

                {isDropdownOpen && (
                  <div className="absolute right-0 mt-2 w-48 bg-white text-black rounded-md shadow-lg">
                    <div className="px-4 py-2 text-sm">
                      Olá, {user?.firstName || "Usuário"}
                    </div>
                    <button
                      onClick={handleLogout}
                      className="block px-4 py-2 text-sm w-full text-left hover:bg-gray-100"
                    >
                      Logout
                    </button>
                  </div>
                )}
              </div>
            </>
          ) : (
            <Link to="/login" className="hover:text-blue-500 mt-1">
              Login
            </Link>
          )}
        </nav>
      </div>
    </header>
  );
};

export default Header;

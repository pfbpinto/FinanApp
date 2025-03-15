import React, { createContext, useContext, useState, useEffect } from "react";
import axios from "axios";

const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  const [isLoggedIn, setIsLoggedIn] = useState(null); // Authentication status
  const [user, setUser] = useState(null); // User full data
  const [loading, setLoading] = useState(true); // Loading indicator

  useEffect(() => {
    const fetchAuthStatus = async () => {
      try {
        // Making the API call to check authentication status
        const response = await axios.get("/api/auth-status", {
          withCredentials: true, // Ensures cookies are sent with the request
        });

        const data = response.data;

        // Updates the state based on the response
        setIsLoggedIn(data.authenticated); // `authenticated` from the backend
        setUser(data); // Stores all user data
      } catch (error) {
        // Checks if the error is an unauthorized (401) response
        if (error.response && error.response.status === 401) {
          setIsLoggedIn(false); // Sets the user as not logged in
          setUser(null); // Clears user data
        } else {
          // Handles other errors and logs them
          console.log("Error checking auth status", error);
          setIsLoggedIn(false); // Sets the user as not logged in
          setUser(null); // Clears user data
        }
      } finally {
        setLoading(false); // Marks loading as complete
      }
    };

    // Fetch authentication status on component mount
    fetchAuthStatus();
  }, []); // Empty dependency array means this effect runs once on mount

  return (
    <AuthContext.Provider
      value={{
        isLoggedIn, // Authentication status
        user, // Full user data
        loading, // Loading status
        setIsLoggedIn, // Function to update authentication status
        setUser, // Function to update user data
      }}
    >
      {children} {/* Renders the children components */}
    </AuthContext.Provider>
  );
};

// Custom hook to access the AuthContext
export const useAuth = () => useContext(AuthContext);

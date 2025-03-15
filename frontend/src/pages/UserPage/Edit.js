import React, { useState, useEffect } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { useAuth } from "../../components/AuthContext";
import toastr from "toastr";
import "toastr/build/toastr.min.css";

const EditProfile = () => {
  const { isLoggedIn, user, loading } = useAuth();
  const navigate = useNavigate();
  const { userId } = useParams();

  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [dateOfBirth, setDateOfBirth] = useState("");

  // Verificação inicial - se não houver userId, exibe erro e redireciona
  useEffect(() => {
    if (!userId) {
      toastr.error("User ID is missing!");
      navigate("/user");
    }
  }, [userId, navigate]);

  // Verifica se o usuário está logado e redireciona se necessário
  useEffect(() => {
    if (!loading && !isLoggedIn) {
      navigate("/login");
    }
  }, [isLoggedIn, loading, navigate]);

  const handleSubmit = async (e) => {
    e.preventDefault();

    const userData = {
      firstName,
      lastName,
      dateOfBirth,
      userId: parseInt(userId),
    };

    try {
      const response = await fetch("/api/user-edit", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(userData),
        credentials: "include",
      });

      // Verifique se a resposta não foi bem-sucedida
      if (!response.ok) {
        let errorData = await response
          .json()
          .catch(() => ({ message: ["Failed to update profile"] }));

        const errorMessage =
          errorData.message && Array.isArray(errorData.message)
            ? errorData.message.join(", ")
            : "Failed to update profile";

        toastr.error(errorMessage);
        return;
      }

      const data = await response.json();
      if (data.status === "success") {
        toastr.success("Profile updated successfully!");
        navigate("/user");
      }
    } catch (error) {
      console.error("Profile update failed:", error);
      toastr.error("An error occurred. Please try again.");
    }
  };

  if (loading || !user) {
    return <p>Loading...</p>;
  }

  return (
    <div className="max-w-md mx-auto p-6 mt-10 bg-white rounded-lg shadow-lg">
      <h1 className="text-2xl font-bold text-center text-blue-600 mb-6">
        Edit Profile
      </h1>
      <form onSubmit={handleSubmit}>
        {/* First Name Input */}
        <div className="mb-4">
          <label
            htmlFor="firstName"
            className="block text-sm font-medium text-gray-700"
          >
            First Name
          </label>
          <input
            type="text"
            id="firstName"
            name="firstName"
            placeholder={user.firstName}
            value={firstName}
            onChange={(e) => setFirstName(e.target.value)}
            className="mt-1 block w-full px-4 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
            required
          />
        </div>

        {/* Last Name Input */}
        <div className="mb-4">
          <label
            htmlFor="lastName"
            className="block text-sm font-medium text-gray-700"
          >
            Last Name
          </label>
          <input
            type="text"
            id="lastName"
            name="lastName"
            placeholder={user.lastName}
            value={lastName}
            onChange={(e) => setLastName(e.target.value)}
            className="mt-1 block w-full px-4 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
            required
          />
        </div>

        {/* Date of Birth Input */}
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
            placeholder={user.dataOfBirth}
            value={user.dataOfBirth}
            onChange={(e) => setDateOfBirth(e.target.value)}
            required
          />
        </div>

        {/* Submit Button */}
        <button
          type="submit"
          className="w-full py-2 px-4 bg-blue-600 text-white font-semibold rounded-md hover:bg-blue-700"
        >
          Save Changes
        </button>
      </form>
    </div>
  );
};

export default EditProfile;

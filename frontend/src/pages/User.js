import React, { useEffect, useState } from "react";
import { useCallback } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../components/AuthContext";
import TaxForm from "../components/TaxForm";
import GroupForm from "../components/GroupForm";
import ModalLarge from "../components/ModalLarge";
import { Link } from "react-router-dom";

function User() {
  const [userDashboard, setUserDashboard] = useState(null);
  const [isModalTaxOpen, setIsModalTaxOpen] = useState(false);
  const [isModalGroupOpen, setIsModalGroupOpen] = useState(false);
  const [currentUser, setCurrentUser] = useState(null);
  const { isLoggedIn, loading } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    if (loading) return;
    if (!isLoggedIn) navigate("/login");
  }, [isLoggedIn, loading, navigate]);

  useEffect(() => {
    if (isLoggedIn && !userDashboard) {
      fetchUserDashboard();
    }
  }, [isLoggedIn, userDashboard]);

  // Fetch User Dashboard info
  const fetchUserDashboard = () => {
    fetch("/api/user", {
      method: "GET",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
    })
      .then((response) => response.json())
      .then((data) => {
        setUserDashboard(data); // Update user data
        setCurrentUser(data.user.userProfileID);
      })
      .catch((error) => console.error("Error fetching user data:", error));
  };

  // Open the Tax Modal
  const openTaxModal = useCallback(() => {
    setIsModalTaxOpen(true); // Open the tax modal
  }, []);

  const openGroupModal = useCallback(() => {
    setIsModalGroupOpen(true); // Open Group modal
  }, []);

  // Close the Tax Modal (reset state)
  const closeTaxModal = useCallback(() => {
    setIsModalTaxOpen(false); // Close the modal
  }, []);

  // Close the Group Modal (reset state)
  const closeGroupModal = useCallback(() => {
    setIsModalGroupOpen(false); // Close the modal
  }, []);

  if (loading) {
    return <p className="text-gray-600 text-center mt-6">Loading...</p>;
  }

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-3xl font-semibold text-gray-800">Your Dashboard</h1>

      {userDashboard ? (
        <div className="bg-gray-100 p-6 rounded-lg shadow-md mt-4 flex flex-col md:flex-row gap-6">
          {/* Card do Usuário */}
          <div className="w-full md:w-1/2 bg-white p-6 rounded-lg shadow-md">
            <div className="flex items-center justify-between">
              <div className="flex items-center">
                <div className="w-20 h-20 bg-gray-300 rounded-full"></div>
                <div className="ml-4">
                  <h2 className="text-xl font-medium text-gray-900">
                    {userDashboard.user.firstName} {userDashboard.user.lastName}
                  </h2>
                </div>
              </div>
              <Link
                to={`/user-page/edit/${userDashboard.user.userProfileID}`}
                className="px-4 py-2 text-sm font-medium text-white bg-blue-500 rounded-md hover:bg-blue-600 transition"
              >
                Edit Profile
              </Link>
            </div>
            <div className="space-y-4 mt-4">
              <InfoRow
                label="Name"
                value={`${userDashboard.user.firstName} ${userDashboard.user.lastName}`}
              />
              <InfoRow label="Email" value={userDashboard.user.emailAddress} />

              <InfoRow
                label="Date of Birth"
                value={formatDate(userDashboard.user.dob)}
              />
              <InfoRow
                label="Account Created"
                value={formatDate(userDashboard.user.createdAt)}
              />
              <InfoRow
                label="User Subscription"
                value={
                  userDashboard.user.userSubscription ? "Active" : "Inactive"
                }
              />
            </div>
          </div>

          {/* Outra Div */}
          <div className="w-full md:w-1/2 bg-white p-6 rounded-lg shadow-md">
            <h3 className="text-lg font-medium text-gray-900 mb-4">
              Additional Info
            </h3>
            <p className="text-gray-600">
              Aqui você pode adicionar mais informações sobre o usuário,
              estatísticas ou qualquer outra seção relevante.
            </p>
          </div>
        </div>
      ) : (
        <p className="text-gray-600 text-center mt-6">
          Loading user dashboard...
        </p>
      )}

      <div className="bg-gray-100 p-6 rounded-lg shadow-md mt-4">
        <div className="flex flex-wrap gap-2">
          <button
            className="px-4 py-2 mb-3 text-sm font-medium text-white bg-green-500 rounded-md hover:bg-green-600 transition"
            onClick={() => openTaxModal()}
          >
            Setup Taxes
          </button>
          <button
            className="px-4 py-2 mb-3 text-sm font-medium text-white bg-green-500 rounded-md hover:bg-green-600 transition"
            onClick={() => openGroupModal()}
          >
            Setup Groups
          </button>

          <Link
            to={`/user-income`}
            className="px-4 py-2 mb-3 text-sm font-medium text-white bg-blue-500 rounded-md hover:bg-blue-600 transition"
          >
            Incomes
          </Link>
          <Link
            to={`/user-expense`}
            className="px-4 py-2 mb-3 text-sm font-medium text-white bg-blue-500 rounded-md hover:bg-blue-600 transition"
          >
            Expenses
          </Link>
          <Link
            to={`/user-asset`}
            className="px-4 py-2 mb-3 text-sm font-medium text-white bg-blue-500 rounded-md hover:bg-blue-600 transition"
          >
            Assets
          </Link>
        </div>
      </div>

      {isModalTaxOpen && (
        <ModalLarge onClose={closeTaxModal} title={`Setup Taxes`}>
          <TaxForm onClose={closeTaxModal} user={currentUser} />
        </ModalLarge>
      )}

      {isModalGroupOpen && (
        <ModalLarge onClose={closeGroupModal} title={`Setup Groups`}>
          <GroupForm onClose={closeGroupModal} user={currentUser} />
        </ModalLarge>
      )}
    </div>
  );
}

function InfoRow({ label, value, className = "" }) {
  return (
    <div className="flex justify-between items-center">
      <span className="text-gray-700 font-medium">{label}:</span>
      <span className={`text-gray-800 font-semibold ${className}`}>
        {value}
      </span>
    </div>
  );
}

function formatDate(dateString) {
  return new Date(dateString).toLocaleDateString();
}

export default User;

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
                    {userDashboard.user.first_name}{" "}
                    {userDashboard.user.last_name}
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
                value={`${userDashboard.user.first_name} ${userDashboard.user.last_name}`}
              />
              <InfoRow label="Email" value={userDashboard.user.email_address} />

              <InfoRow
                label="Date of Birth"
                value={formatDate(userDashboard.user.date_of_birth)}
              />
              <InfoRow
                label="Account Created"
                value={formatDate(userDashboard.user.created_at)}
              />
              <InfoRow
                label="User Subscription"
                value={
                  userDashboard.user.user_subscription ? "Active" : "Inactive"
                }
              />
            </div>
          </div>

          {/* Side div */}
          <div className="w-full md:w-1/2 bg-white p-6 rounded-lg shadow-md">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="relative w-full max-w-md bg-white p-6 rounded-lg shadow-md border-2 border-gray-900 text-center mx-auto">
                <span className="absolute -top-4 left-1/2 transform -translate-x-1/2 bg-white px-4 text-gray-900 font-bold text-lg uppercase tracking-wide">
                  Incomes
                </span>
                <div className="flex flex-col items-center justify-center gap-3 mt-6">
                  <div className="flex flex-wrap gap-2 w-full justify-center">
                    <Link
                      to="/user-income"
                      state={{ userID: currentUser }}
                      className="flex-1 sm:flex-none px-4 py-2 text-sm font-medium text-white bg-blue-500 rounded-md hover:bg-blue-600 transition text-center"
                    >
                      Manage
                    </Link>
                  </div>
                </div>
              </div>

              <div className="relative w-full max-w-md bg-white p-6 rounded-lg shadow-md border-2 border-gray-900 text-center mx-auto">
                <span className="absolute -top-4 left-1/2 transform -translate-x-1/2 bg-white px-4 text-gray-900 font-bold text-lg uppercase tracking-wide">
                  Assets
                </span>
                <div className="flex flex-col items-center justify-center gap-3 mt-6">
                  <div className="flex flex-wrap gap-2 w-full justify-center">
                    <Link
                      to="/user-asset-forecast"
                      className="flex-1 sm:flex-none px-4 py-2 text-sm font-medium text-white bg-blue-500 rounded-md hover:bg-blue-600 transition text-center"
                    >
                      Forecast
                    </Link>
                    <Link
                      to="/user-asset-actuals"
                      className="flex-1 sm:flex-none px-4 py-2 text-sm font-medium text-white bg-blue-500 rounded-md hover:bg-blue-600 transition text-center"
                    >
                      Actuals
                    </Link>
                  </div>
                </div>
              </div>

              <div className="relative w-full max-w-md bg-white p-6 rounded-lg shadow-md border-2 border-gray-900 text-center mx-auto">
                <span className="absolute -top-4 left-1/2 transform -translate-x-1/2 bg-white px-4 text-gray-900 font-bold text-lg uppercase tracking-wide">
                  Expenses
                </span>
                <div className="flex flex-col items-center justify-center gap-3 mt-6">
                  <div className="flex flex-wrap gap-2 w-full justify-center">
                    <Link
                      to="/user-expense-forecast"
                      className="flex-1 sm:flex-none px-4 py-2 text-sm font-medium text-white bg-blue-500 rounded-md hover:bg-blue-600 transition text-center"
                    >
                      Forecast
                    </Link>
                    <Link
                      to="/user-expense-actuals"
                      className="flex-1 sm:flex-none px-4 py-2 text-sm font-medium text-white bg-blue-500 rounded-md hover:bg-blue-600 transition text-center"
                    >
                      Actuals
                    </Link>
                  </div>
                </div>
              </div>

              <div className="relative w-full max-w-md bg-white p-6 rounded-lg shadow-md border-2 border-gray-900 text-center mx-auto">
                <span className="absolute -top-4 left-1/2 transform -translate-x-1/2 bg-white px-4 text-gray-900 font-bold text-lg uppercase tracking-wide">
                  Groups
                </span>
                <div className="flex flex-col items-center justify-center gap-3 mt-6">
                  <div className="flex flex-wrap gap-2 w-full justify-center">
                    <button
                      onClick={() => openGroupModal()}
                      className="flex-1 sm:flex-none px-4 py-2 text-sm font-medium text-white bg-blue-500 rounded-md hover:bg-blue-600 transition text-center"
                    >
                      Manage Groups
                    </button>
                  </div>
                </div>
              </div>
            </div>
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

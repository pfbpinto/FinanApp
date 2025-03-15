import React, { useEffect, useState } from "react";
import { useCallback } from "react";

import { useNavigate } from "react-router-dom";
import { useAuth } from "../components/AuthContext";
import AssetForm from "../components/AssetForm";
import TaxForm from "../components/TaxForm";
import CategoryForm from "../components/CategoryForm";
import GroupForm from "../components/GroupForm";
import IncomeForm from "../components/IncomeForm";
import ExpenseForm from "../components/ExpenseForm";
import Modal from "../components/Modal";
import ModalLarge from "../components/ModalLarge";
import toastr from "toastr";
import { Link } from "react-router-dom";
import { ChevronUpIcon } from "@heroicons/react/24/solid";

function User() {
  const [userDashboard, setUserDashboard] = useState(null);
  const [isModalAssetOpen, setIsModalAssetOpen] = useState(false);
  const [isModalCategoryOpen, setisModalCategoryOpen] = useState(false);
  const [isModalIncomeOpen, setIsModalIncomeOpen] = useState(false);
  const [isModalExpenditureOpen, setisModalExpenditureOpen] = useState(false);
  const [isModalTaxOpen, setIsModalTaxOpen] = useState(false);
  const [isModalGroupOpen, setIsModalGroupOpen] = useState(false);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [deleteType, setDeleteType] = useState("");
  const [itemToDelete, setItemToDelete] = useState(null);
  const [currentUser, setCurrentUser] = useState(null);
  const [currentAsset, setCurrentAsset] = useState(null);
  const [currentCategory, setCurrentCategory] = useState(null);
  const [currentIncome, setcurrentIncome] = useState(null);
  const [currentExpenditure, setcurrentExpenditure] = useState(null);
  const [openToggleAsset, setopenToggleAsset] = useState(true);
  const [openToggleIncome, setopenToggleIncome] = useState(true);
  const [openToggleExpenditure, setopenToggleExpenditure] = useState(true);
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
        setCurrentUser(data.user.ID);
      })
      .catch((error) => console.error("Error fetching user data:", error));
  };
  // Toggle the asset panel (open/close)
  const toggleAssetPanel = useCallback(() => {
    setopenToggleAsset((prevState) => !prevState); // Toggle state
  }, []);

  // Toggle the Income panel (open/close)
  const toggleIncomePanel = useCallback(() => {
    setopenToggleIncome((prevState) => !prevState); // Toggle state
  }, []);

  // Toggle the Expenditure panel (open/close)
  const toggleExpenditurePanel = useCallback(() => {
    setopenToggleExpenditure((prevState) => !prevState); // Toggle state
  }, []);

  // Open the Asset Create Modal (reset the asset to null)
  const openCreateAssetModal = useCallback(() => {
    setCurrentAsset(null); // Ensure the current Asset is empty
    setIsModalAssetOpen(true); // Open the modal
  }, []);

  // Open the Asset Update Modal (set the current asset to be edited)
  const openUpdateAssetModal = useCallback((asset) => {
    setCurrentAsset(asset); // Set the Asset to be updated
    setIsModalAssetOpen(true); // Open the modal
  }, []);

  // Open the Tax Modal
  const openTaxModal = useCallback(() => {
    setIsModalTaxOpen(true); // Open the tax modal
  }, []);

  const openGroupModal = useCallback(() => {
    setIsModalGroupOpen(true); // Open Group modal
  }, []);

  // Open the Category Modal (pass the category type to set)
  const openCategoryModal = useCallback((type) => {
    setCurrentCategory(type); // Set the category type
    setisModalCategoryOpen(true); // Open the modal
  }, []);

  // Open the Income Create Modal (reset the Income to null)
  const openCreateIncomeModal = useCallback(() => {
    setcurrentIncome(null); // Ensure the current Income is empty
    setIsModalIncomeOpen(true); // Open the income create modal
  }, []);

  // Open the Income Update Modal (set the current Income to be edited)
  const openUpdateIncomeModal = useCallback((income) => {
    setcurrentIncome(income); // Set the Income to be updated
    setIsModalIncomeOpen(true); // Open the income edit modal
  }, []);

  // Open the Expenditure Create Modal (reset the Expenditure to null)
  const openCreateExpenditureModal = useCallback(() => {
    setcurrentExpenditure(null); // Ensure the current Expense is empty
    setisModalExpenditureOpen(true); // Open the expense modal
  }, []);

  // Open the Expenditure Create Modal (reset the Expenditure to null)
  const openUpdateExpenditureModal = useCallback((expense) => {
    setcurrentExpenditure(expense); // Set the Expense to be updated
    setisModalExpenditureOpen(true); // Open the modal
  }, []);

  // Open the Delete Modal (set the item and type to delete)
  const openDeleteModal = useCallback((item, type) => {
    if (item && item.ID) {
      setItemToDelete(item); // Set the item to delete
      setDeleteType(type); // Set the type of item to delete
      setIsDeleteModalOpen(true); // Open the delete modal
    } else {
      console.error("Invalid item or missing ID"); // Log error if item or ID is missing
    }
  }, []);

  // Close the Delete Modal (reset state)
  const closeDeleteModal = useCallback(() => {
    setIsDeleteModalOpen(false); // Close the modal
    setItemToDelete(null); // Reset the item to delete
  }, []);

  // Close the Asset Modal (reset state)
  const closeAssetModal = useCallback(() => {
    setIsModalAssetOpen(false); // Close the modal
    setCurrentAsset(null); // Reset the asset
  }, []);

  // Close the Tax Modal (reset state)
  const closeTaxModal = useCallback(() => {
    setIsModalTaxOpen(false); // Close the modal
  }, []);

  // Close the Group Modal (reset state)
  const closeGroupModal = useCallback(() => {
    setIsModalGroupOpen(false); // Close the modal
  }, []);

  // Close the Category Modal (reset state)
  const closeCategoryModal = useCallback(() => {
    setisModalCategoryOpen(false); // Close the modal
    setCurrentCategory(null); // Reset the category type
  }, []);

  // Close the Income Modal (reset state)
  const closeIncomeModal = useCallback(() => {
    setIsModalIncomeOpen(false); // Close the modal
    setcurrentIncome(null); // Reset the category type
  }, []);

  // Close the Expenditure Modal (reset state)
  const closeExpenditureModal = useCallback(() => {
    setisModalExpenditureOpen(false); // Close the modal
    setcurrentExpenditure(null); // Reset the category type
  }, []);

  // Asset Create and Update handle
  const handleAssetSubmit = (data) => {
    // Determine the HTTP method (POST for new asset, PUT for updating an existing one)
    const method = currentAsset ? "PUT" : "POST";

    // Set the URL for the API request (use asset ID for updating)
    const url = currentAsset ? `/api/assets/${currentAsset.ID}` : "/api/assets";

    // If there are no taxes, the taxes array will be empty, but we ensure that 'Taxes' is present
    const dataToSend = currentAsset
      ? { ...data, AssetValue: parseFloat(data.AssetValue) }
      : { ...data, userID: userDashboard.user.ID };

    // Send the data to the server
    fetch(url, {
      method: method,
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(dataToSend),
      credentials: "include", // Include credentials like cookies for authentication
    })
      .then((response) => {
        // Check if the response status is ok (status code 200-299)
        if (!response.ok) {
          // If response is not ok, extract error message from the JSON response
          return response.json().then((errorData) => {
            return Promise.reject(errorData);
          });
        }
        // If successful, return the response as JSON
        return response.json();
      })
      .then((data) => {
        if (!data) return;
        // Show a success message based on whether it's a new asset or an update
        toastr.success(currentAsset ? "Asset updated!" : "Asset created!");
        setIsModalAssetOpen(false); // Close the modal
        fetchUserDashboard(); // Refresh the user dashboard data
      })
      .catch((error) => {
        // Log any error during the save process
        console.error("Error saving asset:", error);
        toastr.error(`Error: ${error.message}`);
      });
  };

  // Income Create and Update handle
  const handleIncomeSubmit = (data) => {
    // Determine the HTTP method (POST for new asset, PUT for updating an existing one)
    const method = currentIncome ? "PUT" : "POST";

    // Set the URL for the API request (use asset ID for updating)
    const url = currentIncome
      ? `/api/income/${currentIncome.ID}`
      : "/api/income";

    // If there are no taxes, the taxes array will be empty, but we ensure that 'Taxes' is present
    const dataToSend = currentIncome
      ? { ...data, IncomeValue: parseFloat(data.IncomeValue) }
      : { ...data, userID: userDashboard.user.ID };

    // Send the data to the server
    fetch(url, {
      method: method,
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(dataToSend),
      credentials: "include", // Include credentials like cookies for authentication
    })
      .then((response) => {
        // Check if the response status is ok (status code 200-299)
        if (!response.ok) {
          // If response is not ok, extract error message from the JSON response
          return response.json().then((errorData) => {
            return Promise.reject(errorData);
          });
        }
        // If successful, return the response as JSON
        return response.json();
      })
      .then(() => {
        // Show a success message based on whether it's a new asset or an update
        toastr.success(currentIncome ? "Income updated!" : "Income created!");
        setIsModalIncomeOpen(false); // Close the modal
        fetchUserDashboard(); // Refresh the user dashboard data
      })
      .catch((error) => {
        // Log any error during the save process
        console.error("Error saving income:", error);
        toastr.error(`Error: ${error.message}`);
      });
  };

  // Expense Create and Update handle
  const handleExpenseSubmit = (data) => {
    // Determine the HTTP method (POST for new asset, PUT for updating an existing one)
    const method = currentExpenditure ? "PUT" : "POST";

    // Set the URL for the API request (use asset ID for updating)
    const url = currentExpenditure
      ? `/api/expense/${currentExpenditure.ID}`
      : "/api/expense";

    // If there are no taxes, the taxes array will be empty, but we ensure that 'Taxes' is present
    const dataToSend = currentExpenditure
      ? { ...data, ExpenditureValue: parseFloat(data.ExpenditureValue) }
      : { ...data, userID: userDashboard.user.ID };

    // Send the data to the server
    fetch(url, {
      method: method,
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(dataToSend),
      credentials: "include", // Include credentials like cookies for authentication
    })
      .then((response) => {
        // Check if the response status is ok (status code 200-299)
        if (!response.ok) {
          // If response is not ok, extract error message from the JSON response
          return response.json().then((errorData) => {
            return Promise.reject(errorData);
          });
        }
        // If successful, return the response as JSON
        return response.json();
      })
      .then(() => {
        // Show a success message based on whether it's a new asset or an update
        toastr.success(
          currentExpenditure ? "Expense updated!" : "Expense created!"
        );
        setisModalExpenditureOpen(false); // Close the modal
        fetchUserDashboard(); // Refresh the user dashboard data
      })
      .catch((error) => {
        // Log any error during the save process
        console.error("Error saving Expense:", error);
        toastr.error(`Error: ${error.message}`);
      });
  };

  // Item Delete handle
  const handleDeleteItem = (item) => {
    if (!item || !item.ID) return;
    fetch(`/api/delete-${deleteType}/${item.ID}`, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ itemId: item.ID }),
      credentials: "include",
    })
      .then((response) => response.json())
      .then(() => {
        toastr.success(
          `${
            deleteType.charAt(0).toUpperCase() + deleteType.slice(1)
          } Successfully Deleted!`
        );
        fetchUserDashboard();
        closeDeleteModal();
        setDeleteType(null);
      })
      .catch((error) => {
        console.error("Error deleting item:", error);
        toastr.error(`Erro ao deletar o ${deleteType}.`);
        closeDeleteModal();
        setDeleteType(null);
      });
  };

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
                    {userDashboard.user.FirstName} {userDashboard.user.LastName}
                  </h2>
                  <p className="text-gray-500">
                    {userDashboard.user.UserType
                      ? userDashboard.user.UserType.Name
                      : "N/A"}
                  </p>
                </div>
              </div>
              <Link
                to={`/UserPage/Edit/${userDashboard.user.ID}`}
                className="px-4 py-2 text-sm font-medium text-white bg-blue-500 rounded-md hover:bg-blue-600 transition"
              >
                Edit Profile
              </Link>
            </div>
            <div className="space-y-4 mt-4">
              <InfoRow
                label="Name"
                value={`${userDashboard.user.FirstName} ${userDashboard.user.LastName}`}
              />
              <InfoRow label="Email" value={userDashboard.user.EmailAddress} />
              <InfoRow
                label="Role"
                value={
                  userDashboard.user.UserType
                    ? userDashboard.user.UserType.Name
                    : "N/A"
                }
              />
              <InfoRow
                label="Date of Birth"
                value={formatDate(userDashboard.user.DataOfBirth)}
              />
              <InfoRow
                label="Account Created"
                value={formatDate(userDashboard.user.CreatedAt)}
              />
              <InfoRow
                label="Last Login"
                value={formatDate(userDashboard.user.LastLogin)}
              />
              <InfoRow
                label="Status"
                value={userDashboard.user.IsActive ? "Active" : "Inactive"}
                className={
                  userDashboard.user.IsActive
                    ? "text-green-600"
                    : "text-red-600"
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
        </div>
      </div>

      <div className="space-y-6">
        <div className="bg-white shadow-md rounded-lg p-4 mt-2">
          <div
            className="flex justify-between w-full px-4 py-2 text-left text-sm font-medium text-white bg-blue-500 rounded-md hover:bg-blue-600 transition cursor-pointer"
            onClick={toggleAssetPanel}
          >
            <span>My Assets</span>
            <ChevronUpIcon
              className={`w-5 h-5 transition ${
                openToggleAsset ? "rotate-180" : ""
              }`}
            />
          </div>

          {openToggleAsset && (
            <div className="mt-4">
              <button
                className="px-4 py-2 mb-3 text-sm font-medium text-white bg-green-500 rounded-md hover:bg-green-600 transition"
                onClick={() => openCreateAssetModal()}
              >
                New Asset
              </button>

              <button
                className="px-4 py-2 ml-2 mb-3 text-sm font-medium text-white bg-green-500 rounded-md hover:bg-green-600 transition"
                onClick={() => openCategoryModal("asset")}
              >
                New Asset Category
              </button>

              {userDashboard?.userAsset?.length > 0 ? (
                <div className="overflow-x-auto">
                  {" "}
                  <table className="w-full bg-white rounded-lg shadow-md border">
                    <thead>
                      <tr className="bg-gray-100">
                        <th className="px-4 py-2 border-b">Name</th>
                        <th className="px-4 py-2 border-b">Type</th>
                        <th className="px-4 py-2 border-b">Value</th>
                        <th className="px-4 py-2 border-b">Acquisition Date</th>
                        <th className="px-4 py-2 border-b">Shared</th>
                        <th className="px-4 py-2 border-b">Tax</th>
                        <th className="px-4 py-2 border-b"></th>
                      </tr>
                    </thead>
                    <tbody>
                      {userDashboard.userAsset.map((asset, index) => (
                        <tr key={index} className="text-center">
                          <td className="px-4 py-2 border-b">
                            {asset.AssetName || "N/A"}
                          </td>
                          <td className="px-4 py-2 border-b">
                            {asset.AssetType.AssetTypeName || "N/A"}
                          </td>
                          <td className="px-4 py-2 border-b">
                            {asset.AssetValue || "N/A"}
                          </td>
                          <td className="px-4 py-2 border-b">
                            {asset.AssetAquisitionDate
                              ? new Date(
                                  asset.AssetAquisitionDate
                                ).toLocaleDateString()
                              : "N/A"}
                          </td>
                          <td className="px-4 py-2 border-b">
                            {asset.SharedAsset ? "Yes" : "No"}
                          </td>
                          <td className="px-4 py-2 border-b">
                            {asset.UserAssetTaxes &&
                            asset.UserAssetTaxes.length > 0 ? (
                              <div className="flex flex-wrap gap-2">
                                {asset.UserAssetTaxes.map(
                                  (userAssetTax, idx) => (
                                    <div
                                      key={idx}
                                      className="p-1 bg-blue-100 rounded-lg shadow-sm text-sm font-medium text-blue-800"
                                    >
                                      {userAssetTax.Tax?.TaxName || "Tax"}
                                    </div>
                                  )
                                )}
                              </div>
                            ) : (
                              <span className="text-gray-500">No Taxes</span>
                            )}
                          </td>

                          <td className="px-4 py-2 border-b space-x-2">
                            <button
                              className="px-3 py-1 text-xs font-medium text-white bg-yellow-500 rounded-md hover:bg-yellow-600"
                              onClick={() => openUpdateAssetModal(asset)}
                            >
                              Edit
                            </button>
                            <button
                              className="px-3 py-1 text-xs font-medium text-white bg-red-500 rounded-md hover:bg-red-600"
                              onClick={() => openDeleteModal(asset, "assets")}
                            >
                              Delete
                            </button>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              ) : (
                <div className="w-full bg-gray-100 rounded-lg p-4 text-center text-gray-500">
                  You don't have Assets yet
                </div>
              )}
            </div>
          )}
        </div>
      </div>

      <div className="space-y-6">
        <div className="bg-white shadow-md rounded-lg p-4 mt-2">
          <div
            className="flex justify-between w-full px-4 py-2 text-left text-sm font-medium text-white bg-blue-500 rounded-md hover:bg-blue-600 transition cursor-pointer"
            onClick={toggleIncomePanel}
          >
            <span>My Incomes</span>
            <ChevronUpIcon
              className={`w-5 h-5 transition ${
                openToggleIncome ? "rotate-180" : ""
              }`}
            />
          </div>

          {openToggleIncome && (
            <div className="mt-4">
              <button
                className="px-4 py-2 mb-3 text-sm font-medium text-white bg-green-500 rounded-md hover:bg-green-600 transition"
                onClick={() => openCreateIncomeModal()}
              >
                New Income
              </button>

              <button
                className="px-4 py-2 ml-2 mb-3 text-sm font-medium text-white bg-green-500 rounded-md hover:bg-green-600 transition"
                onClick={() => openCategoryModal("income")}
              >
                New Income Category
              </button>

              {userDashboard?.userIncome?.length > 0 ? (
                <div className="overflow-x-auto">
                  {" "}
                  <table className="w-full bg-white rounded-lg shadow-md border">
                    <thead>
                      <tr className="bg-gray-100">
                        <th className="px-4 py-2 border-b">Name</th>
                        <th className="px-4 py-2 border-b">Type</th>
                        <th className="px-4 py-2 border-b">Value</th>
                        <th className="px-4 py-2 border-b">Recurrence</th>
                        <th className="px-4 py-2 border-b">StartDate</th>
                        <th className="px-4 py-2 border-b">Shared</th>
                        <th className="px-4 py-2 border-b">Tax</th>
                        <th className="px-4 py-2 border-b"></th>
                      </tr>
                    </thead>
                    <tbody>
                      {userDashboard.userIncome.map((income, index) => (
                        <tr key={index} className="text-center">
                          <td className="px-4 py-2 border-b">
                            {income.IncomeName || "N/A"}
                          </td>
                          <td className="px-4 py-2 border-b">
                            {income.IncomeType.IncomeTypeName || "N/A"}
                          </td>
                          <td className="px-4 py-2 border-b">
                            {income.IncomeValue || "N/A"}
                          </td>
                          <td className="px-4 py-2 border-b">
                            {income.IncomeRecurrence || "N/A"}
                          </td>
                          <td className="px-4 py-2 border-b">
                            {income.IncomeStartDate
                              ? new Date(
                                  income.IncomeStartDate
                                ).toLocaleDateString()
                              : "N/A"}
                          </td>
                          <td className="px-4 py-2 border-b">
                            {income.SharedIncome ? "Yes" : "No"}
                          </td>
                          <td className="px-4 py-2 border-b">
                            {income.UserTaxes && income.UserTaxes.length > 0 ? (
                              <div className="flex flex-wrap gap-2">
                                {income.UserTaxes.map((UserTax, idx) => (
                                  <div
                                    key={idx}
                                    className="p-1 bg-blue-100 rounded-lg shadow-sm text-sm font-medium text-blue-800"
                                  >
                                    {UserTax.Tax?.TaxName || "Tax Name"}
                                  </div>
                                ))}
                              </div>
                            ) : (
                              <span className="text-gray-500">No Taxes</span>
                            )}
                          </td>

                          <td className="px-4 py-2 border-b space-x-2">
                            <button
                              className="px-3 py-1 text-xs font-medium text-white bg-yellow-500 rounded-md hover:bg-yellow-600"
                              onClick={() => openUpdateIncomeModal(income)}
                            >
                              Edit
                            </button>
                            <button
                              className="px-3 py-1 text-xs font-medium text-white bg-red-500 rounded-md hover:bg-red-600"
                              onClick={() => openDeleteModal(income, "income")}
                            >
                              Delete
                            </button>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              ) : (
                <div className="w-full bg-gray-100 rounded-lg p-4 text-center text-gray-500">
                  You don't have Income yet
                </div>
              )}
            </div>
          )}
        </div>
      </div>

      <div className="space-y-6 mb-10">
        <div className="bg-white shadow-md rounded-lg p-4 mt-2">
          <div
            className="flex justify-between w-full px-4 py-2 text-left text-sm font-medium text-white bg-blue-500 rounded-md hover:bg-blue-600 transition cursor-pointer"
            onClick={toggleExpenditurePanel}
          >
            <span>My Expenditure</span>
            <ChevronUpIcon
              className={`w-5 h-5 transition ${
                openToggleExpenditure ? "rotate-180" : ""
              }`}
            />
          </div>

          {openToggleExpenditure && (
            <div className="mt-4">
              <button
                className="px-4 py-2 mb-3 text-sm font-medium text-white bg-green-500 rounded-md hover:bg-green-600 transition"
                onClick={() => openCreateExpenditureModal()}
              >
                New Expense
              </button>

              <button
                className="px-4 py-2 ml-2 mb-3 text-sm font-medium text-white bg-green-500 rounded-md hover:bg-green-600 transition"
                onClick={() => openCategoryModal("expenditure")}
              >
                New Expense Category
              </button>

              {userDashboard?.userExpense?.length > 0 ? (
                <div className="overflow-x-auto">
                  {" "}
                  <table className="w-full bg-white rounded-lg shadow-md border">
                    <thead>
                      <tr className="bg-gray-100">
                        <th className="px-4 py-2 border-b">Name</th>
                        <th className="px-4 py-2 border-b">Type</th>
                        <th className="px-4 py-2 border-b">Value</th>
                        <th className="px-4 py-2 border-b">Recurrence</th>
                        <th className="px-4 py-2 border-b">StartDate</th>
                        <th className="px-4 py-2 border-b">Shared</th>
                        <th className="px-4 py-2 border-b"></th>
                      </tr>
                    </thead>
                    <tbody>
                      {userDashboard.userExpense.map((expense, index) => (
                        <tr key={index} className="text-center">
                          <td className="px-4 py-2 border-b">
                            {expense.ExpenditureName || "N/A"}
                          </td>
                          <td className="px-4 py-2 border-b">
                            {expense.Expenditure.ExpenditureTypeName || "N/A"}
                          </td>
                          <td className="px-4 py-2 border-b">
                            {expense.ExpenditureValue || "N/A"}
                          </td>
                          <td className="px-4 py-2 border-b">
                            {expense.ExpenditureRecurrence || "N/A"}
                          </td>
                          <td className="px-4 py-2 border-b">
                            {expense.ExpenditureStartDate
                              ? new Date(
                                  expense.ExpenditureStartDate
                                ).toLocaleDateString()
                              : "N/A"}
                          </td>
                          <td className="px-4 py-2 border-b">
                            {expense.SharedExpenditure ? "Yes" : "No"}
                          </td>

                          <td className="px-4 py-2 border-b space-x-2">
                            <button
                              className="px-3 py-1 text-xs font-medium text-white bg-yellow-500 rounded-md hover:bg-yellow-600"
                              onClick={() =>
                                openUpdateExpenditureModal(expense)
                              }
                            >
                              Edit
                            </button>
                            <button
                              className="px-3 py-1 text-xs font-medium text-white bg-red-500 rounded-md hover:bg-red-600"
                              onClick={() =>
                                openDeleteModal(expense, "expense")
                              }
                            >
                              Delete
                            </button>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              ) : (
                <div className="w-full bg-gray-100 rounded-lg p-4 text-center text-gray-500">
                  You don't have Expenditure yet
                </div>
              )}
            </div>
          )}
        </div>
      </div>

      {isModalAssetOpen && (
        <Modal
          onClose={closeAssetModal}
          title={currentAsset ? "Edit Asset" : "New Asset"}
        >
          <AssetForm
            onSubmit={handleAssetSubmit}
            asset={currentAsset}
            onClose={closeAssetModal}
          />
        </Modal>
      )}

      {isModalIncomeOpen && (
        <Modal
          onClose={closeIncomeModal}
          title={currentIncome ? "Edit Income" : "New Income"}
        >
          <IncomeForm
            onSubmit={handleIncomeSubmit}
            income={currentIncome}
            onClose={closeIncomeModal}
          />
        </Modal>
      )}

      {isModalExpenditureOpen && (
        <Modal
          onClose={closeExpenditureModal}
          title={currentExpenditure ? "Edit Expenses" : "New Expense"}
        >
          <ExpenseForm
            onSubmit={handleExpenseSubmit}
            expense={currentExpenditure}
            onClose={closeExpenditureModal}
          />
        </Modal>
      )}

      {isDeleteModalOpen && (
        <Modal
          onClose={closeDeleteModal}
          title={`Confirm Deletion of ${
            deleteType.charAt(0).toUpperCase() + deleteType.slice(1)
          }`}
          onDelete={handleDeleteItem}
          item={itemToDelete}
        >
          <div>
            <p>
              Are you sure you want to delete the {deleteType}:{" "}
              {itemToDelete ? itemToDelete.IncomeName : "No item to delete"}?
            </p>
          </div>
        </Modal>
      )}

      {isModalTaxOpen && (
        <ModalLarge onClose={closeTaxModal} title={`Setup Taxes`}>
          <TaxForm onClose={closeTaxModal} user={currentUser} />
        </ModalLarge>
      )}
      {isModalCategoryOpen && (
        <ModalLarge onClose={closeCategoryModal} title={`Setup Categories`}>
          <CategoryForm
            onClose={closeCategoryModal}
            category={currentCategory}
          />
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

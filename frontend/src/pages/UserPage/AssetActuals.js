import React, { useEffect, useState } from "react";
import { useCallback } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../../components/AuthContext";

import CategoryForm from "../../components/CategoryForm";
import AssetForm from "../../components/AssetForm";
import Modal from "../../components/Modal";
import ModalLarge from "../../components/ModalLarge";
import toastr from "toastr";
import { Link } from "react-router-dom";
import { ChevronUpIcon } from "@heroicons/react/24/solid";

const UserAsset = () => {
  const { isLoggedIn, loading } = useAuth();
  const [userDashboard, setUserDashboard] = useState(null);
  const [isModalAssetOpen, setIsModalAssetOpen] = useState(false);
  const [isModalCategoryOpen, setisModalCategoryOpen] = useState(false);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [deleteType, setDeleteType] = useState("");
  const [itemToDelete, setItemToDelete] = useState(null);
  const [currentUser, setCurrentUser] = useState(null);
  const [currentAsset, setCurrentAsset] = useState(null);
  const [currentCategory, setCurrentCategory] = useState(null);
  const [openToggleAsset, setopenToggleAsset] = useState(true);

  const navigate = useNavigate();

  // Check if the user is logged in
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

  // Open the Category Modal (pass the category type to set)
  const openCategoryModal = useCallback((type) => {
    setCurrentCategory(type); // Set the category type
    setisModalCategoryOpen(true); // Open the modal
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

  // Close the Category Modal (reset state)
  const closeCategoryModal = useCallback(() => {
    setisModalCategoryOpen(false); // Close the modal
    setCurrentCategory(null); // Reset the category type
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
      : { ...data, userID: currentUser };

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
      <div className="mx-auto p-6 mt-4 bg-white rounded-lg shadow-lg">
        <div className="flex items-center justify-between">
          <Link
            to={`/user`}
            className="px-4 py-2 text-sm font-medium text-white bg-blue-500 rounded-md hover:bg-blue-600 transition"
          >
            Back
          </Link>
          <h1 className="text-2xl font-bold text-blue-600 w-full text-center">
            Assets Actuals
          </h1>
        </div>
      </div>
      <div className="space-y-6">
        <div className="bg-white shadow-md rounded-lg p-4 mt-2">
          <div
            className="flex justify-between w-full px-4 py-2 text-left text-sm font-medium text-white bg-blue-500 rounded-md hover:bg-blue-600 transition cursor-pointer"
            onClick={toggleAssetPanel}
          >
            <span>My Actuals Assets</span>
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
                New Actuals
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
      {isModalCategoryOpen && (
        <ModalLarge onClose={closeCategoryModal} title={`Setup Categories`}>
          <CategoryForm
            onClose={closeCategoryModal}
            category={currentCategory}
          />
        </ModalLarge>
      )}
    </div>
  );
};

export default UserAsset;

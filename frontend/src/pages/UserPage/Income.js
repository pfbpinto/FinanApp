import React, { useEffect, useState, useCallback } from "react";
import { useLocation } from "react-router-dom";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../../components/AuthContext";
import { Link } from "react-router-dom";
import IncomeForm from "../../components/IncomeForm";
import IncomePage from "../UserPage/IncomeItem";
import CategoryForm from "../../components/CategoryForm";
import Modal from "../../components/Modal";
import ModalLarge from "../../components/ModalLarge";
import toastr from "toastr";
import { ChevronUpIcon } from "@heroicons/react/24/solid";

const UserIncome = () => {
  const { isLoggedIn, loading } = useAuth();
  const { state } = useLocation();
  const currentUser = state?.userID || null;
  const [userFinancialItem, setUserFinancialItem] = useState(null);

  const [userUserCagegory, setUserUserCagegory] = useState(null);
  const [userIncomeTypes, setUserIncomeTypes] = useState(null);
  const [currency, setCurrency] = useState(null);
  const [recurrency, setRecurrency] = useState(null);

  const [isModalIncomeOpen, setIsModalIncomeOpen] = useState(false);
  const [isModalIncomePageOpen, setIsModalIncomePageOpen] = useState(false);

  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [isModalCategoryOpen, setisModalCategoryOpen] = useState(false);

  const [deleteType, setDeleteType] = useState("");
  const [itemToDelete, setItemToDelete] = useState(null);
  const [currentIncome, setcurrentIncome] = useState(null);
  const [currentCategory, setCurrentCategory] = useState(null);
  const [openToggleIncome, setOpenToggleIncome] = useState(true);

  const navigate = useNavigate();

  // Check if the user is logged in
  useEffect(() => {
    if (loading) return;
    if (!isLoggedIn) navigate("/login");
  }, [isLoggedIn, loading, navigate]);

  useEffect(() => {
    if (isLoggedIn && !currentUser) {
      fetchUserIncome();
    }
  }, [isLoggedIn, currentUser]);

  // Fetch User Dashboard info
  const fetchUserIncome = () => {
    fetch("/api/user-income", {
      method: "GET",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
    })
      .then((response) => response.json())
      .then((data) => {
        setCurrency(data.currency);
        setRecurrency(data.recurrency);
        setUserIncomeTypes(data.income_type);
        setUserUserCagegory(data.user_categories);
        setUserFinancialItem(data.financial_user_items);
      })
      .catch((error) => console.error("Error fetching user data:", error));
  };

  // Toggle the Income panel (open/close)
  const toggleIncomePanel = useCallback(() => {
    setOpenToggleIncome((prevState) => !prevState); // Toggle state
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

  // Open the Income Page Modal
  const openViewIncomeModal = useCallback((income) => {
    setcurrentIncome(income); // Ensure the current Income is empty
    setIsModalIncomePageOpen(true); // Open the income create modal
  }, []);

  // Open the Delete Modal (set the item and type to delete)
  const openDeleteModal = useCallback((item, type) => {
    console.log(item);
    if (item) {
      setItemToDelete(item); // Set the item to delete
      setDeleteType(type); // Set the type of item to delete
      setIsDeleteModalOpen(true); // Open the delete modal
    } else {
      console.error("Invalid item or missing ID"); // Log error if item or ID is missing
    }
  }, []);

  // Open the Category Modal (pass the category type to set)
  const openCategoryModal = useCallback((type) => {
    setCurrentCategory(type); // Set the category type
    setisModalCategoryOpen(true); // Open the modal
  }, []);

  // Close the Delete Modal (reset state)
  const closeDeleteModal = useCallback(() => {
    setIsDeleteModalOpen(false); // Close the modal
    setItemToDelete(null); // Reset the item to delete
  }, []);

  // Close the Income Modal (reset state)
  const closeIncomeModal = useCallback(() => {
    setIsModalIncomeOpen(false); // Close the modal
    setcurrentIncome(null); // Reset the category type
  }, []);

  // Close the Income Modal (reset state)
  const closeIncomePageModal = useCallback(() => {
    setIsModalIncomePageOpen(false); // Close the modal
    setcurrentIncome(null); // Reset the category type
  }, []);

  // Close the Category Modal (reset state)
  const closeCategoryModal = useCallback(() => {
    setisModalCategoryOpen(false); // Close the modal
    setCurrentCategory(null); // Reset the category type
  }, []);

  // Income Create and Update handle
  const handleIncomeSubmit = (data) => {
    const method = currentIncome ? "PUT" : "POST";
    const url = currentIncome
      ? `/api/income-update/${data.FinancialUserItemId}`
      : "/api/income";

    // 👇 Corrigir ID vazio para null (importante para evitar erro no backend)
    const cleanData = {
      ...data,
      FinancialUserItemId:
        data.FinancialUserItemId === "" ? null : data.FinancialUserItemId,
    };

    // 👇 Adicionar o IncomeValue só no update
    const dataToSend = currentIncome
      ? { ...cleanData, IncomeValue: parseFloat(cleanData.IncomeValue) }
      : cleanData;

    fetch(url, {
      method: method,
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(dataToSend),
      credentials: "include",
    })
      .then((response) => {
        if (!response.ok) {
          return response.json().then((errorData) => {
            return Promise.reject(errorData);
          });
        }
        return response.json();
      })
      .then(() => {
        toastr.success(currentIncome ? "Income updated!" : "Income created!");
        setIsModalIncomeOpen(false);
        fetchUserIncome();
      })
      .catch((error) => {
        console.error("Error saving income:", error);
        toastr.error(`Error: ${error.message}`);
      });
  };

  // Item Delete handle
  const handleDeleteItem = (item) => {
    if (!item || !item.financialUserItemId) return;
    fetch(`/api/delete-${deleteType}/${item.financialUserItemId}`, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ itemId: item.financialUserItemId }),
      credentials: "include",
    })
      .then((response) => response.json())
      .then(() => {
        toastr.success(
          `${
            deleteType.charAt(0).toUpperCase() + deleteType.slice(1)
          } Successfully Deleted!`
        );
        fetchUserIncome();
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
            My Incomes
          </h1>
        </div>
      </div>

      <div className="space-y-6">
        <div className="bg-white shadow-md rounded-lg p-4 mt-2">
          <div
            className="flex justify-between w-full px-4 py-2 text-left text-sm font-medium text-white bg-blue-500 rounded-md hover:bg-blue-600 transition cursor-pointer"
            onClick={toggleIncomePanel}
          >
            <span>Incomes</span>
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

              {userFinancialItem && userFinancialItem.length > 0 ? (
                <div className="overflow-x-auto">
                  <table className="w-full bg-white rounded-lg shadow-md border">
                    <thead>
                      <tr className="bg-gray-100">
                        <th className="px-4 py-2 border-b"></th>
                        <th className="px-4 py-2 border-b">Name</th>
                        <th className="px-4 py-2 border-b">Type</th>
                        <th className="px-4 py-2 border-b">Recurrency</th>
                        <th className="px-4 py-2 border-b">Start Date</th>

                        <th className="px-4 py-2 border-b"></th>
                      </tr>
                    </thead>
                    <tbody>
                      {userFinancialItem.map((income, index) => (
                        <tr key={index} className="text-center">
                          <td className="px-4 py-2 border-b space-x-2">
                            <button
                              className="px-3 py-1 text-xs font-medium text-white bg-blue-500 rounded-md hover:bg-blue-600"
                              onClick={() => openViewIncomeModal(income)}
                            >
                              Open
                            </button>
                          </td>

                          <td className="px-4 py-2 border-b">
                            {income.financialUserItemName || "N/A"}
                          </td>
                          <td className="px-4 py-2 border-b">
                            {income.incomeTypeName || "N/A"}
                          </td>
                          <td className="px-4 py-2 border-b">
                            {income.recurrencyName || "N/A"}
                          </td>

                          <td className="px-4 py-2 border-b">
                            {income.createdAt
                              ? new Date(income.createdAt).toLocaleDateString()
                              : "N/A"}
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

      {isModalIncomeOpen && (
        <Modal
          onClose={closeIncomeModal}
          title={currentIncome ? "Edit Income" : "New Income"}
        >
          <IncomeForm
            onSubmit={handleIncomeSubmit}
            income={currentIncome}
            userCategory={userUserCagegory}
            currency={currency}
            recurrency={recurrency}
            onClose={closeIncomeModal}
          />
        </Modal>
      )}
      {isModalIncomePageOpen && (
        <ModalLarge onClose={closeIncomePageModal}>
          <IncomePage income={currentIncome} onClose={closeIncomePageModal} />
        </ModalLarge>
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
        <ModalLarge onClose={closeCategoryModal} title={`Setup Category`}>
          <CategoryForm
            onClose={closeCategoryModal}
            category={currentCategory}
            entityTypes={userIncomeTypes}
            userCategory={userUserCagegory}
          />
        </ModalLarge>
      )}
    </div>
  );
};

export default UserIncome;

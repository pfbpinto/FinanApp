import React, { useEffect, useState } from "react";
import { useAuth } from "../components/AuthContext";
import PropTypes from "prop-types";
import ConfirmDeleteModal from "../components/ModalConfirm";
import toastr from "toastr";

const CategoryForm = ({ category, onClose }) => {
  const [name, setName] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [categoryData, setCategoryData] = useState(null);
  // States confirmation modal and to Tax delete
  const [isConfirmDeleteModalOpen, setIsConfirmDeleteModalOpen] =
    useState(false);
  const [CategoryToDelete, setCategoryToDelete] = useState(null); // Para armazenar o ID da taxa a ser excluída
  const [deleteMessage, setDeleteMessage] = useState("");

  // Getting authentication state from context
  const { isLoggedIn } = useAuth();

  // Mapping Category behavior attributes
  const categoryTextMapping = {
    asset: {
      label: "Asset Type Name",
      fieldName: "AssetTypeName",
      name: "AssetTypeName",
      jsonData: "assetType",
    },
    tax: {
      label: "Tax Type Name",
      fieldName: "TaxTypeName",
      name: "TaxTypeName",
      jsonData: "taxTypes",
    },
    expenditure: {
      label: "Expenditure Type Name",
      fieldName: "ExpenditureTypeName",
      name: "ExpenditureTypeName",
      jsonData: "expenditureType",
    },
    file: {
      label: "File Type Name",
      fieldName: "FileTypeName",
      name: "FileTypeName",
      jsonData: "fileType",
    },
    group: {
      label: "Group Type Name",
      fieldName: "GroupTypeName",
      name: "GroupTypeName",
      jsonData: "groupType",
    },
    income: {
      label: "Income Type Name",
      fieldName: "IncomeTypeName",
      name: "IncomeTypeName",
      jsonData: "incomeType",
    },
  };

  // Set choosed category
  const dynamicTexts = categoryTextMapping[category];

  // useEffect to fetch tax types and taxes when the user is logged in
  useEffect(() => {
    if (isLoggedIn) {
      fetch("/api/categories", {
        method: "GET",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
      })
        .then((response) => response.json())
        .then((data) => {
          setCategoryData(data[dynamicTexts.jsonData] || []); // Storing existing taxes
        })
        .catch((error) => console.error("Error fetching data:", error));
    }
  }, [isLoggedIn, dynamicTexts.jsonData]); // Re-fetch data if login status changes

  const handleCategorySubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    // Send Payload to API
    const payload = {
      category,
      name: name,
      model: dynamicTexts.jsonData,
      field: dynamicTexts.fieldName,
    };

    fetch(`/api/create-category`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
      credentials: "include", // Include credentials like cookies for authentication
    })
      .then((response) => {
        // Check if the response status is ok (status code 200-299)
        if (!response.ok) {
          // If response is not ok, extract error message from the JSON response
          return response.json().then((errorData) => {
            setLoading(false);
            return Promise.reject(errorData);
          });
        }
        // If successful, return the response as JSON
        return response.json();
      })
      .then((data) => {
        // If the data is null (error occurred), skip the success handler
        if (!data) return;
        // Show a success message if the tax was created successfully
        toastr.success("Tax Successfully Created.");
        // Update table with the new value
        setCategoryData((prevTaxes) => [...prevTaxes, data.category]);
        // stop loading
        setLoading(false);
      })
      .catch((error) => {
        // Log the error and display an error message using toastr
        console.error("Error saving Tax:", error);
        toastr.error(`Error: ${error.message}`);
      });
  };

  // Handle confirmation modal for Tax delete
  const handleCategoryDelete = (cat) => {
    setCategoryToDelete(cat);
    setDeleteMessage(
      `Are you sure you want to delete the ${category} type "${
        cat[dynamicTexts.fieldName]
      }"?`
    );
    setIsConfirmDeleteModalOpen(true);
  };

  // Handle tax delete
  const deleteCategory = async () => {
    try {
      const response = await fetch(
        `/api/delete-category/${CategoryToDelete.ID}`,
        {
          method: "DELETE",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            categoryID: CategoryToDelete.ID,
            model: dynamicTexts.jsonData,
          }),
          credentials: "include",
        }
      );

      if (!response.ok) {
        throw new Error("Failed to delete tax");
      }

      // // Update table removing the deleted Tax
      setCategoryData((prevTaxes) =>
        prevTaxes.filter((tax) => tax.ID !== CategoryToDelete.ID)
      );

      toastr.success("Category Successfully Deleted!");
    } catch (error) {
      console.error("Error deleting category:", error);
      toastr.error("Failed to delete category");
    }
  };

  return (
    <div>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
        <div className="max-w-none md:max-w-md">
          <h2 className="text-xl font-semibold mb-4">
            Insert new {category} Category
          </h2>
          <form onSubmit={handleCategorySubmit} className="space-y-4">
            <div>
              <label
                htmlFor="categoryName"
                className="block text-sm font-medium text-gray-700"
              >
                {dynamicTexts.label}
              </label>
              <input
                id="categoryName"
                type="text"
                name={dynamicTexts.name}
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder={`Insert ${dynamicTexts.label.toLowerCase()}`}
                required
                className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
              />
            </div>

            {error && <div className="text-red-500 text-sm mt-2">{error}</div>}

            <div className="flex justify-start">
              <button
                type="button"
                onClick={onClose}
                className="mr-2 px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-gray-500"
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={loading}
                className="px-4 py-2 text-sm font-medium text-white bg-blue-500 rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                {loading ? "Saving..." : "Save"}
              </button>
            </div>
          </form>
        </div>

        {/* Section to display existing taxes */}
        <div className="min-w-80">
          <h2 className="text-xl font-semibold mb-4">Current {category}s</h2>
          <div
            id="tax-table-container"
            className="overflow-x-auto max-h-80 overflow-y-auto border rounded-md"
          >
            <table className="min-w-full table-auto border-collapse">
              <thead className="sticky top-0 bg-white shadow-md">
                <tr>
                  <th className="px-4 py-2 text-left border-b">ID</th>
                  <th className="px-4 py-2 text-left border-b">Name</th>
                  <th></th>
                </tr>
              </thead>
              <tbody>
                {categoryData?.map((category) => (
                  <tr key={category.ID}>
                    <td className="px-4 py-2 border-b">{category.ID}</td>
                    <td className="px-4 py-2 border-b">
                      {category[dynamicTexts.fieldName]}
                    </td>

                    <td className="px-4 py-2 border-b space-x-2">
                      <button
                        className="px-3 py-1 text-xs font-medium text-white bg-red-500 rounded-md hover:bg-red-600"
                        onClick={() => handleCategoryDelete(category)}
                      >
                        X
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>
      <br></br>
      {/* Cancel Button at the bottom */}
      <div className="flex justify-end">
        <button
          type="button"
          onClick={onClose}
          className="mr-2 px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-gray-500"
        >
          Cancel
        </button>
      </div>
      <div>
        <ConfirmDeleteModal
          isOpen={isConfirmDeleteModalOpen}
          onClose={() => setIsConfirmDeleteModalOpen(false)}
          message={deleteMessage}
          onConfirm={deleteCategory}
        />
      </div>
    </div>
  );
};

CategoryForm.propTypes = {
  category: PropTypes.string.isRequired,
  onClose: PropTypes.func.isRequired,
};

export default CategoryForm;

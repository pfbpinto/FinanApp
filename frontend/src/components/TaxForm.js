import React, { useEffect, useState } from "react";
import { useAuth } from "../components/AuthContext";
import ConfirmDeleteModal from "../components/ModalConfirm";
import toastr from "toastr";

// TaxForm component for inserting new taxes
function TaxForm({ onClose, user }) {
  // Getting authentication state from context
  const { isLoggedIn } = useAuth();
  // State to store form data
  const [formData, setFormData] = useState({
    UserID: 0,
    TaxName: "",
    TaxTypeID: "",
    TaxPercentage: "",
    TaxApplicableCycle: "",
  });

  // States to store tax types and existing taxes
  const [taxesTypes, setTaxesTypes] = useState([]);
  const [taxes, setTaxes] = useState([]);
  // States confirmation modal and to Tax delete
  const [isConfirmDeleteModalOpen, setIsConfirmDeleteModalOpen] =
    useState(false);
  const [taxToDelete, setTaxToDelete] = useState(null); // Para armazenar o ID da taxa a ser excluída
  const [deleteMessage, setDeleteMessage] = useState("");

  useEffect(() => {
    if (user) {
      setFormData((prevData) => ({ ...prevData, UserID: user }));
    }
  }, [user]);

  // useEffect to fetch tax types and taxes when the user is logged in
  useEffect(() => {
    if (isLoggedIn) {
      fetch("/api/get-taxes", {
        method: "GET",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
      })
        .then((response) => response.json())
        .then((data) => {
          setTaxesTypes(data.taxTypes || []); // Storing tax types
          setTaxes(data.taxes || []); // Storing existing taxes
        })
        .catch((error) => console.error("Error fetching data:", error));
    }
  }, [isLoggedIn]); // Re-fetch data if login status changes

  // handleChange function to handle input changes
  const handleChange = (e) => {
    const { name, value, type, checked } = e.target;
    // Update the form data state with the new value from the input field
    setFormData((prevData) => ({
      ...prevData,
      [name]: type === "checkbox" ? checked : value, // Check if the input is a checkbox
    }));
  };

  // handleSubmit function to handle form submission
  const handleTaxSubmit = (e) => {
    // Send the data to the server
    e.preventDefault(); // Prevent page reload on form submission

    fetch(`/api/create-taxes`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(formData),
      credentials: "include", // Include credentials like cookies for authentication
    })
      .then((response) => {
        // Check if the response status is ok (status code 200-299)
        if (!response.ok) {
          // If response is not ok, extract error message from the JSON response
          return response.json().then((errorData) => {
            // Display error message using toastr
            toastr.error(errorData.message || "Unknown error");
            // Return null to prevent further execution of the next .then() block
            return null;
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
        setTaxes((prevTaxes) => [...prevTaxes, data.taxes]);
      })
      .catch((error) => {
        // Log the error and display an error message using toastr
        console.error("Error saving Tax:", error);
        toastr.error(`Error: ${error.message}`);
      });
  };

  // Handle confirmation modal for Tax delete
  const handleTaxDelete = (tax) => {
    setTaxToDelete(tax);
    setDeleteMessage(
      `Are you sure you want to delete the tax "${tax.TaxName}"?`
    );
    setIsConfirmDeleteModalOpen(true);
  };
  // Handle tax delete
  const deleteTax = async () => {
    try {
      const response = await fetch(`/api/delete-tax/${taxToDelete.ID}`, {
        method: "DELETE",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ taxID: taxToDelete.ID }),
        credentials: "include",
      });

      if (!response.ok) {
        throw new Error("Failed to delete tax");
      }

      // // Update table removing the deleted Tax
      setTaxes((prevTaxes) =>
        prevTaxes.filter((tax) => tax.ID !== taxToDelete.ID)
      );

      toastr.success("Tax Successfully Deleted!");
    } catch (error) {
      console.error("Error deleting tax:", error);
      toastr.error("Failed to delete tax");
    }
  };

  return (
    <div>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
        <div className="max-w-none md:max-w-md">
          <h2 className="text-xl font-semibold mb-4">Insert new Tax</h2>
          <form onSubmit={handleTaxSubmit}>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-4">
              {/* Input for Tax Name */}
              <div>
                <label className="block text-sm font-medium text-gray-700">
                  Tax Name
                </label>
                <input
                  required
                  type="text"
                  name="TaxName"
                  onChange={handleChange} // handleChange is used here to update state
                  className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                />
              </div>

              {/* Input for Tax Type */}
              <div>
                <label className="block text-sm font-medium text-gray-700">
                  Tax Type
                </label>
                <select
                  required
                  name="TaxTypeID"
                  onChange={handleChange} // handleChange is also used here
                  className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                >
                  <option value="">Select Tax Type</option>
                  {taxesTypes?.map((type) => (
                    <option key={type.ID} value={type.ID}>
                      {type.TaxTypeName}
                    </option>
                  ))}
                </select>
              </div>
            </div>

            {/* Grid for Tax Percentage, Range, and Cycle inputs */}
            <div className="grid grid-cols-3 gap-4 mb-4">
              {/* Input for Tax Percentage */}
              <div>
                <label className="block text-sm font-medium text-gray-700">
                  Tax Percentage
                </label>
                <input
                  required
                  max={100}
                  type="number"
                  name="TaxPercentage"
                  onChange={handleChange} // handleChange is used for number inputs as well
                  className="mt-1 block w-full max-w-xs px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                />
              </div>

              {/* Input for Tax Percentage Range */}
              <div>
                <label className="block text-sm font-medium text-gray-700">
                  Tax % Range
                </label>
                <input
                  max={100}
                  type="number"
                  name="TaxPercentageRange"
                  onChange={handleChange} // handleChange is used here
                  className="mt-1 block w-full max-w-xs px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                />
              </div>

              {/* Input for Tax Applicable Cycle */}
              <div>
                <label className="block text-sm font-medium text-gray-700">
                  Tax Cycle
                </label>
                <input
                  type="text"
                  name="TaxApplicableCycle"
                  onChange={handleChange} // handleChange handles the input changes here
                  className="mt-1 block w-full max-w-xs px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                />
              </div>
            </div>

            <div className="flex justify-start">
              {/* Cancel Button */}
              <button
                type="button"
                onClick={onClose}
                className="mr-2 px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-gray-500"
              >
                Cancel
              </button>
              {/* Submit Button */}
              <button
                type="submit"
                className="px-4 py-2 text-sm font-medium text-white bg-blue-500 rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                Save
              </button>
            </div>
          </form>
        </div>

        {/* Section to display existing taxes */}
        <div className="min-w-80">
          <h2 className="text-xl font-semibold mb-4">Current Taxes</h2>
          <div
            id="tax-table-container"
            className="overflow-x-auto max-h-80 overflow-y-auto border rounded-md"
          >
            <table className="min-w-full table-auto border-collapse">
              <thead className="sticky top-0 bg-white shadow-md">
                <tr>
                  <th className="px-4 py-2 text-left border-b">Tax Name</th>
                  <th className="px-4 py-2 text-left border-b">Tax Type</th>
                  <th className="px-4 py-2 text-left border-b">%</th>
                  <th className="px-4 py-2 text-left border-b">Cycle</th>
                  <th></th>
                </tr>
              </thead>
              <tbody>
                {taxes
                  ?.slice()
                  .sort((a, b) => b.ID - a.ID) // Sort by ID descending
                  .map((tax) => (
                    <tr key={tax.ID}>
                      <td className="px-4 py-2 border-b">{tax.TaxName}</td>
                      <td className="px-4 py-2 border-b">
                        {taxesTypes?.find((type) => type.ID === tax.TaxTypeID)
                          ?.TaxTypeName || "Unknown"}
                      </td>
                      <td className="px-4 py-2 border-b">
                        {tax.TaxPercentage}%
                      </td>
                      <td className="px-4 py-2 border-b">
                        {tax.TaxApplicableCycle}
                      </td>
                      <td className="px-4 py-2 border-b space-x-2">
                        <button
                          className="px-3 py-1 text-xs font-medium text-white bg-red-500 rounded-md hover:bg-red-600"
                          onClick={() => handleTaxDelete(tax)}
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
          onConfirm={deleteTax}
        />
      </div>
    </div>
  );
}

export default TaxForm;

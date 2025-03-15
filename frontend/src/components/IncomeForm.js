import React, { useEffect, useState } from "react";
import { useAuth } from "../components/AuthContext";

function IncomeForm({ onSubmit, income, onClose }) {
  const [formData, setFormData] = useState(income || { UserTaxes: [] });
  const [incomeTypes, setIncomeTypes] = useState([]);
  const [taxes, setTaxes] = useState([]);
  const { isLoggedIn } = useAuth();

  useEffect(() => {
    if (isLoggedIn) {
      fetch("/api/income-type", {
        method: "GET",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
      })
        .then((response) => response.json())
        .then((data) => {
          setIncomeTypes(data.incomeTypes || []);
          setTaxes(data.taxes || []);
        })
        .catch((error) => console.error("Error fetching data:", error));
    }
  }, [isLoggedIn]);

  const handleChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData((prevData) => ({
      ...prevData,
      [name]: type === "checkbox" ? checked : value,
    }));
  };

  console.log(income);

  const handleTaxChange = (e, tax) => {
    const { checked } = e.target;

    setFormData((prevData) => {
      const updatedTaxes = checked
        ? [
            ...(prevData.UserTaxes || []),
            {
              TaxID: tax.ID,
              TaxName: tax.TaxName,
              TaxPercentage: tax.TaxPercentage,
            },
          ]
        : (prevData.UserTaxes || []).filter((t) => t.TaxID !== tax.ID);

      return {
        ...prevData,
        UserTaxes: updatedTaxes,
      };
    });
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    onSubmit(formData);
  };

  const handleTypeChange = (event) => {
    const value = event.target.value;

    if (value === "") {
      // Se o valor for vazio, limpe o campo de IncomeType
      setFormData((prevData) => ({
        ...prevData,
        IncomeType: null, // Ou um valor padrão
      }));
    } else {
      // Caso contrário, atualize o valor de IncomeType com o ID selecionado
      setFormData((prevData) => ({
        ...prevData,
        IncomeType: { ID: value },
        IncomeTypeID: value,
      }));
    }
  };

  console.log(formData);
  return (
    <form onSubmit={handleSubmit}>
      <div className="mb-4">
        <label className="block text-sm font-medium text-gray-700">
          Income Name
        </label>
        <input
          required
          type="text"
          name="IncomeName"
          value={formData.IncomeName || ""}
          onChange={handleChange}
          className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
        />
      </div>

      <div className="mb-4">
        <label className="block text-sm font-medium text-gray-700">
          Income Value
        </label>
        <input
          required
          type="number"
          name="IncomeValue"
          value={formData.IncomeValue || ""}
          onChange={handleChange}
          className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
        />
      </div>

      <div className="mb-4">
        <label className="block text-sm font-medium text-gray-700">
          Income Type
        </label>
        <select
          required
          name="IncomeType"
          value={formData?.IncomeType?.ID ? String(formData.IncomeType.ID) : ""}
          onChange={handleTypeChange}
          className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
        >
          <option value="">Select Income Type</option>
          {incomeTypes?.map((type) => (
            <option key={type.ID} value={type.ID}>
              {type.IncomeTypeName}
            </option>
          ))}
        </select>
      </div>

      <div className="mb-4">
        <label className="block text-sm font-medium text-gray-700">
          Recurrence
        </label>
        <input
          required
          type="text"
          name="IncomeRecurrence"
          value={formData.IncomeRecurrence || ""}
          onChange={handleChange}
          className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
        />
      </div>

      <div className="mb-4">
        <label className="block text-sm font-medium text-gray-700">
          Start Date
        </label>
        <input
          type="date"
          name="IncomeStartDate"
          value={
            formData.IncomeStartDate
              ? new Date(formData.IncomeStartDate).toISOString().split("T")[0]
              : ""
          }
          onChange={handleChange}
          className="mt-1 block w-full px-4 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
          required
        />
      </div>

      <div className="mb-4">
        <div className="mt-1 border border-gray-300 rounded-md p-2">
          <input
            type="checkbox"
            name="SharedIncome"
            checked={formData.SharedIncome || false}
            onChange={handleChange}
            className="mr-2"
          />
          <label>Shared Income</label>
        </div>
      </div>

      <div className="mb-4">
        <label className="block text-sm font-medium text-gray-700">
          Owning Percentage
        </label>
        <input
          required
          type="number"
          name="OwningPercentage"
          value={formData.OwningPercentage || ""}
          onChange={handleChange}
          className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
        />
      </div>

      {/* Taxes Checkbox */}
      <div className="mb-4">
        <label className="block text-sm font-medium text-gray-700">
          Taxes (Optional)
        </label>
        <div className="mt-1 border border-gray-300 rounded-md p-2">
          {taxes?.map((tax) => (
            <div key={tax.ID} className="flex items-center mb-1">
              <input
                type="checkbox"
                checked={
                  formData.UserTaxes?.some((t) => t.TaxID === tax.ID) || false
                }
                onChange={(e) => handleTaxChange(e, tax)}
                className="mr-2"
              />
              <label>
                {tax.TaxName} ({tax.TaxPercentage}%)
              </label>
            </div>
          ))}
        </div>
      </div>

      <div className="flex justify-end">
        <button
          type="button"
          onClick={onClose}
          className="mr-2 px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-gray-500"
        >
          Cancel
        </button>
        <button
          type="submit"
          className="px-4 py-2 text-sm font-medium text-white bg-blue-500 rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          Save
        </button>
      </div>
    </form>
  );
}

export default IncomeForm;

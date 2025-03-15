import React, { useEffect, useState } from "react";
import { useAuth } from "../components/AuthContext";

function ExpenseForm({ onSubmit, expense, onClose }) {
  const [formData, setFormData] = useState(expense || []);
  const [expenseTypes, setExpenseTypes] = useState([]);
  const { isLoggedIn } = useAuth();

  useEffect(() => {
    if (isLoggedIn) {
      fetch("/api/expense-type", {
        method: "GET",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
      })
        .then((response) => response.json())
        .then((data) => {
          setExpenseTypes(data.expenseTypes || []);
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
        ExpenditureID: null, // Ou um valor padrão
      }));
    } else {
      // Caso contrário, atualize o valor de IncomeType com o ID selecionado
      setFormData((prevData) => ({
        ...prevData,
        ExpenditureID: value,
        Expenditure_ID: value,
      }));
    }
  };

  console.log(formData);

  return (
    <form onSubmit={handleSubmit}>
      <div className="mb-4">
        <label className="block text-sm font-medium text-gray-700">
          Expense Name
        </label>
        <input
          required
          type="text"
          name="ExpenditureName"
          value={formData.ExpenditureName || ""}
          onChange={handleChange}
          className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
        />
      </div>

      <div className="mb-4">
        <label className="block text-sm font-medium text-gray-700">
          Expense Value
        </label>
        <input
          required
          type="number"
          name="ExpenditureValue"
          value={formData.ExpenditureValue || ""}
          onChange={handleChange}
          className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
        />
      </div>

      <div className="mb-4">
        <label className="block text-sm font-medium text-gray-700">
          Expense Type
        </label>
        <select
          required
          name="ExpenditureID"
          value={formData?.ExpenditureID ? String(formData.ExpenditureID) : ""}
          onChange={handleTypeChange}
          className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
        >
          <option value="">Select Expense Type</option>
          {expenseTypes?.map((type) => (
            <option key={type.ID} value={type.ID}>
              {type.ExpenditureTypeName}
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
          name="ExpenditureRecurrence"
          value={formData.ExpenditureRecurrence || ""}
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
          name="ExpenditureStartDate"
          value={
            formData.ExpenditureStartDate
              ? new Date(formData.ExpenditureStartDate)
                  .toISOString()
                  .split("T")[0]
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
            name="SharedExpenditure"
            checked={formData.SharedExpenditure || false}
            onChange={handleChange}
            className="mr-2"
          />
          <label>Shared Expense</label>
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

export default ExpenseForm;

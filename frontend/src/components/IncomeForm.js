import React, { useEffect, useState } from "react";

function IncomeForm({ onSubmit, income, userCategory, onClose }) {
  const [formData, setFormData] = useState(
    income || {
      user_financial_forecast_name: "",
      user_financial_forecast_amount: "",
      entity_type_id: "",
      user_financial_forecast_begin_date: "",
    }
  );

  const [incomeTypes, setIncomeTypes] = useState([]);

  useEffect(() => {
    setIncomeTypes(userCategory);
  }, [userCategory]);

  useEffect(() => {
    if (income) {
      setFormData(income);
    }
  }, [income]);

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

    setFormData((prevData) => ({
      ...prevData,
      income_type_id: value || null, // Se vazio, define como null
    }));
  };

  return (
    <form onSubmit={handleSubmit}>
      <div className="mb-4">
        <label className="block text-sm font-medium text-gray-700">
          Income Name
        </label>
        <input
          required
          type="text"
          name="user_financial_forecast_name"
          value={formData.user_financial_forecast_name || ""}
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
          name="user_financial_forecast_amount"
          value={formData.user_financial_forecast_amount || ""}
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
          name="income_type_id"
          value={
            formData?.income_type_id ? String(formData.income_type_id) : ""
          }
          onChange={handleTypeChange}
          className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
        >
          <option value="">Select Income Category</option>
          {incomeTypes?.map((type) => (
            <option key={type.user_category_id} value={type.user_category_id}>
              {type.user_category_name}
            </option>
          ))}
        </select>
      </div>

      <div className="mb-4">
        <label className="block text-sm font-medium text-gray-700">
          Begin Date
        </label>
        <input
          type="date"
          name="user_financial_forecast_begin_date"
          value={
            formData.user_financial_forecast_begin_date
              ? new Date(formData.user_financial_forecast_begin_date)
                  .toISOString()
                  .split("T")[0]
              : ""
          }
          onChange={handleChange}
          className="mt-1 block w-full px-4 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
          required
        />
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

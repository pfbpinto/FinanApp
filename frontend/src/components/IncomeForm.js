import React, { useEffect, useState } from "react";

function IncomeForm({ onSubmit, income, currency, recurrency, onClose }) {
  const [formData, setFormData] = useState({
    FinancialUserItemId: "",
    FinancialUserItemName: "",
    RecurrencyID: "",
    CurrencyID: "",
    amount: "",
  });

  useEffect(() => {
    if (income) {
      setFormData({
        FinancialUserItemId: income.financialUserItemId || "",
        FinancialUserItemName: income.financialUserItemName || "",
        RecurrencyID: income.recurrencyId || "",
        CurrencyID: income.currencyId || "",
        amount: income.amount || "",
      });
    }
  }, [income]);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prevData) => ({
      ...prevData,
      [name]: value,
    }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    onSubmit(formData);
  };

  return (
    <form onSubmit={handleSubmit}>
      <div className="mb-4">
        <label className="block text-sm font-medium text-gray-700">
          Financial Item Name
        </label>
        <input
          type="text"
          name="FinancialUserItemName"
          value={formData.FinancialUserItemName}
          onChange={handleChange}
          required
          className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
        />
      </div>

      <div>
        <div className="mb-4 flex gap-4">
          {/* Amount Input */}
          <div className="w-1/2">
            <label className="block text-sm font-medium text-gray-700">
              Amount
            </label>
            <input
              type="number"
              name="amount"
              value={formData.amount}
              onChange={handleChange}
              required
              className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            />
          </div>

          {/* Currency Select */}
          <div className="w-1/2">
            <label className="block text-sm font-medium text-gray-700">
              Currency
            </label>
            <select
              name="CurrencyID"
              value={formData.CurrencyID}
              onChange={handleChange}
              required
              className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="">Select Currency</option>
              {Array.isArray(currency) &&
                currency.map((curr, index) => (
                  <option
                    key={curr.currency_id || index}
                    value={curr.currency_id}
                  >
                    {curr.currency_name}
                  </option>
                ))}
            </select>
          </div>
        </div>

        <div className="mb-4 flex gap-4">
          {/* Recurrency Select */}
          <div className="w-1/2">
            <label className="block text-sm font-medium text-gray-700">
              Recurrency
            </label>
            <select
              name="RecurrencyID"
              value={formData.RecurrencyID}
              onChange={handleChange}
              required
              className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="">Select Recurrency</option>
              {Array.isArray(recurrency) &&
                recurrency.map((rec, index) => (
                  <option
                    key={rec.recurrency_id || index}
                    value={rec.recurrency_id}
                  >
                    {rec.recurrency_name}
                  </option>
                ))}
            </select>
          </div>
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

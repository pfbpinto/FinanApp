import React, { useEffect, useState, useCallback } from "react";
import { ChevronUpIcon } from "@heroicons/react/24/solid";

const IncomeItem = ({ income, onClose }) => {
  const [userIncomeForecast, setUserIncomeForecast] = useState(null);
  const [userIncomeActuals, setUserIncomeActuals] = useState(null);
  const [loading, setLoading] = useState(true);

  const [incomeData, setIncomeData] = useState([]);

  const [openToggleForecast, setOpenToggleForecast] = useState(true);
  const [openToggleActuals, setOpenToggleActuals] = useState(true);

  useEffect(() => {
    setIncomeData(income);
  }, [income]);

  const fetchUserIncome = useCallback(() => {
    setLoading(true);
    fetch(`/api/income-item/${incomeData.financialUserItemId}`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ itemId: incomeData.financialUserItemId }),
      credentials: "include",
    })
      .then((response) => response.json())
      .then((data) => {
        setUserIncomeForecast(data.user_financial_forecasts);
        setUserIncomeActuals(data.user_financial_actuals);
      })
      .catch((error) => console.error("Error fetching user data:", error))
      .finally(() => setLoading(false));
  }, [incomeData?.financialUserItemId]);

  useEffect(() => {
    if (incomeData?.financialUserItemId) {
      fetchUserIncome();
    }
  }, [incomeData, fetchUserIncome]);

  // Toggle the Forecast panel (open/close)
  const toggleForecastPanel = useCallback(() => {
    setOpenToggleForecast((prevState) => !prevState);
  }, []);

  // Toggle the Actuals panel (open/close)
  const toggleActualsPanel = useCallback(() => {
    setOpenToggleActuals((prevState) => !prevState);
  }, []);

  if (loading) {
    return <p>Loading...</p>;
  } else {
    return (
      <div className="container mx-100 p-6 bg-white rounded-lg shadow-lg">
        <div className="space-y-6">
          <div className="bg-white shadow-md rounded-lg p-6 mt-2">
            <div className="flex justify-between items-center mb-6">
              <h2 className="text-2xl font-semibold text-gray-800">
                {incomeData.financialUserItemName}
              </h2>
              <span
                className={`px-4 py-2 text-sm font-medium rounded-full ${
                  incomeData.isActive
                    ? "bg-green-100 text-green-600"
                    : "bg-gray-100 text-gray-600"
                }`}
              >
                {incomeData.isActive ? "Ativo" : "Inativo"}
              </span>
            </div>

            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
              <div>
                <p className="text-sm text-gray-500">Income Type</p>
                <p className="text-lg font-semibold text-gray-800">
                  {incomeData.incomeTypeName}
                </p>
              </div>

              <div>
                <p className="text-sm text-gray-500">Recurrency</p>
                <p className="text-lg font-semibold text-gray-800">
                  {incomeData.recurrencyName}
                </p>
              </div>

              <div>
                <p className="text-sm text-gray-500">Entity</p>
                <p className="text-lg font-semibold text-gray-800">
                  {incomeData.entityType}
                </p>
              </div>

              <div>
                <p className="text-sm text-gray-500">Create At</p>
                <p className="text-lg font-semibold text-gray-800">
                  {new Date(incomeData.createdAt).toLocaleDateString()}
                </p>
              </div>
            </div>

            <div className="mt-6 flex justify-end gap-4">
              <button className="px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700">
                Confirm Monthly Forecast
              </button>
            </div>
          </div>
        </div>

        <div className="space-y-6">
          <div className="bg-white shadow-md rounded-lg p-4 mt-2">
            <div
              className="flex justify-between w-full px-4 py-2 text-left text-sm font-medium text-white bg-blue-500 rounded-md hover:bg-blue-600 transition cursor-pointer"
              onClick={toggleForecastPanel}
            >
              <span>Forecast</span>
              <ChevronUpIcon
                className={`w-5 h-5 transition ${
                  openToggleForecast ? "rotate-180" : ""
                }`}
              />
            </div>

            {openToggleForecast && (
              <div className="mt-4">
                {userIncomeForecast && userIncomeForecast.length > 0 ? (
                  <div className="overflow-x-auto">
                    <table className="min-w-full table-auto border border-gray-300 rounded-lg">
                      <thead>
                        <tr className="bg-gray-100">
                          <th className="px-4 py-2 text-left">
                            Financial Item
                          </th>
                          <th className="px-4 py-2 text-left">Category</th>
                          <th className="px-4 py-2 text-left">Amount</th>
                          <th className="px-4 py-2 text-left">Begin Date</th>
                          <th className="px-4 py-2 text-left">Currency</th>
                        </tr>
                      </thead>
                      <tbody>
                        {userIncomeForecast &&
                          userIncomeForecast.map((forecast) => (
                            <tr
                              key={forecast.UserFinancialForecastID}
                              className="hover:bg-gray-50"
                            >
                              <td className="px-4 py-2">
                                {forecast.FinancialUserItemName}
                              </td>
                              <td className="px-4 py-2">
                                {forecast.UserCategoryName}
                              </td>
                              <td className="px-4 py-2">
                                {forecast.UserFinancialForecastAmount}
                              </td>
                              <td className="px-4 py-2">
                                {new Date(
                                  forecast.UserFinancialForecastBeginDate
                                ).toLocaleDateString()}
                              </td>
                              <td className="px-4 py-2">
                                {forecast.CurrencyName}
                              </td>
                            </tr>
                          ))}
                      </tbody>
                    </table>
                  </div>
                ) : (
                  <div className="w-full bg-gray-100 rounded-lg p-4 text-center text-gray-500">
                    You don't have Forecast yet
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
              onClick={toggleActualsPanel}
            >
              <span>Actuals</span>
              <ChevronUpIcon
                className={`w-5 h-5 transition ${
                  openToggleActuals ? "rotate-180" : ""
                }`}
              />
            </div>

            {openToggleActuals && (
              <div className="mt-4">
                {userIncomeActuals && userIncomeActuals.length > 0 ? (
                  <div className="overflow-x-auto">
                    <table className="min-w-full table-auto border border-gray-300 rounded-lg">
                      <thead>
                        <tr className="bg-gray-100">
                          <th className="px-4 py-2 text-left">
                            Financial Item
                          </th>
                          <th className="px-4 py-2 text-left">Category</th>
                          <th className="px-4 py-2 text-left">Amount</th>
                          <th className="px-4 py-2 text-left">Begin Date</th>
                          <th className="px-4 py-2 text-left">Currency</th>
                          <th className="px-4 py-2 text-left">Note</th>
                        </tr>
                      </thead>
                      <tbody>
                        {userIncomeActuals &&
                          userIncomeActuals.map((actual) => (
                            <tr
                              key={actual.UserFinancialActualID}
                              className="hover:bg-gray-50"
                            >
                              <td className="px-4 py-2">
                                {actual.FinancialUserItemName}
                              </td>
                              <td className="px-4 py-2">
                                {actual.UserCategoryName}
                              </td>
                              <td className="px-4 py-2">
                                {actual.UserFinancialActualAmount}
                              </td>
                              <td className="px-4 py-2">
                                {new Date(
                                  actual.UserFinancialActualtBeginDate
                                ).toLocaleDateString()}
                              </td>
                              <td className="px-4 py-2">
                                {actual.CurrencyName}
                              </td>
                              <td className="px-4 py-2">
                                {actual.Note || "N/A"}
                              </td>
                            </tr>
                          ))}
                      </tbody>
                    </table>
                  </div>
                ) : (
                  <div className="w-full bg-gray-100 rounded-lg p-4 text-center text-gray-500">
                    You don't have Actuals yet
                  </div>
                )}
              </div>
            )}
          </div>
        </div>
        <div className="flex justify-end mt-3">
          <button
            type="button"
            onClick={onClose}
            className="mr-2 px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-gray-500"
          >
            Cancel
          </button>
        </div>
      </div>
    );
  }
};

export default IncomeItem;

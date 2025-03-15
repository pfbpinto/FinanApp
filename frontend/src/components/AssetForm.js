import React, { useEffect, useState } from "react";
import { useAuth } from "../components/AuthContext";

function AssetForm({ onSubmit, asset, onClose }) {
  const [formData, setFormData] = useState(asset || { UserAssetTaxes: [] });
  const [assetTypes, setAssetTypes] = useState([]);
  const [taxes, setTaxes] = useState([]);
  const { isLoggedIn } = useAuth();

  useEffect(() => {
    if (isLoggedIn) {
      fetch("/api/asset-type", {
        method: "GET",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
      })
        .then((response) => response.json())
        .then((data) => {
          setAssetTypes(data.assetTypes || []);
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

  const handleTaxChange = (e, tax) => {
    const { checked } = e.target;

    setFormData((prevData) => {
      const updatedTaxes = checked
        ? [
            ...(prevData.UserAssetTaxes || []),
            {
              TaxID: tax.ID,
              TaxName: tax.TaxName,
              TaxPercentage: tax.TaxPercentage,
            },
          ]
        : (prevData.UserAssetTaxes || []).filter((t) => t.TaxID !== tax.ID);

      return {
        ...prevData,
        UserAssetTaxes: updatedTaxes,
      };
    });
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    onSubmit(formData);
  };

  return (
    <form onSubmit={handleSubmit}>
      <div className="mb-4">
        <label className="block text-sm font-medium text-gray-700">
          Asset Name
        </label>
        <input
          required
          type="text"
          name="AssetName"
          value={formData.AssetName || ""}
          onChange={handleChange}
          className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
        />
      </div>

      <div className="mb-4">
        <label className="block text-sm font-medium text-gray-700">
          Asset Value
        </label>
        <input
          required
          type="number"
          name="AssetValue"
          value={formData.AssetValue || ""}
          onChange={handleChange}
          className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
        />
      </div>

      <div className="mb-4">
        <label className="block text-sm font-medium text-gray-700">
          Asset Type
        </label>
        <select
          required
          name="AssetTypeID"
          value={formData.AssetTypeID || ""}
          onChange={handleChange}
          className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
        >
          <option value="">Select Asset Type</option>
          {assetTypes?.map((type) => (
            <option key={type.ID} value={type.ID}>
              {type.AssetTypeName}
            </option>
          ))}
        </select>
      </div>

      <div className="mb-4">
        <label className="block text-sm font-medium text-gray-700">
          Date of Acquisition
        </label>
        <input
          type="date"
          name="AssetAquisitionDate"
          value={
            formData.AssetAquisitionDate
              ? new Date(formData.AssetAquisitionDate)
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
            name="SharedAsset"
            checked={formData.SharedAsset || false}
            onChange={handleChange}
            className="mr-2"
          />
          <label>Shared Asset</label>
        </div>
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
                  formData.UserAssetTaxes?.some((t) => t.TaxID === tax.ID) ||
                  false
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

export default AssetForm;

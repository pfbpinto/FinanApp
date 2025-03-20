import React, { useEffect, useState } from "react";

const CategoryForm = ({ category, entityTypes, userCategory, onClose }) => {
  const [categoryData, setCategoryData] = useState([]);
  const [incomeTypeData, setincomeTypeData] = useState([]);
  const [name, setName] = useState("");
  const [cat, setCat] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  // States confirmation modal and to Tax delete

  useEffect(() => {
    setCategoryData(userCategory);
  }, [userCategory]);

  useEffect(() => {
    setincomeTypeData(entityTypes);
  }, [entityTypes]);

  const handleCategorySubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
  };

  return (
    <div>
      <div className="container">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* Category Form */}
          <div className="p-4 border rounded-lg shadow-md bg-white">
            <h2 className="text-xl font-semibold mb-4">Insert New Category</h2>
            <form onSubmit={handleCategorySubmit} className="space-y-4">
              <div>
                <label
                  htmlFor="categoryName"
                  className="block text-sm font-medium text-gray-700"
                >
                  Category Name
                </label>
                <input
                  id="categoryName"
                  type="text"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder="Insert category name"
                  name="user_category_name"
                  required
                  className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                />
              </div>

              <div className="mb-2">
                <label className="block text-sm font-medium text-gray-700">
                  Income Type
                </label>
                <select
                  required
                  value={cat}
                  name="item_type_name_id"
                  onChange={(e) => setCat(e.target.value)}
                  className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                >
                  <option value="">Select Income Category</option>
                  {incomeTypeData?.map((type) => (
                    <option
                      key={type.income_type_id}
                      value={type.income_type_id}
                    >
                      {type.income_type_name}
                    </option>
                  ))}
                </select>
              </div>

              {error && <div className="text-red-500 text-sm">{error}</div>}

              <div className="flex justify-start">
                <button
                  type="button"
                  onClick={onClose}
                  className="mr-2 px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={loading}
                  className="px-4 py-2 text-sm font-medium text-white bg-blue-500 rounded-md hover:bg-blue-600"
                >
                  {loading ? "Saving..." : "Save"}
                </button>
              </div>
            </form>
          </div>

          {/* Category Table */}

          <div
            id="cat-table-container"
            className="overflow-x-auto max-h-80 overflow-y-auto border rounded-md"
          >
            <table className="min-w-full table-auto border-collapse">
              <thead>
                <tr>
                  <th className="border px-4 py-2 text-left">ID</th>
                  <th className="border px-4 py-2 text-left">Category</th>
                </tr>
              </thead>
              <tbody>
                {categoryData.length > 0 ? (
                  categoryData.map((cat) => (
                    <tr key={cat.user_category_id} className="hover:bg-gray-50">
                      <td className="border border-gray-300 px-4 py-2">
                        {cat.user_category_id}
                      </td>
                      <td className="border border-gray-300 px-4 py-2">
                        {cat.user_category_name}
                      </td>
                    </tr>
                  ))
                ) : (
                  <tr>
                    <td colSpan="2" className="text-center py-3">
                      No categories found.
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </div>
        <br></br>
        {/* Section to display existing categories */}
        <div className="min-w-80">
          <h2 className="text-xl font-semibold mb-4">
            Current {category}s Types
          </h2>
          <div
            id="cat-table-container"
            className="overflow-x-auto max-h-80 overflow-y-auto border rounded-md"
          >
            <table className="min-w-full table-auto border-collapse">
              <thead className="sticky top-0 bg-white shadow-md">
                <tr>
                  <th className="px-4 py-2 text-left border-b">ID</th>
                  <th className="px-4 py-2 text-left border-b">Name</th>
                  <th className="px-4 py-2 text-left border-b">Description</th>
                </tr>
              </thead>
              <tbody>
                {incomeTypeData?.map((income) => (
                  <tr key={income.income_type_id}>
                    <td className="px-4 py-2 border-b">
                      {income.income_type_id}
                    </td>
                    <td className="px-4 py-2 border-b">
                      {income.income_type_name}
                    </td>
                    <td className="px-4 py-2 border-b">
                      {income.income_description}
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
    </div>
  );
};

export default CategoryForm;

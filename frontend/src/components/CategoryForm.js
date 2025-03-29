import React, { useEffect, useState } from "react";
import toastr from "toastr";

const CategoryForm = ({ userCategory, onClose }) => {
  const [categoryData, setCategoryData] = useState([]);
  const [name, setName] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  // States confirmation modal and to Tax delete
  useEffect(() => {
    setCategoryData(userCategory);
  }, [userCategory]);

  console.log(categoryData);
  const handleCategorySubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      const response = await fetch("/api/income-category", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          user_category_name: name,
        }),
      });

      if (!response.ok) {
        throw new Error("Failed to save category");
      }

      const { user_category_id } = await response.json();

      if (!user_category_id) {
        throw new Error("Invalid response from server");
      }

      // Construindo manualmente o objeto completo
      const newCategory = {
        user_category_id,
        user_category_name: name, // Pegando o nome do input
      };

      setCategoryData((prev) => [...prev, newCategory]); // Atualiza a tabela com a nova categoria
      setName(""); // Limpa o input
      toastr.success("Income Category successfully added");
    } catch (error) {
      setError(error.message);
      toastr.error("Failed to save category");
    } finally {
      setLoading(false);
    }
  };

  const handleDeleteCategory = async (categoryId) => {
    if (!window.confirm("Are you sure you want to delete this category?")) {
      return;
    }

    try {
      const response = await fetch(`/api/delete-income-category`, {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ user_category_id: categoryId }),
      });

      if (!response.ok) {
        throw new Error("Failed to delete category");
      }

      // Remove a categoria do estado local
      setCategoryData((prev) =>
        prev.filter((cat) => cat.user_category_id !== categoryId)
      );

      toastr.success("Category deleted successfully");
    } catch (error) {
      toastr.error("Failed to delete category");
    }
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
                  <th className="border px-4 py-2 text-left">Category</th>
                  <th className="border px-4 py-2 text-left"></th>
                </tr>
              </thead>
              <tbody>
                {categoryData.length > 0 ? (
                  categoryData.map((cat) => (
                    <tr key={cat.user_category_id} className="hover:bg-gray-50">
                      <td className="border border-gray-300 px-4 py-2">
                        {cat.user_category_name}
                      </td>
                      <td className="border border-gray-300 px-4 py-2">
                        <button
                          onClick={() =>
                            handleDeleteCategory(cat.user_category_id)
                          }
                          className="px-3 py-1 text-sm text-white bg-red-500 rounded-md hover:bg-red-600"
                        >
                          X
                        </button>
                      </td>
                    </tr>
                  ))
                ) : (
                  <tr>
                    <td colSpan="3" className="text-center py-3">
                      No categories found.
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </div>
        <br></br>
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

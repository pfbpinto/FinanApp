import React, { useEffect, useState } from "react";
import { useAuth } from "../components/AuthContext";
import ConfirmDeleteModal from "../components/ModalConfirm";
import toastr from "toastr";

// GroupForm component for inserting new Group
function GroupForm({ onClose, user }) {
  // Getting authentication state from context
  const { isLoggedIn } = useAuth();
  // State to store form data
  const [formData, setFormData] = useState({
    GroupName: "",
    GroupTypeID: "",
    UserID: "",
  });

  const [formGroup, setFormGroup] = useState({
    GroupID: "",
    GroupItemSelected: "",
  });

  const [formGroupInvite, setformGroupInvite] = useState({
    GroupID: "",
    InviteEmail: "",
  });

  const [userGroup, setUserGroup] = useState([]);
  const [userExpense, setUserExpense] = useState([]);
  const [userIncome, setUserIncome] = useState([]);
  const [groupType, setGroupType] = useState([]);

  // States confirmation modal and to Group delete
  const [isConfirmDeleteModalOpen, setIsConfirmDeleteModalOpen] =
    useState(false);
  const [isGroupModalOpen, setIsGroupModalOpen] = useState(false);

  const [groupToDelete, setGroupToDelete] = useState(null);
  const [deleteMessage, setDeleteMessage] = useState("");

  const [groupToManage, setgroupToManage] = useState(null);

  useEffect(() => {
    if (user) {
      setFormData((prevData) => ({ ...prevData, UserID: user }));
    }
  }, [user]);

  useEffect(() => {
    if (isLoggedIn) {
      fetchUserGroup();
    }
  }, [isLoggedIn]);

  // Fetch User Dashboard info
  const fetchUserGroup = () => {
    fetch("/api/user-group", {
      method: "GET",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
    })
      .then((response) => response.json())
      .then((data) => {
        setUserGroup(data.userGroup || []);
        setUserIncome(data.userIncome || []);
        setUserExpense(data.userExpense || []);
        setGroupType(data.groupType || []);
      })
      .catch((error) => console.error("Error fetching data:", error));
  };
  // handleChange function to handle input changes
  const handleChange = (e) => {
    const { name, value } = e.target;
    // Update the form data state with the new value from the input field
    setFormData((prevData) => ({
      ...prevData,
      [name]: value, // Check if the input is a checkbox
    }));
  };

  const handleGroupChange = (e) => {
    const { name, value } = e.target;
    setFormGroup((prevGroup) => ({
      ...prevGroup,
      [name]: value,
    }));
  };

  const handleGroupInvite = (e) => {
    const { name, value } = e.target;
    setformGroupInvite((prevGroupInv) => ({
      ...prevGroupInv,
      [name]: value,
    }));
  };

  // handleSubmit function to handle form submission (create group)
  const handleGroupSubmit = (e) => {
    // Send the data to the server
    e.preventDefault(); // Prevent page reload on form submission

    fetch(`/api/create-group`, {
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
            return Promise.reject(errorData);
          });
        }
        // If successful, return the response as JSON
        return response.json();
      })
      .then((data) => {
        // If the data is null (error occurred), skip the success handler
        console.log("API: " + data);
        if (!data) return;
        // Show a success message if the tax was created successfully
        toastr.success("Group Successfully Created.");
        // Update table with the new value
        setUserGroup((prevData) => [...prevData, data.group]);
      })
      .catch((error) => {
        // Log the error and display an error message using toastr
        console.error("Error saving Tax:", error);
        toastr.error(`Error: ${error.message}`);
      });
  };
  // handleSubmit function to handle form submission (add item to group)
  const handleGroupItemSubmit = (e) => {
    e.preventDefault(); // Prevent page reload on form submission
    fetch(`/api/create-group-item`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(formGroup),
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
        toastr.success("Item Successfully Added to the Group.");
        // Update table with the new value
        fetchUserGroup();
      })
      .catch((error) => {
        // Log the error and display an error message using toastr
        console.error("Error saving Tax:", error);
        toastr.error(`Error: ${error.message}`);
      });
  };

  // handleSubmit function to handle form submission (add item to group)
  const handleGroupInviteSubmit = (e) => {
    e.preventDefault(); // Prevent page reload on form submission
    fetch(`/api/create-group-invite`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(formGroupInvite),
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
        toastr.success("Group Invite Sent.");
        // Update table with the new value
        fetchUserGroup();
      })
      .catch((error) => {
        // Log the error and display an error message using toastr
        console.error("Error saving Tax:", error);
        toastr.error(`Error: ${error.message}`);
      });
  };

  // Handle confirmation modal for Group delete
  const handleGroupDelete = (group) => {
    setGroupToDelete(group);
    setDeleteMessage(
      `Are you sure you want to delete the Group "${group.GroupName}"?`
    );
    setIsConfirmDeleteModalOpen(true);
  };
  // Handle Group delete
  const deleteGroup = async () => {
    try {
      const response = await fetch(`/api/delete-group/${groupToDelete.ID}`, {
        method: "DELETE",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ groupID: groupToDelete.ID }),
        credentials: "include",
      });

      if (!response.ok) {
        throw new Error("Failed to delete group");
      }

      // // Update table removing the deleted group
      setUserGroup((prevDataDel) =>
        prevDataDel.filter((group) => group.ID !== groupToDelete.ID)
      );

      toastr.success("Group Successfully Deleted!");
    } catch (error) {
      console.error("Error deleting group:", error);
      toastr.error("Failed to delete group");
    }
  };

  // Handle manage modal for Groups
  const handleGroupInfo = (group) => {
    setgroupToManage(group);
    setIsGroupModalOpen(true);
  };

  console.log(userGroup);

  return (
    <div>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
        <div className="max-w-none md:max-w-md">
          <h2 className="text-xl font-semibold mb-4">Create new Group</h2>
          <form onSubmit={handleGroupSubmit}>
            <div className="grid grid-cols-1 sm:grid-cols-1 gap-2 mb-4">
              {/* Input for Group Name */}
              <div>
                <label className="block text-sm font-medium text-gray-700">
                  Group Name
                </label>
                <input
                  required
                  type="text"
                  name="GroupName"
                  onChange={handleChange} // handleChange is used here to update state
                  className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                />
              </div>

              {/* Input for Group Type */}
              <div>
                <label className="block text-sm font-medium text-gray-700">
                  Group Type
                </label>
                <select
                  required
                  name="GroupTypeID"
                  onChange={handleChange} // handleChange is also used here
                  className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                >
                  <option value="">Select Group Type</option>
                  {groupType?.map((type) => (
                    <option key={type.ID} value={type.ID}>
                      {type.GroupTypeName}
                    </option>
                  ))}
                </select>
              </div>
            </div>

            <div className="flex justify-start mt-5">
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

        {/* Section to display existing Group */}

        {userGroup.length > 0 && (
          <div>
            <div className="min-w-80">
              <h2 className="text-xl font-semibold mb-4">
                Share Item with Group
              </h2>
              <form onSubmit={handleGroupItemSubmit}>
                <div className="grid grid-cols-1 sm:grid-cols-1 gap-2 mb-4">
                  {/* Input for Group */}
                  <div>
                    <label className="block text-sm font-medium text-gray-700">
                      Select a Group
                    </label>
                    <select
                      required
                      name="GroupID"
                      onChange={handleGroupChange}
                      className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                    >
                      <option value="">Select Group</option>
                      {userGroup?.map((type) => (
                        <option key={type.ID} value={type.ID}>
                          {type.GroupName} ( {type.GroupType.GroupTypeName} )
                        </option>
                      ))}
                    </select>
                  </div>

                  {/* Input for Asset, Income and Expense */}
                  <div>
                    <label className="block text-sm font-medium text-gray-700">
                      Add to Group
                    </label>
                    <select
                      required
                      name="GroupItemSelected"
                      onChange={handleGroupChange}
                      className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                    >
                      <option value="">Select Item</option>
                      {userIncome?.map((inc) => (
                        <option key={inc.ID} value={`income_${inc.ID}`}>
                          🟢 {inc.IncomeName} {"- Income"}
                        </option>
                      ))}
                      {userExpense?.map((exp) => (
                        <option key={exp.ID} value={`expense_${exp.ID}`}>
                          🔴 {exp.ExpenditureName} {"- Expense"}
                        </option>
                      ))}
                    </select>
                  </div>
                </div>

                <div className="flex justify-start mt-5">
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
          </div>
        )}

        {userGroup.length > 0 && (
          <div>
            <div className="min-w-80">
              <h2 className="text-xl font-semibold mb-4">Group Invite</h2>
              <form onSubmit={handleGroupInviteSubmit}>
                <div className="grid grid-cols-1 sm:grid-cols-1 gap-2 mb-4">
                  {/* Input for Group */}
                  <div>
                    <label className="block text-sm font-medium text-gray-700">
                      Select a Group
                    </label>
                    <select
                      required
                      name="GroupID"
                      onChange={handleGroupInvite}
                      className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                    >
                      <option value="">Select Group</option>
                      {userGroup?.map((type) => (
                        <option key={type.ID} value={type.ID}>
                          {type.GroupName} ( {type.GroupType.GroupTypeName} )
                        </option>
                      ))}
                    </select>
                  </div>

                  {/* Input for User Email Invite */}
                  <div>
                    <label className="block text-sm font-medium text-gray-700">
                      User Email
                    </label>
                    <input
                      required
                      type="text"
                      name="InviteEmail"
                      onChange={handleGroupInvite} // handleChange is used here to update state
                      className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                    />
                  </div>
                </div>

                <div className="flex justify-start mt-5">
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
          </div>
        )}
      </div>

      <hr className="mt-5"></hr>
      {/* Section to display existing Groups */}
      <div className="min-w-80 mt-5">
        <h2 className="text-xl font-semibold mb-4">Current Groups</h2>
        <div
          id="tax-table-container"
          className="overflow-x-auto max-h-80 overflow-y-auto border rounded-md"
        >
          <table className="min-w-full table-auto border-collapse">
            <thead className="sticky top-0 bg-white shadow-md">
              <tr>
                <th></th>
                <th className="px-4 py-2 text-left border-b">Group Name</th>
                <th className="px-4 py-2 text-left border-b">Group Type</th>
                <th className="px-4 py-2 text-left border-b">Shared Assets</th>
                <th className="px-4 py-2 text-left border-b">Shared Incomes</th>
                <th className="px-4 py-2 text-left border-b">
                  Shared Espenses
                </th>
                <th className="px-4 py-2 text-left border-b">Members</th>
                <th className="px-4 py-2 text-left border-b">
                  Pending Invites
                </th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {userGroup
                ?.slice()
                .sort((a, b) => b.ID - a.ID) // Sort by ID descending
                .map((ug) => (
                  <tr key={ug.ID}>
                    <td className="px-4 py-2 border-b space-x-2">
                      <button
                        className="px-3 py-1 text-xs font-medium text-white bg-blue-500 rounded-md hover:bg-black-600"
                        onClick={() => handleGroupInfo(ug)}
                      >
                        ?
                      </button>
                    </td>

                    <td className="px-4 py-2 border-b">{ug.GroupName}</td>
                    <td className="px-4 py-2 border-b">
                      {ug.GroupType.GroupTypeName}
                    </td>
                    <td className="px-4 py-2 border-b">Assets</td>
                    <td className="px-4 py-2 border-b">
                      {Array.isArray(ug.UserGroupIncomes) &&
                      ug.UserGroupIncomes.length > 0 ? (
                        <div className="flex flex-wrap gap-2">
                          {ug.UserGroupIncomes.map((userGroupInc, idx) => (
                            <div
                              key={idx}
                              className="p-1 bg-blue-100 rounded-lg shadow-sm text-sm font-medium text-blue-800"
                            >
                              {userGroupInc?.UserIncome?.IncomeName ||
                                "Unnamed Income"}
                            </div>
                          ))}
                        </div>
                      ) : (
                        <span className="text-gray-500">No Incomes</span>
                      )}
                    </td>

                    <td className="px-4 py-2 border-b">Expenses</td>
                    <td className="px-4 py-2 border-b">
                      {Array.isArray(ug.GroupMembers) &&
                      ug.GroupMembers.length > 0 ? (
                        <div className="flex flex-wrap gap-2">
                          {ug.GroupMembers.length}
                        </div>
                      ) : (
                        <span className="text-gray-500">0</span>
                      )}
                    </td>
                    <td className="px-4 py-2 border-b">
                      {Array.isArray(ug.GroupInvites) &&
                      ug.GroupInvites.length > 0 ? (
                        <div className="flex flex-wrap gap-2">
                          {ug.GroupInvites.length}
                        </div>
                      ) : (
                        <span className="text-gray-500">0</span>
                      )}
                    </td>

                    <td className="px-4 py-2 border-b space-x-2">
                      <button
                        className="px-3 py-1 text-xs font-medium text-white bg-red-500 rounded-md hover:bg-red-600"
                        onClick={() => handleGroupDelete(ug)}
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
          onConfirm={deleteGroup}
        />
      </div>
      <div>
        <ConfirmDeleteModal
          isOpen={isGroupModalOpen}
          onClose={() => setIsGroupModalOpen(false)}
          message={"Group Management"}
          group={groupToManage}
        />
      </div>
    </div>
  );
}

export default GroupForm;

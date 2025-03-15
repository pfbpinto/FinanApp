import React from "react";

function Modal({ onClose, title, onDelete, item, children }) {
  return (
    <div className="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50 z-50">
      <div className="bg-white rounded-lg shadow-lg w-full max-w-md p-6">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-semibold">{title}</h2>
          <button
            onClick={onClose}
            className="text-gray-500 hover:text-gray-700 focus:outline-none"
          >
            ×
          </button>
        </div>
        {children}
        {onDelete && item && (
          <div className="flex justify-end mt-4">
            <button
              onClick={onClose}
              className="mr-2 bg-gray-300 px-4 py-2 rounded"
            >
              Cancel
            </button>
            <button
              onClick={() => onDelete(item)} // Passando o item para a função onDelete
              className="bg-red-600 text-white px-4 py-2 rounded"
            >
              Delete
            </button>
          </div>
        )}
      </div>
    </div>
  );
}

export default Modal;

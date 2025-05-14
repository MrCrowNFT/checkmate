import { useState } from "react";
import { useDeployments } from "../hooks";
import type { platformCredentialInput } from "../types";

export const CredentialsButton = () => {
  const { credentials, isLoading, error, newCredentials } = useDeployments();
  const [showCredentials, setShowCredentials] = useState(false);
  const [newCredential, setNewCredential] = useState<platformCredentialInput>({
    platform: "",
    name: "",
    apiKey: "",
  });

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setNewCredential((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleAddCredential = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await newCredentials(newCredential);
      // Reset form after successful addition
      setNewCredential({
        platform: "",
        name: "",
        apiKey: "",
      });
    } catch (err) {
      console.error("Failed to add credential:", err);
    }
  };

  return (
    <div className="w-full max-w-md mx-auto font-sans">
      {/* Primary button using Checkmate style */}
      <button
        onClick={() => setShowCredentials(!showCredentials)}
        className="w-full bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded transition-colors duration-200 ease-in-out mb-4 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-opacity-50"
      >
        {showCredentials ? "Hide Credentials" : "Show Credentials"}
      </button>

      {showCredentials && (
        <div className="mt-4 bg-white dark:bg-gray-900 rounded shadow-sm border border-gray-200 dark:border-gray-700 p-6">
          <h2 className="text-xl font-bold mb-4 text-gray-900 dark:text-white">
            Your Credentials
          </h2>

          {isLoading && (
            <p className="text-gray-500 dark:text-gray-400">
              Loading credentials...
            </p>
          )}

          {error && (
            <div className="mb-4 p-3 bg-red-50 dark:bg-red-900/20 border-l-4 border-red-500 dark:border-red-400 text-red-700 dark:text-red-400">
              {error}
            </div>
          )}

          {!isLoading && credentials.length === 0 ? (
            <p className="text-gray-500 dark:text-gray-400 mb-4">
              No credentials found. Add a new one below.
            </p>
          ) : (
            <ul className="mb-6 divide-y divide-gray-200 dark:divide-gray-700">
              {credentials.map((cred) => (
                <li key={cred.id} className="py-3">
                  <div className="flex justify-between">
                    <div>
                      <p className="font-medium text-gray-900 dark:text-white">
                        {cred.name}
                      </p>
                      <div className="flex items-center mt-1">
                        <span className="inline-block w-2 h-2 rounded-full bg-green-500 dark:bg-green-400 mr-2"></span>
                        <p className="text-sm text-gray-500 dark:text-gray-400">
                          {cred.platform}
                        </p>
                      </div>
                    </div>
                    <div className="text-xs text-gray-400 dark:text-gray-500">
                      {new Date(cred.created_at).toLocaleDateString()}
                    </div>
                  </div>
                </li>
              ))}
            </ul>
          )}

          <div className="pt-4 border-t border-gray-200 dark:border-gray-700">
            <h3 className="text-lg font-medium mb-3 text-gray-900 dark:text-white">
              Add New Credential
            </h3>
            <form onSubmit={handleAddCredential}>
              <div className="mb-3">
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Platform
                </label>
                <input
                  type="text"
                  name="platform"
                  value={newCredential.platform}
                  onChange={handleInputChange}
                  className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded shadow-sm 
                           bg-white dark:bg-gray-800 text-gray-800 dark:text-white
                           focus:outline-none focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400 focus:border-blue-500"
                  required
                />
              </div>

              <div className="mb-3">
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Name
                </label>
                <input
                  type="text"
                  name="name"
                  value={newCredential.name}
                  onChange={handleInputChange}
                  className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded shadow-sm 
                           bg-white dark:bg-gray-800 text-gray-800 dark:text-white
                           focus:outline-none focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400 focus:border-blue-500"
                  required
                />
              </div>

              <div className="mb-4">
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  API Key
                </label>
                <input
                  type="password"
                  name="apiKey"
                  value={newCredential.apiKey}
                  onChange={handleInputChange}
                  className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded shadow-sm 
                           bg-white dark:bg-gray-800 text-gray-800 dark:text-white
                           focus:outline-none focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400 focus:border-blue-500"
                  required
                />
              </div>

              <button
                type="submit"
                disabled={isLoading}
                className="w-full bg-green-500 hover:bg-green-600 dark:bg-green-600 dark:hover:bg-green-700 
                         text-white font-medium py-2 px-4 rounded transition-colors duration-200 ease-in-out
                         disabled:bg-gray-300 dark:disabled:bg-gray-700 disabled:text-gray-500 dark:disabled:text-gray-400
                         focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-opacity-50"
              >
                {isLoading ? "Adding..." : "Add Credential"}
              </button>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

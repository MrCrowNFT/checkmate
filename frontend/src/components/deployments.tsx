import { useEffect, useState } from "react";
import { useAuth } from "../contexts/auth-context";
import api from "../api/api";
import type { User, platformCredentialInput } from "../types";
import { useDeployments, useThemeStore } from "../hooks";
import DeploymentCard from "./deployment-card";
import { motion, AnimatePresence } from "framer-motion";
import { Sun, Moon } from "lucide-react";

const DarkModeToggle = () => {
  const { darkMode, toggleDarkMode } = useThemeStore();

  const toggleVariants = {
    light: {
      backgroundColor: "#D1D5DB",
      boxShadow: "0 2px 8px rgba(0, 0, 0, 0.1)",
    },
    dark: {
      backgroundColor: "#374151",
      boxShadow: "0 2px 8px rgba(0, 0, 0, 0.3)",
    },
  };

  const springConfig = {
    type: "spring",
    stiffness: 700,
    damping: 30,
  };

  return (
    <motion.button
      onClick={toggleDarkMode}
      className="relative flex items-center w-16 h-8 rounded-full p-1 cursor-pointer"
      variants={toggleVariants}
      animate={darkMode ? "dark" : "light"}
      initial={false}
      transition={{ duration: 0.2 }}
      whileHover={{ scale: 1.05 }}
      whileTap={{ scale: 0.95 }}
    >
      <AnimatePresence mode="wait">
        <motion.div
          className="absolute left-1"
          key={darkMode ? "moon" : "sun"}
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          exit={{ opacity: 0, y: 20 }}
          transition={{ duration: 0.2 }}
        >
          {darkMode ? (
            <Moon className="w-5 h-5 text-gray-100" />
          ) : (
            <Sun className="w-5 h-5 text-yellow-500" />
          )}
        </motion.div>
      </AnimatePresence>

      <motion.div
        className="w-6 h-6 rounded-full bg-white dark:bg-gray-800"
        animate={{
          x: darkMode ? 32 : 0,
          backgroundColor: darkMode ? "#1F2937" : "#FFFFFF",
        }}
        transition={springConfig}
        style={{
          boxShadow: "0 1px 3px rgba(0, 0, 0, 0.1)",
        }}
      />

      {/* Focus ring */}
      <div className="absolute inset-0 rounded-full ring-0 focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 focus:outline-none" />
    </motion.button>
  );
};

const CredentialsSection = () => {
  const { credentials, isLoading, error, newCredentials, getCredentials } =
    useDeployments();
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
      // Refresh credentials list
      await getCredentials();
    } catch (err) {
      console.error("Failed to add credential:", err);
    }
  };

  return (
    <div className="w-full max-w-md mx-auto font-sans mb-8">
      {/* Primary button using Checkmate style */}
      <button
        onClick={() => setShowCredentials(!showCredentials)}
        className="w-full bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded transition-colors duration-200 ease-in-out mb-4 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-opacity-50 dark:bg-blue-700 dark:hover:bg-blue-800"
      >
        {showCredentials ? "Hide Credentials" : "Show Credentials"}
      </button>

      <AnimatePresence>
        {showCredentials && (
          <motion.div
            initial={{ opacity: 0, height: 0 }}
            animate={{ opacity: 1, height: "auto" }}
            exit={{ opacity: 0, height: 0 }}
            transition={{ duration: 0.3 }}
            className="mt-4 bg-white dark:bg-gray-800 rounded shadow-sm border border-gray-200 dark:border-gray-700 p-6 overflow-hidden"
          >
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
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
};

const Deployments = () => {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string>("");
  const { currentUser } = useAuth();
  const { darkMode } = useThemeStore();

  //zustand store
  const {
    deployments,
    isLoading: deploymentsLoading,
    error: deploymentsError,
    getCredentials,
    startPolling,
    stopPolling,
  } = useDeployments();

  //todo probably want to hide this logic somewhere else
  useEffect(() => {
    const fetchUserProfile = async () => {
      if (!currentUser) return;

      try {
        await new Promise((resolve) => setTimeout(resolve, 500));

        console.log("Making API request with user:", currentUser.email);

        // call backend
        const response = await api.get("/");
        console.log("API response received:", response.data);
        setUser(response.data);
      } catch (err: any) {
        console.error("Error fetching user profile:", err);
        if (err.response) {
          console.error("Error response data:", err.response.data);
          console.error("Error response status:", err.response.status);
          setError(
            `Failed to load user profile (Status: ${err.response.status})`
          );
        } else if (err.request) {
          // request made but no response received
          console.error("No response received from server");
          setError("No response from server. Please check your connection.");
        } else {
          setError(`Request error: ${err.message}`);
        }
      } finally {
        setLoading(false);
      }
    };

    fetchUserProfile();
  }, [currentUser]);

  useEffect(() => {
    if (currentUser) {
      // fetch credentials, which will trigger getDeployments
      getCredentials().catch((error) => {
        console.error("Error fetching credentials:", error);
      });
      // start polling for deployments every 30 seconds
      startPolling();
      // clean up on component unmount
      return () => {
        stopPolling();
      };
    }
  }, [currentUser, getCredentials, startPolling, stopPolling]);

  if (loading)
    return (
      <div className="p-8 text-center dark:text-white">
        Loading user profile...
      </div>
    );
  if (!currentUser)
    return (
      <div className="p-8 text-center dark:text-white">
        Please log in to view your profile
      </div>
    ); //todo redirect to login
  if (error)
    return <div className="p-8 text-center text-red-500">Error: {error}</div>;

  return (
    <div className={`container mx-auto px-4 py-8 ${darkMode ? "dark" : ""}`}>
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-2xl font-bold dark:text-white">
          Welcome, {user?.displayName || user?.email || "User"}!
        </h1>
        <DarkModeToggle />
      </div>

      {/* Add Credentials Section */}
      <CredentialsSection />

      {/* Deployment status */}
      {deploymentsLoading && (
        <div className="text-center py-4 dark:text-gray-300">
          Loading deployments...
        </div>
      )}
      {deploymentsError && (
        <div className="text-red-500 py-4">
          Error loading deployments: {deploymentsError}
        </div>
      )}

      {/* Deployments grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {deployments.map((deployment) => (
          <DeploymentCard
            key={deployment.id}
            deployment={deployment}
            onClick={() => {}}
          />
        ))}
      </div>

      {/* Show message when no deployments */}
      {!deploymentsLoading && deployments.length === 0 && (
        <div className="text-center py-8 text-gray-500 dark:text-gray-400">
          No deployments found. Add credentials to start deploying your
          applications.
        </div>
      )}
    </div>
  );
};

export default Deployments;

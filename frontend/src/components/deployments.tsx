import { useEffect, useState } from "react";
import { useAuth } from "../contexts/auth-context";
import api from "../api/api";
import type { User } from "../types";
import { useDeployments } from "../hooks";
import DeploymentCard from "./deployment-card";

const Deployments = () => {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string>("");
  const { currentUser } = useAuth();

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

  if (loading) return <div>Loading user profile...</div>;
  if (!currentUser) return <div>Please log in to view your profile</div>; //todo redirect to login
  if (error) return <div>Error: {error}</div>;

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-2xl font-bold mb-6">
        Welcome, {user?.displayName || user?.email || "User"}!
      </h1>

      {/* Deployment status */}
      {deploymentsLoading && (
        <div className="text-center py-4">Loading deployments...</div>
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
        <div className="text-center py-8 text-gray-500">
          No deployments found. Add credentials to start deploying your
          applications.
        </div>
      )}
    </div>
  );
};

export default Deployments;

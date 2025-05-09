import { useEffect, useState } from "react";
import { useAuth } from "../contexts/auth-context";
import api from "../api/api";

const Deployments = () => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const { currentUser } = useAuth();

  //todo Fix the api call to get user
  //todo Add typing of user

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
      } catch (err) {
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

  if (loading) return <div>Loading user profile...</div>;
  if (!currentUser) return <div>Please log in to view your profile</div>; //todo redirect to login
  if (error) return <div>Error: {error}</div>;

  return (
    <div>
      <h1>Welcome, {user?.displayName || user?.email || "User"}!</h1>
      {/* here need to deconstruct the deployents into deployment cards */}
    </div>
  );
};

export default Deployments;

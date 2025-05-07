import { useEffect, useState } from "react";
import { useAuth } from "../contexts/auth-context";
import api from "../api/api";

const Deployments = () => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const { currentUser } = useAuth();

  useEffect(() => {
    const fetchUserProfile = async () => {
      if (!currentUser) return;

      try {
        // call backend
        const response = await api.get("/auth");
        setUser(response.data);
      } catch (err) {
        console.error("Error fetching user profile:", err);
        setError("Failed to load user profile");
      } finally {
        setLoading(false);
      }
    };

    fetchUserProfile();
  }, [currentUser]);

  if (loading) return <div>Loading user profile...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <div>
      <h1>Welcome, {user?.displayName || user?.email || "User"}!</h1>
      {/* here need to deconstruct the deployents into deployment cards */}
    </div>
  );
};

export default Deployments;

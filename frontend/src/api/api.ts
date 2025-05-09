import axios from "axios";
import { auth } from "../firebase";

const API_URL = import.meta.env.VITE_API_URL;

const api = axios.create({
  baseURL: API_URL,
});

// automatically add auth token to all requests
api.interceptors.request.use(
  async (config) => {
    try {
      const user = auth.currentUser;

      // proceed only if user
      if (user) {
        // token refresh if needed with getIdToken(true)
        const token = await user.getIdToken(true);

        //  "Bearer " format is consistent with backend
        config.headers.Authorization = `Bearer ${token}`;

        // debugging
        console.log("Token attached to request");
      } else {
        console.warn("No current user found when making API request");
      }

      return config;
    } catch (error) {
      console.error("Error in request interceptor:", error);
      return Promise.reject(error);
    }
  },
  (error) => {
    return Promise.reject(error);
  }
);

export default api;

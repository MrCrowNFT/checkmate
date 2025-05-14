import { create } from "zustand";
import { persist } from "zustand/middleware";
import type {
  Deployment,
  platformCredentialInput,
  safeCredential,
} from "../types";
import {
  getCredentials,
  newCredential,
  updateCredential,
  deleteCredential,
} from "../api";
import { getDeployments as fetchDeployments } from "../api";
import axios from "axios";

type DeploymentsList = {
  credentials: safeCredential[];
  deployments: Deployment[];
  isLoading: boolean;
  error: string | null;
  pollingInterval: number | null; //to track intervals

  //credentials
  getCredentials: () => Promise<safeCredential[]>;
  newCredentials: (
    credential: platformCredentialInput
  ) => Promise<safeCredential>;
  updateCredentials: (id: number, updateCred: platformCredentialInput) => void;
  deleteCredentials: (id: number) => void;

  //deployments
  getDeployments: () => Promise<Deployment[]>;
  startPolling: () => void;
  stopPolling: () => void;
};

//* naming is not great, this is for deployment and credentials
export const useDeployments = create<DeploymentsList>()(
  persist(
    (set, get) => ({
      credentials: [],
      deployments: [],
      isLoading: false,
      error: null,
      pollingInterval: null,

      //credentials

      //this one will be the first method called after loging in
      getCredentials: async () => {
        set({ isLoading: true, error: null });
        try {
          const response = await getCredentials();
          set({
            credentials: response,
            isLoading: false,
          });
          //fetch deployments after getting the creds
          get().getDeployments();
          return response;
        } catch (error) {
          const errorMessage = axios.isAxiosError(error)
            ? error.response?.data?.message
            : "Error fetching credentials";
          set({ isLoading: false, error: errorMessage });
          throw error;
        }
      },

      newCredentials: async (credential) => {
        set({ isLoading: true, error: null });
        try {
          const response = await newCredential(credential);
          set((state) => ({
            credentials: [...state.credentials, response],
            isLoading: false,
          }));

          //refetch deployments after adding new creds
          get().getDeployments();

          return response;
        } catch (error) {
          const errorMessage = axios.isAxiosError(error)
            ? error.response?.data?.message
            : "Error creating credentials";
          set({ isLoading: false, error: errorMessage });
          throw error;
        }
      },

      updateCredentials: async (id, updateCred) => {
        set({ isLoading: true, error: null });
        try {
          // Optimistic update
          set((state) => ({
            credentials: state.credentials.map((cred) =>
              cred.id === id ? { ...cred, ...updateCred } : cred
            ),
            isLoading: false,
          }));

          // API call
          await updateCredential(id, updateCred);

          //refetch deployments after updating creds
          get().getDeployments();
        } catch (error) {
          // rollback on error
          const errorMessage = axios.isAxiosError(error)
            ? error.response?.data?.message
            : "Error updating credentials";
          set({ isLoading: false, error: errorMessage });

          // refresh credentials to restore correct state
          get().getCredentials();
          throw error;
        }
      },

      deleteCredentials: async (id) => {
        set({ isLoading: true, error: null });
        try {
          // Optimistic delete
          set((state) => ({
            credentials: state.credentials.filter((cred) => cred.id !== id),
            isLoading: false,
          }));

          await deleteCredential(id);

          //refetch deployments after deliting creds
          get().getDeployments();
        } catch (error) {
          // rollback on error
          const errorMessage = axios.isAxiosError(error)
            ? error.response?.data?.message
            : "Error deleting credentials";
          set({ isLoading: false, error: errorMessage });

          // refresh credentials to restore correct state
          get().getCredentials();
          throw error;
        }
      },

      //deployments
      getDeployments: async () => {
        set({ isLoading: true, error: null });
        try {
          const deployments = await fetchDeployments();
          set({
            deployments,
            isLoading: false,
          });
          return deployments;
        } catch (error) {
          const errorMessage = axios.isAxiosError(error)
            ? error.response?.data?.message
            : "Error fetching deployments";
          set({ isLoading: false, error: errorMessage });
          throw error;
        }
      },
      startPolling: () => {
        if (get().pollingInterval !== null) {
          window.clearInterval(get().pollingInterval as number);
        }
        const intervalId = window.setInterval(() => {
          get().getDeployments();
        }, 30000); // 30 seconds

        set({ pollingInterval: intervalId });
      },
      stopPolling: () => {
        if (get().pollingInterval !== null) {
          window.clearInterval(get().pollingInterval as number);
          set({ pollingInterval: null });
        }
      },
    }),
    {
      name: "deployments",
      partialize: (state) => ({
        credentials: state.credentials,
        deployments: state.deployments,
        isLoading: state.isLoading,
        error: state.error,
      }),
    }
  )
);

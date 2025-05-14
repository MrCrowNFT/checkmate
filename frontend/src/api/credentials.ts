import type { platformCredentialInput, safeCredential } from "../types";
import api from "./api";

//*note i will use this on the zustand store so try catch block will be applied there
export const getCredentials = async (): Promise<safeCredential[]> => {
  const res = await api.get("/credentials");
  return res.data.credentials;
};

export const newCredential = async (
  credential: platformCredentialInput
): Promise<safeCredential> => {
  const res = await api.post("/credentials/new", credential);
  return res.data;
};

//i will just make an optimistic update on zustand
export const updateCredential = async (
  id: number,
  updateCred: platformCredentialInput
) => {
  const res = await api.put(`/credentials/update?id=${id}`, updateCred);
  return res.data;
};

export const deleteCredential = async (id: number) => {
  const res = await api.delete(`/credentials/delete?id=${id}`);
  return res.data;
};

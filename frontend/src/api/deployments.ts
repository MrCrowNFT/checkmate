import api from "./api";
import type { deployment } from "../types";

export const getDeployments = async (): Promise<deployment[]> => {
  const res = await api.get("/deployments");
  return res.data;
};

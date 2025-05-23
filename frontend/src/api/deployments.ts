import api from "./api";
import type { Deployment } from "../types";

export const getDeployments = async (): Promise<Deployment[]> => {
  const res = await api.get("/deployments");
  console.log(res);
  return res.data.deployments;
};

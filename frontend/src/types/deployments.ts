import type { safeCredential } from "./credentials";

export type DeploymentStatus =
  | "live"
  | "deploying"
  | "canceled"
  | "failed"
  | "unknown";

export interface Deployment {
  id: string;
  platformCredentialID: number;
  name: string;
  status: DeploymentStatus;
  url: string;
  lastDeployedAt: string | null;
  branch: string;
  serviceType: string;
  framework: string;
  lastUpdatedAt: string;
  metadata: Record<string, any>;
}

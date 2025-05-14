export type deploymentStatus =
  | "live"
  | "deploying"
  | "canceled"
  | "failed"
  | "unknown";

export interface deployment {
  id: string;
  platformCredentialID: number;
  name: string;
  status: deploymentStatus;
  url: string;
  lastDeployedAt: string | null;
  branch: string;
  serviceType: string;
  framework: string;
  lastUpdatedAt: string;
  metadata: Record<string, any>;
}

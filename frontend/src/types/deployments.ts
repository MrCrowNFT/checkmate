export type DeploymentStatus =
  | "live"
  | "deploying"
  | "canceled"
  | "failed"
  | "unknown";

export interface PlatformCredential {
  id: string;
  name?: string;
  type?: string;
}

export interface Deployment {
  id: string;
  platformCredential: PlatformCredential | null;
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

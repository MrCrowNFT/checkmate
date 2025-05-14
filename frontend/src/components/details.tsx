import { Clock, ExternalLink, Server, Code, GitBranch } from "lucide-react";
import StatusBadge from "./status-badge";
import type { Deployment, DeploymentStatus } from "../types";
import { formatDate } from "../helpers";

const DeploymentDetails = ({
  deployment,
  onBack,
}: {
  deployment: Deployment;
  onBack: () => void;
}) => {
  const getBorderColor = (status: DeploymentStatus) => {
    switch (status) {
      case "live":
        return "border-green-500";
      case "deploying":
        return "border-blue-500";
      case "canceled":
        return "border-gray-400";
      case "failed":
        return "border-red-500";
      default:
        return "border-yellow-500";
    }
  };

  return (
    <div
      className={`border-l-4 ${getBorderColor(
        deployment.status
      )} rounded bg-white dark:bg-gray-900 shadow-sm dark:shadow-md p-4 font-['Inter']`}
    >
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-xl font-medium text-gray-900 dark:text-white">
          {deployment.name}
        </h2>
        <StatusBadge status={deployment.status} />
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="space-y-4">
          <div>
            <p className="text-xs text-gray-500 dark:text-gray-400 mb-1">
              Deployment ID
            </p>
            <p className="font-mono text-sm bg-gray-100 dark:bg-gray-800 p-2 rounded">
              {deployment.id}
            </p>
          </div>

          <div>
            <p className="text-xs text-gray-500 dark:text-gray-400 mb-1">
              Service Type
            </p>
            <p className="flex items-center gap-2 text-gray-900 dark:text-white">
              <Server size={16} />
              {deployment.serviceType}
            </p>
          </div>

          <div>
            <p className="text-xs text-gray-500 dark:text-gray-400 mb-1">
              Framework
            </p>
            <p className="flex items-center gap-2 text-gray-900 dark:text-white">
              <Code size={16} />
              {deployment.framework}
            </p>
          </div>

          <div>
            <p className="text-xs text-gray-500 dark:text-gray-400 mb-1">
              Branch
            </p>
            <p className="flex items-center gap-2 text-gray-900 dark:text-white">
              <GitBranch size={16} />
              {deployment.branch}
            </p>
          </div>

          <div>
            <p className="text-xs text-gray-500 dark:text-gray-400 mb-1">URL</p>
            <a
              href={deployment.url}
              target="_blank"
              rel="noopener noreferrer"
              className="text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300 flex items-center gap-2"
            >
              {deployment.url}
              <ExternalLink size={16} />
            </a>
          </div>
        </div>

        <div className="space-y-4">
          <div>
            <p className="text-xs text-gray-500 dark:text-gray-400 mb-1">
              Platform Credential ID
            </p>
            <p className="text-gray-900 dark:text-white">
              {deployment.platformCredentialID}
            </p>
          </div>

          <div>
            <p className="text-xs text-gray-500 dark:text-gray-400 mb-1">
              Last Deployed At
            </p>
            <p className="flex items-center gap-2 text-gray-900 dark:text-white">
              <Clock size={16} />
              {formatDate(deployment.lastDeployedAt)}
            </p>
          </div>

          <div>
            <p className="text-xs text-gray-500 dark:text-gray-400 mb-1">
              Last Updated At
            </p>
            <p className="flex items-center gap-2 text-gray-900 dark:text-white">
              <Clock size={16} />
              {formatDate(deployment.lastUpdatedAt)}
            </p>
          </div>

          <div>
            <p className="text-xs text-gray-500 dark:text-gray-400 mb-1">
              Metadata
            </p>
            <div className="bg-gray-100 dark:bg-gray-800 p-2 rounded font-mono text-xs overflow-auto max-h-32">
              <pre className="text-gray-900 dark:text-white">
                {JSON.stringify(deployment.metadata, null, 2)}
              </pre>
            </div>
          </div>
        </div>
      </div>

      <div className="mt-8 pt-4 border-t border-gray-200 dark:border-gray-700">
        <button
          onClick={onBack}
          className="bg-transparent hover:bg-gray-100 dark:hover:bg-gray-800 px-4 py-2 rounded text-gray-600 dark:text-gray-400 border border-gray-400 dark:border-gray-600 transition-colors"
        >
          Back to Card View
        </button>
      </div>
    </div>
  );
};

export default DeploymentDetails;

import {
  Clock,
  ExternalLink,
  Info,
  Server,
  Code,
  GitBranch,
} from "lucide-react";
import type { Deployment } from "../types";
import StatusBadge from "./status-badge";
import { formatDate } from "../helpers";

//todo on click should redirect to details with the id of the project
// deployment Card Component
const DeploymentCard = ({
  deployment,
  onClick,
}: {
  deployment: Deployment;
  onClick: () => void;
}) => {
  return (
    <div
      className="border border-gray-200 dark:border-gray-700 rounded bg-white dark:bg-gray-900 shadow-sm p-4 hover:shadow transition-shadow duration-200 cursor-pointer"
      onClick={onClick}
    >
      <div className="flex justify-between items-start mb-4">
        <h3 className="font-medium text-base text-gray-900 dark:text-white truncate">
          {deployment.name}
        </h3>
        <StatusBadge status={deployment.status} />
      </div>

      <div className="space-y-2 text-sm text-gray-600 dark:text-gray-300">
        <div className="flex items-center gap-2">
          <Server size={16} className="text-gray-500 dark:text-gray-400" />
          <span>{deployment.serviceType}</span>
        </div>

        <div className="flex items-center gap-2">
          <Code size={16} className="text-gray-500 dark:text-gray-400" />
          <span>{deployment.framework}</span>
        </div>

        <div className="flex items-center gap-2">
          <GitBranch size={16} className="text-gray-500 dark:text-gray-400" />
          <span>{deployment.branch}</span>
        </div>

        <div className="flex items-center gap-2">
          <Clock size={16} className="text-gray-500 dark:text-gray-400" />
          <span>Last updated: {formatDate(deployment.lastUpdatedAt)}</span>
        </div>
      </div>

      <div className="mt-4 pt-3 border-t border-gray-200 dark:border-gray-700 flex justify-between items-center">
        <button
          className="text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300 flex items-center gap-1 text-sm font-medium"
          onClick={(e) => {
            e.stopPropagation();
            window.open(deployment.url, "_blank");
          }}
        >
          <ExternalLink size={14} />
          Visit Site
        </button>

        <button
          className="text-gray-600 dark:text-gray-400 hover:text-gray-800 dark:hover:text-gray-300 flex items-center gap-1 text-sm font-medium"
          onClick={(e) => {
            e.stopPropagation();
            onClick();
          }}
        >
          <Info size={14} />
          Details
        </button>
      </div>
    </div>
  );
};

export default DeploymentCard;

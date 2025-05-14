import {
  Clock,
  ExternalLink,
  Info,
  Server,
  Code,
  GitBranch,
} from "lucide-react";
import type { Deployment, DeploymentStatus } from "../types";

// Helper function to format dates
const formatDate = (dateString: string | null) => {
  if (!dateString) return "Never";
  const date = new Date(dateString);
  return date.toLocaleString();
};

// Status Badge Component
const StatusBadge = ({ status }: { status: DeploymentStatus }) => {
  const getStatusStyles = () => {
    switch (status) {
      case "live":
        return "bg-green-100 text-green-800 border-green-300";
      case "deploying":
        return "bg-blue-100 text-blue-800 border-blue-300";
      case "canceled":
        return "bg-gray-100 text-gray-800 border-gray-300";
      case "failed":
        return "bg-red-100 text-red-800 border-red-300";
      default:
        return "bg-yellow-100 text-yellow-800 border-yellow-300";
    }
  };

  return (
    <span
      className={`px-2 py-1 text-xs font-medium rounded-full border ${getStatusStyles()}`}
    >
      {status.charAt(0).toUpperCase() + status.slice(1)}
    </span>
  );
};

// Deployment Card Component
const DeploymentCard = ({
  deployment,
  onClick,
}: {
  deployment: Deployment;
  onClick: () => void;
}) => {
  return (
    <div
      className="border rounded-lg shadow-sm p-4 hover:shadow-md transition-shadow cursor-pointer"
      onClick={onClick}
    >
      <div className="flex justify-between items-start mb-3">
        <h3 className="font-medium text-lg truncate">{deployment.name}</h3>
        <StatusBadge status={deployment.status} />
      </div>

      <div className="space-y-2 text-sm text-gray-600">
        <div className="flex items-center gap-2">
          <Server size={16} />
          <span>{deployment.serviceType}</span>
        </div>

        <div className="flex items-center gap-2">
          <Code size={16} />
          <span>{deployment.framework}</span>
        </div>

        <div className="flex items-center gap-2">
          <GitBranch size={16} />
          <span>{deployment.branch}</span>
        </div>

        <div className="flex items-center gap-2">
          <Clock size={16} />
          <span>Last updated: {formatDate(deployment.lastUpdatedAt)}</span>
        </div>
      </div>

      <div className="mt-3 pt-3 border-t flex justify-between items-center">
        <button
          className="text-blue-600 hover:text-blue-800 flex items-center gap-1 text-sm"
          onClick={(e) => {
            e.stopPropagation();
            window.open(deployment.url, "_blank");
          }}
        >
          <ExternalLink size={14} />
          Visit Site
        </button>

        <button
          className="text-gray-600 hover:text-gray-800 flex items-center gap-1 text-sm"
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

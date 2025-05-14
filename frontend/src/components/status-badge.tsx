import type { DeploymentStatus } from "../types";

const StatusBadge = ({ status }: { status: DeploymentStatus }) => {
  const getStatusStyles = () => {
    switch (status) {
      case "live":
        return "bg-green-100 text-green-500 dark:bg-opacity-10 dark:text-green-300";
      case "deploying":
        return "bg-blue-100 text-blue-600 dark:bg-opacity-10 dark:text-blue-300";
      case "canceled":
        return "bg-gray-100 text-gray-600 dark:bg-opacity-10 dark:text-gray-400";
      case "failed":
        return "bg-red-100 text-red-500 dark:bg-opacity-10 dark:text-red-300";
      default:
        return "bg-yellow-100 text-yellow-600 dark:bg-opacity-10 dark:text-yellow-300";
    }
  };

  return (
    <span
      className={`px-2 py-1 text-xs font-medium rounded-full ${getStatusStyles()}`}
    >
      {status.charAt(0).toUpperCase() + status.slice(1)}
    </span>
  );
};

export default StatusBadge;

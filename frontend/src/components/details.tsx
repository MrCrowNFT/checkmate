const DeploymentDetails = ({
  deployment,
  onBack,
}: {
  deployment: Deployment;
  onBack: () => void;
}) => {
  const getStatusColor = () => {
    switch (deployment.status) {
      case "live":
        return "border-green-500";
      case "deploying":
        return "border-blue-500";
      case "canceled":
        return "border-gray-500";
      case "failed":
        return "border-red-500";
      default:
        return "border-yellow-500";
    }
  };

  return (
    <div
      className={`border-l-4 ${getStatusColor()} rounded-lg shadow-md p-6 bg-white`}
    >
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-2xl font-semibold">{deployment.name}</h2>
        <StatusBadge status={deployment.status} />
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div className="space-y-4">
          <div>
            <p className="text-sm text-gray-500 mb-1">Deployment ID</p>
            <p className="font-mono text-sm bg-gray-100 p-2 rounded">
              {deployment.id}
            </p>
          </div>

          <div>
            <p className="text-sm text-gray-500 mb-1">Service Type</p>
            <p className="flex items-center gap-2">
              <Server size={18} />
              {deployment.serviceType}
            </p>
          </div>

          <div>
            <p className="text-sm text-gray-500 mb-1">Framework</p>
            <p className="flex items-center gap-2">
              <Code size={18} />
              {deployment.framework}
            </p>
          </div>

          <div>
            <p className="text-sm text-gray-500 mb-1">Branch</p>
            <p className="flex items-center gap-2">
              <GitBranch size={18} />
              {deployment.branch}
            </p>
          </div>

          <div>
            <p className="text-sm text-gray-500 mb-1">URL</p>
            <a
              href={deployment.url}
              target="_blank"
              rel="noopener noreferrer"
              className="text-blue-600 hover:text-blue-800 flex items-center gap-2"
            >
              {deployment.url}
              <ExternalLink size={16} />
            </a>
          </div>
        </div>

        <div className="space-y-4">
          <div>
            <p className="text-sm text-gray-500 mb-1">Platform Credential ID</p>
            <p>{deployment.platformCredentialID}</p>
          </div>

          <div>
            <p className="text-sm text-gray-500 mb-1">Last Deployed At</p>
            <p className="flex items-center gap-2">
              <Clock size={18} />
              {formatDate(deployment.lastDeployedAt)}
            </p>
          </div>

          <div>
            <p className="text-sm text-gray-500 mb-1">Last Updated At</p>
            <p className="flex items-center gap-2">
              <Clock size={18} />
              {formatDate(deployment.lastUpdatedAt)}
            </p>
          </div>

          <div>
            <p className="text-sm text-gray-500 mb-1">Metadata</p>
            <div className="bg-gray-100 p-2 rounded font-mono text-xs overflow-auto max-h-32">
              <pre>{JSON.stringify(deployment.metadata, null, 2)}</pre>
            </div>
          </div>
        </div>
      </div>

      <div className="mt-8 pt-4 border-t">
        <button
          onClick={onBack}
          className="bg-gray-100 hover:bg-gray-200 px-4 py-2 rounded-md text-gray-700 transition-colors"
        >
          Back to Card View
        </button>
      </div>
    </div>
  );
};

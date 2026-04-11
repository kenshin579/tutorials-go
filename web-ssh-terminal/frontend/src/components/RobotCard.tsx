import { Link } from 'react-router-dom';
import StatusBadge from './StatusBadge';

interface RobotCardProps {
  id: string;
  name: string;
  host: string;
  port: number;
  description?: string;
  isOnline: boolean;
}

export default function RobotCard({
  id,
  name,
  host,
  port,
  description,
  isOnline,
}: RobotCardProps) {
  return (
    <div className="border border-gray-200 rounded-xl p-6 bg-white shadow-sm hover:shadow-md transition-shadow">
      <div className="flex items-center justify-between mb-3">
        <h3 className="text-lg font-semibold text-gray-900">{name}</h3>
        <StatusBadge isOnline={isOnline} />
      </div>
      <p className="text-sm text-gray-500 font-mono">
        {host}:{port}
      </p>
      {description && (
        <p className="text-sm text-gray-600 mt-2">{description}</p>
      )}
      <div className="mt-4">
        {isOnline ? (
          <Link
            to={`/terminal/${id}`}
            className="inline-flex items-center px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-lg hover:bg-blue-700 transition-colors"
          >
            Connect
          </Link>
        ) : (
          <span className="inline-flex items-center px-4 py-2 bg-gray-200 text-gray-400 text-sm font-medium rounded-lg cursor-not-allowed">
            Offline
          </span>
        )}
      </div>
    </div>
  );
}

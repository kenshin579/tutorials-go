import { useEffect, useState } from 'react';
import { fetchRobots, type Robot } from '../api/robots';
import RobotList from '../components/RobotList';

export default function HomePage() {
  const [robots, setRobots] = useState<Robot[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchRobots()
      .then(setRobots)
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, []);

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white border-b border-gray-200 px-6 py-4">
        <h1 className="text-2xl font-bold text-gray-900">Robot Manager</h1>
        <p className="text-sm text-gray-500 mt-1">
          Manage and connect to your robots via SSH
        </p>
      </header>

      <main className="max-w-6xl mx-auto px-6 py-8">
        {loading && (
          <div className="text-center text-gray-500 py-12">Loading robots...</div>
        )}
        {error && (
          <div className="text-center text-red-500 py-12">Error: {error}</div>
        )}
        {!loading && !error && <RobotList robots={robots} />}
      </main>
    </div>
  );
}

import { useEffect, useState } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import { fetchRobots, type Robot } from '../api/robots';
import Terminal from '../components/Terminal';
import StatusBadge from '../components/StatusBadge';

export default function TerminalPage() {
  const { robotId } = useParams<{ robotId: string }>();
  const navigate = useNavigate();
  const [robot, setRobot] = useState<Robot | null>(null);
  const [connected, setConnected] = useState(true);
  const [elapsed, setElapsed] = useState(0);

  useEffect(() => {
    fetchRobots().then((robots) => {
      const found = robots.find((r) => r.id === robotId);
      if (found) setRobot(found);
    });
  }, [robotId]);

  useEffect(() => {
    if (!connected) return;
    const interval = setInterval(() => setElapsed((s) => s + 1), 1000);
    return () => clearInterval(interval);
  }, [connected]);

  const formatTime = (seconds: number) => {
    const h = Math.floor(seconds / 3600);
    const m = Math.floor((seconds % 3600) / 60);
    const s = seconds % 60;
    return [h, m, s].map((v) => String(v).padStart(2, '0')).join(':');
  };

  if (!robotId) return null;

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col">
      <header className="bg-white border-b border-gray-200 px-6 py-3 flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Link
            to="/"
            className="text-sm text-blue-600 hover:text-blue-800 flex items-center gap-1"
          >
            &larr; Back to List
          </Link>
          <h1 className="text-lg font-bold text-gray-900">Robot Manager</h1>
        </div>
        <button
          onClick={() => navigate('/')}
          className="px-3 py-1.5 bg-red-500 text-white text-sm rounded-lg hover:bg-red-600 transition-colors"
        >
          Disconnect
        </button>
      </header>

      <main className="flex-1 flex flex-col p-4 gap-3">
        {robot && (
          <div className="flex items-center gap-3 bg-white rounded-lg px-4 py-2 border border-gray-200">
            <StatusBadge isOnline={connected} />
            <span className="font-semibold text-gray-900">{robot.name}</span>
            <span className="text-sm text-gray-500 font-mono">
              ({robot.host}:{robot.port})
            </span>
          </div>
        )}

        <div className="flex-1 min-h-0">
          <Terminal
            robotId={robotId}
            onDisconnect={() => setConnected(false)}
          />
        </div>

        <div className="bg-gray-100 rounded-lg px-4 py-2 text-sm text-gray-600 flex items-center gap-4">
          <span>
            SSH: {robot?.username}@{robot?.host}:{robot?.port}
          </span>
          <span>Session: {formatTime(elapsed)}</span>
        </div>
      </main>
    </div>
  );
}

import type { Robot } from '../api/robots';
import RobotCard from './RobotCard';

interface RobotListProps {
  robots: Robot[];
}

export default function RobotList({ robots }: RobotListProps) {
  if (robots.length === 0) {
    return (
      <div className="text-center text-gray-500 py-12">
        No robots configured.
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {robots.map((robot) => (
        <RobotCard
          key={robot.id}
          id={robot.id}
          name={robot.name}
          host={robot.host}
          port={robot.port}
          description={robot.description}
          isOnline={robot.isOnline}
        />
      ))}
    </div>
  );
}

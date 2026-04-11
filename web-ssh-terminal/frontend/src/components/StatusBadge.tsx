interface StatusBadgeProps {
  isOnline: boolean;
}

export default function StatusBadge({ isOnline }: StatusBadgeProps) {
  return (
    <span
      className={`inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-medium ${
        isOnline ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
      }`}
    >
      <span
        className={`w-2 h-2 rounded-full ${isOnline ? 'bg-green-500' : 'bg-red-500'}`}
      />
      {isOnline ? 'Online' : 'Offline'}
    </span>
  );
}

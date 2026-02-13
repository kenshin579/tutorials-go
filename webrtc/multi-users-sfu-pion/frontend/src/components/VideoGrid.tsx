import { useEffect, useRef, type RefObject } from 'react';

interface VideoGridProps {
  localVideoRef: RefObject<HTMLVideoElement | null>;
  remoteStreams: Map<string, MediaStream>;
}

function RemoteVideo({ stream }: { stream: MediaStream }) {
  const ref = useRef<HTMLVideoElement>(null);

  useEffect(() => {
    if (ref.current) {
      ref.current.srcObject = stream;
    }
  }, [stream]);

  return (
    <div style={styles.videoBox}>
      <video ref={ref} autoPlay playsInline style={styles.video} />
      <span style={styles.label}>Remote</span>
    </div>
  );
}

export function VideoGrid({ localVideoRef, remoteStreams }: VideoGridProps) {
  const totalCount = 1 + remoteStreams.size;

  const columns = totalCount <= 1 ? 1 : totalCount <= 2 ? 2 : totalCount <= 4 ? 2 : 3;

  return (
    <div style={{ ...styles.grid, gridTemplateColumns: `repeat(${columns}, 1fr)` }}>
      <div style={styles.videoBox}>
        <video ref={localVideoRef} autoPlay muted playsInline style={styles.video} />
        <span style={styles.label}>나 (Local)</span>
      </div>
      {Array.from(remoteStreams.entries()).map(([id, stream]) => (
        <RemoteVideo key={id} stream={stream} />
      ))}
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  grid: {
    display: 'grid',
    gap: 12,
    flex: 1,
    alignContent: 'center',
  },
  videoBox: {
    position: 'relative',
    background: '#0f0f23',
    borderRadius: 12,
    overflow: 'hidden',
    border: '1px solid #2a2a4a',
    aspectRatio: '16 / 9',
  },
  video: {
    width: '100%',
    height: '100%',
    objectFit: 'cover',
  },
  label: {
    position: 'absolute',
    bottom: 12,
    left: 12,
    background: 'rgba(0,0,0,0.6)',
    padding: '4px 10px',
    borderRadius: 6,
    fontSize: 12,
    color: '#eee',
  },
};

import { useEffect, useRef, type RefObject } from 'react';
import type { RemoteTrack } from 'livekit-client';

interface TrackInfo {
  track: RemoteTrack;
  participantIdentity: string;
}

interface VideoGridProps {
  localVideoRef: RefObject<HTMLVideoElement | null>;
  remoteTracks: Map<string, TrackInfo>;
}

function RemoteVideo({ trackInfo }: { trackInfo: TrackInfo }) {
  const ref = useRef<HTMLVideoElement>(null);

  useEffect(() => {
    if (ref.current) {
      const el = trackInfo.track.attach();
      ref.current.srcObject = (el as HTMLVideoElement).srcObject;
    }
    return () => {
      trackInfo.track.detach();
    };
  }, [trackInfo.track]);

  return (
    <div style={styles.videoBox}>
      <video ref={ref} autoPlay playsInline style={styles.video} />
      <span style={styles.label}>{trackInfo.participantIdentity}</span>
    </div>
  );
}

export function VideoGrid({ localVideoRef, remoteTracks }: VideoGridProps) {
  const totalCount = 1 + remoteTracks.size;
  const columns = totalCount <= 1 ? 1 : totalCount <= 4 ? 2 : 3;

  return (
    <div style={{ ...styles.grid, gridTemplateColumns: `repeat(${columns}, 1fr)` }}>
      <div style={styles.videoBox}>
        <video ref={localVideoRef} autoPlay muted playsInline style={styles.video} />
        <span style={styles.label}>나 (Local)</span>
      </div>
      {Array.from(remoteTracks.entries()).map(([sid, trackInfo]) => (
        <RemoteVideo key={sid} trackInfo={trackInfo} />
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

import type { RefObject } from 'react';

interface VideoPanelProps {
  localVideoRef: RefObject<HTMLVideoElement | null>;
  remoteVideoRef: RefObject<HTMLVideoElement | null>;
}

export function VideoPanel({ localVideoRef, remoteVideoRef }: VideoPanelProps) {
  return (
    <div style={styles.container}>
      <div style={styles.videoBox}>
        <video ref={localVideoRef} autoPlay muted playsInline style={styles.video} />
        <span style={styles.label}>나 (Local)</span>
      </div>
      <div style={styles.videoBox}>
        <video ref={remoteVideoRef} autoPlay playsInline style={styles.video} />
        <span style={styles.label}>상대방 (Remote)</span>
      </div>
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  container: {
    display: 'flex',
    gap: 16,
    flex: 1,
  },
  videoBox: {
    flex: 1,
    position: 'relative',
    background: '#0f0f23',
    borderRadius: 12,
    overflow: 'hidden',
    border: '1px solid #2a2a4a',
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

import { useState } from 'react';
import { useSignaling } from './hooks/useSignaling';
import { useWebRTC } from './hooks/useWebRTC';
import { VideoGrid } from './components/VideoGrid';
import { ChatPanel } from './components/ChatPanel';

function App() {
  const [roomId, setRoomId] = useState<string | null>(null);
  const [userName, setUserName] = useState('');
  const [inputRoomId, setInputRoomId] = useState('room-123');
  const [inputName, setInputName] = useState('');

  const { connected, send, sendChat, chatMessages, setOnMessage } = useSignaling({
    roomId,
    userName,
  });
  const { localVideoRef, remoteStreams } = useWebRTC({ send, setOnMessage, connected });

  const handleJoin = () => {
    const rid = inputRoomId.trim();
    const name = inputName.trim() || `User-${Math.random().toString(36).slice(2, 6)}`;
    if (!rid) return;
    setUserName(name);
    setRoomId(rid);
  };

  if (!roomId) {
    return (
      <div style={styles.entry}>
        <h1 style={styles.title}>SFU 다자간 화상 통화</h1>
        <p style={styles.subtitle}>Room ID와 이름을 입력하고 입장하세요</p>
        <div style={styles.form}>
          <input
            style={styles.input}
            value={inputName}
            onChange={(e) => setInputName(e.target.value)}
            placeholder="이름 (선택)"
          />
          <input
            style={styles.input}
            value={inputRoomId}
            onChange={(e) => setInputRoomId(e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && handleJoin()}
            placeholder="Room ID (예: room-123)"
          />
          <button style={styles.button} onClick={handleJoin}>
            입장
          </button>
        </div>
      </div>
    );
  }

  return (
    <div style={styles.app}>
      <header style={styles.header}>
        <h1 style={styles.headerTitle}>SFU 다자간 화상 통화</h1>
        <span style={styles.roomBadge}>Room: {roomId}</span>
        <div style={styles.headerRight}>
          <span
            style={{ ...styles.statusDot, background: connected ? '#4ecca3' : '#e74c3c' }}
          />
          <span style={{ fontSize: 13, color: connected ? '#4ecca3' : '#e74c3c' }}>
            {connected ? 'Connected' : 'Connecting...'}
          </span>
          <span style={styles.peerCount}>
            참가자: {1 + remoteStreams.size}명
          </span>
        </div>
      </header>
      <div style={styles.main}>
        <div style={styles.videoSection}>
          <VideoGrid localVideoRef={localVideoRef} remoteStreams={remoteStreams} />
        </div>
        <ChatPanel messages={chatMessages} onSend={sendChat} disabled={!connected} />
      </div>
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  entry: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    minHeight: '100vh',
    gap: 24,
  },
  title: { fontSize: 28 },
  subtitle: { color: '#888', fontSize: 14 },
  form: { display: 'flex', gap: 8 },
  input: {
    padding: '10px 16px',
    borderRadius: 8,
    border: '1px solid #0f3460',
    background: '#16213e',
    color: '#eee',
    fontSize: 15,
    width: 200,
    outline: 'none',
  },
  button: {
    padding: '10px 20px',
    borderRadius: 8,
    border: 'none',
    background: '#4ecca3',
    color: '#1a1a2e',
    fontWeight: 600,
    fontSize: 15,
    cursor: 'pointer',
  },
  app: { display: 'flex', flexDirection: 'column', flex: 1 },
  header: {
    background: '#16213e',
    padding: '12px 24px',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
    borderBottom: '1px solid #0f3460',
  },
  headerTitle: { fontSize: 18, fontWeight: 600 },
  roomBadge: {
    background: '#0f3460',
    padding: '4px 12px',
    borderRadius: 12,
    fontSize: 13,
    color: '#7ec8e3',
  },
  headerRight: { display: 'flex', alignItems: 'center', gap: 8 },
  statusDot: { width: 8, height: 8, borderRadius: '50%', display: 'inline-block' },
  peerCount: {
    marginLeft: 8,
    fontSize: 13,
    color: '#7ec8e3',
  },
  main: { display: 'flex', flex: 1, height: 'calc(100vh - 49px)' },
  videoSection: { flex: 1, display: 'flex', flexDirection: 'column', padding: 16 },
};

export default App;

import { useState } from 'react';
import { useSignaling } from './hooks/useSignaling';
import { useWebRTC } from './hooks/useWebRTC';
import { VideoPanel } from './components/VideoPanel';
import { ChatPanel } from './components/ChatPanel';

function App() {
  const [roomId, setRoomId] = useState<string | null>(null);
  const [inputValue, setInputValue] = useState('room-123');

  const { connected, send, setOnMessage } = useSignaling(roomId);
  const { localVideoRef, remoteVideoRef, chatMessages, dcOpen, sendChat, createOffer } =
    useWebRTC({ send, setOnMessage, connected });

  if (!roomId) {
    return (
      <div style={styles.entry}>
        <h1 style={styles.title}>WebRTC 1:1 화상 통화</h1>
        <p style={styles.subtitle}>Room ID를 입력하고 입장하세요</p>
        <div style={styles.form}>
          <input
            style={styles.input}
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && inputValue.trim() && setRoomId(inputValue.trim())}
            placeholder="Room ID 입력 (예: room-123)"
          />
          <button
            style={styles.button}
            onClick={() => inputValue.trim() && setRoomId(inputValue.trim())}
          >
            입장
          </button>
        </div>
      </div>
    );
  }

  return (
    <div style={styles.app}>
      <header style={styles.header}>
        <h1 style={styles.headerTitle}>WebRTC 1:1 화상 통화</h1>
        <span style={styles.roomBadge}>Room: {roomId}</span>
        <div style={styles.headerRight}>
          <span style={{ ...styles.statusDot, background: connected ? '#4ecca3' : '#e74c3c' }} />
          <span style={{ fontSize: 13, color: connected ? '#4ecca3' : '#e74c3c' }}>
            {connected ? 'Connected' : 'Connecting...'}
          </span>
          {connected && (
            <button style={styles.callButton} onClick={createOffer}>
              통화 시작
            </button>
          )}
        </div>
      </header>
      <div style={styles.main}>
        <div style={styles.videoSection}>
          <VideoPanel localVideoRef={localVideoRef} remoteVideoRef={remoteVideoRef} />
        </div>
        <ChatPanel messages={chatMessages} onSend={sendChat} disabled={!dcOpen} />
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
    width: 240,
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
  callButton: {
    marginLeft: 8,
    padding: '6px 14px',
    borderRadius: 8,
    border: 'none',
    background: '#4ecca3',
    color: '#1a1a2e',
    fontWeight: 600,
    fontSize: 13,
    cursor: 'pointer',
  },
  main: { display: 'flex', flex: 1, height: 'calc(100vh - 49px)' },
  videoSection: { flex: 1, display: 'flex', flexDirection: 'column', padding: 16 },
};

export default App;

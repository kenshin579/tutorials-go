import { useEffect, useRef, useState } from 'react';

interface ChatMessage {
  text: string;
  mine: boolean;
}

interface ChatPanelProps {
  messages: ChatMessage[];
  onSend: (text: string) => void;
  disabled: boolean;
}

export function ChatPanel({ messages, onSend, disabled }: ChatPanelProps) {
  const [input, setInput] = useState('');
  const listRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (listRef.current) {
      listRef.current.scrollTop = listRef.current.scrollHeight;
    }
  }, [messages]);

  const handleSend = () => {
    const text = input.trim();
    if (!text) return;
    onSend(text);
    setInput('');
  };

  return (
    <div style={styles.container}>
      <div style={styles.header}>DataChannel Chat</div>
      <div ref={listRef} style={styles.messages}>
        {messages.map((msg, i) => (
          <div key={i} style={{ ...styles.msg, ...(msg.mine ? styles.mine : styles.peer) }}>
            {msg.text}
          </div>
        ))}
      </div>
      <div style={styles.inputRow}>
        <input
          style={styles.input}
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={(e) => e.key === 'Enter' && handleSend()}
          placeholder={disabled ? '연결 대기 중...' : '메시지 입력...'}
          disabled={disabled}
        />
        <button style={styles.button} onClick={handleSend} disabled={disabled}>
          전송
        </button>
      </div>
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  container: {
    width: 320,
    background: '#16213e',
    borderLeft: '1px solid #0f3460',
    display: 'flex',
    flexDirection: 'column',
  },
  header: {
    padding: '12px 16px',
    borderBottom: '1px solid #0f3460',
    fontSize: 14,
    fontWeight: 600,
    color: '#7ec8e3',
  },
  messages: {
    flex: 1,
    padding: 12,
    overflowY: 'auto',
    display: 'flex',
    flexDirection: 'column',
    gap: 8,
  },
  msg: {
    maxWidth: '85%',
    padding: '8px 12px',
    borderRadius: 12,
    fontSize: 13,
    lineHeight: 1.4,
    color: '#e0e0e0',
  },
  mine: {
    alignSelf: 'flex-end',
    background: '#0f3460',
    borderBottomRightRadius: 4,
  },
  peer: {
    alignSelf: 'flex-start',
    background: '#2a2a4a',
    borderBottomLeftRadius: 4,
  },
  inputRow: {
    display: 'flex',
    gap: 8,
    padding: 12,
    borderTop: '1px solid #0f3460',
  },
  input: {
    flex: 1,
    padding: '8px 12px',
    borderRadius: 8,
    border: '1px solid #0f3460',
    background: '#0f0f23',
    color: '#eee',
    fontSize: 13,
    outline: 'none',
  },
  button: {
    padding: '8px 16px',
    borderRadius: 8,
    border: 'none',
    background: '#4ecca3',
    color: '#1a1a2e',
    fontWeight: 600,
    fontSize: 13,
    cursor: 'pointer',
  },
};

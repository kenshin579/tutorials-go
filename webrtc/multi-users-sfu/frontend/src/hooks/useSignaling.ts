import { useCallback, useEffect, useRef, useState } from 'react';

export interface SignalingMessage {
  type: 'offer' | 'answer' | 'ice' | 'join' | 'leave' | 'chat';
  senderId?: string;
  senderName?: string;
  message?: string;
  payload?: unknown;
}

export interface ChatMessage {
  text: string;
  mine: boolean;
  senderName?: string;
}

interface UseSignalingParams {
  roomId: string | null;
  userName: string;
}

export function useSignaling({ roomId, userName }: UseSignalingParams) {
  const wsRef = useRef<WebSocket | null>(null);
  const [connected, setConnected] = useState(false);
  const [chatMessages, setChatMessages] = useState<ChatMessage[]>([]);
  const onMessageRef = useRef<((msg: SignalingMessage) => void) | null>(null);

  useEffect(() => {
    if (!roomId) return;

    const ws = new WebSocket(`ws://localhost:8080/ws?roomId=${roomId}`);
    wsRef.current = ws;

    ws.onopen = () => {
      console.log('[Signaling] connected');
      setConnected(true);
    };

    ws.onmessage = (event) => {
      const msg: SignalingMessage = JSON.parse(event.data);
      console.log('[Signaling] received:', msg.type);

      if (msg.type === 'chat') {
        setChatMessages((prev) => [
          ...prev,
          { text: msg.message ?? '', mine: false, senderName: msg.senderName },
        ]);
      } else {
        onMessageRef.current?.(msg);
      }
    };

    ws.onclose = () => {
      console.log('[Signaling] disconnected');
      setConnected(false);
    };

    ws.onerror = (err) => {
      console.error('[Signaling] error:', err);
    };

    return () => {
      ws.close();
      wsRef.current = null;
    };
  }, [roomId]);

  const send = useCallback((msg: SignalingMessage) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      console.log('[Signaling] send:', msg.type);
      wsRef.current.send(JSON.stringify(msg));
    }
  }, []);

  const sendChat = useCallback(
    (text: string) => {
      send({ type: 'chat', senderName: userName, message: text });
      setChatMessages((prev) => [...prev, { text, mine: true }]);
    },
    [send, userName],
  );

  const setOnMessage = useCallback((handler: (msg: SignalingMessage) => void) => {
    onMessageRef.current = handler;
  }, []);

  return { connected, send, sendChat, chatMessages, setOnMessage };
}

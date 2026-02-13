import { useCallback, useEffect, useRef, useState } from 'react';

export interface SignalingMessage {
  type: 'offer' | 'answer' | 'ice';
  payload: unknown;
}

export function useSignaling(roomId: string | null) {
  const wsRef = useRef<WebSocket | null>(null);
  const [connected, setConnected] = useState(false);
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
      onMessageRef.current?.(msg);
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

  const setOnMessage = useCallback((handler: (msg: SignalingMessage) => void) => {
    onMessageRef.current = handler;
  }, []);

  return { connected, send, setOnMessage };
}

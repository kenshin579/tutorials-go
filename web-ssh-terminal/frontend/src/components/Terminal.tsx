import { useEffect, useRef, useCallback } from 'react';
import { Terminal as XTerm } from '@xterm/xterm';
import { FitAddon } from '@xterm/addon-fit';
import '@xterm/xterm/css/xterm.css';

interface TerminalProps {
  robotId: string;
  onDisconnect?: () => void;
}

export default function Terminal({ robotId, onDisconnect }: TerminalProps) {
  const terminalRef = useRef<HTMLDivElement>(null);
  const cleanupRef = useRef<(() => void) | null>(null);

  const connect = useCallback(() => {
    if (!terminalRef.current) return;

    const term = new XTerm({
      cursorBlink: true,
      fontSize: 14,
      fontFamily: 'Menlo, Monaco, "Courier New", monospace',
      theme: {
        background: '#1e1e2e',
        foreground: '#cdd6f4',
        cursor: '#f5e0dc',
      },
    });

    const fitAddon = new FitAddon();
    term.loadAddon(fitAddon);
    term.open(terminalRef.current);
    fitAddon.fit();

    term.writeln('Connecting to robot...');

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const ws = new WebSocket(
      `${protocol}//${window.location.host}/ws/terminal?robotId=${robotId}`
    );

    ws.onopen = () => {
      term.writeln('WebSocket connected. Waiting for SSH...');
    };

    ws.onmessage = (event) => {
      const data = event.data;
      try {
        const parsed = JSON.parse(data);
        if (parsed.type === 'status') {
          term.writeln(`SSH ${parsed.message}\r\n`);
          return;
        }
        if (parsed.type === 'error') {
          term.writeln(`\r\nError: ${parsed.message}`);
          return;
        }
      } catch {
        // Not JSON — raw terminal output
      }
      term.write(data);
    };

    ws.onclose = () => {
      term.writeln('\r\n\r\nConnection closed.');
      onDisconnect?.();
    };

    ws.onerror = () => {
      term.writeln('\r\nWebSocket error.');
    };

    term.onData((data) => {
      if (ws.readyState === WebSocket.OPEN) {
        ws.send(data);
      }
    });

    const handleResize = () => {
      fitAddon.fit();
      if (ws.readyState === WebSocket.OPEN) {
        ws.send(
          JSON.stringify({ type: 'resize', cols: term.cols, rows: term.rows })
        );
      }
    };
    window.addEventListener('resize', handleResize);

    cleanupRef.current = () => {
      window.removeEventListener('resize', handleResize);
      ws.close();
      term.dispose();
    };
  }, [robotId, onDisconnect]);

  useEffect(() => {
    connect();
    return () => {
      cleanupRef.current?.();
    };
  }, [connect]);

  return (
    <div
      ref={terminalRef}
      className="w-full h-full min-h-[400px] bg-[#1e1e2e] rounded-lg p-1"
    />
  );
}

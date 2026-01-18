import { useState, useEffect, useCallback } from 'react';
import mqtt from 'mqtt';
import type { MqttClient } from 'mqtt';

interface DeviceState {
  deviceId: string;
  status: 'idle' | 'running';
  temperature: number;
  timestamp: number;
}

// 로그 엔트리 타입
interface LogEntry {
  id: number;
  time: string;
  direction: 'in' | 'out';  // in: 수신, out: 송신
  type: 'STATE' | 'CMD';
  message: string;
}

const MAX_LOG_SIZE = 50;  // 최대 로그 개수

export function useMqtt(brokerUrl: string) {
  const [client, setClient] = useState<MqttClient | null>(null);
  const [connected, setConnected] = useState(false);
  const [deviceState, setDeviceState] = useState<DeviceState | null>(null);
  const [logs, setLogs] = useState<LogEntry[]>([]);

  // 로그 추가 함수
  const addLog = useCallback((direction: 'in' | 'out', type: 'STATE' | 'CMD', message: string) => {
    const now = new Date();
    const time = now.toLocaleTimeString('ko-KR', { hour12: false });

    setLogs(prev => {
      const newLog: LogEntry = {
        id: Date.now(),
        time,
        direction,
        type,
        message,
      };
      const updated = [newLog, ...prev];
      return updated.slice(0, MAX_LOG_SIZE);  // 최대 개수 제한
    });
  }, []);

  // 로그 초기화 함수
  const clearLogs = useCallback(() => {
    setLogs([]);
  }, []);

  useEffect(() => {
    const mqttClient = mqtt.connect(brokerUrl, {
      protocolVersion: 5,
      reconnectPeriod: 1000,
    });

    mqttClient.on('connect', () => {
      console.log('[MQTT] Connected');
      setConnected(true);
      mqttClient.subscribe('device/1/state', (err) => {
        if (err) {
          console.error('[MQTT] Subscribe failed:', err);
        } else {
          console.log('[MQTT] Subscribed to device/1/state');
        }
      });
    });

    mqttClient.on('close', () => {
      console.log('[MQTT] Disconnected');
      setConnected(false);
    });

    mqttClient.on('error', (err) => {
      console.error('[MQTT] Error:', err);
    });

    mqttClient.on('message', (topic, payload) => {
      if (topic === 'device/1/state') {
        try {
          const state: DeviceState = JSON.parse(payload.toString());
          setDeviceState(state);
          // 수신 로그 추가
          addLog('in', 'STATE', `${state.status} ${state.temperature.toFixed(1)}°C`);
        } catch (err) {
          console.error('[MQTT] Failed to parse message:', err);
        }
      }
    });

    setClient(mqttClient);

    return () => {
      mqttClient.end();
    };
  }, [brokerUrl, addLog]);

  const sendCommand = useCallback((action: 'start' | 'stop') => {
    if (client && connected) {
      client.publish('device/1/command', JSON.stringify({ action }), { qos: 1 });
      // 송신 로그 추가
      addLog('out', 'CMD', action);
      console.log('[MQTT] Command sent:', action);
    }
  }, [client, connected, addLog]);

  return { connected, deviceState, logs, sendCommand, clearLogs };
}

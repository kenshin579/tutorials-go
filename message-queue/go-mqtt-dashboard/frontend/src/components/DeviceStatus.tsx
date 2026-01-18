import { useMqtt } from '../hooks/useMqtt';
import styles from './DeviceStatus.module.css';

export function DeviceStatus() {
  const { connected, deviceState, logs, sendCommand, clearLogs } = useMqtt('ws://localhost:9001');

  return (
    <div className={styles.container}>
      <h1 className={styles.title}>🖥️ Device Dashboard</h1>

      <div className={styles.connectionStatus}>
        <span>Connection: </span>
        <span className={connected ? styles.connected : styles.disconnected}>
          {connected ? '🟢 Connected' : '🔴 Disconnected'}
        </span>
      </div>

      {deviceState && (
        <table className={styles.stateTable}>
          <tbody>
            <tr>
              <td>Status</td>
              <td className={deviceState.status === 'running' ? styles.running : styles.idle}>
                {deviceState.status === 'running' ? '🔵' : '⚪'} {deviceState.status}
              </td>
            </tr>
            <tr>
              <td>Temperature</td>
              <td>{deviceState.temperature.toFixed(1)}°C</td>
            </tr>
          </tbody>
        </table>
      )}

      <div className={styles.buttonGroup}>
        <button
          className={styles.startButton}
          onClick={() => sendCommand('start')}
          disabled={!connected}
        >
          ▶ Start
        </button>
        <button
          className={styles.stopButton}
          onClick={() => sendCommand('stop')}
          disabled={!connected}
        >
          ⏹ Stop
        </button>
      </div>

      {/* 메시지 로그 영역 */}
      <div className={styles.logSection}>
        <div className={styles.logHeader}>
          <span>📋 Message Log</span>
          <button className={styles.clearButton} onClick={clearLogs}>
            Clear
          </button>
        </div>
        <div className={styles.logList}>
          {logs.length === 0 ? (
            <div className={styles.logEmpty}>No messages yet</div>
          ) : (
            logs.map(log => (
              <div
                key={log.id}
                className={`${styles.logEntry} ${log.direction === 'in' ? styles.logIn : styles.logOut}`}
              >
                <span className={styles.logTime}>{log.time}</span>
                <span className={styles.logDirection}>
                  {log.direction === 'in' ? '←' : '→'}
                </span>
                <span className={styles.logType}>{log.type}</span>
                <span className={styles.logMessage}>{log.message}</span>
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  );
}

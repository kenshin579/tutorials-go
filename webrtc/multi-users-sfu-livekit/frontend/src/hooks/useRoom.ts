import { useCallback, useEffect, useRef, useState } from 'react';
import {
  Room,
  RoomEvent,
  type RemoteTrack,
  type RemoteTrackPublication,
  type RemoteParticipant,
  Track,
} from 'livekit-client';

const LIVEKIT_URL = 'ws://localhost:7880';

export interface ChatMessage {
  text: string;
  mine: boolean;
  senderName?: string;
}

interface TrackInfo {
  track: RemoteTrack;
  participantIdentity: string;
}

interface UseRoomParams {
  roomId: string | null;
  userName: string;
}

export function useRoom({ roomId, userName }: UseRoomParams) {
  const roomRef = useRef<Room | null>(null);
  const [connected, setConnected] = useState(false);
  const [chatMessages, setChatMessages] = useState<ChatMessage[]>([]);
  const [remoteTracks, setRemoteTracks] = useState<Map<string, TrackInfo>>(new Map());
  const [participantCount, setParticipantCount] = useState(1);
  const localVideoRef = useRef<HTMLVideoElement>(null);

  const updateParticipantCount = useCallback(() => {
    const room = roomRef.current;
    if (!room) return;
    setParticipantCount(1 + room.remoteParticipants.size);
  }, []);

  useEffect(() => {
    if (!roomId || !userName) return;

    const room = new Room({
      adaptiveStream: true,
      dynacast: true,
    });
    roomRef.current = room;

    const handleTrackSubscribed = (
      track: RemoteTrack,
      _publication: RemoteTrackPublication,
      participant: RemoteParticipant,
    ) => {
      if (track.kind === Track.Kind.Video) {
        setRemoteTracks((prev) => {
          const next = new Map(prev);
          next.set(track.sid!, { track, participantIdentity: participant.identity });
          return next;
        });
      }
      if (track.kind === Track.Kind.Audio) {
        const el = track.attach();
        document.body.appendChild(el);
      }
      updateParticipantCount();
    };

    const handleTrackUnsubscribed = (
      track: RemoteTrack,
      _publication: RemoteTrackPublication,
      _participant: RemoteParticipant,
    ) => {
      track.detach();
      if (track.kind === Track.Kind.Video) {
        setRemoteTracks((prev) => {
          const next = new Map(prev);
          next.delete(track.sid!);
          return next;
        });
      }
    };

    const handleParticipantDisconnected = (_participant: RemoteParticipant) => {
      updateParticipantCount();
    };

    const handleParticipantConnected = (_participant: RemoteParticipant) => {
      updateParticipantCount();
    };

    const handleDataReceived = (
      payload: Uint8Array,
      participant?: RemoteParticipant,
    ) => {
      const decoder = new TextDecoder();
      const message = JSON.parse(decoder.decode(payload));
      if (message.type === 'chat') {
        setChatMessages((prev) => [
          ...prev,
          { text: message.message, mine: false, senderName: participant?.identity },
        ]);
      }
    };

    const handleDisconnect = () => {
      console.log('[LiveKit] disconnected');
      setConnected(false);
    };

    room
      .on(RoomEvent.TrackSubscribed, handleTrackSubscribed)
      .on(RoomEvent.TrackUnsubscribed, handleTrackUnsubscribed)
      .on(RoomEvent.ParticipantDisconnected, handleParticipantDisconnected)
      .on(RoomEvent.ParticipantConnected, handleParticipantConnected)
      .on(RoomEvent.DataReceived, handleDataReceived as (...args: unknown[]) => void)
      .on(RoomEvent.Disconnected, handleDisconnect);

    const connect = async () => {
      try {
        const res = await fetch(
          `http://localhost:8080/token?roomId=${encodeURIComponent(roomId)}&userName=${encodeURIComponent(userName)}`,
        );
        const { token } = await res.json();

        await room.connect(LIVEKIT_URL, token);
        console.log('[LiveKit] connected to room:', room.name);
        setConnected(true);

        await room.localParticipant.enableCameraAndMicrophone();

        const camPub = room.localParticipant.getTrackPublication(Track.Source.Camera);
        if (camPub?.track && localVideoRef.current) {
          const el = camPub.track.attach();
          localVideoRef.current.srcObject = (el as HTMLVideoElement).srcObject;
        }

        updateParticipantCount();
      } catch (err) {
        console.error('[LiveKit] connection failed:', err);
      }
    };

    connect();

    return () => {
      room.disconnect();
      roomRef.current = null;
      setRemoteTracks(new Map());
      setConnected(false);
      setParticipantCount(1);
    };
  }, [roomId, userName, updateParticipantCount]);

  const sendChat = useCallback(
    (text: string) => {
      const room = roomRef.current;
      if (!room) return;

      const encoder = new TextEncoder();
      const data = encoder.encode(JSON.stringify({ type: 'chat', message: text }));
      room.localParticipant.publishData(data, { reliable: true });
      setChatMessages((prev) => [...prev, { text, mine: true }]);
    },
    [],
  );

  return { connected, localVideoRef, remoteTracks, chatMessages, sendChat, participantCount };
}

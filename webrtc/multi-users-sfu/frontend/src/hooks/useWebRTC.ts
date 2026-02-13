import { useCallback, useEffect, useRef, useState } from 'react';
import type { SignalingMessage } from './useSignaling';

const ICE_SERVERS = [{ urls: 'stun:stun.l.google.com:19302' }];

interface UseWebRTCParams {
  send: (msg: SignalingMessage) => void;
  setOnMessage: (handler: (msg: SignalingMessage) => void) => void;
  connected: boolean;
}

export function useWebRTC({ send, setOnMessage, connected }: UseWebRTCParams) {
  const pcRef = useRef<RTCPeerConnection | null>(null);
  const localVideoRef = useRef<HTMLVideoElement>(null);
  const localStreamRef = useRef<MediaStream | null>(null);
  const [remoteStreams, setRemoteStreams] = useState<Map<string, MediaStream>>(new Map());

  const startCall = useCallback(async () => {
    const pc = new RTCPeerConnection({ iceServers: ICE_SERVERS });
    pcRef.current = pc;

    pc.onicecandidate = (event) => {
      if (event.candidate) {
        console.log('[ICE] sending candidate');
        send({ type: 'ice', payload: event.candidate.toJSON() });
      }
    };

    pc.oniceconnectionstatechange = () => {
      console.log('[ICE] state:', pc.iceConnectionState);
    };

    pc.ontrack = (event) => {
      const stream = event.streams[0];
      if (!stream) return;
      console.log('[WebRTC] remote track received, stream:', stream.id);
      setRemoteStreams((prev) => new Map(prev).set(stream.id, stream));

      stream.onremovetrack = () => {
        if (stream.getTracks().length === 0) {
          console.log('[WebRTC] stream removed:', stream.id);
          setRemoteStreams((prev) => {
            const next = new Map(prev);
            next.delete(stream.id);
            return next;
          });
        }
      };
    };

    // Get local media
    const stream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true });
    localStreamRef.current = stream;
    if (localVideoRef.current) {
      localVideoRef.current.srcObject = stream;
    }
    stream.getTracks().forEach((track) => pc.addTrack(track, stream));

    // Set up signaling message handler
    setOnMessage(async (msg: SignalingMessage) => {
      if (msg.type === 'offer') {
        // Renegotiation offer from SFU (new tracks added)
        console.log('[WebRTC] received offer (renegotiation)');
        await pc.setRemoteDescription(
          new RTCSessionDescription(msg.payload as RTCSessionDescriptionInit),
        );
        const answer = await pc.createAnswer();
        await pc.setLocalDescription(answer);
        send({ type: 'answer', payload: answer });
      } else if (msg.type === 'answer') {
        console.log('[WebRTC] received answer');
        await pc.setRemoteDescription(
          new RTCSessionDescription(msg.payload as RTCSessionDescriptionInit),
        );
      } else if (msg.type === 'ice') {
        console.log('[ICE] received candidate');
        await pc.addIceCandidate(new RTCIceCandidate(msg.payload as RTCIceCandidateInit));
      } else if (msg.type === 'leave') {
        console.log('[WebRTC] peer left:', msg.senderId);
      }
    });

    // Create initial offer to SFU
    const offer = await pc.createOffer();
    await pc.setLocalDescription(offer);
    send({ type: 'offer', payload: offer });
    console.log('[WebRTC] initial offer sent to SFU');
  }, [send, setOnMessage]);

  useEffect(() => {
    if (connected) {
      startCall();
    }
    return () => {
      localStreamRef.current?.getTracks().forEach((t) => t.stop());
      pcRef.current?.close();
      pcRef.current = null;
      setRemoteStreams(new Map());
    };
  }, [connected, startCall]);

  return { localVideoRef, remoteStreams };
}

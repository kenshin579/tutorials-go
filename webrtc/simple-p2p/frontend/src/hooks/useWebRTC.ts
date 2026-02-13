import { useCallback, useEffect, useRef, useState } from 'react';
import type { SignalingMessage } from './useSignaling';

const ICE_SERVERS = [{ urls: 'stun:stun.l.google.com:19302' }];

interface ChatMessage {
  text: string;
  mine: boolean;
}

interface UseWebRTCParams {
  send: (msg: SignalingMessage) => void;
  setOnMessage: (handler: (msg: SignalingMessage) => void) => void;
  connected: boolean;
}

export function useWebRTC({ send, setOnMessage, connected }: UseWebRTCParams) {
  const pcRef = useRef<RTCPeerConnection | null>(null);
  const dcRef = useRef<RTCDataChannel | null>(null);
  const localVideoRef = useRef<HTMLVideoElement>(null);
  const remoteVideoRef = useRef<HTMLVideoElement>(null);
  const [chatMessages, setChatMessages] = useState<ChatMessage[]>([]);
  const [dcOpen, setDcOpen] = useState(false);
  const isCallerRef = useRef(false);
  const localStreamRef = useRef<MediaStream | null>(null);

  const setupDataChannel = useCallback((dc: RTCDataChannel) => {
    dcRef.current = dc;
    dc.onopen = () => {
      console.log('[DataChannel] open');
      setDcOpen(true);
    };
    dc.onclose = () => {
      console.log('[DataChannel] close');
      setDcOpen(false);
    };
    dc.onmessage = (event) => {
      setChatMessages((prev) => [...prev, { text: event.data, mine: false }]);
    };
  }, []);

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
      console.log('[WebRTC] remote track received');
      if (remoteVideoRef.current) {
        remoteVideoRef.current.srcObject = event.streams[0];
      }
    };

    pc.ondatachannel = (event) => {
      console.log('[WebRTC] remote data channel received');
      setupDataChannel(event.channel);
    };

    const stream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true });
    localStreamRef.current = stream;
    if (localVideoRef.current) {
      localVideoRef.current.srcObject = stream;
    }
    stream.getTracks().forEach((track) => pc.addTrack(track, stream));

    setOnMessage(async (msg: SignalingMessage) => {
      if (msg.type === 'offer') {
        console.log('[WebRTC] received offer, creating answer');
        await pc.setRemoteDescription(new RTCSessionDescription(msg.payload as RTCSessionDescriptionInit));
        const answer = await pc.createAnswer();
        await pc.setLocalDescription(answer);
        send({ type: 'answer', payload: answer });
      } else if (msg.type === 'answer') {
        console.log('[WebRTC] received answer');
        await pc.setRemoteDescription(new RTCSessionDescription(msg.payload as RTCSessionDescriptionInit));
      } else if (msg.type === 'ice') {
        console.log('[ICE] received candidate');
        await pc.addIceCandidate(new RTCIceCandidate(msg.payload as RTCIceCandidateInit));
      }
    });
  }, [send, setOnMessage, setupDataChannel]);

  const createOffer = useCallback(async () => {
    const pc = pcRef.current;
    if (!pc) return;

    isCallerRef.current = true;
    const dc = pc.createDataChannel('chat');
    setupDataChannel(dc);

    const offer = await pc.createOffer();
    await pc.setLocalDescription(offer);
    send({ type: 'offer', payload: offer });
    console.log('[WebRTC] offer sent');
  }, [send, setupDataChannel]);

  const sendChat = useCallback((text: string) => {
    if (dcRef.current?.readyState === 'open') {
      dcRef.current.send(text);
      setChatMessages((prev) => [...prev, { text, mine: true }]);
    }
  }, []);

  useEffect(() => {
    if (connected) {
      startCall();
    }
    return () => {
      localStreamRef.current?.getTracks().forEach((t) => t.stop());
      pcRef.current?.close();
      pcRef.current = null;
    };
  }, [connected, startCall]);

  return {
    localVideoRef,
    remoteVideoRef,
    chatMessages,
    dcOpen,
    sendChat,
    createOffer,
  };
}

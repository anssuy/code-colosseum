"use client";

import {
  createContext,
  ReactNode,
  useCallback,
  useContext,
  useRef,
  useState,
} from "react";

type SubmissionResult = {
  submissionId: string;
  status: string;
  passedTests: number;
  totalTests: number;
};

type MatchContextValue = {
  connected: boolean;
  isQueued: boolean;
  currentMatchId: string | null;
  matchStarted: boolean;
  lastResult: SubmissionResult | null;
  winnerId: string | null;
  abandoned: boolean;
  error: string | null;

  joinQueue: () => void;
  leaveQueue: () => void;
  sendReady: () => void;
  submit: (language: string, sourceCode: string) => void;
};

const MatchContext = createContext<MatchContextValue | null>(null);

export function MatchProvider({ children }: { children: ReactNode }) {
  const wsRef = useRef<WebSocket | null>(null);

  const [connected, setConnected] = useState(false);
  const [isQueued, setIsQueued] = useState(false);
  const [currentMatchId, setCurrentMatchId] = useState<string | null>(null);
  const [matchStarted, setMatchStarted] = useState(false);
  const [lastResult, setLastResult] = useState<SubmissionResult | null>(null);
  const [winnerId, setWinnerId] = useState<string | null>(null);
  const [abandoned, setAbandoned] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const ensureConnected = useCallback(() => {
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      return wsRef.current;
    }

    const ws = new WebSocket(
      `${process.env.NEXT_PUBLIC_WS_URL ?? "ws://localhost:8080"}/api/ws`,
    );

    ws.onopen = () => setConnected(true);
    ws.onclose = () => setConnected(false);

    ws.onmessage = (event) => {
      const msg = JSON.parse(event.data);

      console.log(msg);

      switch (msg.type) {
        case "queue_joined":
          setIsQueued(true);
          break;
        case "queue_left":
          setIsQueued(false);
          break;
        case "match_found":
          setIsQueued(false);
          setCurrentMatchId(msg.payload.matchId);
          setMatchStarted(false);
          setLastResult(null);
          setWinnerId(null);
          setAbandoned(false);
          break;
        case "match_started":
          setMatchStarted(true);
          break;
        case "submission_result":
          setLastResult(msg.payload);
          break;
        case "match_finished":
          setWinnerId(msg.payload?.winnerId ?? null);
          break;
        case "match_abandoned":
          setAbandoned(true);
          break;
        case "error":
          setError(msg.payload.message);
          break;
      }
    };

    wsRef.current = ws;
    return ws;
  }, []);

  const send = useCallback(
    (type: string, payload: unknown = {}) => {
      const ws = ensureConnected();

      const doSend = () => ws.send(JSON.stringify({ type, payload }));

      if (ws.readyState === WebSocket.OPEN) {
        doSend();
      } else {
        ws.addEventListener("open", doSend, { once: true });
      }
    },
    [ensureConnected],
  );

  const joinQueue = useCallback(() => send("queue_join"), [send]);
  const leaveQueue = useCallback(() => send("queue_left"), [send]);
  const sendReady = useCallback(() => send("ready"), [send]);
  const submit = useCallback(
    (language: string, sourceCode: string) =>
      send("submit", { language, sourceCode }),
    [send],
  );

  return (
    <MatchContext.Provider
      value={{
        connected,
        isQueued,
        currentMatchId,
        matchStarted,
        lastResult,
        winnerId,
        abandoned,
        error,
        joinQueue,
        leaveQueue,
        sendReady,
        submit,
      }}
    >
      {children}
    </MatchContext.Provider>
  );
}

export function useMatch() {
  const ctx = useContext(MatchContext);
  if (!ctx) {
    throw new Error("useMatch must be used within a MatchProvider");
  }
  return ctx;
}

"use client";

import { CircleUserRound, Search } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useRef } from "react";

import { SidebarTrigger } from "@/components/ui/sidebar";
import { toast } from "@/components/ui/toast";
import { getActiveMatch } from "@/lib/api/matches";
import { useAuth } from "@/providers/AuthProvider";
import { useMatch } from "@/providers/MatchProvider";

export default function Header() {
  const { user } = useAuth();
  const { isQueued, currentMatchId, joinQueue, leaveQueue, sendReady } =
    useMatch();

  const router = useRouter();
  const toastIdRef = useRef<string | number | null>(null);
  const activeMatchToastIdRef = useRef<string | number | null>(null);
  const readySentRef = useRef(false);

  useEffect(() => {
    getActiveMatch().then((matchData) => {
      if (!matchData?.match) return;

      toast.close();

      activeMatchToastIdRef.current = toast.add({
        title: "You have an active match",
        description: "You already have a match in progress.",
        type: "loading",
        actionProps: {
          children: "Go to match",
          onClick: () => {
            toast.close(activeMatchToastIdRef.current?.toString());
            activeMatchToastIdRef.current = null;
            router.push(`/matches/${matchData.match.id}`);
          },
        },
      });
    });
  }, [router]);

  useEffect(() => {
    if (!isQueued) return;

    toastIdRef.current = toast.add({
      title: "Searching for a match...",
      description: "This may take a moment",
      type: "loading",
      actionProps: { children: "Cancel", onClick: () => leaveQueue() },
    });

    return () => {
      if (toastIdRef.current) {
        toast.close(toastIdRef.current.toString());
        toastIdRef.current = null;
      }
    };
  }, [isQueued, leaveQueue]);

  useEffect(() => {
    if (currentMatchId && toastIdRef.current) {
      toast.update(toastIdRef.current.toString(), {
        type: "success",
        title: "Match found!",
      });
      toastIdRef.current = null;
    }

    if (currentMatchId && !readySentRef.current) {
      readySentRef.current = true;
      sendReady();
      router.push(`/matches/${currentMatchId}`);
    }

    if (!currentMatchId) {
      readySentRef.current = false;
    }
  }, [currentMatchId, sendReady, router]);

  const handleFindMatch = () => {
    joinQueue();
  };

  return (
    <header className="flex h-16 items-center justify-between border-zinc-200 border-b bg-white px-4">
      <div className="flex items-center gap-3">
        <SidebarTrigger />
      </div>

      <button
        className="flex items-center gap-2 rounded-lg bg-zinc-900 px-4 py-2 font-medium text-sm text-white transition hover:bg-zinc-800 hover:shadow-sm hover:ring-zinc-800/20 active:scale-[0.98]"
        disabled={isQueued}
        onClick={handleFindMatch}
        type="button"
      >
        <Search className="size-4" />
        {isQueued ? "Searching..." : "Find Match"}
      </button>

      <Link
        aria-label="Profile"
        className="flex items-center gap-2 rounded-lg px-2 py-1.5 transition hover:bg-zinc-100"
        href="/dashboard"
        type="button"
      >
        <CircleUserRound className="size-5 text-zinc-600" />
        <span className="font-medium text-sm text-zinc-800">
          {user?.username}
        </span>
      </Link>
    </header>
  );
}

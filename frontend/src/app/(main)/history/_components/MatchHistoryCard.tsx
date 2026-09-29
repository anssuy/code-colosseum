"use client";

import Link from "next/link";
import { useEffect, useState } from "react";

import { User } from "@/lib/api/auth";
import { getMatch, Match, MatchData } from "@/lib/api/matches";

interface MatchHistoryCardProps {
  match: Match;
  user: User | null;
}

export default function MatchHistoryCard({
  match,
  user,
}: MatchHistoryCardProps) {
  const [data, setData] = useState<MatchData | undefined>(undefined);

  useEffect(() => {
    getMatch(match.id).then((data) => setData(data));
  }, [match]);

  if (!data || !user) return null;

  const isWinner = match.winnerId === user?.id;
  const isDraw = match.status === "finished" && !match.winnerId;

  const result =
    match.status === "finished"
      ? isWinner
        ? {
            label: "Won",
            className: "bg-green-50 text-green-700",
          }
        : isDraw
          ? {
              label: "Draw",
              className: "bg-zinc-100 text-zinc-600",
            }
          : {
              label: "Lost",
              className: "bg-red-50 text-red-700",
            }
      : match.status === "active"
        ? {
            label: "Active",
            className: "bg-blue-50 text-blue-700",
          }
        : {
            label: "Abandoned",
            className: "bg-zinc-100 text-zinc-500",
          };

  return (
    <Link
      className="group flex items-center justify-between gap-4 border-zinc-100 border-b px-5 py-4 last:border-b-0 hover:bg-zinc-50"
      href={`/history/${match.id}`}
      key={match.id}
    >
      <div className="flex min-w-0 items-center gap-4">
        <div className="flex size-10 shrink-0 items-center justify-center rounded-full bg-zinc-100 font-semibold text-sm text-zinc-600">
          VS
        </div>

        <div className="min-w-0">
          <div className="font-medium text-zinc-900">{data.problem.title}</div>
          <div className="mt-0.5 text-xs text-zinc-400">{match.status}</div>
        </div>
      </div>

      <div
        className={`shrink-0 rounded-full px-2.5 py-1 font-medium text-xs ${result.className}`}
      >
        {result.label}
      </div>
    </Link>
  );
}

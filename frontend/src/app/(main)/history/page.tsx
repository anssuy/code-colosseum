"use client";

import Link from "next/link";
import { useEffect, useState } from "react";

import { getMatches, Match } from "@/lib/api/matches";
import { useAuth } from "@/providers/AuthProvider";

export default function HistoryPage() {
  const { user } = useAuth();

  const [matches, setMatches] = useState<Match[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!user) return;

    getMatches()
      .then((response) => setMatches(response.data))
      .finally(() => setLoading(false));
  }, [user]);

  if (loading) {
    return <div className="p-6 text-zinc-500">Loading matches...</div>;
  }

  return (
    <div className="mx-auto max-w-4xl p-6">
      <h1 className="mb-4 font-semibold text-xl">Matches</h1>

      {matches.length === 0 ? (
        <div className="rounded-lg border border-zinc-200 p-6 text-center text-zinc-500">
          No matches found.
        </div>
      ) : (
        <div className="flex flex-col gap-2">
          {matches.map((match) => {
            const isPlayerOne = match.playerOneId === user?.id;
            const opponentId = isPlayerOne
              ? match.playerTwoId
              : match.playerOneId;

            const isWinner = match.winnerId === user?.id;
            const isDraw = match.status === "finished" && !match.winnerId;

            return (
              <Link
                className="flex items-center justify-between rounded-lg border border-zinc-200 p-4 transition hover:bg-zinc-50"
                href={`/matches/${match.id}`}
                key={match.id}
              >
                <div>
                  <div className="font-medium">vs {opponentId}</div>

                  <div className="text-sm text-zinc-500">{match.status}</div>
                </div>

                <div className="font-medium text-sm">
                  {match.status === "finished" &&
                    (isWinner ? (
                      <span className="text-green-600">Won</span>
                    ) : isDraw ? (
                      <span className="text-zinc-500">Draw</span>
                    ) : (
                      <span className="text-red-600">Lost</span>
                    ))}

                  {match.status === "active" && (
                    <span className="text-blue-600">Active</span>
                  )}

                  {match.status === "abandoned" && (
                    <span className="text-zinc-500">Abandoned</span>
                  )}
                </div>
              </Link>
            );
          })}
        </div>
      )}
    </div>
  );
}

"use client";

import { Trophy } from "lucide-react";
import { useEffect, useState } from "react";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { getLeaderboard, type LeaderboardEntry } from "@/lib/api/leaderboard";

export default function LeaderboardPage() {
  const [players, setPlayers] = useState<LeaderboardEntry[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    getLeaderboard()
      .then(setPlayers)
      .catch(() => setError("Failed to load leaderboard."))
      .finally(() => setLoading(false));
  }, []);

  return (
    <main className="mx-auto flex w-full max-w-6xl flex-col gap-6 p-6 text-zinc-900">
      <div>
        <h1 className="flex items-center gap-2 font-bold text-3xl">
          <Trophy className="size-7" />
          Leaderboard
        </h1>

        <p className="mt-1 text-zinc-500">
          Top players ranked by rating.
        </p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Top Players</CardTitle>
        </CardHeader>

        <CardContent>
          {loading ? (
            <p className="py-8 text-center text-zinc-500">
              Loading leaderboard...
            </p>
          ) : error ? (
            <p className="py-8 text-center text-red-500">{error}</p>
          ) : players.length === 0 ? (
            <p className="py-8 text-center text-zinc-500">
              No players found.
            </p>
          ) : (
            <div className="overflow-hidden rounded-lg border">
              <div className="grid grid-cols-[70px_1fr_120px_100px_100px] border-b bg-zinc-50 px-4 py-3 text-sm font-medium text-zinc-500">
                <span>Rank</span>
                <span>Player</span>
                <span>Rating</span>
                <span>Wins</span>
                <span>Losses</span>
              </div>

              {players.map((player) => (
                <div
                  key={player.username}
                  className="grid grid-cols-[70px_1fr_120px_100px_100px] items-center border-b px-4 py-4 last:border-0"
                >
                  <span className="font-semibold">#{player.rank}</span>
                  <span className="font-medium">{player.username}</span>
                  <span className="font-semibold">{player.rating}</span>
                  <span className="text-green-600">{player.wins}</span>
                  <span className="text-red-500">{player.losses}</span>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </main>
  );
}
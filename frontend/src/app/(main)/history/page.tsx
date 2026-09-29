"use client";

import Link from "next/link";
import { useEffect, useState } from "react";

import { getMatch, getMatches, Match } from "@/lib/api/matches";
import { useAuth } from "@/providers/AuthProvider";
import MatchHistoryCard from "./_components/MatchHistoryCard";

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
    <div className="mx-auto max-w-3xl p-6">
      <div className="mb-6">
        <h1 className="font-semibold text-2xl tracking-tight">Match history</h1>
        <p className="mt-1 text-sm text-zinc-500">
          Review your previous and active matches.
        </p>
      </div>

      {matches.length === 0 ? (
        <div className="rounded-xl border border-zinc-300 border-dashed p-10 text-center">
          <p className="font-medium text-zinc-700">No matches yet</p>
          <p className="mt-1 text-sm text-zinc-500">
            Your matches will appear here.
          </p>
        </div>
      ) : (
        <div className="overflow-hidden rounded-xl border border-zinc-200 bg-white">
          {matches.map((match) => (
            <MatchHistoryCard key={match.id} match={match} user={user} />
          ))}
        </div>
      )}
    </div>
  );
}

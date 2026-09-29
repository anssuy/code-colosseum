"use client";

import { ArrowLeft, Check, Clock, Code2, X } from "lucide-react";
import Link from "next/link";
import { useParams } from "next/navigation";
import { useEffect, useState } from "react";

import {
  getMatch,
  getSubmissions,
  MatchData,
  SubmissionResponse,
} from "@/lib/api/matches";
import { useAuth } from "@/providers/AuthProvider";

export default function MatchHistoryPage() {
  const { id } = useParams<{ id: string }>();
  const { user } = useAuth();

  const [data, setData] = useState<MatchData | null>(null);
  const [submissions, setSubmissions] = useState<SubmissionResponse[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    getMatch(id)
      .then(setData)
      .finally(() => setLoading(false));
    getSubmissions(id)
      .then(setSubmissions)
      .finally(() => setLoading(false));
  }, [id]);

  if (loading) {
    return <div className="p-6 text-zinc-500">Loading match...</div>;
  }

  if (!data) {
    return <div className="p-6 text-red-600">Match not found.</div>;
  }

  const { match, problem } = data;

  match.playerOneId === user?.id ? match.playerTwoId : match.playerOneId;

  const isWinner = match.winnerId === user?.id;
  const isDraw = match.status === "finished" && !match.winnerId;
  const isAbandoned = match.status === "abandoned";

  const result = isAbandoned
    ? "Abandoned"
    : isDraw
      ? "Draw"
      : isWinner
        ? "Victory"
        : "Defeat";

  const resultClass = isAbandoned
    ? "bg-zinc-100 text-zinc-600"
    : isDraw
      ? "bg-zinc-100 text-zinc-700"
      : isWinner
        ? "bg-green-50 text-green-700"
        : "bg-red-50 text-red-700";

  return (
    <main className="mx-auto max-w-4xl p-6">
      <Link
        className="mb-8 inline-flex items-center gap-2 text-sm text-zinc-500 transition hover:text-zinc-900"
        href="/history"
      >
        <ArrowLeft className="size-4" />
        Match history
      </Link>

      <section className="mb-8">
        <div className="mb-6 flex items-start justify-between gap-4">
          <div>
            <p className="mb-1 text-sm text-zinc-500">Match result</p>
            <h1 className="font-semibold text-3xl tracking-tight">
              {data.problem.title}
            </h1>
            <p className="mt-2 text-sm text-zinc-500">{problem.title}</p>
          </div>

          <span
            className={`rounded-full px-3 py-1.5 font-medium text-sm ${resultClass}`}
          >
            {result}
          </span>
        </div>

        <div className="grid grid-cols-3 gap-3">
          <div className="rounded-xl bg-zinc-50 px-4 py-3">
            <div className="text-xs text-zinc-400">Slug</div>
            <div className="mt-1 font-medium text-sm">{problem.slug}</div>
          </div>

          <div className="rounded-xl bg-zinc-50 px-4 py-3">
            <div className="text-xs text-zinc-400">Difficulty</div>
            <div className="mt-1 font-medium text-sm capitalize">
              {problem.difficulty}
            </div>
          </div>

          <div className="rounded-xl bg-zinc-50 px-4 py-3">
            <div className="text-xs text-zinc-400">Status</div>
            <div className="mt-1 font-medium text-sm capitalize">
              {match.status}
            </div>
          </div>
        </div>
      </section>

      <section className="overflow-hidden rounded-xl bg-white shadow-sm ring-1 ring-zinc-200/70">
        <div className="flex items-center justify-between px-5 py-4">
          <div>
            <h2 className="font-semibold">Submissions</h2>
            <p className="mt-0.5 text-sm text-zinc-500">
              Your attempts during this match
            </p>
          </div>

          <span className="text-sm text-zinc-400">
            {submissions.length}{" "}
            {submissions.length === 1 ? "attempt" : "attempts"}
          </span>
        </div>

        {submissions.length === 0 ? (
          <div className="border-zinc-100 border-t px-5 py-10 text-center">
            <Code2 className="mx-auto mb-3 size-8 text-zinc-300" />
            <p className="font-medium text-sm text-zinc-600">No submissions</p>
            <p className="mt-1 text-sm text-zinc-400">
              No code was submitted during this match.
            </p>
          </div>
        ) : (
          <div className="divide-y divide-zinc-100">
            {submissions.map((submission, index) => {
              const accepted = submission.status === "accepted";

              return (
                <div
                  className="flex items-center justify-between gap-4 px-5 py-4 transition hover:bg-zinc-50"
                  key={submission.id}
                >
                  <div className="flex min-w-0 items-center gap-3">
                    <div
                      className={`flex size-9 shrink-0 items-center justify-center rounded-full ${
                        accepted
                          ? "bg-green-50 text-green-600"
                          : "bg-red-50 text-red-600"
                      }`}
                    >
                      {accepted ? (
                        <Check className="size-4" />
                      ) : (
                        <X className="size-4" />
                      )}
                    </div>

                    <div className="min-w-0">
                      <div className="font-medium text-sm">
                        Attempt {index + 1}
                      </div>

                      <div className="mt-0.5 flex items-center gap-2 text-xs text-zinc-400">
                        <Code2 className="size-3" />
                        <span className="capitalize">
                          {submission.language}
                        </span>
                      </div>
                    </div>
                  </div>

                  <div className="flex shrink-0 items-center gap-6 text-right">
                    <div>
                      <div className="font-medium text-sm">
                        {submission.passedTests}/{submission.totalTests}
                      </div>
                      <div className="text-xs text-zinc-400">tests passed</div>
                    </div>

                    <div
                      className={`min-w-20 font-medium text-sm ${
                        accepted ? "text-green-600" : "text-red-600"
                      }`}
                    >
                      {submission.status}
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </section>

      <div className="mt-4 flex items-center gap-2 text-xs text-zinc-400">
        <Clock className="size-3.5" />
        Match ID: {match.id}
      </div>
    </main>
  );
}

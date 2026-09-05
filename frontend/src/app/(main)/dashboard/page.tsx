"use client";

import { useAuth } from "@/providers/AuthProvider";

export default function DashboardPage() {
  const { user, loading } = useAuth();

  if (loading || !user) {
    return (
      <p className="min-h-screen p-8 text-center text-zinc-900">Loading...</p>
    );
  }

  return (
    <main className="flex items-center justify-center p-6 text-zinc-900">
      <section className="w-full max-w-md rounded-xl border border-zinc-200 bg-white p-8 shadow-lg">
        <h1 className="mb-6 font-bold text-2xl">Welcome, {user.username}</h1>

        <div className="mb-6 space-y-2 text-zinc-700">
          <p>Email: {user.email}</p>
          <p>Rating: {user.rating}</p>
          <p>Wins: {user.wins}</p>
          <p>Losses: {user.losses}</p>
        </div>
      </section>
    </main>
  );
}

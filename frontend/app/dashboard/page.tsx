"use client";

import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

import { useAuth } from "@/providers/AuthProvider";

export default function DashboardPage() {
  const router = useRouter();
  const { user, loading, logout } = useAuth();
  const [loggingOut, setLoggingOut] = useState(false);

  useEffect(() => {
    if (!loading && !user) router.replace("/login");
  }, [loading, user, router]);

  async function handleLogout() {
    setLoggingOut(true);

    try {
      await logout();
      router.replace("/login");
    } finally {
      setLoggingOut(false);
    }
  }

  if (loading || !user) {
    return (
      <p className="min-h-screen bg-zinc-100 p-8 text-center text-zinc-900">
        Loading...
      </p>
    );
  }

  return (
    <main className="flex min-h-screen items-center justify-center bg-zinc-100 p-6 text-zinc-900">
      <section className="w-full max-w-md rounded-xl border border-zinc-200 bg-white p-8 shadow-lg">
        <h1 className="mb-6 font-bold text-2xl">Welcome, {user.username}</h1>

        <div className="mb-6 space-y-2 text-zinc-700">
          <p>Email: {user.email}</p>
          <p>Rating: {user.rating}</p>
          <p>Wins: {user.wins}</p>
          <p>Losses: {user.losses}</p>
        </div>

        <button
          className="w-full rounded-md bg-zinc-900 px-4 py-2.5 font-semibold text-white hover:bg-zinc-800 disabled:opacity-60"
          disabled={loggingOut}
          onClick={handleLogout}
          type="button"
        >
          {loggingOut ? "Logging out..." : "Logout"}
        </button>
      </section>
    </main>
  );
}

"use client";

import { Flame, Medal, TrendingUp, Trophy } from "lucide-react";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useAuth } from "@/providers/AuthProvider";

export default function DashboardPage() {
  const { user, loading } = useAuth();

  if (loading || !user) {
    return (
      <p className="min-h-screen p-8 text-center text-zinc-900">Loading...</p>
    );
  }

  const totalMatches = user.wins + user.losses;
  const winRate =
    totalMatches > 0 ? Math.round((user.wins / totalMatches) * 100) : 0;

  const stats = [
    {
      label: "Rating",
      value: user.rating,
      icon: TrendingUp,
    },
    {
      label: "Wins",
      value: user.wins,
      icon: Trophy,
    },
    {
      label: "Losses",
      value: user.losses,
      icon: Medal,
    },
    {
      label: "Win Rate",
      value: `${winRate}%`,
      icon: Flame,
    },
  ];

  return (
    <main className="mx-auto flex w-full max-w-6xl flex-col gap-6 p-6 text-zinc-900">
      <Card>
        <CardContent className="flex items-center justify-between p-8">
          <div className="space-y-2">
            <p className="font-medium text-sm text-zinc-500">
              Ready to compete?
            </p>

            <h1 className="font-bold text-3xl">Find your next opponent</h1>

            <p className="text-zinc-600">
              Match against another player and put your rating on the line.
            </p>
          </div>
        </CardContent>
      </Card>

      <section className="grid grid-cols-2 gap-4 lg:grid-cols-4">
        {stats.map(({ label, value, icon: Icon }) => (
          <Card key={label}>
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <CardTitle className="font-medium text-sm text-zinc-500">
                {label}
              </CardTitle>
              <Icon className="size-4 text-zinc-500" />
            </CardHeader>

            <CardContent>
              <p className="font-bold text-2xl">{value}</p>
            </CardContent>
          </Card>
        ))}
      </section>
    </main>
  );
}

"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { type SubmitEvent, useState } from "react";

import { useAuth } from "@/providers/AuthProvider";

export default function RegisterPage() {
  const router = useRouter();
  const { register } = useAuth();

  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);

  async function handleSubmit(event: SubmitEvent<HTMLFormElement>) {
    event.preventDefault();
    setError("");
    setSubmitting(true);

    try {
      await register(username, email, password);
      router.replace("/login");
    } catch (error) {
      setError(error instanceof Error ? error.message : "Registration failed");
    } finally {
      setSubmitting(false);
    }
  }

  const inputClass =
    "rounded-md border border-zinc-300 px-3 py-2 outline-none focus:border-zinc-600";

  return (
    <main className="flex min-h-screen items-center justify-center bg-zinc-100 p-6 text-zinc-900">
      <section className="w-full max-w-sm rounded-xl border border-zinc-200 bg-white p-8 shadow-lg">
        <h1 className="mb-6 text-center font-bold text-2xl">Register</h1>

        <form className="flex flex-col gap-5" onSubmit={handleSubmit}>
          <div className="flex flex-col gap-2">
            <label className="font-semibold text-sm" htmlFor="username">
              Username
            </label>
            <input
              className={inputClass}
              id="username"
              maxLength={30}
              minLength={3}
              onChange={(event) => setUsername(event.target.value)}
              required
              value={username}
            />
          </div>

          <div className="flex flex-col gap-2">
            <label className="font-semibold text-sm" htmlFor="email">
              Email
            </label>
            <input
              autoComplete="email"
              className={inputClass}
              id="email"
              name="email"
              onChange={(event) => setEmail(event.target.value)}
              required
              type="email"
              value={email}
            />
          </div>

          <div className="flex flex-col gap-2">
            <label className="font-semibold text-sm" htmlFor="password">
              Password
            </label>
            <input
              autoComplete="new-password"
              className={inputClass}
              id="password"
              maxLength={72}
              minLength={8}
              name="password"
              onChange={(event) => setPassword(event.target.value)}
              required
              type={showPassword ? "text" : "password"}
              value={password}
            />

            <button
              className="self-start font-semibold text-[13px] text-blue-600 hover:cursor-pointer"
              onClick={() => setShowPassword(!showPassword)}
              type="button"
            >
              {showPassword ? "Hide password" : "Show password"}
            </button>
          </div>

          {error && <p className="text-red-600 text-sm">{error}</p>}

          <button
            className="rounded-md bg-zinc-900 px-4 py-2.5 font-semibold text-white hover:bg-zinc-800 disabled:cursor-not-allowed disabled:opacity-60"
            disabled={submitting}
            type="submit"
          >
            {submitting ? "Creating account..." : "Register"}
          </button>
        </form>

        <p className="mt-6 text-center text-sm text-zinc-600">
          Already registered?{" "}
          <Link
            className="font-semibold text-blue-600 hover:underline"
            href="/login"
          >
            Login
          </Link>
        </p>
      </section>
    </main>
  );
}

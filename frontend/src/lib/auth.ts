import { ApiError, apiFetch } from "./api";

export type User = {
  id: string;
  username: string;
  email: string;
  rating: number;
  wins: number;
  losses: number;
  created_at: string;
};

type AuthResponse = { user: User };

export async function register(
  username: string,
  email: string,
  password: string,
) {
  const response = await apiFetch<AuthResponse>(
    "/api/auth/register",
    {
      method: "POST",
      body: JSON.stringify({ username, email, password }),
    },
    false,
  );

  return response.user;
}

export async function login(email: string, password: string) {
  const response = await apiFetch<AuthResponse>(
    "/api/auth/login",
    {
      method: "POST",
      body: JSON.stringify({ email, password }),
    },
    false,
  );

  return response.user;
}

export async function getCurrentUser() {
  try {
    return (await apiFetch<AuthResponse>("/api/auth/me")).user;
  } catch (error) {
    if (error instanceof ApiError && error.status === 401) return null;
    throw error;
  }
}

export async function logout() {
  await apiFetch<void>("/api/auth/logout", { method: "POST" }, false);
}

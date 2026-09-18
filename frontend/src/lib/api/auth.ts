import { ApiError, apiFetch } from "./index";

export type User = {
  id: string;
  username: string;
  email: string;
  rating: number;
  wins: number;
  losses: number;
  created_at: string;
};

type AuthResponse = {
  user: User;
};

export async function register(
  username: string,
  email: string,
  password: string,
) {
  const { user } = await apiFetch<AuthResponse>(
    "/auth/register",
    {
      method: "POST",
      body: JSON.stringify({ username, email, password }),
    },
    false,
  );

  return user;
}

export async function login(email: string, password: string) {
  const { user } = await apiFetch<AuthResponse>(
    "/auth/login",
    {
      method: "POST",
      body: JSON.stringify({ email, password }),
    },
    false,
  );

  return user;
}

export async function getCurrentUser() {
  try {
    const { user } = await apiFetch<AuthResponse>("/auth/me");
    return user;
  } catch (error) {
    if (error instanceof ApiError && error.status === 401) {
      return null;
    }

    throw error;
  }
}

export async function logout() {
  await apiFetch<void>("/auth/logout", { method: "POST" }, false);
}

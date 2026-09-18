const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080/api";

export class ApiError extends Error {
  constructor(
    message: string,
    public readonly status: number,
  ) {
    super(message);
    this.name = "ApiError";
  }
}

async function readResponse(response: Response) {
  if (response.status === 204) return null;
  return response.json().catch(() => null);
}

let refreshPromise: Promise<boolean> | null = null;

const refreshSession = async () => {
  if (!refreshPromise) {
    refreshPromise = fetch(`${API_URL}/auth/refresh`, {
      method: "POST",
      credentials: "include",
    })
      .then((response) => response.ok)
      .finally(() => {
        refreshPromise = null;
      });
  }

  return refreshPromise;
};

export async function apiFetch<T>(
  path: string,
  options: RequestInit = {},
  retryOnUnauthorized = true,
): Promise<T> {
  const headers = new Headers(options.headers);

  if (options.body && !headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json");
  }

  const response = await fetch(`${API_URL}${path}`, {
    ...options,
    headers,
    credentials: "include",
  });

  if (response.status === 401 && retryOnUnauthorized) {
    if (await refreshSession()) {
      return apiFetch<T>(path, options, false);
    }
  }

  const data = await readResponse(response);

  if (!response.ok) {
    const error = data as { error?: string } | null;

    throw new ApiError(error?.error ?? "Request failed", response.status);
  }

  return data as T;
}

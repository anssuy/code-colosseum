const API_URL = process.env.API_URL ?? "http://localhost:8080/api";

export const login = async (): Promise<string> => {
  const response = await fetch(`${API_URL}/auth/login`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      email: "test@gmail.com",
      password: "Test123!",
    }),
  });

  if (!response.ok) throw new Error(`Login failed: ${response.status}`);

  const cookie = response.headers.get("set-cookie");

  if (!cookie) throw new Error("Missing session cookie");

  return cookie.split(";")[0] as string;
};

export const apiFetch = <T>(
  path: string,
  cookie: string,
  options: RequestInit = {},
): Promise<T> => {
  return fetch(`${API_URL}${path}`, {
    ...options,
    headers: {
      ...options.headers,
      Cookie: cookie,
    },
  })
    .then((res) => res.json())
    .then((data) => data as T);
};

import { apiFetch } from "./api";

export type RunResponse = { output: string };

export async function runCode(code: string, language: string) {
  const response = await apiFetch<RunResponse>("/api/sandbox/run", {
    method: "POST",
    body: JSON.stringify({ code, language }),
  });

  return response.output;
}

"use client";

import { Editor } from "@monaco-editor/react";
import { type SubmitEvent, useState } from "react";

import { runCode } from "@/lib/sandbox";

const languages = [
  { value: "javascript", label: "JavaScript" },
  { value: "typescript", label: "TypeScript" },
  { value: "python", label: "Python" },
];

export default function SandboxPage() {
  const [language, setLanguage] = useState("javascript");
  const [code, setCode] = useState('console.log("Hello!");');
  const [output, setOutput] = useState("");
  const [error, setError] = useState("");
  const [running, setRunning] = useState(false);

  async function handleSubmit(event: SubmitEvent<HTMLFormElement>) {
    event.preventDefault();
    setError("");
    setOutput("");
    setRunning(true);

    try {
      const result = await runCode(code, language);
      setOutput(result);
    } catch (error) {
      setError(error instanceof Error ? error.message : "Run failed");
    } finally {
      setRunning(false);
    }
  }

  return (
    <main className="flex min-h-screen items-start justify-center bg-zinc-100 p-6 pt-10 text-zinc-900">
      <section className="w-full max-w-4xl rounded-xl border border-zinc-200 bg-white p-8 shadow-lg">
        <h1 className="mb-6 text-center font-bold text-2xl">Sandbox</h1>

        <form className="flex flex-col gap-5" onSubmit={handleSubmit}>
          <div className="flex flex-col gap-2">
            <label className="font-semibold text-sm" htmlFor="language">
              Language
            </label>
            <select
              className="rounded-md border border-zinc-300 px-3 py-2 outline-none focus:border-zinc-600"
              id="language"
              onChange={(event) => setLanguage(event.target.value)}
              value={language}
            >
              {languages.map((lang) => (
                <option key={lang.value} value={lang.value}>
                  {lang.label}
                </option>
              ))}
            </select>
          </div>

          <div className="flex flex-col gap-2">
            <label className="sr-only" htmlFor="code">
              Add your code
            </label>
            <Editor
              height="50vh"
              language={language}
              onChange={(value) => setCode(value ?? "")}
              value={code}
            />
          </div>

          {error && <p className="text-red-600 text-sm">{error}</p>}

          <button
            className="rounded-md bg-zinc-900 px-4 py-2.5 font-semibold text-white hover:bg-zinc-800 disabled:cursor-not-allowed disabled:opacity-60"
            disabled={running}
            type="submit"
          >
            {running ? "Running..." : "Run"}
          </button>
        </form>

        {output && (
          <pre className="mt-6 whitespace-pre-wrap rounded-md bg-zinc-900 p-4 text-sm text-zinc-100">
            {output}
          </pre>
        )}
      </section>
    </main>
  );
}

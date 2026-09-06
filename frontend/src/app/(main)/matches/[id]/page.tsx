"use client";

import Editor from "@monaco-editor/react";
import { useParams } from "next/navigation";
import { useEffect, useState } from "react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";

import { apiFetch } from "@/lib/api";
import { LANGUAGES, LanguageValue } from "@/lib/constants";
import { useMatch } from "@/providers/MatchProvider";

type MatchData = {
  match: {
    id: string;
    status: string;
    playerOneId: string;
    playerTwoId: string;
    winnerId: string | null;
  };
  problem: {
    id: string;
    title: string;
    slug: string;
    difficulty: string;
    description: string;
  };
};

type ProblemResponse = {
  problem: {
    id: string;
    title: string;
    slug: string;
    difficulty: string;
    description: string;
    functionName: string;
    stubs: Record<LanguageValue, string>;
  };
  testCases: TestCase[];
};

type TestCase = {
  id: string;
  input: string;
  expectedOutput: string;
  isSample: boolean;
};

export default function MatchPage() {
  const { id } = useParams<{ id: string }>();
  const { matchStarted, lastResult, winnerId, abandoned, submit } = useMatch();

  const [data, setData] = useState<MatchData | null>(null);
  const [testCases, setTestCases] = useState<TestCase[]>([]);
  const [stubs, setStubs] = useState<Record<LanguageValue, string> | null>(
    null,
  );
  const [language, setLanguage] = useState<LanguageValue>("python");
  const [code, setCode] = useState("");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    apiFetch<MatchData>(`/api/matches/${id}`)
      .then((matchData) => {
        setData(matchData);
        return apiFetch<ProblemResponse>(
          `/api/problems/slug/${matchData.problem.slug}`,
        );
      })
      .then((p) => {
        console.log(p);
        setTestCases(p.testCases);
        setStubs(p.problem.stubs);
        setCode(p.problem.stubs.python);
      })
      .finally(() => setLoading(false));
  }, [id]);

  const handleLanguageChange = (newLanguage: string) => {
    setLanguage(newLanguage);
    if (stubs) {
      setCode(stubs[newLanguage]);
    }
  };

  const handleSubmit = () => {
    submit(language, code);
  };

  if (loading) {
    return <div className="p-6 text-zinc-500">Loading match...</div>;
  }

  if (!data) {
    return <div className="p-6 text-red-600">Match not found.</div>;
  }

  const isActive = matchStarted || data.match.status === "active";
  const isFinished = !!winnerId || data.match.status === "finished";
  const isAbandoned = abandoned || data.match.status === "abandoned";
  const matchWinnerId = winnerId || data.match.winnerId;

  return (
    <div className="grid h-[calc(100vh-4rem)] grid-cols-2 gap-4 p-4">
      <div className="flex flex-col gap-4 overflow-y-auto rounded-lg border border-zinc-200 p-4">
        <div>
          <h1 className="font-semibold text-lg">{data.problem.title}</h1>
          <span className="text-sm text-zinc-500">
            {data.problem.difficulty}
          </span>
        </div>

        <div className="prose prose-sm max-w-none">
          <ReactMarkdown remarkPlugins={[remarkGfm]}>
            {data.problem.description}
          </ReactMarkdown>
        </div>

        {testCases.length > 0 && (
          <div className="flex flex-col gap-2">
            <h2 className="font-medium text-sm text-zinc-700">
              Sample Test Cases
            </h2>
            {testCases.map((tc, i) => (
              <div
                className="rounded-lg border border-zinc-200 p-3 text-sm"
                key={tc.id}
              >
                <div className="mb-2 font-medium text-zinc-500">
                  Example {i + 1}
                </div>
                <div className="grid grid-cols-2 gap-3">
                  <div>
                    <div className="text-xs text-zinc-400">Input</div>
                    <pre className="whitespace-pre-wrap rounded bg-zinc-50 p-2 font-mono text-xs">
                      {tc.input}
                    </pre>
                  </div>
                  <div>
                    <div className="text-xs text-zinc-400">Output</div>
                    <pre className="whitespace-pre-wrap rounded bg-zinc-50 p-2 font-mono text-xs">
                      {tc.expectedOutput}
                    </pre>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}

        <div className="mt-auto rounded-lg bg-zinc-50 p-3 text-sm">
          {!isActive && !isFinished && !isAbandoned && (
            <span className="text-zinc-500">Waiting for opponent...</span>
          )}
          {isActive && !isFinished && !isAbandoned && (
            <span className="text-green-600">Match in progress</span>
          )}
          {isFinished && (
            <span className="font-medium">
              Match finished — winner: {matchWinnerId}
            </span>
          )}
          {isAbandoned && (
            <span className="text-zinc-500">Match abandoned</span>
          )}
        </div>

        {lastResult && (
          <div className="rounded-lg border border-zinc-200 p-3 text-sm">
            <div className="font-medium">{lastResult.status}</div>
            <div className="text-zinc-500">
              {lastResult.passedTests}/{lastResult.totalTests} tests passed
            </div>
          </div>
        )}
      </div>

      <div className="flex flex-col gap-3">
        <div className="flex items-center justify-between">
          <select
            className="rounded-lg border border-zinc-200 px-3 py-1.5 text-sm"
            onChange={(e) => handleLanguageChange(e.target.value)}
            value={language}
          >
            {LANGUAGES.map((l) => (
              <option key={l.value} value={l.value}>
                {l.label}
              </option>
            ))}
          </select>

          <button
            className="rounded-lg bg-zinc-900 px-4 py-1.5 font-medium text-sm text-white transition hover:bg-zinc-700 disabled:opacity-50"
            disabled={isFinished || isAbandoned}
            onClick={handleSubmit}
            type="button"
          >
            Submit
          </button>
        </div>

        <div className="flex-1 overflow-hidden rounded-lg border border-zinc-200">
          <Editor
            language={language}
            onChange={(v) => setCode(v ?? "")}
            options={{ minimap: { enabled: false }, fontSize: 14 }}
            theme="vs-light"
            value={code}
          />
        </div>
      </div>
    </div>
  );
}

import { LanguageValue } from "../constants";
import { apiFetch } from "./index";

export type TestCase = {
  id: string;
  input: string;
  expectedOutput: string;
  isSample: boolean;
};

export type Problem = {
  id: string;
  title: string;
  slug: string;
  difficulty: string;
  description: string;
  functionName: string;
  stubs: Record<LanguageValue, string>;
};

type ProblemResponse = {
  problem: Problem;
  testCases: TestCase[];
};

export const getProblem = (slug: string) =>
  apiFetch<ProblemResponse>(`/problems/slug/${slug}`);

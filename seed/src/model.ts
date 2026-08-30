import { InferenceClient } from "@huggingface/inference";

import type { Param, TestCase } from "./schema";

const client = new InferenceClient(process.env.HF_TOKEN);
const MODEL = "Qwen/Qwen2.5-Coder-32B-Instruct";

interface Signature {
  functionName: string;
  params: Param[];
  returnType: string;
}

const stripFences = (raw: string) => raw.replace(/```json\n?|```/g, "").trim();

const callModel = async (prompt: string, retries = 2) => {
  const response = await client.chatCompletion({
    model: MODEL,
    messages: [{ role: "user", content: prompt }],
    max_tokens: 8192,
  });

  const choice = response.choices[0];
  const clean = stripFences(choice?.message.content ?? "");

  try {
    return JSON.parse(clean);
  } catch {
    if (retries > 0) {
      console.warn("Invalid JSON, retrying...", clean.slice(-100));
      return callModel(prompt, retries - 1);
    }
    throw new Error(`Model returned invalid JSON after retries: ${clean}`);
  }
};

export const inferSignature = async (
  title: string,
  description: string,
): Promise<Signature> => {
  return callModel(`Given this problem, infer a function signature (LeetCode-style, camelCase).

Return ONLY JSON, no fences:
{"functionName":"...", "params":[{"name":"...", "type":"..."}], "returnType":"..."}

type and returnType must be one of these:
- int
- str
- bool
- int[]
- str[]
- bool[]

Title: ${title}
Description: ${description}`);
};

const buildTestCasePrompt = (description: string, signature: Signature) => {
  const sig = `${signature.functionName}(${signature.params
    .map((p) => `${p.name}: ${p.type}`)
    .join(", ")}) -> ${signature.returnType}`;

  const paramTypesList = signature.params
    .map((p) => `  - ${p.name}: ${p.type}`)
    .join("\n");

  return `Generate test cases for this programming problem.

Function signature: ${sig}

Param types:
${paramTypesList}
Return type: ${signature.returnType}

CRITICAL "input" format:
- Input is ALWAYS a single JSON array, positional, one entry per param, in declared order.
- NEVER use object/dict form like {"paramName": value}. Only array form.
- int[] param -> nested array: [2,7,11,15]
- string[] param -> nested array: ["a","bb","acd"]
- bool[] param -> nested array: [true,false,true]
- int param -> plain number: 9
- string param -> plain string: "abc"
- bool param -> plain boolean: true

Correct example, params (nums: int[], target: int):
"input": "[[2,7,11,15], 9]"

WRONG (do not do this):
"input": "{\\"nums\\":[2,7,11,15], \\"target\\":9}"

"expectedOutput" is the return value only, stringified directly, no wrapping object:
"expectedOutput": "[0,1]"

Rules:
- Generate exactly 6 test cases total (2 sample, 4 hidden).
- Include edge cases (empty, single element, duplicates, boundary values).
- No trailing commas in JSON.

Return ONLY valid JSON, no markdown fences, no explanation:
{"testCases":[{"input":"...", "expectedOutput":"...", "isSample": true}]}

Problem:
${description}`;
};

export const generateTestCases = async (
  description: string,
  signature: Signature,
): Promise<TestCase[]> => {
  const result = await callModel(buildTestCasePrompt(description, signature));
  return result.testCases;
};

export const processScrapedProblem = async (problem: {
  title: string;
  description: string;
}) => {
  const signature = await inferSignature(problem.title, problem.description);
  const testCases = await generateTestCases(problem.description, signature);

  return {
    functionName: signature.functionName,
    params: signature.params,
    returnType: signature.returnType,
    testCases,
  };
};

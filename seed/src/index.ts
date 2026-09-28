import "dotenv/config";

import { readFile } from "node:fs/promises";
import path from "node:path";

import { apiFetch, login } from "./api";
import { getRandomProblem, getRandomProblems } from "./leetcode";
import { processScrapedProblem } from "./model";
import type { Problem, ProblemTag, TestCase } from "./schema";

type ProblemResponse = {
  problem: Problem & { id: string };
};

type TagResponse = {
  tag: ProblemTag & { id: string };
};

const main = async () => {
  const cookie = await login();

  const command = process.argv[2];

  switch (command) {
    case "list": {
      const problems = await getRandomProblems();
      console.log(problems);
      break;
    }
    case "data": {
      const generated: {
        problem: Problem;
        tags: ProblemTag[];
        testCases: TestCase[];
      }[] = JSON.parse(
        await readFile(path.join(__dirname, "..", "data/data.json"), {
          encoding: "utf-8",
        }),
      );

      for (const { problem, tags, testCases } of generated) {
        try {
          const {
            problem: { id: problemId },
          } = await apiFetch<ProblemResponse>("/problems", cookie, {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify(problem),
          });

          tags.forEach(async (tag) => {
            try {
              const {
                tag: { id: tagId },
              } = await apiFetch<TagResponse>("/tags", cookie, {
                method: "POST",
                headers: {
                  "Content-Type": "application/json",
                },
                body: JSON.stringify(tag),
              });

              await apiFetch(`/problems/${problemId}/tags/${tagId}`, cookie, {
                method: "POST",
                headers: {
                  "Content-Type": "application/json",
                },
                body: JSON.stringify(tag),
              });
            } catch (err) {
              console.error(err);
            }
          });

          testCases.forEach(async (testCase) => {
            await apiFetch(`/problems/${problemId}/testcases`, cookie, {
              method: "POST",
              headers: {
                "Content-Type": "application/json",
              },
              body: JSON.stringify(testCase),
            });
          });
        } catch (err) {
          console.error(err);
        }
      }
      break;
    }
    default: {
      const problem = await getRandomProblem();

      const { functionName, params, returnType, testCases } =
        await processScrapedProblem({
          title: problem.title,
          description: problem.description,
        });

      problem.functionName = functionName;
      problem.params = params;
      problem.returnType = returnType;

      try {
        const {
          problem: { id },
        } = await apiFetch<ProblemResponse>("/problems", cookie, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(problem),
        });

        testCases.forEach(async (testCase) => {
          await apiFetch(`/problems/${id}/testcases`, cookie, {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify(testCase),
          });
        });
      } catch {}
    }
  }
};

main().catch((error) => {
  console.error(error);
  process.exit(1);
});

import "dotenv/config";

import { apiFetch, login } from "./api";
import { getRandomProblem, getRandomProblems } from "./leetcode";
import { processScrapedProblem } from "./model";
import type { Problem } from "./schema";

type ProblemResponse = {
  problem: Problem & { id: string };
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

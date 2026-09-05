import { Leetcode } from "@codingsnack/leetcode-api";

import { htmlToMarkdown } from "./markdown";
import { type Problem, ProblemSchema } from "./schema";

const csrfToken = process.env.LC_CSRF_TOKEN;
const session = process.env.LC_SESSION;

if (!csrfToken || !session) throw new Error("Missing LeetCode credentials");

const lc = new Leetcode({
  csrfToken,
  session,
});

export const getRandomProblem = async (): Promise<Problem> => {
  const problem = await lc.getRandomQuestion();

  return ProblemSchema.parse({
    title: problem.title,
    slug: problem.titleSlug,
    difficulty: problem.difficulty,
    description: htmlToMarkdown(problem.content),
    tags: problem.topicTags.map((tag) => ({
      slug: tag.slug,
      name: tag.name,
    })),
    functionName: "",
    params: [],
    returnType: "",
  });
};

export const getRandomProblems = async () => {
  const problems = await lc.getProblems({
    skip: Math.floor(Math.random() * 100),
    filters: {
      premiumOnly: false,
    },
    categorySlug: "algorithms",
    limit: 10,
  });
  return problems.questions.map((problem) => {
    return {
      title: problem.title,
      slug: problem.titleSlug,
      description: htmlToMarkdown(problem.content),
    };
  });
};

import z from "zod";

export const ParamSchema = z.object({
  name: z.string(),
  type: z.string(),
});

export const ProblemTagSchema = z.object({
  slug: z.string(),
  name: z.string(),
});

export const TestCaseSchema = z.object({
  input: z.string(),
  expectedOutput: z.string(),
  isSample: z.boolean(),
});

export const ProblemSchema = z.object({
  title: z.string(),
  slug: z.string(),
  difficulty: z.string().toLowerCase(),
  description: z.string(),
  tags: z.array(ProblemTagSchema),
  functionName: z.string(),
  params: z.array(ParamSchema),
  returnType: z.string(),
});

export type Problem = z.infer<typeof ProblemSchema>;
export type ProblemTag = z.infer<typeof ProblemTagSchema>;
export type TestCase = z.infer<typeof TestCaseSchema>;
export type Param = z.infer<typeof ParamSchema>;

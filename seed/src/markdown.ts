import TurndownService from "turndown";

const turndown = new TurndownService({
  headingStyle: "atx",
  bulletListMarker: "-",
  codeBlockStyle: "fenced",
});

const preprocessExamples = (html: string): string =>
  html.replace(/<pre>([\s\S]*?)<\/pre>/gi, (_, content: string) => {
    const normalized = content
      .replace(/<strong>\s*Input:\s*<\/strong>/gi, "Input:")
      .replace(/<strong>\s*Output:\s*<\/strong>/gi, "Output:")
      .replace(/<strong>\s*Explaination:\s*<\/strong>/gi, "Explanation:")
      .replace(/<strong>\s*Explanation:\s*<\/strong>/gi, "Explanation:");

    return `<pre>${normalized}</pre>`;
  });

export const htmlToMarkdown = (html: string): string => {
  const preprocessed = preprocessExamples(
    html
      .replace(/<p>\s*(?:&nbsp;|\s)*<\/p>/gi, "")
      .replace(
        /<p>\s*<strong[^>]*class=["']example["'][^>]*>\s*Example\s+(\d+):\s*<\/strong>\s*<\/p>/gi,
        "<h3>Example $1</h3>",
      )
      .replace(
        /<p>\s*<strong>\s*Constraints:\s*<\/strong>\s*<\/p>/gi,
        "<h3>Constraints</h3>",
      )
      .replace(/Explaination/gi, "Explanation")
      .trim(),
  );

  return turndown
    .turndown(preprocessed)
    .replace(/\n{3,}/g, "\n\n")
    .trim();
};

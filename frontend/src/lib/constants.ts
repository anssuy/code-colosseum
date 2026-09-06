type Language = {
  value: string;
  label: string;
};

export const LANGUAGES: Language[] = [
  { value: "java", label: "Java" },
  { value: "python", label: "Python" },
  { value: "javascript", label: "JavaScript" },
  { value: "typescript", label: "TypeScript" },
  { value: "cpp", label: "C++" },
];

export type LanguageValue = Language["value"];

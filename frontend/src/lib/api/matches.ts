import { apiFetch } from "./index";

export type SubmissionResult = {
  passedTests: number;
  status: string;
  submissionId: string;
  totalTests: number;
  userId: string;
};

export type Match = {
  id: string;
  status: string;
  playerOneId: string;
  playerTwoId: string;
  winnerId: string | null;
};

export type MatchData = {
  match: Match;
  problem: {
    id: string;
    title: string;
    slug: string;
    difficulty: string;
    description: string;
  };
};

type MatchesResponse = {
  data: Match[];
  totalCount: number;
  limit: number;
  offset: number;
};

export const getMatch = (id: string) => apiFetch<MatchData>(`/matches/${id}`);

export const getActiveMatch = () => apiFetch<MatchData>(`/matches/active`);

export const getMatches = (limit = 20, offset = 0) =>
  apiFetch<MatchesResponse>(`/matches?limit=${limit}&offset=${offset}`);

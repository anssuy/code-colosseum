import {apiFetch} from "@/lib/api/index";

export type LeaderboardEntry = {
    rank: number;
    username: string;
    rating: number;
    wins: number;
    losses: number;
};

export const getLeaderboard = (limit = 20, offset = 0) =>
    apiFetch<LeaderboardEntry[]>(
        `/leaderboard?limit=${limit}&offset=${offset}`,
    );
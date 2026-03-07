export interface LeaderboardEntry {
    user_id: string;
    username: string;
    total_winners_correct: number;
    winner_accuracy: number;
    team_score_mse: number;
    total_runs_mse: number;
    leaderboard_score: number;
    rank: number;
}
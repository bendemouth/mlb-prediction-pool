export interface LeaderboardEntry {
    user_id: string;
    username: string;
    total_winners_correct: number;
    winner_accuracy: number;
    total_score_error: number;
    total_runs_error: number;
    rank: number;
    update_time: string;
}
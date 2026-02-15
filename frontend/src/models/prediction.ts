export interface Prediction {
    user_id: string;
    game_id: string;
    home_score_predicted: number;
    away_score_predicted: number;
    total_score_predicted: number;
    confidence: number;
    predicted_winner_id: string;
    actual_winner_id: string;
    winner_correct: boolean;
    home_score_error: number;
    away_score_error: number;
    total_score_error: number;
    submitted_at: string;
}
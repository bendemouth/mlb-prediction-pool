import { Game } from './game';

export interface GamePredictionSummary extends Game {
    prediction_count: number;
    avg_home_score_predicted: number;
    avg_away_score_predicted: number;
    avg_total_score_predicted: number;
    avg_confidence: number;
}

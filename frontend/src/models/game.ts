export interface Game {
    game_id: string;
    date: string;
    home_team: string;
    home_team_id: string;
    away_team: string;
    away_team_id: string;
    home_score: number | null;
    away_score: number | null;
    status: string;
    winner: string | null;
}
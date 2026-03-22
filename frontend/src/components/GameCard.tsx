import React from 'react';
import {
  Box,
  Card,
  CardContent,
  Chip,
  Divider,
  Tooltip,
  Typography,
} from '@mui/material';
import SportBaseballIcon from '@mui/icons-material/SportsBaseball';
import PeopleIcon from '@mui/icons-material/People';
import { GamePredictionSummary } from '../models/game_prediction_summary';
import StatBox from './StatBox';

function formatGameDate(dateStr: string): string {
  const date = new Date(dateStr);
  return date.toLocaleDateString('en-US', {
    weekday: 'short',
    month: 'short',
    day: 'numeric',
    hour: 'numeric',
    minute: '2-digit',
    timeZoneName: 'short',
  });
}

export interface GameCardProps {
  game: GamePredictionSummary;
  onViewPredictions: (game: GamePredictionSummary) => void;
}

function GameCard({ game, onViewPredictions }: GameCardProps) {
  const hasPredictions = game.prediction_count > 0;

  return (
    <Card elevation={3} sx={{ borderRadius: 2, overflow: 'visible' }}>
      <CardContent sx={{ pb: '16px !important' }}>
        {/* Teams header */}
        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mb: 1 }}>
          <Box sx={{ flex: 1, textAlign: 'center' }}>
            <Typography variant="h6" fontWeight="bold">
              {game.away_team}
            </Typography>
            <Typography variant="caption" color="text.secondary">Away</Typography>
          </Box>

          <Box sx={{ px: 2, textAlign: 'center' }}>
            <SportBaseballIcon sx={{ color: 'secondary.main', fontSize: 28 }} />
            <Typography variant="caption" display="block" color="text.secondary" sx={{ mt: 0.25 }}>
              @
            </Typography>
          </Box>

          <Box sx={{ flex: 1, textAlign: 'center' }}>
            <Typography variant="h6" fontWeight="bold">
              {game.home_team}
            </Typography>
            <Typography variant="caption" color="text.secondary">Home</Typography>
          </Box>
        </Box>

        {/* Date */}
        <Typography variant="body2" color="text.secondary" textAlign="center" sx={{ mb: 2 }}>
          {formatGameDate(game.date)}
        </Typography>

        <Divider sx={{ mb: 2 }} />

        {/* Community predictions */}
        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mb: 1.5 }}>
          <Typography variant="subtitle2" fontWeight="bold" color="text.primary">
            Community Predictions
          </Typography>
          <Tooltip title={hasPredictions ? 'Click to view all predictions' : 'No predictions yet'}>
            <span>
              <Chip
                icon={<PeopleIcon />}
                label={`${game.prediction_count} prediction${game.prediction_count !== 1 ? 's' : ''}`}
                size="small"
                variant={hasPredictions ? 'filled' : 'outlined'}
                color={hasPredictions ? 'primary' : 'default'}
                clickable={hasPredictions}
                onClick={hasPredictions ? () => onViewPredictions(game) : undefined}
              />
            </span>
          </Tooltip>
        </Box>

        {hasPredictions ? (
          <Box
            sx={{
              display: 'flex',
              justifyContent: 'space-around',
              alignItems: 'flex-start',
              bgcolor: 'background.default',
              borderRadius: 1,
              py: 1.5,
              px: 1,
            }}
          >
            <StatBox label="Avg Away Score" value={game.avg_away_score_predicted} />
            <Divider orientation="vertical" flexItem />
            <StatBox label="Avg Home Score" value={game.avg_home_score_predicted} />
            <Divider orientation="vertical" flexItem />
            <StatBox label="Avg Total Runs" value={game.avg_total_score_predicted} />
            <Divider orientation="vertical" flexItem />
            <Tooltip title="Average confidence across all predictions (0–100)">
              <Box>
                <StatBox label="Avg Confidence" value={game.avg_confidence} unit="%" />
              </Box>
            </Tooltip>
          </Box>
        ) : (
          <Box
            sx={{
              textAlign: 'center',
              py: 2,
              bgcolor: 'background.default',
              borderRadius: 1,
            }}
          >
            <Typography variant="body2" color="text.secondary">
              No predictions yet — be the first!
            </Typography>
          </Box>
        )}
      </CardContent>
    </Card>
  );
}

export default GameCard;

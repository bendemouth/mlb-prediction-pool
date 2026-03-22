import React, { useEffect, useMemo, useState } from 'react';
import {
  Alert,
  Box,
  Chip,
  CircularProgress,
  Container,
  Grid,
  Typography,
} from '@mui/material';
import SportBaseballIcon from '@mui/icons-material/SportsBaseball';
import useAuth from '../hooks/useAuth';
import { GamePredictionSummary } from '../models/game_prediction_summary';
import GameCard from '../components/GameCard';
import GamePredictionsDialog from '../components/GamePredictionsDialog';

function toLocalDateKey(dateStr: string): string {
  const d = new Date(dateStr);
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`;
}

function formatDateChip(dateKey: string): string {
  const [year, month, day] = dateKey.split('-').map(Number);
  const d = new Date(year, month - 1, day);
  return d.toLocaleDateString('en-US', { weekday: 'short', month: 'short', day: 'numeric' });
}

function Predictions() {
  const { getToken } = useAuth();
  const [games, setGames] = useState<GamePredictionSummary[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedDate, setSelectedDate] = useState<string | null>(null);
  const [dialogGame, setDialogGame] = useState<GamePredictionSummary | null>(null);

  useEffect(() => {
    fetchGames();
  }, []);

  const fetchGames = async () => {
    try {
      setLoading(true);
      const token = await getToken();
      if (!token) throw new Error('Not authenticated');

      const response = await fetch('/games/upcoming', {
        headers: { Authorization: `Bearer ${token}` },
      });

      if (!response.ok) {
        throw new Error('Failed to fetch upcoming games');
      }

      const data: GamePredictionSummary[] = await response.json();
      setGames(data ?? []);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An unknown error occurred');
    } finally {
      setLoading(false);
    }
  };

  // Derive sorted unique date keys from the games list
  const uniqueDates = useMemo(() => {
    const keys = Array.from(new Set(games.map((g) => toLocalDateKey(g.date))));
    keys.sort();
    return keys;
  }, [games]);

  const filteredGames = useMemo(() => {
    if (!selectedDate) return games;
    return games.filter((g) => toLocalDateKey(g.date) === selectedDate);
  }, [games, selectedDate]);

  if (loading) {
    return (
      <Container maxWidth="lg" sx={{ mt: 4, textAlign: 'center' }}>
        <CircularProgress />
        <Typography variant="h6" sx={{ mt: 2 }}>
          Loading upcoming games...
        </Typography>
      </Container>
    );
  }

  if (error) {
    return (
      <Container maxWidth="lg" sx={{ mt: 4 }}>
        <Alert severity="error">{error}</Alert>
      </Container>
    );
  }

  return (
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
      <Box sx={{ display: 'flex', alignItems: 'center', mb: 3 }}>
        <SportBaseballIcon sx={{ fontSize: 40, mr: 2, color: 'secondary.main' }} />
        <Typography variant="h3" component="h1">
          Upcoming Games
        </Typography>
      </Box>

      <Typography variant="body1" color="text.secondary" sx={{ mb: 3 }}>
        Community prediction averages for upcoming games. Scores shown are averages across all submitted predictions.
      </Typography>

      {/* Date filter chips */}
      {uniqueDates.length > 1 && (
        <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1, mb: 3 }}>
          <Chip
            label="All Dates"
            onClick={() => setSelectedDate(null)}
            color={selectedDate === null ? 'primary' : 'default'}
            variant={selectedDate === null ? 'filled' : 'outlined'}
          />
          {uniqueDates.map((dateKey) => (
            <Chip
              key={dateKey}
              label={formatDateChip(dateKey)}
              onClick={() => setSelectedDate(dateKey === selectedDate ? null : dateKey)}
              color={selectedDate === dateKey ? 'primary' : 'default'}
              variant={selectedDate === dateKey ? 'filled' : 'outlined'}
            />
          ))}
        </Box>
      )}

      {filteredGames.length === 0 ? (
        <Alert severity="info">No upcoming games found.</Alert>
      ) : (
        <Grid container spacing={3}>
          {filteredGames.map((game) => (
            <Grid size={{ xs: 12, sm: 6, lg: 4 }} key={game.game_id}>
              <GameCard game={game} onViewPredictions={(g) => setDialogGame(g)} />
            </Grid>
          ))}
        </Grid>
      )}

      <GamePredictionsDialog
        game={dialogGame}
        onClose={() => setDialogGame(null)}
        getToken={getToken}
      />
    </Container>
  );
}

export default Predictions;
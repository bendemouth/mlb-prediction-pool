import React, { useEffect, useState } from 'react';
import {
  Alert,
  Box,
  CircularProgress,
  Dialog,
  DialogContent,
  DialogTitle,
  IconButton,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography,
} from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';
import { GamePredictionSummary } from '../models/game_prediction_summary';
import { Prediction } from '../models/prediction';
import User from '../models/user';

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

export interface GamePredictionsDialogProps {
  game: GamePredictionSummary | null;
  onClose: () => void;
  getToken: () => Promise<string | null>;
}

function GamePredictionsDialog({ game, onClose, getToken }: GamePredictionsDialogProps) {
  const [predictions, setPredictions] = useState<Prediction[]>([]);
  const [usersMap, setUsersMap] = useState<Record<string, string>>({});
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!game) return;

    let cancelled = false;

    const load = async () => {
      setLoading(true);
      setError(null);
      try {
        const token = await getToken();
        if (!token) throw new Error('Not authenticated');

        const [predsRes, usersRes] = await Promise.all([
          fetch(`/predictions/game?gameId=${encodeURIComponent(game.game_id)}`, {
            headers: { Authorization: `Bearer ${token}` },
          }),
          fetch('/users/listUsers', {
            headers: { Authorization: `Bearer ${token}` },
          }),
        ]);

        if (!predsRes.ok) throw new Error('Failed to fetch predictions');
        if (!usersRes.ok) throw new Error('Failed to fetch users');

        const [predsData, usersData]: [Prediction[], User[]] = await Promise.all([
          predsRes.json(),
          usersRes.json(),
        ]);

        if (!cancelled) {
          setPredictions(predsData ?? []);
          const map: Record<string, string> = {};
          for (const u of usersData ?? []) {
            map[u.id] = u.username;
          }
          setUsersMap(map);
        }
      } catch (err) {
        if (!cancelled) setError(err instanceof Error ? err.message : 'Failed to load predictions');
      } finally {
        if (!cancelled) setLoading(false);
      }
    };

    load();
    return () => { cancelled = true; };
  }, [game, getToken]);

  const predictedWinnerName = (p: Prediction): string => {
    if (!game) return p.predicted_winner_id;
    if (p.predicted_winner_id === game.home_team_id) return game.home_team;
    if (p.predicted_winner_id === game.away_team_id) return game.away_team;
    return p.predicted_winner_id;
  };

  return (
    <Dialog open={!!game} onClose={onClose} maxWidth="md" fullWidth>
      <DialogTitle sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Box>
          <Typography variant="h6" component="span">
            {game ? `${game.away_team} @ ${game.home_team}` : ''}
          </Typography>
          {game && (
            <Typography variant="body2" color="text.secondary">
              {formatGameDate(game.date)}
            </Typography>
          )}
        </Box>
        <IconButton onClick={onClose} size="small">
          <CloseIcon />
        </IconButton>
      </DialogTitle>

      <DialogContent dividers>
        {loading && (
          <Box sx={{ textAlign: 'center', py: 4 }}>
            <CircularProgress />
            <Typography variant="body2" sx={{ mt: 1 }}>Loading predictions…</Typography>
          </Box>
        )}

        {error && <Alert severity="error">{error}</Alert>}

        {!loading && !error && predictions.length === 0 && (
          <Alert severity="info">No predictions have been submitted for this game yet.</Alert>
        )}

        {!loading && !error && predictions.length > 0 && (
          <TableContainer component={Paper} variant="outlined">
            <Table size="small">
              <TableHead sx={{ bgcolor: 'primary.main' }}>
                <TableRow>
                  <TableCell sx={{ color: 'white', fontWeight: 'bold' }}>User</TableCell>
                  <TableCell align="right" sx={{ color: 'white', fontWeight: 'bold' }}>Away Score</TableCell>
                  <TableCell align="right" sx={{ color: 'white', fontWeight: 'bold' }}>Home Score</TableCell>
                  <TableCell align="right" sx={{ color: 'white', fontWeight: 'bold' }}>Total Runs</TableCell>
                  <TableCell align="right" sx={{ color: 'white', fontWeight: 'bold' }}>Confidence</TableCell>
                  <TableCell sx={{ color: 'white', fontWeight: 'bold' }}>Predicted Winner</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {predictions.map((p, i) => (
                  <TableRow key={i} sx={{ '&:nth-of-type(odd)': { bgcolor: 'action.hover' } }}>
                    <TableCell>{usersMap[p.user_id] ?? p.user_id}</TableCell>
                    <TableCell align="right">{p.away_score_predicted.toFixed(1)}</TableCell>
                    <TableCell align="right">{p.home_score_predicted.toFixed(1)}</TableCell>
                    <TableCell align="right">{p.total_score_predicted.toFixed(1)}</TableCell>
                    <TableCell align="right">{p.confidence.toFixed(1)}%</TableCell>
                    <TableCell>{predictedWinnerName(p)}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        )}
      </DialogContent>
    </Dialog>
  );
}

export default GamePredictionsDialog;

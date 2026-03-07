import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { LeaderboardEntry } from "../models/leaderboard_entry";
import EmojeEventsIcon from '@mui/icons-material/EmojiEvents';
import {
  Box,
  Container,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography,
  CircularProgress,
  Alert,
  Chip,
} from '@mui/material';
import { toProfilePath } from '../utils/profileRoute';

function Leaderboard() {
    const [leaderboard, setLeaderboard] = useState<LeaderboardEntry[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const navigate = useNavigate();

    useEffect(() => {
        fetchLeaderboard();
    }, []);

    const fetchLeaderboard = async () => {
        try {
            setLoading(true);
            const response = await fetch("/leaderboard");
            if (!response.ok) {
                throw new Error("Failed to fetch leaderboard");
            }

            const data = await response.json();
            setLeaderboard(data || []);
        } catch (error) {
            setError(error instanceof Error ? error.message : "An unknown error occurred");
        } finally {
            setLoading(false);
        }
    };

    const handleUserClick = (username: string) => {
      navigate(toProfilePath(username));
    }

    const getRankColor = (rank: number) => {
        if (rank === 1) return "gold";
        if (rank === 2) return "silver";
        if (rank === 3) return "#ce8946"; // bronze
        return "inherit";
    };

    if (loading) {
        return (
            <Container maxWidth="lg"
                sx={{ mt: 4, textAlign: "center" }}>
                    <CircularProgress />
                    <Typography variant="h6" sx={{ mt: 2 }}>
                        Loading leaderboard...
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
        <EmojeEventsIcon sx={{ fontSize: 40, mr: 2, color: 'gold' }} />
        <Typography variant="h3" component="h1">
          Leaderboard
        </Typography>
      </Box>

      {leaderboard.length === 0 ? (
        <Alert severity="info">No leaderboard data available yet.</Alert>
      ) : (
        <TableContainer component={Paper} elevation={3}>
          <Table>
            <TableHead sx={{ bgcolor: 'primary.main' }}>
              <TableRow>
                <TableCell sx={{ color: 'white', fontWeight: 'bold' }}>Rank</TableCell>
                <TableCell sx={{ color: 'white', fontWeight: 'bold' }}>User</TableCell>
                <TableCell align="right" sx={{ color: 'white', fontWeight: 'bold' }}>
                  Winners Correct
                </TableCell>
                <TableCell align="right" sx={{ color: 'white', fontWeight: 'bold' }}>
                  Winner Accuracy
                </TableCell>
                <TableCell align="right" sx={{ color: 'white', fontWeight: 'bold' }}>
                  Team Score RMSE
                </TableCell>
                <TableCell align="right" sx={{ color: 'white', fontWeight: 'bold' }}>
                  Total Runs RMSE
                </TableCell>
                <TableCell align="right" sx={{ color: 'white', fontWeight: 'bold' }}>
                  Leaderboard Score
                </TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {leaderboard.map((entry: LeaderboardEntry) => (
                <TableRow
                  key={entry.user_id}
                  onClick={() => handleUserClick(entry.username)}
                  sx={{
                    cursor: 'pointer',
                    '&:hover': { bgcolor: 'action.hover' },
                    bgcolor: entry.rank <= 3 ? 'action.selected' : 'inherit',
                  }}
                >
                  <TableCell>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <Typography
                        sx={{
                          fontWeight: 'bold',
                          color: getRankColor(entry.rank),
                          fontSize: entry.rank <= 3 ? '1.2rem' : '1rem',
                        }}
                      >
                        #{entry.rank}
                      </Typography>
                      {entry.rank <= 3 && <EmojeEventsIcon sx={{ color: getRankColor(entry.rank) }} />}
                    </Box>
                  </TableCell>
                  <TableCell>
                    <Typography sx={{ fontWeight: entry.rank <= 3 ? 'bold' : 'normal' }}>
                      {entry.username}
                    </Typography>
                  </TableCell>
                  <TableCell align="right">
                    <Chip label={entry.total_winners_correct} color="primary" size="small" />
                  </TableCell>
                  <TableCell align="right">
                    <Typography sx={{ fontWeight: 'medium' }}>
                      {(entry.winner_accuracy * 100).toFixed(1)}%
                    </Typography>
                  </TableCell>
                  <TableCell align="right">{entry.team_score_mse.toFixed(2)}</TableCell>
                  <TableCell align="right">{entry.total_runs_mse.toFixed(2)}</TableCell>
                  <TableCell align="right">{entry.leaderboard_score.toFixed(2)}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      )}
    </Container>
  );
}

export default Leaderboard;
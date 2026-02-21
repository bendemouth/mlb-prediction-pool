import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Box,
  Button,
  Card,
  CardContent,
  Container,
  Typography,
  CircularProgress,
  Alert,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Chip,
} from '@mui/material';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import CancelIcon from '@mui/icons-material/Cancel';
import UserStats from '../models/user_stats';
import { Prediction } from '../models/prediction';

function UserProfile() {
  const { userId } = useParams<{ userId?: string }>();
  const navigate = useNavigate();
  const [stats, setStats] = useState<UserStats | null>(null);
  const [predictions, setPredictions] = useState<Prediction[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (userId) {
      fetchUserData(userId);
    } else {
      setError('No user ID provided');
      setLoading(false);
    }
  }, [userId]);

  const fetchUserData = async (id: string) => {
    try {
      setLoading(true);
      
      // Fetch user stats
      const statsResponse = await fetch(`/users/stats?user_id=${id}`);
      if (!statsResponse.ok) throw new Error('Failed to fetch user stats');
      const statsData = await statsResponse.json();
      setStats(statsData);

      // Fetch user predictions
      const predsResponse = await fetch(`/predictions?userId=${id}`);
      if (!predsResponse.ok) throw new Error('Failed to fetch predictions');
      const predsData = await predsResponse.json();
      setPredictions(predsData || []);
    } catch (error) {
      console.error('Error fetching user data:', error);
      setError('Failed to load user data');
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <Container maxWidth="lg" sx={{ mt: 4, textAlign: 'center' }}>
        <CircularProgress />
      </Container>
    );
  }

  if (error || !stats) {
    return (
      <Container maxWidth="lg" sx={{ mt: 4 }}>
        <Alert severity="error">{error || 'User not found'}</Alert>
        <Button startIcon={<ArrowBackIcon />} onClick={() => navigate('/')} sx={{ mt: 2 }}>
          Back to Home
        </Button>
      </Container>
    );
  }

  return (
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
      <Button startIcon={<ArrowBackIcon />} onClick={() => navigate('/leaderboard')} sx={{ mb: 3 }}>
        Back to Leaderboard
      </Button>

      <Typography variant="h3" component="h1" gutterBottom>
        {stats.username}
      </Typography>

      <Box
        sx={{
          display: 'flex',
          flexWrap: 'wrap',
          gap: 3,
          mb: 4,
        }}
      >
        <Box sx={{ flex: '1 1 200px', minWidth: 0 }}>
          <Card elevation={3}>
            <CardContent>
              <Typography color="text.secondary" gutterBottom>
                Rank
              </Typography>
              <Typography variant="h4" component="div">
                #{stats.rank}
              </Typography>
            </CardContent>
          </Card>
        </Box>

        <Box sx={{ flex: '1 1 200px', minWidth: 0 }}>
          <Card elevation={3}>
            <CardContent>
              <Typography color="text.secondary" gutterBottom>
                Winner Accuracy
              </Typography>
              <Typography variant="h4" component="div">
                {(stats.winner_accuracy * 100).toFixed(1)}%
              </Typography>
            </CardContent>
          </Card>
        </Box>

        <Box sx={{ flex: '1 1 200px', minWidth: 0 }}>
          <Card elevation={3}>
            <CardContent>
              <Typography color="text.secondary" gutterBottom>
                Winners Correct
              </Typography>
              <Typography variant="h4" component="div">
                {stats.total_winners_correct}
              </Typography>
            </CardContent>
          </Card>
        </Box>

        <Box sx={{ flex: '1 1 200px', minWidth: 0 }}>
          <Card elevation={3}>
            <CardContent>
              <Typography color="text.secondary" gutterBottom>
                Total Predictions
              </Typography>
              <Typography variant="h4" component="div">
                {predictions.length}
              </Typography>
            </CardContent>
          </Card>
        </Box>
      </Box>

      <Typography variant="h5" component="h2" gutterBottom sx={{ mt: 4 }}>
        Prediction History
      </Typography>

      {predictions.length === 0 ? (
        <Alert severity="info">No predictions made yet.</Alert>
      ) : (
        <TableContainer component={Paper} elevation={3}>
          <Table>
            <TableHead sx={{ bgcolor: 'primary.main' }}>
              <TableRow>
                <TableCell sx={{ color: 'white' }}>Game ID</TableCell>
                <TableCell align="center" sx={{ color: 'white' }}>Result</TableCell>
                <TableCell align="right" sx={{ color: 'white' }}>Home Score</TableCell>
                <TableCell align="right" sx={{ color: 'white' }}>Away Score</TableCell>
                <TableCell align="right" sx={{ color: 'white' }}>Confidence</TableCell>
                <TableCell align="right" sx={{ color: 'white' }}>Score Error</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {predictions.map((pred) => (
                <TableRow key={pred.game_id}>
                  <TableCell>{pred.game_id}</TableCell>
                  <TableCell align="center">
                    {pred.winner_correct === undefined ? (
                      <Chip label="Pending" size="small" />
                    ) : pred.winner_correct ? (
                      <CheckCircleIcon color="success" />
                    ) : (
                      <CancelIcon color="error" />
                    )}
                  </TableCell>
                  <TableCell align="right">{pred.home_score_predicted.toFixed(1)}</TableCell>
                  <TableCell align="right">{pred.away_score_predicted.toFixed(1)}</TableCell>
                  <TableCell align="right">{(pred.confidence * 100).toFixed(0)}%</TableCell>
                  <TableCell align="right">
                    {pred.total_score_error !== undefined ? pred.total_score_error.toFixed(2) : '-'}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      )}
    </Container>
  );
}

export default UserProfile;
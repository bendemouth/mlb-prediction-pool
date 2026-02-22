import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Box,
  Card,
  CardContent,
  Container,
  Typography,
  CircularProgress,
  Alert,
  Chip,
} from '@mui/material';
import SportsBaseballIcon from '@mui/icons-material/SportsBaseball';
import LeaderboardIcon from '@mui/icons-material/Leaderboard';
import PredictionsIcon from '@mui/icons-material/Psychology';
import User from '../models/user';
import { Icon, Trophy} from 'lucide-react';
import { baseball } from '@lucide/lab';

function Home() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    fetchUsers();
  }, []);

  const fetchUsers = async () => {
    try {
      const response = await fetch('/users/listUsers');
      if (!response.ok) throw new Error('Failed to fetch users');
      const data = await response.json();
      setUsers(data || []);
    } catch (error) {
      console.error('Error fetching users:', error);
    } finally {
      setLoading(false);
    }
  };


  return (
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
      <Box sx={{ textAlign: 'center', mb: 6 }}>
        {/*<SportsBaseballIcon sx={{ fontSize: 80, color: 'primary.main', mb: 2 }} /> / */}
        <Icon iconNode={baseball} size={100} />
        <Typography variant="h2" component="h1" gutterBottom>
          Welcome to (ML)B Predictions
        </Typography>
        <Typography variant="h5" color="text.secondary" sx={{ mb: 4 }}>
          Compete with machine learning models to predict MLB game outcomes
        </Typography>
      </Box>

      <Box
        sx={{
          display: 'flex',
          flexWrap: 'wrap',
          gap: 4,
          mb: 6,
        }}
      >
        <Box sx={{ flex: '1 1 300px', minWidth: 0 }}>
          <Card
            elevation={3}
            sx={{
              height: '100%',
              cursor: 'pointer',
              transition: 'transform 0.2s',
              '&:hover': { transform: 'translateY(-8px)' },
            }}
            onClick={() => navigate('/leaderboard')}
          >
            <CardContent sx={{ textAlign: 'center', p: 4 }}>
              {/*<LeaderboardIcon sx={{ fontSize: 60, color: 'primary.main', mb: 2 }} />*/}
              <Trophy size={50} />
              <Typography variant="h5" component="h2" gutterBottom>
                Leaderboard
              </Typography>
              <Typography color="text.secondary">
                See how your ML model stacks up against the competition
              </Typography>
            </CardContent>
          </Card>
        </Box>

        <Box sx={{ flex: '1 1 300px', minWidth: 0 }}>
          <Card
            elevation={3}
            sx={{
              height: '100%',
              cursor: 'pointer',
              transition: 'transform 0.2s',
              '&:hover': { transform: 'translateY(-8px)' },
            }}
            onClick={() => navigate('/predictions')}
          >
            <CardContent sx={{ textAlign: 'center', p: 4 }}>
              <PredictionsIcon sx={{ fontSize: 60, color: 'primary.main', mb: 2 }} />
              <Typography variant="h5" component="h2" gutterBottom>
                Make Predictions
              </Typography>
              <Typography color="text.secondary">
                Submit your model's predictions for upcoming games
              </Typography>
            </CardContent>
          </Card>
        </Box>

        <Box sx={{ flex: '1 1 300px', minWidth: 0 }}>
          <Card
            elevation={3}
            sx={{
              height: '100%',
              cursor: 'pointer',
              transition: 'transform 0.2s',
              '&:hover': { transform: 'translateY(-8px)' },
            }}
            onClick={() => navigate('/leaderboard')}
          >
            <CardContent sx={{ textAlign: 'center', p: 4 }}>
              {/*<SportsBaseballIcon sx={{ fontSize: 60, color: 'primary.main', mb: 2 }} /> */}
               <Icon iconNode={baseball} size={50} style={{ marginBottom: 16 }} />
              <Typography variant="h5" component="h2" gutterBottom>
                View Stats
              </Typography>
              <Typography color="text.secondary">
                View prediction history and performance metrics for all users
              </Typography>
            </CardContent>
          </Card>
        </Box>
      </Box>

      <Card elevation={3} sx={{ p: 3 }}>
        <Typography variant="h5" component="h2" gutterBottom>
          Active Participants
        </Typography>
        {loading ? (
          <Box sx={{ textAlign: 'center', py: 3 }}>
            <CircularProgress />
          </Box>
        ) : users.length === 0 ? (
          <Alert severity="info">No users registered yet.</Alert>
        ) : (
          <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1, mt: 2 }}>
            {users.map((user) => (
              <Chip
                key={user.id}
                label={user.username}
                onClick={() => navigate(`/profile/${user.id}`)}
                color="primary"
                variant="outlined"
                sx={{ cursor: 'pointer' }}
              />
            ))}
          </Box>
        )}
      </Card>
    </Container>
  );
}

export default Home;
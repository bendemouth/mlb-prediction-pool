import { useState, useEffect } from 'react';
import {  Routes, Route } from 'react-router-dom';
import Home from './pages/Home';
import Predictions from './pages/Predictions';
import UserProfile from './pages/UserProfile';
import Leaderboard from './pages/Leaderboard';
import Box from '@mui/material/Box';
import Container from '@mui/material/Container';
import AppNavbar from './components/Navbar';
import { fetchAuthSession } from 'aws-amplify/auth';
import ProtectedRoute from './components/ProtectedRoute';
import Login from './pages/Login';
import SetupProfile from './pages/SetupProfile';
import useAuth from './hooks/useAuth';
import User from './models/user';
import { toProfilePath } from './utils/profileRoute';

// Define types that match your Go backend structs
interface HealthStatus {
  service: string;
  database: string;
}

function App() {
  const [healthStatus, setHealthStatus] = useState<HealthStatus | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [profilePath, setProfilePath] = useState<string>('/setup-profile');
  const { isAuthenticated, signOut, user, getToken } = useAuth();

  // Check backend health on mount
  useEffect(() => {
    checkBackendHealth();
    fetchAuthSession()
      .then(session => console.log('Auth session:', session))
      .catch(err => console.error('Auth session error:', err));
  }, []);

  useEffect(() => {
    let cancelled = false;

    const resolveProfilePath = async () => {
      if (!isAuthenticated || !user?.userId) {
        setProfilePath('/login');
        return;
      }

      try {
        const token = await getToken();
        if (!token) {
          setProfilePath('/login');
          return;
        }

        const response = await fetch(`/users?user_id=${encodeURIComponent(user.userId)}`, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (!response.ok) {
          setProfilePath('/setup-profile');
          return;
        }

        const profile: User = await response.json();
        if (!cancelled) {
          setProfilePath(toProfilePath(profile.username));
        }
      } catch {
        if (!cancelled) {
          setProfilePath('/setup-profile');
        }
      }
    };

    resolveProfilePath();

    return () => {
      cancelled = true;
    };
  }, [isAuthenticated, user?.userId]);

  const checkBackendHealth = async (): Promise<void> => {
    try {
      setLoading(true);
      setError(null);
      
      const response = await fetch('/health');
      
      // Check if response is OK
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      
      const data: HealthStatus = await response.json();
      console.log('Health check response:', data);
      
      setHealthStatus(data);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Unknown error';
      setError(`Failed to connect to backend: ${errorMessage}`);
      console.error('Backend health check failed:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleSignOut = async () => {
    await signOut();
    window.location.assign('/login');
  };

  const navItems = [
    { label: 'Home', to: '/' },
    { label: 'Profile', to: profilePath },
    { label: 'Leaderboard', to: '/leaderboard' },
    { label: 'Predictions', to: '/predictions' },
  ];

  return (
    <Box sx={{ minHeight: "100vh" }}>
      <AppNavbar
        title={
          <Box component="span" sx={{ display: "inline-flex", alignItems: "baseline" }}>
            (
            <Box component="span" sx={{ color: "#e9422f" }}>
              ML
            </Box>
            )B Predictions
          </Box>
        }
        navItems={navItems}
        isAuthenticated={isAuthenticated}
        onSignOut={handleSignOut}
      />
      <Container maxWidth="lg" sx={{ py: 3 }}>
        <Routes>
          {/* Public */}
          <Route path="/login" element={<Login />} />

          {/* Needs Cognito auth but no profile yet (setup step) */}
          <Route path="/setup-profile" element={
            <ProtectedRoute allowWithoutProfile>
              <SetupProfile />
            </ProtectedRoute>
          } />

          {/* Fully protected (auth + profile) */}
          <Route path="/" element={<ProtectedRoute><Home /></ProtectedRoute>} />
          <Route path="/predictions" element={<ProtectedRoute><Predictions /></ProtectedRoute>} />
          <Route path="/leaderboard" element={<ProtectedRoute><Leaderboard /></ProtectedRoute>} />
          <Route path="/profile/:username" element={<ProtectedRoute><UserProfile /></ProtectedRoute>} />
        </Routes>
      </Container>
    </Box>
  );
}

export default App;
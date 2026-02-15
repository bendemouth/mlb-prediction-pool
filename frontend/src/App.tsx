import React, { useState, useEffect } from 'react';
import './App.css';
import { BrowserRouter, Link, Routes, Route } from 'react-router-dom';
import Leaderboad from './pages/Leaderboard';
import Home from './pages/Home';
import Predictions from './pages/Predictions';
import UserProfile from './pages/UserProfile';
import Leaderboard from './pages/Leaderboard';
import Navbar from './components/Navbar';


// Define types that match your Go backend structs
interface HealthStatus {
  service: string;
  database: string;
}

function App() {
  const [healthStatus, setHealthStatus] = useState<HealthStatus | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  // Check backend health on mount
  useEffect(() => {
    checkBackendHealth();
  }, []);

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

  return (
    <BrowserRouter>
      <div className="App">
        <Navbar />

        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/predictions" element={<Predictions />} />
          <Route path="/leaderboard" element={<Leaderboard />} />
          <Route path="/profile" element={<UserProfile />} />
          <Route path="/profile/:userId" element={<UserProfile />} />
        </Routes>
      </div>
    </BrowserRouter>
  );
}

export default App;
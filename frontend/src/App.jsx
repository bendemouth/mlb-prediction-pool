import React, { useState, useEffect } from 'react';
import './App.css';

function App() {
  const [healthStatus, setHealthStatus] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Check backend health on mount
  useEffect(() => {
    checkBackendHealth();
  }, []);

  const checkBackendHealth = async () => {
    try {
      setLoading(true);
      setError(null);
      
      const response = await fetch('http://localhost:8080/health');
      
      // Check if response is OK
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      
      const data = await response.json();
      console.log('Health check response:', data); // Debug log
      
      setHealthStatus(data);
    } catch (err) {
      setError(`Failed to connect to backend: ${err.message}`);
      console.error('Backend health check failed:', err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="App">
      <header className="App-header">
        <h1>MLB Prediction Pool</h1>
        
        {/* Backend Status Indicator */}
        <div className="status-indicator">
          {loading && (
            <div className="status-loading">
              <p>🔄 Checking backend connection...</p>
            </div>
          )}
          
          {!loading && error && (
            <div className="status-error">
              <p>❌ {error}</p>
              <button onClick={checkBackendHealth}>Retry Connection</button>
            </div>
          )}
          
          {!loading && !error && healthStatus && (
            <div className="status-success">
              <h2>✅ Backend Connected</h2>
              <div className="status-details">
                <div className="status-item">
                  <span className="status-label">Service:</span>
                  <span className={`status-value ${healthStatus.service === 'healthy' ? 'healthy' : 'unhealthy'}`}>
                    {healthStatus.service || 'unknown'}
                  </span>
                </div>
                <div className="status-item">
                  <span className="status-label">Database:</span>
                  <span className={`status-value ${healthStatus.database === 'healthy' ? 'healthy' : 'unhealthy'}`}>
                    {healthStatus.database || 'unknown'}
                  </span>
                </div>
              </div>
              <button onClick={checkBackendHealth} className="refresh-button">
                Refresh Status
              </button>
            </div>
          )}
        </div>

        <nav className="main-nav">
          <p>Coming soon: Leaderboard, Predictions, and more!</p>
        </nav>
      </header>
    </div>
  );
}

export default App;
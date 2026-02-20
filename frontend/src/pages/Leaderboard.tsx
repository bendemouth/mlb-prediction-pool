import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { LeaderboardEntry } from '../models/leaderboard_entry';
import './leaderboard.css';

function Leaderboard() {
    const [leaderboard, setLeaderboard] = useState([]);
    const [loading, setLoading] = useState(true);
    const navigate = useNavigate();

    useEffect(() => {
        fetchLeaderboard();
    }, []);

    const fetchLeaderboard = async () => {
        try {
            const response = await fetch('/api/leaderboard');
            const data = await response.json();
            setLeaderboard(data);
        } catch (error) {
            console.error('Error fetching leaderboard:', error);
        } finally {
            setLoading(false);
        }
    };

    const handleUserClick = (userId: string) => {
        navigate(`/profile/${userId}`);
    }

    if (loading) {
        return <div className="leaderboard-container"><div className="loading">Loading...</div></div>;
    }

    return (
        <div className="leaderboard-container">
            <h2>Leaderboard</h2>
            <table>
                <thead>
                    <tr>
                        <th>Rank</th>
                        <th>User</th>
                        <th>Total Runs Error</th>
                        <th>Total Score Error</th>
                        <th>Winner Accuracy</th>
                    </tr>
                </thead>
                <tbody>
                    {leaderboard.map((entry: LeaderboardEntry) => (
                        <tr key={entry.user_id} onClick={() => handleUserClick(entry.user_id)}>
                            <td>{entry.rank}</td>
                            <td>{entry.username}</td>
                            <td>{entry.total_runs_error.toFixed(2)}</td>
                            <td>{entry.total_score_error.toFixed(2)}</td>
                            <td>{(entry.winner_accuracy * 100).toFixed(2)}%</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
}

export default Leaderboard;
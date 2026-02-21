import * as process from 'process';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

export const api = {
  getLeaderboard: async () => {
    const response = await fetch(`${API_BASE_URL}/leaderboard`);
    if (!response.ok) {
      throw new Error('Failed to fetch leaderboard');
    }
    return response.json();
  },

  getUserPredictions: async (userId: string) => {
    const response = await fetch(`${API_BASE_URL}/predictions?userId=${userId}`);
    if (!response.ok) {
      throw new Error('Failed to fetch user predictions');
    }
    return response.json();
  },

  getUsers: async () => {
    const response = await fetch(`${API_BASE_URL}/users/listUsers`);
    if (!response.ok) {
      throw new Error('Failed to fetch users');
    }
    return response.json();
  }
}
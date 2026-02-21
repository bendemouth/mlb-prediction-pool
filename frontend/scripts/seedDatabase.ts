import { DynamoDBClient } from '@aws-sdk/client-dynamodb';
import { DynamoDBDocumentClient, PutCommand, BatchWriteCommand } from '@aws-sdk/lib-dynamodb';
import * as process from 'process';

const client = new DynamoDBClient({
  region: process.env.DYNAMODB_REGION || 'us-east-1',
  endpoint: process.env.DYNAMODB_ENDPOINT || 'http://localhost:8000',
});

const docClient = DynamoDBDocumentClient.from(client);

const USERS_TABLE = process.env.DYNAMODB_USERS_TABLE || 'mlb-prediction-pool-dev-users';
const GAMES_TABLE = process.env.DYNAMODB_GAMES_TABLE || 'mlb-prediction-pool-dev-games';
const PREDICTIONS_TABLE = process.env.DYNAMODB_PREDICTIONS_TABLE || 'mlb-prediction-pool-dev-predictions';

async function seedUsers() {
  const users = [
    { userId: 'user-1', username: 'BaseballBot2000', email: 'bot@example.com', createdAt: new Date().toISOString() },
    { userId: 'user-2', username: 'DiamondPredictor', email: 'diamond@example.com', createdAt: new Date().toISOString() },
    { userId: 'user-3', username: 'StatsGuru', email: 'stats@example.com', createdAt: new Date().toISOString() },
    { userId: 'user-4', username: 'MLBOracle', email: 'oracle@example.com', createdAt: new Date().toISOString() },
    { userId: 'user-5', username: 'PitchPerfect', email: 'pitch@example.com', createdAt: new Date().toISOString() },
  ];

  console.log(`Seeding ${users.length} users into ${USERS_TABLE}...`);
  
  for (const user of users) {
    await docClient.send(new PutCommand({
      TableName: USERS_TABLE,
      Item: user,
    }));
  }

  console.log(`Seeded ${users.length} users`);
}

async function seedGames() {
  const teams = [
    { id: '147', name: 'New York Yankees' },
    { id: '111', name: 'Boston Red Sox' },
    { id: '137', name: 'San Francisco Giants' },
    { id: '119', name: 'Los Angeles Dodgers' },
    { id: '145', name: 'Chicago White Sox' },
    { id: '112', name: 'Chicago Cubs' },
  ];

  const games = [];
  const today = new Date();

  console.log(`Generating games for ${GAMES_TABLE}...`);

  for (let i = 0; i < 20; i++) {
    const homeTeam = teams[Math.floor(Math.random() * teams.length)];
    let awayTeam = teams[Math.floor(Math.random() * teams.length)];
    while (awayTeam.id === homeTeam.id) {
      awayTeam = teams[Math.floor(Math.random() * teams.length)];
    }

    const gameDate = new Date(today);
    gameDate.setDate(today.getDate() + Math.floor(i / 5) - 2);

    const isPast = gameDate < today;
    const homeScore = isPast ? Math.floor(Math.random() * 10) : 0;
    const awayScore = isPast ? Math.floor(Math.random() * 10) : 0;
    const winner = isPast ? (homeScore > awayScore ? homeTeam.id : awayTeam.id) : '';

    games.push({
      gameId: `game-${String(i + 1).padStart(3, '0')}`,
      date: gameDate.toISOString(),
      homeTeam: homeTeam.name,
      homeTeamId: homeTeam.id,
      awayTeam: awayTeam.name,
      awayTeamId: awayTeam.id,
      homeScore,
      awayScore,
      status: isPast ? 'completed' : 'upcoming',
      ...(winner && { winner }),
    });
  }

  for (const game of games) {
    await docClient.send(new PutCommand({
      TableName: GAMES_TABLE,
      Item: game,
    }));
  }

  console.log(`Seeded ${games.length} games`);
}

async function seedPredictions() {
  const userIds = ['user-1', 'user-2', 'user-3', 'user-4', 'user-5'];
  const predictions = [];

  console.log(`Generating predictions for ${PREDICTIONS_TABLE}...`);

  for (let gameNum = 1; gameNum <= 15; gameNum++) {
    for (const userId of userIds) {
      const homeScore = 2 + Math.random() * 6;
      const awayScore = 2 + Math.random() * 6;
      const predictedWinner = homeScore > awayScore ? '147' : '111';

      predictions.push({
        userId,
        gameId: `game-${String(gameNum).padStart(3, '0')}`,
        homeScorePredicted: parseFloat(homeScore.toFixed(1)),
        awayScorePredicted: parseFloat(awayScore.toFixed(1)),
        totalScorePredicted: parseFloat((homeScore + awayScore).toFixed(1)),
        confidence: parseFloat((0.5 + Math.random() * 0.5).toFixed(2)),
        predictedWinnerId: predictedWinner,
        submittedAt: new Date().toISOString(),
      });
    }
  }

  // Batch write predictions (DynamoDB limit is 25 items per batch)
  for (let i = 0; i < predictions.length; i += 25) {
    const batch = predictions.slice(i, i + 25);
    await docClient.send(new BatchWriteCommand({
      RequestItems: {
        [PREDICTIONS_TABLE]: batch.map(item => ({
          PutRequest: { Item: item }
        }))
      }
    }));
  }

  console.log(`Seeded ${predictions.length} predictions`);
}

async function seed() {
  console.log('Seeding database...\n');
  console.log(`Endpoint: ${process.env.DYNAMODB_ENDPOINT}`);
  console.log(`Region: ${process.env.DYNAMODB_REGION}\n`);
  
  try {
    await seedUsers();
    await seedGames();
    await seedPredictions();
    
    console.log('\nDatabase seeded successfully!');
    process.exit(0);
  } catch (error) {
    console.error('Error seeding database:', error);
    process.exit(1);
  }
}

seed();
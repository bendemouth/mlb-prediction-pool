import { DynamoDBClient } from '@aws-sdk/client-dynamodb';
import {
  DynamoDBDocumentClient,
  BatchWriteCommand,
  BatchWriteCommandInput,
} from '@aws-sdk/lib-dynamodb';
import * as process from 'process';

// ---------- Config ----------
const REGION = process.env.DYNAMODB_REGION || 'us-east-1';
const ENDPOINT = process.env.DYNAMODB_ENDPOINT || 'http://localhost:8000';

// Match backend defaults unless explicitly overridden
const USERS_TABLE = process.env.DYNAMODB_USERS_TABLE || 'mlb-prediction-pool-users';
const GAMES_TABLE = process.env.DYNAMODB_GAMES_TABLE || 'mlb-prediction-pool-games';
const PREDICTIONS_TABLE =
  process.env.DYNAMODB_PREDICTIONS_TABLE || 'mlb-prediction-pool-predictions';

// A datasetId makes reruns predictable and prevents mixing “old seed” with “new seed”.
const DATASET_ID = process.env.SEED_DATASET_ID || 'dev-dataset';

// Controls deterministic random generation
const RANDOM_SEED = parseInt(process.env.SEED_RANDOM_SEED || '42', 10);

// How much data to generate
const NUM_GAMES = parseInt(process.env.SEED_NUM_GAMES || '75', 10);
const NUM_PREDICTION_GAMES = parseInt(process.env.SEED_NUM_PREDICTION_GAMES || '40', 10);

// ---------- DynamoDB ----------
const client = new DynamoDBClient({ region: REGION, endpoint: ENDPOINT });
const docClient = DynamoDBDocumentClient.from(client);

// ---------- Deterministic RNG ----------
function mulberry32(seed: number) {
  let a = seed >>> 0;
  return () => {
    a |= 0;
    a = (a + 0x6D2B79F5) | 0;
    let t = Math.imul(a ^ (a >>> 15), 1 | a);
    t = (t + Math.imul(t ^ (t >>> 7), 61 | t)) ^ t;
    return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
  };
}

const rand = mulberry32(RANDOM_SEED);

function randInt(min: number, max: number): number {
  // inclusive
  return Math.floor(rand() * (max - min + 1)) + min;
}

function randFloat(min: number, max: number): number {
  return rand() * (max - min) + min;
}

function clamp(n: number, min: number, max: number): number {
  return Math.max(min, Math.min(max, n));
}

// Approx normal-ish noise using sum of uniforms
function noise(mean = 0, stdev = 1): number {
  const u =
    (rand() + rand() + rand() + rand() + rand() + rand() - 3) * (stdev / 1);
  return mean + u;
}

// ---------- Batch write with retry ----------
async function batchWriteAll(
  tableName: string,
  items: any[],
  batchSize = 25,
  maxRetries = 8
) {
  console.log(`Writing ${items.length} items to ${tableName}...`);

  for (let i = 0; i < items.length; i += batchSize) {
    const batch = items.slice(i, i + batchSize);

    let request: BatchWriteCommandInput = {
      RequestItems: {
        [tableName]: batch.map((Item) => ({ PutRequest: { Item } })),
      },
    };

    for (let attempt = 0; attempt <= maxRetries; attempt++) {
      const res = await docClient.send(new BatchWriteCommand(request));
      const unprocessed = res.UnprocessedItems?.[tableName] || [];

      if (unprocessed.length === 0) break;

      if (attempt === maxRetries) {
        throw new Error(
          `BatchWrite max retries exceeded. Unprocessed items: ${unprocessed.length}`
        );
      }

      const backoffMs = 50 * Math.pow(2, attempt) + randInt(0, 100);
      await new Promise((r) => setTimeout(r, backoffMs));

      request = {
        RequestItems: {
          [tableName]: unprocessed,
        },
      };
    }
  }
}

// ---------- Domain generation ----------
type Team = { id: string; name: string };

const teams: Team[] = [
  { id: '147', name: 'New York Yankees' },
  { id: '111', name: 'Boston Red Sox' },
  { id: '137', name: 'San Francisco Giants' },
  { id: '119', name: 'Los Angeles Dodgers' },
  { id: '145', name: 'Chicago White Sox' },
  { id: '112', name: 'Chicago Cubs' },
];

type SeedUser = { userId: string; username: string; email: string; createdAt: string };

type SeedGame = {
  gameId: string;
  date: string;
  homeTeam: string;
  homeTeamId: string;
  awayTeam: string;
  awayTeamId: string;
  homeScore: number;
  awayScore: number;
  status: 'completed' | 'upcoming';
  winner?: string;
};

type Persona = {
  name: string;
  // Probability of picking correct winner for completed games
  winnerSkill: number; // 0..1
  // Score noise in runs
  scoreStdev: number;
  // confidence range
  confMin: number;
  confMax: number;
  // If true, tends to be overconfident
  overconfident?: boolean;
};

const personas: Persona[] = [
  { name: 'SharpModel', winnerSkill: 0.72, scoreStdev: 1.2, confMin: 0.65, confMax: 0.95 },
  { name: 'SolidModel', winnerSkill: 0.62, scoreStdev: 1.8, confMin: 0.55, confMax: 0.85 },
  { name: 'CoinFlip', winnerSkill: 0.50, scoreStdev: 2.5, confMin: 0.50, confMax: 0.70 },
  { name: 'Overconfident', winnerSkill: 0.53, scoreStdev: 2.2, confMin: 0.80, confMax: 0.99, overconfident: true },
  { name: 'WildCard', winnerSkill: 0.58, scoreStdev: 3.0, confMin: 0.40, confMax: 0.90 },
];

function makeUsers(): SeedUser[] {
  const createdAt = new Date().toISOString();

  return personas.map((p, idx) => ({
    userId: `${DATASET_ID}-user-${idx + 1}`,
    username: p.name,
    email: `${p.name.toLowerCase()}@example.com`,
    createdAt,
  }));
}

function makeGames(): SeedGame[] {
  const games: SeedGame[] = [];
  const today = new Date();

  // Spread across -5..+4 days (some completed, some upcoming)
  for (let i = 0; i < NUM_GAMES; i++) {
    const homeTeam = teams[randInt(0, teams.length - 1)];
    let awayTeam = teams[randInt(0, teams.length - 1)];
    while (awayTeam.id === homeTeam.id) awayTeam = teams[randInt(0, teams.length - 1)];

    const gameDate = new Date(today);
    gameDate.setDate(today.getDate() + Math.floor(i / 6) - 5);

    const isPast = gameDate.getTime() < today.getTime();

    let homeScore = 0;
    let awayScore = 0;
    let winner: string | undefined;

    if (isPast) {
      // Completed games have results so you can compute accuracy immediately
      homeScore = clamp(Math.round(randFloat(1, 8) + noise(0, 2.0)), 0, 18);
      awayScore = clamp(Math.round(randFloat(1, 8) + noise(0, 2.0)), 0, 18);

      // Avoid ties for winner logic (baseball doesn’t tie)
      if (homeScore === awayScore) homeScore += 1;

      winner = homeScore > awayScore ? homeTeam.id : awayTeam.id;
    }

    games.push({
      gameId: `${DATASET_ID}-game-${String(i + 1).padStart(4, '0')}`,
      date: gameDate.toISOString(),
      homeTeam: homeTeam.name,
      homeTeamId: homeTeam.id,
      awayTeam: awayTeam.name,
      awayTeamId: awayTeam.id,
      homeScore,
      awayScore,
      status: isPast ? 'completed' : 'upcoming',
      ...(winner ? { winner } : {}),
    });
  }

  return games;
}

function pickPredictedWinner(game: SeedGame, persona: Persona): string {
  if (game.status !== 'completed' || !game.winner) {
    // for upcoming games, choose based on a simple “home advantage” coin flip
    return rand() < 0.52 ? game.homeTeamId : game.awayTeamId;
  }

  const pickCorrect = rand() < persona.winnerSkill;
  if (pickCorrect) return game.winner;

  // pick the other team
  return game.winner === game.homeTeamId ? game.awayTeamId : game.homeTeamId;
}

function makePredictionScores(game: SeedGame, persona: Persona) {
  // If the game is completed, center around actual scores; else use plausible priors
  const baseHome = game.status === 'completed' ? game.homeScore : randFloat(2.5, 5.5);
  const baseAway = game.status === 'completed' ? game.awayScore : randFloat(2.5, 5.5);

  const home = clamp(baseHome + noise(0, persona.scoreStdev), 0, 20);
  const away = clamp(baseAway + noise(0, persona.scoreStdev), 0, 20);

  // Round to 1 decimal like your current schema
  const home1 = parseFloat(home.toFixed(1));
  const away1 = parseFloat(away.toFixed(1));

  return {
    homeScorePredicted: home1,
    awayScorePredicted: away1,
    totalScorePredicted: parseFloat((home1 + away1).toFixed(1)),
  };
}

function makeConfidence(game: SeedGame, persona: Persona, predictedWinnerId: string): number {
  // If completed, “confidence” loosely follows correctness (helps UX/testing)
  if (game.status === 'completed' && game.winner) {
    const correct = predictedWinnerId === game.winner;
    if (persona.overconfident) {
      return parseFloat(randFloat(0.85, 0.99).toFixed(2));
    }
    return parseFloat(
      randFloat(correct ? persona.confMin : 0.45, correct ? persona.confMax : 0.65).toFixed(2)
    );
  }
  return parseFloat(randFloat(persona.confMin, persona.confMax).toFixed(2));
}

type SeedPrediction = {
  userId: string;
  gameId: string;
  homeScorePredicted: number;
  awayScorePredicted: number;
  totalScorePredicted: number;
  confidence: number;
  predictedWinnerId: string;
  submittedAt: string;
  // Scored fields — only present for completed games
  actualWinnerId?: string;
  winnerCorrect?: boolean;
  homeScoreError?: number;
  awayScoreError?: number;
  totalScoreError?: number;
};

function scorePrediction(pred: SeedPrediction, game: SeedGame) {
  if (game.status !== 'completed' || game.winner === undefined) return pred;

  const winnerCorrect = pred.predictedWinnerId === game.winner;
  const homeScoreError = Math.abs(pred.homeScorePredicted - game.homeScore);
  const awayScoreError = Math.abs(pred.awayScorePredicted - game.awayScore);
  const totalScoreError = Math.abs(pred.totalScorePredicted - (game.homeScore + game.awayScore));

  return {
    ...pred,
    actualWinnerId: game.winner,
    winnerCorrect,
    homeScoreError: parseFloat(homeScoreError.toFixed(4)),
    awayScoreError: parseFloat(awayScoreError.toFixed(4)),
    totalScoreError: parseFloat(totalScoreError.toFixed(4)),
  };
}

function makePredictions(users: SeedUser[], games: SeedGame[]): SeedPrediction[] {
  const predictions: SeedPrediction[] = [];
  const submittedAt = new Date().toISOString();

  const gamesToPredict = games.slice(0, Math.min(NUM_PREDICTION_GAMES, games.length));

  for (const game of gamesToPredict) {
    users.forEach((u) => {
      const persona = personas.find((p) => p.name === u.username) ?? personas[0];

      const predictedWinnerId = pickPredictedWinner(game, persona);
      const scores = makePredictionScores(game, persona);
      const confidence = makeConfidence(game, persona, predictedWinnerId);

      const prediction: SeedPrediction = {
        userId: u.userId,
        gameId: game.gameId,
        ...scores,
        confidence,
        predictedWinnerId,
        submittedAt,
      };

      predictions.push(scorePrediction(prediction, game));
    });
  }

  return predictions;
}

// ---------- Seed runner ----------
async function seed() {
  console.log('Seeding database...\n');
  console.log(`Endpoint: ${ENDPOINT}`);
  console.log(`Region: ${REGION}`);
  console.log(`Dataset: ${DATASET_ID}`);
  console.log(`Random seed: ${RANDOM_SEED}`);
  console.log(`Tables:`);
  console.log(`  USERS_TABLE=${USERS_TABLE}`);
  console.log(`  GAMES_TABLE=${GAMES_TABLE}`);
  console.log(`  PREDICTIONS_TABLE=${PREDICTIONS_TABLE}\n`);

  const users = makeUsers();
  const games = makeGames();
  const predictions = makePredictions(users, games);

  // Write in batches + retry
  await batchWriteAll(USERS_TABLE, users);
  await batchWriteAll(GAMES_TABLE, games);
  await batchWriteAll(PREDICTIONS_TABLE, predictions);

  console.log('\nSeed summary:');
  console.log(`  users: ${users.length}`);
  console.log(`  games: ${games.length} (completed: ${games.filter((g) => g.status === 'completed').length})`);
  console.log(`  predictions: ${predictions.length}`);
  console.log('\nDatabase seeded successfully!');
}

seed().then(
  () => process.exit(0),
  (err) => {
    console.error('Error seeding database:', err);
    process.exit(1);
  }
);
import json
from urllib import response
import boto3
import requests
from datetime import datetime, timedelta
from typing import List, Dict
import statsapi
from bs4 import BeautifulSoup

dynamodb = boto3.client('dynamodb')
s3 = boto3.client('s3')

GAMES_TABLE = "mlb-prediction-pool-games"
DATA_BUCKET = "mlb-pool-data"

def lambda_handler(event, context):
    '''
    Fetches daily MLB data and stores it in DynamoDB and S3.
    '''
    try:
        games = fetch_upcoming_games()
        print(f"Fetched {len(games)} upcoming games.")

        store_games_in_dynamodb(games)
        print("Stored games in DynamoDB.")

        team_stats = fetch_team_stats()
        print("Fetched team stats.")

        store_data_in_s3(team_stats)
        print("Stored team stats in S3.")

        return {
            "statusCode": 200,
            "body": json.dumps({
                "message": "Data ingestion completed successfully",
                "games_fetched": len(games),
                "teams_updated": len(team_stats),
            })
        }
    except Exception as e:
        print(f"Error fetching games: {e}")
        return {"statusCode": 500, "body": "Error fetching games"}
    

def fetch_upcoming_games() -> List[Dict]:
    """
    Fetches games for the next 3 days from MLB Stats API.
    """
    today = datetime.now()
    end_date = today + timedelta(days=3)
    
    game_data: List[Dict] = statsapi.schedule(
        start_date=today.strftime("%Y-%m-%d"),
        end_date=end_date.strftime("%Y-%m-%d"),
        include_series_status=False
    )

    games = []

    for game in game_data:
        game_info = {
            'gameId': str(game.get('game_id')),
            'date': game.get('game_date'),
            'homeTeamId': game.get('home_id'),
            'homeTeam': game.get('home_name'),
            'awayTeamId': game.get('away_id'),
            'awayTeam': game.get('away_name'),
            'status': str(game.get('status'))
        }

        games.append(game_info)

    return games


def fetch_team_stats() -> Dict:
    stats_url = "https://www.baseball-reference.com/leagues/majors/2025.shtml"
    response = requests.get(stats_url)

    soup = BeautifulSoup(response.content, 'html.parser')

    batting_table = soup.find('table', {'id': 'teams_standard_batting'})
    pitching_table = soup.find('table', {'id': 'teams_standard_pitching'})
    fielding_table = soup.find('table', {'id': 'teams_standard_fielding'})

    team_stats = {}

    for table, stat_type in [(batting_table, 'batting'), (pitching_table, 'pitching'), (fielding_table, 'fielding')]:
        if table:
            rows = table.find('tbody').find_all('tr') # type: ignore
            for row in rows:
                team_name_cell = row.find('th', {'data-stat': 'team_name'})
                if team_name_cell:
                    team_name = team_name_cell.text.strip()
                    stats = {}
                    for cell in row.find_all('td'):
                        stat_name = cell['data-stat']
                        stat_value = cell.text.strip()
                        stats[stat_name] = stat_value
                    if team_name not in team_stats:
                        team_stats[team_name] = {}
                    team_stats[team_name][stat_type] = stats

    return team_stats

def fetch_standings() -> Dict:
    standings = statsapi.standings_data()

    return standings

def store_games_in_dynamodb(games: List[Dict]):
    """Store fetched games in DynamoDB"""
    table = dynamodb.Table(GAMES_TABLE)

    with table.batch_writer() as batch:
        for game in games:
            batch.put_item(Item=game)

def store_data_in_s3(team_stats: Dict):
    date_str = datetime.now().strftime("%Y-%m-%d")
    # Store each day
    s3.put_object(
        Bucket = DATA_BUCKET,
        Key = f"team-stats/{date_str}.json",
        Body = json.dumps(team_stats, indent=2),    
        ContentType = "application/json"
    )

    # Store as latest for easier lookup
    s3.put_object(
        Bucket=DATA_BUCKET,
        Key = "team-stats/latest.json",
        Body = json.dumps(team_stats, indent=2),
        ContentType = "application/json"
    )

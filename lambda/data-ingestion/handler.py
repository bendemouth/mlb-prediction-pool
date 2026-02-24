import json
import os
import boto3
import requests
from datetime import datetime, timedelta
from typing import List, Dict
import statsapi

dynamodb = boto3.resource('dynamodb')
s3 = boto3.client('s3')

GAMES_TABLE = os.environ.get("GAMES_TABLE", "mlb-prediction-pool-dev-games")
DATA_BUCKET = os.environ.get("DATA_BUCKET", "mlb-prediction-pool-dev-mlb-data")
MLB_API_BASE = os.environ.get("MLB_API_BASE", "https://statsapi.mlb.com/api/v1")

def lambda_handler(event, context):
    '''
    Fetches daily MLB data and stores it in DynamoDB and S3.
    '''
    try:
        games = fetch_upcoming_games()
        print(f"Fetched {len(games)} upcoming games.")

        store_games_in_dynamodb(games)
        print("Stored games in DynamoDB.")

        fetch_and_store_all_stats()
        print("Fetched and stored team stats in S3.")

        return {
            "statusCode": 200,
            "body": json.dumps({
                "message": "Data ingestion completed successfully",
                "games_fetched": len(games),
                "teams_updated": len(fetch_team_hitting_stats()),
            })
        }
    except Exception as e:
        print(f"Error fetching games: {e}")
        return {"statusCode": 500, "body": f"Error fetching games: {e}"}
    

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

def fetch_team_hitting_stats() -> Dict:
    """Fetches team hitting stats for the current season from MLB Stats API."""
    url = build_mlb_api_url("teams/stats", {
        "season": "2025",
        "group": "hitting",
        "sportIds": "1"
    })

    response = requests.get(url)

    if response.status_code != 200:
        raise Exception(f"Failed to fetch team hitting stats: {response.status_code} - {response.text}")

    data = response.json()

    team_stats = parse_team_stats_response(data)
    print(f"Fetched hitting stats for {len(team_stats)} teams.")

    return team_stats

def fetch_team_pitching_stats() -> Dict:
    """Fetches team pitching stats for the current season from MLB Stats API."""
    url = build_mlb_api_url("teams/stats", {
        "season": "2025",
        "group": "pitching",
        "sportIds": "1"
    })

    response = requests.get(url)

    if response.status_code != 200:
        raise Exception(f"Failed to fetch team pitching stats: {response.status_code} - {response.text}")
    
    data = response.json()

    team_stats = parse_team_stats_response(data)
    print(f"Fetched pitching stats for {len(team_stats)} teams.")

    return team_stats

def fetch_team_fielding_stats() -> Dict:
    """Fetches team fielding stats for the current season from MLB Stats API."""
    url = build_mlb_api_url("teams/stats", {
        "season": "2025",
        "group": "fielding",
        "sportIds": "1"
    })

    response = requests.get(url)
    if response.status_code != 200:
        raise Exception(f"Failed to fetch team fielding stats: {response.status_code} - {response.text}")
    
    data = response.json()

    team_stats = parse_team_stats_response(data)
    print(f"Fetched fielding stats for {len(team_stats)} teams.")

    return team_stats

def fetch_team_catching_stats() -> Dict:
    """Fetches team catching stats for the current season from MLB Stats API."""
    url = build_mlb_api_url("teams/stats", {
        "season": "2025",
        "group": "catching",
        "sportIds": "1"
    })

    response = requests.get(url)
    if response.status_code != 200:
        raise Exception(f"Failed to fetch team catching stats: {response.status_code} - {response.text}")
    
    data = response.json()

    team_stats = parse_team_stats_response(data)
    print(f"Fetched catching stats for {len(team_stats)} teams.")

    return team_stats

def fetch_and_store_all_stats():
    """Fetches all team stats and stores them in S3."""
    team_hitting_stats = fetch_team_hitting_stats()
    team_pitching_stats = fetch_team_pitching_stats()
    team_fielding_stats = fetch_team_fielding_stats()
    team_catching_stats = fetch_team_catching_stats()

    store_data_in_s3(team_hitting_stats, "hitting")
    store_data_in_s3(team_pitching_stats, "pitching")
    store_data_in_s3(team_fielding_stats, "fielding")
    store_data_in_s3(team_catching_stats, "catching")

def fetch_standings() -> Dict:
    """Fetches current standings from MLB Stats API."""
    standings = statsapi.standings_data()

    return standings

def store_games_in_dynamodb(games: List[Dict]):
    """Store fetched games in DynamoDB"""
    table = dynamodb.Table(GAMES_TABLE) # type: ignore

    with table.batch_writer() as batch:
        for game in games:
            batch.put_item(Item=game)

def store_data_in_s3(team_stats: Dict, category: str):
    date_str = datetime.now().strftime("%Y-%m-%d")
    # Store each day
    s3.put_object(
        Bucket = DATA_BUCKET,
        Key = f"team-stats/{category}/{date_str}.json",
        Body = json.dumps(team_stats, indent=2),    
        ContentType = "application/json"
    )

    # Store as latest for easier lookup
    s3.put_object(
        Bucket=DATA_BUCKET,
        Key = f"team-stats/{category}/latest.json",
        Body = json.dumps(team_stats, indent=2),
        ContentType = "application/json"
    )

def build_mlb_api_url(endpoint: str, params: Dict[str, str]) -> str:
    query_string = '&'.join([f"{key}={value}" for key, value in params.items()])
    return f"{MLB_API_BASE}/{endpoint}?{query_string}"

def parse_team_stats_response(data: Dict) -> Dict:
    team_stats = {}
    stats_blocks = data.get("stats", [])
    splits = stats_blocks[0].get("splits", []) if stats_blocks else []

    for split in splits:
        team = split.get("team", {})
        stat = split.get("stat", {})
        team_id = team.get("id")

        if team_id:
            team_stats[team_id] = {
                "teamId": team_id,
                "teamName": team.get("name"),
                "season": split.get("season"),
                "rank": split.get("rank"),
                **stat
            }

    return team_stats
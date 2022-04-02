docker-compose down
docker stop nba-simulation
docker rm nba-simulation
docker rmi nba-simulation
docker build -t nba-simulation .
docker-compose up
docker build --file=frontend/Dockerfile  -t social-network-frontend .
docker build --file=backend/Dockerfile  -t social-network-backend .
docker-compose -f docker-compose.yml up
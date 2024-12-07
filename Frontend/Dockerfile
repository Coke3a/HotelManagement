FROM node:20-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm install

COPY . .

FROM node:20-alpine
WORKDIR /app

# Development stage
FROM builder AS development
ENV NODE_ENV=development
EXPOSE 9500
CMD ["npm", "start"]

# Production stage
FROM builder AS production
ENV NODE_ENV=production
RUN npm run build
EXPOSE 9500
CMD ["npm", "start"]

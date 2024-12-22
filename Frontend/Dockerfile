FROM node:20-alpine AS builder

# Add timezone data
RUN apk add --no-cache tzdata

# Set timezone
ENV TZ=Asia/Bangkok

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm install

COPY . .

FROM node:20-alpine
WORKDIR /app

# Development stage
FROM builder AS development
ENV NODE_ENV=development
ENV TZ=Asia/Bangkok
RUN apk add --no-cache tzdata
EXPOSE 9500
CMD ["npm", "start"]

# Production stage
FROM builder AS production
ENV NODE_ENV=production
ENV TZ=Asia/Bangkok
RUN apk add --no-cache tzdata
RUN npm run build
EXPOSE 9500
CMD ["npm", "start"]

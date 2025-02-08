FROM node:22.11.0-slim AS build

WORKDIR /app

COPY package*.json .

# RUN npm ci
RUN npm install

COPY . .
RUN npm run build

FROM node:22.11.0-slim AS run

ENV NODE_ENV=production

WORKDIR /app
COPY --from=build /app/build/ ./build
COPY --from=build /app/package.json ./package.json
COPY --from=build /app/node_modules ./node_modules
RUN ulimit -c unlimited

EXPOSE 80
ENTRYPOINT ["node", "build", "-l", "80"]
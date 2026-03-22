ARG PLAYWRIGHT_IMAGE_TAG=v1.58.2-noble
FROM mcr.microsoft.com/playwright:${PLAYWRIGHT_IMAGE_TAG}

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .

CMD ["npm", "run", "test"]

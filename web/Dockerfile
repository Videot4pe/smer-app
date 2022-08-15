FROM node:17

WORKDIR /app

COPY package.json ./
COPY . .

RUN yarn
RUN yarn build

CMD ["yarn", "serve"]

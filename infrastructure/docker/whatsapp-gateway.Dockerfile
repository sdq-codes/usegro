FROM node:20-alpine
WORKDIR /app
COPY services/whatsapp-gateway/package*.json ./
RUN npm install --omit=dev
COPY services/whatsapp-gateway/src/ ./src/
COPY proto/ ./proto/
EXPOSE 3000
CMD ["node", "src/index.js"]

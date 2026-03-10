# ---------- Dev Stage ----------
FROM node:22-alpine

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install all dependencies including dev
RUN npm install --legacy-peer-deps

# Copy the app
COPY . .

# Expose Vite and Nuxt ports
EXPOSE 3000 5173 8098

# Environment
ENV NODE_ENV=development
ENV NUXT_HOST=0.0.0.0
ENV NUXT_PORT=3000

# Start Nuxt dev server (enables HMR + Vue Devtools)
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]

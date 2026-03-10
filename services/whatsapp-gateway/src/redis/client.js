const Redis = require('ioredis');
const { logger } = require('../utils/logger');

let client;

async function connectRedis() {
  client = new Redis({
    host: process.env.REDIS_HOST || 'localhost',
    port: parseInt(process.env.REDIS_PORT) || 6379,
    password: process.env.REDIS_PASSWORD || undefined,
    retryStrategy: (times) => Math.min(times * 100, 3000),
  });

  client.on('error', (err) => logger.error('Redis error', { error: err.message }));
  client.on('connect', () => logger.info('Redis connected'));

  await client.ping();
  return client;
}

function getClient() {
  if (!client) throw new Error('Redis not connected');
  return client;
}

module.exports = { connectRedis, getClient };

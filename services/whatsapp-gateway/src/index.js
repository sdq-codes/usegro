require('dotenv').config();
const express = require('express');
const helmet = require('helmet');
const rateLimit = require('express-rate-limit');
const { logger } = require('./utils/logger');
const webhookRouter = require('./webhook/router');
const { connectRedis } = require('./redis/client');
const { startGrpcServer } = require('./grpc/server');

const app = express();
const PORT = process.env.WHATSAPP_GATEWAY_PORT || 3000;

// ─── Middleware ────────────────────────────────────────────
app.use(helmet());
app.use(express.json());

app.use(rateLimit({
  windowMs: 60 * 1000,
  max: 500,
  message: 'Too many requests from this IP'
}));

// ─── Routes ────────────────────────────────────────────────
app.use('/webhook', webhookRouter);

app.get('/health', (req, res) => {
  res.json({ status: 'ok', service: 'whatsapp-gateway', ts: new Date() });
});

// ─── Boot ──────────────────────────────────────────────────
async function boot() {
  try {
    await connectRedis();
    logger.info('Redis connected');

    app.listen(PORT, () => {
      logger.info(`WhatsApp Gateway HTTP listening on :${PORT}`);
    });

    await startGrpcServer();
    logger.info('gRPC server started');
  } catch (err) {
    logger.error('Boot failed', { error: err.message });
    process.exit(1);
  }
}

boot();

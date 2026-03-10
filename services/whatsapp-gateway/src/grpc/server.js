const path = require('path');
const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');
const { sendMessage, buildTextPayload, buildTemplatePayload } = require('../sender/meta');
const { logger } = require('../utils/logger');

const PROTO_PATH = path.join(__dirname, '../../proto/whatsapp.proto');

async function startGrpcServer() {
  const packageDef = protoLoader.loadSync(PROTO_PATH, {
    keepCase: true,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true,
  });

  const server = new grpc.Server();
  const port = process.env.WHATSAPP_GRPC_PORT || 4000;

  await new Promise((resolve, reject) => {
    server.bindAsync(
      `0.0.0.0:${port}`,
      grpc.ServerCredentials.createInsecure(),
      (err, port) => {
        if (err) return reject(err);
        logger.info(`WhatsApp Gateway gRPC listening on :${port}`);
        resolve(port);
      }
    );
  });

  return server;
}

module.exports = { startGrpcServer };

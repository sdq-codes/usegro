const { getClient } = require('./client');

// Channel naming convention: usegro:{event_type}:{waba_id}
// Go services subscribe to these channels

const CHANNELS = {
  INBOUND_MESSAGE: (wabaId) => `usegro:inbound:message:${wabaId}`,
  STATUS_UPDATE:   (wabaId) => `usegro:inbound:status:${wabaId}`,
  SEND_MESSAGE:    (wabaId) => `usegro:outbound:send:${wabaId}`,  // Go → Node.js
};

async function publishInboundMessage(wabaId, event) {
  const redis = getClient();
  await redis.publish(CHANNELS.INBOUND_MESSAGE(wabaId), JSON.stringify(event));
}

async function publishStatusUpdate(wabaId, event) {
  const redis = getClient();
  await redis.publish(CHANNELS.STATUS_UPDATE(wabaId), JSON.stringify(event));
}

module.exports = { publishInboundMessage, publishStatusUpdate, CHANNELS };

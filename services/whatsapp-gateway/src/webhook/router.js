const express = require('express');
const crypto = require('crypto');
const { logger } = require('../utils/logger');
const { publishInboundMessage, publishStatusUpdate } = require('../redis/publisher');

const router = express.Router();

// ─── Meta Webhook Verification (GET) ──────────────────────
router.get('/', (req, res) => {
  const mode      = req.query['hub.mode'];
  const token     = req.query['hub.verify_token'];
  const challenge = req.query['hub.challenge'];

  if (mode === 'subscribe' && token === process.env.WHATSAPP_WEBHOOK_VERIFY_TOKEN) {
    logger.info('Webhook verified by Meta');
    return res.status(200).send(challenge);
  }

  logger.warn('Webhook verification failed', { mode, token });
  return res.sendStatus(403);
});

// ─── Inbound Events (POST) ─────────────────────────────────
router.post('/', verifySignature, async (req, res) => {
  // Always respond 200 immediately — Meta will retry if we don't
  res.sendStatus(200);

  try {
    const body = req.body;
    if (body.object !== 'whatsapp_business_account') return;

    for (const entry of body.entry || []) {
      const wabaId = entry.id;

      for (const change of entry.changes || []) {
        if (change.field !== 'messages') continue;

        const value = change.value;
        const phoneNumberId = value.metadata?.phone_number_id;

        // ── Inbound messages ────────────────────────────────
        for (const msg of value.messages || []) {
          const contact = value.contacts?.find(c => c.wa_id === msg.from);

          const event = {
            message_id:      msg.id,
            waba_id:         wabaId,
            phone_number_id: phoneNumberId,
            from:            msg.from,
            to:              value.metadata?.display_phone_number,
            type:            msg.type,
            content:         extractContent(msg),
            timestamp:       new Date(parseInt(msg.timestamp) * 1000).toISOString(),
            contact: {
              wa_id: contact?.wa_id || msg.from,
              name:  contact?.profile?.name || '',
            }
          };

          logger.info('Inbound message', { messageId: msg.id, from: msg.from, type: msg.type });
          await publishInboundMessage(wabaId, event);
        }

        // ── Status updates ──────────────────────────────────
        for (const status of value.statuses || []) {
          const event = {
            message_id:   status.id,
            waba_id:      wabaId,
            status:       status.status,  // sent, delivered, read, failed
            timestamp:    new Date(parseInt(status.timestamp) * 1000).toISOString(),
            recipient_id: status.recipient_id,
          };

          logger.info('Status update', { messageId: status.id, status: status.status });
          await publishStatusUpdate(wabaId, event);
        }
      }
    }
  } catch (err) {
    logger.error('Error processing webhook', { error: err.message, stack: err.stack });
  }
});

// ─── Helpers ───────────────────────────────────────────────

function extractContent(msg) {
  switch (msg.type) {
    case 'text':
      return { text: { body: msg.text?.body } };
    case 'image':
    case 'video':
    case 'audio':
    case 'document':
      return {
        media: {
          media_id:  msg[msg.type]?.id,
          mime_type: msg[msg.type]?.mime_type,
          sha256:    msg[msg.type]?.sha256,
          caption:   msg[msg.type]?.caption,
          filename:  msg[msg.type]?.filename,
        }
      };
    case 'location':
      return {
        location: {
          latitude:  msg.location?.latitude,
          longitude: msg.location?.longitude,
          name:      msg.location?.name,
          address:   msg.location?.address,
        }
      };
    case 'interactive':
      if (msg.interactive?.type === 'button_reply') {
        return { button: { id: msg.interactive.button_reply.id, title: msg.interactive.button_reply.title } };
      }
      return { interactive: { id: msg.interactive?.list_reply?.id, title: msg.interactive?.list_reply?.title } };
    default:
      return {};
  }
}

function verifySignature(req, res, next) {
  const sig = req.headers['x-hub-signature-256'];
  if (!sig) {
    logger.warn('Missing signature header');
    return res.sendStatus(401);
  }

  const expected = 'sha256=' + crypto
    .createHmac('sha256', process.env.WHATSAPP_APP_SECRET)
    .update(JSON.stringify(req.body))
    .digest('hex');

  if (!crypto.timingSafeEqual(Buffer.from(sig), Buffer.from(expected))) {
    logger.warn('Invalid webhook signature');
    return res.sendStatus(401);
  }

  next();
}

module.exports = router;

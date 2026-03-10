const axios = require('axios');
const { logger } = require('../utils/logger');

const BASE_URL = `https://graph.facebook.com/${process.env.WHATSAPP_API_VERSION || 'v19.0'}`;

/**
 * Send a message via Meta Cloud API
 * @param {string} phoneNumberId - The sender's phone number ID
 * @param {string} accessToken   - The WABA access token (per tenant)
 * @param {object} payload       - Cloud API message payload
 */
async function sendMessage(phoneNumberId, accessToken, payload) {
  try {
    const response = await axios.post(
      `${BASE_URL}/${phoneNumberId}/messages`,
      payload,
      {
        headers: {
          Authorization: `Bearer ${accessToken}`,
          'Content-Type': 'application/json',
        },
        timeout: 10000,
      }
    );

    logger.info('Message sent', {
      phoneNumberId,
      messageId: response.data?.messages?.[0]?.id,
    });

    return response.data;
  } catch (err) {
    const errData = err.response?.data?.error;
    logger.error('Failed to send message', {
      phoneNumberId,
      errorCode: errData?.code,
      errorMessage: errData?.message,
    });
    throw err;
  }
}

// ─── Payload builders ──────────────────────────────────────

function buildTextPayload(to, text) {
  return { messaging_product: 'whatsapp', recipient_type: 'individual', to, type: 'text', text: { body: text, preview_url: false } };
}

function buildTemplatePayload(to, templateName, languageCode, components = []) {
  return {
    messaging_product: 'whatsapp', recipient_type: 'individual', to, type: 'template',
    template: { name: templateName, language: { code: languageCode }, components }
  };
}

function buildImagePayload(to, mediaId, caption = '') {
  return { messaging_product: 'whatsapp', recipient_type: 'individual', to, type: 'image', image: { id: mediaId, caption } };
}

function buildInteractiveButtonPayload(to, bodyText, buttons) {
  return {
    messaging_product: 'whatsapp', recipient_type: 'individual', to, type: 'interactive',
    interactive: {
      type: 'button',
      body: { text: bodyText },
      action: {
        buttons: buttons.map(b => ({ type: 'reply', reply: { id: b.id, title: b.title } }))
      }
    }
  };
}

function buildInteractiveListPayload(to, bodyText, buttonText, sections) {
  return {
    messaging_product: 'whatsapp', recipient_type: 'individual', to, type: 'interactive',
    interactive: {
      type: 'list',
      body: { text: bodyText },
      action: { button: buttonText, sections }
    }
  };
}

// ─── Media upload ──────────────────────────────────────────

async function uploadMedia(phoneNumberId, accessToken, fileBuffer, mimeType) {
  const FormData = require('form-data');
  const form = new FormData();
  form.append('file', fileBuffer, { contentType: mimeType });
  form.append('messaging_product', 'whatsapp');

  const response = await axios.post(
    `${BASE_URL}/${phoneNumberId}/media`,
    form,
    {
      headers: { ...form.getHeaders(), Authorization: `Bearer ${accessToken}` }
    }
  );

  return response.data.id;
}

module.exports = {
  sendMessage,
  buildTextPayload,
  buildTemplatePayload,
  buildImagePayload,
  buildInteractiveButtonPayload,
  buildInteractiveListPayload,
  uploadMedia,
};

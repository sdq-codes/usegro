<!-- /oauth/callback — opened inside a popup by the Google login flow -->
<script lang="ts" setup>
// No middleware — this page must be publicly accessible
definePageMeta({ middleware: [] })

if (process.client) {
  const params = new URLSearchParams(window.location.search)
  const accessToken = params.get('access_token')

  // Strip tokens from the URL immediately to avoid leaking them in history
  window.history.replaceState({}, '', window.location.pathname)

  if (accessToken && window.opener) {
    // The refresh token is in an HttpOnly cookie set by the server redirect —
    // it does not need to be passed through postMessage.
    window.opener.postMessage(
      { type: 'GOOGLE_AUTH_SUCCESS', accessToken },
      window.location.origin
    )
  } else if (window.opener) {
    const error = params.get('error') || 'auth_failed'
    window.opener.postMessage(
      { type: 'GOOGLE_AUTH_ERROR', error },
      window.location.origin
    )
  }

  window.close()
}
</script>

<template>
  <div />
</template>

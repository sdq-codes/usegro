<!-- /oauth/callback — opened inside a popup by the Google login flow -->
<script lang="ts" setup>
// No middleware — this page must be publicly accessible
definePageMeta({ middleware: [] })

if (process.client) {
  const params = new URLSearchParams(window.location.search)
  const accessToken = params.get('access_token')
  const refreshToken = params.get('refresh_token')

  // Strip tokens from the URL immediately to avoid leaking them in history
  window.history.replaceState({}, '', window.location.pathname)

  if (accessToken && refreshToken && window.opener) {
    window.opener.postMessage(
      { type: 'GOOGLE_AUTH_SUCCESS', accessToken, refreshToken },
      window.location.origin
    )
  } else if (window.opener) {
    window.opener.postMessage(
      { type: 'GOOGLE_AUTH_ERROR' },
      window.location.origin
    )
  }

  window.close()
}
</script>

<template>
  <div />
</template>

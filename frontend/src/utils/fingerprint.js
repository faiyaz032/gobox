import FingerprintJS from '@fingerprintjs/fingerprintjs';

let cachedFingerprint = null;

/**
 * Get a stable browser fingerprint using FingerprintJS.
 * The fingerprint persists across browser restarts and PC restarts.
 * Result is cached in memory for the session to avoid re-computation.
 */
export async function getFingerprint() {
  if (cachedFingerprint) {
    return cachedFingerprint;
  }

  const fp = await FingerprintJS.load();
  const result = await fp.get();
  cachedFingerprint = result.visitorId;

  console.log('[GoBox] Browser fingerprint generated:', cachedFingerprint);
  return cachedFingerprint;
}

const RAW_BASE = (import.meta.env.VITE_PUBLIC_WEB_URL ?? "").trim();

const LOCAL_HOSTS = /^https?:\/\/(localhost|127\.0\.0\.1|\[::1\])(:\d+)?(\/|$)/i;

// Normaliza la base pública configurada. Devuelve null si falta o no es una URL.
export function getPublicBaseUrl(): string | null {
  if (!RAW_BASE) return null;
  const normalized = RAW_BASE.replace(/\/+$/, "");
  try {
    const url = new URL(normalized);
    if (url.protocol !== "https:" && url.protocol !== "http:") return null;
    return normalized;
  } catch {
    return null;
  }
}

// En producción un QR que apunta a localhost es inservible: se bloquean la
// descarga y la impresión, pero se puede seguir copiando el enlace en desarrollo.
export function canPrintQr(): boolean {
  const base = getPublicBaseUrl();
  if (!base) return false;
  if (import.meta.env.PROD && LOCAL_HOSTS.test(base)) return false;
  return true;
}

export function getQrUrlWarning(): string | null {
  const base = getPublicBaseUrl();
  if (!base) {
    return "Configura VITE_PUBLIC_WEB_URL para descargar o imprimir el QR.";
  }
  if (!canPrintQr()) {
    return "VITE_PUBLIC_WEB_URL apunta a localhost en producción. El QR impreso no funcionaría.";
  }
  return null;
}

export function buildQrUrl(token: string): string | null {
  const base = getPublicBaseUrl();
  if (!base) return null;
  return `${base}/q/${token}`;
}

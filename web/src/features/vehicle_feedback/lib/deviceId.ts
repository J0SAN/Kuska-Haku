const STORAGE_KEY = "kuskahaku_device_id";
const DEVICE_ID_PATTERN = /^[a-f0-9]{32}$/;

function generateDeviceId(): string {
  const bytes = new Uint8Array(16);
  // getRandomValues existe también fuera de contextos seguros; randomUUID no.
  crypto.getRandomValues(bytes);
  return Array.from(bytes, (byte) => byte.toString(16).padStart(2, "0")).join(
    "",
  );
}

// Identificador anónimo del navegador para evitar calificar dos veces el mismo
// vehículo el mismo día. En navegación privada sin localStorage se usa un valor
// temporal que solo dura la sesión de la página.
export function getDeviceId(): string {
  try {
    const stored = localStorage.getItem(STORAGE_KEY);
    if (stored && DEVICE_ID_PATTERN.test(stored)) {
      return stored;
    }
    const generated = generateDeviceId();
    localStorage.setItem(STORAGE_KEY, generated);
    return generated;
  } catch {
    return generateDeviceId();
  }
}

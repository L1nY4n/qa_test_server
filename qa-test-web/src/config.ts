const { hostname, port } = window.location;
const url =
  import.meta.env.DEV ? import.meta.env.VITE_APP_BASE_URL:`${hostname}:${port}`
export const base_url = `http://${url}`
export const requestTimeout = 30000
export const websocket_url = `ws://${url}`


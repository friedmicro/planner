const BASE = 'http://' + window.location.hostname + ":" + window.location.port;

async function req(method, path, body) {
  const res = await fetch(`${BASE}${path}`, {
    method,
    headers: body ? { 'Content-Type': 'application/json' } : {},
    body: body ? JSON.stringify(body) : undefined,
  })
  if (!res.ok) {
    let msg = `${method} ${path} failed (${res.status})`
    try {
      const j = await res.json()
      if (j.error) msg = j.error
    } catch (_) {}
    throw new Error(msg)
  }
  const text = await res.text()
  return text ? JSON.parse(text) : null
}

export const api = {
  config: {
    get: () => req('GET', '/api/config'),
    save: (c) => req('PUT', '/api/config', c),
  },
  sleep: {
    list: () => req('GET', '/api/sleep'),
    add: (b) => req('POST', '/api/sleep', b),
    del: (id) => req('DELETE', `/api/sleep/${id}`),
  },
  blacklist: {
    list: () => req('GET', '/api/blacklist'),
    add: (b) => req('POST', '/api/blacklist', b),
    del: (id) => req('DELETE', `/api/blacklist/${id}`),
  },
  events: {
    list: () => req('GET', '/api/events'),
    add: (e) => req('POST', '/api/events', e),
    update: (e) => req('PUT', `/api/events/${e.id}`, e),
    del: (id) => req('DELETE', `/api/events/${id}`),
  },
  extras: {
    list: () => req('GET', '/api/extras'),
    add: (x) => req('POST', '/api/extras', x),
    del: (id) => req('DELETE', `/api/extras/${id}`),
  },
  solve: () => req('POST', '/api/solve'),
  schedule: () => req('GET', '/api/schedule'),
}

export const DAYS = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']

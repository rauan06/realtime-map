// Realtime map client.
// Connects to /ws and renders 5 layers: flights, ships, transport, roads, OBU.
// Marker identity is keyed by (layer, id) where id is the source-specific ID
// from the upstream record (icao24 / mmsi / vehicle_id / road_id / device_id).

(function () {
  const map = L.map('map').setView([48.0, 66.0], 4); // Kazakhstan centred default
  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '© OpenStreetMap contributors',
    maxZoom: 19,
  }).addTo(map);

  const groups = {
    flight: L.layerGroup().addTo(map),
    ship: L.layerGroup().addTo(map),
    transport: L.layerGroup().addTo(map),
    road: L.layerGroup().addTo(map),
    obu: L.layerGroup().addTo(map),
  };

  const features = { flight: new Map(), ship: new Map(), transport: new Map(), road: new Map(), obu: new Map() };
  const counts = { flight: 0, ship: 0, transport: 0, road: 0, obu: 0 };

  // Layer-specific icons (simple emoji DIV icons keep us tile-server-free).
  function iconFor(layer) {
    const map = { flight: '✈️', ship: '🚢', transport: '🚌', obu: '📍' };
    const emoji = map[layer] || '•';
    return L.divIcon({
      html: `<div style="font-size:18px;line-height:18px;text-shadow:0 1px 2px rgba(0,0,0,0.4)">${emoji}</div>`,
      className: '',
      iconSize: [22, 22],
      iconAnchor: [11, 11],
    });
  }

  const severityColor = {
    critical: '#c0392b',
    high: '#e67e22',
    medium: '#f1c40f',
    low: '#2ecc71',
  };

  function popupTable(rows) {
    const tr = rows.filter(r => r[1] != null && r[1] !== '').map(([k, v]) => `<tr><td>${k}</td><td>${v}</td></tr>`).join('');
    return `<table>${tr}</table>`;
  }

  function upsertPoint(layer, id, lat, lng, popupHTML) {
    const existing = features[layer].get(id);
    if (existing) {
      existing.setLatLng([lat, lng]);
      if (popupHTML) existing.bindPopup(popupHTML);
      return;
    }
    const marker = L.marker([lat, lng], { icon: iconFor(layer) });
    if (popupHTML) marker.bindPopup(popupHTML);
    marker.addTo(groups[layer]);
    features[layer].set(id, marker);
    counts[layer]++;
    document.getElementById('count-' + layer).textContent = counts[layer];
  }

  function upsertPolyline(layer, id, coords, color, popupHTML) {
    const existing = features[layer].get(id);
    const latlngs = coords.map(c => [c[1], c[0]]); // coords are [lng,lat]
    if (existing) {
      existing.setLatLngs(latlngs);
      existing.setStyle({ color });
      if (popupHTML) existing.bindPopup(popupHTML);
      return;
    }
    const line = L.polyline(latlngs, { color, weight: 4, opacity: 0.8 });
    if (popupHTML) line.bindPopup(popupHTML);
    line.addTo(groups[layer]);
    features[layer].set(id, line);
    counts[layer]++;
    document.getElementById('count-' + layer).textContent = counts[layer];
  }

  function handle(layer, payload) {
    switch (layer) {
      case 'flight': {
        const id = payload.icao24;
        if (id == null || payload.lat == null || payload.lng == null) return;
        upsertPoint('flight', id, payload.lat, payload.lng, popupTable([
          ['Callsign', payload.callsign || id],
          ['ICAO24', id],
          ['Altitude (m)', payload.altitude],
          ['Velocity (m/s)', payload.velocity],
          ['Heading', payload.heading],
          ['On ground', payload.on_ground],
        ]));
        return;
      }
      case 'ship': {
        const id = payload.mmsi;
        if (id == null || payload.lat == null || payload.lng == null) return;
        upsertPoint('ship', id, payload.lat, payload.lng, popupTable([
          ['MMSI', id],
          ['SOG (kn)', payload.sog],
          ['COG', payload.cog],
          ['Heading', payload.heading],
          ['Nav status', payload.nav_stat],
        ]));
        return;
      }
      case 'transport': {
        const id = payload.vehicle_id;
        if (id == null || payload.lat == null || payload.lng == null) return;
        upsertPoint('transport', id, payload.lat, payload.lng, popupTable([
          ['Vehicle', payload.label || id],
          ['Status', payload.status],
          ['Speed (m/s)', payload.speed],
          ['Bearing', payload.bearing],
        ]));
        return;
      }
      case 'road': {
        const id = payload.road_id;
        const coords = payload.coords;
        if (!id || !Array.isArray(coords) || coords.length < 2) return;
        const sev = (payload.severity || 'medium').toLowerCase();
        const color = severityColor[sev] || severityColor.medium;
        const reasons = Array.isArray(payload.reason) ? payload.reason.join(', ') : '';
        upsertPolyline('road', id, coords, color, popupTable([
          ['Road', payload.road_name || id],
          ['Restriction', payload.restriction_type],
          ['Severity', sev],
          ['Confidence', payload.confidence != null ? payload.confidence.toFixed(2) : ''],
          ['Reasons', reasons],
          ['KM', `${payload.start_km}–${payload.end_km}`],
        ]));
        return;
      }
      case 'obu': {
        const id = payload.session_id || payload.device_id;
        const lat = payload.lat ?? payload.latitude;
        const lng = payload.long ?? payload.lng ?? payload.longitude;
        if (id == null || lat == null || lng == null) return;
        upsertPoint('obu', id, lat, lng, popupTable([
          ['Session', id],
          ['Updated', payload.created_at || payload.timestamp],
        ]));
        return;
      }
    }
  }

  // Layer toggle wiring
  document.querySelectorAll('#panel input[type=checkbox]').forEach(cb => {
    cb.addEventListener('change', () => {
      const layer = cb.dataset.layer;
      if (cb.checked) map.addLayer(groups[layer]); else map.removeLayer(groups[layer]);
    });
  });

  // WebSocket with exponential backoff
  const status = document.getElementById('status');
  let backoff = 500;
  function connect() {
    const wsProto = location.protocol === 'https:' ? 'wss:' : 'ws:';
    const ws = new WebSocket(`${wsProto}//${location.host}/ws`);
    ws.onopen = () => { status.textContent = 'live'; backoff = 500; };
    ws.onclose = () => {
      status.textContent = `reconnecting in ${backoff}ms…`;
      setTimeout(connect, backoff);
      backoff = Math.min(backoff * 2, 10000);
    };
    ws.onerror = () => { ws.close(); };
    ws.onmessage = ev => {
      try {
        const env = JSON.parse(ev.data);
        if (!env || !env.layer || !env.payload) return;
        const payload = typeof env.payload === 'string' ? JSON.parse(env.payload) : env.payload;
        handle(env.layer, payload);
      } catch (e) {
        console.warn('bad message', e);
      }
    };
  }
  connect();
})();

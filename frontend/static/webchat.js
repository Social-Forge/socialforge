/**
 * Social Forge — embeddable webchat widget.
 *
 * Self-contained, dependency-free floating bubble. Embed with:
 *   <script src="https://your-host/webchat.js"
 *           data-sf-channel="<channelId>" data-sf-host="https://your-host" defer></script>
 *
 * Talks to the public webchat API (session / messages / poll). Uses
 * application/x-www-form-urlencoded so requests stay CORS-"simple" (no preflight)
 * and polls every few seconds for new agent/bot replies.
 */
;(function () {
  var script = document.currentScript
  if (!script) return
  var CHANNEL = script.getAttribute('data-sf-channel')
  var HOST = (script.getAttribute('data-sf-host') || '').replace(/\/$/, '')
  if (!CHANNEL || !HOST) return

  var STORAGE_KEY = 'sf_webchat_' + CHANNEL
  var POLL_MS = 3000
  var visitorId = localStorage.getItem(STORAGE_KEY) || null
  var cursor = null
  var seen = {}
  var open = false
  var pollTimer = null
  var started = false

  // --- styles --------------------------------------------------------------
  var css =
    '.sfw-bubble{position:fixed;bottom:20px;right:20px;width:56px;height:56px;border-radius:50%;' +
    'background:#4f46e5;color:#fff;border:none;cursor:pointer;box-shadow:0 6px 20px rgba(0,0,0,.25);' +
    'display:flex;align-items:center;justify-content:center;z-index:2147483000;font-size:24px}' +
    '.sfw-panel{position:fixed;bottom:88px;right:20px;width:360px;max-width:calc(100vw - 32px);' +
    'height:520px;max-height:calc(100vh - 120px);background:#fff;border-radius:14px;overflow:hidden;' +
    'box-shadow:0 12px 40px rgba(0,0,0,.28);display:none;flex-direction:column;z-index:2147483000;' +
    'font-family:system-ui,-apple-system,Segoe UI,Roboto,sans-serif}' +
    '.sfw-panel.sfw-open{display:flex}' +
    '.sfw-head{background:#4f46e5;color:#fff;padding:14px 16px;font-weight:600;font-size:15px}' +
    '.sfw-body{flex:1;overflow-y:auto;padding:14px;background:#f6f7f9;display:flex;flex-direction:column;gap:8px}' +
    '.sfw-msg{max-width:78%;padding:8px 12px;border-radius:14px;font-size:14px;line-height:1.4;white-space:pre-wrap;word-wrap:break-word}' +
    '.sfw-visitor{align-self:flex-end;background:#4f46e5;color:#fff;border-bottom-right-radius:4px}' +
    '.sfw-bot,.sfw-agent{align-self:flex-start;background:#fff;color:#111;border:1px solid #e5e7eb;border-bottom-left-radius:4px}' +
    '.sfw-foot{display:flex;gap:8px;padding:10px;border-top:1px solid #eee;background:#fff}' +
    '.sfw-input{flex:1;border:1px solid #d1d5db;border-radius:10px;padding:9px 12px;font-size:14px;outline:none}' +
    '.sfw-send{background:#4f46e5;color:#fff;border:none;border-radius:10px;padding:0 16px;cursor:pointer;font-size:14px}' +
    '.sfw-send:disabled{opacity:.5;cursor:default}'
  var style = document.createElement('style')
  style.textContent = css
  document.head.appendChild(style)

  // --- dom ------------------------------------------------------------------
  var bubble = el('button', 'sfw-bubble', '💬')
  var panel = el('div', 'sfw-panel')
  var head = el('div', 'sfw-head', 'Chat with us')
  var body = el('div', 'sfw-body')
  var foot = el('div', 'sfw-foot')
  var input = el('input', 'sfw-input')
  input.type = 'text'
  input.placeholder = 'Type a message…'
  var sendBtn = el('button', 'sfw-send', 'Send')
  foot.appendChild(input)
  foot.appendChild(sendBtn)
  panel.appendChild(head)
  panel.appendChild(body)
  panel.appendChild(foot)
  document.body.appendChild(bubble)
  document.body.appendChild(panel)

  bubble.addEventListener('click', toggle)
  sendBtn.addEventListener('click', send)
  input.addEventListener('keydown', function (e) {
    if (e.key === 'Enter') send()
  })

  // --- behavior -------------------------------------------------------------
  function toggle() {
    open = !open
    panel.classList.toggle('sfw-open', open)
    bubble.textContent = open ? '✕' : '💬'
    if (open && !started) start()
    if (open) startPolling()
    else stopPolling()
  }

  function start() {
    started = true
    post('/webchat/' + CHANNEL + '/session', { visitorId: visitorId || '' })
      .then(function (data) {
        if (!data) return
        visitorId = data.visitorId
        localStorage.setItem(STORAGE_KEY, visitorId)
        if (data.channel && data.channel.name) head.textContent = data.channel.name
        render(data.messages || [])
      })
      .catch(noop)
  }

  function send() {
    var text = input.value.trim()
    if (!text || !visitorId) return
    input.value = ''
    // Optimistic render.
    append({ id: 'local-' + Date.now(), role: 'visitor', body: text })
    post('/webchat/' + CHANNEL + '/messages', { visitorId: visitorId, body: text }).catch(noop)
  }

  function startPolling() {
    stopPolling()
    pollTimer = setInterval(poll, POLL_MS)
  }
  function stopPolling() {
    if (pollTimer) clearInterval(pollTimer)
    pollTimer = null
  }

  function poll() {
    if (!visitorId) return
    var url = '/webchat/' + CHANNEL + '/messages?visitorId=' + encodeURIComponent(visitorId)
    if (cursor) url += '&after=' + encodeURIComponent(cursor)
    fetch(HOST + url)
      .then(function (r) {
        return r.ok ? r.json() : null
      })
      .then(function (data) {
        if (data && data.messages) render(data.messages)
      })
      .catch(noop)
  }

  function render(messages) {
    messages.forEach(function (m) {
      append(m)
    })
  }

  function append(m) {
    // Dedup server echoes against optimistic + prior renders.
    if (m.id && seen[m.id]) return
    if (m.id) seen[m.id] = true
    if (m.at) cursor = m.at
    var node = el('div', 'sfw-msg sfw-' + (m.role || 'agent'), m.body || '')
    body.appendChild(node)
    body.scrollTop = body.scrollHeight
  }

  // --- helpers --------------------------------------------------------------
  function post(path, params) {
    return fetch(HOST + path, {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: encode(params),
    }).then(function (r) {
      return r.ok
        ? r.json().catch(function () {
            return {}
          })
        : null
    })
  }

  function encode(obj) {
    return Object.keys(obj)
      .map(function (k) {
        return encodeURIComponent(k) + '=' + encodeURIComponent(obj[k])
      })
      .join('&')
  }

  function el(tag, className, text) {
    var node = document.createElement(tag)
    if (className) node.className = className
    if (text != null) node.textContent = text
    return node
  }

  function noop() {}
})()

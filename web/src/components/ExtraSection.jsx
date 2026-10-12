import React, { useState } from 'react'

export default function ExtraSection({ items, onAdd, onDel }) {
  const blank = { name: '', time: 0, max_time: 0, start_hour: 0, end_hour: 24 }
  const [f, setF] = useState(blank)
  const [msg, setMsg] = useState('')
  const set = (k, v) => setF((p) => ({ ...p, [k]: v }))

  const submit = async (e) => {
    e.preventDefault()
    try {
      await onAdd({ ...f, time: Number(f.time) || 0, max_time: Number(f.max_time) || 0, start_hour: Number(f.start_hour) || 0, end_hour: Number(f.end_hour) || 24 })
      setF(blank); setMsg('')
    } catch (err) { setMsg(err.message) }
  }

  return (
    <div className="panel" style={{ borderLeft: '3px solid var(--extra)' }}>
      <h2>Optional activities</h2>
      <form className="row" onSubmit={submit}>
        <input type="text" placeholder="Name" value={f.name} onChange={(e) => set('name', e.target.value)} required />
        <label>Hrs</label><input type="number" min="0" step="0.5" value={f.time} onChange={(e) => set('time', e.target.value)} />
        <label>Max</label><input type="number" min="0" step="0.5" value={f.max_time} onChange={(e) => set('max_time', e.target.value)} />
        <label>Window</label>
        <input type="number" min="0" max="24" step="0.5" value={f.start_hour} onChange={(e) => set('start_hour', e.target.value)} />
        <span>–</span>
        <input type="number" min="0" max="24" step="0.5" value={f.end_hour} onChange={(e) => set('end_hour', e.target.value)} />
        <button className="btn" type="submit">Add</button>
      </form>
      {msg && <div className="banner warn">{msg}</div>}
      <div className="list">
        {items.length === 0 && <div className="item"><span className="meta">None</span></div>}
        {items.map((it) => (
          <div className="item" key={it.id}>
            <div>
              <span className="name">{it.name}</span>
              <span className="meta">  {it.time}h {it.start_hour || it.end_hour ? `${it.start_hour}–${it.end_hour}` : ''}</span>
            </div>
            <button className="del" onClick={() => onDel(it.id)}>✕</button>
          </div>
        ))}
      </div>
    </div>
  )
}
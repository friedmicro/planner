import React, { useState } from 'react'
import { DAYS } from '../api.js'

export default function EventSection({ items, onAdd, onDel, onUpdate }) {
  const blank = { name: '', time: 0, min_block: 0, max_block: 0, days: [], fixed_day: 0, fixed_start: 0, fixed_end: 0, hour_earliest: 0, hour_latest: 0, hour_start: 0, hour_end: 0 }
  const [f, setF] = useState(blank)
  const [msg, setMsg] = useState('')
  const [editId, setEditId] = useState(null)

  const set = (k, v) => setF((p) => ({ ...p, [k]: v }))

  const submit = async (e) => {
    e.preventDefault()
    const payload = {
      name: f.name,
      time: Number(f.time) || 0,
      min_block: Number(f.min_block) || 0,
      max_block: Number(f.max_block) || 0,
      days: f.fixed_day ? [] : (f.days || []),
      fixed_day: Number(f.fixed_day) || 0,
      fixed_start: Number(f.fixed_start) || 0,
      fixed_end: Number(f.fixed_end) || 0,
      hour_earliest: Number(f.hour_earliest) || 0,
      hour_latest: Number(f.hour_latest) || 0,
      hour_start: Number(f.hour_start) || 0,
      hour_end: Number(f.hour_end) || 0,
    }
    try {
      if (editId != null) { await onUpdate({ ...payload, id: editId }); setEditId(null) }
      else await onAdd(payload)
      setF(blank); setMsg('')
    } catch (err) { setMsg(err.message) }
  }

  const toggleDay = (d) => set('days', (f.days || []).includes(d) ? f.days.filter((x) => x !== d) : [...(f.days || []), d])

  return (
    <div className="panel" style={{ borderLeft: '3px solid var(--event)' }}>
      <h2>Commitments (required events)</h2>
      <form onSubmit={submit}>
        <div className="row"><label>Name</label><input type="text" value={f.name} onChange={(e) => set('name', e.target.value)} required /></div>
        <div className="row">
          <label>Hours/wk</label><input type="number" min="0" step="0.5" value={f.time} onChange={(e) => set('time', e.target.value)} />
          <label>Min block</label><input type="number" min="0" step="0.5" value={f.min_block} onChange={(e) => set('min_block', e.target.value)} />
          <label>Max block</label><input type="number" min="0" step="0.5" value={f.max_block} onChange={(e) => set('max_block', e.target.value)} />
        </div>
        <div className="row">
          <label>Fixed slot?</label>
          <select value={f.fixed_day} onChange={(e) => set('fixed_day', e.target.value)}>
            <option value={0}>No</option>
            {DAYS.map((d, i) => <option key={i + 1} value={i + 1}>{d}</option>)}
          </select>
          {Number(f.fixed_day) > 0 && <>
            <label>Start</label><input type="number" min="0" max="24" step="0.5" value={f.fixed_start} onChange={(e) => set('fixed_start', e.target.value)} />
            <label>End</label><input type="number" min="0" max="24" step="0.5" value={f.fixed_end} onChange={(e) => set('fixed_end', e.target.value)} />
          </>}
        </div>
        <div className="row">
          <label>Window hrs</label><input type="number" min="0" max="24" step="0.5" value={f.hour_earliest} onChange={(e) => set('hour_earliest', e.target.value)} placeholder="earliest" />
          <input type="number" min="0" max="24" step="0.5" value={f.hour_latest} onChange={(e) => set('hour_latest', e.target.value)} placeholder="latest" />
        </div>
        {Number(f.fixed_day) === 0 && (
          <div className="row">
            <label>Allowed days</label>
            {DAYS.map((d, i) => (
              <label key={i} style={{ minWidth: 0 }}>
                <input type="checkbox" checked={(f.days || []).includes(i + 1)} onChange={() => toggleDay(i + 1)} /> {d}
              </label>
            ))}
          </div>
        )}
        <div className="row">
          <button className="btn" type="submit">{editId != null ? 'Update' : 'Add commitment'}</button>
          {editId != null && <button className="btn ghost" type="button" onClick={() => { setF(blank); setEditId(null) }}>Cancel</button>}
        </div>
      </form>
      {msg && <div className="banner warn">{msg}</div>}
      <div className="list">
        {items.length === 0 && <div className="item"><span className="meta">None</span></div>}
        {items.map((it) => (
          <div className="item" key={it.id}>
            <div>
              <span className="name">{it.name}</span>
              <span className="meta">  {it.time}h{it.days?.length ? ' on ' + it.days.map((d) => DAYS[d - 1]).join(',') : ''}{it.fixed_day ? ` fixed ${DAYS[it.fixed_day - 1]} ${it.fixed_start}-${it.fixed_end}` : ''}</span>
            </div>
            <span>
              <button className="del" onClick={() => { setF({ ...blank, ...it, days: it.days || [] }); setEditId(it.id) }}>✎</button>
              <button className="del" onClick={() => onDel(it.id)}>✕</button>
            </span>
          </div>
        ))}
      </div>
    </div>
  )
}
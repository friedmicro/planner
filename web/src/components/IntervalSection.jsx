import React, { useState } from 'react'
import { DAYS } from '../api.js'

export default function IntervalSection({ title, items, onAdd, onDel, kindColor }) {
  const [day, setDay] = useState(1)
  const [start, setStart] = useState(22)
  const [end, setEnd] = useState(8)
  const [msg, setMsg] = useState('')

  const submit = async (e) => {
    e.preventDefault()
    const s = Number(start)
    const en = Number(end)
    let normS = s, normE = en
    if (en < s) normE = en + 24
    if (s > 12 && en <= 12 && normE - normS > 20) normE = en // guard: avoid weird wrap
    try {
      await onAdd({ day: Number(day), start: normS, end: en < s && en <= 12 ? en : normE })
      setMsg('')
    } catch (err) {
      setMsg(err.message)
    }
  }

  return (
    <div className="panel" style={{ borderLeft: `3px solid ${kindColor}` }}>
      <h2>{title}</h2>
      <form className="row" onSubmit={submit}>
        <label>Day</label>
        <select value={day} onChange={(e) => setDay(e.target.value)}>
          {DAYS.map((d, i) => <option key={i + 1} value={i + 1}>{d}</option>)}
        </select>
        <label>Start</label>
        <input type="number" min="0" max="24" step="0.5" value={start} onChange={(e) => setStart(e.target.value)} />
        <label>End</label>
        <input type="number" min="0" max="24" step="0.5" value={end} onChange={(e) => setEnd(e.target.value)} />
        <button className="btn" type="submit">Add</button>
      </form>
      {msg && <div className="banner warn">{msg}</div>}
      <div className="list">
        {items.length === 0 && <div className="item"><span className="meta">None</span></div>}
        {items.map((it) => (
          <div className="item" key={it.id}>
            <div>
              <span className="name">{DAYS[it.day - 1]}</span>
              <span className="meta">  {it.start}h–{it.end}h</span>
            </div>
            <button className="del" onClick={() => onDel(it.id)}>✕</button>
          </div>
        ))}
      </div>
    </div>
  )
}
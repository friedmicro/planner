import React from 'react'
import { DAYS } from '../api.js'

const KIND_COLOR = {
  event: 'var(--event)',
  extra: 'var(--extra)',
  sleep: 'var(--sleep)',
  nap: 'var(--nap)',
  blackout: 'var(--blackout)',
}

const H = 480
const px = (hour) => (hour / 24) * H

function HourGrid() {
  const marks = []
  for (let h = 0; h <= 24; h += 4) {
    marks.push(
      <div key={h}>
        <div className="gline" style={{ top: px(h) }} />
        <span className="glabel" style={{ top: px(h) - 2 }}>{String(h % 24).padStart(2, '0')}</span>
      </div>,
    )
  }
  return marks
}

function fmt(h) {
  const hh = Math.floor(h)
  const mm = Math.round((h - hh) * 60)
  return `${String(hh % 24).padStart(2, '0')}:${String(mm).padStart(2, '0')}`
}

export default function Timeline({ schedule }) {
  if (!schedule) {
    return <div className="banner info">No schedule yet. Press "Solve week".</div>
  }
  const placements = schedule.placements || []
  const byDay = {}
  for (let d = 1; d <= 7; d++) byDay[d] = []
  for (const p of placements) {
    if (p.day >= 1 && p.day <= 7) byDay[p.day].push(p)
  }

  return (
    <div className="timeline">
      {DAYS.map((name, i) => {
        const day = i + 1
        const blocks = byDay[day]
        return (
          <div className="day" key={day}>
            <h3>{name}</h3>
            <div className="track">
              <HourGrid />
              {blocks.map((b, idx) => {
                const top = px(b.start)
                const hgt = Math.max(px(b.end) - px(b.start), 8)
                return (
                  <div
                    key={idx}
                    className="block"
                    style={{ top, height: hgt, background: KIND_COLOR[b.kind] || 'var(--event)' }}
                    title={`${b.name} · ${fmt(b.start)}–${fmt(b.end)}`}
                  >
                    <span className="t">{b.name}</span>
                    {hgt > 22 && <div>{fmt(b.start)}–{fmt(b.end)}</div>}
                  </div>
                )
              })}
            </div>
          </div>
        )
      })}
    </div>
  )
}
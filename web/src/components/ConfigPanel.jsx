import React from 'react'

export default function ConfigPanel({ config, setConfig, banner }) {
  const cfg = (k, step = 0.5) => (
    <input type="number" step={step} value={config[k] ?? ''}
      onChange={(e) => setConfig((p) => ({ ...p, [k]: e.target.value }))} />
  )

  return (
    <>
      <div className="row">
        <label>Sleep workday</label>{cfg('sleep_per_day_work_week')}
        <label>Sleep weekend</label>{cfg('sleep_per_day_weekend')}
        <label>Min block</label>{cfg('default_min_block')}
      </div>
      <div className="row">
        <label>Rests</label>{cfg('preferred_rests', 1)}
        <label className="meta" style={{ minWidth: 0 }}>
          <input type="checkbox" checked={config.prefer_morning_nap}
            onChange={(e) => setConfig((p) => ({ ...p, prefer_morning_nap: e.target.checked }))} /> morning nap
        </label>
        <label className="meta" style={{ minWidth: 0 }}>
          <input type="checkbox" checked={config.allow_short_rests}
            onChange={(e) => setConfig((p) => ({ ...p, allow_short_rests: e.target.checked }))} /> short rests
        </label>
        <label className="meta" style={{ minWidth: 0 }}>
          <input type="checkbox" checked={config.schedule_extra}
            onChange={(e) => setConfig((p) => ({ ...p, schedule_extra: e.target.checked }))} /> schedule extras
        </label>
        <label className="meta" style={{ minWidth: 0 }}>
          <input type="checkbox" checked={config.prefer_spacing_required_work}
            onChange={(e) => setConfig((p) => ({ ...p, prefer_spacing_required_work: e.target.checked }))} /> space required
        </label>
      </div>
      {banner && <div className={`banner ${banner.type}`}>{banner.text}</div>}
    </>
  )
}

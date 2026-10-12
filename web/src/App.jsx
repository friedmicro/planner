import React, { useEffect, useState } from 'react'
import { api } from './api.js'
import Timeline from './components/Timeline.jsx'
import ConfigPanel from './components/ConfigPanel.jsx'
import IntervalSection from './components/IntervalSection.jsx'
import EventSection from './components/EventSection.jsx'
import ExtraSection from './components/ExtraSection.jsx'

const DEFAULT_CONFIG = {
  sleep_per_day_work_week: 7,
  sleep_per_day_weekend: 8,
  prefer_morning_nap: false,
  preferred_rests: 2,
  allow_short_rests: false,
  default_min_block: 1,
  schedule_extra: false,
  prefer_spacing_required_work: false,
}

export default function App() {
  const [config, setConfig] = useState(DEFAULT_CONFIG)
  const [sleep, setSleep] = useState([])
  const [blacklist, setBlacklist] = useState([])
  const [events, setEvents] = useState([])
  const [extras, setExtras] = useState([])
  const [schedule, setSchedule] = useState(null)
  const [loading, setLoading] = useState(true)
  const [solving, setSolving] = useState(false)
  const [banner, setBanner] = useState(null)

  const load = async () => {
    try {
      const [c, s, b, e, x] = await Promise.all([
        api.config.get(), api.sleep.list(), api.blacklist.list(),
        api.events.list(), api.extras.list(),
      ])
      setConfig(c || DEFAULT_CONFIG)
      setSleep(s || []); setBlacklist(b || []); setEvents(e || []); setExtras(x || [])
      const sch = await api.schedule().catch(() => null)
      setSchedule(sch)
    } catch (err) {
      setBanner({ type: 'warn', text: `Failed to load: ${err.message}` })
    } finally {
      setLoading(false)
    }
  }

  const flash = (text, type = 'info') => setBanner({ type, text })

  const saveConfig = async () => {
    try {
      const saved = await api.config.save({
        ...config,
        sleep_per_day_work_week: Number(config.sleep_per_day_work_week),
        sleep_per_day_weekend: Number(config.sleep_per_day_weekend),
        preferred_rests: Number(config.preferred_rests),
        default_min_block: Number(config.default_min_block),
      })
      setConfig(saved)
      flash('Config saved')
    } catch (err) { flash(err.message, 'warn') }
  }

  const solve = async () => {
    setSolving(true); setBanner(null)
    try {
      const sch = await api.solve()
      setSchedule(sch)
      flash('Week solved')
    } catch (err) { flash(err.message, 'warn') }
    finally { setSolving(false) }
  }

  const interval = (setter, bucket) => async (item) => {
    const created = await bucket.add(item)
    setter((prev) => [...prev, { ...item, id: created.id }])
    flash('Added')
  }

  const delInterval = (setter, bucket) => async (id) => {
    await bucket.del(id)
    setter((prev) => prev.filter((x) => x.id !== id))
    flash('Removed')
  }

  const addEvent = async (e) => {
    const created = await api.events.add(e)
    setEvents((prev) => [...prev, { ...e, id: created.id }]); flash('Commitment added')
  }

  const updateEvent = async (e) => {
    const updated = await api.events.update(e)
    setEvents((prev) => prev.map((x) => (x.id === e.id ? updated : x)))
    flash('Commitment updated')
  }

  const delEvent = async (id) => {
    await api.events.del(id)
    setEvents((prev) => prev.filter((x) => x.id !== id)); flash('Commitment removed')
  }

  const addExtra = async (x) => {
    const created = await api.extras.add(x)
    setExtras((prev) => [...prev, { ...x, id: created.id }]); flash('Activity added')
  }

  const delExtra = async (id) => {
    await api.extras.del(id)
    setExtras((prev) => prev.filter((x) => x.id !== id)); flash('Activity removed')
  }

  useEffect(() => { load() }, [])

  if (loading) return <div className="banner info">Loading…</div>

  return (
    <>
      <header>
        <div>
          <h1>planner</h1>
          <div className="sub">Weekly schedule solver</div>
        </div>
        <div className="toolbar">
          <button className="btn ghost" onClick={saveConfig}>Save config</button>
          <button className="btn" onClick={solve} disabled={solving}>{solving ? 'Solving…' : 'Solve week'}</button>
        </div>
      </header>

      <main>
        <section className="panel schedule">
          <h2>Week</h2>
          <ConfigPanel config={config} setConfig={setConfig} banner={banner} />
          <Timeline schedule={schedule} />
        </section>

        <section className="editor">
          <IntervalSection title="Sleep" items={sleep} kindColor="var(--sleep)"
            onAdd={interval(setSleep, api.sleep)} onDel={delInterval(setSleep, api.sleep)} />
          <IntervalSection title="Blackout (unavailable)" items={blacklist} kindColor="var(--blackout)"
            onAdd={interval(setBlacklist, api.blacklist)} onDel={delInterval(setBlacklist, api.blacklist)} />
          <EventSection items={events} onAdd={addEvent} onUpdate={updateEvent} onDel={delEvent} />
          <ExtraSection items={extras} onAdd={addExtra} onDel={delExtra} />
        </section>
      </main>
    </>
  )
}

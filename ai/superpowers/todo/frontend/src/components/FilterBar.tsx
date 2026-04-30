import type { Order, Query, SortKey, Status } from '../types'

interface Props {
  query: Query
  counts: { all: number; active: number; completed: number }
  onChange: (next: Query) => void
}

const segments: Array<{ value: Status; label: string }> = [
  { value: 'all', label: '전체' },
  { value: 'active', label: '미완료' },
  { value: 'completed', label: '완료' },
]

export function FilterBar({ query, counts, onChange }: Props) {
  return (
    <div className="filter-bar" role="toolbar" aria-label="필터/정렬">
      <fieldset className="filter-bar__segments">
        <legend className="visually-hidden">상태</legend>
        {segments.map((s) => (
          <label key={s.value} className="filter-bar__segment">
            <input
              type="radio"
              name="status"
              value={s.value}
              checked={query.status === s.value}
              onChange={() => onChange({ ...query, status: s.value })}
            />
            <span className="filter-bar__segment-label">
              {s.label}
              <span className="filter-bar__count" aria-hidden="true">{counts[s.value]}</span>
            </span>
          </label>
        ))}
      </fieldset>

      <div className="filter-bar__sort">
        <select
          className="filter-bar__sort-select"
          aria-label="정렬"
          value={query.sort}
          onChange={(e) => onChange({ ...query, sort: e.target.value as SortKey })}
        >
          <option value="createdAt">생성일</option>
          <option value="dueDate">마감일</option>
          <option value="priority">우선순위</option>
        </select>
        <button
          type="button"
          className="filter-bar__order"
          aria-label="정렬 방향 토글"
          onClick={() =>
            onChange({ ...query, order: (query.order === 'asc' ? 'desc' : 'asc') as Order })
          }
        >
          {query.order === 'asc' ? '↑' : '↓'}
        </button>
      </div>
    </div>
  )
}

import type { Order, Query, SortKey, Status } from '../types'

interface Props {
  query: Query
  onChange: (next: Query) => void
}

export function FilterBar({ query, onChange }: Props) {
  return (
    <div role="toolbar" aria-label="필터/정렬">
      <fieldset>
        <legend>상태</legend>
        {(['all', 'active', 'completed'] as Status[]).map((s) => (
          <label key={s}>
            <input
              type="radio"
              name="status"
              value={s}
              checked={query.status === s}
              onChange={() => onChange({ ...query, status: s })}
            />
            {s === 'all' ? '전체' : s === 'active' ? '미완료' : '완료'}
          </label>
        ))}
      </fieldset>
      <label>
        정렬
        <select
          value={query.sort}
          onChange={(e) => onChange({ ...query, sort: e.target.value as SortKey })}
        >
          <option value="createdAt">생성일</option>
          <option value="dueDate">마감일</option>
          <option value="priority">우선순위</option>
        </select>
      </label>
      <button
        type="button"
        aria-label="정렬 방향 토글"
        onClick={() => onChange({ ...query, order: (query.order === 'asc' ? 'desc' : 'asc') as Order })}
      >
        {query.order === 'asc' ? '↑' : '↓'}
      </button>
    </div>
  )
}

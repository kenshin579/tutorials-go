# Wails Todo

Wails v2 + React + TypeScript로 구현한 데스크톱 Todo 애플리케이션입니다.

## 기술 스택

- **Backend**: Go 1.23, Wails v2.11
- **Frontend**: React 18, TypeScript, Vite
- **데이터 저장**: JSON 파일 기반 로컬 스토리지 (`todos.json`)

## 주요 기능

- Todo 추가, 완료 토글, 삭제 (삭제 시 확인 다이얼로그)
- 완료 현황 표시 (N개 중 M개 완료)
- JSON 파일로 Todo 목록 내보내기/불러오기
- 네이티브 메뉴 지원 (파일 > 불러오기/내보내기/종료)
- 다크 테마 UI

## 프로젝트 구조

```
wails-todo/
├── main.go                  # Wails 앱 엔트리포인트
├── backend/
│   ├── app.go               # App 구조체 및 바인딩 메서드 (CRUD, Import/Export)
│   ├── todo.go              # Todo 모델 및 TodoStore (JSON 파일 읽기/쓰기)
│   └── menu.go              # 네이티브 메뉴 설정
├── frontend/
│   ├── src/
│   │   ├── App.tsx           # 메인 컴포넌트 (상태 관리, 이벤트 리스너)
│   │   ├── App.css           # 스타일 (다크 테마)
│   │   └── components/
│   │       ├── TodoInput.tsx  # 입력 폼
│   │       ├── TodoItem.tsx   # 개별 Todo 항목
│   │       └── TodoList.tsx   # Todo 목록
│   └── wailsjs/              # Wails 자동 생성 바인딩
├── build/                    # 빌드 설정 (아이콘, 플랫폼별 설정)
├── wails.json                # Wails 프로젝트 설정
└── go.mod
```

## 사전 요구사항

- Go 1.23+
- Node.js 16+
- [Wails CLI v2](https://wails.io/docs/gettingstarted/installation)
  ```bash
  go install github.com/wailsapp/wails/v2/cmd/wails@latest
  ```

## 실행 방법

### 개발 모드

```bash
cd desktop/wails-todo
wails dev
```

핫 리로드가 지원되며, `http://localhost:34115`에서 브라우저로도 접근 가능합니다.

### 프로덕션 빌드

```bash
wails build
```

빌드 결과물은 `build/bin/` 디렉터리에 생성됩니다.

## Backend API (바인딩 메서드)

| 메서드 | 설명 |
|--------|------|
| `GetTodos()` | 전체 Todo 목록 조회 |
| `AddTodo(title)` | 새 Todo 추가 |
| `ToggleTodo(id)` | 완료 상태 토글 |
| `DeleteTodo(id)` | Todo 삭제 (확인 다이얼로그 포함) |
| `ExportTodos()` | 파일 저장 다이얼로그를 통해 JSON 내보내기 |
| `ImportTodos()` | 파일 열기 다이얼로그를 통해 JSON 불러오기 |

## 단축키

| 단축키 | 기능 |
|--------|------|
| `Cmd/Ctrl + O` | Todo 목록 불러오기 |
| `Cmd/Ctrl + S` | Todo 목록 내보내기 |
| `Cmd/Ctrl + Q` | 앱 종료 |

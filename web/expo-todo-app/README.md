# expo-todo-app

[Expo](https://expo.dev/)(blank TypeScript 템플릿)로 만든 간단한 Todo 앱 예제입니다.
React Native 입문 + Expo 개발 흐름을 익히기 위한 실습용 코드입니다.

> 📝 블로그 글: [Expo로 시작하는 React Native 앱 개발: Todo 앱 만들기](https://blog-v2.advenoh.pe.kr)

## 기능

- ✅ 할 일 추가 (`TextInput` + 추가 버튼)
- ✅ 완료 토글 (항목을 탭하면 완료/취소, 취소선 표시)
- ✅ 삭제
- ✅ **AsyncStorage**로 데이터 영속화 — 앱을 껐다 켜도 목록이 유지됨

## 기술 스택

| 항목 | 버전 |
|------|------|
| Expo SDK | ~56.0 |
| React | 19.2 |
| React Native | 0.85 |
| TypeScript | ~6.0 |
| @react-native-async-storage/async-storage | 2.2 |

> 버전은 `package.json` 기준입니다. 새로 시작한다면 최신 SDK로 만들어도 코드는 거의 그대로 동작합니다.

## 사전 준비

- [Node.js](https://nodejs.org/) LTS 이상 (`node -v`로 확인)
- 실행 방법에 따라 아래 중 하나
  - **실기기**: [Expo Go](https://expo.dev/go) 앱 (가장 간편, 추가 설치 불필요)
  - **iOS 시뮬레이터**: macOS + Xcode
  - **Android 에뮬레이터**: Android Studio

## 시작하기

```bash
# 1. 의존성 설치
npm install

# 2. 개발 서버 실행
npx expo start
```

`expo start`를 실행하면 터미널에 QR 코드와 실행 옵션 메뉴가 나타납니다.

| 키 | 동작 | 필요 환경 |
|----|------|-----------|
| `i` | iOS 시뮬레이터에서 실행 | macOS + Xcode |
| `a` | Android 에뮬레이터에서 실행 | Android Studio |
| `w` | 웹 브라우저에서 실행 | - |
| QR 스캔 | 실기기에서 실행 | Expo Go 앱 |

## iOS 시뮬레이터가 없다면?

iOS 시뮬레이터는 **macOS + Xcode**가 있어야만 쓸 수 있습니다. Xcode가 없거나(또는 Windows/Linux 사용 중이거나) 설치가 번거롭다면, 시뮬레이터 없이도 충분히 실행할 수 있습니다.

### 방법 1. 실기기 + Expo Go (가장 권장)

시뮬레이터/에뮬레이터 없이 **내 폰에서 바로** 실행하는 가장 쉬운 방법입니다.

1. 스마트폰에 [Expo Go](https://expo.dev/go) 앱 설치 (iOS 앱스토어 / Android 플레이스토어)
2. `npx expo start` 실행
3. 터미널에 표시된 QR 코드를 스캔
   - **iOS**: 기본 카메라 앱으로 QR을 찍으면 Expo Go로 연결
   - **Android**: Expo Go 앱 안의 스캐너로 QR 스캔
4. 단, **PC와 스마트폰이 같은 Wi-Fi 네트워크**에 있어야 합니다. 회사망 등으로 연결이 안 되면 터널 모드로 우회하세요.

```bash
npx expo start --tunnel
```

### 방법 2. 웹 브라우저로 실행

가장 빠르게 UI를 확인하는 방법입니다. 아무 것도 설치할 필요가 없습니다.

```bash
npx expo start --web
# 또는 expo start 실행 후 터미널에서 w 키
```

이 Todo 앱은 웹에서도 동작하며, AsyncStorage가 웹에서는 브라우저 `localStorage`로 처리되어 **새로고침 후에도 목록이 유지**됩니다. 다만 웹은 React Native Web으로 렌더링되므로 실제 모바일과 미세한 차이가 있을 수 있습니다.

### 방법 3. Android 에뮬레이터

[Android Studio](https://developer.android.com/studio)가 설치되어 있다면, AVD(Android Virtual Device)를 만들고 `a` 키로 실행할 수 있습니다. iOS 시뮬레이터 없이 가상 기기에서 테스트하고 싶을 때 좋습니다.

### 방법 4. iOS 시뮬레이터를 설치하고 싶다면 (macOS 한정)

1. Mac 앱스토어에서 **Xcode** 설치
2. Xcode를 한 번 실행해 라이선스 동의 및 추가 구성요소 설치
3. **Xcode → Settings → Components**(버전에 따라 Platforms)에서 iOS 시뮬레이터 런타임 다운로드
4. 커맨드라인 도구 경로 지정 (필요 시):

```bash
sudo xcode-select -s /Applications/Xcode.app/Contents/Developer
```

5. 이후 `npx expo start` → `i` 키로 시뮬레이터 실행

> 정리: **빠르게 확인만 하려면 방법 2(웹)**, **실제 모바일 감각이 필요하면 방법 1(실기기 + Expo Go)**을 추천합니다. iOS 시뮬레이터 설치는 선택 사항입니다.

## npm 스크립트

| 명령 | 설명 |
|------|------|
| `npm start` | `expo start` (개발 서버) |
| `npm run ios` | iOS 시뮬레이터로 바로 실행 |
| `npm run android` | Android 에뮬레이터로 바로 실행 |
| `npm run web` | 웹으로 바로 실행 |

## 프로젝트 구조

```
expo-todo-app/
├── App.tsx          # 앱 진입점이자 전체 로직 (이 파일 하나로 동작)
├── index.ts         # registerRootComponent(App) — 실제 엔트리
├── app.json         # 앱 이름·아이콘·스플래시 등 Expo 설정
├── package.json     # 의존성 및 스크립트
├── tsconfig.json    # expo/tsconfig.base 확장 (strict 모드)
└── assets/          # 아이콘, 스플래시 이미지
```

거의 모든 코드는 **`App.tsx`** 한 파일에 들어 있습니다.

## 코드 동작 방식

### 상태 관리

- `todos`: 할 일 배열 (`Todo[]`). 각 항목은 `{ id, text, done }`
- `text`: 입력창의 현재 값
- `loaded`: AsyncStorage 최초 로드 완료 여부 플래그

### 주요 컴포넌트

- `View` / `Text`: 웹의 `div` / `p`에 해당하는 기본 빌딩 블록
- `TextInput`: 텍스트 입력 (`onChangeText`로 값 수신)
- `FlatList`: 가상화 리스트 — 화면에 보이는 항목만 렌더링
- `TouchableOpacity`: 누르면 살짝 흐려지는 터치 영역(버튼 역할)
- `SafeAreaView`: 노치/상태 표시줄을 피해 안전 영역에 렌더링

### 데이터 영속화 (AsyncStorage)

두 개의 `useEffect`로 처리합니다.

1. **불러오기**: 앱이 마운트될 때(`[]` 의존성) 저장된 JSON을 읽어 `todos`에 복원하고, 끝나면 `loaded`를 `true`로 설정
2. **저장하기**: `todos`가 바뀔 때마다 JSON으로 직렬화해 저장

```ts
useEffect(() => {
  if (!loaded) return; // 최초 로드 전에는 저장하지 않음
  AsyncStorage.setItem(STORAGE_KEY, JSON.stringify(todos));
}, [todos, loaded]);
```

`loaded` 플래그가 핵심입니다. 이 가드가 없으면, 저장된 데이터를 불러오기 **전에** 저장 effect가 빈 배열(`[]`)로 먼저 실행되어 기존 데이터를 덮어쓰게 됩니다.

> 네이티브 모듈인 AsyncStorage는 `npm install`이 아니라 **`npx expo install`**로 설치했습니다. Expo가 현재 SDK 버전과 호환되는 버전을 자동으로 골라줍니다.

## 타입 체크

```bash
npx tsc --noEmit
```

에러 없이 종료되면 정상입니다.

## 라이선스

[MIT](./LICENSE)

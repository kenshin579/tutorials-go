#!/bin/bash
#
# Ralph Loop 실행 스크립트
#
# ✅ DO: 모드와 반복 횟수를 파라미터로 제어
# ❌ DON'T: 하드코딩된 무한 루프만 사용
#
# Usage:
#   ./loop.sh              # Build 모드, 무제한
#   ./loop.sh plan         # Plan 모드, 무제한
#   ./loop.sh plan 3       # Plan 모드, 최대 3회
#   ./loop.sh build 10     # Build 모드, 최대 10회
#

set -euo pipefail

# === 설정 ===
MODE="build"
PROMPT_FILE="PROMPT_build.md"
MAX_ITERATIONS=0
ITERATION=0
CURRENT_BRANCH=$(git branch --show-current)

# === 인자 파싱 ===
if [ "${1:-}" = "plan" ]; then
  MODE="plan"
  PROMPT_FILE="PROMPT_plan.md"
  MAX_ITERATIONS=${2:-0}
elif [[ "${1:-}" =~ ^[0-9]+$ ]]; then
  MAX_ITERATIONS=$1
elif [ "${1:-}" = "build" ]; then
  MAX_ITERATIONS=${2:-0}
fi

# === 사전 검증 ===
# ✅ DO: 프롬프트 파일 존재 여부 확인
# ❌ DON'T: 파일 없이 실행하여 에러 발생
if [ ! -f "$PROMPT_FILE" ]; then
  echo "Error: $PROMPT_FILE 파일이 없습니다."
  exit 1
fi

echo "=== Ralph Loop ==="
echo "Mode: $MODE"
echo "Prompt: $PROMPT_FILE"
echo "Max iterations: $([ $MAX_ITERATIONS -gt 0 ] && echo $MAX_ITERATIONS || echo 'unlimited')"
echo "Branch: $CURRENT_BRANCH"
echo "==================="
echo ""

# === 메인 루프 ===
while true; do
  # 반복 횟수 제한 체크
  if [ $MAX_ITERATIONS -gt 0 ] && [ $ITERATION -ge $MAX_ITERATIONS ]; then
    echo "최대 반복 횟수($MAX_ITERATIONS)에 도달했습니다."
    break
  fi

  ITERATION=$((ITERATION + 1))
  echo "--- Iteration $ITERATION ---"

  # ✅ DO: headless 모드(-p)로 실행 + 구조화된 출력
  # ❌ DON'T: 대화형 모드로 실행 (자동화 불가)
  cat "$PROMPT_FILE" | claude -p \
    --dangerously-skip-permissions \
    --output-format=stream-json \
    --model opus \
    --verbose

  # ✅ DO: 매 반복마다 push하여 진행 상태를 원격에 보관
  # ❌ DON'T: 로컬에만 커밋하여 작업 유실 위험
  git push origin "$CURRENT_BRANCH" 2>/dev/null || true

  echo "--- Iteration $ITERATION 완료 ---"
  echo ""
done

echo "=== Ralph Loop 종료 (총 ${ITERATION}회 실행) ==="

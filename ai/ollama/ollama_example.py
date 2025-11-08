#!/usr/bin/env python3
"""
Ollama API 사용 예제
"""

import requests
import json
import sys

OLLAMA_URL = "http://localhost:30025"
DEFAULT_MODEL = "qwen2.5:0.5b"

def check_ollama_status():
    """Ollama 서버 상태 확인"""
    try:
        response = requests.get(f"{OLLAMA_URL}/api/tags")
        return response.status_code == 200
    except:
        return False

def list_models():
    """설치된 모델 목록 조회"""
    response = requests.get(f"{OLLAMA_URL}/api/tags")
    models = response.json().get('models', [])
    
    if not models:
        print("설치된 모델이 없습니다.")
        print(f"다음 명령어로 모델을 설치하세요: make model-pull")
        return None
    
    print("설치된 모델 목록:")
    for model in models:
        size_mb = model['size'] / 1024 / 1024
        print(f"  - {model['name']} ({size_mb:.0f}MB)")
    
    return models[0]['name'] if models else None

def generate_text(prompt, model=DEFAULT_MODEL):
    """텍스트 생성"""
    payload = {
        "model": model,
        "prompt": prompt,
        "stream": False
    }
    
    print(f"\n모델 {model}로 응답 생성 중...")
    response = requests.post(f"{OLLAMA_URL}/api/generate", json=payload)
    
    if response.status_code == 200:
        return response.json()['response']
    else:
        return f"오류 발생: {response.status_code}"

def chat_mode():
    """대화형 채팅 모드"""
    model = list_models()
    if not model:
        return
    
    print(f"\n채팅 모드 시작 (모델: {model})")
    print("종료하려면 'exit' 또는 'quit'를 입력하세요.\n")
    
    while True:
        prompt = input("You: ").strip()
        
        if prompt.lower() in ['exit', 'quit']:
            print("채팅을 종료합니다.")
            break
        
        if not prompt:
            continue
        
        response = generate_text(prompt, model)
        print(f"\nOllama: {response}\n")

def main():
    """메인 함수"""
    print("Ollama API 예제")
    print("=" * 50)
    
    # 서버 상태 확인
    if not check_ollama_status():
        print("❌ Ollama 서버에 연결할 수 없습니다.")
        print("다음을 확인하세요:")
        print("  1. Ollama가 배포되었는지: make status")
        print("  2. 포트 30025가 열려있는지")
        sys.exit(1)
    
    print("✅ Ollama 서버 연결 성공")
    
    # 예제 실행
    if len(sys.argv) > 1:
        # 명령줄 인자가 있으면 단일 질문으로 처리
        prompt = " ".join(sys.argv[1:])
        model = list_models()
        if model:
            response = generate_text(prompt, model)
            print(f"\n응답: {response}")
    else:
        # 대화형 모드
        chat_mode()

if __name__ == "__main__":
    main() 
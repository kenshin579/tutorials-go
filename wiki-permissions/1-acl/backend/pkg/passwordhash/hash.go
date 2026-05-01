package passwordhash

import "golang.org/x/crypto/bcrypt"

// Hash는 평문 비밀번호를 bcrypt 기본 cost로 해시한다.
func Hash(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Verify는 평문 비밀번호가 bcrypt 해시와 일치하는지 확인한다.
// 인자 순서 주의: (평문, 해시) — bcrypt 내부 API는 (해시, 평문)이지만 호출 가독성을 위해 뒤집었다.
func Verify(plain, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}

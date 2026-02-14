package go1_25

import "testing"

func Test_Attr_테스트_속성_기록(t *testing.T) {
	// Go 1.25: t.Attr()로 테스트 메타데이터 기록
	// -json 플래그로 실행 시 === ATTR 형태로 출력됨
	t.Attr("version", "1.25")
	t.Attr("category", "runtime")
	t.Attr("priority", "high")

	// 테스트 로직
	t.Log("테스트 속성이 기록되었습니다")
}

func Test_Attr_여러속성(t *testing.T) {
	t.Attr("feature", "synctest")
	t.Attr("type", "unit")

	// go test -v -json ./golang/go1_25/ 으로 실행하면 속성 확인 가능
	t.Log("여러 속성을 기록할 수 있습니다")
}

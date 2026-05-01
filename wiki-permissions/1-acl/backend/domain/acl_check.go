package domain

// EvaluateACL은 page에 대해 userID가 want action을 수행할 수 있는지 결정한다.
//
// 평가 규칙(우선순위 순):
//  1. page == nil 이면 false (방어적 처리).
//  2. page.OwnerID == userID 이면 모든 action 허용.
//  3. entries에 (page, user, want) 매칭 entry가 있으면 허용.
//  4. want == ActionRead 이고 사용자가 ActionEdit 권한을 가진 경우, edit이 read를 함의하므로 허용.
//  5. 위 어디에도 해당하지 않으면 false.
//
// 본 함수는 도메인 로직(순수 함수)으로 외부 부수효과가 없고 테스트가 쉽다.
// usecase 계층에서 호출되며, page와 entries를 미리 조회한 뒤 의사결정만 위임한다.
func EvaluateACL(page *Page, userID uint, want Action, entries []ACLEntry) bool {
	if page == nil {
		return false
	}
	if page.OwnerID == userID {
		return true
	}
	for _, e := range entries {
		if e.UserID != userID || e.PageID != page.ID {
			continue
		}
		if e.Action == want {
			return true
		}
		if want == ActionRead && e.Action == ActionEdit {
			return true
		}
	}
	return false
}

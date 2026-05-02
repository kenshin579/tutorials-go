package domain

// Action은 페이지에 대한 액션이다.
type Action string

const (
	ActionRead Action = "read"
	ActionEdit Action = "edit"
)

// Decision은 정책 평가 결과를 표현한다.
//
//	Allowed: 액션 허용 여부
//	Reason : 사용자에게 보여줄 결정 이유 (UX/감사 로그용)
//	Policy : 어떤 정책이 결정을 내렸는지 식별자 (디버깅/감사용)
//
// 1·2편의 평가는 bool만 반환했지만, ABAC은 결정 이유까지 풍부하게 표현 가능하다.
// 사용자에게 "왜 허용/거부됐는지"를 그대로 보여줄 수 있어 UX/감사가 향상된다.
type Decision struct {
	Allowed bool   `json:"allowed"`
	Reason  string `json:"reason"`
	Policy  string `json:"policy"`
}

// EvaluateABAC는 (사용자, 페이지, 액션) 트리플에 대해 ABAC 정책을 평가한다.
//
// 평가 우선순위 (먼저 매칭되는 정책이 결정):
//
//  1. Owner 정책 — page.OwnerID == user.ID 면 모든 액션 허용 (RBAC가 잃었던 owner 개념을 ABAC가 회복).
//  2. Public 정책 — page.Confidentiality == public 이면 누구나 read 허용 (edit은 거부).
//  3. Department-match 가드 — internal/confidential 페이지는 같은 부서가 아니면 거부.
//  4. Internal 정책 — internal + 같은 부서 → read+edit 허용.
//  5. Confidential 정책 — confidential + 같은 부서 + 정규직(fulltime) → read+edit 허용. contract는 거부.
//  6. Default Deny — 위 어디에도 매칭 안 되면 거부.
//
// 외부 정책 엔진(OPA/Cedar) 미사용 — Go 함수로 직접 구현해 ABAC의 본질(속성 기반 평가)에 집중한다.
// 운영에서는 정책 수가 늘어나면 OPA Rego 등 정책 언어를 도입하는 것이 일반적이다.
func EvaluateABAC(user *User, page *Page, action Action) Decision {
	if user == nil || page == nil {
		return Decision{Allowed: false, Reason: "user or page is nil", Policy: "guard"}
	}

	// 1) Owner — 항상 모든 액션 허용
	if page.OwnerID == user.ID {
		return Decision{Allowed: true, Reason: "owner of the page", Policy: "owner"}
	}

	// 2) Public — read만 허용, edit은 거부
	if page.Confidentiality == ConfidentialityPublic {
		if action == ActionRead {
			return Decision{Allowed: true, Reason: "public page", Policy: "public"}
		}
		return Decision{Allowed: false, Reason: "public page is read-only for non-owners", Policy: "public"}
	}

	// 3) Department-match 가드 — internal/confidential 페이지는 부서 매칭 필수
	if page.DepartmentID == nil || *page.DepartmentID != user.DepartmentID {
		return Decision{Allowed: false, Reason: "different department", Policy: "department-match"}
	}

	// 4) Internal — 같은 부서 사용자에게 read+edit 모두 허용
	if page.Confidentiality == ConfidentialityInternal {
		return Decision{Allowed: true, Reason: "same department, internal page", Policy: "internal"}
	}

	// 5) Confidential — 같은 부서 + 정규직 조건
	if page.Confidentiality == ConfidentialityConfidential {
		if user.EmploymentType != EmploymentFulltime {
			return Decision{Allowed: false, Reason: "confidential pages require fulltime employment", Policy: "confidential"}
		}
		return Decision{Allowed: true, Reason: "same department, fulltime, confidential page", Policy: "confidential"}
	}

	// 6) Default deny
	return Decision{Allowed: false, Reason: "no policy matched", Policy: "default-deny"}
}

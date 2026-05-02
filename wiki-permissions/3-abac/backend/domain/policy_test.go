package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 시드 시나리오의 12 케이스 (4 사용자 × 3 페이지)를 정책 단위로 검증한다.

var (
	engID uint = 1
	mktID uint = 2
)

func newUser(id, deptID uint, et EmploymentType) *User {
	return &User{ID: id, DepartmentID: deptID, EmploymentType: et}
}

func newPage(id, ownerID uint, c Confidentiality, deptID *uint) *Page {
	return &Page{ID: id, OwnerID: ownerID, Confidentiality: c, DepartmentID: deptID}
}

func TestEvaluateABAC_NilGuard(t *testing.T) {
	d := EvaluateABAC(nil, nil, ActionRead)
	assert.False(t, d.Allowed)
	assert.Equal(t, "guard", d.Policy)
}

// alice (Eng/fulltime)이 자기 페이지(Eng Roadmap, internal/Eng) 모든 액션
func TestEvaluateABAC_OwnerAllAccess(t *testing.T) {
	alice := newUser(1, engID, EmploymentFulltime)
	engRoadmap := newPage(10, 1, ConfidentialityInternal, &engID)

	for _, action := range []Action{ActionRead, ActionEdit} {
		d := EvaluateABAC(alice, engRoadmap, action)
		assert.True(t, d.Allowed, "owner should allow %s", action)
		assert.Equal(t, "owner", d.Policy)
	}
}

// public 페이지 — read 허용, edit은 owner가 아니면 거부
func TestEvaluateABAC_PublicReadButNotEdit(t *testing.T) {
	bob := newUser(2, engID, EmploymentFulltime)
	onboarding := newPage(30, 1, ConfidentialityPublic, nil) // owner = alice

	r := EvaluateABAC(bob, onboarding, ActionRead)
	assert.True(t, r.Allowed)
	assert.Equal(t, "public", r.Policy)

	e := EvaluateABAC(bob, onboarding, ActionEdit)
	assert.False(t, e.Allowed)
	assert.Equal(t, "public", e.Policy)
}

// bob (Eng/fulltime)이 Engineering Roadmap (internal/Eng) → 같은 부서, read+edit 허용
func TestEvaluateABAC_InternalSameDepartment(t *testing.T) {
	bob := newUser(2, engID, EmploymentFulltime)
	engRoadmap := newPage(10, 1, ConfidentialityInternal, &engID)

	for _, action := range []Action{ActionRead, ActionEdit} {
		d := EvaluateABAC(bob, engRoadmap, action)
		assert.True(t, d.Allowed)
		assert.Equal(t, "internal", d.Policy)
	}
}

// carol (Mkt)이 Engineering Roadmap (internal/Eng) → 다른 부서 거부
func TestEvaluateABAC_InternalDifferentDepartment(t *testing.T) {
	carol := newUser(3, mktID, EmploymentFulltime)
	engRoadmap := newPage(10, 1, ConfidentialityInternal, &engID)

	d := EvaluateABAC(carol, engRoadmap, ActionRead)
	assert.False(t, d.Allowed)
	assert.Equal(t, "department-match", d.Policy)
}

// carol (Mkt/fulltime)이 Q4 Marketing Plan (confidential/Mkt) → 같은 부서 + 정규직 → 허용
// (단, carol은 owner라 owner 정책이 먼저 매칭됨. 따라서 다른 fulltime을 시뮬레이션)
func TestEvaluateABAC_ConfidentialFulltime(t *testing.T) {
	otherMktFulltime := newUser(99, mktID, EmploymentFulltime) // 가상의 Marketing 정규직
	q4 := newPage(20, 3, ConfidentialityConfidential, &mktID)  // owner = carol(3)

	for _, action := range []Action{ActionRead, ActionEdit} {
		d := EvaluateABAC(otherMktFulltime, q4, action)
		assert.True(t, d.Allowed)
		assert.Equal(t, "confidential", d.Policy)
	}
}

// dave (Mkt/contract)가 Q4 (confidential/Mkt) → 같은 부서이지만 contract → 거부
func TestEvaluateABAC_ConfidentialContractRejected(t *testing.T) {
	dave := newUser(4, mktID, EmploymentContract)
	q4 := newPage(20, 3, ConfidentialityConfidential, &mktID)

	d := EvaluateABAC(dave, q4, ActionRead)
	assert.False(t, d.Allowed)
	assert.Equal(t, "confidential", d.Policy)
	assert.Contains(t, d.Reason, "fulltime")
}

// alice (Eng)가 Q4 (confidential/Mkt) → 다른 부서 → department-match에서 거부
func TestEvaluateABAC_ConfidentialDifferentDepartment(t *testing.T) {
	alice := newUser(1, engID, EmploymentFulltime)
	q4 := newPage(20, 3, ConfidentialityConfidential, &mktID)

	d := EvaluateABAC(alice, q4, ActionRead)
	assert.False(t, d.Allowed)
	assert.Equal(t, "department-match", d.Policy)
}

// public 페이지에 owner가 접근 — owner 정책이 우선 적용됨
func TestEvaluateABAC_OwnerOverridesPublicEditRestriction(t *testing.T) {
	alice := newUser(1, engID, EmploymentFulltime)
	onboarding := newPage(30, 1, ConfidentialityPublic, nil) // owner = alice

	d := EvaluateABAC(alice, onboarding, ActionEdit)
	assert.True(t, d.Allowed)
	assert.Equal(t, "owner", d.Policy) // public이 아니라 owner로 결정
}

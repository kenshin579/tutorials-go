# Keycloak Database Schema 분석

## 개요

**Keycloak**은 오픈소스 IAM(Identity and Access Management) 솔루션으로, SSO(Single Sign-On), OAuth 2.0, OpenID Connect, SAML 2.0을 지원합니다.

- **버전**: Keycloak 기반 MySQL 스키마
- **데이터베이스**: MySQL 8.0.25+
- **테이블 수**: 약 90개

> **주의**: 이 스키마 덤프 파일에는 실제 사용자 데이터가 포함되어 있습니다. 학습 목적으로만 사용하고, 민감한 정보는 삭제 후 사용하세요.

---

## 핵심 도메인 구조

```
┌─────────────────────────────────────────────────────────────────────┐
│                           REALM (영역)                               │
│  - 보안 정책, 인증 설정의 최상위 격리 단위                              │
└───────────────────────────────┬─────────────────────────────────────┘
                                │
        ┌───────────────────────┼───────────────────────┐
        │                       │                       │
        ▼                       ▼                       ▼
┌───────────────┐     ┌───────────────┐     ┌───────────────────┐
│    CLIENT     │     │  USER_ENTITY  │     │   KEYCLOAK_ROLE   │
│ (애플리케이션)  │     │    (사용자)    │     │      (역할)        │
└───────┬───────┘     └───────┬───────┘     └─────────┬─────────┘
        │                     │                       │
        │    ┌────────────────┼────────────────┐      │
        │    │                │                │      │
        ▼    ▼                ▼                ▼      ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│  CLIENT_SCOPE   │  │ USER_ATTRIBUTE  │  │  ROLE_MAPPING   │
│   (스코프)       │  │  (사용자 속성)   │  │  (역할 매핑)     │
└─────────────────┘  └─────────────────┘  └─────────────────┘
```

---

## ERD (Entity Relationship Diagram)

### 전체 ERD

```
┌─────────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                              REALM                                                   │
│                                         (최상위 격리 단위)                                            │
└───────────────┬─────────────────────────────────────┬───────────────────────────────────────────────┘
                │                                     │
    ┌───────────┴───────────┬─────────────┬──────────┴────────────┬─────────────────────┐
    │                       │             │                       │                     │
    ▼                       ▼             ▼                       ▼                     ▼
┌─────────────┐     ┌─────────────┐  ┌─────────────┐     ┌─────────────────┐    ┌─────────────────┐
│   CLIENT    │     │ USER_ENTITY │  │KEYCLOAK_ROLE│     │ KEYCLOAK_GROUP  │    │IDENTITY_PROVIDER│
│ (애플리케이션)│     │   (사용자)   │  │    (역할)    │     │     (그룹)       │    │    (소셜로그인)  │
└──────┬──────┘     └──────┬──────┘  └──────┬──────┘     └────────┬────────┘    └────────┬────────┘
       │                   │                │                     │                     │
       │      ┌────────────┼────────────────┼─────────────────────┼─────────────────────┘
       │      │            │                │                     │
       ▼      ▼            ▼                │                     ▼
┌───────────────────────────────────┐       │            ┌─────────────────┐
│       USER_ROLE_MAPPING           │◄──────┘            │ GROUP_ATTRIBUTE │
│   (사용자 ↔ 역할 N:M 매핑)          │                    └─────────────────┘
└───────────────────────────────────┘                              │
                                                                   ▼
                                                          ┌─────────────────┐
                                                          │GROUP_ROLE_MAPPING│
                                                          │(그룹 ↔ 역할 매핑) │
                                                          └─────────────────┘
```

### 사용자 관련 ERD

```
                                    ┌────────────────────┐
                                    │    USER_ENTITY     │
                                    ├────────────────────┤
                                    │ ID (PK, varchar36) │
                                    │ EMAIL              │
                                    │ USERNAME           │
                                    │ REALM_ID (FK)      │
                                    │ CREATED_TIMESTAMP  │
                                    └─────────┬──────────┘
                                              │
         ┌────────────────┬───────────────────┼───────────────────┬────────────────┐
         │                │                   │                   │                │
         ▼                ▼                   ▼                   ▼                ▼
┌─────────────────┐ ┌──────────────┐ ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
│ USER_ATTRIBUTE  │ │  CREDENTIAL  │ │USER_ROLE_MAPPING│ │USER_GROUP_MEMBER│ │FEDERATED_IDENTITY│
├─────────────────┤ ├──────────────┤ ├─────────────────┤ ├─────────────────┤ ├─────────────────┤
│ ID (PK)         │ │ ID (PK)      │ │ USER_ID (FK)    │ │ USER_ID (FK)    │ │ USER_ID (FK)    │
│ USER_ID (FK)    │ │ USER_ID (FK) │ │ ROLE_ID (FK)    │ │ GROUP_ID (FK)   │ │ IDENTITY_PROVIDER│
│ NAME            │ │ TYPE         │ └────────┬────────┘ └────────┬────────┘ │ FEDERATED_USER  │
│ VALUE           │ │ SECRET_DATA  │          │                   │          └─────────────────┘
│ REALM_ID        │ │ CREDENTIAL   │          ▼                   ▼
└─────────────────┘ │ _DATA        │ ┌─────────────────┐ ┌─────────────────┐
                    └──────────────┘ │  KEYCLOAK_ROLE  │ │ KEYCLOAK_GROUP  │
                                     ├─────────────────┤ ├─────────────────┤
                                     │ ID (PK)         │ │ ID (PK)         │
                                     │ NAME            │ │ NAME            │
                                     │ REALM_ID (FK)   │ │ REALM_ID (FK)   │
                                     │ CLIENT (FK)     │ │ PARENT_GROUP    │
                                     │ CLIENT_ROLE     │ └─────────────────┘
                                     └─────────────────┘
```

### 클라이언트 관련 ERD

```
                                    ┌────────────────────┐
                                    │       CLIENT       │
                                    ├────────────────────┤
                                    │ ID (PK, varchar36) │
                                    │ CLIENT_ID          │
                                    │ SECRET             │
                                    │ REALM_ID (FK)      │
                                    │ PROTOCOL           │
                                    │ PUBLIC_CLIENT      │
                                    └─────────┬──────────┘
                                              │
         ┌────────────────┬───────────────────┼───────────────────┬────────────────┐
         │                │                   │                   │                │
         ▼                ▼                   ▼                   ▼                ▼
┌─────────────────┐ ┌──────────────┐ ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
│CLIENT_ATTRIBUTES│ │REDIRECT_URIS │ │  WEB_ORIGINS    │ │ SCOPE_MAPPING   │ │PROTOCOL_MAPPER  │
├─────────────────┤ ├──────────────┤ ├─────────────────┤ ├─────────────────┤ ├─────────────────┤
│ CLIENT_ID (FK)  │ │ CLIENT_ID(FK)│ │ CLIENT_ID (FK)  │ │ CLIENT_ID (FK)  │ │ ID (PK)         │
│ NAME            │ │ VALUE        │ │ VALUE           │ │ ROLE_ID (FK)    │ │ CLIENT_ID (FK)  │
│ VALUE           │ └──────────────┘ └─────────────────┘ └─────────────────┘ │ CLIENT_SCOPE_ID │
└─────────────────┘                                                          │ NAME            │
                                                                             │ PROTOCOL        │
       ┌─────────────────────────────────────────────────────────────────────┴─────────────────┐
       │                                                                                       │
       ▼                                                                                       ▼
┌─────────────────────┐                                                      ┌─────────────────────┐
│CLIENT_SCOPE_CLIENT  │                                                      │PROTOCOL_MAPPER_CONFIG│
├─────────────────────┤                                                      ├─────────────────────┤
│ CLIENT_ID (FK)      │◄───────────────────┐                                 │PROTOCOL_MAPPER_ID(FK)│
│ SCOPE_ID (FK)       │                    │                                 │ NAME                │
│ DEFAULT_SCOPE       │                    │                                 │ VALUE               │
└──────────┬──────────┘                    │                                 └─────────────────────┘
           │                               │
           ▼                               │
┌─────────────────────┐                    │
│    CLIENT_SCOPE     │────────────────────┘
├─────────────────────┤
│ ID (PK)             │
│ NAME                │
│ REALM_ID (FK)       │
│ DESCRIPTION         │
│ PROTOCOL            │
└─────────────────────┘
```

### 인증 흐름 ERD

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                                    REALM                                     │
└───────────────────────────────────────┬─────────────────────────────────────┘
                                        │
                                        ▼
                            ┌───────────────────────┐
                            │  AUTHENTICATION_FLOW  │
                            ├───────────────────────┤
                            │ ID (PK)               │
                            │ ALIAS                 │
                            │ REALM_ID (FK)         │
                            │ PROVIDER_ID           │
                            │ TOP_LEVEL             │
                            │ BUILT_IN              │
                            └───────────┬───────────┘
                                        │
                                        ▼
                         ┌─────────────────────────────┐
                         │  AUTHENTICATION_EXECUTION   │
                         ├─────────────────────────────┤
                         │ ID (PK)                     │
                         │ FLOW_ID (FK)                │
                         │ REALM_ID (FK)               │
                         │ AUTHENTICATOR              │
                         │ AUTH_FLOW_ID               │◄──── (Self Reference)
                         │ REQUIREMENT                │
                         │ PRIORITY                   │
                         │ AUTHENTICATOR_FLOW         │
                         │ FLOW_ID                    │
                         └─────────────┬───────────────┘
                                       │
                                       ▼
                          ┌─────────────────────────┐
                          │  AUTHENTICATOR_CONFIG   │
                          ├─────────────────────────┤
                          │ ID (PK)                 │
                          │ REALM_ID (FK)           │
                          │ ALIAS                   │
                          └─────────────┬───────────┘
                                        │
                                        ▼
                       ┌─────────────────────────────────┐
                       │  AUTHENTICATOR_CONFIG_ENTRY    │
                       ├─────────────────────────────────┤
                       │ AUTHENTICATOR_ID (FK)          │
                       │ NAME                           │
                       │ VALUE                          │
                       └─────────────────────────────────┘
```

### Authorization Services ERD

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                               RESOURCE_SERVER                                │
│                            (리소스 서버 = CLIENT)                             │
├─────────────────────────────────────────────────────────────────────────────┤
│ ID (PK)  │  ALLOW_RS_REMOTE_MGMT  │  POLICY_ENFORCE_MODE  │  DECISION_STRAT │
└────────────────────────────────────┬────────────────────────────────────────┘
                                     │
           ┌─────────────────────────┼─────────────────────────┐
           │                         │                         │
           ▼                         ▼                         ▼
┌─────────────────────┐   ┌─────────────────────┐   ┌─────────────────────┐
│RESOURCE_SERVER_SCOPE│   │RESOURCE_SERVER_POLICY│   │RESOURCE_SERVER_     │
│       (스코프)       │   │       (정책)         │   │    RESOURCE (리소스) │
├─────────────────────┤   ├─────────────────────┤   ├─────────────────────┤
│ ID (PK)             │   │ ID (PK)             │   │ ID (PK)             │
│ NAME                │   │ NAME                │   │ NAME                │
│ DISPLAY_NAME        │   │ TYPE                │   │ TYPE                │
│ ICON_URI            │   │ LOGIC               │   │ URI                 │
│ RESOURCE_SERVER_ID  │   │ DECISION_STRATEGY   │   │ ICON_URI            │
│       (FK)          │   │ RESOURCE_SERVER_ID  │   │ OWNER               │
└─────────┬───────────┘   │       (FK)          │   │ RESOURCE_SERVER_ID  │
          │               └─────────┬───────────┘   │       (FK)          │
          │                         │               └─────────┬───────────┘
          │                         │                         │
          │           ┌─────────────┼─────────────┐           │
          │           │             │             │           │
          │           ▼             ▼             ▼           │
          │   ┌─────────────┐ ┌──────────────┐ ┌─────────────┐│
          │   │POLICY_CONFIG│ │SCOPE_POLICY  │ │RESOURCE_    ││
          │   ├─────────────┤ ├──────────────┤ │   POLICY    ││
          │   │POLICY_ID(FK)│ │POLICY_ID(FK) │ ├─────────────┤│
          │   │ NAME        │ │SCOPE_ID (FK) │◄┤POLICY_ID(FK)││
          │   │ VALUE       │ └──────────────┘ │RESOURCE_ID  │◄┘
          │   └─────────────┘        ▲         │    (FK)     │
          │                          │         └─────────────┘
          └──────────────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │  RESOURCE_SCOPE │
                    ├─────────────────┤
                    │ RESOURCE_ID(FK) │
                    │ SCOPE_ID (FK)   │
                    └─────────────────┘
```

---

## 주요 테이블 카테고리

### 1. Realm 관리 (핵심)

| 테이블 | 설명 |
|--------|------|
| `REALM` | 영역 기본 설정 (SSO 타임아웃, 토큰 설정 등) |
| `REALM_ATTRIBUTE` | 영역별 추가 속성 (Key-Value) |
| `REALM_REQUIRED_CREDENTIAL` | 필수 인증 수단 설정 |
| `REALM_SMTP_CONFIG` | 이메일 SMTP 설정 |
| `REALM_EVENTS_LISTENERS` | 이벤트 리스너 설정 |
| `REALM_SUPPORTED_LOCALES` | 지원 언어 설정 |
| `REALM_LOCALIZATIONS` | 다국어 메시지 |

**설계 포인트**:
- `REALM` 테이블이 전체 스키마의 중심
- 대부분의 테이블이 `REALM_ID`로 FK 참조
- **Multi-tenancy**: 여러 조직을 하나의 Keycloak에서 관리

### 2. 사용자 관리

| 테이블 | 설명 |
|--------|------|
| `USER_ENTITY` | 사용자 기본 정보 |
| `USER_ATTRIBUTE` | 사용자 커스텀 속성 |
| `USER_REQUIRED_ACTION` | 필수 액션 (비밀번호 변경 등) |
| `USER_ROLE_MAPPING` | 사용자-역할 매핑 |
| `USER_GROUP_MEMBERSHIP` | 사용자-그룹 매핑 |
| `USER_CONSENT` | 사용자 동의 정보 |
| `CREDENTIAL` | 인증 정보 (비밀번호 해시 등) |

```sql
-- USER_ENTITY 구조
CREATE TABLE USER_ENTITY (
  ID varchar(36) PRIMARY KEY,
  EMAIL varchar(255),
  EMAIL_VERIFIED tinyint DEFAULT 0,
  ENABLED tinyint DEFAULT 1,
  FIRST_NAME varchar(255),
  LAST_NAME varchar(255),
  REALM_ID varchar(36),           -- FK to REALM
  USERNAME varchar(255),
  CREATED_TIMESTAMP bigint,
  SERVICE_ACCOUNT_CLIENT_LINK varchar(255),
  NOT_BEFORE int DEFAULT 0
);
```

**설계 포인트**:
- UUID 형태의 ID 사용 (`varchar(36)`)
- `USER_ATTRIBUTE`로 확장 속성 관리 (EAV 패턴)
- `CREDENTIAL`과 분리하여 인증 정보 별도 관리

### 3. 클라이언트 (애플리케이션) 관리

| 테이블 | 설명 |
|--------|------|
| `CLIENT` | 클라이언트 앱 정보 |
| `CLIENT_ATTRIBUTES` | 클라이언트 추가 속성 |
| `CLIENT_SCOPE` | 클라이언트 스코프 정의 |
| `CLIENT_SCOPE_CLIENT` | 클라이언트-스코프 매핑 |
| `REDIRECT_URIS` | 허용된 리다이렉트 URI |
| `WEB_ORIGINS` | CORS 허용 origin |

```sql
-- CLIENT 주요 컬럼
CREATE TABLE CLIENT (
  ID varchar(36) PRIMARY KEY,
  CLIENT_ID varchar(255),              -- 클라이언트 식별자
  ENABLED tinyint DEFAULT 0,
  SECRET varchar(255),                 -- Client Secret
  PROTOCOL varchar(255),               -- openid-connect, saml
  PUBLIC_CLIENT tinyint DEFAULT 0,     -- Public vs Confidential
  REALM_ID varchar(36),
  CONSENT_REQUIRED tinyint DEFAULT 0,
  STANDARD_FLOW_ENABLED tinyint DEFAULT 1,
  DIRECT_ACCESS_GRANTS_ENABLED tinyint DEFAULT 0
);
```

**설계 포인트**:
- OAuth 2.0 Grant Type을 Boolean 플래그로 관리
- Public/Confidential 클라이언트 구분

### 4. 역할 및 권한 (RBAC)

| 테이블 | 설명 |
|--------|------|
| `KEYCLOAK_ROLE` | 역할 정의 |
| `KEYCLOAK_GROUP` | 그룹 정의 |
| `COMPOSITE_ROLE` | 복합 역할 (역할의 역할) |
| `SCOPE_MAPPING` | 클라이언트 스코프-역할 매핑 |
| `GROUP_ROLE_MAPPING` | 그룹-역할 매핑 |

```sql
-- 역할 구조
CREATE TABLE KEYCLOAK_ROLE (
  ID varchar(36) PRIMARY KEY,
  NAME varchar(255),
  REALM_ID varchar(36),
  CLIENT varchar(36),                  -- NULL이면 Realm Role
  CLIENT_REALM_CONSTRAINT varchar(255),
  CLIENT_ROLE tinyint DEFAULT 0,
  DESCRIPTION varchar(255)
);
```

**설계 포인트**:
- **Realm Role vs Client Role**: `CLIENT` 컬럼으로 구분
- **Composite Role**: 역할 상속 지원
- 그룹을 통한 역할 일괄 할당

### 5. 인증 흐름 (Authentication Flow)

| 테이블 | 설명 |
|--------|------|
| `AUTHENTICATION_FLOW` | 인증 흐름 정의 |
| `AUTHENTICATION_EXECUTION` | 흐름 내 실행 단계 |
| `AUTHENTICATOR_CONFIG` | 인증기 설정 |
| `REQUIRED_ACTION_PROVIDER` | 필수 액션 제공자 |

**설계 포인트**:
- 인증 프로세스를 모듈화된 단계로 구성
- 각 단계의 순서(`PRIORITY`)와 필수 여부(`REQUIREMENT`) 관리
- 커스텀 인증 흐름 구성 가능

### 6. Identity Provider (소셜 로그인 등)

| 테이블 | 설명 |
|--------|------|
| `IDENTITY_PROVIDER` | IdP 설정 (Google, GitHub 등) |
| `IDENTITY_PROVIDER_CONFIG` | IdP 상세 설정 |
| `IDENTITY_PROVIDER_MAPPER` | 클레임 매핑 설정 |
| `FEDERATED_IDENTITY` | 연동된 외부 계정 정보 |

```sql
CREATE TABLE IDENTITY_PROVIDER (
  INTERNAL_ID varchar(36) PRIMARY KEY,
  ENABLED tinyint DEFAULT 0,
  PROVIDER_ALIAS varchar(255),         -- 'google', 'github'
  PROVIDER_ID varchar(255),
  REALM_ID varchar(36),
  TRUST_EMAIL tinyint DEFAULT 0,
  FIRST_BROKER_LOGIN_FLOW_ID varchar(36)
);
```

### 7. 세션 관리

| 테이블 | 설명 |
|--------|------|
| `USER_SESSION` | 사용자 세션 |
| `CLIENT_SESSION` | 클라이언트별 세션 |
| `OFFLINE_USER_SESSION` | 오프라인 세션 (Remember Me) |
| `OFFLINE_CLIENT_SESSION` | 오프라인 클라이언트 세션 |

### 8. Authorization Services (Fine-Grained Authorization)

| 테이블 | 설명 |
|--------|------|
| `RESOURCE_SERVER` | 리소스 서버 설정 |
| `RESOURCE_SERVER_RESOURCE` | 보호 대상 리소스 |
| `RESOURCE_SERVER_SCOPE` | 리소스 스코프 |
| `RESOURCE_SERVER_POLICY` | 권한 정책 |
| `RESOURCE_SERVER_PERM_TICKET` | 권한 티켓 |

**설계 포인트**:
- UMA 2.0 기반 Fine-Grained Authorization
- Resource, Scope, Policy 조합으로 세밀한 권한 제어

### 9. 이벤트 및 감사

| 테이블 | 설명 |
|--------|------|
| `EVENT_ENTITY` | 사용자 이벤트 (로그인 등) |
| `ADMIN_EVENT_ENTITY` | 관리자 이벤트 |

```sql
CREATE TABLE EVENT_ENTITY (
  ID varchar(36) PRIMARY KEY,
  CLIENT_ID varchar(255),
  DETAILS_JSON text,                   -- JSON 형태의 상세 정보
  ERROR varchar(255),
  IP_ADDRESS varchar(255),
  REALM_ID varchar(255),
  EVENT_TIME bigint,
  TYPE varchar(255),                   -- LOGIN, LOGOUT, REGISTER 등
  USER_ID varchar(255)
);
```

### 10. Protocol Mapper

| 테이블 | 설명 |
|--------|------|
| `PROTOCOL_MAPPER` | 토큰 클레임 매핑 정의 |
| `PROTOCOL_MAPPER_CONFIG` | 매퍼 상세 설정 |

---

## 스키마 설계 패턴 및 배울 점

### 1. UUID 기반 Primary Key
```sql
ID varchar(36) PRIMARY KEY  -- UUID 형식
```
- 분산 환경에서 ID 충돌 방지
- 외부 노출 시 예측 불가

### 2. EAV (Entity-Attribute-Value) 패턴
```
USER_ENTITY (1) ──────── (*) USER_ATTRIBUTE
REALM (1) ──────── (*) REALM_ATTRIBUTE
CLIENT (1) ──────── (*) CLIENT_ATTRIBUTES
```
- 동적 속성 추가 가능
- 스키마 변경 없이 확장

### 3. Multi-Tenancy (Realm 기반 격리)
```
모든 주요 테이블 → REALM_ID (FK)
```
- 하나의 DB에서 여러 조직 관리
- 데이터 격리 보장

### 4. Soft Delete vs Hard Delete
- Keycloak은 주로 **Hard Delete** 사용
- 감사 로그는 별도 이벤트 테이블로 관리

### 5. 시간 저장 방식
```sql
CREATED_TIMESTAMP bigint  -- Unix timestamp (milliseconds)
```
- 타임존 독립적
- 정수 비교로 빠른 쿼리

### 6. Boolean 저장
```sql
ENABLED tinyint DEFAULT 0  -- 0/1 사용
```
- MySQL의 BOOLEAN은 내부적으로 TINYINT(1)

### 7. JSON 데이터 저장
```sql
DETAILS_JSON text  -- 유연한 구조의 데이터
```
- 스키마리스 데이터 저장
- 복잡한 중첩 구조 지원

---

## 주요 FK 관계

```
REALM ─┬─► USER_ENTITY
       ├─► CLIENT
       ├─► KEYCLOAK_ROLE
       ├─► KEYCLOAK_GROUP
       ├─► IDENTITY_PROVIDER
       ├─► AUTHENTICATION_FLOW
       └─► CLIENT_SCOPE

USER_ENTITY ─┬─► USER_ATTRIBUTE
             ├─► USER_ROLE_MAPPING ─► KEYCLOAK_ROLE
             ├─► USER_GROUP_MEMBERSHIP ─► KEYCLOAK_GROUP
             ├─► CREDENTIAL
             └─► FEDERATED_IDENTITY

CLIENT ─┬─► CLIENT_ATTRIBUTES
        ├─► REDIRECT_URIS
        ├─► CLIENT_SCOPE_CLIENT ─► CLIENT_SCOPE
        └─► PROTOCOL_MAPPER
```

---

## 파일 구조

```
keycloak/
├── README.md               # 이 문서
└── keycloak_schema.sql     # 전체 스키마 덤프
    ├── DDL (CREATE TABLE)
    └── DML (INSERT) - 민감정보 포함 주의!
```

---

## 참고 자료

- [Keycloak 공식 문서](https://www.keycloak.org/documentation)
- [Keycloak GitHub](https://github.com/keycloak/keycloak)
- [Keycloak Database Schema](https://www.keycloak.org/docs/latest/server_development/#_database)

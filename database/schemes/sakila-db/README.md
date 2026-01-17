# Sakila Database Schema 분석

## 개요

**Sakila**는 MySQL의 공식 샘플 데이터베이스로, DVD 대여점을 모델링한 스키마입니다. 실무에서 자주 사용되는 다양한 데이터베이스 기능(View, Stored Procedure, Trigger, Function)을 포함하고 있어 중급 이상의 SQL 학습에 적합합니다.

- **출처**: [MySQL Sample Database - Sakila](https://dev.mysql.com/doc/sakila/en/)
- **다운로드**: [MySQL Sakila Database](https://downloads.mysql.com/docs/sakila-db.zip)
- **라이선스**: BSD License (Oracle)
- **MySQL 버전**: 5.7+ (8.0 권장)
- **테이블 수**: 16개
- **뷰**: 7개
- **Stored Procedure/Function**: 6개

---

## ERD (Entity Relationship Diagram)

### 핵심 비즈니스 ERD

```
┌─────────────┐                                    ┌─────────────┐
│   country   │                                    │   language  │
├─────────────┤                                    ├─────────────┤
│ country_id  │◄───┐                               │ language_id │◄──────────┐
│ country     │    │                               │ name        │           │
└─────────────┘    │                               └─────────────┘           │
                   │                                                         │
                   │    ┌─────────────┐                                      │
                   │    │    city     │                                      │
                   │    ├─────────────┤                                      │
                   └────┤ country_id  │                                      │
                        │ city_id     │◄───┐                                 │
                        │ city        │    │                                 │
                        └─────────────┘    │                                 │
                                           │                                 │
                               ┌───────────┴───────────┐                     │
                               │                       │                     │
                               ▼                       ▼                     ▼
                    ┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
                    │     address     │     │      store      │     │      film       │
                    ├─────────────────┤     ├─────────────────┤     ├─────────────────┤
                    │ address_id (PK) │◄────┤ address_id (FK) │     │ film_id (PK)    │
                    │ city_id (FK)    │     │ store_id (PK)   │◄─┐  │ title           │
                    │ address         │     │ manager_staff_id│  │  │ language_id(FK) │
                    │ phone           │     └────────┬────────┘  │  │ rental_duration │
                    │ location (GEO)  │              │           │  │ rental_rate     │
                    └────────┬────────┘              │           │  │ rating (ENUM)   │
                             │                       │           │  │ special_features│
         ┌───────────────────┼───────────────────────┘           │  └────────┬────────┘
         │                   │                                   │           │
         ▼                   ▼                                   │           │
┌─────────────────┐ ┌─────────────────┐                          │           │
│     staff       │ │    customer     │                          │           │
├─────────────────┤ ├─────────────────┤                          │           │
│ staff_id (PK)   │ │ customer_id(PK) │                          │           │
│ address_id (FK) │ │ store_id (FK)   │──────────────────────────┘           │
│ store_id (FK)   │ │ address_id (FK) │                                      │
│ first_name      │ │ first_name      │                                      │
│ email           │ │ email           │                                      │
│ username        │ │ active          │                                      │
│ password        │ │ create_date     │                                      │
│ picture (BLOB)  │ └────────┬────────┘                                      │
└────────┬────────┘          │                                               │
         │                   │                                               │
         │                   │      ┌─────────────────┐                      │
         │                   │      │   inventory     │                      │
         │                   │      ├─────────────────┤                      │
         │                   │      │ inventory_id(PK)│◄──────────┐          │
         │                   │      │ film_id (FK)    │───────────┼──────────┘
         │                   │      │ store_id (FK)   │           │
         │                   │      └────────┬────────┘           │
         │                   │               │                    │
         │                   │               │                    │
         │                   ▼               ▼                    │
         │           ┌───────────────────────────────┐            │
         │           │          rental               │            │
         │           ├───────────────────────────────┤            │
         │           │ rental_id (PK)                │            │
         └──────────►│ staff_id (FK)                 │            │
                     │ customer_id (FK)              │            │
                     │ inventory_id (FK)             │────────────┘
                     │ rental_date                   │
                     │ return_date                   │
                     └───────────────┬───────────────┘
                                     │
                                     ▼
                            ┌─────────────────┐
                            │    payment      │
                            ├─────────────────┤
                            │ payment_id (PK) │
                            │ customer_id(FK) │
                            │ staff_id (FK)   │
                            │ rental_id (FK)  │
                            │ amount          │
                            │ payment_date    │
                            └─────────────────┘
```

### 영화/배우 관계 ERD

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│     actor       │     │   film_actor    │     │      film       │
├─────────────────┤     ├─────────────────┤     ├─────────────────┤
│ actor_id (PK)   │◄────┤ actor_id (FK)   │     │ film_id (PK)    │
│ first_name      │     │ film_id (FK)    │────►│ title           │
│ last_name       │     │ last_update     │     │ description     │
│ last_update     │     └─────────────────┘     │ release_year    │
└─────────────────┘              N:M             │ language_id     │
                                                │ rating          │
┌─────────────────┐     ┌─────────────────┐     └────────┬────────┘
│    category     │     │  film_category  │              │
├─────────────────┤     ├─────────────────┤              │
│ category_id(PK) │◄────┤ category_id(FK) │              │
│ name            │     │ film_id (FK)    │──────────────┘
│ last_update     │     │ last_update     │
└─────────────────┘     └─────────────────┘
                                N:M

                        ┌─────────────────┐
                        │   film_text     │     (FULLTEXT 검색용)
                        ├─────────────────┤
                        │ film_id (PK)    │◄─── Trigger로 동기화
                        │ title           │
                        │ description     │
                        │ [FULLTEXT IDX]  │
                        └─────────────────┘
```

---

## 테이블 상세

### 1. 마스터 테이블

#### `actor` - 배우 정보
| 컬럼 | 타입 | 설명 |
|------|------|------|
| actor_id | SMALLINT UNSIGNED (PK) | 배우 ID |
| first_name | VARCHAR(45) | 이름 |
| last_name | VARCHAR(45) | 성 (INDEX) |
| last_update | TIMESTAMP | 최종 수정일 |

#### `category` - 영화 카테고리
| 컬럼 | 타입 | 설명 |
|------|------|------|
| category_id | TINYINT UNSIGNED (PK) | 카테고리 ID |
| name | VARCHAR(25) | 카테고리명 (Action, Comedy 등) |

#### `language` - 언어
| 컬럼 | 타입 | 설명 |
|------|------|------|
| language_id | TINYINT UNSIGNED (PK) | 언어 ID |
| name | CHAR(20) | 언어명 (English, Italian 등) |

### 2. 지역 테이블

#### `country` → `city` → `address`
계층 구조로 주소 정규화

```sql
-- 전체 주소 조회
SELECT a.address, c.city, co.country
FROM address a
JOIN city c ON a.city_id = c.city_id
JOIN country co ON c.country_id = co.country_id;
```

#### `address` 특징
- **SPATIAL INDEX**: `location` 컬럼에 지리 데이터 저장
- MySQL 5.7.5+ 에서 GEOMETRY 타입 지원

### 3. 핵심 비즈니스 테이블

#### `film` - 영화 정보
```sql
CREATE TABLE film (
  film_id SMALLINT UNSIGNED NOT NULL AUTO_INCREMENT,
  title VARCHAR(128) NOT NULL,
  description TEXT,
  release_year YEAR,
  language_id TINYINT UNSIGNED NOT NULL,        -- FK
  original_language_id TINYINT UNSIGNED,         -- FK (NULL 허용)
  rental_duration TINYINT UNSIGNED DEFAULT 3,    -- 대여 기간
  rental_rate DECIMAL(4,2) DEFAULT 4.99,         -- 대여 요금
  length SMALLINT UNSIGNED,                       -- 상영 시간(분)
  replacement_cost DECIMAL(5,2) DEFAULT 19.99,   -- 분실 시 비용
  rating ENUM('G','PG','PG-13','R','NC-17'),     -- 등급
  special_features SET('Trailers','Commentaries','Deleted Scenes','Behind the Scenes'),
  PRIMARY KEY (film_id)
);
```

**설계 포인트**:
- **ENUM**: 등급 제한
- **SET**: 다중 선택 가능한 특별 기능
- **YEAR**: 연도만 저장하는 타입

#### `inventory` - 재고 (DVD 개별 복사본)
```sql
-- 하나의 영화에 여러 개의 DVD 재고 존재
SELECT film_id, store_id, COUNT(*) as copies
FROM inventory
GROUP BY film_id, store_id;
```

#### `rental` - 대여 기록
| 컬럼 | 타입 | 설명 |
|------|------|------|
| rental_id | INT (PK) | 대여 ID |
| rental_date | DATETIME | 대여일 |
| inventory_id | MEDIUMINT (FK) | 재고 ID |
| customer_id | SMALLINT (FK) | 고객 ID |
| return_date | DATETIME | 반납일 (NULL: 미반납) |
| staff_id | TINYINT (FK) | 담당 직원 |

**UNIQUE 제약**: `(rental_date, inventory_id, customer_id)` - 동시 중복 대여 방지

#### `payment` - 결제
| 컬럼 | 타입 | 설명 |
|------|------|------|
| payment_id | SMALLINT (PK) | 결제 ID |
| customer_id | SMALLINT (FK) | 고객 |
| staff_id | TINYINT (FK) | 담당 직원 |
| rental_id | INT (FK) | 대여 ID (NULL 허용 - 기타 결제) |
| amount | DECIMAL(5,2) | 결제 금액 |
| payment_date | DATETIME | 결제일 |

### 4. 운영 테이블

#### `store` - 매장
```sql
CREATE TABLE store (
  store_id TINYINT UNSIGNED PRIMARY KEY,
  manager_staff_id TINYINT UNSIGNED NOT NULL UNIQUE,  -- 1:1 관계
  address_id SMALLINT UNSIGNED NOT NULL
);
```
- 매니저는 한 매장만 관리 가능 (UNIQUE)

#### `staff` - 직원
- `picture BLOB`: 프로필 사진 저장
- `password VARCHAR(40)`: SHA1 해시 저장 (레거시)

#### `customer` - 고객
- `active BOOLEAN`: 활성/비활성 상태
- `create_date DATETIME`: 가입일

---

## View 목록

| View | 설명 | 용도 |
|------|------|------|
| `customer_list` | 고객 + 주소 조인 | 고객 목록 조회 |
| `staff_list` | 직원 + 주소 조인 | 직원 목록 조회 |
| `film_list` | 영화 + 카테고리 + 배우 | 영화 상세 조회 |
| `actor_info` | 배우별 출연 영화 | 배우 정보 조회 |
| `sales_by_store` | 매장별 매출 집계 | 매출 분석 |
| `sales_by_film_category` | 카테고리별 매출 | 매출 분석 |
| `nicer_but_slower_film_list` | 영화 목록 (대소문자 정리) | UI 표시용 |

---

## Stored Procedure/Function

### Procedures

```sql
-- 재고 확인
CALL film_in_stock(1, 1, @count);
SELECT @count;

-- 미반납 재고 확인
CALL film_not_in_stock(1, 1, @count);

-- 우수 고객 리포트
CALL rewards_report(5, 25.00, @count);
```

### Functions

```sql
-- 고객 잔액 계산 (대여료 + 연체료 - 결제액)
SELECT get_customer_balance(1, NOW());

-- 재고 보유 고객 확인
SELECT inventory_held_by_customer(10);

-- 재고 가용 여부
SELECT inventory_in_stock(10);
```

---

## Trigger

```sql
-- film 테이블 변경 시 film_text 자동 동기화
ins_film  -- INSERT 후 film_text에 추가
upd_film  -- UPDATE 후 film_text 갱신
del_film  -- DELETE 후 film_text에서 삭제
```

**목적**: `film_text`의 FULLTEXT INDEX를 활용한 전문 검색 지원

---

## 학습 포인트

### 1. 복합 조인
```sql
-- 고객별 대여 내역 (5개 테이블 조인)
SELECT c.first_name, c.last_name, f.title,
       r.rental_date, r.return_date, p.amount
FROM customer c
JOIN rental r ON c.customer_id = r.customer_id
JOIN inventory i ON r.inventory_id = i.inventory_id
JOIN film f ON i.film_id = f.film_id
JOIN payment p ON r.rental_id = p.rental_id
WHERE c.customer_id = 1;
```

### 2. 집계 및 분석
```sql
-- 카테고리별 평균 대여 기간
SELECT c.name, AVG(f.rental_duration) as avg_duration
FROM category c
JOIN film_category fc ON c.category_id = fc.category_id
JOIN film f ON fc.film_id = f.film_id
GROUP BY c.name
ORDER BY avg_duration DESC;

-- 월별 매출 추이
SELECT DATE_FORMAT(payment_date, '%Y-%m') as month,
       SUM(amount) as total_sales
FROM payment
GROUP BY month
ORDER BY month;
```

### 3. FULLTEXT 검색
```sql
-- 영화 제목/설명에서 검색
SELECT * FROM film_text
WHERE MATCH(title, description) AGAINST('action' IN NATURAL LANGUAGE MODE);
```

### 4. ENUM/SET 활용
```sql
-- R등급 영화 조회
SELECT * FROM film WHERE rating = 'R';

-- 삭제 장면이 포함된 영화
SELECT * FROM film
WHERE FIND_IN_SET('Deleted Scenes', special_features) > 0;
```

---

## 스키마 설계 특징

### 1. 적절한 데이터 타입 선택
- `TINYINT`: category, language (소량 데이터)
- `SMALLINT`: actor, customer (중간 규모)
- `MEDIUMINT`: inventory (대용량)
- `INT`: rental, payment (트랜잭션)

### 2. last_update 패턴
```sql
last_update TIMESTAMP NOT NULL
  DEFAULT CURRENT_TIMESTAMP
  ON UPDATE CURRENT_TIMESTAMP
```
모든 테이블에 자동 갱신 타임스탬프 적용

### 3. ON DELETE/UPDATE 전략
- `RESTRICT`: 참조 무결성 위반 시 거부
- `CASCADE`: 상위 테이블 변경 시 하위도 변경
- `SET NULL`: 상위 삭제 시 NULL 설정 (payment.rental_id)

### 4. 인덱스 전략
- FK 컬럼에 자동 INDEX
- 검색 빈도 높은 컬럼(last_name, title)에 추가 INDEX
- 복합 INDEX: `(store_id, film_id)` - 자주 같이 조회

---

## 파일 구조

```
sakila-db/
├── README.md           # 이 문서
├── sakila-schema.sql   # DDL (테이블, 뷰, 프로시저, 함수, 트리거)
├── sakila-data.sql     # DML (INSERT 문)
└── sakila.mwb          # MySQL Workbench 모델 파일
```

---

## 설치 방법

```bash
# 1. 스키마 생성
mysql -u root -p < sakila-schema.sql

# 2. 데이터 로드
mysql -u root -p < sakila-data.sql

# 3. 확인
mysql -u root -p -e "USE sakila; SHOW TABLES;"
```

---

## 참고 자료

- [Sakila Sample Database](https://dev.mysql.com/doc/sakila/en/)
- [MySQL Workbench Data Modeling](https://dev.mysql.com/doc/workbench/en/wb-data-modeling.html)

# Menagerie Database Schema 분석

## 개요

**Menagerie**는 MySQL 공식 튜토리얼에서 사용되는 가장 기본적인 학습용 데이터베이스입니다. 반려동물과 그 이벤트를 관리하는 단순한 구조로, SQL 기초를 배우기에 적합합니다.

- **출처**: [MySQL Reference Manual - Tutorial](https://dev.mysql.com/doc/refman/8.0/en/tutorial.html)
- **다운로드**: [MySQL Sample Database - Menagerie](https://downloads.mysql.com/docs/menagerie-db.zip)
- **용도**: SQL 기초 학습, CRUD 연습
- **테이블 수**: 2개 (pet, event)

---

## ERD (Entity Relationship Diagram)

```
┌─────────────────────┐          ┌─────────────────────┐
│        pet          │          │       event         │
├─────────────────────┤          ├─────────────────────┤
│  name    VARCHAR(20)│◄─ ─ ─ ─ ─│  name    VARCHAR(20)│
│  owner   VARCHAR(20)│   약한    │  date    DATE       │
│  species VARCHAR(20)│   참조    │  type    VARCHAR(15)│
│  sex     CHAR(1)    │          │  remark  VARCHAR(255)│
│  birth   DATE       │          └─────────────────────┘
│  death   DATE       │
└─────────────────────┘

* 명시적 FK/PK 없음 - 학습용 단순 설계
* pet.name과 event.name으로 논리적 연결
```

---

## 테이블 상세

### 1. `pet` - 반려동물 정보

```sql
CREATE TABLE pet (
  name    VARCHAR(20),   -- 이름 (논리적 식별자)
  owner   VARCHAR(20),   -- 주인 이름
  species VARCHAR(20),   -- 종류 (cat, dog, bird, snake)
  sex     CHAR(1),       -- 성별 (m: 수컷, f: 암컷, NULL: 미확인)
  birth   DATE,          -- 생년월일
  death   DATE           -- 사망일 (NULL: 생존)
);
```

**샘플 데이터**:
| name | owner | species | sex | birth | death |
|------|-------|---------|-----|-------|-------|
| Fluffy | Harold | cat | f | 1993-02-04 | NULL |
| Claws | Gwen | cat | m | 1994-03-17 | NULL |
| Buffy | Harold | dog | f | 1989-05-13 | NULL |
| Fang | Benny | dog | m | 1990-08-27 | NULL |
| Bowser | Diane | dog | m | 1979-08-31 | 1995-07-29 |
| Chirpy | Gwen | bird | f | 1998-09-11 | NULL |
| Whistler | Gwen | bird | NULL | 1997-12-09 | NULL |
| Slim | Benny | snake | m | 1996-04-29 | NULL |

### 2. `event` - 반려동물 이벤트 기록

```sql
CREATE TABLE event (
  name   VARCHAR(20),    -- 반려동물 이름 (pet 참조)
  date   DATE,           -- 이벤트 날짜
  type   VARCHAR(15),    -- 이벤트 유형
  remark VARCHAR(255)    -- 비고
);
```

**이벤트 유형**:
- `litter`: 출산 기록
- `vet`: 동물병원 방문
- `kennel`: 애완동물 호텔 이용
- `birthday`: 생일

**샘플 데이터**:
| name | date | type | remark |
|------|------|------|--------|
| Fluffy | 1995-05-15 | litter | 4 kittens, 3 female, 1 male |
| Buffy | 1993-06-23 | litter | 5 puppies, 2 female, 3 male |
| Chirpy | 1999-03-21 | vet | needed beak straightened |
| Slim | 1997-08-03 | vet | broken rib |
| Fang | 1998-08-28 | birthday | Gave him a new chew toy |

---

## 스키마 설계 특징

### 1. 의도적으로 단순한 설계
- **Primary Key 없음**: 학습 초기 단계에서 복잡성 제거
- **Foreign Key 없음**: 논리적 관계만 존재 (name 컬럼으로 연결)
- **인덱스 없음**: 소량 데이터에서 불필요

### 2. NULL 값 활용
```sql
-- 생존 중인 동물
SELECT * FROM pet WHERE death IS NULL;

-- 성별 미확인 동물
SELECT * FROM pet WHERE sex IS NULL;
```

### 3. DATE 타입 활용
```sql
-- 나이 계산
SELECT name, TIMESTAMPDIFF(YEAR, birth, CURDATE()) AS age
FROM pet WHERE death IS NULL;

-- 다음 달 생일인 동물
SELECT name, birth FROM pet
WHERE MONTH(birth) = MONTH(DATE_ADD(CURDATE(), INTERVAL 1 MONTH));
```

---

## 학습 포인트

### 1. 기본 CRUD 연습

```sql
-- CREATE
INSERT INTO pet VALUES ('Puffball', 'Diane', 'hamster', 'f', '1999-03-30', NULL);

-- READ
SELECT * FROM pet WHERE species = 'dog';
SELECT owner, COUNT(*) FROM pet GROUP BY owner;

-- UPDATE
UPDATE pet SET death = '2020-01-01' WHERE name = 'Bowser';

-- DELETE
DELETE FROM pet WHERE name = 'Puffball';
```

### 2. 조인 연습 (논리적 조인)

```sql
-- 각 동물의 이벤트 조회
SELECT p.name, p.species, e.date, e.type, e.remark
FROM pet p
LEFT JOIN event e ON p.name = e.name
ORDER BY p.name, e.date;

-- 이벤트가 있는 동물만
SELECT DISTINCT p.*
FROM pet p
INNER JOIN event e ON p.name = e.name;
```

### 3. 집계 함수 연습

```sql
-- 종류별 동물 수
SELECT species, COUNT(*) as count
FROM pet
GROUP BY species;

-- 주인별 동물 수 (2마리 이상)
SELECT owner, COUNT(*) as count
FROM pet
GROUP BY owner
HAVING count >= 2;

-- 출산 기록 통계
SELECT name, COUNT(*) as litter_count
FROM event
WHERE type = 'litter'
GROUP BY name;
```

### 4. 서브쿼리 연습

```sql
-- 출산 기록이 있는 동물
SELECT * FROM pet
WHERE name IN (
  SELECT DISTINCT name FROM event WHERE type = 'litter'
);

-- 가장 나이 많은 동물
SELECT * FROM pet
WHERE birth = (SELECT MIN(birth) FROM pet WHERE death IS NULL);
```

---

## 설계 개선 제안 (실무 적용 시)

### 1. Primary Key 추가
```sql
CREATE TABLE pet (
  pet_id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(20) NOT NULL,
  ...
);

CREATE TABLE event (
  event_id INT AUTO_INCREMENT PRIMARY KEY,
  pet_id INT NOT NULL,
  ...
  FOREIGN KEY (pet_id) REFERENCES pet(pet_id)
);
```

### 2. 정규화
```sql
-- 종류 테이블 분리
CREATE TABLE species (
  species_id INT PRIMARY KEY,
  name VARCHAR(20) NOT NULL,
  category ENUM('mammal', 'bird', 'reptile', 'fish')
);

-- 주인 테이블 분리
CREATE TABLE owner (
  owner_id INT PRIMARY KEY,
  name VARCHAR(50) NOT NULL,
  email VARCHAR(100),
  phone VARCHAR(20)
);
```

### 3. 이벤트 유형 ENUM 또는 별도 테이블
```sql
CREATE TABLE event_type (
  type_id INT PRIMARY KEY,
  name VARCHAR(20) NOT NULL,
  description VARCHAR(100)
);
```

---

## 파일 구조

```
menagerie-db/
├── README.txt          # 원본 설치 가이드
├── README.md           # 이 문서
├── cr_pet_tbl.sql      # pet 테이블 DDL
├── cr_event_tbl.sql    # event 테이블 DDL
├── pet.txt             # pet 초기 데이터 (TSV)
├── event.txt           # event 초기 데이터 (TSV)
├── load_pet_tbl.sql    # pet 데이터 로드 스크립트
└── ins_puff_rec.sql    # Puffball 추가 INSERT
```

---

## 설치 방법

```bash
# 1. 데이터베이스 생성
mysql -u root -p -e "CREATE DATABASE menagerie;"

# 2. 테이블 생성
mysql menagerie < cr_pet_tbl.sql
mysql menagerie < cr_event_tbl.sql

# 3. 데이터 로드
mysql menagerie -e "LOAD DATA LOCAL INFILE 'pet.txt' INTO TABLE pet;"
mysql menagerie -e "LOAD DATA LOCAL INFILE 'event.txt' INTO TABLE event;"

# 또는 mysqlimport 사용
mysqlimport --local menagerie pet.txt
mysqlimport --local menagerie event.txt
```

---

## 참고 자료

- [MySQL Tutorial](https://dev.mysql.com/doc/refman/8.0/en/tutorial.html)
- [MySQL Sample Databases](https://dev.mysql.com/doc/index-other.html)

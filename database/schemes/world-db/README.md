# World Database Schema 분석

## 개요

**World**는 MySQL의 가장 기본적인 공식 샘플 데이터베이스로, 전 세계 국가, 도시, 언어 정보를 담고 있습니다. 단순한 구조로 SQL 기초 학습과 조인 연습에 적합합니다.

- **출처**: [MySQL Sample Database - World](https://dev.mysql.com/doc/index-other.html)
- **다운로드**: [MySQL World Database](https://downloads.mysql.com/docs/world-db.zip)
- **라이선스**: Public Domain
- **MySQL 버전**: 5.0+
- **테이블 수**: 3개
- **데이터**: 239개 국가, 4,079개 도시

---

## ERD (Entity Relationship Diagram)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                               country                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│ Code (PK, CHAR 3)  │ ISO 3166-1 Alpha-3 국가 코드 (예: KOR, USA, JPN)        │
│ Name               │ 국가명 (예: South Korea)                               │
│ Continent          │ 대륙 (ENUM: 7개 대륙)                                   │
│ Region             │ 지역 (예: Eastern Asia)                                │
│ SurfaceArea        │ 면적 (km²)                                             │
│ IndepYear          │ 독립 연도 (NULL: 미독립)                                │
│ Population         │ 인구                                                   │
│ LifeExpectancy     │ 기대 수명                                              │
│ GNP                │ 국민총생산                                             │
│ GNPOld             │ 이전 GNP                                               │
│ LocalName          │ 현지어 국가명                                          │
│ GovernmentForm     │ 정부 형태                                              │
│ HeadOfState        │ 국가 원수                                              │
│ Capital            │ 수도 도시 ID (city.ID 참조, FK 아님)                    │
│ Code2              │ ISO 3166-1 Alpha-2 코드 (예: KR, US, JP)               │
└──────────┬──────────────────────────────────────────────────────────────────┘
           │
           │ Code (PK)
           │
     ┌─────┴─────┐
     │           │
     ▼           ▼
┌─────────────────────┐     ┌─────────────────────────────────────────────────┐
│        city         │     │              countrylanguage                     │
├─────────────────────┤     ├─────────────────────────────────────────────────┤
│ ID (PK, AUTO_INC)   │     │ CountryCode (PK, FK) │ 국가 코드                 │
│ Name                │     │ Language (PK)        │ 언어명                    │
│ CountryCode (FK)    │────►│ IsOfficial           │ 공용어 여부 (T/F)          │
│ District            │     │ Percentage           │ 사용 비율 (%)              │
│ Population          │     └─────────────────────────────────────────────────┘
└─────────────────────┘
```

---

## 테이블 상세

### 1. `country` - 국가 정보

```sql
CREATE TABLE country (
  Code CHAR(3) NOT NULL DEFAULT '' PRIMARY KEY,  -- ISO 3166-1 Alpha-3
  Name CHAR(52) NOT NULL DEFAULT '',
  Continent ENUM('Asia','Europe','North America','Africa',
                 'Oceania','Antarctica','South America') NOT NULL DEFAULT 'Asia',
  Region CHAR(26) NOT NULL DEFAULT '',
  SurfaceArea DECIMAL(10,2) NOT NULL DEFAULT '0.00',
  IndepYear SMALLINT DEFAULT NULL,               -- NULL: 독립하지 않은 지역
  Population INT NOT NULL DEFAULT '0',
  LifeExpectancy DECIMAL(3,1) DEFAULT NULL,
  GNP DECIMAL(10,2) DEFAULT NULL,
  GNPOld DECIMAL(10,2) DEFAULT NULL,
  LocalName CHAR(45) NOT NULL DEFAULT '',
  GovernmentForm CHAR(45) NOT NULL DEFAULT '',
  HeadOfState CHAR(60) DEFAULT NULL,
  Capital INT DEFAULT NULL,                       -- city.ID 참조 (명시적 FK 없음)
  Code2 CHAR(2) NOT NULL DEFAULT ''              -- ISO 3166-1 Alpha-2
);
```

**설계 포인트**:
- **Natural Key**: ISO 표준 코드를 PK로 사용
- **ENUM**: 7개 대륙으로 제한
- **Capital**: city.ID를 참조하지만 명시적 FK 제약 없음 (순환 참조 방지)

### 2. `city` - 도시 정보

```sql
CREATE TABLE city (
  ID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  Name CHAR(35) NOT NULL DEFAULT '',
  CountryCode CHAR(3) NOT NULL DEFAULT '',
  District CHAR(20) NOT NULL DEFAULT '',         -- 행정구역 (도/시/주)
  Population INT NOT NULL DEFAULT '0',
  CONSTRAINT city_ibfk_1
    FOREIGN KEY (CountryCode) REFERENCES country (Code)
);
```

**샘플 데이터**:
| ID | Name | CountryCode | District | Population |
|----|------|-------------|----------|------------|
| 2331 | Seoul | KOR | Seoul | 9,981,619 |
| 2332 | Pusan | KOR | Pusan | 3,804,522 |
| 3793 | New York | USA | New York | 8,008,278 |
| 1532 | Tokyo | JPN | Tokyo-to | 7,980,230 |

### 3. `countrylanguage` - 국가별 언어

```sql
CREATE TABLE countrylanguage (
  CountryCode CHAR(3) NOT NULL DEFAULT '',
  Language CHAR(30) NOT NULL DEFAULT '',
  IsOfficial ENUM('T','F') NOT NULL DEFAULT 'F',  -- T: 공용어, F: 비공용어
  Percentage DECIMAL(4,1) NOT NULL DEFAULT '0.0',
  PRIMARY KEY (CountryCode, Language),            -- 복합 PK
  CONSTRAINT countryLanguage_ibfk_1
    FOREIGN KEY (CountryCode) REFERENCES country (Code)
);
```

**샘플 데이터**:
| CountryCode | Language | IsOfficial | Percentage |
|-------------|----------|------------|------------|
| KOR | Korean | T | 99.9 |
| USA | English | T | 86.2 |
| USA | Spanish | F | 7.5 |
| CHE | German | T | 63.6 |
| CHE | French | T | 19.2 |

---

## 스키마 설계 특징

### 1. Natural Key vs Surrogate Key
```
country: Code (Natural Key) - ISO 표준 코드 사용
city: ID (Surrogate Key) - AUTO_INCREMENT 사용
```
- Natural Key: 비즈니스 의미가 있는 키 (변경 가능성 낮음)
- Surrogate Key: 시스템이 생성하는 키 (변경 없음)

### 2. 복합 Primary Key
```sql
PRIMARY KEY (CountryCode, Language)
```
- 한 국가에 여러 언어 존재 가능
- 언어-국가 조합은 유일해야 함

### 3. ENUM 타입 활용
```sql
Continent ENUM('Asia','Europe',...) -- 제한된 값만 허용
IsOfficial ENUM('T','F')            -- Boolean 대안
```

### 4. NULL 허용 설계
```sql
IndepYear SMALLINT DEFAULT NULL    -- 독립하지 않은 지역
LifeExpectancy DECIMAL DEFAULT NULL -- 데이터 없음
HeadOfState CHAR(60) DEFAULT NULL   -- 공화국 등
```

---

## 학습 포인트

### 1. 기본 조인
```sql
-- 국가와 수도 조회
SELECT co.Name AS Country, ci.Name AS Capital
FROM country co
JOIN city ci ON co.Capital = ci.ID
ORDER BY co.Name;

-- 결과:
-- Afghanistan - Kabul
-- South Korea - Seoul
-- United States - Washington
```

### 2. 집계 함수
```sql
-- 대륙별 국가 수, 총 인구
SELECT Continent,
       COUNT(*) AS Countries,
       SUM(Population) AS TotalPopulation
FROM country
GROUP BY Continent
ORDER BY TotalPopulation DESC;

-- 인구 1억 이상 국가
SELECT Name, Population
FROM country
WHERE Population >= 100000000
ORDER BY Population DESC;
```

### 3. 서브쿼리
```sql
-- 평균 인구보다 많은 국가
SELECT Name, Population
FROM country
WHERE Population > (SELECT AVG(Population) FROM country);

-- 공용어가 2개 이상인 국가
SELECT co.Name, COUNT(*) AS OfficialLanguages
FROM country co
JOIN countrylanguage cl ON co.Code = cl.CountryCode
WHERE cl.IsOfficial = 'T'
GROUP BY co.Code
HAVING COUNT(*) >= 2;
```

### 4. 다중 테이블 조인
```sql
-- 국가, 수도, 공용어 조회
SELECT co.Name AS Country,
       ci.Name AS Capital,
       cl.Language,
       cl.Percentage
FROM country co
JOIN city ci ON co.Capital = ci.ID
JOIN countrylanguage cl ON co.Code = cl.CountryCode
WHERE cl.IsOfficial = 'T'
ORDER BY co.Name;
```

### 5. 분석 쿼리
```sql
-- 대륙별 기대수명 평균
SELECT Continent,
       ROUND(AVG(LifeExpectancy), 1) AS AvgLifeExpectancy
FROM country
WHERE LifeExpectancy IS NOT NULL
GROUP BY Continent
ORDER BY AvgLifeExpectancy DESC;

-- 인구 밀도 상위 10개국
SELECT Name,
       Population,
       SurfaceArea,
       ROUND(Population / SurfaceArea, 2) AS Density
FROM country
WHERE SurfaceArea > 0
ORDER BY Density DESC
LIMIT 10;
```

---

## 데이터 통계

| 항목 | 수량 |
|------|------|
| 국가 수 | 239 |
| 도시 수 | 4,079 |
| 언어 수 | 984 (중복 포함) |
| 대륙 수 | 7 |
| 총 인구 (합계) | ~6,078,749,450 |

### 대륙별 분포
| 대륙 | 국가 수 | 도시 수 |
|------|---------|---------|
| Africa | 58 | 366 |
| Asia | 51 | 1,766 |
| Europe | 46 | 841 |
| North America | 37 | 581 |
| South America | 14 | 470 |
| Oceania | 28 | 55 |
| Antarctica | 5 | 0 |

---

## 파일 구조

```
world-db/
├── README.md       # 이 문서
└── world.sql       # DDL + DML (전체 스키마 및 데이터)
```

---

## 설치 방법

```bash
# 전체 설치 (스키마 + 데이터)
mysql -u root -p < world.sql

# 확인
mysql -u root -p -e "USE world; SELECT COUNT(*) FROM country;"
# 결과: 239

mysql -u root -p -e "USE world; SELECT COUNT(*) FROM city;"
# 결과: 4079
```

---

## 참고 자료

- [MySQL Example Databases](https://dev.mysql.com/doc/index-other.html)
- [ISO 3166-1 Country Codes](https://en.wikipedia.org/wiki/ISO_3166-1)

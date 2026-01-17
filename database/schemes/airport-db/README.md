# Airport Database Schema 분석

## 개요

**Airport DB**는 MySQL용 샘플 데이터베이스로, 항공 예약 시스템을 모델링한 스키마입니다. Oracle Cloud Infrastructure (OCI)의 MySQL DB System과 HeatWave에서 테스트용으로 사용됩니다.

- **출처**: [MySQL Sample Database](https://dev.mysql.com/doc/airportdb/en/)
- **다운로드**: [Oracle Cloud Infrastructure - Airport DB](https://downloads.mysql.com/docs/airport-db.zip)
- **라이선스**: CC BY 4.0 (Stefan Pröll, Eva Zangerle, Wolfgang Gassler)
- **MySQL 버전**: 8.0.25+

---

## ERD (Entity Relationship Diagram)

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│  airplane_type  │     │     airline     │     │     airport     │
├─────────────────┤     ├─────────────────┤     ├─────────────────┤
│ PK type_id      │◄────┤ FK base_airport │────►│ PK airport_id   │
│    identifier   │     │ PK airline_id   │     │    iata (UQ)    │
│    description  │     │    iata (UQ)    │     │    icao (UQ)    │
│   [FULLTEXT]    │     │    airlinename  │     │    name         │
└────────┬────────┘     └────────┬────────┘     └────────┬────────┘
         │                       │                       │
         │              ┌────────┴────────┐              │
         │              │                 │              │
         ▼              ▼                 │              ▼
┌─────────────────┐  ┌─────────────────┐  │   ┌─────────────────┐
│    airplane     │  │ flightschedule  │  │   │   airport_geo   │
├─────────────────┤  ├─────────────────┤  │   ├─────────────────┤
│ PK airplane_id  │  │ PK flightno     │  │   │ PK airport_id   │
│ FK type_id      │  │ FK from         │──┼──►│    name         │
│ FK airline_id   │  │ FK to           │──┤   │    city         │
│    capacity     │  │ FK airline_id   │  │   │    country      │
└────────┬────────┘  │    departure    │  │   │    latitude     │
         │           │    arrival      │  │   │    longitude    │
         │           │    monday..sun  │  │   │    geolocation  │
         │           └────────┬────────┘  │   │   [SPATIAL]     │
         │                    │           │   └─────────────────┘
         │                    │           │
         │                    ▼           │   ┌─────────────────┐
         │           ┌─────────────────┐  │   │airport_reachable│
         │           │     flight      │  │   ├─────────────────┤
         │           ├─────────────────┤  │   │ PK airport_id   │
         └──────────►│ PK flight_id    │  │   │    hops         │
                     │ FK flightno     │──┘   └─────────────────┘
                     │ FK from         │
                     │ FK to           │
                     │ FK airline_id   │
                     │ FK airplane_id  │
                     │    departure    │
                     │    arrival      │
                     └────────┬────────┘
                              │
         ┌────────────────────┼────────────────────┐
         │                    │                    │
         ▼                    ▼                    ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│   flight_log    │  │     booking     │  │   weatherdata   │
├─────────────────┤  ├─────────────────┤  ├─────────────────┤
│ PK flight_log_id│  │ PK booking_id   │  │ PK (log_date,   │
│ FK flight_id    │  │ FK flight_id    │  │    time,station)│
│    log_date     │  │ FK passenger_id │  │    temp         │
│    user         │  │    seat (UQ)    │  │    humidity     │
│    *_old/*_new  │  │    price        │  │    airpressure  │
│    comment      │  └────────┬────────┘  │    wind         │
└─────────────────┘           │           │    weather      │
                              │           │    winddirection│
                              ▼           └─────────────────┘
                     ┌─────────────────┐
                     │    passenger    │
                     ├─────────────────┤
                     │ PK passenger_id │
                     │    passportno   │
                     │    firstname    │
                     │    lastname     │
                     └────────┬────────┘
                              │
                              ▼
                     ┌─────────────────┐
                     │passengerdetails │
                     ├─────────────────┤
                     │ PK passenger_id │
                     │    birthdate    │
                     │    sex          │
                     │    street       │
                     │    city/zip     │
                     │    country      │
                     │    emailaddress │
                     │    telephoneno  │
                     └─────────────────┘

┌─────────────────┐
│    employee     │
├─────────────────┤
│ PK employee_id  │
│    firstname    │
│    lastname     │
│    birthdate    │
│    sex          │
│    address      │
│    emailaddress │
│    telephoneno  │
│    salary       │
│    department   │
│    username(UQ) │
│    password     │
└─────────────────┘
```

---

## 테이블 상세 설명

### 1. 마스터 테이블 (Master Tables)

#### `airport` - 공항 정보
| 컬럼 | 타입 | 설명 |
|------|------|------|
| airport_id | SMALLINT (PK) | 공항 ID (AUTO_INCREMENT) |
| iata | CHAR(3) | IATA 코드 (e.g., ICN, JFK) |
| icao | CHAR(4) | ICAO 코드 (UNIQUE) |
| name | VARCHAR(50) | 공항명 |

**설계 포인트**:
- `icao`에 UNIQUE 제약으로 데이터 무결성 보장
- `iata`와 `name`에 INDEX 추가로 검색 성능 최적화

#### `airport_geo` - 공항 지리 정보
| 컬럼 | 타입 | 설명 |
|------|------|------|
| airport_id | SMALLINT (PK, FK) | 공항 ID |
| latitude | DECIMAL(11,8) | 위도 |
| longitude | DECIMAL(11,8) | 경도 |
| geolocation | POINT | 공간 데이터 |

**설계 포인트**:
- **1:1 관계**: `airport`와 분리하여 정규화
- **SPATIAL INDEX**: 지리적 쿼리 (근처 공항 검색 등) 최적화
- MySQL의 GIS 기능 활용

#### `airline` - 항공사 정보
| 컬럼 | 타입 | 설명 |
|------|------|------|
| airline_id | SMALLINT (PK) | 항공사 ID |
| iata | CHAR(2) | IATA 항공사 코드 (UNIQUE) |
| airlinename | VARCHAR(30) | 항공사명 |
| base_airport | SMALLINT (FK) | 본사 공항 |

#### `airplane_type` - 항공기 기종
| 컬럼 | 타입 | 설명 |
|------|------|------|
| type_id | INT (PK) | 기종 ID |
| identifier | VARCHAR(50) | 기종 식별자 (e.g., B737) |
| description | TEXT | 기종 설명 |

**설계 포인트**:
- **FULLTEXT INDEX**: `identifier`와 `description`에 전문 검색 인덱스
- 항공기 기종 검색 시 자연어 검색 가능

---

### 2. 운영 테이블 (Operational Tables)

#### `airplane` - 항공기
| 컬럼 | 타입 | 설명 |
|------|------|------|
| airplane_id | INT (PK) | 항공기 ID |
| capacity | MEDIUMINT UNSIGNED | 좌석 수 |
| type_id | INT (FK) | 기종 ID |
| airline_id | INT (FK) | 소속 항공사 |

#### `flightschedule` - 운항 스케줄 (정기편)
| 컬럼 | 타입 | 설명 |
|------|------|------|
| flightno | CHAR(8) (PK) | 편명 (e.g., KE001) |
| from/to | SMALLINT (FK) | 출발/도착 공항 |
| departure/arrival | TIME | 출발/도착 시간 |
| monday..sunday | TINYINT(1) | 요일별 운항 여부 |

**설계 포인트**:
- **Boolean 컬럼**: 요일별 운항 여부를 7개의 개별 컬럼으로 관리
- `flightno`를 Primary Key로 사용 (Natural Key)

#### `flight` - 실제 운항 (인스턴스)
| 컬럼 | 타입 | 설명 |
|------|------|------|
| flight_id | INT (PK) | 운항 ID |
| flightno | CHAR(8) (FK) | 스케줄 편명 참조 |
| departure/arrival | DATETIME | 실제 출발/도착 시각 |

**설계 포인트**:
- `flightschedule`은 템플릿, `flight`은 실제 인스턴스
- 다수의 INDEX로 검색 최적화 (from, to, departure, arrival, airline, airplane)

---

### 3. 고객 테이블 (Customer Tables)

#### `passenger` - 승객 기본 정보
| 컬럼 | 타입 | 설명 |
|------|------|------|
| passenger_id | INT (PK) | 승객 ID |
| passportno | CHAR(9) (UNIQUE) | 여권번호 |
| firstname/lastname | VARCHAR(100) | 이름 |

#### `passengerdetails` - 승객 상세 정보
| 컬럼 | 타입 | 설명 |
|------|------|------|
| passenger_id | INT (PK, FK) | 승객 ID |
| birthdate | DATE | 생년월일 |
| sex | CHAR(1) | 성별 |
| address (street, city, zip, country) | VARCHAR | 주소 |
| emailaddress | VARCHAR(120) | 이메일 |
| telephoneno | VARCHAR(30) | 전화번호 |

**설계 포인트**:
- **테이블 분리 (Vertical Partitioning)**:
  - `passenger`: 필수 정보 (항상 필요)
  - `passengerdetails`: 선택적 정보 (필요 시만 조회)
- **ON DELETE CASCADE**: 승객 삭제 시 상세정보도 자동 삭제

---

### 4. 트랜잭션 테이블

#### `booking` - 예약
| 컬럼 | 타입 | 설명 |
|------|------|------|
| booking_id | INT (PK) | 예약 ID |
| flight_id | INT (FK) | 운항 ID |
| seat | CHAR(4) | 좌석번호 |
| passenger_id | INT (FK) | 승객 ID |
| price | DECIMAL(10,2) | 가격 |

**설계 포인트**:
- **복합 UNIQUE 제약**: `(flight_id, seat)` - 동일 편에 중복 좌석 방지
- 약 5,500만 건의 데이터 (AUTO_INCREMENT 55099799)

---

### 5. 감사/로그 테이블

#### `flight_log` - 운항 변경 이력
| 컬럼 | 타입 | 설명 |
|------|------|------|
| flight_log_id | INT UNSIGNED (PK) | 로그 ID |
| log_date | DATETIME | 변경 시각 |
| user | VARCHAR(100) | 변경자 |
| *_old / *_new | 각 타입 | 변경 전/후 값 |
| comment | VARCHAR(200) | 변경 사유 |

**설계 포인트**:
- **Audit Trail 패턴**: 모든 변경 사항을 old/new 쌍으로 기록
- 운항 정보 변경에 대한 완전한 추적 가능

---

### 6. 부가 테이블

#### `employee` - 직원 정보
| 컬럼 | 타입 | 설명 |
|------|------|------|
| employee_id | INT (PK) | 직원 ID |
| department | ENUM | 부서 (Marketing, Buchhaltung, Management, Logistik, Flugfeld) |
| username | VARCHAR(20) (UNIQUE) | 로그인 ID |
| password | CHAR(32) | 비밀번호 (MD5 해시) |

**설계 포인트**:
- **ENUM 타입**: 부서를 제한된 값으로 관리 (독일어 사용)
- 비밀번호 32자 = MD5 해시 저장 (보안상 권장하지 않음, 예제용)

#### `weatherdata` - 날씨 데이터
| 컬럼 | 타입 | 설명 |
|------|------|------|
| (log_date, time, station) | PK | 복합 키 |
| temp, humidity, airpressure, wind | DECIMAL | 기상 수치 |
| weather | ENUM | 날씨 상태 (독일어) |
| winddirection | SMALLINT | 풍향 (각도) |

**설계 포인트**:
- **복합 Primary Key**: (날짜, 시간, 관측소) 조합
- 시계열 데이터 저장 패턴

---

## 스키마 설계 특징 및 배울 점

### 1. 정규화와 비정규화의 균형
- **정규화**: `airport`와 `airport_geo` 분리 (1:1 관계)
- **비정규화**: `flightschedule`의 요일별 Boolean 컬럼

### 2. 적절한 인덱스 전략
- **B-Tree INDEX**: FK 컬럼, 검색 빈도 높은 컬럼
- **UNIQUE INDEX**: 비즈니스 키 (IATA/ICAO 코드, 여권번호)
- **FULLTEXT INDEX**: 텍스트 검색용 (`airplane_type.description`)
- **SPATIAL INDEX**: 지리 데이터 (`airport_geo.geolocation`)

### 3. 데이터 타입 선택
- `SMALLINT` vs `INT`: 예상 데이터 크기에 따른 선택
- `DECIMAL(10,2)`: 금액은 정확한 소수점 타입
- `CHAR` vs `VARCHAR`: 고정 길이(IATA/ICAO 코드)는 CHAR

### 4. 감사 추적 (Audit Trail)
- `flight_log` 테이블로 변경 이력 관리
- old/new 값 쌍으로 저장하는 패턴

### 5. 복합 키와 자연 키
- `flightschedule`: 편명을 자연 키(Natural Key)로 사용
- `weatherdata`: 복합 키로 시계열 데이터 관리

---

## 사용 예시

```sql
-- 인천공항에서 출발하는 모든 항공편 조회
SELECT f.flightno, f.departure, a.name as destination
FROM flight f
JOIN airport a ON f.to = a.airport_id
JOIN airport dep ON f.from = dep.airport_id
WHERE dep.iata = 'ICN';

-- 특정 위치에서 가장 가까운 공항 찾기 (SPATIAL 쿼리)
SELECT a.name, ST_Distance_Sphere(
  geo.geolocation,
  ST_GeomFromText('POINT(126.4406 37.4602)')
) as distance
FROM airport a
JOIN airport_geo geo ON a.airport_id = geo.airport_id
ORDER BY distance
LIMIT 5;

-- 승객의 전체 예약 내역 조회
SELECT p.firstname, p.lastname, b.seat, b.price,
       f.flightno, f.departure
FROM passenger p
JOIN booking b ON p.passenger_id = b.passenger_id
JOIN flight f ON b.flight_id = f.flight_id
WHERE p.passportno = 'M12345678';
```

---

## 파일 구조

```
airport-db/
├── README.md               # 스키마 분석 문서
├── README.txt              # 원본 설명
├── @.sql, @.post.sql       # 스키마 초기화/후처리
├── @.json, @.done.json     # 메타데이터
├── @.manifest.json         # 매니페스트
├── airportdb.sql           # 데이터베이스 생성
└── airportdb@{table}.sql   # 각 테이블 DDL
```

> **참고**: 원본 데이터 파일(`.tsv.zst`, `.idx`, `.json`)은 용량 문제로 제외되었습니다.
> 전체 데이터가 필요한 경우 [다운로드 링크](https://downloads.mysql.com/docs/airport-db.zip)에서 받을 수 있습니다.

#create database
CREATE DATABASE IF NOT EXISTS study_db default CHARACTER SET UTF8;
SHOW DATABASES;

USE study_db;

# mysql에서 지원하는 coordinate 시스템 
SELECT `srs_name`, `srs_id`
FROM INFORMATION_SCHEMA.ST_SPATIAL_REFERENCE_SYSTEMS;


SET @point_p := POINT(3, 4);
SELECT @point_p;
# BLOB은 geometries를 위한 MYSQL의 기본 저장 포멧

# 거리 계산
SET @point_o := POINT(0, 0); # origin
SET @point_p := POINT(3, 4);
SELECT ST_Distance(@point_o, @point_p);

# WKT로 Geo 생성
SET @point_o := ST_GeomFromText('POINT(0 0)');
SET @point_p := ST_GeomFromText('POINT(3 4)');

# 거리 계산 값 & 현재 사용하는 SRID 값
SELECT ST_DISTANCE(@point_o, @point_p) AS distance,
       ST_SRID(@point_o)               AS _srid;

# convert geometric shape of the type point -> WKT, WKB, JSON format
SET @point_o := ST_GeomFromText('POINT(0 0)');

SELECT ST_AsText(@point_o)    AS 'wkt_value',
       ST_AsBinary(@point_o)  AS 'wkb_value',
       ST_AsGeoJson(@point_o) AS 'geo_json_value';

# convert WKT, WKB, Json -> geometric point
SET @point_o := ST_GeomFromText('POINT(0 0)'); # geometry from WKT
SET @point_o_wkb := ST_AsBinary(@point_o); #wkb of a geometry
SET @point_o_from_wkb := ST_GeomFromWKB(@point_o_wkb); # geometry from WKB
SET @point_o_from_geoJSON := ST_GeomFromGeoJSON('{"type": "Point", "coordinates": [0.0, 0.0]}'); # geometry from geoJSON

SELECT ST_ASTEXT(@point_o_from_wkb)     AS 'point_o_from_wkb',
       ST_ASTEXT(@point_o_from_geoJSON) AS 'point_o_from_geoJSON';

# x, y을 값
SET @point_o := ST_GeomFromText('POINT(3 4)'); # geometry from WKT
SELECT ST_X(@point_o) AS 'point_o_x',
       ST_Y(@point_o) AS 'point_o_y';

# 2d cartesion coordinates system 테이블 생성
CREATE TABLE IF NOT EXISTS `locations_flat`
(
    `id`       INT   NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `name`     VARCHAR(100),
    `position` POINT NOT NULL SRID 0
);
SHOW COLUMNS FROM `locations_flat`;

# point 샘플 데이터 insert 하기
INSERT INTO `locations_flat`(`name`, `position`)
VALUES ('point_1', ST_GeomFromText('POINT( 1 1 )', 0)),
       ('point_2', ST_GeomFromText('POINT( 2 2 )', 0)),
       ('point_3', ST_GeomFromText('POINT( 3 3 )', 0));

#
SELECT *,
       ST_ASTEXT(`position`) AS `pos_wkt`
FROM `locations_flat`
LIMIT 10;

# 거리가 100 이하인 장소 찾기
SET @user_location = ST_GeomFromText('POINT(0 0)');
SELECT *,
       ST_AsText(`position`)                   AS `pos_wkt`,
       ST_Distance(`position`, @user_location) AS `distance`
FROM `locations_flat`
WHERE ST_Distance(`position`, @user_location) <= 100;

# (0,0) 100 의 원형안에 있는 장소인지
# ST_Buffer는 SRID 0에서만 사용할 수 있음
SET @user_location = ST_GeomFromText('POINT(0 0)');
SET @area_to_search = ST_Buffer(@user_location, 100);
SELECT *,
       ST_AsText(`position`)                   AS `pos_wkt`,
       ST_Distance(`position`, @user_location) AS `distance`
FROM `locations_flat`
WHERE ST_Within(`position`, @area_to_search) = 1;

# SPATIAL INDEX 추가
ALTER TABLE `locations_flat`
    ADD SPATIAL INDEX (`position`);
SHOW INDEXES FROM `locations_flat`;

# 성능 체크
# full-table scan을 해서 1000 rows가 설치가 되었음
SET @user_location = ST_GeomFromText('POINT(0 0)');
SET @area_to_search = ST_Buffer(@user_location, 100);
EXPLAIN
SELECT *,
       ST_AsText(`position`)                   AS `pos_wkt`,
       ST_Distance(`position`, @user_location) AS `distance`
FROM `locations_flat`
         IGNORE INDEX (`position`)
WHERE ST_Within(`position`, @area_to_search) = 1;

# db 생성 - geographic coordinate system 4326
# POINT(lat, long)
CREATE TABLE IF NOT EXISTS `locations_earth`
(
    `id`       INT   NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `name`     VARCHAR(100),
    `position` POINT NOT NULL SRID 4326,
    SPATIAL INDEX (`position`)
);
SHOW COLUMNS FROM `locations_earth`;

# 인도 IndiaGate -> India Temlate 거리 => 7210 meters
SET @lotus_temple := ST_GeomFromText('POINT(28.553298 77.259221)', 4326, 'axis-order=lat-long');
SET @india_gate := ST_GeomFromText('POINT(28.612849 77.229883)', 4326);
SELECT ST_Latitude(@lotus_temple)                     AS `lat_lotus_temple`,
       ST_Longitude(@lotus_temple)                    AS `long_lotus_temple`,
       ST_Latitude(@india_gate)                       AS `lat_india_gate`,
       ST_Longitude(@india_gate)                      AS `long_india_gate`,
       ST_Distance_Sphere(@lotus_temple, @india_gate) AS `distance`;

# spherical coordinate system에서는 circular geometry를 생성할 수 없어 대신 polygon 을 생성하여 검색을 시도함
# - query가 index를 사용하는지 확인
SET @poly_o = ST_GeomFromText('POLYGON(( 30 40, 40 50, 30 60, 20 50, 30 40 ))', 4326);
SET @user_location = ST_GeomFromText('POINT(30 50)', 4326);

EXPLAIN
SELECT *,
       ST_AsText(`position`)                          AS `pos_wkt`,
       ST_Distance_Sphere(`position`, @user_location) AS `distance`
FROM `locations_earth`
WHERE ST_Within(`position`, @poly_o);

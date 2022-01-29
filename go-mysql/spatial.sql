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

#
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

INSERT INTO `locations_flat`(`name`, `position`)
VALUES ('point_4', ST_GeomFromText('POINT( 4 4 )', 0));


SELECT * from `study_db`.`locations_flat`;

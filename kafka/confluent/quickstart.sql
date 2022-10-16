
-- 1.create stream
CREATE STREAM pageviews_stream
  WITH (KAFKA_TOPIC='pageviews', VALUE_FORMAT='AVRO');

-- 2. create table
CREATE TABLE users_table (id VARCHAR PRIMARY KEY)
    WITH (KAFKA_TOPIC='users', VALUE_FORMAT='AVRO');

-- create transient query
SELECT pageid FROM pageviews_stream EMIT CHANGES LIMIT 3;

-- 3.join pageviews stream and users table
CREATE STREAM user_pageviews
  AS SELECT users_table.id AS userid, pageid, regionid, gender
     FROM pageviews_stream
              LEFT JOIN users_table ON pageviews_stream.userid = users_table.id
         EMIT CHANGES;

-- 4.filter stream by region field
CREATE STREAM pageviews_region_like_89
  WITH (KAFKA_TOPIC='pageviews_filtered_r8_r9', VALUE_FORMAT='AVRO')
    AS SELECT * FROM user_pageviews
       WHERE regionid LIKE '%_8' OR regionid LIKE '%_9'
           EMIT CHANGES;

-- 5.create windowed view
CREATE TABLE pageviews_per_region_89 WITH (KEY_FORMAT='JSON')
AS SELECT userid, gender, regionid, COUNT(*) AS numusers
   FROM pageviews_region_like_89
            WINDOW TUMBLING (SIZE 30 SECOND)
   GROUP BY userid, gender, regionid
   HAVING COUNT(*) > 1
       EMIT CHANGES;

-- 6.lookup table data by pull query
SELECT * FROM PAGEVIEWS_PER_REGION_89
WHERE userid = 'User_1' AND gender='FEMALE' AND regionid='Region_9';

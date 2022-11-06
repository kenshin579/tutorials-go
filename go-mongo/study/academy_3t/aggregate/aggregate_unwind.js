use test;
//1.하나의 도큐먼트를 sizes의 배열 요소로 분리함
db.unwind1.insertMany([{ "_id" : 1, "item" : "ABC1", sizes: [ "S", "M", "L"] }]);
db.unwind1.find({});
db.unwind1.aggregate([{$unwind: "$sizes"}]);
db.unwind1.drop();

//2.sizes 배열의 요소의 순번 표시하기
db.unwind2.insertMany([{ "_id" : 1, "item" : "ABC", "sizes": [ "S", "M", "L"] },
{ "_id" : 2, "item" : "EFG", "sizes" : [ ] },
{ "_id" : 3, "item" : "IJK", "sizes": "M" },
{ "_id" : 4, "item" : "LMN" },
{ "_id" : 5, "item" : "XYZ", "sizes" : null }]);
db.unwind2.find({});
db.unwind2.aggregate([{$unwind: { path: "$sizes", includeArrayIndex: "arrayIndex"}}])

//3.sizes의 원소가 없거나 null이어도 출력하라
db.unwind2.aggregate([{$unwind: { path: "$sizes", preserveNullAndEmptyArrays: true}}])
db.unwind2.drop();



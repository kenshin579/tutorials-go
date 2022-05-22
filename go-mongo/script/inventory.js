// https://www.mongodb.com/docs/manual/tutorial/query-documents/

//login으로 들어가려
// use admin
// db.auth('mongoadmin', 'mongopassword')

//show database
show dbs

//create database
use go-mongo

//drop database
//db.dropDtabase()

//create collection
db.createCollection("inventory")
show collections

//////////////////////////////////////////////////////////////////////////////
// Insert Documents
//////////////////////////////////////////////////////////////////////////////
db.inventory.insertOne(
    {item: "canvas", qty: 100, tags: ["cotton"], size: {h: 28, w: 35.5, uom: "cm"}}
)

db.inventory.insertMany([
    {item: "journal", qty: 25, tags: ["blank", "red"], size: {h: 14, w: 21, uom: "cm"}},
    {item: "mat", qty: 85, tags: ["gray"], size: {h: 27.9, w: 35.5, uom: "cm"}},
    {item: "mousepad", qty: 25, tags: ["gel", "blue"], size: {h: 19, w: 22.85, uom: "cm"}}
])


//////////////////////////////////////////////////////////////////////////////
//Query Documents
//////////////////////////////////////////////////////////////////////////////
db.inventory.insertMany([
    {item: "journal", qty: 25, size: {h: 14, w: 21, uom: "cm"}, status: "A"},
    {item: "notebook", qty: 50, size: {h: 8.5, w: 11, uom: "in"}, status: "A"},
    {item: "paper", qty: 100, size: {h: 8.5, w: 11, uom: "in"}, status: "D"},
    {item: "planner", qty: 75, size: {h: 22.85, w: 30, uom: "cm"}, status: "D"},
    {item: "postcard", qty: 45, size: {h: 10, w: 15.25, uom: "cm"}, status: "A"}
])

/*
Select All Documents in a Collection
SELECT * FROM inventory
 */
db.inventory.find({})


/*
Specify Equality Condition
db.inventory.find( { status: "D" } )
SELECT * FROM inventory WHERE status = "D"
 */
db.inventory.find({status: "D"})


/*
Specify Conditions Using Query Operators
SELECT * FROM inventory WHERE status in ("A", "D")
 */
db.inventory.find({status: {$in: ["A", "D"]}})

/*
Specify AND Conditions

SELECT * FROM inventory WHERE status = "A" AND qty < 30
 */
db.inventory.find({status: "A", qty: {$lt: 30}})

/*
Specify OR Conditions
SELECT * FROM inventory WHERE status = "A" OR qty < 30
 */
db.inventory.find({$or: [{status: "A"}, {qty: {$lt: 30}}]})

/*
Specify AND as well as OR Conditions
SELECT * FROM inventory WHERE status = "A" AND ( qty < 30 OR item LIKE "p%")
 */
db.inventory.find({
    status: "A",
    $or: [{qty: {$lt: 30}}, {item: /^p/}]
})

// 1 Query on Embedded/Nested Documents
db.inventory.insertMany([
    {item: "journal", qty: 25, size: {h: 14, w: 21, uom: "cm"}, status: "A"},
    {item: "notebook", qty: 50, size: {h: 8.5, w: 11, uom: "in"}, status: "A"},
    {item: "paper", qty: 100, size: {h: 8.5, w: 11, uom: "in"}, status: "D"},
    {item: "planner", qty: 75, size: {h: 22.85, w: 30, uom: "cm"}, status: "D"},
    {item: "postcard", qty: 45, size: {h: 10, w: 15.25, uom: "cm"}, status: "A"}
]);

/*
1.1 Match an Embedded/Nested Document
{ <field>: <value> }
 */
db.inventory.find({size: {h: 14, w: 21, uom: "cm"}})

/*
1.2 Query on Nested Field
 */
//1.2.1 Specify Equality Match on a Nested Field (dot)
db.inventory.find({"size.uom": "in"})

/*
1.2.2 Specify Match using Query Operator (ex. lt)
{ <field1>: { <operator1>: <value1> }, ... }
 */
db.inventory.find({"size.h": {$lt: 15}})

//1.2.3 Specify AND Condition
db.inventory.find({"size.h": {$lt: 15}, "size.uom": "in", status: "D"})

/*
2. Query an Array
 */
db.inventory.insertMany([
    {item: "journal", qty: 25, tags: ["blank", "red"], dim_cm: [14, 21]},
    {item: "notebook", qty: 50, tags: ["red", "blank"], dim_cm: [14, 21]},
    {item: "paper", qty: 100, tags: ["red", "blank", "plain"], dim_cm: [14, 21]},
    {item: "planner", qty: 75, tags: ["blank", "red"], dim_cm: [22.85, 30]},
    {item: "postcard", qty: 45, tags: ["blue"], dim_cm: [10, 15.25]}
]);

/*
2.1 Match an Array
an array with exactly two elements, "red" and "blank", in the specified order
 */
db.inventory.find({tags: ["red", "blank"]})

//find an array that contains both the elements "red" and "blank", without regard to order
db.inventory.find({tags: {$all: ["red", "blank"]}})

/*
2.2 Query an Array for an Element
{ <field>: <value> }
 */
// all documents where tags is an array that contains the string "red" as one of its elements:
db.inventory.find({tags: "red"})

/*
{ <array field>: { <operator1>: <value1>, ... } }
all documents where the array dim_cm contains at least one element whose value is greater than 25.
 */
db.inventory.find({dim_cm: {$gt: 25}})


//2.3 Query an Array by Array Length
db.inventory.find({"tags": {$size: 3}})

/*
3. Project Fields to Return from Query
 */
db.inventory.insertMany([
    {item: "journal", status: "A", size: {h: 14, w: 21, uom: "cm"}, instock: [{warehouse: "A", qty: 5}]},
    {item: "notebook", status: "A", size: {h: 8.5, w: 11, uom: "in"}, instock: [{warehouse: "C", qty: 5}]},
    {item: "paper", status: "D", size: {h: 8.5, w: 11, uom: "in"}, instock: [{warehouse: "A", qty: 60}]},
    {item: "planner", status: "D", size: {h: 22.85, w: 30, uom: "cm"}, instock: [{warehouse: "A", qty: 40}]},
    {
        item: "postcard",
        status: "A",
        size: {h: 10, w: 15.25, uom: "cm"},
        instock: [{warehouse: "B", qty: 15}, {warehouse: "C", qty: 35}]
    }
]);

/*
3.1 Return All Fields in Matching Documents
SELECT * from inventory WHERE status = "A"
 */
db.inventory.find({status: "A"})

/*
3.2 Return the Specified Fields and the _id Field Only
projection can explicitly include several fields by setting the <field> to 1 in the projection document.
by default, the _id fields return in the matching documents.
SELECT _id, item, status from inventory WHERE status = "A"
 */
db.inventory.find({status: "A"}, {item: 1, status: 1})

/*
3.3 Suppress _id Field
You can remove the _id field from the results by setting it to 0 in the projection,

SELECT item, status from inventory WHERE status = "A"
 */
db.inventory.find({status: "A"}, {item: 1, status: 1, _id: 0})

/*
3.4 Return All But the Excluded Fields
 */
db.inventory.find({status: "A"}, {status: 0, instock: 0})

/*
3.5 Project Specific Array Elements in the Returned Array
operator : $elemMatch, $slice, and $.
return the last element in the instock array:
 */
db.inventory.find({status: "A"}, {item: 1, status: 1, instock: {$slice: -1}})

/*
4. Query for Null or Missing Fields
 */
db.inventory.insertMany([
    {_id: 1, item: null},
    {_id: 2}
])

//4.1 Equality Filter
db.inventory.find({item: null})

/*
4.2 Type Check
item : { $type: 10 } } query matches only documents that contain the item field whose value is null
i.e. the value of the item field is of BSON Type Null (type number 10) :
 */
db.inventory.find({item: {$type: 10}})

/*
4.3 Existence Check

 */
db.inventory.find({item: {$exists: false}})

//////////////////////////////////////////////////////////////////////////////
//Update Documents
//////////////////////////////////////////////////////////////////////////////

//////////////////////////////////////////////////////////////////////////////
//Delete Documents
//////////////////////////////////////////////////////////////////////////////


//////////////////////////////////////////////////////////////////////////////
//Aggregation Operations
//////////////////////////////////////////////////////////////////////////////


//////////////////////////////////////////////////////////////////////////////
//Transactions
//////////////////////////////////////////////////////////////////////////////

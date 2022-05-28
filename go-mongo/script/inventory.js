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
db.inventory.insertMany([
    {item: "canvas", qty: 100, size: {h: 28, w: 35.5, uom: "cm"}, status: "A"},
    {item: "journal", qty: 25, size: {h: 14, w: 21, uom: "cm"}, status: "A"},
    {item: "mat", qty: 85, size: {h: 27.9, w: 35.5, uom: "cm"}, status: "A"},
    {item: "mousepad", qty: 25, size: {h: 19, w: 22.85, uom: "cm"}, status: "P"},
    {item: "notebook", qty: 50, size: {h: 8.5, w: 11, uom: "in"}, status: "P"},
    {item: "paper", qty: 100, size: {h: 8.5, w: 11, uom: "in"}, status: "D"},
    {item: "planner", qty: 75, size: {h: 22.85, w: 30, uom: "cm"}, status: "D"},
    {item: "postcard", qty: 45, size: {h: 10, w: 15.25, uom: "cm"}, status: "A"},
    {item: "sketchbook", qty: 80, size: {h: 14, w: 21, uom: "cm"}, status: "A"},
    {item: "sketch pad", qty: 95, size: {h: 22.85, w: 30.5, uom: "cm"}, status: "A"}
]);


/*
Update Documents in a Collection
- Some update operators, such as $set, will create the field if the field does not exist.
{
  <update operator>: { <field1>: <value1>, ... },
  <update operator>: { <field2>: <value2>, ... },
  ...
}
 */

//1.Update a Single Document - update the first document where item equals to "paper"
db.inventory.updateOne(
    {
        item: "paper"
    },
    {
        $set: {"size.uom": "cm", status: "P"},
        $currentDate: {lastModified: true}
    }
)


//2.Update a multiple documents - update all documents where qty is less than 50
db.inventory.updateMany(
    {"qty": {$lt: 50}},
    {
        $set: {"size.uom": "in", status: "P"},
        $currentDate: {lastModified: true}
    }
)

//replace One - matching 되는 데이터를 교체함
db.inventory.replaceOne(
    {item: "paper"},
    {item: "paper", instock: [{warehouse: "A", qty: 60}, {warehouse: "B", qty: 40}]}
)


/*
Updates with Aggregation Pipeline:
- Using the aggregation pipeline allows for a more expressive update statement,
such as expressing conditional updates based on current field values
or updating one field using the value of another field(s).

$addFields
$set
$project
$unset
$replaceRoot
$replaceWith
 */

db.students.insertMany([
    {_id: 1, test1: 95, test2: 92, test3: 90, modified: new Date("2020-01-05")},
    {_id: 2, test1: 98, test2: 100, test3: 102, modified: new Date("2020-01-05")},
    {_id: 3, test1: 95, test2: 110, modified: new Date("2020-01-04")}
])

db.students.find()

/*
Example 1 - $$NOW 값은 current datetime 값을 얻어 올 수 있다
1.updateOne() operation uses an aggregation pipeline to update the document with _id: 3:
Specifically, the pipeline consists of a $set stage which adds the test3 field
(and sets its value to 98) to the document and sets the modified field to the current datetime.
The operation uses the aggregation variable NOW for the current datetime.
To access the variable, prefix with $$ and enclose in quotes.
 */
db.students.updateOne({_id: 3}, [{$set: {"test3": 98, modified: "$$NOW"}}])
db.students.find().pretty()


/*
Example 2 - unset 되어 있는 field에 zero로 세팅해줌
 */
db.students2.insertMany([
    {"_id": 1, quiz1: 8, test2: 100, quiz2: 9, modified: new Date("2020-01-05")},
    {"_id": 2, quiz2: 5, test1: 80, test2: 89, modified: new Date("2020-01-05")},
])

/*
- $replaceRoot stage with a $mergeObjects expression to set default values for the quiz1, quiz2, test1 and test2 fields.
The aggregation variable ROOT refers to the current document being modified.
To access the variable, prefix with $$ and enclose in quotes.
The current document fields will override the default values.
- $set stage to update the modified field to the current datetime.
The operation uses the aggregation variable NOW for the current datetime.
To access the variable, prefix with $$ and enclose in quotes.
 */
db.students2.updateMany({},
    [
        {
            $replaceRoot: {
                newRoot:
                    {$mergeObjects: [{quiz1: 0, quiz2: 0, test1: 0, test2: 0}, "$$ROOT"]}
            }
        },
        {$set: {modified: "$$NOW"}}
    ]
)

/*
Example 3 - average, grade를 추가로 넣음
 */
db.students3.insertMany([
    {"_id": 1, "tests": [95, 92, 90], "modified": ISODate("2019-01-01T00:00:00Z")},
    {"_id": 2, "tests": [94, 88, 90], "modified": ISODate("2019-01-01T00:00:00Z")},
    {"_id": 3, "tests": [70, 75, 82], "modified": ISODate("2019-01-01T00:00:00Z")}
]);

/*
- $set stage to calculate the truncated average value of the tests array elements and to update the modified field to the current datetime.
To calculate the truncated average, the stage uses the $avg and $trunc expressions.
The operation uses the aggregation variable NOW for the current datetime.
To access the variable, prefix with $$ and enclose in quotes.
- $set stage to add the grade field based on the average using the $switch expression.
 */
db.students3.updateMany(
    {},
    [
        {$set: {average: {$trunc: [{$avg: "$tests"}, 0]}, modified: "$$NOW"}},
        {
            $set: {
                grade: {
                    $switch: {
                        branches: [
                            {case: {$gte: ["$average", 90]}, then: "A"},
                            {case: {$gte: ["$average", 80]}, then: "B"},
                            {case: {$gte: ["$average", 70]}, then: "C"},
                            {case: {$gte: ["$average", 60]}, then: "D"}
                        ],
                        default: "F"
                    }
                }
            }
        }
    ]
)

/*
Example 4 - quizzes list에 값을 추가함
 */
db.students4.insertMany([
    {"_id": 1, "quizzes": [4, 6, 7]},
    {"_id": 2, "quizzes": [5]},
    {"_id": 3, "quizzes": [10, 10, 10]}
])

db.students4.updateOne({_id: 2},
    [{$set: {quizzes: {$concatArrays: ["$quizzes", [8, 6]]}}}]
)

/*x
Example 5 -
 */
db.temperatures.insertMany([
    {"_id": 1, "date": ISODate("2019-06-23"), "tempsC": [4, 12, 17]},
    {"_id": 2, "date": ISODate("2019-07-07"), "tempsC": [14, 24, 11]},
    {"_id": 3, "date": ISODate("2019-10-30"), "tempsC": [18, 6, 8]}
])

/*
Specifically, the pipeline consists of an $addFields stage to add a new array field tempsF that contains the temperatures in Fahrenheit.
To convert each celsius temperature in the tempsC array to Fahrenheit,
the stage uses the $map expression with $add and $multiply expressions.
 */
db.temperatures.updateMany({},
    [
        {
            $addFields: {
                "tempsF": {
                    $map: {
                        input: "$tempsC",
                        as: "celsius",
                        in: {$add: [{$multiply: ["$$celsius", 9 / 5]}, 32]}
                    }
                }
            }
        }
    ]
)

//////////////////////////////////////////////////////////////////////////////
//Delete Documents
//////////////////////////////////////////////////////////////////////////////


//////////////////////////////////////////////////////////////////////////////
//Aggregation Operations
//////////////////////////////////////////////////////////////////////////////


//////////////////////////////////////////////////////////////////////////////
//Transactions
//////////////////////////////////////////////////////////////////////////////


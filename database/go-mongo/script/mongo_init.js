//https://www.tutorialspoint.com/mongodb/mongodb_create_database.htm
//login으로 들어가려
use admin
db.auth('mongoadmin', 'mongopassword')

//show database
show dbs


//create database
use mydb

db.movie.insert({ "name": "tutorials point" })

//drop database
//db.dropDtabase()

//create collection
db.createCollection("mycollection")
show collections

//capped: fixed size collection; automatically overwrites its oldest entries when it reaches its maximum size
db.createCollection("mycol", { capped: true, autoIndexID: true, size: 6142800, max: 10000 })
db.tutorialspoint.insert({ "name": "tutorialspoint" })

//drop collection
db.mycollection.drop()
show collections

//insert document
//insert()
db.users.insert({
    _id: ObjectId("507f191e810c19729de860ea"),
    title: "MongoDB Overview",
    description: "MongoDB is no sql database",
    by: "tutorials point",
    url: "http://www.tutorialspoint.com",
    tags: ['mongodb', 'database', 'NoSQL'],
    likes: 100
})

//post
db.createCollection("post")
db.post.insert([
    {
        title: "MongoDB Overview",
        description: "MongoDB is no SQL database",
        by: "tutorials point",
        url: "http://www.tutorialspoint.com",
        tags: ["mongodb", "database", "NoSQL"],
        likes: 100
    },
    {
        title: "NoSQL Database",
        description: "NoSQL database doesn't have tables",
        by: "tutorials point",
        url: "http://www.tutorialspoint.com",
        tags: ["mongodb", "database", "NoSQL"],
        likes: 20,
        comments: [
            {
                user: "user1",
                message: "My first comment",
                dateCreated: new Date(2013, 11, 10, 2, 35),
                like: 0
            }
        ]
    }
])

//insertOne()
db.createCollection("empDetails")
db.empDetails.insertOne(
    {
        First_Name: "Radhika",
        Last_Name: "Sharma",
        Date_Of_Birth: "1995-09-26",
        e_mail: "radhika_sharma.123@gmail.com",
        phone: "9848022338"
    })

//insertMany()
db.empDetails.insertMany(
    [
        {
            First_Name: "Radhika",
            Last_Name: "Sharma",
            Date_Of_Birth: "1995-09-26",
            e_mail: "radhika_sharma.123@gmail.com",
            phone: "9000012345"
        },
        {
            First_Name: "Rachel",
            Last_Name: "Christopher",
            Date_Of_Birth: "1990-02-16",
            e_mail: "Rachel_Christopher.123@gmail.com",
            phone: "9000054321"
        },
        {
            First_Name: "Fathima",
            Last_Name: "Sheik",
            Date_Of_Birth: "1990-02-16",
            e_mail: "Fathima_Sheik.123@gmail.com",
            phone: "9000054321"
        }
    ]
)


//query document
//find()
use sampleDB
db.createCollection("mycol")

db.mycol.insert([
    {
        title: "MongoDB Overview",
        description: "MongoDB is no SQL database",
        by: "tutorials point",
        url: "http://www.tutorialspoint.com",
        tags: ["mongodb", "database", "NoSQL"],
        likes: 100
    },
    {
        title: "NoSQL Database",
        description: "NoSQL database doesn't have tables",
        by: "tutorials point",
        url: "http://www.tutorialspoint.com",
        tags: ["mongodb", "database", "NoSQL"],
        likes: 20,
        comments: [
            {
                user: "user1",
                message: "My first comment",
                dateCreated: new Date(2013, 11, 10, 2, 35),
                like: 0
            }
        ]
    }
])

//pretty()
db.mycol.find()
db.mycol.find().pretty()

//findOne()
db.mycol.findOne({ title: "MongoDB Overview" })

//where by = 'tutorials point'
db.mycol.find({ "by": "tutorials point" })

//where likes < 50
db.mycol.find({ "likes": { $lt: 50 } })

//where likes <= 50
db.mycol.find({ "likes": { $lte: 50 } })

//where likes > 50
db.mycol.find({ "likes": { $gt: 50 } })

//where likes >= 50
db.mycol.find({ "likes": { $gte: 50 } })

//where likes != 50
db.mycol.find({ "likes": { $ne: 50 } })

//in :["Raj", "Ram", "Raghu"]
db.mycol.find({ "name": { $in: ["Raj", "Ram", "Raghu"] } })

//not in :["Ramu", "Raghav"] 
db.mycol.find({ "name": { $nin: ["Ramu", "Raghav"] } })

//AND
//where by = 'tutorials point' AND title = 'MongoDB Overview'
db.mycol.find({ $and: [{ "by": "tutorials point" }, { "title": "MongoDB Overview" }] })

//OR
db.mycol.find({ $or: [{ "by": "tutorials point" }, { "title": "MongoDB Overview" }] })

//Using AND and OR Together
//where likes>10 AND (by = 'tutorials point' OR title = 'MongoDB Overview')
db.mycol.find({
    "likes": { $gt: 10 },
    $or: [{ "by": "tutorials point" }, { "title": "MongoDB Overview" }]
})

//NOR in MongoDB
//$nor performs a logical NOR operation on an array of one or more query expression and selects the documents 
//that fail all the query expressions in the array
db.empDetails.insertMany(
    [
        {
            First_Name: "Radhika",
            Last_Name: "Sharma",
            Age: "26",
            e_mail: "radhika_sharma.123@gmail.com",
            phone: "9000012345"
        },
        {
            First_Name: "Rachel",
            Last_Name: "Christopher",
            Age: "27",
            e_mail: "Rachel_Christopher.123@gmail.com",
            phone: "9000054321"
        },
        {
            First_Name: "Fathima",
            Last_Name: "Sheik",
            Age: "24",
            e_mail: "Fathima_Sheik.123@gmail.com",
            phone: "9000054321"
        }
    ]
)

db.empDetails.find(
    {
        $nor: [
            { "First_Name": "Radhika" },
            { "Last_Name": "Christopher" }
        ]
    }
)

//NOT in MongoDB
db.empDetails.find({ "Age": { $not: { $gt: "25" } } })


//update document
//update()
//db.COLLECTION_NAME.update(SELECTION_CRITERIA, UPDATED_DATA)
db.mycol.update({ 'title': 'MongoDB Overview' }, { $set: { 'title': 'New MongoDB Tutorial' } })
db.mycol.find()


//By default, MongoDB will update only a single document. 
//To update multiple documents, you need to set a parameter 'multi' to true.
db.mycol.update({ 'title': 'MongoDB Overview' },
    { $set: { 'title': 'New MongoDB Tutorial' } }, { multi: true })

//save()
//db.COLLECTION_NAME.save({_id:ObjectId(),NEW_DATA})
//save() method replaces the existing document with the new document passed in the save() method.
db.mycol.save(
    {
        "_id": ObjectId("507f191e810c19729de860e3"),
        "title": "Tutorials Point New Topic",
        "by": "Tutorials Point"
    }
)

db.mycol.find()

//findOneAndUpdate()
//db.COLLECTION_NAME.findOneAndUpdate(SELECTIOIN_CRITERIA, UPDATED_DATA)
//updates the values in the existing document.
db.empDetails.insertMany(
    [
        {
            First_Name: "Radhika",
            Last_Name: "Sharma",
            Age: "26",
            e_mail: "radhika_sharma.123@gmail.com",
            phone: "9000012345"
        },
        {
            First_Name: "Rachel",
            Last_Name: "Christopher",
            Age: "27",
            e_mail: "Rachel_Christopher.123@gmail.com",
            phone: "9000054321"
        },
        {
            First_Name: "Fathima",
            Last_Name: "Sheik",
            Age: "24",
            e_mail: "Fathima_Sheik.123@gmail.com",
            phone: "9000054321"
        }
    ]
)

db.empDetails.findOneAndUpdate(
    { First_Name: 'Radhika' },
    { $set: { Age: '30', e_mail: 'radhika_newemail@gmail.com' } }
)

//updateOne()
//updates a single document which matches the given filter
//db.COLLECTION_NAME.updateOne(<filter>, <update>)
db.empDetails.updateOne(
    { First_Name: 'Radhika' },
    { $set: { Age: '30', e_mail: 'radhika_newemail@gmail.com' } }
)

//updateMany()
//db.COLLECTION_NAME.update(<filter>, <update>)
//updates all the documents that matches the given filter.
db.empDetails.updateMany(
    { Age: { $gt: "25" } },
    { $set: { Age: '00' } }
)


//delete document
//remove()
//method is used to remove a document from the collection.
//db.COLLECTION_NAME.remove(DELLETION_CRITTERIA)
//remove all the documents whose title is 'MongoDB Overview'.
db.mycol.remove({ 'title': 'MongoDB Overview' })


//remove justOne
//db.COLLECTION_NAME.remove(DELETION_CRITERIA,1)

//rmeove all documents
db.mycol.remove({})

//project document
//find()
//1: show the field 
//0: hide the fields.
//db.COLLECTION_NAME.find({},{KEY:1})
db.mycol.find()
db.mycol.find({}, { "title": 1, "description": 1, _id: 0 })


//limiting records
//limit()
//db.COLLECTION_NAME.find().limit(NUMBER)
db.mycol.find({}, { "title": 1, _id: 0 }).limit(2)

//skip()
//which also accepts number type argument and is used to skip the number of documents.
//default value in skip() method is 0.
//db.COLLECTION_NAME.find().limit(NUMBER).skip(NUMBER)
db.mycol.find({}, { "title": 1, _id: 0 }).limit(1).skip(1)


//sorting records
//sort()
//1:ascending order
//-1:descending order
//sorted by title in the descending order.
db.mycol.find({}, { "title": 1, _id: 0 }).sort({ "title": -1 })


//indexing
//createIndex()
//1:ascending order
//-1:descending order
//db.COLLECTION_NAME.createIndex({KEY:1})
db.mycol.createIndex({ "title": 1 })

//pass multiple fields, to create index on multiple fields
db.mycol.createIndex({ "title": 1, "description": -1 })
db.mycol.find({}, { "title": 1, _id: 0 })

//dropIndex()
//db.COLLECTION_NAME.dropIndex({KEY:1})
db.mycol.dropIndex({ "title": 1 })

//getIndexes() 
//returns the description of all the indexes int the collection.
//db.COLLECTION_NAME.getIndexes()
db.mycol.getIndexes()

//aggregation
//Aggregation operations group values from multiple documents together, and can perform 
//a variety of operations on the grouped data to return a single result
//db.COLLECTION_NAME.aggregate(AGGREGATE_OPERATION)
db.mycol.insertMany([
    {
        title: 'MongoDB Overview',
        description: 'MongoDB is no sql database',
        by_user: 'tutorials point',
        url: 'http://www.tutorialspoint.com',
        tags: ['mongodb', 'database', 'NoSQL'],
        likes: 100
    },
    {
        title: 'NoSQL Overview',
        description: 'No sql database is very fast',
        by_user: 'tutorials point',
        url: 'http://www.tutorialspoint.com',
        tags: ['mongodb', 'database', 'NoSQL'],
        likes: 10
    },
    {
        title: 'Neo4j Overview',
        description: 'Neo4j is no sql database',
        by_user: 'Neo4j',
        url: 'http://www.neo4j.com',
        tags: ['neo4j', 'database', 'NoSQL'],
        likes: 750
    }
])

//select by_user, count(*) from mycol group by by_user
db.mycol.aggregate([{ $group: { _id: "$by_user", num_tutorial: { $sum: 1 } } }])

//$sum
//Sums up the defined value from all documents in the collection.
db.mycol.aggregate([{ $group: { _id: "$by_user", num_tutorial: { $sum: "$likes" } } }])

//$avg
//Calculates the average of all given values from all documents in the collection.
db.mycol.aggregate([{ $group: { _id: "$by_user", num_tutorial: { $avg: "$likes" } } }])

//$min
//Gets the minimum of the corresponding values from all documents in the collection.
db.mycol.aggregate([{ $group: { _id: "$by_user", num_tutorial: { $min: "$likes" } } }])

//$max
//Gets the maximum of the corresponding values from all documents in the collection.
db.mycol.aggregate([{ $group: { _id: "$by_user", num_tutorial: { $max: "$likes" } } }])

//$push ??
//Inserts the value to an array in the resulting document.
db.mycol.aggregate([{ $group: { _id: "$by_user", url: { $push: "$url" } } }])


//$addToSet ??
//Inserts the value to an array in the resulting document but does not create duplicates.
db.mycol.aggregate([{ $group: { _id: "$by_user", url: { $addToSet: "$url" } } }])

//$first
//Gets the first document from the source documents according to the grouping. 
//Typically this makes only sense together with some previously applied “$sort”-stage.
db.mycol.aggregate([{ $group: { _id: "$by_user", first_url: { $first: "$url" } } }])

//$last
//Gets the last document from the source documents according to the grouping. 
//Typically this makes only sense together with some previously applied “$sort”-stage.
db.mycol.aggregate([{ $group: { _id: "$by_user", last_url: { $last: "$url" } } }])




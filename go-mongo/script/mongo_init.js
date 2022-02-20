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
db.mycol.update({ 'title': 'MongoDB Overview' }, { $set: { 'title': 'New MongoDB Tutorial' } })
db.mycol.find()


/
//


//delete document

//project document

//limiting records

//sorting records

//indexing

//aggregation



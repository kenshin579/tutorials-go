use clone
db.inventory.insertMany([
    {item: "journal", qty: 25, tags: ["blank", "red"], size: {h: 14, w: 21, uom: "cm"}},
    {item: "mat", qty: 85, tags: ["gray"], size: {h: 27.9, w: 35.5, uom: "cm"}},
    {item: "mousepad", qty: 25, tags: ["gel", "blue"], size: {h: 19, w: 22.85, uom: "cm"}}
])

//1.db.collection.find().forEach() command
db.inventory.find().forEach(
    function(docs){
        db.inventory2.insert(docs);
    })

//2.db.collecion.aggregate() command
db.inventory.aggregate([{ $match: {} }, { $out: "inventory3" }])



db.terrier.remove({})
for (let i = 0; i < 10000; i++) {
    db.terrier.insert({
        "_id": "5d8a342ba9e109cb56" + i,
        "DogName": "June",
        "lastupdated": ISODate("2020-09-23T12:55:23.902Z"),
        "Activity": "Sleeping",
        "Publish_status": false
    });
}
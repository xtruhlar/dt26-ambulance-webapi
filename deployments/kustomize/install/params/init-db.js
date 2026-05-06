const mongoHost = process.env.AMBULANCE_API_MONGODB_HOST
const mongoPort = process.env.AMBULANCE_API_MONGODB_PORT

const mongoUser = process.env.AMBULANCE_API_MONGODB_USERNAME
const mongoPassword = process.env.AMBULANCE_API_MONGODB_PASSWORD

const database = process.env.AMBULANCE_API_MONGODB_DATABASE
const collection = process.env.AMBULANCE_API_MONGODB_COLLECTION

const retrySeconds = parseInt(process.env.RETRY_CONNECTION_SECONDS || "5") || 5;

// try to connect to mongoDB until it is available
let connection;
while(true) {
    try {
        connection = Mongo(`mongodb://${mongoUser}:${mongoPassword}@${mongoHost}:${mongoPort}`);
        break;
    } catch (exception) {
        print(`Cannot connect to mongoDB: ${exception}`);
        print(`Will retry after ${retrySeconds} seconds`)
        sleep(retrySeconds * 1000);
    }
}

// if database and collection already exist, exit with success
const databases = connection.getDBNames()
if (databases.includes(database)) {
    const dbInstance = connection.getDB(database)
    const collections = dbInstance.getCollectionNames()
    if (collections.includes(collection)) {
        print(`Collection '${collection}' already exists in database '${database}'`)
        process.exit(0);
    }
}

// initialize database and collection
const db = connection.getDB(database)
db.createCollection(collection)

// create index on id field
db[collection].createIndex({ "id": 1 })

// insert sample ambulance with consultation entries
let result = db[collection].insertMany([
    {
        "id": "dt26-ambulance",
        "consultationEntries": [
            {
                "id": "c001",
                "patientId": "460527-jan-novak",
                "patientName": "Ján Novák",
                "condition": "Hypertenzia",
                "status": "active",
                "createdAt": new Date("2038-12-24T10:05:00Z")
            },
            {
                "id": "c002",
                "patientId": "780907-maria-kovacova",
                "patientName": "Mária Kováčová",
                "condition": "Diabetes typu 2",
                "status": "pending",
                "createdAt": new Date("2038-12-24T10:25:00Z")
            }
        ]
    }
]);

if (result.writeError) {
    console.error(result)
    print(`Error when writing the data: ${result.errmsg}`)
}

process.exit(0);

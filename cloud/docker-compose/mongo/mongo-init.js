db.createUser(
    {
        user: "muser",
        pwd: "password",
        roles: [
            {
                role: "readWrite",
                db: "go-mongo"
            }
        ]
    }
);

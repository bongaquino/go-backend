db = db.getSiblingDB('admin');

// Ensure root user exists
if (db.system.users.find({ user: 'root' }).count() === 0) {
  db.Create({
    user: 'root',
    pwd: 'password',
    roles: [{ role: 'root', db: 'admin' }]
  });
}

// Create 'bongaquino' database and user
db = db.getSiblingDB('bongaquino');

if (db.system.users.find({ user: 'bongaquino_user' }).count() === 0) {
  db.createUser({
    user: 'bongaquino_user',
    pwd: 'bongaquino_password',
    roles: [{ role: 'readWrite', db: 'bongaquino' }]
  });
}

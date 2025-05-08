db = db.getSiblingDB('admin');

// Ensure root user exists
if (db.system.users.find({ user: 'root' }).count() === 0) {
  db.Create({
    user: 'root',
    pwd: 'password',
    roles: [{ role: 'root', db: 'admin' }]
  });
}

// Create 'koneksi' database and user
db = db.getSiblingDB('koneksi');

if (db.system.users.find({ user: 'koneksi_user' }).count() === 0) {
  db.createUser({
    user: 'koneksi_user',
    pwd: 'koneksi_password',
    roles: [{ role: 'readWrite', db: 'koneksi' }]
  });
}

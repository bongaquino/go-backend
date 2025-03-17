db = db.getSiblingDB('admin');

if (db.system.users.find({ user: 'root' }).count() === 0) {
  db.createUser({
    user: 'root',
    pwd: 'password',
    roles: [{ role: 'root', db: 'admin' }]
  });
}

db = db.getSiblingDB('koneksi');

db.createUser({
  user: 'koneksi_user',
  pwd: 'koneksi_password',
  roles: [{ role: 'readWrite', db: 'koneksi' }]
});

db.sampleCollection.insertMany([
  { name: 'Sample Data 1', value: 100 },
  { name: 'Sample Data 2', value: 200 },
  { name: 'Sample Data 3', value: 300 }
]);
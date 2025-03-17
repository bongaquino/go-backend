db = db.getSiblingDB('admin');

// Ensure root user exists
if (db.system.users.find({ user: 'root' }).count() === 0) {
  db.createUser({
    user: 'root',
    pwd: 'password',
    roles: [{ role: 'root', db: 'admin' }]
  });
}

// Define sample data
const sampleData = [
  { name: 'Sample Data 1', value: 100 },
  { name: 'Sample Data 2', value: 200 },
  { name: 'Sample Data 3', value: 300 }
];

// Create 'koneksi' database and user
db = db.getSiblingDB('koneksi');
db.createUser({
  user: 'koneksi_user',
  pwd: 'koneksi_password',
  roles: [{ role: 'readWrite', db: 'koneksi' }]
});

// Insert data into structured collections
db.account_service.insertMany(sampleData);
db.dashboard_service.insertMany(sampleData);
db.backup_service.insertMany(sampleData);

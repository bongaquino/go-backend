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

// Create 'account_service' database and user
db = db.getSiblingDB('account_service');
db.createUser({
  user: 'account_service_user',
  pwd: 'account_service_password',
  roles: [{ role: 'readWrite', db: 'account_service' }]
});
db.data.insertMany(sampleData);

// Create 'dashboard_service' database and user
db = db.getSiblingDB('dashboard_service');
db.createUser({
  user: 'dashboard_service_user',
  pwd: 'dashboard_service_password',
  roles: [{ role: 'readWrite', db: 'dashboard_service' }]
});
db.data.insertMany(sampleData);

// Create 'backup_service' database and user
db = db.getSiblingDB('backup_service');
db.createUser({
  user: 'backup_service_user',
  pwd: 'backup_service_password',
  roles: [{ role: 'readWrite', db: 'backup_service' }]
});
db.data.insertMany(sampleData);

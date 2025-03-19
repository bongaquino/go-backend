## IAM Service

### **Overview**

We will implement authentication and authorization for two types of accounts:

1. **User Accounts** – Follow **Role-Based Access Control (RBAC)**.
2. **Service Accounts** – Follow **Policy-Based Access Control (PBAC)**.

The application requires fine-grained control over both user and service account permissions to securely handle data within the system.

### **Authentication Flow**

- **Users** authenticate via **email/password**.
- **Service accounts** authenticate using **JSON key files**.
- Both use **JWT** for session management.

### **Authorization Flow**

- **User accounts** use **roles** and **permissions**.
- **Service accounts** use **JSON-based policies** with conditions.

---

## **High-Level Access Control Design**

### **Role-Based Access Control (RBAC) for Users**

**Roles** are predefined sets of **permissions**, and users are assigned roles to determine what actions they can perform.

#### **Relationships**
- **User** → **has one or more** → **Roles**
- **Role** → **grants one or more** → **Permissions**
- **Permissions** → **define allowed actions**

#### **Example**
- **Admin Role**: Can **upload, download, list**
- **User Role**: Can **upload, download, list**

```
[User] ---> [User_Role] ---> [Role] ---> [Role_Permission] ---> [Permission]
```

### **Policy-Based Access Control (PBAC) for Service Accounts**

**Policies** define access rules dynamically and are linked to **service accounts** via a policy identifier.

#### **Relationships**
- **Service Account** → **assigned one** → **Policy**
- **Policy** → **defines** → **Resource access rules**
- **Rules** → **specify allowed actions and conditions**

#### **Example**
- **Backup Agent Policy**: Allows **upload, download, list**
- **Analytics Engine Policy**: Allows **download, list**

```
[Service Account] ---> [Policy] ---> [Policy_Permission] ---> [Permission]
```

### **Differences Between Policies and Roles**

| Aspect       | RBAC (Roles & Permissions) | PBAC (Policies) |
|--------------|----------------------------|-----------------|
| **Scope**    | User access control        | Service account control |
| **Structure** | Predefined roles with fixed permissions | Flexible policies with dynamic rules |
| **Flexibility** | Less flexible, requires role updates for new permissions | Highly flexible, policies can be customized per service |
| **Use Case** | Users accessing UI/API endpoints | Service accounts interacting with APIs |

---

## **Implementation Plan**

1. **User Authentication**
   - Register/login users with password hashing and JWT issuance.
   - Assign roles to users.
   - Validate user permissions during API requests.

2. **Service Authentication**
   - Generate and store **JSON key files** containing credentials.
   - The JSON file includes:
     ```json
     {
       "client_id": "abcdef123456",
       "client_secret": "hashed_secret",
       "private_key": "-----BEGIN PRIVATE KEY-----...-----END PRIVATE KEY-----",
       "policy_id": "65fcd89a89c9a8f123456789"
     }
     ```
   - Services must use this key file to authenticate API requests.
   - Use JWT with **policy_id** for authorization.

3. **Access Enforcement**
   - **Users:** Match roles to permissions.
   - **Services:** Validate requests against policies before granting access.

---

## **MongoDB Schema Design**

### **Users Collection**

| Field           | Type     | Description                         |
| --------------- | -------- | ----------------------------------- |
| `_id`           | ObjectId | Unique user identifier              |
| `email`         | String   | User's email address                |
| `password_hash` | String   | Hashed password for authentication  |
| `created_at`    | Date     | Timestamp when the user was created |
| `updated_at`    | Date     | Timestamp when the user was updated |

### **User Roles Collection**

| Field    | Type     | Description                           |
| -------- | -------- | ------------------------------------  |
| `_id`    | ObjectId | Unique user-role identifier           |
| `user_id`| ObjectId | Reference to the user                 |
| `role_id`| ObjectId | Reference to the assigned role        |

### **Roles Collection**

| Field        | Type     | Description                         |
|------------- | -------- | ----------------------------------- |
| `_id`        | ObjectId | Unique role identifier              |
| `name`       | String   | Name of the role                    |
| `created_at` | Date     | Timestamp when the role was created |

### **Role Permissions Collection**

| Field          | Type     | Description                            |
|--------------- | -------- | -------------------------------------- |
| `_id`          | ObjectId | Unique role-permission identifier      |
| `role_id`      | ObjectId | Reference to a role                    |
| `permission_id`| ObjectId | Reference to a permission              |

### **Permissions Collection**

| Field        | Type     | Description                               |
|------------- | -------- | ----------------------------------------- |
| `_id`        | ObjectId | Unique permission identifier              |
| `name`       | String   | Name of the permission                    |
| `created_at` | Date     | Timestamp when the permission was created |

### **Service Accounts Collection**

| Field                | Type     | Description                       |
| -------------------- | -------- | --------------------------------- |
| `_id`                | ObjectId | Unique service account identifier |
| `name`               | String   | Name of the service account       |
| `description`        | String   | Description of the service account|
| `client_key`         | String   | Unique client key                 |
| `client_secret_hash` | String   | Hashed secret for authentication  |
| `policy_id`          | ObjectId | Reference to a policy document    |
| `created_at`         | Date     | Timestamp when created            |
| `updated_at`         | Date     | Timestamp when updated            |

### **Policies Collection**

| Field        | Type     | Description                           |
| ------------ | -------- | ------------------------------------- |
| `_id`        | ObjectId | Unique policy identifier              |
| `name`       | String   | Name of the policy                    |
| `created_at` | Date     | Timestamp when the policy was created |

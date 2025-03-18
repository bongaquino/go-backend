## Account Service

### **Overview**

We will implement authentication and authorization for two types of accounts:

1. **User Accounts** – Follow **Role-Based Access Control (RBAC)**.
2. **Service Accounts** – Follow **Policy-Based Access Control (PBAC)**.

### **Authentication Flow**

- **Users** authenticate via **email/password**.
- **Service accounts** authenticate using **JSON key files**.
- Both use **JWT** for session management.

### **Authorization Flow**

- **User accounts** use **roles** and **permissions**.
- **Service accounts** use **JSON-based policies** with conditions.

---

## **High-Level Access Control Design**

### **Role-Based Access Control (RBAC)**

**Roles** are predefined sets of **permissions**, and users are assigned roles to determine what actions they can perform.

#### **Relationships**
- **User** → **has one or more** → **Roles**
- **Role** → **grants one or more** → **Permissions**

#### **Example**
- **Admin Role**: Can **upload, download, list**
- **User Role**: Can **upload, download, list**

```
[User] ---> [Role] ---> [Permissions]
       |          |         |
       |          |         --> Upload
       |          |         --> Download
       |          |         --> List
       |          --> Admin
       |          --> User
```

### **Policy-Based Access Control (PBAC)**

**Policies** define access rules dynamically and are linked to **service accounts** via a policy identifier.

#### **Relationships**
- **Service Account** → **assigned one** → **Policy**
- **Policy** → **defines** → **Resource access rules**

#### **Example**
- **Backup Agent Policy**: Allows **upload, download, list**
- **Analytics Engine Policy**: Allows **download, list**

```
[Service Account] ---> [Policy] ---> [Access Rules]
         |                 |             |
         |                 |             --> Upload ✅
         |                 |             --> Download ✅
         |                 |             --> List ✅
         |                 --> Backup Agent Policy
         |                 --> Analytics Engine Policy
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

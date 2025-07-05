# Emergency Response Application: Full Technical Process Flow (Enhanced with H3 & Updated DB)

## 1. User Registration and Profile Creation

- **User Action:**
  - Opens the mobile app and registers an account.
  - Inputs personal and medical information.

- **Data Captured:**
  - Full name, phone number, email.
  - Medical details: allergies, blood type, chronic conditions, medications.
  - Emergency contacts (optional).

- **Database Tables:**
  - `users` (id, name, phone, email, hashed_password, created_at).
  - `medical_profiles` (id, user_id, blood_type, allergies, conditions, medications, emergency_contacts).
  - `clients` (id, name, organization_details, created_at).
  - `zones` (id, client_id, name, geojson_boundary).

- **Backend Systems Used:**
  - **AWS Cognito** for secure user authentication.
  - **Aurora PostgreSQL** with **PostGIS** for storing user profiles and spatial data.
  - Medical details encrypted at the database level using AWS KMS.
  - Medical profiles accessible via API for paramedic apps (authorized by user_id).

---

## 2. Emergency Creation Workflow

- **Trigger:**
  - User presses "Request Emergency Assistance" in the app.
  - Alternatively, an external system can trigger an emergency via API.

- **Data Flow:**
  - API call to Fargate-hosted container service.
  - Emergency details include:
    - User ID reference.
    - GPS coordinates.
    - H3 Index (calculated server-side).
    - Type of emergency.
    - Timestamp.

- **Database Tables:**
  - `incidents` (id, user_id, client_id, zone_id, type, status, h3_index, created_at, location_geometry).

- **Event System:**
  - Backend publishes an event to **AWS Kinesis Event Bus**, classified as `incident.created`.

---

## 3. Event Bus, Dispatch System, and Real-time Updates

- **AWS Kinesis Event Bus:**
  - Routes events to:
    - **Dispatch Dashboard** (web-based, uses OpenStreetMap tiles directly).
    - Notification Service.
    - Analytics or audit systems.

- **Dispatch Dashboard:**
  - Built with React and WebSockets.
  - Displays incidents on a map using **OpenStreetMap** tiles.
  - Secure API calls fetch:
    - Patient data (`medical_profiles`).
    - Available paramedics from `paramedics` table.
    - Zone boundaries from `zones` table.
  - Uses PostGIS and `h3_index` for fast spatial queries.

- **Paramedic Availability:**
  - Live locations stored in **Redis**, updated every 15 seconds.
  - Paramedic records include latest `h3_index`.
  - Dashboard displays real-time responder positions.

---

## 4. Paramedic Assignment and Notification Flow

- **Dispatcher Action:**
  - Selects closest available paramedic using H3 proximity and PostGIS spatial filters.

- **Backend Interaction:**
  - Updates `incidents` table with assigned responder.
  - Sends push notification to paramedic via AWS SNS.

- **Paramedic Mobile App:**
  - Flutter or React Native.
  - Securely fetches incident and patient medical details.
  - Embeds **Google Maps SDK** for maps display.

---

## 5. Navigation, Directions, and ETA Calculation

- **Maps & Navigation:**
  - Google Maps SDK provides map tiles.
  - App uses **Google Directions API** for routes.

- **ETA Flow:**
  - App calculates ETA.
  - Sends ETA to backend.
  - Backend updates `incidents` table with ETA.
  - Dashboard reflects ETA in real-time.

- **Live Tracking:**
  - Paramedic app sends live location to Redis.
  - Backend updates dashboard via WebSockets.

---

## 6. Chat Functionality (Dispatcher ↔ Paramedic, Dispatcher ↔ Client)

- **Chat System:**
  - WebSockets with Redis Pub/Sub for real-time delivery.
  - Messages stored in `chat_messages` table (id, incident_id, sender_id, recipient_id, message, created_at).

- **Client ↔ Dispatcher Chat:**
  - Chat initiated upon emergency request.
  - Optional AI Chatbot (AWS Lex) integrated into conversation.

---

## 7. Incident Completion and Data Capture

- **Arrival:**
  - Paramedic marks incident as `on scene`.
  - Location and timestamp updated.

- **Completion:**
  - Incident status `completed`.
  - Generates report with:
    - Timeline, Route details, ETA vs. actual.
    - Chat transcript.
    - Medical data accessed.

---

## 8. System Architecture Overview

| Component             | Technology/Service          |
| --------------------- | --------------------------- |
| User Authentication   | AWS Cognito                 |
| Database              | Aurora PostgreSQL + PostGIS |
| Realtime Cache        | AWS ElastiCache Redis       |
| Event Bus             | AWS Kinesis                 |
| Backend Compute       | AWS Fargate containers      |
| Paramedic Maps        | Google Maps SDK             |
| Directions & ETA      | Google Directions API       |
| Dispatcher Maps       | OpenStreetMap tiles         |
| AI Chatbot (optional) | AWS Lex                     |
| Notifications         | AWS SNS                     |

---

## 9. Notifications and Error Handling

- **Push Notifications:**
  - AWS SNS for mobile push alerts.

- **Error Logging:**
  - Centralized logs via **CloudWatch**.
  - API errors, formatting issues, and system faults captured.

- **Development & Testing:**
  - Separate dev environment.
  - API responses logged.
  - Event Bus messages validated.
  - Integration tests simulate emergency flows.

---

## 10. Cost Comparison (1000 Emergencies per Day)

### Google Costs
- **Maps SDK Tiles:** Free
- **Directions API:**
  - $5 per 1000 requests → $5/day → $150/month.

### AWS Costs
- **Aurora PostgreSQL:**
  - ~db.t3.medium: $120/month.
  - Storage: $0.10/GB.
- **ElastiCache Redis:** ~ $15/month.
- **Kinesis Event Bus:** Minimal ($0.02/month).
- **Fargate Backend:** ~ $100/month.
- **SNS Notifications:** ~ $5/month.

### Additional Costs
- CloudWatch: $10-50/month.
- Security:
  - AWS Shield Basic: Free.
  - AWS WAF (optional): $20-50/month.

**Total Estimate:**
- Best Case: ~$500/month.
- Worst Case: $1000-2000/month (growth, advanced routing, increased emergencies).

---

## 11. Security and Compliance
- Data encrypted at rest (KMS).
- Role-based API permissions.
- Access auditing.
- POPIA-compliant (South Africa).
- Optional:
  - AWS GuardDuty.
  - PrivateLink/VPN.

---

## 12. Scalability & Optimizations

- H3 index stored in incidents & paramedics tables for fast proximity searches.
- Regional data sharding.
- Redis clustering.
- Aurora read replicas.
- CDN for static content.
- Tile caching for paramedic zones.
- Option to separate spatial-heavy PostGIS workloads to distinct DB instance.

---

**End of Updated Technical Document with H3 & Database Enhancements**


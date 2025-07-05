# Cursor AI Backend Development Context for Emergency Response System

## Project Scope: API & Database Layer Only

This document provides targeted technical context to guide AI tools (such as Cursor.ai) in generating reliable, production-ready API logic, database schemas, and backend services for the Emergency Response Application.

**Strict Boundaries:**
- This scope is **backend only**.
- AI should focus on:
  - Database schemas.
  - API endpoint structures.
  - Core service logic for incidents, profiles, chat.
  - Event-based workflows.
  - Real-time Redis usage.
- AI should **not** generate:
  - Frontend code (web or mobile apps).
  - Google Maps SDK or mobile integrations.
  - Complex AI chatbot dialogue.
  - DevOps infrastructure beyond service references.

---

## System Architecture Relevant to Backend

| Component           | Tech/Service                          |
|--------------------|---------------------------------------|
| Database           | AWS Aurora PostgreSQL + PostGIS       |
| Real-time Cache    | AWS ElastiCache Redis (managed Redis) |
| Event Bus          | AWS Kinesis                           |
| Backend Compute    | AWS Fargate (containers)              |
| Spatial Indexing   | H3 (hexagonal grid system)            |

---

## Required Database Tables (Initial)

- `users` (id, name, phone, email, hashed_password, created_at)
- `medical_profiles` (id, user_id, blood_type, allergies, conditions, medications, emergency_contacts)
- `clients` (id, org_name, contact_person, contact_info)
- `zones` (id, client_id, name, geojson_boundary)
- `paramedics` (id, user_id, zone_id, active_status)
- `incidents` (id, user_id, zone_id, type, status, h3_index, location_geometry, created_at, assigned_paramedic_id, eta)
- `chat_messages` (id, incident_id, sender_id, receiver_type, message, created_at)

---

## Core API Endpoints (Example Specification)

- `POST /api/register` → User signup.
- `POST /api/medical` → Create/update medical profile.
- `POST /api/emergency` → Trigger new incident.
- `GET /api/incidents` → Dispatcher fetches incidents.
- `POST /api/assign` → Assign paramedic to incident.
- `POST /api/chat` → Send chat message.
- `GET /api/chat/:incident_id` → Fetch chat history.

Additional endpoints can include auth, zone uploads, client onboarding.

---

## Event-Driven Workflow Expectations

- **Incident Created:**
  - API call inserts `incidents` record.
  - Event emitted to AWS Kinesis (`incident.created`).
- **Matching Service Logic:**
  - Calculates nearest paramedic using:
    - H3 index filtering.
    - PostGIS `ST_Distance` if refined search needed.
  - Updates incident with assigned paramedic.
- **Real-Time Data:**
  - Redis Pub/Sub broadcasts updates.
  - Live locations stored in Redis.

---

## Redis Responsibilities

- Real-time paramedic locations (lat/lon by responder id).
- Chat message delivery channels.
- Presence tracking for online status.
- Optional caching of zone boundaries for fast lookup.

---

## Research Guidance for AI Assistance

Cursor or similar tools should:
- Review H3 usage patterns for hex-based proximity search.
- Understand PostGIS spatial functions (e.g., `ST_Contains`, `ST_Distance`).
- Familiarize with AWS Kinesis event publishing models.
- Apply best practices for:
  - Secure API design.
  - Role-based access for paramedics/dispatchers.
  - Data encryption at rest for sensitive fields.
- Consider PostgreSQL indexing strategies, especially for location queries.

---

## Special Considerations

- Medical profiles must be securely fetched by paramedic apps.
- Incident zone assignment enforces responder visibility boundaries.
- All chat, event, and real-time features rely on unified Redis layer.
- H3 index stored **within `incidents` table** for efficient matching.
- Errors, API failures, or data issues should log to centralized system (e.g., CloudWatch).

---

## Out of Scope (Do Not Generate)

- Frontend application code (React, Flutter, etc.).
- Google Maps SDK integration specifics.
- Full AI chatbot conversations.
- Infrastructure-as-Code (Terraform, CloudFormation).
- Third-party login providers unless specifically instructed.

---

## Next Steps

Cursor AI can proceed with:
- Generating SQL schema migrations.
- Building API handler logic.
- Structuring service layer for incidents, chat, profiles.
- Integrating Redis patterns.
- Preparing for event emissions (Kinesis stubs).

**Confirm before attempting production-ready mobile or frontend integrations.**

---

*End of Cursor Backend Context Document*


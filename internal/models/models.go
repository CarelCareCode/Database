package models

import (
	"time"
)

type User struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Phone     string    `json:"phone" db:"phone"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"hashed_password"`
	UserType  string    `json:"user_type" db:"user_type"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type MedicalProfile struct {
	ID                string    `json:"id" db:"id"`
	UserID            string    `json:"user_id" db:"user_id"`
	BloodType         string    `json:"blood_type" db:"blood_type"`
	Allergies         string    `json:"allergies" db:"allergies"`
	Conditions        string    `json:"conditions" db:"conditions"`
	Medications       string    `json:"medications" db:"medications"`
	EmergencyContacts string    `json:"emergency_contacts" db:"emergency_contacts"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

type Client struct {
	ID               string    `json:"id" db:"id"`
	Name             string    `json:"name" db:"name"`
	OrganizationName string    `json:"organization_name" db:"organization_name"`
	ContactPerson    string    `json:"contact_person" db:"contact_person"`
	ContactInfo      string    `json:"contact_info" db:"contact_info"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

type Zone struct {
	ID              string    `json:"id" db:"id"`
	ClientID        string    `json:"client_id" db:"client_id"`
	Name            string    `json:"name" db:"name"`
	GeojsonBoundary string    `json:"geojson_boundary" db:"geojson_boundary"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

type Paramedic struct {
	ID           string    `json:"id" db:"id"`
	UserID       string    `json:"user_id" db:"user_id"`
	ZoneID       string    `json:"zone_id" db:"zone_id"`
	ActiveStatus bool      `json:"active_status" db:"active_status"`
	CurrentLat   float64   `json:"current_lat" db:"current_lat"`
	CurrentLng   float64   `json:"current_lng" db:"current_lng"`
	H3Index      string    `json:"h3_index" db:"h3_index"`
	LastSeen     time.Time `json:"last_seen" db:"last_seen"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type Incident struct {
	ID                  string    `json:"id" db:"id"`
	UserID              string    `json:"user_id" db:"user_id"`
	ClientID            string    `json:"client_id" db:"client_id"`
	ZoneID              string    `json:"zone_id" db:"zone_id"`
	Type                string    `json:"type" db:"type"`
	Status              string    `json:"status" db:"status"`
	H3Index             string    `json:"h3_index" db:"h3_index"`
	LocationGeometry    string    `json:"location_geometry" db:"location_geometry"`
	LocationLat         float64   `json:"location_lat" db:"location_lat"`
	LocationLng         float64   `json:"location_lng" db:"location_lng"`
	AssignedParamedicID *string   `json:"assigned_paramedic_id" db:"assigned_paramedic_id"`
	ETA                 *int      `json:"eta" db:"eta"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

type ChatMessage struct {
	ID           string    `json:"id" db:"id"`
	IncidentID   string    `json:"incident_id" db:"incident_id"`
	SenderID     string    `json:"sender_id" db:"sender_id"`
	ReceiverID   string    `json:"receiver_id" db:"receiver_id"`
	ReceiverType string    `json:"receiver_type" db:"receiver_type"`
	Message      string    `json:"message" db:"message"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// Request/Response structs
type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	UserType string `json:"user_type" validate:"required,oneof=client paramedic dispatcher"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	UserID   string `json:"user_id"`
	UserType string `json:"user_type"`
}

type EmergencyRequest struct {
	Type        string  `json:"type" validate:"required"`
	Latitude    float64 `json:"latitude" validate:"required"`
	Longitude   float64 `json:"longitude" validate:"required"`
	Description string  `json:"description"`
}

type LocationUpdate struct {
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}

type ChatMessageRequest struct {
	IncidentID   string `json:"incident_id" validate:"required"`
	ReceiverID   string `json:"receiver_id" validate:"required"`
	ReceiverType string `json:"receiver_type" validate:"required,oneof=dispatcher paramedic client"`
	Message      string `json:"message" validate:"required"`
}

type AssignParamedicRequest struct {
	ParamedicID string `json:"paramedic_id" validate:"required"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=pending assigned en_route on_scene completed cancelled"`
}

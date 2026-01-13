package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents system users with role-based access
type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email        string         `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"not null" json:"-"`
	Role         string         `gorm:"not null;index" json:"role"` // patient, clinic_doctor, clinic_manager, admin
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Patient *Patient `gorm:"foreignKey:UserID" json:"patient,omitempty"`
	Clinic  *Clinic  `gorm:"foreignKey:UserID" json:"clinic,omitempty"`
}

// Patient represents a patient in the system
type Patient struct {
	ID                    uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID                uuid.UUID      `gorm:"type:uuid;index" json:"user_id"`
	FirstName             string         `gorm:"not null" json:"first_name"`
	LastName              string         `gorm:"not null" json:"last_name"`
	DateOfBirth           *time.Time     `json:"date_of_birth"`
	Phone                 string         `json:"phone"`
	DiagnocatPatientID    string         `gorm:"uniqueIndex" json:"diagnocat_patient_id,omitempty"`
	PreferredCity         string         `json:"preferred_city"`
	PreferredDistrict     string         `json:"preferred_district"`
	PreferredPriceSegment string         `json:"preferred_price_segment"` // economy, business, premium
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User    User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Studies []Study `gorm:"foreignKey:PatientID" json:"studies,omitempty"`
}

// Study represents a DICOM study (one imaging session)
type Study struct {
	ID                  uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PatientID           uuid.UUID      `gorm:"type:uuid;not null;index" json:"patient_id"`
	DiagnocatStudyUID   string         `gorm:"uniqueIndex" json:"diagnocat_study_uid,omitempty"`
	Status              string         `gorm:"not null;index;default:'created'" json:"status"` // created, uploading, uploaded, processing, completed, failed
	Modality            string         `json:"modality"`                                       // CBCT, CT, etc.
	StudyDate           *time.Time     `json:"study_date"`
	UploadedAt          *time.Time     `json:"uploaded_at"`
	CompletedAt         *time.Time     `json:"completed_at"`
	ErrorMessage        string         `json:"error_message,omitempty"`
	DiagnocatResultJSON map[string]any `gorm:"type:jsonb" json:"diagnocat_result,omitempty"`
	DiagnocatReportPDF  []byte         `json:"-"` // Store PDF binary
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`

	// Relationships
	Patient      Patient       `gorm:"foreignKey:PatientID" json:"patient,omitempty"`
	PlanVersions []PlanVersion `gorm:"foreignKey:StudyID" json:"plan_versions,omitempty"`
}

// PlanVersion represents an immutable snapshot of a treatment plan
type PlanVersion struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	StudyID   uuid.UUID `gorm:"type:uuid;not null;index" json:"study_id"`
	Version   int       `gorm:"not null" json:"version"`           // Auto-increment per study
	Source    string    `gorm:"default:'diagnocat'" json:"source"` // diagnocat, manual, modified
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	Study     Study      `gorm:"foreignKey:StudyID" json:"study,omitempty"`
	PlanItems []PlanItem `gorm:"foreignKey:PlanVersionID" json:"plan_items,omitempty"`
}

// PlanItem represents a single procedure in the treatment plan
type PlanItem struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PlanVersionID uuid.UUID      `gorm:"type:uuid;not null;index" json:"plan_version_id"`
	ToothNumber   *int           `gorm:"check:tooth_number >= 11 AND tooth_number <= 48" json:"tooth_number,omitempty"` // FDI notation
	Specialty     string         `gorm:"not null;index" json:"specialty"`                                               // therapy, orthopedics, surgery, hygiene, periodontics
	ProcedureCode string         `json:"procedure_code"`
	ProcedureName string         `gorm:"not null" json:"procedure_name"`
	Diagnosis     string         `json:"diagnosis"`
	Quantity      int            `gorm:"default:1" json:"quantity"`
	Notes         string         `json:"notes"`
	Metadata      map[string]any `gorm:"type:jsonb" json:"metadata,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`

	// Relationships
	PlanVersion PlanVersion `gorm:"foreignKey:PlanVersionID" json:"-"`
}

// Clinic represents a dental clinic
type Clinic struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID          uuid.UUID      `gorm:"type:uuid;index" json:"user_id"` // Manager account
	Name            string         `gorm:"not null" json:"name"`
	LegalName       string         `json:"legal_name"`
	LicenseNumber   string         `gorm:"uniqueIndex" json:"license_number"`
	YearEstablished int            `json:"year_established"`
	City            string         `gorm:"not null;index" json:"city"`
	District        string         `gorm:"index" json:"district"`
	Address         string         `json:"address"`
	Phone           string         `json:"phone"`
	Email           string         `json:"email"`
	Website         string         `json:"website"`
	PriceSegment    string         `json:"price_segment"`                  // economy, business, premium
	IsActive        bool           `gorm:"default:false" json:"is_active"` // Requires admin approval
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	PriceListItems []PriceListItem `gorm:"foreignKey:ClinicID" json:"price_list_items,omitempty"`
	Offers         []Offer         `gorm:"foreignKey:ClinicID" json:"offers,omitempty"`
}

// PriceListItem represents a single service price in clinic's pricelist
type PriceListItem struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ClinicID      uuid.UUID `gorm:"type:uuid;not null;index" json:"clinic_id"`
	Specialty     string    `gorm:"not null;index" json:"specialty"`
	ProcedureCode string    `gorm:"index" json:"procedure_code"`
	ProcedureName string    `gorm:"not null" json:"procedure_name"`
	PriceFrom     float64   `json:"price_from"` // Minimum price
	PriceTo       float64   `json:"price_to"`   // Maximum price (optional range)
	IsActive      bool      `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Relationships
	Clinic Clinic `gorm:"foreignKey:ClinicID" json:"-"`
}

// OfferRequest represents a patient's request for offers from clinics
type OfferRequest struct {
	ID                uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PatientID         uuid.UUID   `gorm:"type:uuid;not null;index" json:"patient_id"`
	PlanVersionID     uuid.UUID   `gorm:"type:uuid;not null;index" json:"plan_version_id"`
	SelectedItemIDs   []uuid.UUID `gorm:"type:uuid[];not null" json:"selected_item_ids"` // Which plan items to include
	PreferredCity     string      `json:"preferred_city"`
	PreferredDistrict string      `json:"preferred_district"`
	PriceSegment      string      `json:"price_segment"`
	Status            string      `gorm:"default:'open'" json:"status"` // open, closed
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`

	// Relationships
	Patient     Patient     `gorm:"foreignKey:PatientID" json:"patient,omitempty"`
	PlanVersion PlanVersion `gorm:"foreignKey:PlanVersionID" json:"plan_version,omitempty"`
	Offers      []Offer     `gorm:"foreignKey:OfferRequestID" json:"offers,omitempty"`
}

// Offer represents a clinic's proposal for a treatment plan
type Offer struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OfferRequestID   uuid.UUID `gorm:"type:uuid;not null;index" json:"offer_request_id"`
	ClinicID         uuid.UUID `gorm:"type:uuid;not null;index" json:"clinic_id"`
	TotalPrice       float64   `gorm:"not null" json:"total_price"`
	DiscountPercent  float64   `json:"discount_percent"`
	HasInstallment   bool      `gorm:"default:false" json:"has_installment"`
	InstallmentTerms string    `json:"installment_terms"`
	SpecialOffer     string    `json:"special_offer"`
	EstimatedDays    int       `json:"estimated_days"`
	Status           string    `gorm:"default:'pending'" json:"status"` // pending, accepted, rejected
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	// Relationships
	OfferRequest OfferRequest `gorm:"foreignKey:OfferRequestID" json:"offer_request,omitempty"`
	Clinic       Clinic       `gorm:"foreignKey:ClinicID" json:"clinic,omitempty"`
	Orders       []Order      `gorm:"foreignKey:OfferID" json:"orders,omitempty"`
}

// Order represents an accepted offer and tracks treatment progress
type Order struct {
	ID                 uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OfferID            uuid.UUID  `gorm:"type:uuid;not null;index" json:"offer_id"`
	PatientID          uuid.UUID  `gorm:"type:uuid;not null;index" json:"patient_id"`
	ClinicID           uuid.UUID  `gorm:"type:uuid;not null;index" json:"clinic_id"`
	Status             string     `gorm:"not null;default:'new'" json:"status"` // new, consultation_scheduled, in_progress, completed, cancelled
	ConsultationDate   *time.Time `json:"consultation_date"`
	TreatmentStarted   *time.Time `json:"treatment_started"`
	TreatmentCompleted *time.Time `json:"treatment_completed"`
	CancellationReason string     `json:"cancellation_reason"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`

	// Relationships
	Offer   Offer   `gorm:"foreignKey:OfferID" json:"offer,omitempty"`
	Patient Patient `gorm:"foreignKey:PatientID" json:"patient,omitempty"`
	Clinic  Clinic  `gorm:"foreignKey:ClinicID" json:"clinic,omitempty"`
	Slots   []Slot  `gorm:"foreignKey:OrderID" json:"slots,omitempty"`
}

// Slot represents available appointment times
type Slot struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ClinicID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"clinic_id"`
	DoctorName  string     `json:"doctor_name"`
	Specialty   string     `json:"specialty"`
	StartTime   time.Time  `gorm:"not null;index" json:"start_time"`
	EndTime     time.Time  `gorm:"not null" json:"end_time"`
	IsAvailable bool       `gorm:"default:true" json:"is_available"`
	OrderID     *uuid.UUID `gorm:"type:uuid;index" json:"order_id,omitempty"` // Null if available
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Relationships
	Clinic Clinic `gorm:"foreignKey:ClinicID" json:"-"`
	Order  *Order `gorm:"foreignKey:OrderID" json:"order,omitempty"`
}

// JobQueue for simple DB-based async processing
type JobQueue struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	JobType     string         `gorm:"not null;index"` // upload_to_diagnocat, poll_results, etc.
	Payload     map[string]any `gorm:"type:jsonb;not null"`
	Status      string         `gorm:"not null;default:'pending';index"` // pending, processing, completed, failed
	Attempts    int            `gorm:"default:0"`
	MaxAttempts int            `gorm:"default:3"`
	Error       string         `gorm:"type:text"`
	ScheduledAt time.Time      `gorm:"not null;index"`
	StartedAt   *time.Time
	CompletedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// AuditLog for tracking all important actions
type AuditLog struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID     *uuid.UUID     `gorm:"type:uuid;index"`
	Action     string         `gorm:"not null;index"` // create_study, upload_file, create_offer, etc.
	EntityType string         `gorm:"not null"`       // study, offer, order, etc.
	EntityID   *uuid.UUID     `gorm:"type:uuid"`
	Details    map[string]any `gorm:"type:jsonb"`
	IPAddress  string
	CreatedAt  time.Time `gorm:"index"`
}

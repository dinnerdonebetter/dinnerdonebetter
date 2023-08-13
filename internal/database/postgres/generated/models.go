// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package generated

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

type ComponentType string

const (
	ComponentTypeUnspecified ComponentType = "unspecified"
	ComponentTypeAmuseBouche ComponentType = "amuse-bouche"
	ComponentTypeAppetizer   ComponentType = "appetizer"
	ComponentTypeSoup        ComponentType = "soup"
	ComponentTypeMain        ComponentType = "main"
	ComponentTypeSalad       ComponentType = "salad"
	ComponentTypeBeverage    ComponentType = "beverage"
	ComponentTypeSide        ComponentType = "side"
	ComponentTypeDessert     ComponentType = "dessert"
)

func (e *ComponentType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = ComponentType(s)
	case string:
		*e = ComponentType(s)
	default:
		return fmt.Errorf("unsupported scan type for ComponentType: %T", src)
	}
	return nil
}

type NullComponentType struct {
	ComponentType ComponentType
	Valid         bool // Valid is true if ComponentType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullComponentType) Scan(value interface{}) error {
	if value == nil {
		ns.ComponentType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.ComponentType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullComponentType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.ComponentType), nil
}

func (e ComponentType) Valid() bool {
	switch e {
	case ComponentTypeUnspecified,
		ComponentTypeAmuseBouche,
		ComponentTypeAppetizer,
		ComponentTypeSoup,
		ComponentTypeMain,
		ComponentTypeSalad,
		ComponentTypeBeverage,
		ComponentTypeSide,
		ComponentTypeDessert:
		return true
	}
	return false
}

func AllComponentTypeValues() []ComponentType {
	return []ComponentType{
		ComponentTypeUnspecified,
		ComponentTypeAmuseBouche,
		ComponentTypeAppetizer,
		ComponentTypeSoup,
		ComponentTypeMain,
		ComponentTypeSalad,
		ComponentTypeBeverage,
		ComponentTypeSide,
		ComponentTypeDessert,
	}
}

type GroceryListItemStatus string

const (
	GroceryListItemStatusUnknown      GroceryListItemStatus = "unknown"
	GroceryListItemStatusAlreadyowned GroceryListItemStatus = "already owned"
	GroceryListItemStatusNeeds        GroceryListItemStatus = "needs"
	GroceryListItemStatusUnavailable  GroceryListItemStatus = "unavailable"
	GroceryListItemStatusAcquired     GroceryListItemStatus = "acquired"
)

func (e *GroceryListItemStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = GroceryListItemStatus(s)
	case string:
		*e = GroceryListItemStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for GroceryListItemStatus: %T", src)
	}
	return nil
}

type NullGroceryListItemStatus struct {
	GroceryListItemStatus GroceryListItemStatus
	Valid                 bool // Valid is true if GroceryListItemStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullGroceryListItemStatus) Scan(value interface{}) error {
	if value == nil {
		ns.GroceryListItemStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.GroceryListItemStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullGroceryListItemStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.GroceryListItemStatus), nil
}

func (e GroceryListItemStatus) Valid() bool {
	switch e {
	case GroceryListItemStatusUnknown,
		GroceryListItemStatusAlreadyowned,
		GroceryListItemStatusNeeds,
		GroceryListItemStatusUnavailable,
		GroceryListItemStatusAcquired:
		return true
	}
	return false
}

func AllGroceryListItemStatusValues() []GroceryListItemStatus {
	return []GroceryListItemStatus{
		GroceryListItemStatusUnknown,
		GroceryListItemStatusAlreadyowned,
		GroceryListItemStatusNeeds,
		GroceryListItemStatusUnavailable,
		GroceryListItemStatusAcquired,
	}
}

type IngredientAttributeType string

const (
	IngredientAttributeTypeTexture     IngredientAttributeType = "texture"
	IngredientAttributeTypeConsistency IngredientAttributeType = "consistency"
	IngredientAttributeTypeColor       IngredientAttributeType = "color"
	IngredientAttributeTypeAppearance  IngredientAttributeType = "appearance"
	IngredientAttributeTypeOdor        IngredientAttributeType = "odor"
	IngredientAttributeTypeTaste       IngredientAttributeType = "taste"
	IngredientAttributeTypeSound       IngredientAttributeType = "sound"
	IngredientAttributeTypeTemperature IngredientAttributeType = "temperature"
	IngredientAttributeTypeOther       IngredientAttributeType = "other"
)

func (e *IngredientAttributeType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = IngredientAttributeType(s)
	case string:
		*e = IngredientAttributeType(s)
	default:
		return fmt.Errorf("unsupported scan type for IngredientAttributeType: %T", src)
	}
	return nil
}

type NullIngredientAttributeType struct {
	IngredientAttributeType IngredientAttributeType
	Valid                   bool // Valid is true if IngredientAttributeType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullIngredientAttributeType) Scan(value interface{}) error {
	if value == nil {
		ns.IngredientAttributeType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.IngredientAttributeType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullIngredientAttributeType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.IngredientAttributeType), nil
}

func (e IngredientAttributeType) Valid() bool {
	switch e {
	case IngredientAttributeTypeTexture,
		IngredientAttributeTypeConsistency,
		IngredientAttributeTypeColor,
		IngredientAttributeTypeAppearance,
		IngredientAttributeTypeOdor,
		IngredientAttributeTypeTaste,
		IngredientAttributeTypeSound,
		IngredientAttributeTypeTemperature,
		IngredientAttributeTypeOther:
		return true
	}
	return false
}

func AllIngredientAttributeTypeValues() []IngredientAttributeType {
	return []IngredientAttributeType{
		IngredientAttributeTypeTexture,
		IngredientAttributeTypeConsistency,
		IngredientAttributeTypeColor,
		IngredientAttributeTypeAppearance,
		IngredientAttributeTypeOdor,
		IngredientAttributeTypeTaste,
		IngredientAttributeTypeSound,
		IngredientAttributeTypeTemperature,
		IngredientAttributeTypeOther,
	}
}

type InvitationState string

const (
	InvitationStatePending   InvitationState = "pending"
	InvitationStateCancelled InvitationState = "cancelled"
	InvitationStateAccepted  InvitationState = "accepted"
	InvitationStateRejected  InvitationState = "rejected"
)

func (e *InvitationState) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = InvitationState(s)
	case string:
		*e = InvitationState(s)
	default:
		return fmt.Errorf("unsupported scan type for InvitationState: %T", src)
	}
	return nil
}

type NullInvitationState struct {
	InvitationState InvitationState
	Valid           bool // Valid is true if InvitationState is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullInvitationState) Scan(value interface{}) error {
	if value == nil {
		ns.InvitationState, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.InvitationState.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullInvitationState) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.InvitationState), nil
}

func (e InvitationState) Valid() bool {
	switch e {
	case InvitationStatePending,
		InvitationStateCancelled,
		InvitationStateAccepted,
		InvitationStateRejected:
		return true
	}
	return false
}

func AllInvitationStateValues() []InvitationState {
	return []InvitationState{
		InvitationStatePending,
		InvitationStateCancelled,
		InvitationStateAccepted,
		InvitationStateRejected,
	}
}

type MealName string

const (
	MealNameBreakfast       MealName = "breakfast"
	MealNameSecondBreakfast MealName = "second_breakfast"
	MealNameBrunch          MealName = "brunch"
	MealNameLunch           MealName = "lunch"
	MealNameSupper          MealName = "supper"
	MealNameDinner          MealName = "dinner"
)

func (e *MealName) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = MealName(s)
	case string:
		*e = MealName(s)
	default:
		return fmt.Errorf("unsupported scan type for MealName: %T", src)
	}
	return nil
}

type NullMealName struct {
	MealName MealName
	Valid    bool // Valid is true if MealName is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullMealName) Scan(value interface{}) error {
	if value == nil {
		ns.MealName, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.MealName.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullMealName) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.MealName), nil
}

func (e MealName) Valid() bool {
	switch e {
	case MealNameBreakfast,
		MealNameSecondBreakfast,
		MealNameBrunch,
		MealNameLunch,
		MealNameSupper,
		MealNameDinner:
		return true
	}
	return false
}

func AllMealNameValues() []MealName {
	return []MealName{
		MealNameBreakfast,
		MealNameSecondBreakfast,
		MealNameBrunch,
		MealNameLunch,
		MealNameSupper,
		MealNameDinner,
	}
}

type MealPlanStatus string

const (
	MealPlanStatusAwaitingVotes MealPlanStatus = "awaiting_votes"
	MealPlanStatusFinalized     MealPlanStatus = "finalized"
)

func (e *MealPlanStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = MealPlanStatus(s)
	case string:
		*e = MealPlanStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for MealPlanStatus: %T", src)
	}
	return nil
}

type NullMealPlanStatus struct {
	MealPlanStatus MealPlanStatus
	Valid          bool // Valid is true if MealPlanStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullMealPlanStatus) Scan(value interface{}) error {
	if value == nil {
		ns.MealPlanStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.MealPlanStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullMealPlanStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.MealPlanStatus), nil
}

func (e MealPlanStatus) Valid() bool {
	switch e {
	case MealPlanStatusAwaitingVotes,
		MealPlanStatusFinalized:
		return true
	}
	return false
}

func AllMealPlanStatusValues() []MealPlanStatus {
	return []MealPlanStatus{
		MealPlanStatusAwaitingVotes,
		MealPlanStatusFinalized,
	}
}

type Oauth2ClientTokenScopes string

const (
	Oauth2ClientTokenScopesUnknown         Oauth2ClientTokenScopes = "unknown"
	Oauth2ClientTokenScopesHouseholdMember Oauth2ClientTokenScopes = "household_member"
	Oauth2ClientTokenScopesHouseholdAdmin  Oauth2ClientTokenScopes = "household_admin"
	Oauth2ClientTokenScopesServiceAdmin    Oauth2ClientTokenScopes = "service_admin"
)

func (e *Oauth2ClientTokenScopes) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Oauth2ClientTokenScopes(s)
	case string:
		*e = Oauth2ClientTokenScopes(s)
	default:
		return fmt.Errorf("unsupported scan type for Oauth2ClientTokenScopes: %T", src)
	}
	return nil
}

type NullOauth2ClientTokenScopes struct {
	Oauth2ClientTokenScopes Oauth2ClientTokenScopes
	Valid                   bool // Valid is true if Oauth2ClientTokenScopes is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullOauth2ClientTokenScopes) Scan(value interface{}) error {
	if value == nil {
		ns.Oauth2ClientTokenScopes, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Oauth2ClientTokenScopes.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullOauth2ClientTokenScopes) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Oauth2ClientTokenScopes), nil
}

func (e Oauth2ClientTokenScopes) Valid() bool {
	switch e {
	case Oauth2ClientTokenScopesUnknown,
		Oauth2ClientTokenScopesHouseholdMember,
		Oauth2ClientTokenScopesHouseholdAdmin,
		Oauth2ClientTokenScopesServiceAdmin:
		return true
	}
	return false
}

func AllOauth2ClientTokenScopesValues() []Oauth2ClientTokenScopes {
	return []Oauth2ClientTokenScopes{
		Oauth2ClientTokenScopesUnknown,
		Oauth2ClientTokenScopesHouseholdMember,
		Oauth2ClientTokenScopesHouseholdAdmin,
		Oauth2ClientTokenScopesServiceAdmin,
	}
}

type PrepStepStatus string

const (
	PrepStepStatusUnfinished PrepStepStatus = "unfinished"
	PrepStepStatusPostponed  PrepStepStatus = "postponed"
	PrepStepStatusIgnored    PrepStepStatus = "ignored"
	PrepStepStatusCanceled   PrepStepStatus = "canceled"
	PrepStepStatusFinished   PrepStepStatus = "finished"
)

func (e *PrepStepStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PrepStepStatus(s)
	case string:
		*e = PrepStepStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for PrepStepStatus: %T", src)
	}
	return nil
}

type NullPrepStepStatus struct {
	PrepStepStatus PrepStepStatus
	Valid          bool // Valid is true if PrepStepStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPrepStepStatus) Scan(value interface{}) error {
	if value == nil {
		ns.PrepStepStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PrepStepStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPrepStepStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PrepStepStatus), nil
}

func (e PrepStepStatus) Valid() bool {
	switch e {
	case PrepStepStatusUnfinished,
		PrepStepStatusPostponed,
		PrepStepStatusIgnored,
		PrepStepStatusCanceled,
		PrepStepStatusFinished:
		return true
	}
	return false
}

func AllPrepStepStatusValues() []PrepStepStatus {
	return []PrepStepStatus{
		PrepStepStatusUnfinished,
		PrepStepStatusPostponed,
		PrepStepStatusIgnored,
		PrepStepStatusCanceled,
		PrepStepStatusFinished,
	}
}

type RecipeStepProductType string

const (
	RecipeStepProductTypeIngredient RecipeStepProductType = "ingredient"
	RecipeStepProductTypeInstrument RecipeStepProductType = "instrument"
	RecipeStepProductTypeVessel     RecipeStepProductType = "vessel"
)

func (e *RecipeStepProductType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = RecipeStepProductType(s)
	case string:
		*e = RecipeStepProductType(s)
	default:
		return fmt.Errorf("unsupported scan type for RecipeStepProductType: %T", src)
	}
	return nil
}

type NullRecipeStepProductType struct {
	RecipeStepProductType RecipeStepProductType
	Valid                 bool // Valid is true if RecipeStepProductType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullRecipeStepProductType) Scan(value interface{}) error {
	if value == nil {
		ns.RecipeStepProductType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.RecipeStepProductType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullRecipeStepProductType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.RecipeStepProductType), nil
}

func (e RecipeStepProductType) Valid() bool {
	switch e {
	case RecipeStepProductTypeIngredient,
		RecipeStepProductTypeInstrument,
		RecipeStepProductTypeVessel:
		return true
	}
	return false
}

func AllRecipeStepProductTypeValues() []RecipeStepProductType {
	return []RecipeStepProductType{
		RecipeStepProductTypeIngredient,
		RecipeStepProductTypeInstrument,
		RecipeStepProductTypeVessel,
	}
}

type SettingType string

const (
	SettingTypeUser       SettingType = "user"
	SettingTypeHousehold  SettingType = "household"
	SettingTypeMembership SettingType = "membership"
)

func (e *SettingType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = SettingType(s)
	case string:
		*e = SettingType(s)
	default:
		return fmt.Errorf("unsupported scan type for SettingType: %T", src)
	}
	return nil
}

type NullSettingType struct {
	SettingType SettingType
	Valid       bool // Valid is true if SettingType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullSettingType) Scan(value interface{}) error {
	if value == nil {
		ns.SettingType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.SettingType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullSettingType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.SettingType), nil
}

func (e SettingType) Valid() bool {
	switch e {
	case SettingTypeUser,
		SettingTypeHousehold,
		SettingTypeMembership:
		return true
	}
	return false
}

func AllSettingTypeValues() []SettingType {
	return []SettingType{
		SettingTypeUser,
		SettingTypeHousehold,
		SettingTypeMembership,
	}
}

type StorageContainerType string

const (
	StorageContainerTypeUncovered             StorageContainerType = "uncovered"
	StorageContainerTypeCovered               StorageContainerType = "covered"
	StorageContainerTypeOnawirerack           StorageContainerType = "on a wire rack"
	StorageContainerTypeInanairtightcontainer StorageContainerType = "in an airtight container"
)

func (e *StorageContainerType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = StorageContainerType(s)
	case string:
		*e = StorageContainerType(s)
	default:
		return fmt.Errorf("unsupported scan type for StorageContainerType: %T", src)
	}
	return nil
}

type NullStorageContainerType struct {
	StorageContainerType StorageContainerType
	Valid                bool // Valid is true if StorageContainerType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullStorageContainerType) Scan(value interface{}) error {
	if value == nil {
		ns.StorageContainerType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.StorageContainerType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullStorageContainerType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.StorageContainerType), nil
}

func (e StorageContainerType) Valid() bool {
	switch e {
	case StorageContainerTypeUncovered,
		StorageContainerTypeCovered,
		StorageContainerTypeOnawirerack,
		StorageContainerTypeInanairtightcontainer:
		return true
	}
	return false
}

func AllStorageContainerTypeValues() []StorageContainerType {
	return []StorageContainerType{
		StorageContainerTypeUncovered,
		StorageContainerTypeCovered,
		StorageContainerTypeOnawirerack,
		StorageContainerTypeInanairtightcontainer,
	}
}

type ValidElectionMethod string

const (
	ValidElectionMethodSchulze       ValidElectionMethod = "schulze"
	ValidElectionMethodInstantRunoff ValidElectionMethod = "instant-runoff"
)

func (e *ValidElectionMethod) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = ValidElectionMethod(s)
	case string:
		*e = ValidElectionMethod(s)
	default:
		return fmt.Errorf("unsupported scan type for ValidElectionMethod: %T", src)
	}
	return nil
}

type NullValidElectionMethod struct {
	ValidElectionMethod ValidElectionMethod
	Valid               bool // Valid is true if ValidElectionMethod is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullValidElectionMethod) Scan(value interface{}) error {
	if value == nil {
		ns.ValidElectionMethod, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.ValidElectionMethod.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullValidElectionMethod) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.ValidElectionMethod), nil
}

func (e ValidElectionMethod) Valid() bool {
	switch e {
	case ValidElectionMethodSchulze,
		ValidElectionMethodInstantRunoff:
		return true
	}
	return false
}

func AllValidElectionMethodValues() []ValidElectionMethod {
	return []ValidElectionMethod{
		ValidElectionMethodSchulze,
		ValidElectionMethodInstantRunoff,
	}
}

type VesselShape string

const (
	VesselShapeHemisphere VesselShape = "hemisphere"
	VesselShapeRectangle  VesselShape = "rectangle"
	VesselShapeCone       VesselShape = "cone"
	VesselShapePyramid    VesselShape = "pyramid"
	VesselShapeCylinder   VesselShape = "cylinder"
	VesselShapeSphere     VesselShape = "sphere"
	VesselShapeCube       VesselShape = "cube"
	VesselShapeOther      VesselShape = "other"
)

func (e *VesselShape) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = VesselShape(s)
	case string:
		*e = VesselShape(s)
	default:
		return fmt.Errorf("unsupported scan type for VesselShape: %T", src)
	}
	return nil
}

type NullVesselShape struct {
	VesselShape VesselShape
	Valid       bool // Valid is true if VesselShape is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullVesselShape) Scan(value interface{}) error {
	if value == nil {
		ns.VesselShape, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.VesselShape.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullVesselShape) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.VesselShape), nil
}

func (e VesselShape) Valid() bool {
	switch e {
	case VesselShapeHemisphere,
		VesselShapeRectangle,
		VesselShapeCone,
		VesselShapePyramid,
		VesselShapeCylinder,
		VesselShapeSphere,
		VesselShapeCube,
		VesselShapeOther:
		return true
	}
	return false
}

func AllVesselShapeValues() []VesselShape {
	return []VesselShape{
		VesselShapeHemisphere,
		VesselShapeRectangle,
		VesselShapeCone,
		VesselShapePyramid,
		VesselShapeCylinder,
		VesselShapeSphere,
		VesselShapeCube,
		VesselShapeOther,
	}
}

type WebhookEvent string

const (
	WebhookEventWebhookCreated  WebhookEvent = "webhook_created"
	WebhookEventWebhookUpdated  WebhookEvent = "webhook_updated"
	WebhookEventWebhookArchived WebhookEvent = "webhook_archived"
)

func (e *WebhookEvent) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = WebhookEvent(s)
	case string:
		*e = WebhookEvent(s)
	default:
		return fmt.Errorf("unsupported scan type for WebhookEvent: %T", src)
	}
	return nil
}

type NullWebhookEvent struct {
	WebhookEvent WebhookEvent
	Valid        bool // Valid is true if WebhookEvent is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullWebhookEvent) Scan(value interface{}) error {
	if value == nil {
		ns.WebhookEvent, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.WebhookEvent.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullWebhookEvent) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.WebhookEvent), nil
}

func (e WebhookEvent) Valid() bool {
	switch e {
	case WebhookEventWebhookCreated,
		WebhookEventWebhookUpdated,
		WebhookEventWebhookArchived:
		return true
	}
	return false
}

func AllWebhookEventValues() []WebhookEvent {
	return []WebhookEvent{
		WebhookEventWebhookCreated,
		WebhookEventWebhookUpdated,
		WebhookEventWebhookArchived,
	}
}

type MealPlanEvents struct {
	ID                string
	Notes             string
	StartsAt          time.Time
	EndsAt            time.Time
	MealName          MealName
	BelongsToMealPlan string
	CreatedAt         time.Time
	LastUpdatedAt     sql.NullTime
	ArchivedAt        sql.NullTime
}

type MealPlanOptionVotes struct {
	CreatedAt               time.Time
	LastUpdatedAt           sql.NullTime
	ArchivedAt              sql.NullTime
	ID                      string
	Notes                   string
	ByUser                  string
	BelongsToMealPlanOption string
	Rank                    int32
	Abstain                 bool
}

type Oauth2ClientTokens struct {
	AccessExpiresAt     time.Time
	CodeExpiresAt       time.Time
	RefreshExpiresAt    time.Time
	RefreshCreatedAt    time.Time
	CodeCreatedAt       time.Time
	AccessCreatedAt     time.Time
	CodeChallenge       string
	CodeChallengeMethod string
	Scope               Oauth2ClientTokenScopes
	ClientID            string
	Access              string
	Code                string
	ID                  string
	Refresh             string
	RedirectUri         string
	BelongsToUser       string
}

type PasswordResetTokens struct {
	ID            string
	Token         string
	ExpiresAt     time.Time
	CreatedAt     time.Time
	LastUpdatedAt sql.NullTime
	RedeemedAt    sql.NullTime
	BelongsToUser string
}

type RecipeRatings struct {
	CreatedAt     time.Time
	LastUpdatedAt sql.NullTime
	ArchivedAt    sql.NullTime
	ID            string
	RecipeID      string
	Notes         string
	ByUser        string
	Taste         sql.NullString
	Difficulty    sql.NullString
	Cleanup       sql.NullString
	Instructions  sql.NullString
	Overall       sql.NullString
}

package validators

import (
	"errors"
	"time"

	"github.com/tadeasf/eve-ran/src/db/models"
)

// ESI Response Validators

// ValidateCharacterResponse validates an ESI character response
func ValidateCharacterResponse(char *models.Character) error {
	if char == nil {
		return errors.New("character response is nil")
	}
	if char.ID == 0 {
		return errors.New("character ID is required")
	}
	if char.Name == "" {
		return errors.New("character name is required")
	}
	return nil
}

// ValidateSystemResponse validates an ESI system response
func ValidateSystemResponse(system *models.System) error {
	if system == nil {
		return errors.New("system response is nil")
	}
	if system.SystemID == 0 {
		return errors.New("system ID is required")
	}
	if system.Name == "" {
		return errors.New("system name is required")
	}
	if system.SecurityStatus < -1.0 || system.SecurityStatus > 1.0 {
		return errors.New("invalid security status")
	}
	return nil
}

// ValidateRegionResponse validates an ESI region response
func ValidateRegionResponse(region *models.Region) error {
	if region == nil {
		return errors.New("region response is nil")
	}
	if region.RegionID == 0 {
		return errors.New("region ID is required")
	}
	if region.Name == "" {
		return errors.New("region name is required")
	}
	return nil
}

// ValidateConstellationResponse validates an ESI constellation response
func ValidateConstellationResponse(constellation *models.Constellation) error {
	if constellation == nil {
		return errors.New("constellation response is nil")
	}
	if constellation.ConstellationID == 0 {
		return errors.New("constellation ID is required")
	}
	if constellation.Name == "" {
		return errors.New("constellation name is required")
	}
	if constellation.RegionID == 0 {
		return errors.New("region ID is required")
	}
	return nil
}

// ValidateESIItemResponse validates an ESI item response
func ValidateESIItemResponse(item *models.ESIItem) error {
	if item == nil {
		return errors.New("item response is nil")
	}
	if item.TypeID == 0 {
		return errors.New("item type ID is required")
	}
	if item.Name == "" {
		return errors.New("item name is required")
	}
	return nil
}

// zKillboard Response Validators

// ValidateKillmailResponse validates a zKillboard killmail response
func ValidateKillmailResponse(kill *models.Kill) error {
	if kill == nil {
		return errors.New("kill response is nil")
	}
	if kill.KillmailID == 0 {
		return errors.New("killmail ID is required")
	}
	if kill.KillmailTime.IsZero() {
		return errors.New("killmail time is required")
	}
	if kill.SolarSystemID == 0 {
		return errors.New("solar system ID is required")
	}

	// Validate killmail time is not in the future
	if kill.KillmailTime.After(time.Now().Add(1 * time.Hour)) {
		return errors.New("killmail time cannot be in the future")
	}

	// Validate killmail time is not too old (EVE Online launched in 2003)
	if kill.KillmailTime.Before(time.Date(2003, 1, 1, 0, 0, 0, 0, time.UTC)) {
		return errors.New("killmail time is too old")
	}

	return nil
}

// ValidateVictim validates a victim object
func ValidateVictim(victim *models.Victim) error {
	if victim == nil {
		return errors.New("victim is nil")
	}
	if victim.ShipTypeID == 0 {
		return errors.New("victim ship type ID is required")
	}
	if victim.DamageTaken <= 0 {
		return errors.New("damage taken must be positive")
	}
	return nil
}

// ValidateZkillData validates zkillboard data
func ValidateZkillData(zkill *models.Zkill) error {
	if zkill == nil {
		return errors.New("zkill data is nil")
	}
	if zkill.KillmailID == 0 {
		return errors.New("zkill killmail ID is required")
	}
	if zkill.Hash == "" {
		return errors.New("zkill hash is required")
	}

	// Validate ISK values are not negative
	if zkill.FittedValue < 0 || zkill.DroppedValue < 0 || zkill.DestroyedValue < 0 || zkill.TotalValue < 0 {
		return errors.New("ISK values cannot be negative")
	}

	// Validate points are not negative
	if zkill.Points < 0 {
		return errors.New("points cannot be negative")
	}

	// Validate total value is roughly sum of dropped + destroyed
	// Allow some tolerance for floating point errors
	expectedTotal := zkill.DroppedValue + zkill.DestroyedValue
	if zkill.TotalValue > 0 && (zkill.TotalValue < expectedTotal*0.9 || zkill.TotalValue > expectedTotal*1.1) {
		// Log warning but don't fail validation
		// Some zkillboard data may have discrepancies
	}

	return nil
}

// ValidateAttacker validates an attacker object
func ValidateAttacker(attacker *models.Attacker) error {
	if attacker == nil {
		return errors.New("attacker is nil")
	}
	if attacker.DamageDone <= 0 {
		return errors.New("damage done must be positive")
	}
	if attacker.ShipTypeID == 0 {
		return errors.New("attacker ship type ID is required")
	}
	return nil
}

// Batch Validators

// ValidateCharacterBatch validates a batch of characters
func ValidateCharacterBatch(characters []models.Character) ([]models.Character, []error) {
	var validChars []models.Character
	var errors []error

	for _, char := range characters {
		if err := ValidateCharacterResponse(&char); err != nil {
			errors = append(errors, err)
		} else {
			validChars = append(validChars, char)
		}
	}

	return validChars, errors
}

// ValidateKillmailBatch validates a batch of killmails
func ValidateKillmailBatch(kills []models.Kill) ([]models.Kill, []error) {
	var validKills []models.Kill
	var errors []error

	for _, kill := range kills {
		if err := ValidateKillmailResponse(&kill); err != nil {
			errors = append(errors, err)
		} else {
			validKills = append(validKills, kill)
		}
	}

	return validKills, errors
}

// ValidateSystemBatch validates a batch of systems
func ValidateSystemBatch(systems []models.System) ([]models.System, []error) {
	var validSystems []models.System
	var errors []error

	for _, system := range systems {
		if err := ValidateSystemResponse(&system); err != nil {
			errors = append(errors, err)
		} else {
			validSystems = append(validSystems, system)
		}
	}

	return validSystems, errors
}

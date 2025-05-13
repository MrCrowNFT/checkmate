package utils

import (
	"checkmate/api/internal/model"
)

// ConvertToSafeCredential converts a PlatformCredential to SafeCredential
func ConvertToSafeCredential(cred *model.PlatformCredential) model.SafeCredential {
	return model.SafeCredential{
		ID:        cred.ID,
		UserID:    cred.UserID,
		Platform:  cred.Platform,
		Name:      cred.Name,
		CreatedAt: cred.CreatedAt,
	}
}

// ConvertToSafeCredentials converts a slice of PlatformCredential to SafeCredential
func ConvertToSafeCredentials(creds []*model.PlatformCredential) []model.SafeCredential {
	safeCreds := make([]model.SafeCredential, len(creds))
	for i, cred := range creds {
		safeCreds[i] = ConvertToSafeCredential(cred)
	}
	return safeCreds
}

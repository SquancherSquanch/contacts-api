package connectors

import (
	"github.com/squanchersquanch/contacts/lib"
	"github.com/squanchersquanch/contacts/models"
)

// GenerateNewContact generates a new Contact fixture
func GenerateNewContact(filePath string) *models.Contact {
	contact := &models.Contact{}
	lib.MustUnMarshalFromFile(filePath, contact)
	return contact
}

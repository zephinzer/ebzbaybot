package collection

import (
	"fmt"
	"strings"

	"github.com/zephinzer/ebzbaybot/pkg/constants"
)

func GetCollectionByIdentifier(identifier string) (*Collection, error) {
	collectionNames := constants.CollectionByAddress[identifier]
	collectionAddresses := constants.CollectionByName[identifier]
	isCollectionNameValid := collectionNames != nil && len(collectionNames) > 0
	isCollectionAddressValid := collectionAddresses != nil && len(collectionAddresses) > 0

	aliases := []string{}
	var collectionInstance *Collection = nil
	if !isCollectionNameValid && !isCollectionAddressValid {
		return nil, fmt.Errorf("failed to receive a known identifier")
	} else if isCollectionNameValid { // address provided
		if len(collectionNames) > 1 {
			aliases = collectionNames[1:]
		}
		collectionInstance = &Collection{
			ID:    identifier,
			Label: collectionNames[0],
		}
	} else if isCollectionAddressValid { // name provided
		address := collectionAddresses[0]
		primaryName := constants.CollectionByAddress[address][0]
		if len(constants.CollectionByAddress[address]) > 1 {
			aliases = constants.CollectionByAddress[address][1:]
		}
		collectionInstance = &Collection{
			ID:    address,
			Label: primaryName,
		}
	}

	collectionInstance.Aliases = strings.Join(aliases, ",")
	if collectionInstance != nil {
		return collectionInstance, nil
	}
	return nil, fmt.Errorf("failed to behave as expected")
}

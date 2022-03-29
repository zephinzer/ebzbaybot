package collection

import (
	"fmt"

	"github.com/zephinzer/ebzbaybot/pkg/constants"
)

type Collection struct {
	Name    string
	Address string
}

func GetCollectionByIdentifier(identifier string) (*Collection, error) {
	collectionNames := constants.CollectionByAddress[identifier]
	collectionAddresses := constants.CollectionByName[identifier]
	isCollectionNameValid := collectionNames != nil && len(collectionNames) > 0
	isCollectionAddressValid := collectionAddresses != nil && len(collectionAddresses) > 0
	fmt.Println(collectionNames)
	fmt.Println(collectionAddresses)
	if !isCollectionNameValid && !isCollectionAddressValid {
		return nil, fmt.Errorf("failed to receive a known identifier")
	} else if isCollectionNameValid { // address provided
		return &Collection{
			Name:    collectionNames[0],
			Address: identifier,
		}, nil
	} else if isCollectionAddressValid { // name provided
		address := collectionAddresses[0]
		primaryName := constants.CollectionByAddress[address][0]
		return &Collection{
			Address: address,
			Name:    primaryName,
		}, nil
	}
	return nil, fmt.Errorf("failed to behave as expected")
}

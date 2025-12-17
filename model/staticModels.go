package model

import "github.com/greeneg/system-confd/globals"

func GetRootPath() (MetaData, error) {
	metaData := MetaData{
		Name:      globals.NAME,
		Version:   globals.VERSION,
		Copyright: globals.COPYRIGHT,
		License:   globals.LICENSE,
		Author:    globals.AUTHOR,
	}
	return metaData, nil
}

func GetVersion() string {
	return globals.VERSION
}

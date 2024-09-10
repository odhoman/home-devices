package common

const (
	ErrDeviceNotCreatedErrorCode    = "DEVICE_NO_CREATED"
	ErrDeviceNotCreatedErrorMessage = "An error occurred creating new device"

	ErrDeviceAlreadyExistsCode    = "DEVICE_ALREADY_EXISTS"
	ErrDeviceAlreadyExistsMessage = "Device already exist for that mac and homeId"

	ErrGettingDeviceCode    = "ERROR_GETTING_DEVICE"
	ErrGettingDeviceMessage = "An error ocurred getting the device"

	ErrNoFieldToUpdateCode    = "NO_FIELDS_TO_UPDATE"
	ErrNoFieldToUpdateMessage = "Please enter a value property to update"

	ErrUpdatingDeviceCode    = "ERROR_UPDATING_DEVICE"
	ErrUpdatingDeviceMessage = "An error occurred updating a device"

	ErrDeviceNotFoundCode    = "ERROR_DEVICE_NOT_FOUND"
	ErrDeviceNotFoundMessage = "Device Not Found"

	ErrDeletingDeviceCode    = "ERROR_DELETING_DEVICE"
	ErrDeletingDeviceMessage = "An error occurred deleting a device"

	ErrGettingConfigCode    = "ERROR_GETTING_CONFIG"
	ErrGettingConfigMessage = "An error occurred deleting a device"

	InternalServerErrorDefaultBodyResponse = "{\"errors\": [\"Internal Server Error\"]}"

	ResponseOKWithMessageTemplate = "{\"message\": \"%v\"}"

	TableNameHomeDevicesProperty = "HOME_DEVICE_TABLE_NAME"
	MacHomeIdIndexNameProperty   = "MAC_HOMEID_INDEX_NAME"
)

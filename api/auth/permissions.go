package auth

type PermissionFlags int64

const (
	// SYSTEM
	ADMIN PermissionFlags = 1 << 0

	// SERVICES
	READ_SERVICES         PermissionFlags = 1 << 1
	EDIT_SERVICES_CONFIGS PermissionFlags = 1 << 2
	CREATE_NEW_SERVICES   PermissionFlags = 1 << 3
	DELETE_SERVICES       PermissionFlags = 1 << 4
	MANAGE_SERVICES_STATE PermissionFlags = 1 << 5
	UPDATE_SERVICES       PermissionFlags = 1 << 6

	// FILES
	READ_ALL_FILES   PermissionFlags = 1 << 7
	EDIT_ALL_FILES   PermissionFlags = 1 << 8
	DELETE_ALL_FILES PermissionFlags = 1 << 9

	// SETTINGS
	READ_APP_CONFIG PermissionFlags = 1 << 10
	EDIT_APP_CONFIG PermissionFlags = 1 << 11

	// DATABASE
	READ_ALL_DATABASES   PermissionFlags = 1 << 12
	UPDATE_ALL_DATABASES PermissionFlags = 1 << 13
	DELETE_ALL_DATABASES PermissionFlags = 1 << 14

	// USERS
	SEE_USERS_DETAILS  PermissionFlags = 1 << 15
	EDIT_USERS_DETAILS PermissionFlags = 1 << 16
	CREATE_USERS       PermissionFlags = 1 << 17
	REMOVE_USERS       PermissionFlags = 1 << 18
	MANAGE_USERS       PermissionFlags = 1 << 19

	// ROLES
	EDIT_ROLES   PermissionFlags = 1 << 20
	CREATE_ROLES PermissionFlags = 1 << 21
	DELETE_ROLES PermissionFlags = 1 << 22
)

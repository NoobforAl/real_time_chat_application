package taskTypes

const (
	ActionCreate      = "create"
	ActionArchive     = "archive"
	ActionReportDaily = "report_daily"

	TypeMessage            = "message:"
	TypeMessageSave        = TypeMessage + ActionCreate
	TypeMessageArchive     = TypeMessage + ActionArchive
	TypeMessageReportDaily = TypeMessage + ActionReportDaily

	TypeRoom     = "room:"
	TypeRoomSave = TypeRoom + ActionCreate

	TypeNotification     = "notification:"
	TypeNotificationSave = TypeNotification + ActionCreate
)

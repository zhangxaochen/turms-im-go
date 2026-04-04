package property

type TurmsProperties struct {
	Service ServiceProperties
}

type ConferenceProperties struct {
	Meeting MeetingProperties
}

type MeetingProperties struct {
	AllowCancel bool
	ID          MeetingIdProperties
	Name        MeetingFieldProperties
	Intro       MeetingFieldProperties
	Password    MeetingPasswordProperties
	Quota       MeetingQuotaProperties
	Scheduling  MeetingSchedulingProperties
}

type MeetingIdProperties struct {
	Type int // 0 for DIGIT_9, 1 for DIGIT_10
}

type MeetingFieldProperties struct {
	MinLength int
	MaxLength int
}

type MeetingPasswordProperties struct {
	MinLength int
	MaxLength int
}

type MeetingQuotaProperties struct {
	MaxActiveMeetingCountPerUser int
}

type MeetingSchedulingProperties struct {
	MaxAllowedStartDateOffsetSeconds int64
	AllowCancel                      bool
}

type ServiceProperties struct {
	AdminApi     AdminApiProperties
	Conference   ConferenceProperties
	Conversation NotificationConversationProperties
	Notification NotificationProperties
}

type NotificationConversationProperties struct {
	ReadReceipt  ReadReceiptProperties
	TypingStatus TypingStatusProperties
}

type ReadReceiptProperties struct {
	Enabled                  bool
	AllowMoveReadDateForward bool
	UseServerTime            bool
}

type TypingStatusProperties struct {
	Enabled bool
}

type NotificationProperties struct {
	MeetingCanceled          NotificationMeetingCanceledProperties
	MeetingUpdated           NotificationMeetingUpdatedProperties
	MeetingInvitationUpdated NotificationMeetingInvitationUpdatedProperties
	PrivateConversationSettingDeleted NotificationPrivateConversationSettingDeletedProperties
	PrivateConversationSettingUpdated NotificationPrivateConversationSettingUpdatedProperties
	GroupConversationSettingDeleted   NotificationGroupConversationSettingDeletedProperties
	GroupConversationSettingUpdated   NotificationGroupConversationSettingUpdatedProperties
	PrivateConversationReadDateUpdated NotificationPrivateConversationReadDateUpdatedProperties
	GroupConversationReadDateUpdated   NotificationGroupConversationReadDateUpdatedProperties
}

type NotificationPrivateConversationReadDateUpdatedProperties struct {
	NotifyRequesterOtherOnlineSessions bool
	NotifyContact                      bool
}

type NotificationGroupConversationReadDateUpdatedProperties struct {
	NotifyRequesterOtherOnlineSessions bool
	NotifyOtherGroupMembers           bool
}

type NotificationPrivateConversationSettingDeletedProperties struct {
	NotifyRequesterOtherOnlineSessions bool
}

type NotificationPrivateConversationSettingUpdatedProperties struct {
	NotifyRequesterOtherOnlineSessions bool
}

type NotificationGroupConversationSettingDeletedProperties struct {
	NotifyRequesterOtherOnlineSessions bool
}

type NotificationGroupConversationSettingUpdatedProperties struct {
	NotifyRequesterOtherOnlineSessions bool
}

type NotificationMeetingCanceledProperties struct {
	NotifyRequesterOtherOnlineSessions bool
	NotifyMeetingParticipants          bool
}

type NotificationMeetingUpdatedProperties struct {
	NotifyRequesterOtherOnlineSessions bool
	NotifyMeetingParticipants          bool
}

type NotificationMeetingInvitationUpdatedProperties struct {
	NotifyRequesterOtherOnlineSessions bool
	NotifyMeetingParticipants          bool
}

type AdminApiProperties struct {
	MaxDayDifferencePerRequest              int
	MaxHourDifferencePerCountRequest        int
	MaxDayDifferencePerCountRequest         int
	MaxMonthDifferencePerCountRequest       int
	MaxAvailableRecordsPerRequest           int
	MaxAvailableOnlineUsersStatusPerRequest int
	DefaultAvailableRecordsPerRequest       int
}

type TurmsPropertiesManager struct {
	properties TurmsProperties
}

func NewTurmsPropertiesManager() *TurmsPropertiesManager {
	return &TurmsPropertiesManager{
		properties: TurmsProperties{
			Service: ServiceProperties{
				AdminApi: AdminApiProperties{
					MaxDayDifferencePerRequest:        90,
					MaxHourDifferencePerCountRequest:  24,
					MaxDayDifferencePerCountRequest:   31,
					MaxMonthDifferencePerCountRequest: 12,
					MaxAvailableRecordsPerRequest:     1000,
					DefaultAvailableRecordsPerRequest: 10,
				},
				Conference: ConferenceProperties{
					Meeting: MeetingProperties{
						AllowCancel: true,
						ID: MeetingIdProperties{
							Type: 0,
						},
						Name: MeetingFieldProperties{
							MinLength: 1,
							MaxLength: 64,
						},
						Intro: MeetingFieldProperties{
							MinLength: 0,
							MaxLength: 255,
						},
						Password: MeetingPasswordProperties{
							MinLength: 0,
							MaxLength: 64,
						},
						Quota: MeetingQuotaProperties{
							MaxActiveMeetingCountPerUser: 10,
						},
						Scheduling: MeetingSchedulingProperties{
							MaxAllowedStartDateOffsetSeconds: 86400,
							AllowCancel:                      true,
						},
					},
				},
				Conversation: NotificationConversationProperties{
					ReadReceipt: ReadReceiptProperties{
						Enabled:                  true,
						AllowMoveReadDateForward: true,
						UseServerTime:            true,
					},
					TypingStatus: TypingStatusProperties{
						Enabled: true,
					},
				},
				Notification: NotificationProperties{
					MeetingCanceled: NotificationMeetingCanceledProperties{
						NotifyRequesterOtherOnlineSessions: true,
						NotifyMeetingParticipants:          true,
					},
					MeetingUpdated: NotificationMeetingUpdatedProperties{
						NotifyRequesterOtherOnlineSessions: true,
						NotifyMeetingParticipants:          true,
					},
					MeetingInvitationUpdated: NotificationMeetingInvitationUpdatedProperties{
						NotifyRequesterOtherOnlineSessions: true,
						NotifyMeetingParticipants:          true,
					},
					PrivateConversationSettingDeleted: NotificationPrivateConversationSettingDeletedProperties{
						NotifyRequesterOtherOnlineSessions: true,
					},
					PrivateConversationSettingUpdated: NotificationPrivateConversationSettingUpdatedProperties{
						NotifyRequesterOtherOnlineSessions: true,
					},
					GroupConversationSettingDeleted: NotificationGroupConversationSettingDeletedProperties{
						NotifyRequesterOtherOnlineSessions: true,
					},
					GroupConversationSettingUpdated: NotificationGroupConversationSettingUpdatedProperties{
						NotifyRequesterOtherOnlineSessions: true,
					},
					PrivateConversationReadDateUpdated: NotificationPrivateConversationReadDateUpdatedProperties{
						NotifyRequesterOtherOnlineSessions: true,
						NotifyContact:                      true,
					},
					GroupConversationReadDateUpdated: NotificationGroupConversationReadDateUpdatedProperties{
						NotifyRequesterOtherOnlineSessions: true,
						NotifyOtherGroupMembers:           true,
					},
				},
			},
		},
	}
}

func (m *TurmsPropertiesManager) GetLocalProperties() *TurmsProperties {
	return &m.properties
}

func (m *TurmsPropertiesManager) AddGlobalPropertiesChangeListener(listener func(*TurmsProperties)) {
	listener(&m.properties)
}

func (m *TurmsPropertiesManager) NotifyAndAddGlobalPropertiesChangeListener(listener func(*TurmsProperties)) {
	listener(&m.properties)
}

func (m *TurmsPropertiesManager) GetGlobalProperties() *TurmsProperties {
	return &m.properties
}

func (m *TurmsPropertiesManager) UpdateLocalProperties(reset bool, turmsProperties map[string]interface{}) error {
	// mock implementation
	return nil
}

func (m *TurmsPropertiesManager) UpdateGlobalProperties(reset bool, turmsProperties map[string]interface{}) error {
	// mock implementation
	return nil
}

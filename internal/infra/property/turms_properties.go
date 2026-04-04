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
	Notification NotificationProperties
}

type NotificationProperties struct {
	MeetingCanceled          NotificationMeetingCanceledProperties
	MeetingUpdated           NotificationMeetingUpdatedProperties
	MeetingInvitationUpdated NotificationMeetingInvitationUpdatedProperties
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

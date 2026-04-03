package permission

type AdminPermission string

const (
	FlightRecordingCreate AdminPermission = "FLIGHT_RECORDING_CREATE"
	FlightRecordingDelete AdminPermission = "FLIGHT_RECORDING_DELETE"
	FlightRecordingUpdate AdminPermission = "FLIGHT_RECORDING_UPDATE"
	FlightRecordingQuery  AdminPermission = "FLIGHT_RECORDING_QUERY"

	// region business
	UserCreate AdminPermission = "USER_CREATE"
	UserDelete AdminPermission = "USER_DELETE"
	UserUpdate AdminPermission = "USER_UPDATE"
	UserQuery  AdminPermission = "USER_QUERY"

	UserRelationshipCreate AdminPermission = "USER_RELATIONSHIP_CREATE"
	UserRelationshipDelete AdminPermission = "USER_RELATIONSHIP_DELETE"
	UserRelationshipUpdate AdminPermission = "USER_RELATIONSHIP_UPDATE"
	UserRelationshipQuery  AdminPermission = "USER_RELATIONSHIP_QUERY"

	UserRelationshipGroupCreate AdminPermission = "USER_RELATIONSHIP_GROUP_CREATE"
	UserRelationshipGroupDelete AdminPermission = "USER_RELATIONSHIP_GROUP_DELETE"
	UserRelationshipGroupUpdate AdminPermission = "USER_RELATIONSHIP_GROUP_UPDATE"
	UserRelationshipGroupQuery  AdminPermission = "USER_RELATIONSHIP_GROUP_QUERY"

	UserFriendRequestCreate AdminPermission = "USER_FRIEND_REQUEST_CREATE"
	UserFriendRequestDelete AdminPermission = "USER_FRIEND_REQUEST_DELETE"
	UserFriendRequestUpdate AdminPermission = "USER_FRIEND_REQUEST_UPDATE"
	UserFriendRequestQuery  AdminPermission = "USER_FRIEND_REQUEST_QUERY"

	UserRoleCreate AdminPermission = "USER_ROLE_CREATE"
	UserRoleDelete AdminPermission = "USER_ROLE_DELETE"
	UserRoleUpdate AdminPermission = "USER_ROLE_UPDATE"
	UserRoleQuery  AdminPermission = "USER_ROLE_QUERY"

	UserOnlineInfoUpdate AdminPermission = "USER_ONLINE_INFO_UPDATE"
	UserOnlineInfoQuery  AdminPermission = "USER_ONLINE_INFO_QUERY"

	GroupCreate AdminPermission = "GROUP_CREATE"
	GroupDelete AdminPermission = "GROUP_DELETE"
	GroupUpdate AdminPermission = "GROUP_UPDATE"
	GroupQuery  AdminPermission = "GROUP_QUERY"

	GroupBlocklistCreate AdminPermission = "GROUP_BLOCKLIST_CREATE"
	GroupBlocklistDelete AdminPermission = "GROUP_BLOCKLIST_DELETE"
	GroupBlocklistUpdate AdminPermission = "GROUP_BLOCKLIST_UPDATE"
	GroupBlocklistQuery  AdminPermission = "GROUP_BLOCKLIST_QUERY"

	GroupInvitationCreate AdminPermission = "GROUP_INVITATION_CREATE"
	GroupInvitationDelete AdminPermission = "GROUP_INVITATION_DELETE"
	GroupInvitationUpdate AdminPermission = "GROUP_INVITATION_UPDATE"
	GroupInvitationQuery  AdminPermission = "GROUP_INVITATION_QUERY"

	GroupQuestionCreate AdminPermission = "GROUP_QUESTION_CREATE"
	GroupQuestionDelete AdminPermission = "GROUP_QUESTION_DELETE"
	GroupQuestionUpdate AdminPermission = "GROUP_QUESTION_UPDATE"
	GroupQuestionQuery  AdminPermission = "GROUP_QUESTION_QUERY"

	GroupJoinRequestCreate AdminPermission = "GROUP_JOIN_REQUEST_CREATE"
	GroupJoinRequestDelete AdminPermission = "GROUP_JOIN_REQUEST_DELETE"
	GroupJoinRequestUpdate AdminPermission = "GROUP_JOIN_REQUEST_UPDATE"
	GroupJoinRequestQuery  AdminPermission = "GROUP_JOIN_REQUEST_QUERY"

	GroupMemberUpdate AdminPermission = "GROUP_MEMBER_UPDATE"
	GroupMemberCreate AdminPermission = "GROUP_MEMBER_CREATE"
	GroupMemberDelete AdminPermission = "GROUP_MEMBER_DELETE"
	GroupMemberQuery  AdminPermission = "GROUP_MEMBER_QUERY"

	GroupTypeCreate AdminPermission = "GROUP_TYPE_CREATE"
	GroupTypeDelete AdminPermission = "GROUP_TYPE_DELETE"
	GroupTypeUpdate AdminPermission = "GROUP_TYPE_UPDATE"
	GroupTypeQuery  AdminPermission = "GROUP_TYPE_QUERY"

	ConversationQuery  AdminPermission = "CONVERSATION_QUERY"
	ConversationDelete AdminPermission = "CONVERSATION_DELETE"
	ConversationUpdate AdminPermission = "CONVERSATION_UPDATE"

	MessageCreate AdminPermission = "MESSAGE_CREATE"
	MessageDelete AdminPermission = "MESSAGE_DELETE"
	MessageUpdate AdminPermission = "MESSAGE_UPDATE"
	MessageQuery  AdminPermission = "MESSAGE_QUERY"
	// endregion

	// region session
	SessionDelete AdminPermission = "SESSION_DELETE"
	// endregion

	// region business - statistics
	StatisticsUserQuery    AdminPermission = "STATISTICS_USER_QUERY"
	StatisticsGroupQuery   AdminPermission = "STATISTICS_GROUP_QUERY"
	StatisticsMessageQuery AdminPermission = "STATISTICS_MESSAGE_QUERY"
	// endregion

	// region admin
	AdminCreate AdminPermission = "ADMIN_CREATE"
	AdminDelete AdminPermission = "ADMIN_DELETE"
	AdminUpdate AdminPermission = "ADMIN_UPDATE"
	AdminQuery  AdminPermission = "ADMIN_QUERY"

	AdminRoleCreate AdminPermission = "ADMIN_ROLE_CREATE"
	AdminRoleDelete AdminPermission = "ADMIN_ROLE_DELETE"
	AdminRoleUpdate AdminPermission = "ADMIN_ROLE_UPDATE"
	AdminRoleQuery  AdminPermission = "ADMIN_ROLE_QUERY"

	AdminPermissionQuery AdminPermission = "ADMIN_PERMISSION_QUERY"
	// endregion

	// region client - blocklist
	ClientBlocklistCreate AdminPermission = "CLIENT_BLOCKLIST_CREATE"
	ClientBlocklistDelete AdminPermission = "CLIENT_BLOCKLIST_DELETE"
	ClientBlocklistQuery  AdminPermission = "CLIENT_BLOCKLIST_QUERY"
	// endregion

	// region client - request
	ClientRequestCreate AdminPermission = "CLIENT_REQUEST_CREATE"
	// endregion

	// region cluster
	ClusterMemberCreate AdminPermission = "CLUSTER_MEMBER_CREATE"
	ClusterMemberDelete AdminPermission = "CLUSTER_MEMBER_DELETE"
	ClusterMemberUpdate AdminPermission = "CLUSTER_MEMBER_UPDATE"
	ClusterMemberQuery  AdminPermission = "CLUSTER_MEMBER_QUERY"

	ClusterLeaderUpdate AdminPermission = "CLUSTER_LEADER_UPDATE"
	ClusterLeaderQuery  AdminPermission = "CLUSTER_LEADER_QUERY"

	ClusterSettingUpdate AdminPermission = "CLUSTER_SETTING_UPDATE"
	ClusterSettingQuery  AdminPermission = "CLUSTER_SETTING_QUERY"
	// endregion

	// region node - plugin
	PluginCreate AdminPermission = "PLUGIN_CREATE"
	PluginDelete AdminPermission = "PLUGIN_DELETE"
	PluginUpdate AdminPermission = "PLUGIN_UPDATE"
	PluginQuery  AdminPermission = "PLUGIN_QUERY"
	// endregion

	// region node - log
	LogQuery AdminPermission = "LOG_QUERY"
	// endregion

	// region node - others
	Shutdown AdminPermission = "SHUTDOWN"
	// endregion
)

var AllAdminPermissions = []AdminPermission{
	FlightRecordingCreate, FlightRecordingDelete, FlightRecordingUpdate, FlightRecordingQuery,

	UserCreate, UserDelete, UserUpdate, UserQuery,
	UserRelationshipCreate, UserRelationshipDelete, UserRelationshipUpdate, UserRelationshipQuery,
	UserRelationshipGroupCreate, UserRelationshipGroupDelete, UserRelationshipGroupUpdate, UserRelationshipGroupQuery,
	UserFriendRequestCreate, UserFriendRequestDelete, UserFriendRequestUpdate, UserFriendRequestQuery,
	UserRoleCreate, UserRoleDelete, UserRoleUpdate, UserRoleQuery,
	UserOnlineInfoUpdate, UserOnlineInfoQuery,

	GroupCreate, GroupDelete, GroupUpdate, GroupQuery,
	GroupBlocklistCreate, GroupBlocklistDelete, GroupBlocklistUpdate, GroupBlocklistQuery,
	GroupInvitationCreate, GroupInvitationDelete, GroupInvitationUpdate, GroupInvitationQuery,
	GroupQuestionCreate, GroupQuestionDelete, GroupQuestionUpdate, GroupQuestionQuery,
	GroupJoinRequestCreate, GroupJoinRequestDelete, GroupJoinRequestUpdate, GroupJoinRequestQuery,
	GroupMemberUpdate, GroupMemberCreate, GroupMemberDelete, GroupMemberQuery,
	GroupTypeCreate, GroupTypeDelete, GroupTypeUpdate, GroupTypeQuery,
	ConversationQuery, ConversationDelete, ConversationUpdate,
	MessageCreate, MessageDelete, MessageUpdate, MessageQuery,

	SessionDelete,

	StatisticsUserQuery, StatisticsGroupQuery, StatisticsMessageQuery,

	AdminCreate, AdminDelete, AdminUpdate, AdminQuery,
	AdminRoleCreate, AdminRoleDelete, AdminRoleUpdate, AdminRoleQuery,
	AdminPermissionQuery,

	ClientBlocklistCreate, ClientBlocklistDelete, ClientBlocklistQuery,
	ClientRequestCreate,

	ClusterMemberCreate, ClusterMemberDelete, ClusterMemberUpdate, ClusterMemberQuery,
	ClusterLeaderUpdate, ClusterLeaderQuery,
	ClusterSettingUpdate, ClusterSettingQuery,

	PluginCreate, PluginDelete, PluginUpdate, PluginQuery,
	LogQuery,
	Shutdown,
}

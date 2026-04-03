const fs = require('fs');
let content = fs.readFileSync('docs/refactor_progress_report.md', 'utf8');

const mappings = {
  'decode(BerBuffer buffer)': '[internal/infra/ldap/element/elements.go:Decode(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go)',
  'estimateSize()': '[internal/infra/ldap/element/elements.go:EstimateSize()](../internal/infra/ldap/element/elements.go)',
  'writeTo(BerBuffer buffer)': '[internal/infra/ldap/element/elements.go:WriteTo(buffer *asn1.BerBuffer)](../internal/infra/ldap/element/elements.go)',
  'mongoDataGenerator()': '[internal/infra/mongo/mongo_data_generator.go:NewMongoDataGenerator()](../internal/infra/mongo/mongo_data_generator.go)'
};

for (const [key, val] of Object.entries(mappings)) {
    content = content.replace(new RegExp('- \\\\[ \\\\] \\\\`' + key.replace(/(\\.|\\(|\\)|\\*)/g, '\\\\$1') + '\\\\`', 'g'), '- [x] `' + key + '` -> ' + val);
}

const dtoVal = '[internal/domain/common/dto/request_handler_result.go:RequestHandlerResult](../internal/domain/common/dto/request_handler_result.go)';
const reqHandlerResultLines = [
  'of(@NotNull ResponseStatusCode code)',
  'of(@NotNull ResponseStatusCode code, @Nullable String reason)',
  'of(@NotNull TurmsNotification.Data response)',
  'of(boolean forwardNotificationsToRequesterOtherOnlineSessions, @NotNull TurmsRequest notification)',
  'of(boolean forwardNotificationsToRequesterOtherOnlineSessions, @NotNull Long recipientId, @NotNull TurmsRequest notification)',
  'of(@NotNull Long recipientId, @NotNull TurmsRequest notification)',
  'of(@NotEmpty Set<Long> recipientIds, @NotNull TurmsRequest dataForRecipient)',
  'of(boolean forwardNotificationToRequesterOtherOnlineSessions, @NotEmpty Set<Long> recipientIds, @NotNull TurmsRequest notification)',
  'of(TurmsNotification.Data response, @NotEmpty Set<Long> recipientIds, @NotNull TurmsRequest notificationForRecipients, @NotNull TurmsRequest notificationForRequesterOtherOnlineSessions)',
  'of(@NotEmpty Set<Long> recipientIds, @NotNull TurmsRequest notificationForRecipients, @NotNull TurmsRequest notificationForRequesterOtherOnlineSessions)',
  'of(TurmsNotification.Data response, boolean forwardNotificationToRequesterOtherOnlineSessions, @NotEmpty Set<Long> recipientIds, @NotNull TurmsRequest notification)',
  'of(@NotNull ResponseStatusCode code, @NotNull Long recipientId, @NotNull TurmsRequest notification)',
  'of(@NotNull ResponseStatusCode code, boolean forwardNotificationToRequesterOtherOnlineSessions, @NotNull Long recipientId, @NotNull TurmsRequest notification)',
  'of(@NotNull List<Notification> notifications)',
  'of(@NotNull Notification notification)',
  'ofDataLong(@NotNull Long value)',
  'ofDataLong(@NotNull Long value, @NotNull Long recipientId, @NotNull TurmsRequest notification)',
  'ofDataLong(@NotNull Long value, boolean forwardNotificationToRequesterOtherOnlineSessions, @NotNull Long recipientId, @NotNull TurmsRequest notification)',
  'ofDataLong(@NotNull Long value, boolean forwardDataForRecipientsToRequesterOtherOnlineSessions, @NotNull TurmsRequest notification)',
  'ofDataLong(@NotNull Long value, boolean forwardNotificationsToRequesterOtherOnlineSessions, @NotEmpty Set<Long> recipients, TurmsRequest notification)',
  'ofDataLongs(@NotNull Collection<Long> values)',
  'Notification(boolean forwardToRequesterOtherOnlineSessions, Set<Long> recipients, TurmsRequest notification)',
  'of(boolean forwardToRequesterOtherOnlineSessions, Set<Long> recipients, TurmsRequest notification)',
  'of(boolean forwardToRequesterOtherOnlineSessions, Long recipient, TurmsRequest notification)',
  'of(boolean forwardToRequesterOtherOnlineSessions, TurmsRequest notification)'
];

for (const line of reqHandlerResultLines) {
    content = content.replace(new RegExp('- \\\\[ \\\\] \\\\`' + line.replace(/(\\.|\\(|\\)|\\*|\\@|\\<|\\>)/g, '\\\\$1') + '\\\\`', 'g'), '- [x] `' + line + '` -> ' + dtoVal);
}

content = content.replace(/- \[ \] `queryAdmins\(@QueryParam\(required = false\)`/g, '- [x] `queryAdmins(@QueryParam(required = false)` -> [internal/domain/admin/access/admin/controller/admin_controllers.go:QueryAdminsWithQuery](../internal/domain/admin/access/admin/controller/admin_controllers.go)');
content = content.replace(/- \[ \] `queryAdminRoles\(@QueryParam\(required = false\)`/g, '- [x] `queryAdminRoles(@QueryParam(required = false)` -> [internal/domain/admin/access/admin/controller/admin_controllers.go:QueryAdminRolesWithQuery](../internal/domain/admin/access/admin/controller/admin_controllers.go)');

fs.writeFileSync('docs/refactor_progress_report.md', content);

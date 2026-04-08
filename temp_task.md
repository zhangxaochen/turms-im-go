# Local Progress Tracker for Batch 3

- [x] **Repository uses upsert instead of insert**: The Go repository `AddGroupMember` (`group_member_repository.go:34-41`) uses `UpdateOne` with `SetUpsert(true)`, which silently overwrites an existing group member. The Java version uses `groupMemberRepository.insert()`, which fails if a member already exists, preserving data integrity.

---

- [x] **Null role silently defaults instead of erroring**: When `role` is nil in the DTO, the Go controller defaults to `MEMBER` (line 348). In Java, a null role would trigger `Validator.notNull(groupMemberRole, "groupMemberRole")` in the service, returning an error. The Go version silently assigns a default role instead of rejecting the request.

---

- [x] **Missing non-paged query with filters**: Java has a `@GetMapping` (non-paged) variant that accepts filter parameters (`groupIds`, `userIds`, `roles`, date ranges) and passes `page=0`. The Go `QueryGroupMembersWithQuery` accepts a `page` parameter but there is no corresponding non-paged controller method that passes `page=0` with filter parameters like the Java version does.

---

- [x] **Missing count query for paginated endpoint**: The Java paged variant (`@GetMapping("page")`) calls both `countMembers(...)` and `queryGroupMembers(...)`, returning a `PaginationDTO` with total count and data. The Go `QueryGroupMembersWithQuery` only queries members — it never calls a count method, so no total count is returned for pagination.

---

- [x] **Iterates one-by-one instead of batch update**: The Go controller iterates over each key and calls `UpdateGroupMember` individually (lines 318-332). The Java version converts the list to a set and calls `updateGroupMembers` as a single batch operation via `groupMemberRepository.updateGroupMembers(keys, ...)`. This is functionally different — individual updates are not atomic and if one fails mid-way, some members are updated and others are not.

---

- [x] **updateGroupMembersVersion parameter is false instead of true**: The Go controller passes `false` for `updateVersion` (line 327), meaning the group members version is never updated. The Java controller passes `true`, which triggers `groupVersionService.updateMembersVersion(groupId)` to notify clients of the change via the version mechanism.

---

- [x] **Incorrect additional deletion date filter**: The Java version only filters by `eq(Group.Fields.OWNER_ID, ownerId)` with no deletion date check. The Go version at line 133 adds `"dd": bson.M{"$exists": false}`, filtering out deleted groups, which does not match the Java behavior.

---

- [x] **Method not ported**: The Java version filters by both `OWNER_ID` and `TYPE_ID`. No equivalent method exists in the Go repository that accepts both `ownerID` and `groupTypeId` parameters.

---

- [x] **Method signature completely different**: The Java `findGroups` takes 11 nullable parameters (ids, typeIds, creatorIds, ownerIds, isActive, creationDateRange, deletionDateRange, lastUpdatedDateRange, muteEndDateRange, page, size) with no deletion date filter. The Go `FindGroups` at line 37 only takes `groupIDs []int64` and always adds `dd: $exists false`. The Go `QueryGroups` at line 57 takes different parameters (name, lastUpdatedDate, skip, limit). Neither matches the Java method's full filter capabilities.

---

- [x] **Incorrect default deletion date filter**: The Java `findGroups` does NOT filter by deletion date. The Go `FindGroups` always adds `dd: $exists false`, filtering out deleted groups that should be included.

---

- [x] **`null` equality vs `$exists: false` mismatch**: The Java version uses `eq(Group.Fields.DELETION_DATE, null)` which matches documents where the field is either absent or explicitly set to null. The Go version uses `"dd": bson.M{"$exists": false}` which only matches documents where the field is absent entirely, missing documents where `dd` is explicitly `null`.

---

- [x] **Method not ported**: The Java `findAllNames` uses `QueryOptions.include(Group.Fields.NAME)` projection to return only the name field of all groups. No equivalent method exists in the Go `GroupRepository`.

---

- [x] **`null` equality vs `$exists: false` mismatch**: Same as `findNotDeletedGroups` — Java uses `eq(DELETION_DATE, null)` while Go uses `"dd": bson.M{"$exists": false}`, missing documents where `dd` is explicitly `null`.

---

- [x] **Missing `muteEndDate` parameter / uses `time.Now()` instead**: The Java version takes `(Long groupId, Date muteEndDate)` and compares the stored `MUTE_END_DATE` field against the passed-in `muteEndDate`. The Go version at line 290 takes only `groupID int64` and hardcodes `time.Now()` as the comparison value. This changes the semantics from a caller-provided comparison point to always comparing against the current time.

---

- [x] **`null` equality vs `$exists: false` mismatch**: Java uses `eq(Group.Fields.DELETION_DATE, null)` while Go uses `"dd": bson.M{"$exists": false}`, missing documents where `dd` is explicitly `null`.
*Checked methods: updateVersions(String field), updateVersions(@Nullable Set<Long> groupIds, String field), updateVersion(Long groupId, String field), updateVersion(Long groupId, boolean updateMembers, boolean updateBlocklist, boolean joinRequests, boolean joinQuestions), findBlocklist(Long groupId), findInvitations(Long groupId), findJoinRequests(Long groupId), findJoinQuestions(Long groupId), findMembers(Long groupId)*
Now I have all the information needed for a thorough comparison. Let me analyze each method systematically.
**Java** (line 47-51): Updates ALL documents (empty filter) setting the given field to `new Date()`.
**Go** (line 140-148): `UpdateVersions` takes `groupIDs []int64` and `field string`. When `groupIDs` is empty, the filter is `bson.M{}` which matches all documents — correct parity with the Java version that has no filter. However, the Java version has a separate single-parameter overload `updateVersions(String field)` with an empty filter. The Go version requires `groupIDs` to be explicitly passed — if called with an empty slice it behaves correctly (matches all). **This is fine.**
**Java** (line 53-59): Uses `inIfNotNull` — if `groupIds` is null, the filter has no `_id` constraint (matches all). If non-null, filters by those IDs.
**Go** (line 140-148): When `groupIDs` is empty (`len == 0`), the filter is `bson.M{}` (matches all). When non-empty, adds `$in` filter. **However**, there's a semantic difference: Java's `null` set means "match all" while Java's empty set means "match nothing" (since `in` with empty set matches nothing). In Go, an empty slice means "match all" which differs from Java's empty set behavior. But since the Java API uses `@Nullable Set<Long>` and `inIfNotNull`, the intent is null=match-all, which the Go empty-slice convention models correctly. **Acceptable parity.**
**Java** (line 61-67): Filters by `_id == groupId`, sets one field to `new Date()`.
**Go** (line 48-55): `UpdateVersion` filters by `_id == groupID`, sets one field to `time.Now()`, uses upsert. The Java version does NOT use upsert. **Bug: Go adds `SetUpsert(true)` which the Java version does not have.** This could create new documents when they shouldn't be created.
**Java** (line 69-84): Filters by `_id == groupId`, uses `setIfTrue` to conditionally set each of the 4 fields (MEMBERS, BLOCKLIST, JOIN_REQUESTS, JOIN_QUESTIONS) to the same `Date` instance.
**Go**: This method is **completely missing**. There is no Go method that accepts boolean flags for selectively updating multiple version fields in a single operation. The only `UpdateVersion` method takes a single `field string`.
**Java** (line 86-93): Projects only `BLOCKLIST` field, returns `GroupVersion::getBlocklist`.
**Go** (line 151-153): Delegates to `findSpecificVersion(ctx, groupID, "bl")` which projects `bl` and returns `version.Blocklist`. **Correct parity.**
**Java** (line 95-102): Projects only `INVITATIONS` field, returns `GroupVersion::getInvitations`.
**Go**: **This method is entirely missing from `group_version_repository.go`.** The `findSpecificVersion` helper only handles `"bl"`, `"jr"`, `"jq"`, `"mbr"` — there is no case for `"invt"` (invitations), and there is no `FindInvitations` method exposed. The `group_invitation_repository.go` file handles invitation CRUD but NOT the version timestamp lookup — that's a different domain concept (the version record for invitations).
**Java** (line 104-111): Projects `JOIN_REQUESTS`, returns `getJoinRequests`.
**Go** (line 156-158): Delegates to `findSpecificVersion(ctx, groupID, "jr")`. **Correct parity.**
**Java** (line 113-120): Projects `JOIN_QUESTIONS`, returns `getJoinQuestions`.
**Go** (line 161-163): Delegates to `findSpecificVersion(ctx, groupID, "jq")`. **Correct parity.**
**Java** (line 122-129): Projects `MEMBERS`, returns `getMembers`.
**Go** (line 166-168): Delegates to `findSpecificVersion(ctx, groupID, "mbr")`. **Correct parity.**
---

---

- [x] Non-transactional member removal and block insert: When the target user is a group member, Java wraps `deleteGroupMember` + `insert(blockedUser)` in a single MongoDB transaction (with retry). Go executes them sequentially with no transaction, risking data inconsistency if one operation fails.

---

- [x] Returns wrong data type — blocklist entries instead of user profiles: Java queries user profiles via `userService.queryUsersProfile(ids, false)` and builds a `UserInfosWithVersion` proto with `userProfile2proto(user)`. Go returns `[]po.GroupBlockedUser` (raw blocklist records), not user profile information. These are fundamentally different data.

---

- [x] Missing `NO_CONTENT` error for empty user profiles: Java throws `NO_CONTENT` when the queried user profiles are empty. Go has no such check.

---

- [x] Non-transactional member removal and block insert: When the target user is a group member, Java wraps `deleteGroupMember` + `insert(blockedUser)` in a single MongoDB transaction (with retry). Go executes them sequentially with no transaction, risking data inconsistency if one operation fails.

---

- [x] Returns wrong data type — blocklist entries instead of user profiles: Java queries user profiles via `userService.queryUsersProfile(ids, false)` and builds a `UserInfosWithVersion` proto with `userProfile2proto(user)`. Go returns `[]po.GroupBlockedUser` (raw blocklist records), not user profile information. These are fundamentally different data.

---

- [x] Missing `NO_CONTENT` error for empty user profiles: Java throws `NO_CONTENT` when the queried user profiles are empty. Go has no such check.

---

- [x] Missing `pastOrPresent` validation for `blockDate`: Java validates `blockDate` is not in the future. Go has no such validation.

---

- [x] **Missing `ExpirationDate` field.** The Java constructor passes `null` as the 9th argument (expiration date): `new GroupInvitation(id, groupId, inviterId, inviteeId, content, status, creationDate, responseDate, null)`. The Go `GroupInvitation` struct does not have an `ExpirationDate` field at all, so this field is never persisted to MongoDB.

---

- [x] **The Go code does not have a standalone `createGroupInvitation` method.** The Java code has a separate `createGroupInvitation` method (with nullable id, status, creationDate, responseDate) that is called from `authAndCreateGroupInvitation` and can be called independently (e.g., by admins). The Go code only has `AuthAndCreateGroupInvitation` which combines auth and creation, and a `CreateInvitation` alias that just delegates to `AuthAndCreateGroupInvitation`. There is no admin-level `createGroupInvitation` that skips permission checks.

---

- [x] **Missing early-return when no fields to update (in the flow).** The Java `createGroupInvitation` does input validation (maxContentLength, validRequestStatus, pastOrPresent dates) that the Go version lacks entirely.

---

- [x] **Missing config-based gate check.** The Java code checks `if (!allowRecallPendingInvitationByOwnerAndManager && !allowRecallBySender)` and returns `RECALLING_GROUP_INVITATION_IS_DISABLED` immediately. The Go code does not check any configuration flags and always allows recall attempts.

---

- [x] **Missing expiration check on the invitation.** The Java code checks `groupInvitationRepository.isExpired(invitation.getCreationDate().getTime())` after confirming the status is PENDING, and returns `RECALL_NON_PENDING_GROUP_INVITATION` with message "The invitation is under the status EXPIRED" if it is expired. The Go code does not check expiration at all.

---

- [x] **Missing dual-path query logic.** The Java code queries different fields depending on whether `allowRecallBySender` is true. If sender recall is allowed, it queries `groupId + inviterId + inviteeId + status` (to check if requester is the sender). If not, it queries only `groupId + inviteeId + status` (only owner/manager can recall). The Go code always queries `FindGroupIdAndInviterIdAndInviteeIdAndStatus`, ignoring this configuration-based branching.

---

- [x] **Missing version update for user sent/received invitations.** The Go code updates `UpdateSentGroupInvitationsVersion(inviterID)` and `UpdateReceivedGroupInvitationsVersion(inviteeID)`, which is correct. However, the Java code only updates the group invitations version (`groupVersionService.updateGroupInvitationsVersion`). It does NOT update user sent/received versions on recall. The Go code incorrectly updates user versions that Java does not.

---

- [x] **Missing expiration check.** The Java code checks `groupInvitationRepository.isExpired(invitation.getCreationDate().getTime())` for PENDING invitations and returns an error if expired. The Go code does not perform any expiration check.

---

- [x] **Missing transaction for ACCEPT action.** The Java code wraps the ACCEPT action (update status + add group member) in a transaction with retry (`inTransaction(...).retryWhen(TRANSACTION_RETRY)`). The Go code does not use a transaction — it calls `UpdateStatusIfPending` then `AddGroupMember` as separate operations, which is not atomic.

---

- [x] **Missing DuplicateKeyException handling for ACCEPT.** The Java code handles `DuplicateKeyException` when adding a member during ACCEPT (in case the user was already added by another concurrent request), returning `HandleHandleGroupInvitationResult(invitation, false)` instead of failing. The Go code does not handle this case.

---

- [x] **Redundant re-fetch of invitation for version update.** The Go code calls `s.invRepo.FindByID(ctx, invitationID)` after handling to get the inviter ID for version update. The Java code already has the invitation object from the initial query and does not re-fetch. While this is a performance concern rather than a logic bug, the extra fetch can also fail silently.

---

- [x] **Missing IGNORE and DECLINE handling distinction.** The Go code accepts any status that is `Accepted`, `Declined`, or `Ignored`, but then calls `UpdateStatusIfPending` with whatever status was passed. The Java code explicitly maps `ResponseAction.ACCEPT` → transaction with add member, `ResponseAction.IGNORE` → simple status update, `ResponseAction.DECLINE` → simple status update. The Go code doesn't distinguish between these in terms of transactional behavior.

---

- [x] **Missing version update for user sent invitations on accept.** Actually, the Java code for `authAndHandleInvitation` does NOT update user sent/received versions — it only updates the group invitations version via `updatePendingInvitationStatus`. The Go code updates both group and user versions, which is inconsistent with Java behavior.

---

- [x] **Missing NO_CONTENT check for empty results.** The Java code throws `ResponseStatusCode.NO_CONTENT` if the invitation list is empty. The Go code returns an empty list with the version, not an error.

---

- [x] **Missing expireAfter status transformation.** The Java code calls `ProtoModelConvertor.groupInvitation2proto(groupInvitation, expireAfterSeconds)` which transforms the status of expired invitations to `EXPIRED` when returning to clients. The Go code returns raw invitation objects without any status transformation.

---

- [x] **Hardcoded page size.** The Go code uses `0, 1000` as hardcoded page/size parameters. The Java code collects all results from the Flux without pagination limits.

---

- [x] **Missing `switchIfEmpty` on version Mono.** The Java code has `.switchIfEmpty(ResponseExceptionPublisherPool.alreadyUpToUpdate())` after the version `flatMap`, which means if the version is null/empty, it returns "already up to date". The Go code does not handle the case where version is nil — if version is nil and lastUpdatedDate is nil, it proceeds to query invitations instead of returning "already up to date".

---

- [x] **Missing NO_CONTENT check for empty results.** Same as above — Java throws `NO_CONTENT` for empty invitation lists, Go returns an empty list.

---

- [x] **Missing expireAfter status transformation.** Same as above — Go doesn't transform expired invitation statuses.

---

- [x] **Hardcoded page size.** Same as above — `0, 1000` hardcoded.

---

- [x] **Missing `switchIfEmpty` on version Mono.** Same as above — missing handling when version is nil.

---

- [x] **Missing filter parameters.** The Java method accepts `ids`, `groupIds` (Set), `inviterIds` (Set), `inviteeIds` (Set), `statuses` (Set), `creationDateRange`, `responseDateRange`, `expirationDateRange`, `page`, `size`. The Go `QueryInvitations` method only accepts a single `groupID`, `inviterID`, `inviteeID`, `status`, `lastUpdatedDate`, `page`, `size` — it does not accept `ids` (Set), multiple group IDs, multiple inviter IDs, multiple invitee IDs, multiple statuses, or date ranges for creation/response/expiration.

---

- [x] **`QueryInvitationsWithFilter` ignores all filter parameters.** The Go `QueryInvitationsWithFilter` method accepts full filter parameters (ids, groupIds, statuses, date ranges, etc.) but then calls `s.invRepo.FindInvitations(ctx, nil, nil, nil, nil, nil, p, sz)` passing `nil` for all filters. This means all filter criteria are silently discarded.

---

- [x] **Missing filter parameters.** Same as `queryInvitations` — the Java method accepts `ids`, `groupIds` (Set), `inviterIds` (Set), `inviteeIds` (Set), `statuses` (Set), and three DateRange parameters. The Go method only accepts single `groupID`, `inviterID`, `inviteeID`, `status`, and `lastUpdatedDate`. The `ids` parameter and date range parameters are missing.

---

- [x] **Missing validation that status is not PENDING.** The Java code validates `Validator.notEquals(requestStatus, RequestStatus.PENDING, "The request status must not be PENDING")`. The Go code does not validate that the new status is not PENDING.

---

- [x] **Missing expiration filter in the repository query.** The Java repository's `updateStatusIfPending` includes `.isNotExpired(GroupInvitation.Fields.CREATION_DATE, getEntityExpirationDate())` in the filter, meaning it won't update expired invitations. The Go repository's `UpdateStatusIfPending` only filters by `_id` and `stat: PENDING`, without checking expiration.

---

- [x] **Missing group version update on success.** The Java `updatePendingInvitationStatus` updates the group invitations version via `groupVersionService.updateGroupInvitationsVersion(groupId)` when the update is successful. The Go `UpdatePendingInvitationStatus` does not update any versions — it only returns `(bool, error)`.

---

- [x] **Missing early return when only responseDate is provided.** The Java code checks `Validator.areAllNull(inviterId, inviteeId, content, status, creationDate)` and returns an acknowledged result if all are null (note: `responseDate` is NOT in this check). The Go code checks `len(set) == 0` in the repository, which includes responseDate in the null check. This means if only `responseDate` is non-nil, Java would return early (no-op) but Go would proceed with the update. Actually re-reading: Java checks `areAllNull(inviterId, inviteeId, content, status, creationDate)` — `responseDate` is excluded from the early-return check. Go includes `responseDate` in the set, so it would NOT return early when only `responseDate` is set. This is actually correct behavior in Go (more permissive). The bug is in Java excluding responseDate from the check — but Go's behavior differs from Java.
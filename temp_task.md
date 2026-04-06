# Local Progress Tracker for Batch 12

[Context: ## handleQueryJoinedGroupIdsRequest]
- [x] `lastUpdatedDate` is not passed to the service: Java calls `queryJoinedGroupIdsWithVersion(userId, lastUpdatedDate)`, Go calls `QueryUserJoinedGroupIds(ctx, s.UserID)` ignoring the parsed `lastUpdatedDate`

---

[Context: ## handleQueryJoinedGroupIdsRequest]
- [x] Wrong service used: Java delegates to `groupService.queryJoinedGroupIdsWithVersion`, Go uses `groupMemberService.QueryUserJoinedGroupIds`

---

[Context: ## handleQueryJoinedGroupIdsRequest]
- [x] Dead code: The `if lastUpdatedDate != nil && c.groupService != nil` block (lines 202-205) is a no-op comment with no actual version checking logic

---

[Context: ## handleQueryJoinedGroupsRequest]
- [x] Fundamentally different implementation approach: Java calls `groupService.queryJoinedGroupsWithVersion(userId, lastUpdatedDate)` as a single call that handles version checking internally, while Go manually queries group IDs first then queries groups separately — this misses the version comparison optimization

---

[Context: ## handleQueryJoinedGroupsRequest]
- [x] Missing version in response: Java returns `GroupsWithVersion` which includes a `lastUpdatedDate` version field; Go's response never populates `LastUpdatedDate` on `GroupsWithVersion`

---

[Context: ## handleUpdateGroupRequest]
- [x] Missing ownership transfer branch: Java has a critical branch where if `successorId != null`, it calls `authAndTransferGroupOwnership(userId, groupId, successorId, quitAfterTransfer, null)` instead of the regular update — Go only has the regular update path

---

[Context: ## handleUpdateGroupRequest]
- [x] Missing `muteEndDate` field: Java extracts `request.getMuteEndDate()` and passes it to `authAndUpdateGroupInformation`, Go has a TODO comment but does not pass it

---

[Context: ## handleUpdateGroupRequest]
- [x] Missing `userDefinedAttributes` field: Java passes `request.getUserDefinedAttributesMap()`, Go does not

---

[Context: ## handleUpdateGroupRequest]
- [x] Missing `announcement` field: Go passes `updateReq.Announcement` but does not convert nil/empty properly as Java does with the `hasAnnouncement()` check

---

[Context: ## handleUpdateGroupRequest]
- [x] Missing notification logic: Java conditionally notifies group members or requester's other sessions, Go has none

---

[Context: ## handleCreateGroupBlockedUserRequest]
- [x] Parameter order differs from Java: Java calls `authAndBlockUser(userId, groupId, userIdToBlock, null)`, Go calls `BlockUser(ctx, groupId, userId, s.UserID)` — the requester vs blocked-user ordering may be swapped depending on the Go service signature

---

[Context: ## handleCreateGroupBlockedUserRequest]
- [x] Missing notification logic: Java dispatches notifications to group members, blocked user, and requester's other sessions, Go has none

---

[Context: ## handleDeleteGroupBlockedUserRequest]
- [x] Missing `wasBlocked` check: Java checks `if (!wasBlocked) { return RequestHandlerResult.OK; }`, Go does not check this

---

[Context: ## handleDeleteGroupBlockedUserRequest]
- [x] Missing notification logic: Java dispatches notifications to group members, unblocked user, and requester's other sessions, Go has none

---

[Context: ## handleDeleteGroupBlockedUserRequest]
- [x] Missing requester ID parameter: Java calls `unblockUser(userId, groupId, userIdToUnblock, null, true)` passing the requester for auth, Go calls `UnblockUser(ctx, groupId, userId)` without the requester

---

[Context: ## handleQueryGroupBlockedUserIdsRequest]
- [x] Extra auth check not in Java: Java calls `queryGroupBlockedUserIdsWithVersion(groupId, lastUpdatedDate)` without a userId auth check, Go calls `AuthAndQueryGroupBlockedUserIds` with `s.UserID` — this may cause authorization failures if Java intentionally allows unauthenticated queries here

---

[Context: ## handleQueryGroupBlockedUsersInfosRequest]
- [x] Incomplete UserInfo construction: Java lets the service return full `UserInfosWithVersion` protos (which likely include user details), Go manually constructs `UserInfo` with only `Id` field set, losing other user info fields

---

[Context: ## handleQueryGroupBlockedUsersInfosRequest]
- [x] Extra auth check not in Java: Same as handleQueryGroupBlockedUserIdsRequest — Go adds auth that Java doesn't have

---

[Context: ## handleCheckGroupQuestionAnswerRequest]
- [x] Missing `joined`, `questionIds`, `score` fields in response: Java constructs `GroupJoinQuestionsAnswerResult` with `setJoined(joined)`, `addAllQuestionIds(questionIds)`, `setScore(answerResult.score())`, Go returns whatever `CheckGroupJoinQuestionsAnswersAndJoin` returns directly without mapping these fields

---

[Context: ## handleCheckGroupQuestionAnswerRequest]
- [x] Missing notification logic when `joined` is true: Java creates a `CreateGroupMembersRequest` notification and dispatches it to group members or the added member, Go has none

---

[Context: ## handleCreateGroupInvitationRequestRequest]
- [x] Missing auth prefix: Java calls `authAndCreateGroupInvitation`, Go calls `CreateInvitation` — missing the auth check

---

[Context: ## handleCreateGroupInvitationRequestRequest]
- [x] Missing notification logic: Java dispatches notifications to group members, owner/managers, invitee, and requester's other sessions, Go has none

---

[Context: ## handleCreateGroupJoinRequestRequest]
- [x] Missing auth prefix: Java calls `authAndCreateGroupJoinRequest`, Go calls `CreateJoinRequest` — missing the auth check

---

[Context: ## handleCreateGroupJoinRequestRequest]
- [x] Missing notification logic: Java dispatches notifications to group members, owner/managers, and requester's other sessions, Go has none

---

[Context: ## handleCreateGroupQuestionsRequest]
- [x] Non-batched creation: Java calls `authAndCreateGroupJoinQuestions(userId, groupId, questions)` in a single batch call, Go iterates and calls `CreateJoinQuestion` individually — this is not atomic and could leave partial state

---

[Context: ## handleCreateGroupQuestionsRequest]
- [x] Missing auth check: Java calls `authAnd...`, Go calls `CreateJoinQuestion` without auth

---

[Context: ## handleCreateGroupQuestionsRequest]
- [x] Wrong response format: Java returns `RequestHandlerResult.ofDataLongs(questionIds)` (a `LongsWithVersion` with just longs), Go wraps in `LongsWithVersion` which adds a version field not present in Java's response

---

[Context: ## handleDeleteGroupInvitationRequest]
- [x] Missing notification logic: Java dispatches notifications to group members, owner/managers, invitee, and requester's other sessions, Go has none

---

[Context: ## handleUpdateGroupInvitationRequest]
- [x] Missing `reason` parameter: Java passes `request.getReason()` to `authAndHandleInvitation`, Go does not pass it

---

[Context: ## handleUpdateGroupInvitationRequest]
- [x] Missing auth prefix: Java calls `authAndHandleInvitation`, Go calls `ReplyToInvitation` — missing auth

---

[Context: ## handleUpdateGroupInvitationRequest]
- [x] Missing complex multi-notification logic: Java has extremely complex logic that sends separate notifications for invitation updates AND member additions (when invitation is accepted and requester joins), including querying group member IDs, owner/manager IDs — Go has none of this

---

[Context: ## handleDeleteGroupJoinRequestRequest]
- [x] Missing notification logic: Java dispatches notifications to group members, owner/managers, and requester's other sessions, Go has none

---

[Context: ## handleUpdateGroupJoinRequestRequest]
- [x] Missing `reason` parameter: Java passes `request.getReason()`, Go does not

---

[Context: ## handleUpdateGroupJoinRequestRequest]
- [x] Missing auth prefix: Java calls `authAndHandleJoinRequest`, Go calls `ReplyToJoinRequest` — missing auth

---

[Context: ## handleUpdateGroupJoinRequestRequest]
- [x] Missing complex multi-notification logic: Java handles requester-added-as-new-member notifications, querying group members, and sending separate join-request-updated and member-added notifications, Go has none

---

[Context: ## handleDeleteGroupJoinQuestionsRequest]
- [x] Non-batched deletion: Java calls `authAndDeleteGroupJoinQuestions(userId, groupId, questionIdsSet)` as a single batch call, Go iterates and calls `DeleteJoinQuestion` individually — not atomic

---

[Context: ## handleDeleteGroupJoinQuestionsRequest]
- [x] Missing auth check: Java calls `authAnd...`, Go calls `DeleteJoinQuestion` without auth

---

[Context: ## handleDeleteGroupJoinQuestionsRequest]
- [x] Missing `groupId` parameter: Java passes `groupId` for authorization, Go does not pass it

---

[Context: ## handleQueryGroupJoinRequestsRequest]
- [x] Different branching logic: Java calls `authAndQueryGroupJoinRequestsWithVersion(userId, groupId, lastUpdatedDate)` with a single method regardless of whether `groupId` is null, Go branches into two different methods — this may produce different behavior when `groupId` is null

---

[Context: ## handleUpdateGroupJoinQuestionRequest]
- [x] Hardcoded `groupId=0`: Go passes `0` as the second argument to `UpdateJoinQuestion`, Java passes `request.getQuestionId()` to `authAndUpdateGroupJoinQuestion` which uses the question's own groupId internally for auth — the `0` is likely wrong

---

[Context: ## handleUpdateGroupJoinQuestionRequest]
- [x] Missing auth check: Java calls `authAnd...`, Go calls `UpdateJoinQuestion` without auth

---

[Context: ## handleCreateGroupMembersRequest]
- [x] Missing `name` parameter: Java passes `request.getName()` (or null), Go does not pass it

---

[Context: ## handleCreateGroupMembersRequest]
- [x] Missing notification logic: Java conditionally notifies other group members, added members, and requester's other sessions, Go has none

---

[Context: ## handleDeleteGroupMembersRequest]
- [x] Missing empty-deletion check: Java checks `if (deletedUserIds.isEmpty()) { return RequestHandlerResult.OK; }`, Go does not

---

[Context: ## handleDeleteGroupMembersRequest]
- [x] Missing notification logic: Java conditionally notifies other group members, removed members, and requester's other sessions, Go has none

---

[Context: ## handleQueryGroupMembersRequest]
- [x] Missing `memberIds`-based query branch: Java has two paths — if `memberIdsCount > 0`, calls `authAndQueryGroupMembers(userId, groupId, memberIds, withStatus)`; otherwise calls `authAndQueryGroupMembersWithVersion(userId, groupId, lastUpdatedDate, withStatus)`. Go only has the versioned path

---

[Context: ## handleQueryGroupMembersRequest]
- [x] Missing `withStatus` parameter: Java extracts `request.getWithStatus()` and passes it, Go does not extract or pass it

---

[Context: ## handleUpdateGroupMemberRequest]
- [x] Missing notification logic: Java conditionally notifies other group members, updated member, and requester's other sessions, Go has none
*Checked methods: NewGroupQuestion(String question, LinkedHashSet<String> answers, Integer score)*
Now I have the complete picture. Let me compare the Java and Go implementations.
**Java `NewGroupQuestion`**: This is a simple Java record (data carrier) with fields `String question`, `LinkedHashSet<String> answers`, `Integer score`. No validation logic in the constructor itself.
**Java `validNewGroupQuestion`**: Validates that:
1. `answers` is not empty (throws if empty)
2. `score` is not null AND `score >= 0` (throws if null or negative)
**Go `NewGroupQuestion` struct**: Has fields `Question *string`, `Answers []string`, `Score *int`. This is a faithful structural port (slice instead of LinkedHashSet is acceptable since Go has no LinkedHashSet).
**Go `ValidNewGroupQuestion`**: Only checks if the input is nil. Missing the core validation logic from Java.

---

[Context: ## ValidNewGroupQuestion]
- [x] Missing validation that `Answers` slice is not empty. The Java version throws `EMPTY_GROUP_QUESTION_ANSWERS` ("The answers of a new group question should not be empty") when `question.answers().isEmpty()`. The Go version only checks for nil input and does not validate the answers field at all.

---

[Context: ## ValidNewGroupQuestion]
- [x] Missing validation that `Score` is not nil and is >= 0. The Java version throws `ILLEGAL_GROUP_QUESTION_SCORE` ("The score of a new group question should not be null and must be greater than or equal to 0") when `score` is null or negative. The Go version does not validate the score field at all.

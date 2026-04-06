# Local Progress Tracker for Batch 10

[Context: ## queryGroupBlockedUsers (non-paginated @GetMapping)]
- [x] Missing all filter parameters: `groupIds`, `userIds`, `blockDateStart`, `blockDateEnd`, `requesterIds`. The Java endpoint accepts these as optional query parameters and passes them to the service call. The Go method signature only has `page` and `size`.

---

[Context: ## queryGroupBlockedUsers (non-paginated @GetMapping)]
- [x] Discards the queried results (uses `_`). Java returns the collection of `GroupBlockedUser` in the response.

---

[Context: ## queryGroupBlockedUsers (non-paginated @GetMapping)]
- [x] Java explicitly passes `0` as the page parameter for the non-paginated query (hardcoded first page), while Go accepts a `page` parameter, changing the semantics.

---

[Context: ## queryGroupBlockedUsers (non-paginated @GetMapping)]
- [x] Missing `getPageSize(size)` equivalent — Java normalizes the `size` parameter via `getPageSize()` before passing to the service.

---

[Context: ## queryGroupBlockedUsers (paginated @GetMapping("page"))]
- [x] Missing entirely — there is no paginated query endpoint in Go. Java has a separate `@GetMapping("page")` endpoint that calls both `countBlockedUsers()` and `queryBlockedUsers()`, returning a `PaginationDTO` with total count + records.

---

[Context: ## queryGroupBlockedUsers (paginated @GetMapping("page"))]
- [x] `QueryGroupBlockedUsersWithQuery` does not call `countBlockedUsers()` equivalent, so even with filters, there's no total count returned — Java returns `PaginationDTO<GroupBlockedUser>` which includes the count.

---

[Context: ## queryGroupBlockedUsers (paginated @GetMapping("page"))]
- [x] Discards queried results (uses `_`). Java returns records in the response.

---

[Context: ## queryGroupBlockedUsers (paginated @GetMapping("page"))]
- [x] Missing `getPageSize(size)` equivalent.

---

[Context: ## updateGroupBlockedUsers]
- [x] Does not deduplicate keys. Java converts `List<GroupBlockedUser.Key>` to `Set<GroupBlockedUser.Key>` via `CollectionUtil.newSet(keys)` before passing to the service. Go passes the raw slice.

---

[Context: ## updateGroupBlockedUsers]
- [x] Returns empty `RequestHandlerResult{}` instead of `UpdateResultDTO`. Java maps the result through `UpdateResultDTO::get` and returns it, providing info about acknowledged/modified count.

---

[Context: ## deleteGroupBlockedUsers]
- [x] Does not deduplicate keys. Java converts `List<GroupBlockedUser.Key>` to `Set<GroupBlockedUser.Key>` via `CollectionUtil.newSet(keys)` before passing to the service. Go passes the raw slice.

---

[Context: ## deleteGroupBlockedUsers]
- [x] Returns empty `RequestHandlerResult{}` instead of `DeleteResultDTO`. Java maps the result through `DeleteResultDTO::get` and returns it, providing info about deleted count.
*Checked methods: addGroup(@RequestBody AddGroupDTO addGroupDTO), queryGroups(@QueryParam(required = false), queryGroups(@QueryParam(required = false), countGroups(@QueryParam(required = false), updateGroups(Set<Long> ids, @RequestBody UpdateGroupDTO updateGroupDTO), deleteGroups(@QueryParam(required = false)*
Now I have a complete picture. Let me compile the bug report.

---

[Context: ## AddGroup]
- [x] **Missing nil check for `CreatorId` before dereferencing**: At `group_controllers.go:105`, `*addGroupDTO.CreatorId` is dereferenced unconditionally, but if both `CreatorId` and `OwnerId` are `nil`, this will cause a nil pointer panic. The Java code (line 82) calls `addGroupDTO.creatorId()` which returns `null` safely, and the `ownerId` fallback logic (line 83-85) also handles null `creatorId`. In Go, the nil `ownerId` at line 91-96 is only set if `CreatorId` is non-nil (line 94-95), so when both are nil, `ownerId` stays at zero-value `0`, which is wrong but doesn't panic. However, `*addGroupDTO.CreatorId` at line 105 **will panic** if `CreatorId` is nil.

---

[Context: ## AddGroup]
- [x] **`ownerId` defaults to zero when both `CreatorId` and `OwnerId` are nil**: In Java (lines 83-85), when `ownerId` is null, it falls back to `creatorId()`. The Go code (lines 91-96) does the same fallback, but if both are nil, `ownerId` stays at `int64(0)` — an invalid ID. The Java code would pass `null` to the service which handles it differently. This is a behavioral difference.

---

[Context: ## QueryGroups (non-paged, GET /)]
- [x] **Missing `lastUpdatedDateStart`/`lastUpdatedDateEnd` parameters**: The Java `queryGroups` method (lines 110-111) accepts `lastUpdatedDateStart` and `lastUpdatedDateEnd` and passes them as `DateRange.of(lastUpdatedDateStart, lastUpdatedDateEnd)` to `groupService.queryGroups`. The Go `QueryGroupsWithQuery` method signature at line 130 does **not** include these parameters at all — the function only has `creationDateStart, creationDateEnd, deletionDateStart, deletionDateEnd, muteEndDateStart, muteEndDateEnd` but no `lastUpdatedDate` pair.

---

[Context: ## QueryGroups (non-paged, GET /)]
- [x] **Missing `getPageSize` logic for the `size` parameter**: Java (line 116) calls `size = getPageSize(size)` which applies `defaultAvailableRecordsPerRequest` when size is null/<=0 and caps at `maxAvailableRecordsPerRequest`. The Go code passes `size` directly to the service without any defaulting or capping, and does not enforce a maximum page size. This is a behavioral difference that could allow unbounded queries.

---

[Context: ## QueryGroups (paged, GET /page)]
- [x] **Entire paged endpoint is missing**: The Java code has a second `queryGroups` method (lines 132-179) mapped to `GET "page"` that accepts a `page` parameter and returns a `PaginationDTO<Group>` with both a count query and a data query. The Go code has no equivalent — `QueryGroups` (line 122) only does simple pagination, and `QueryGroupsWithQuery` (line 130) passes page/size to the service but never performs the count query needed for pagination metadata.

---

[Context: ## QueryGroups (paged, GET /page)]
- [x] **Missing `lastUpdatedDateStart`/`lastUpdatedDateEnd` parameters** (same as the non-paged variant above).

---

[Context: ## CountGroups (GET /count)]
- [x] **Entire `countGroups` endpoint is missing from the Go controller**: The Java `countGroups` method (lines 181-231) is a complex statistics endpoint with `DivideBy` support, `GroupStatisticsDTO`, conditional counting of created/deleted groups and groups that sent messages, and `checkAndQueryBetweenDate` logic. None of this exists in the Go `GroupController`. There is no `CountGroups` or any statistics method on the controller.

---

[Context: ## UpdateGroups]
- [x] **Nil dereference on `QuitAfterTransfer`**: At line 140, `*updateGroupDTO.QuitAfterTransfer` is dereferenced unconditionally. In Java (line 256), `updateGroupDTO.quitAfterTransfer()` returns a nullable value. If `QuitAfterTransfer` is nil in the Go DTO but `SuccessorId` is non-nil, this will cause a nil pointer panic.

---

[Context: ## UpdateGroups]
- [x] **Missing second `nil` parameter for `GroupVersion` in `UpdateGroupsInformation` call**: The Java code (lines 240-253) calls `groupService.updateGroupsInformation(...)` with two trailing `null` parameters (for `session` and `groupVersion`). The Go code at line 151-163 passes only one trailing `nil` (for `session`). Looking at the Go `UpdateGroupsInformation` signature (lines 484-498), it takes `session mongo.SessionContext` as the last parameter. This is actually correct since Go's service method has one fewer parameter than Java's. **Not a bug** upon closer inspection — the Go service only takes `session`.

---

[Context: ## DeleteGroups]
- [x] **`deleteLogically` parameter is ignored**: The Java `deleteGroups` (lines 262-269) passes `deleteLogically` to `groupService.deleteGroupsAndGroupMembers(ids, deleteLogically)` where it determines whether to perform logical or physical deletion. The Go code (lines 171-179) ignores the `deleteLogical` parameter entirely, always passing `nil` to `DeleteGroupsAndGroupMembers`, which means the service always uses its default behavior (physical delete) instead of respecting the caller's explicit preference.

---

[Context: ## DeleteGroups]
- [x] **Return value does not include `DeleteResultDTO`**: The Java code (line 266-267) maps the `DeleteResult` to `DeleteResultDTO.get()` which contains the count of deleted records, and returns it via `HttpHandlerResult.okIfTruthy`. The Go code returns an empty `RequestHandlerResult{}` with no deleted count information, losing the deletion result data.
*Checked methods: addGroupInvitation(@RequestBody AddGroupInvitationDTO addGroupInvitationDTO), queryGroupInvitations(@QueryParam(required = false), queryGroupInvitations(@QueryParam(required = false), updateGroupInvitations(Set<Long> ids, @RequestBody UpdateGroupInvitationDTO updateGroupInvitationDTO), deleteGroupInvitations(@QueryParam(required = false)*
Now I have the full picture. Let me analyze each method carefully.
**Java** calls `groupInvitationService.createGroupInvitation()` with 8 parameters: `id`, `groupId`, `inviterId`, `inviteeId`, `content`, `status`, `creationDate`, `responseDate`.
**Go** calls `c.groupInvitationService.AuthAndCreateGroupInvitation()` with only 4 parameters: `inviterId`, `groupId`, `inviteeId`, `content`.
The Go version is missing: `id`, `status`, `creationDate`, `responseDate` parameters. The Java `createGroupInvitation` is different from `AuthAndCreateGroupInvitation` - the Java controller calls the admin-oriented `createGroupInvitation` which accepts all fields directly, while the Go code calls the client-facing `AuthAndCreateGroupInvitation` which does auth checks and only accepts 4 fields.
Also, the Java constructor takes `TurmsPropertiesManager` and passes it to the base controller via `super(propertiesManager)` — the Go version does not.
The Java version has two query methods: one without pagination (GET /) and one with pagination (GET /page). Both accept rich filter parameters and call `queryInvitations` with `DateRange` objects, and also call `countInvitations` for the page variant. Both map results to `GroupInvitationDTO` with `getEntityExpirationDate()`.
The Go version has `QueryGroupInvitations` (page, size) and `QueryGroupInvitationsWithQuery` (with all filters). The Go version doesn't call `countInvitations` for the pagination variant, and doesn't map results to include `getEntityExpirationDate()`.
**Java** calls `groupInvitationService.updateInvitations(ids, inviterId, inviteeId, content, status, creationDate, responseDate)`.
**Go** is a **no-op stub** — it returns an empty result without calling any service method.
**Java** calls `groupInvitationService.deleteInvitations(ids)` and maps the result to `DeleteResultDTO`.
**Go** calls `c.groupInvitationService.DeleteInvitations(ctx, ids)` — the service method exists and is called correctly. The result handling differs (Java wraps in `DeleteResultDTO`) but the core delete logic is present.

---

[Context: ## addGroupInvitation]
- [x] **Missing fields passed to service**: Java passes 8 fields (`id`, `groupId`, `inviterId`, `inviteeId`, `content`, `status`, `creationDate`, `responseDate`) to `createGroupInvitation`. Go only passes 4 fields (`inviterId`, `groupId`, `inviteeId`, `content`) to `AuthAndCreateGroupInvitation`. Missing: `id`, `status`, `creationDate`, `responseDate`.

---

[Context: ## addGroupInvitation]
- [x] **Wrong service method**: Java calls `createGroupInvitation` (admin-level, no auth checks). Go calls `AuthAndCreateGroupInvitation` (client-level, with auth checks). The admin controller should bypass auth and directly create the invitation with all provided fields.

---

[Context: ## queryGroupInvitations (non-paged)]
- [x] **Missing `GroupInvitationDTO` response mapping**: Java maps each result to `new GroupInvitationDTO(invitation, groupInvitationService.getEntityExpirationDate())`, attaching the entity expiration date. Go returns raw `RequestHandlerResult{}` without this mapping.

---

[Context: ## queryGroupInvitations (paged)]
- [x] **Missing count query**: Java calls `groupInvitationService.countInvitations(...)` separately to get total count for pagination. Go does not call a count method for the paged variant.

---

[Context: ## queryGroupInvitations (paged)]
- [x] **Missing `GroupInvitationDTO` response mapping**: Same as non-paged — Java maps results to `GroupInvitationDTO` with expiration date. Go does not.

---

[Context: ## updateGroupInvitations]
- [x] **Method is a no-op stub**: The method returns an empty `RequestHandlerResult{}` without calling any service method. Java calls `groupInvitationService.updateInvitations(ids, inviterId, inviteeId, content, status, creationDate, responseDate)`.
This method correctly calls `DeleteInvitations`. However:

---

[Context: ## deleteGroupInvitations]
- [x] **Missing `DeleteResultDTO` response**: Java maps the delete count to `DeleteResultDTO.get`. Go returns an empty `RequestHandlerResult{}` without including the deletion count in the response.
---
Here are the summarized bugs:

---

[Context: ## AddGroupInvitation]
- [x] Missing fields `id`, `status`, `creationDate`, `responseDate` from being passed to the service. Java's `createGroupInvitation` accepts all 8 fields; Go's `AuthAndCreateGroupInvitation` only receives 4.

---

[Context: ## AddGroupInvitation]
- [x] Wrong service method used: Java calls `createGroupInvitation` (admin API, no auth), Go calls `AuthAndCreateGroupInvitation` (client API, with auth checks). Admin controllers should not invoke auth-gated client methods.

---

[Context: ## QueryGroupInvitations (non-paged)]
- [x] Missing `GroupInvitationDTO` response mapping with `getEntityExpirationDate()`. Java wraps each invitation in `GroupInvitationDTO(invitation, expirationDate)`.

---

[Context: ## QueryGroupInvitations (paged)]
- [x] Missing `countInvitations` call for total count. Java calls `countInvitations(...)` separately before querying results for pagination.

---

[Context: ## QueryGroupInvitations (paged)]
- [x] Missing `GroupInvitationDTO` response mapping with `getEntityExpirationDate()`. Same issue as non-paged variant.

---

[Context: ## UpdateGroupInvitations]
- [x] Method is a no-op stub: returns empty result without calling `groupInvitationService.UpdateInvitations`. Java calls `updateInvitations(ids, inviterId, inviteeId, content, status, creationDate, responseDate)` — Go should call the existing `UpdateInvitations` service method with the DTO fields.

---

[Context: ## DeleteGroupInvitations]
- [x] Missing `DeleteResultDTO` response wrapping: Java maps the deleted count to `DeleteResultDTO.get()` to include the count in the response. Go returns an empty `RequestHandlerResult{}` without reporting how many records were deleted.
*Checked methods: addGroupJoinRequest(@RequestBody AddGroupJoinRequestDTO addGroupJoinRequestDTO), queryGroupJoinRequests(@QueryParam(required = false), queryGroupJoinRequests(@QueryParam(required = false), updateGroupJoinRequests(Set<Long> ids, @RequestBody UpdateGroupJoinRequestDTO updateGroupJoinRequestDTO), deleteGroupJoinRequests(@QueryParam(required = false)*
Now I have all the information needed for a thorough comparison. Let me summarize the bugs.

---

[Context: ## addGroupJoinRequest]
- [x] **Wrong service method called**: Java calls `createGroupJoinRequest()` (admin/bypass version that accepts id, groupId, requesterId, responderId, content, status, creationDate, responseDate, responseReason and directly inserts without auth checks). Go calls `AuthAndCreateJoinRequest()` (client-facing version that performs group validation, membership checks, blocklist checks, etc.). The admin endpoint should bypass these client-side validations.

---

[Context: ## addGroupJoinRequest]
- [x] **Missing fields passed to creation**: Java passes 9 fields (id, groupId, requesterId, responderId, content, status, creationDate, responseDate, responseReason). Go only passes 3 fields (requesterId, groupId, content), ignoring id, responderId, status, creationDate, responseDate, and responseReason entirely.

---

[Context: ## addGroupJoinRequest]
- [x] **Missing expiration date in response**: Java wraps the result in `new GroupJoinRequestDTO(request, groupJoinRequestService.getEntityExpirationDate())` which includes an `expirationDate` field. The Go controller returns an empty `RequestHandlerResult{}` without any created entity data or expiration date.

---

[Context: ## queryGroupJoinRequests (non-paginated)]
- [x] **Missing response data**: Java maps each result through `new GroupJoinRequestDTO(request, groupJoinRequestService.getEntityExpirationDate())` to include `expirationDate` in the response. Go discards the query results (assigns to `_`) and returns an empty `RequestHandlerResult{}`, so the client receives no data.

---

[Context: ## queryGroupJoinRequests (paginated - "page" endpoint)]
- [x] **Entire paginated endpoint missing**: Java has a separate `@GetMapping("page")` method with pagination support (accepts `page` parameter, calls `countJoinRequests` for total count, then `queryJoinRequests` with page/size, wrapping via `HttpHandlerResult.page()`). Go has no corresponding paginated query method — the `QueryGroupJoinRequestsWithQuery` method doesn't differentiate between paginated and non-paginated queries.

---

[Context: ## updateGroupJoinRequests]
- [x] **Method body is a no-op**: Java calls `groupJoinRequestService.updateJoinRequests(ids, requesterId, responderId, content, status, creationDate, responseDate)`. Go returns an empty `RequestHandlerResult{}` immediately without calling `UpdateJoinRequests` or any service method at all. All update fields (requesterId, responderId, content, status, creationDate, responseDate) are silently ignored.

---

[Context: ## deleteGroupJoinRequests]
- [x] **Missing response data**: Java maps the result through `DeleteResultDTO::get` to return the count of deleted records. Go discards the returned count (assigns to `_`) and returns an empty `RequestHandlerResult{}`, so the client gets no information about how many records were deleted.
*Checked methods: addGroupMember(@RequestBody AddGroupMemberDTO addGroupMemberDTO), queryGroupMembers(@QueryParam(required = false), queryGroupMembers(@QueryParam(required = false), updateGroupMembers(List<GroupMember.Key> keys, @RequestBody UpdateGroupMemberDTO updateGroupMemberDTO), deleteGroupMembers(@QueryParam(required = false)*
Good — `AddGroupMembers` does accept and pass through `name`. The controller passes `addGroupMemberDTO.Name`, so that's fine. Let me also verify the Java `addGroupMember` service uses `insert` not upsert.
The Java service uses `groupMemberRepository.insert(groupMember, session)` — this is a true insert that fails on duplicate keys.
Now let me compile the final bug list:

---

[Context: ## addGroupMember]
- [x] **Repository uses upsert instead of insert**: The Go repository `AddGroupMember` (`group_member_repository.go:34-41`) uses `UpdateOne` with `SetUpsert(true)`, which silently overwrites an existing group member. The Java version uses `groupMemberRepository.insert()`, which fails if a member already exists, preserving data integrity.

---

[Context: ## addGroupMember]
- [x] **Null role silently defaults instead of erroring**: When `role` is nil in the DTO, the Go controller defaults to `MEMBER` (line 348). In Java, a null role would trigger `Validator.notNull(groupMemberRole, "groupMemberRole")` in the service, returning an error. The Go version silently assigns a default role instead of rejecting the request.

---

[Context: ## queryGroupMembers (non-paged)]
- [x] **Missing non-paged query with filters**: Java has a `@GetMapping` (non-paged) variant that accepts filter parameters (`groupIds`, `userIds`, `roles`, date ranges) and passes `page=0`. The Go `QueryGroupMembersWithQuery` accepts a `page` parameter but there is no corresponding non-paged controller method that passes `page=0` with filter parameters like the Java version does.

---

[Context: ## queryGroupMembers (paged)]
- [x] **Missing count query for paginated endpoint**: The Java paged variant (`@GetMapping("page")`) calls both `countMembers(...)` and `queryGroupMembers(...)`, returning a `PaginationDTO` with total count and data. The Go `QueryGroupMembersWithQuery` only queries members — it never calls a count method, so no total count is returned for pagination.

---

[Context: ## updateGroupMembers]
- [x] **Iterates one-by-one instead of batch update**: The Go controller iterates over each key and calls `UpdateGroupMember` individually (lines 318-332). The Java version converts the list to a set and calls `updateGroupMembers` as a single batch operation via `groupMemberRepository.updateGroupMembers(keys, ...)`. This is functionally different — individual updates are not atomic and if one fails mid-way, some members are updated and others are not.

---

[Context: ## updateGroupMembers]
- [x] **updateGroupMembersVersion parameter is false instead of true**: The Go controller passes `false` for `updateVersion` (line 327), meaning the group members version is never updated. The Java controller passes `true`, which triggers `groupVersionService.updateMembersVersion(groupId)` to notify clients of the change via the version mechanism.
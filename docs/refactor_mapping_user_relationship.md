# User Relationship and Friend Request Subsystem Java to Go Refactor Mapping

## 1. PO / Domain Entities
### `im.turms.service.domain.user.po.UserRelationship`
- **Type**: Code File (Entity / PO)
- **Functions**: 
  - [x] `UserRelationship`: Represents the data model of a user relationship (ownerId, relatedUserId, name, blockDate, establishmentDate). Maps to `internal/domain/user/po/user_relationship.go`.

### `im.turms.service.domain.user.po.UserFriendRequest`
- **Type**: Code File (Entity / PO)
- **Functions**: 
  - [ ] `UserFriendRequest`: Represents the data model of a friend request (id, content, status, reason, creationDate, responseDate, requesterId, recipientId). Maps to `internal/domain/user/po/user_friend_request.go`.

## 2. Repositories (Data Access Layer)
### `im.turms.service.domain.user.repository.UserRelationshipRepository`
- **Type**: Code File (Repository)
- **Functions**:
  - [ ] `deleteAllRelationships`: Delete all relationships of specific owners.
  - [ ] `deleteByIds`: Delete relationships by composite keys.
  - [ ] `deleteById`: Delete relationship by ownerId and relatedUserId.
  - [ ] `findRelatedUserIds`: Find related user IDs (optionally blocked).
  - [ ] `findRelationships`: Query relationship records with pagination and filters.
  - [ ] `countRelationships`: Count relationships.
  - [ ] `upsert`: Update or Insert a one-sided relationship.
  - [x] `insert`: Insert a new relationship. (Implemented as `Insert` in `internal/domain/user/repository/user_relationship_repository.go`)
  - [ ] `isBlocked`: Check if relatedUser is blocked by owner.
  - [x] `hasRelationshipAndNotBlocked`: Check if relationship exists and not blocked. (Implemented as `HasRelationshipAndNotBlocked` in `internal/domain/user/repository/user_relationship_repository.go`)
  - [ ] `updateUserOneSidedRelationships`: Update relationship attributes (alias, blockDate, etc.). (Partially implemented as `UpdateBlockDate`)

### `im.turms.service.domain.user.repository.UserFriendRequestRepository`
- **Type**: Code File (Repository)
- **Functions**:
  - [ ] `deleteExpiredData`: Cleanup expired friend requests.
  - [ ] `hasPendingFriendRequest`: Check if a pending request exists.
  - [ ] `hasPendingOrDeclinedOrIgnoredOrExpiredRequest`: Check historical requests to limit spam.
  - [ ] `insert`: Create a new friend request.
  - [ ] `updateStatusIfPending`: Update request status (accept, decline, etc.) atomically.
  - [ ] `updateFriendRequests`: Bulk update requests.
  - [ ] `findRecipientId`: Projection query to get recipientId.
  - [ ] `findRequesterIdAndRecipientIdAndStatus`: Projection query.
  - [ ] `findRequesterIdAndRecipientIdAndCreationDateAndStatus`: Projection query for auth and validation.
  - [ ] `findFriendRequestsByRecipientId`: Query requests for recipient.
  - [ ] `findFriendRequestsByRequesterId`: Query requests for requester.
  - [ ] `deleteByIds`: Explicitly delete requests.
  - [ ] `findFriendRequests`, `countFriendRequests`: Universal query and count.

## 3. Services (Business Logic Layer)
### `im.turms.service.domain.user.service.UserRelationshipService`
- **Type**: Code File (Service)
- **Functions**:
  - [ ] `deleteAllRelationships`: Cascading delete of user relationships.
  - [ ] `deleteOneSidedRelationships`, `deleteOneSidedRelationship`: Remove unidirectional relationships.
  - [ ] `tryDeleteTwoSidedRelationships`: Safe removal of bi-directional relationship.
  - [ ] `queryRelatedUserIdsWithVersion`, `queryRelationshipsWithVersion`: Returns data with versioning sync.
  - [ ] `queryRelatedUserIds`, `queryRelationships`, `queryMembersRelationships`: Flexible queries.
  - [ ] `countRelationships`: Aggregate counting.
  - [ ] `friendTwoUsers`: Transactional creation of bidirectional friendships.
  - [ ] `upsertOneSidedRelationship`: Create/Update one-sided relationship.
  - [ ] `isBlocked`, `isNotBlocked`: Cache-backed blocked checks.
  - [x] `hasRelationshipAndNotBlocked`: Cache-backed relationship checks. (Implemented as `HasRelationshipAndNotBlocked` in `internal/domain/user/service/user_relationship_service.go`)
  - [ ] `hasOneSidedRelationship`: Basic logic check.
  - [ ] `updateUserOneSidedRelationships`: Modifying relationship states.

### `im.turms.service.domain.user.service.UserFriendRequestService`
- **Type**: Code File (Service)
- **Functions**:
  - [ ] `updateProperties`: Live config reload handling.
  - [ ] `removeAllExpiredFriendRequests`: TTL and cron-triggered cleanup.
  - [ ] `hasPendingFriendRequest`, `hasPendingOrDeclinedOrIgnoredOrExpiredRequest`: Spam protection checks.
  - [ ] `createFriendRequest`, `authAndCreateFriendRequest`: Verified friend request initiation.
  - [ ] `authAndRecallFriendRequest`: Sender cancels an unhandled request.
  - [ ] `updatePendingFriendRequestStatus`, `updateFriendRequests`: Admin/System updates.
  - [ ] `queryRecipientId`, `queryRequesterIdAndRecipientIdAndStatus`, `queryRequesterId...AndStatus`: Inner checks.
  - [ ] `authAndHandleFriendRequest`: Recipient reacting (accept/decline/ignore).
  - [ ] `queryFriendRequestsWithVersion`, `queryFriendRequestsByRecipientId`, `queryFriendRequestsByRequesterId`, `queryFriendRequests`: API reads.
  - [ ] `deleteFriendRequests`, `countFriendRequests`: Deletion and stats.

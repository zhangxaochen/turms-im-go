package service

import (
	"context"
	"errors"
	"math/rand"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"

	"im.turms/server/internal/domain/admin/permission"
	"im.turms/server/internal/domain/admin/po"
	"im.turms/server/internal/domain/admin/repository"
	"im.turms/server/internal/domain/common/infra/idgen"
)

var (
	ErrRequesterNotExist = errors.New("UNAUTHORIZED: requester does not exist")
	ErrPermissionDenied  = errors.New("permission denied")
)

const (
	RootRoleID int64 = 0
	RootAdminID int64 = 0

	// rootRoleRank is the rank of the root role (Integer.MAX_VALUE in Java).
	rootRoleRank = int(^uint(0) >> 1)

	// MinRoleNameLimit and MaxRoleNameLimit match Java's MIN_ROLE_NAME_LIMIT and MAX_ROLE_NAME_LIMIT.
	MinRoleNameLimit = 1
	MaxRoleNameLimit = 32
)

// rootRoleRankValue is a variable copy of rootRoleRank for taking its address.
var rootRoleRankValue = rootRoleRank

// rootRole is the in-memory root admin role (not stored in DB).
var rootRole = &po.AdminRole{
	ID:          RootRoleID,
	Name:        "ROOT",
	Permissions: permission.AllAdminPermissions,
	Rank:        rootRoleRank,
}

// getRootRole returns a copy of the in-memory root admin role.
func getRootRole() *po.AdminRole {
	return &po.AdminRole{
		ID:          RootRoleID,
		Name:        "ROOT",
		Permissions: permission.AllAdminPermissions,
		Rank:        rootRoleRank,
	}
}

// randomAlphabetic generates a random alphabetic string of the given length.
func randomAlphabetic(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// AdminRoleService maps to AdminRoleService in Java.
// @MappedFrom AdminRoleService
type AdminRoleService interface {
	AuthAndAddAdminRole(ctx context.Context, requesterId int64, roleId *int64, name string, permissions []permission.AdminPermission, rank *int) (*po.AdminRole, error)
	AddAdminRole(ctx context.Context, roleId int64, name string, permissions []permission.AdminPermission, rank *int) (*po.AdminRole, error)
	AuthAndDeleteAdminRoles(ctx context.Context, requesterId int64, roleIds []int64) (int64, error)
	DeleteAdminRoles(ctx context.Context, roleIds []int64) (int64, error)
	AuthAndUpdateAdminRoles(ctx context.Context, requesterId int64, roleIds []int64, newName *string, permissions []permission.AdminPermission, rank *int) (int64, error)
	UpdateAdminRole(ctx context.Context, roleIds []int64, newName *string, permissions []permission.AdminPermission, rank *int) (int64, error)
	QueryAdminRoles(ctx context.Context, ids []int64, names []string, includedPermissions []permission.AdminPermission, ranks []int, page *int, size *int) ([]*po.AdminRole, error)
	QueryAndCacheRolesByRoleIdsAndRankGreaterThan(ctx context.Context, roleIds []int64, rankGreaterThan int) ([]*po.AdminRole, error)
	CountAdminRoles(ctx context.Context, ids []int64, names []string, includedPermissions []permission.AdminPermission, ranks []int) (int64, error)
	QueryHighestRankByAdminId(ctx context.Context, adminId int64) (*int, error)
	QueryHighestRankByRoleIds(ctx context.Context, roleIds []int64) (*int, error)
	IsAdminRankHigherThanRank(ctx context.Context, adminId int64, rank int) (bool, error)
	QueryPermissions(ctx context.Context, adminId int64) ([]permission.AdminPermission, error)
	QueryRoleIdsByAdminId(ctx context.Context, adminId int64) ([]int64, error)
	SetAdminService(adminService AdminService)
}

type adminRoleService struct {
	idGen        *idgen.SnowflakeIdGenerator
	repo         repository.AdminRoleRepository
	adminService AdminService

	idToRole map[int64]*po.AdminRole
	mutex    sync.RWMutex
}

func NewAdminRoleService(idGen *idgen.SnowflakeIdGenerator, repo repository.AdminRoleRepository) AdminRoleService {
	return &adminRoleService{
		idGen:    idGen,
		repo:     repo,
		idToRole: make(map[int64]*po.AdminRole),
	}
}

func (s *adminRoleService) SetAdminService(adminService AdminService) {
	s.adminService = adminService
}

// @MappedFrom authAndAddAdminRole
func (s *adminRoleService) AuthAndAddAdminRole(ctx context.Context, requesterId int64, roleId *int64, name string, permissions []permission.AdminPermission, rank *int) (*po.AdminRole, error) {
	if roleId != nil && *roleId == RootRoleID {
		return nil, errors.New("the root role cannot be created")
	}
	if name == "" {
		return nil, errors.New("name must not be blank")
	}
	if strings.Contains(name, " ") || strings.Contains(name, "\t") || strings.Contains(name, "\n") || strings.Contains(name, "\r") {
		return nil, errors.New("name must not contain whitespace")
	}
	if len(name) < MinRoleNameLimit || len(name) > MaxRoleNameLimit {
		return nil, errors.New("name length must be between 1 and 32")
	}
	if len(permissions) == 0 {
		return nil, errors.New("permissions must not be empty")
	}
	if rank == nil {
		return nil, errors.New("rank must not be null")
	}
	higher, err := s.IsAdminRankHigherThanRank(ctx, requesterId, *rank)
	if err != nil {
		return nil, err
	}
	if !higher {
		// Check if requester exists at all (Java: switchIfEmpty(errorRequesterNotExist))
		requesterPerms, permErr := s.QueryPermissions(ctx, requesterId)
		if permErr != nil {
			return nil, permErr
		}
		if len(requesterPerms) == 0 {
			return nil, ErrRequesterNotExist
		}
		return nil, ErrPermissionDenied
	}
	// Verify that the requester has all the requested permissions
	requesterPermissions, err := s.QueryPermissions(ctx, requesterId)
	if err != nil {
		return nil, err
	}
	if len(requesterPermissions) == 0 {
		return nil, ErrRequesterNotExist
	}
	permMap := make(map[permission.AdminPermission]bool)
	for _, p := range requesterPermissions {
		permMap[p] = true
	}
	for _, p := range permissions {
		if !permMap[p] {
			return nil, ErrPermissionDenied
		}
	}

	id := int64(0)
	if roleId != nil {
		id = *roleId
	} else {
		id = s.idGen.NextIncreasingId()
	}
	return s.AddAdminRole(ctx, id, name, permissions, rank)
}

// @MappedFrom addAdminRole
func (s *adminRoleService) AddAdminRole(ctx context.Context, roleId int64, name string, permissions []permission.AdminPermission, rank *int) (*po.AdminRole, error) {
	if roleId == RootRoleID {
		return nil, errors.New("the new role ID cannot be the root role ID")
	}
	if len(permissions) == 0 {
		return nil, errors.New("permissions must not be empty")
	}
	if strings.Contains(name, " ") || strings.Contains(name, "\t") || strings.Contains(name, "\n") || strings.Contains(name, "\r") {
		return nil, errors.New("name must not contain whitespace")
	}
	if len(name) < MinRoleNameLimit || len(name) > MaxRoleNameLimit {
		return nil, errors.New("name length must be between 1 and 32")
	}
	if rank == nil {
		return nil, errors.New("rank must not be null")
	}
	role := &po.AdminRole{
		ID:           roleId,
		Name:         name,
		Permissions:  permissions,
		Rank:         *rank,
		CreationDate: time.Now(),
	}
	err := s.repo.Insert(ctx, role)
	if err != nil {
		return nil, err
	}
	s.mutex.Lock()
	s.idToRole[roleId] = role
	s.mutex.Unlock()
	return role, nil
}

// @MappedFrom authAndDeleteAdminRoles
func (s *adminRoleService) AuthAndDeleteAdminRoles(ctx context.Context, requesterId int64, roleIds []int64) (int64, error) {
	if len(roleIds) == 0 {
		return 0, nil
	}
	for _, id := range roleIds {
		if id == RootRoleID {
			return 0, errors.New("the root admin is reserved and cannot be deleted")
		}
	}
	// Query each target role individually and check rank (Java: checkIfAllowedToManageRoles)
	targetHighest, err := s.QueryHighestRankByRoleIds(ctx, roleIds)
	if err != nil {
		return 0, err
	}
	if targetHighest == nil {
		return 0, nil
	}
	higher, err := s.IsAdminRankHigherThanRank(ctx, requesterId, *targetHighest)
	if err != nil {
		return 0, err
	}
	if !higher {
		// Check if requester exists
		requesterRank, rankErr := s.QueryHighestRankByAdminId(ctx, requesterId)
		if rankErr != nil {
			return 0, rankErr
		}
		if requesterRank == nil {
			return 0, ErrRequesterNotExist
		}
		return 0, ErrPermissionDenied
	}
	return s.DeleteAdminRoles(ctx, roleIds)
}

// @MappedFrom deleteAdminRoles
func (s *adminRoleService) DeleteAdminRoles(ctx context.Context, roleIds []int64) (int64, error) {
	if len(roleIds) == 0 {
		return 0, nil
	}
	for _, id := range roleIds {
		if id == RootRoleID {
			return 0, errors.New("the root admin is reserved and cannot be deleted")
		}
	}
	deleted, err := s.repo.DeleteAdminRoles(ctx, roleIds)
	if err == nil && deleted > 0 {
		s.mutex.Lock()
		for _, id := range roleIds {
			delete(s.idToRole, id)
		}
		s.mutex.Unlock()
	}
	return deleted, err
}

// @MappedFrom authAndUpdateAdminRoles
func (s *adminRoleService) AuthAndUpdateAdminRoles(ctx context.Context, requesterId int64, roleIds []int64, newName *string, permissions []permission.AdminPermission, rank *int) (int64, error) {
	for _, id := range roleIds {
		if id == RootRoleID {
			return 0, errors.New("the root admin is reserved and cannot be updated")
		}
	}
	// Check if requester rank is higher than target roles
	targetHighest, err := s.QueryHighestRankByRoleIds(ctx, roleIds)
	if err != nil {
		return 0, err
	}
	if targetHighest == nil {
		return 0, nil
	}
	higher, err := s.IsAdminRankHigherThanRank(ctx, requesterId, *targetHighest)
	if err != nil {
		return 0, err
	}
	if !higher {
		requesterRank, rankErr := s.QueryHighestRankByAdminId(ctx, requesterId)
		if rankErr != nil {
			return 0, rankErr
		}
		if requesterRank == nil {
			return 0, ErrRequesterNotExist
		}
		return 0, ErrPermissionDenied
	}
	// Check rank if updating rank
	if rank != nil {
		rankHigher, rankErr := s.IsAdminRankHigherThanRank(ctx, requesterId, *rank)
		if rankErr != nil {
			return 0, rankErr
		}
		if !rankHigher {
			return 0, ErrPermissionDenied
		}
	}
	// Verify requester has all permissions being assigned
	if len(permissions) > 0 {
		requesterPermissions, permErr := s.QueryPermissions(ctx, requesterId)
		if permErr != nil {
			return 0, permErr
		}
		if len(requesterPermissions) == 0 {
			return 0, ErrRequesterNotExist
		}
		permMap := make(map[permission.AdminPermission]bool)
		for _, p := range requesterPermissions {
			permMap[p] = true
		}
		for _, p := range permissions {
			if !permMap[p] {
				return 0, ErrPermissionDenied
			}
		}
	}
	return s.UpdateAdminRole(ctx, roleIds, newName, permissions, rank)
}

// @MappedFrom updateAdminRole
func (s *adminRoleService) UpdateAdminRole(ctx context.Context, roleIds []int64, newName *string, permissions []permission.AdminPermission, rank *int) (int64, error) {
	if len(roleIds) == 0 {
		return 0, nil
	}
	for _, id := range roleIds {
		if id == RootRoleID {
			return 0, errors.New("the root admin is reserved and cannot be updated")
		}
	}
	// No-op early return if all update fields are nil/empty (Java: Validator.areAllFalsy)
	if newName == nil && len(permissions) == 0 && rank == nil {
		return 0, nil
	}
	// Validate name if provided
	if newName != nil {
		name := *newName
		if name == "" {
			return 0, errors.New("name must not be blank")
		}
		if strings.Contains(name, " ") || strings.Contains(name, "\t") || strings.Contains(name, "\n") || strings.Contains(name, "\r") {
			return 0, errors.New("name must not contain whitespace")
		}
		if len(name) < MinRoleNameLimit || len(name) > MaxRoleNameLimit {
			return 0, errors.New("name length must be between 1 and 32")
		}
	}
	modified, err := s.repo.UpdateAdminRoles(ctx, roleIds, newName, permissions, rank)
	if err == nil && modified > 0 {
		s.mutex.Lock()
		for _, id := range roleIds {
			delete(s.idToRole, id) // Invalidate cache
		}
		s.mutex.Unlock()
	}
	return modified, err
}

func (s *adminRoleService) QueryRoleIdsByAdminId(ctx context.Context, adminId int64) ([]int64, error) {
	return s.adminService.QueryRoleIdsByAdminIds(ctx, []int64{adminId})
}

func (s *adminRoleService) QueryAdminRoles(ctx context.Context, ids []int64, names []string, includedPermissions []permission.AdminPermission, ranks []int, page *int, size *int) ([]*po.AdminRole, error) {
	roles, err := s.repo.FindAdminRoles(ctx, ids, names, includedPermissions, ranks, page, size)
	if err != nil {
		return nil, err
	}
	// Include root role if qualified (Java: isRootRoleQualified + startWith(getRootRole()))
	if isRootRoleQualified(ids, names, includedPermissions, ranks) {
		// Prepend root role
		result := make([]*po.AdminRole, 0, len(roles)+1)
		result = append(result, getRootRole())
		result = append(result, roles...)
		return result, nil
	}
	return roles, nil
}

// isRootRoleQualified checks whether the root role should be included in query results.
// @MappedFrom isRootRoleQualified in Java
func isRootRoleQualified(ids []int64, names []string, includedPermissions []permission.AdminPermission, ranks []int) bool {
	// If ids is specified, only include root if RootRoleID is in the list
	if len(ids) > 0 {
		found := false
		for _, id := range ids {
			if id == RootRoleID {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	// If names is specified, check if "ROOT" is in the list
	if len(names) > 0 {
		found := false
		for _, n := range names {
			if n == "ROOT" {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	// If includedPermissions is specified, check that root role contains all of them
	if len(includedPermissions) > 0 {
		root := getRootRole()
		rootPermSet := make(map[permission.AdminPermission]bool)
		for _, p := range root.Permissions {
			rootPermSet[p] = true
		}
		for _, p := range includedPermissions {
			if !rootPermSet[p] {
				return false
			}
		}
	}
	// If ranks is specified, check if root role's rank (rootRoleRank) is in the list
	if len(ranks) > 0 {
		found := false
		for _, r := range ranks {
			if r == rootRoleRank {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// @MappedFrom queryAndCacheRolesByRoleIdsAndRankGreaterThan
func (s *adminRoleService) QueryAndCacheRolesByRoleIdsAndRankGreaterThan(ctx context.Context, roleIds []int64, rankGreaterThan int) ([]*po.AdminRole, error) {
	// Bug fix: Missing empty roleIds check — Java returns empty if roleIds is empty
	if len(roleIds) == 0 {
		return nil, nil
	}
	// Bug fix: Missing roleIds immutability protection — create a copy before potentially modifying
	idsCopy := make([]int64, len(roleIds))
	copy(idsCopy, roleIds)

	// Bug fix: Missing root role handling — filter out RootRoleID from DB query (it's not in DB)
	filteredIds := make([]int64, 0, len(idsCopy))
	hasRoot := false
	for _, id := range idsCopy {
		if id == RootRoleID {
			hasRoot = true
		} else {
			filteredIds = append(filteredIds, id)
		}
	}

	var roles []*po.AdminRole
	if len(filteredIds) > 0 {
		var err error
		roles, err = s.repo.FindAdminRolesByIdsAndRankGreaterThan(ctx, filteredIds, &rankGreaterThan)
		if err != nil {
			return nil, err
		}
	}

	// Bug fix: Missing cache update — update idToRole cache for each fetched role
	if len(roles) > 0 {
		s.mutex.Lock()
		for _, role := range roles {
			s.idToRole[role.ID] = role
		}
		s.mutex.Unlock()
	}

	// Include root role if it was in the original roleIds and its rank qualifies
	if hasRoot && rootRoleRank > rankGreaterThan {
		roles = append(roles, rootRole)
	}

	return roles, nil
}

// @MappedFrom countAdminRoles
func (s *adminRoleService) CountAdminRoles(ctx context.Context, ids []int64, names []string, includedPermissions []permission.AdminPermission, ranks []int) (int64, error) {
	// Bug fix: Missing +1 for the root role — Java adds 1 because root role is not in DB
	count, err := s.repo.CountAdminRoles(ctx, ids, names, includedPermissions, ranks)
	if err != nil {
		return 0, err
	}
	return count + 1, nil
}

func (s *adminRoleService) QueryHighestRankByAdminId(ctx context.Context, adminId int64) (*int, error) {
	roleIds, err := s.adminService.QueryRoleIdsByAdminIds(ctx, []int64{adminId})
	if err != nil {
		return nil, err
	}
	return s.QueryHighestRankByRoleIds(ctx, roleIds)
}

// @MappedFrom queryHighestRankByRoleIds
func (s *adminRoleService) QueryHighestRankByRoleIds(ctx context.Context, roleIds []int64) (*int, error) {
	// Bug fix: Missing root role handling — if RootRoleID is in roleIds, return MAX_VALUE rank
	for _, id := range roleIds {
		if id == RootRoleID {
			r := rootRoleRank
			return &r, nil
		}
	}
	return s.repo.FindHighestRankByRoleIds(ctx, roleIds)
}

// @MappedFrom isAdminRankHigherThanRank
func (s *adminRoleService) IsAdminRankHigherThanRank(ctx context.Context, adminId int64, rank int) (bool, error) {
	highest, err := s.QueryHighestRankByAdminId(ctx, adminId)
	if err != nil {
		return false, err
	}
	// Bug fix: When admin has no roles, return error (requester does not exist)
	// instead of silently returning false
	if highest == nil {
		return false, s.adminService.ErrorRequesterNotExist()
	}
	return *highest > rank, nil
}

// @MappedFrom queryPermissions
func (s *adminRoleService) QueryPermissions(ctx context.Context, adminId int64) ([]permission.AdminPermission, error) {
	roleIds, err := s.adminService.QueryRoleIdsByAdminIds(ctx, []int64{adminId})
	if err != nil {
		return nil, err
	}

	// Bug fix: Use in-memory cache first, only query DB for uncached roles
	var uncachedRoleIds []int64
	s.mutex.RLock()
	for _, id := range roleIds {
		if _, ok := s.idToRole[id]; !ok {
			uncachedRoleIds = append(uncachedRoleIds, id)
		}
	}
	s.mutex.RUnlock()

	// Fetch uncached roles from DB
	if len(uncachedRoleIds) > 0 {
		freshRoles, err := s.repo.FindAdminRoles(ctx, uncachedRoleIds, nil, nil, nil, nil, nil)
		if err != nil {
			return nil, err
		}
		s.mutex.Lock()
		for _, role := range freshRoles {
			s.idToRole[role.ID] = role
		}
		s.mutex.Unlock()
	}

	// Bug fix: Include root role permissions — root role has all permissions
	permSet := make(map[permission.AdminPermission]struct{})
	for _, id := range roleIds {
		if id == RootRoleID {
			// Root role has all permissions
			for _, p := range permission.AllAdminPermissions {
				permSet[p] = struct{}{}
			}
			continue
		}
		s.mutex.RLock()
		role, ok := s.idToRole[id]
		s.mutex.RUnlock()
		if ok && role != nil {
			for _, p := range role.Permissions {
				permSet[p] = struct{}{}
			}
		}
	}

	var permissions []permission.AdminPermission
	for p := range permSet {
		permissions = append(permissions, p)
	}
	return permissions, nil
}

// AdminService maps to AdminService in Java.
// @MappedFrom AdminService
type AdminService interface {
	QueryRoleIdsByAdminIds(ctx context.Context, adminIds []int64) ([]int64, error)
	AuthAndAddAdmin(ctx context.Context, requesterId int64, loginName string, rawPassword string, displayName *string, roleIds []int64) (*po.Admin, error)
	AddAdmin(ctx context.Context, id *int64, loginName string, rawPassword string, displayName *string, roleIds []int64, upsert bool, registrationDate *time.Time) (*po.Admin, error)
	QueryAdmins(ctx context.Context, ids []int64, loginNames []string, roleIds []int64, page *int, size *int) ([]*po.Admin, error)
	AuthAndDeleteAdmins(ctx context.Context, requesterId int64, adminIds []int64) (int64, error)
	AuthAndUpdateAdmins(ctx context.Context, requesterId int64, targetAdminIds []int64, rawPassword *string, displayName *string, roleIds []int64) (int64, error)
	UpdateAdmins(ctx context.Context, targetAdminIds []int64, rawPassword *string, displayName *string, roleIds []int64) (int64, error)
	CountAdmins(ctx context.Context, ids []int64, roleIds []int64) (int64, error)
	ErrorRequesterNotExist() error
}

type adminService struct {
	idGen            *idgen.SnowflakeIdGenerator
	repo             repository.AdminRepository
	adminRoleService AdminRoleService

	// idToAdmin is an in-memory cache for admin objects, keyed by admin ID.
	idToAdmin map[int64]*po.Admin
	adminMu   sync.RWMutex
}

func NewAdminService(idGen *idgen.SnowflakeIdGenerator, repo repository.AdminRepository, adminRoleService AdminRoleService) AdminService {
	return &adminService{
		idGen:            idGen,
		repo:             repo,
		adminRoleService: adminRoleService,
		idToAdmin:        make(map[int64]*po.Admin),
	}
}

// @MappedFrom queryRoleIdsByAdminIds
func (s *adminService) QueryRoleIdsByAdminIds(ctx context.Context, adminIds []int64) ([]int64, error) {
	// Bug fix: Check in-memory cache first before querying DB
	var uncachedIds []int64
	s.adminMu.RLock()
	for _, id := range adminIds {
		if _, ok := s.idToAdmin[id]; !ok {
			uncachedIds = append(uncachedIds, id)
		}
	}
	s.adminMu.RUnlock()

	// Fetch uncached admins from DB
	if len(uncachedIds) > 0 {
		admins, err := s.repo.FindAdmins(ctx, uncachedIds, nil, nil, nil, nil)
		if err != nil {
			return nil, err
		}
		s.adminMu.Lock()
		for _, admin := range admins {
			s.idToAdmin[admin.ID] = admin
		}
		s.adminMu.Unlock()
	}

	// Collect role IDs from cache, with deduplication
	roleSet := make(map[int64]struct{})
	for _, id := range adminIds {
		s.adminMu.RLock()
		admin, ok := s.idToAdmin[id]
		s.adminMu.RUnlock()
		if ok && admin != nil {
			for _, rid := range admin.RoleIDs {
				roleSet[rid] = struct{}{}
			}
		}
	}

	// Bug fix: Deduplicate role IDs
	var roles []int64
	for rid := range roleSet {
		roles = append(roles, rid)
	}
	return roles, nil
}

// @MappedFrom authAndAddAdmin
func (s *adminService) AuthAndAddAdmin(ctx context.Context, requesterId int64, loginName string, rawPassword string, displayName *string, roleIds []int64) (*po.Admin, error) {
	// Bug fix: Validate that roleIds does not contain RootRoleID
	for _, id := range roleIds {
		if id == RootRoleID {
			return nil, errors.New("the root role ID is not allowed")
		}
	}

	// Bug fix: Requester-not-exist handling — check that requester actually exists
	requesterRank, err := s.adminRoleService.QueryHighestRankByAdminId(ctx, requesterId)
	if err != nil {
		return nil, err
	}
	if requesterRank == nil {
		return nil, s.ErrorRequesterNotExist()
	}

	if len(roleIds) > 0 {
		// Bug fix: Validate that all requested role IDs actually exist
		roles, err := s.adminRoleService.QueryAndCacheRolesByRoleIdsAndRankGreaterThan(ctx, roleIds, -1)
		if err != nil {
			return nil, err
		}
		foundIds := make(map[int64]struct{})
		for _, role := range roles {
			foundIds[role.ID] = struct{}{}
		}
		for _, id := range roleIds {
			if _, ok := foundIds[id]; !ok {
				return nil, errors.New("one or more role IDs do not exist")
			}
		}

		targetHighest, err := s.adminRoleService.QueryHighestRankByRoleIds(ctx, roleIds)
		if err != nil {
			return nil, err
		}
		if targetHighest != nil {
			if *requesterRank <= *targetHighest {
				return nil, ErrPermissionDenied
			}
		}
	}
	return s.AddAdmin(ctx, nil, loginName, rawPassword, displayName, roleIds, false, nil)
}

// @MappedFrom addAdmin
func (s *adminService) AddAdmin(ctx context.Context, id *int64, loginName string, rawPassword string, displayName *string, roleIds []int64, upsert bool, registrationDate *time.Time) (*po.Admin, error) {
	adminID := s.idGen.NextIncreasingId()
	if id != nil {
		adminID = *id
	}

	// Bug fix: Generate default loginName if empty
	if loginName == "" {
		loginName = randomAlphabetic(16)
	}

	// Bug fix: Generate default password if empty
	if rawPassword == "" {
		rawPassword = randomAlphabetic(10)
	}

	// Bug fix: Default displayName to loginName if nil or empty
	displayNameStr := ""
	if displayName != nil && *displayName != "" {
		displayNameStr = *displayName
	} else {
		displayNameStr = loginName
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Bug fix: Use caller-provided registrationDate, default to time.Now() if nil
	regDate := time.Now()
	if registrationDate != nil {
		regDate = *registrationDate
	}

	admin := &po.Admin{
		ID:               adminID,
		LoginName:        loginName,
		Password:         hashed,
		DisplayName:      &displayNameStr,
		RoleIDs:          roleIds,
		RegistrationDate: regDate,
	}

	if upsert {
		if err := s.repo.Upsert(ctx, admin); err != nil {
			return nil, err
		}
	} else {
		if err := s.repo.Insert(ctx, admin); err != nil {
			return nil, err
		}
	}

	// Bug fix: Update in-memory cache after successful DB write
	s.adminMu.Lock()
	s.idToAdmin[adminID] = admin
	s.adminMu.Unlock()

	return admin, nil
}

func (s *adminService) QueryAdmins(ctx context.Context, ids []int64, loginNames []string, roleIds []int64, page *int, size *int) ([]*po.Admin, error) {
	return s.repo.FindAdmins(ctx, ids, loginNames, roleIds, page, size)
}

// @MappedFrom authAndDeleteAdmins
func (s *adminService) AuthAndDeleteAdmins(ctx context.Context, requesterId int64, adminIds []int64) (int64, error) {
	// Bug fix: Validate adminIds is non-empty
	if len(adminIds) == 0 {
		return 0, errors.New("adminIds must not be empty")
	}
	for _, id := range adminIds {
		if id == RootAdminID {
			return 0, errors.New("the root admin cannot be deleted")
		}
	}
	// Bug fix: Requester-not-exist handling
	requesterRank, err := s.adminRoleService.QueryHighestRankByAdminId(ctx, requesterId)
	if err != nil {
		return 0, err
	}
	if requesterRank == nil {
		return 0, s.ErrorRequesterNotExist()
	}

	targetRoleIds, err := s.QueryRoleIdsByAdminIds(ctx, adminIds)
	if err != nil {
		return 0, err
	}
	targetHighest, err := s.adminRoleService.QueryHighestRankByRoleIds(ctx, targetRoleIds)
	if err != nil {
		return 0, err
	}
	if targetHighest != nil {
		if *requesterRank <= *targetHighest {
			return 0, ErrPermissionDenied
		}
	}
	return s.repo.DeleteAdmins(ctx, adminIds)
}

// @MappedFrom authAndUpdateAdmins
func (s *adminService) AuthAndUpdateAdmins(ctx context.Context, requesterId int64, targetAdminIds []int64, rawPassword *string, displayName *string, roleIds []int64) (int64, error) {
	// Bug fix: Early return when all update parameters are nil/empty
	if rawPassword == nil && displayName == nil && len(roleIds) == 0 {
		return 0, nil
	}

	// Bug fix: Check if requester is trying to update their own role IDs
	if len(roleIds) > 0 {
		for _, id := range targetAdminIds {
			if id == requesterId {
				return 0, errors.New("cannot update own role IDs")
			}
		}
	}

	// Bug fix: Requester-not-exist handling
	requesterRank, err := s.adminRoleService.QueryHighestRankByAdminId(ctx, requesterId)
	if err != nil {
		return 0, err
	}
	if requesterRank == nil {
		return 0, s.ErrorRequesterNotExist()
	}

	targetCurrentRoleIds, err := s.QueryRoleIdsByAdminIds(ctx, targetAdminIds)
	if err != nil {
		return 0, err
	}
	targetHighest, err := s.adminRoleService.QueryHighestRankByRoleIds(ctx, targetCurrentRoleIds)
	if err != nil {
		return 0, err
	}
	if targetHighest != nil {
		if *requesterRank <= *targetHighest {
			return 0, ErrPermissionDenied
		}
	}
	if len(roleIds) > 0 {
		// Bug fix: Validate that requested role IDs actually exist
		roles, err := s.adminRoleService.QueryAndCacheRolesByRoleIdsAndRankGreaterThan(ctx, roleIds, -1)
		if err != nil {
			return 0, err
		}
		foundIds := make(map[int64]struct{})
		for _, role := range roles {
			foundIds[role.ID] = struct{}{}
		}
		for _, id := range roleIds {
			if _, ok := foundIds[id]; !ok {
				return 0, errors.New("one or more role IDs do not exist")
			}
		}

		newTargetHighest, err := s.adminRoleService.QueryHighestRankByRoleIds(ctx, roleIds)
		if err != nil {
			return 0, err
		}
		if newTargetHighest != nil {
			if *requesterRank <= *newTargetHighest {
				return 0, ErrPermissionDenied
			}
		}
		// Bug fix: Match Java behavior — pass nil roleIds to UpdateAdmins when roleIds
		// are present (Java passes null roleIds to updateAdmins in this case)
		return s.UpdateAdmins(ctx, targetAdminIds, rawPassword, displayName, nil)
	}
	return s.UpdateAdmins(ctx, targetAdminIds, rawPassword, displayName, nil)
}

// @MappedFrom updateAdmins
func (s *adminService) UpdateAdmins(ctx context.Context, targetAdminIds []int64, rawPassword *string, displayName *string, roleIds []int64) (int64, error) {
	// Bug fix: Early return when all parameters are nil/empty
	if rawPassword == nil && displayName == nil && len(roleIds) == 0 {
		return 0, nil
	}

	var hashed []byte
	if rawPassword != nil {
		var err error
		hashed, err = bcrypt.GenerateFromPassword([]byte(*rawPassword), bcrypt.DefaultCost)
		if err != nil {
			return 0, err
		}
	}
	modified, err := s.repo.UpdateAdmins(ctx, targetAdminIds, hashed, displayName, roleIds)
	if err != nil {
		return modified, err
	}

	// Bug fix: Invalidate cache on successful update
	if modified > 0 {
		s.adminMu.Lock()
		for _, id := range targetAdminIds {
			delete(s.idToAdmin, id)
		}
		s.adminMu.Unlock()
	}
	return modified, nil
}

func (s *adminService) CountAdmins(ctx context.Context, ids []int64, roleIds []int64) (int64, error) {
	return s.repo.CountAdmins(ctx, ids, roleIds)
}

// @MappedFrom errorRequesterNotExist
func (s *adminService) ErrorRequesterNotExist() error {
	return ErrRequesterNotExist
}

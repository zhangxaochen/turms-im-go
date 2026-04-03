package controller

// AdminController maps to AdminController in Java.
// @MappedFrom AdminController
type AdminController struct {
}

// @MappedFrom checkLoginNameAndPassword()
func (c *AdminController) CheckLoginNameAndPassword() {
}

// @MappedFrom addAdmin(RequestContext requestContext, @RequestBody AddAdminDTO addAdminDTO)
func (c *AdminController) AddAdmin() {
}

// @MappedFrom queryAdmins(@QueryParam(required = false)
func (c *AdminController) QueryAdmins() {
}

// @MappedFrom updateAdmins(RequestContext requestContext, Set<Long> ids, @RequestBody UpdateAdminDTO updateAdminDTO)
func (c *AdminController) UpdateAdmins() {
}

// @MappedFrom deleteAdmins(RequestContext requestContext, Set<Long> ids)
func (c *AdminController) DeleteAdmins() {
}

// AdminPermissionController maps to AdminPermissionController in Java.
// @MappedFrom AdminPermissionController
type AdminPermissionController struct {
}

// @MappedFrom queryAdminPermissions()
func (c *AdminPermissionController) QueryAdminPermissions() {
}

// AdminRoleController maps to AdminRoleController in Java.
// @MappedFrom AdminRoleController
type AdminRoleController struct {
}

// @MappedFrom addAdminRole(RequestContext requestContext, @RequestBody AddAdminRoleDTO addAdminRoleDTO)
func (c *AdminRoleController) AddAdminRole() {
}

// @MappedFrom queryAdminRoles(@QueryParam(required = false)
func (c *AdminRoleController) QueryAdminRoles() {
}

// @MappedFrom updateAdminRole(RequestContext requestContext, Set<Long> ids, @RequestBody UpdateAdminRoleDTO updateAdminRoleDTO)
func (c *AdminRoleController) UpdateAdminRole() {
}

// @MappedFrom deleteAdminRoles(RequestContext requestContext, Set<Long> ids)
func (c *AdminRoleController) DeleteAdminRoles() {
}

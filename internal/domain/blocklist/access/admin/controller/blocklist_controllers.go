package controller

// IpBlocklistController maps to IpBlocklistController.java
// @MappedFrom IpBlocklistController
type IpBlocklistController struct {
}

// @MappedFrom addBlockedIps(@RequestBody AddBlockedIpsDTO addBlockedIpsDTO)
func (c *IpBlocklistController) AddBlockedIps() {
}

// @MappedFrom queryBlockedIps(Set<String> ids)
func (c *IpBlocklistController) QueryBlockedIpsByIds() {
}

// @MappedFrom queryBlockedIps(int page, @QueryParam(required = false)
func (c *IpBlocklistController) QueryBlockedIpsByPage() {
}

// @MappedFrom deleteBlockedIps(@QueryParam(required = false)
func (c *IpBlocklistController) DeleteBlockedIps() {
}

// UserBlocklistController maps to UserBlocklistController.java
// @MappedFrom UserBlocklistController
type UserBlocklistController struct {
}

// @MappedFrom addBlockedUserIds(@RequestBody AddBlockedUserIdsDTO addBlockedUserIdsDTO)
func (c *UserBlocklistController) AddBlockedUserIds() {
}

// @MappedFrom deleteBlockedUserIds(@QueryParam(required = false)
func (c *UserBlocklistController) DeleteBlockedUserIds() {
}

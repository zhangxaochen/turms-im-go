package controller

// MemberController maps to MemberController.java
// @MappedFrom MemberController
type MemberController struct {
}

// @MappedFrom queryMembers()
func (c *MemberController) QueryMembers() {
}

// @MappedFrom removeMembers(List<String> ids)
func (c *MemberController) RemoveMembers() {
}

// @MappedFrom addMember(@RequestBody AddMemberDTO addMemberDTO)
func (c *MemberController) AddMember() {
}

// @MappedFrom updateMember(String id, @RequestBody UpdateMemberDTO updateMemberDTO)
func (c *MemberController) UpdateMember() {
}

// @MappedFrom queryLeader()
func (c *MemberController) QueryLeader() {
}

// @MappedFrom electNewLeader(@QueryParam(required = false)
func (c *MemberController) ElectNewLeader() {
}

// SettingController maps to SettingController.java
// @MappedFrom SettingController
type SettingController struct {
}

// @MappedFrom queryClusterSettings(boolean queryLocalSettings, boolean onlyMutable)
func (c *SettingController) QueryClusterSettings() {
}

// @MappedFrom updateClusterSettings(boolean reset, boolean updateLocalSettings, @RequestBody(required = false)
func (c *SettingController) UpdateClusterSettings() {
}

// @MappedFrom queryClusterConfigMetadata(boolean queryLocalSettings, boolean onlyMutable, boolean withValue)
func (c *SettingController) QueryClusterConfigMetadata() {
}

package profile

func (p *Profile) SetProfileGetter(profileGetter ProfileGetter) {
	getProfileDefinition = profileGetter
}

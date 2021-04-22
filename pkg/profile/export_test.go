package profile

func (p *Profile) SetProfileGetter(profileGetter ProfileGetter) {
	p.getProfileDefinition = profileGetter
}

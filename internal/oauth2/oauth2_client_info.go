package oauth2

type oauth2ClientInfoImpl struct {
	domain string
}

func (i *oauth2ClientInfoImpl) GetID() string {
	return "i.client.ID"
}

func (i *oauth2ClientInfoImpl) GetSecret() string {
	return "string(i.client.ClientSecret)"
}

func (i *oauth2ClientInfoImpl) GetDomain() string {
	return i.domain
}

func (i *oauth2ClientInfoImpl) IsPublic() bool {
	return false
}

func (i *oauth2ClientInfoImpl) GetUserID() string {
	return "i.client.BelongsToUser"
}

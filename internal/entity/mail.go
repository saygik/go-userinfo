package entity

type MailBoxDelegates struct {
	Id                       string `json:"id"`
	MailName                 string `json:"mailName"`
	Mail                     string `json:"mail"`
	DelegateUser             string `json:"delegateUser"`
	DelegateUserName         string `json:"delegateUserName"`
	DelegateUserMail         string `json:"delegateUserMail"`
	DelegateUserAccessRights string `json:"delegateUserAccessRights"`
}

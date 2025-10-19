package utils

var (
	IdVerify               = Rules{"ID": []string{NotEmpty()}}
	ApiVerify              = Rules{"Path": {NotEmpty()}, "Description": {NotEmpty()}, "ApiGroup": {NotEmpty()}, "Method": {NotEmpty()}}
	MenuVerify             = Rules{"Path": {NotEmpty()}, "Name": {NotEmpty()}, "Component": {NotEmpty()}, "Sort": {Ge("0")}}
	MenuMetaVerify         = Rules{"Title": {NotEmpty()}}
	LoginVerify            = Rules{"Username": {NotEmpty()}, "Password": {NotEmpty()}}
	RegisterVerify         = Rules{"Username": {NotEmpty()}, "NickName": {NotEmpty()}, "Password": {NotEmpty()}, "AuthorityId": {NotEmpty()}}
	PageInfoVerify         = Rules{"Page": {NotEmpty()}, "PageSize": {NotEmpty()}}
	CustomerVerify         = Rules{"CustomerName": {NotEmpty()}, "CustomerPhoneData": {NotEmpty()}}
	AutoCodeVerify         = Rules{"Abbreviation": {NotEmpty()}, "StructName": {NotEmpty()}, "PackageName": {NotEmpty()}}
	AutoPackageVerify      = Rules{"PackageName": {NotEmpty()}}
	AuthorityVerify        = Rules{"AuthorityId": {NotEmpty()}, "AuthorityName": {NotEmpty()}}
	AuthorityIdVerify      = Rules{"AuthorityId": {NotEmpty()}}
	OldAuthorityVerify     = Rules{"OldAuthorityId": {NotEmpty()}}
	ChangePasswordVerify   = Rules{"Password": {NotEmpty()}, "NewPassword": {NotEmpty()}}
	SetUserAuthorityVerify = Rules{"AuthorityId": {NotEmpty()}}
)

var PaymentAccountVerify = Rules{
	"Name":           {NotEmpty()},
	"ProviderCode":   {NotEmpty()},
	"AccountType":    {NotEmpty(), InEnum("personal", "enterprise")},
	"Config":         {NotEmpty()},
	"Status":         {InEnum("active", "inactive", "maintenance")},
	"Weight":         {Ge("1"), Le("100")},
	"MaxDailyAmount": {Ge("0")},
}

// PaymentAccountUpdateVerify 支付账号更新验证规则
var PaymentAccountUpdateVerify = Rules{
	"ID":             {NotEmpty()},
	"Name":           {NotEmpty()},
	"Config":         {NotEmpty()},
	"Status":         {InEnum("active", "inactive", "maintenance")},
	"Weight":         {Ge("1"), Le("100")},
	"MaxDailyAmount": {Ge("0")},
}

// PaymentAccountSearchVerify 支付账号搜索验证规则
var PaymentAccountSearchVerify = Rules{
	"Page":     {Gt("0")},
	"PageSize": {Gt("0"), Le("100")},
}

// PaymentAccountDeleteVerify 支付账号删除验证规则
var PaymentAccountDeleteVerify = Rules{
	"ID": {NotEmpty()},
}

// PaymentAccountBatchDeleteVerify 批量删除支付账号验证规则
var PaymentAccountBatchDeleteVerify = Rules{
	"IDs": {NotEmpty()},
}

// PaymentAccountBatchStatusVerify 批量更新状态验证规则
var PaymentAccountBatchStatusVerify = Rules{
	"IDs":    {NotEmpty()},
	"Status": {NotEmpty(), InEnum("active", "inactive", "maintenance", "deleted")},
}

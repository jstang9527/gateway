package public

const (
	// ValidatorKey ...
	ValidatorKey = "ValidatorKey"
	// TranslatorKey ...
	TranslatorKey = "TranslatorKey"
	// AdminSessionInfoKey ...
	AdminSessionInfoKey = "AdminSessionInfoKey"
	// LoadTypeHTTP ...
	LoadTypeHTTP = 0
	// LoadTypeTCP ...
	LoadTypeTCP = 1
	// LoadTypeGRPC ...
	LoadTypeGRPC = 2

	// HTTPRuleTypePrefixURL 前缀接入
	HTTPRuleTypePrefixURL = 0
	// HTTPRuleTypeDomain 域名接入
	HTTPRuleTypeDomain = 1
)

var (
	// LoadTypeMap ..
	LoadTypeMap = map[int]string{
		LoadTypeHTTP: "HTTP",
		// LoadTypeTCP ...
		LoadTypeTCP: "TCP",
		// LoadTypeGRPC ...
		LoadTypeGRPC: "GRPC",
	}
)

package codes

import (
	"blogrpc/core/codes"
	"blogrpc/core/errors"
)

// The canonical error codes used by member gRPC service.
const (
	// MemberPropertyNameInvalid indicates the name of a member property is invalid.
	MemberPropertyNameInvalid           codes.Code = 2000001
	ExceedMemberProperties              codes.Code = 2000002
	UpdateMemberFail                    codes.Code = 2000003
	CardNotFound                        codes.Code = 2000004
	CardIsAutoUpgrade                   codes.Code = 2000005
	MemberNotFound                      codes.Code = 2000007
	MemberBlocked                       codes.Code = 2000008
	InsertMemberFailed                  codes.Code = 2000009
	MembershipCardNameExisted           codes.Code = 2000011
	IllegalMembershipCardScore          codes.Code = 2000012
	MembershipCardScoreOverlaped        codes.Code = 2000013
	DefaultMembershipCard               codes.Code = 2000014
	MembershipCardUsed                  codes.Code = 2000015
	ScoreRuleNotFound                   codes.Code = 2000016
	ScoreRuleRequiredFields             codes.Code = 2000017
	ScoreRuleCouponNotFound             codes.Code = 2000018
	ScoreRuleCouponNotAvailable         codes.Code = 2000019
	ScoreRuleRequiredScore              codes.Code = 2000020
	ScoreRuleRequireTriggerType         codes.Code = 2000021
	ScoreRuleRequireProperties          codes.Code = 2000022
	RuleIsNotEnabled                    codes.Code = 2000023
	InValidMember                       codes.Code = 2000024
	MemberHasExceededLimit              codes.Code = 2000025
	IncScoreFail                        codes.Code = 2000026
	NotSupportedLimitType               codes.Code = 2000027
	ErrorScoreRuleType                  codes.Code = 2000028
	IncRewardCountFail                  codes.Code = 2000029
	InvalidOrigin                       codes.Code = 2000030
	InvalidMemberId                     codes.Code = 2000031
	InvalidChannelId                    codes.Code = 2000032
	InvalidRuleCode                     codes.Code = 2000033
	InvalidRuleName                     codes.Code = 2000034
	MissingCodeOrRuleName               codes.Code = 2000035
	CreateScoreHistoryFail              codes.Code = 2000036
	InvalidMemberIds                    codes.Code = 2000037
	ScoreRuleRequireId                  codes.Code = 2000038
	ScoreRuleDefaultNameUneditable      codes.Code = 2000039
	ScoreRuleInvisibleProperties        codes.Code = 2000040
	ScoreRuleInvalidLimitType           codes.Code = 2000041
	ScoreRuleInvalidName                codes.Code = 2000042
	ScoreRuleCodeNotUnique              codes.Code = 2000043
	ScoreHistoryNotFound                codes.Code = 2000044
	UndeletableScoreRule                codes.Code = 2000045
	InvalidMongoId                      codes.Code = 2000046
	DeleteFailed                        codes.Code = 2000047
	InvalidEmail                        codes.Code = 2000048
	MissingPhoneEmailOpenIdUnionId      codes.Code = 2000049
	PhoneBlocked                        codes.Code = 2000050
	UnknownProperty                     codes.Code = 2000051
	OriginOpenIdExist                   codes.Code = 2000052
	MemberPropertyRequired              codes.Code = 2000053
	PhoneNotUnique                      codes.Code = 2000054
	EmailNotUnique                      codes.Code = 2000055
	MissingOpenId                       codes.Code = 2000056
	OpenIdExist                         codes.Code = 2000057
	UnionIdExist                        codes.Code = 2000058
	CardNumberExist                     codes.Code = 2000059
	CreateBlacklistFail                 codes.Code = 2000060
	BlacklistHasExisted                 codes.Code = 2000061
	DeleteBlacklistFail                 codes.Code = 2000062
	GetBlacklistFail                    codes.Code = 2000063
	MemberPropertyNotFound              codes.Code = 2000064
	UpdateMemberPropertyFail            codes.Code = 2000065
	MemberHasBound                      codes.Code = 2000066
	MemberPropertyExceed100Count        codes.Code = 2000067
	MemberPropertyTypeNotExist          codes.Code = 2000068
	MemberPropertyShouldBeUnique        codes.Code = 2000069
	MemberPropertyIsDefaultShouldBeTrue codes.Code = 2000070
	MemberPropertyAlreadyExist          codes.Code = 2000071
	MemberPropertyPropertyIdInvalid     codes.Code = 2000072
	DefaultTagGroupNotFound             codes.Code = 2000073
	ResourceNotExist                    codes.Code = 2000074
	TagGroupNotFound                    codes.Code = 2000075
	TagNameExist                        codes.Code = 2000076
	TagNotFound                         codes.Code = 2000077
	TagNameRequired                     codes.Code = 2000078
	TagGroupHasExisted                  codes.Code = 2000079
	CreateTagGroupFail                  codes.Code = 2000080
	UpdateTagGroupFail                  codes.Code = 2000081
	DeleteTagGroupFail                  codes.Code = 2000082
	CommonMissingRequiredFields         codes.Code = 2000083
	SuspiciousRuleNotFound              codes.Code = 2000084
	SetSuspiciousRuleInRedisFail        codes.Code = 2000085
	UpdateSuspiciousRuleFail            codes.Code = 2000086
	UpsertSuspiciousRuleFail            codes.Code = 2000087
	BlockedStatusRemarkRequired         codes.Code = 2000088
	InvalidScoreResetType               codes.Code = 2000090
	MemberInfoLogNotFound               codes.Code = 2000091
	SocialFilterRequired                codes.Code = 2000092
	MergeMemberNotFound                 codes.Code = 2000093
	PortalOriginRepeated                codes.Code = 2000095
	GeneratePropertyIdFail              codes.Code = 2000096
	MissingChannelId                    codes.Code = 2000097
	InvalidPhone                        codes.Code = 2000098
	MemberPropertyValueInvalid          codes.Code = 2000099
	InputNotUnique                      codes.Code = 2000100
	CannotModifyUnknownStage            codes.Code = 2000101
	MemberStageNotFound                 codes.Code = 2000102
	UpdateMemberStageFail               codes.Code = 2000103
	InvalidPropertyId                   codes.Code = 2000104
	InvalidInformationRuleValue         codes.Code = 2000105
	InvalidInformationRuleOperator      codes.Code = 2000106
	InvalidMembershipCard               codes.Code = 2000107
	UpdateMemberDisabledStatusFail      codes.Code = 2000108
	MemberAddressNotFound               codes.Code = 2000109
	InvalidEventProperties              codes.Code = 2000110
	MissingMemberId                     codes.Code = 2000111
	GetPhoneFromMiniProgramFail         codes.Code = 2000112
	MissingOrigin                       codes.Code = 2000113
	MissingOpenIdOrUnionId              codes.Code = 2000114
	MainMemberExistsInSubMember         codes.Code = 2000115
	FailedToBindAnonymousToMember       codes.Code = 2000116
	MemberDisabled                      codes.Code = 2000117
	MemberIsProcessing                  codes.Code = 2000118
	WaitingForPreviousJob               codes.Code = 2000119
	FailedToGetMemberScoreSyncRecords   codes.Code = 2000120
	CountMemberForTagFailed             codes.Code = 2000121
	FailedToGetTaobaoPointChangeMsgs    codes.Code = 2000122
	MustContainsSocialInfo              codes.Code = 2000123
	TagNotMoreThanOne                   codes.Code = 2000124
	InvalidPropertyInfo                 codes.Code = 2000125
	FailedToCreateMember                codes.Code = 2000126
	TooManyUnionIds                     codes.Code = 2000127
	InvalidPageSize                     codes.Code = 2000128
	FailedToGetMemberEventLog           codes.Code = 2000129
	InvalidCursor                       codes.Code = 2000130
	ScoreRuleRequiredSpent              codes.Code = 2000131
	InvalidCardId                       codes.Code = 2000132
	InvalidMemberDay                    codes.Code = 2000133
	ScoreRuleRewardHistoryNotFound      codes.Code = 2000134
	MissingUnionId                      codes.Code = 2000135
	MemberDayNeedBonus                  codes.Code = 2000136
	MemberGroupNotFound                 codes.Code = 2000137
)

var codeText = map[codes.Code]string{
	MemberPropertyNameInvalid:           "",
	ExceedMemberProperties:              "",
	UpdateMemberFail:                    "Failed to update member",
	CardNotFound:                        "",
	CardIsAutoUpgrade:                   "",
	MemberNotFound:                      "Member not found",
	MemberBlocked:                       "Member is blocked",
	InsertMemberFailed:                  "Insert member failed",
	MembershipCardNameExisted:           "",
	IllegalMembershipCardScore:          "",
	MembershipCardScoreOverlaped:        "",
	DefaultMembershipCard:               "",
	MembershipCardUsed:                  "",
	ScoreRuleNotFound:                   "Score rule is not found",
	ScoreRuleRequiredFields:             "Lack of required fields for score rule",
	ScoreRuleCouponNotFound:             "Coupon for score rule is not found",
	ScoreRuleCouponNotAvailable:         "Coupon for score rule is not available",
	ScoreRuleRequiredScore:              "Reward type for score requires score field",
	ScoreRuleRequireTriggerType:         "Score rule for birthday requires trigger time field",
	ScoreRuleRequireProperties:          "Score rule for perfect information requires properties field",
	RuleIsNotEnabled:                    "rule is not enabled",
	InValidMember:                       "Members do not conform to the default rules or has been rewarded",
	MemberHasExceededLimit:              "member has exceeded limit",
	IncScoreFail:                        "Increase score to member fail",
	NotSupportedLimitType:               "Get not supported score rule limit type",
	ErrorScoreRuleType:                  "Error score rule reward type",
	IncRewardCountFail:                  "Fail to incr member reward count",
	InvalidOrigin:                       "Validate origin fail by invalid origin",
	InvalidMemberId:                     "Find member fail by invalid memberId",
	InvalidChannelId:                    "Get channel fail by invalid ChannelId",
	InvalidRuleCode:                     "Find scoreRule fail by invalid RuleCode",
	InvalidRuleName:                     "Find scoreRule fail by invalid RuleName",
	MissingCodeOrRuleName:               "Missing parameter code or ruleName",
	CreateScoreHistoryFail:              "Create scoreHistory fail",
	InvalidMemberIds:                    "Find member fail by invalid memberIds",
	ScoreRuleRequireId:                  "Update score rule require ID",
	ScoreRuleDefaultNameUneditable:      "Default score rule name can not be edited",
	ScoreRuleInvisibleProperties:        "Selected properties for information perfection has invisible ones",
	ScoreRuleInvalidLimitType:           "The type of limit is not unlimited, day or total",
	ScoreRuleInvalidName:                "Only default name can be perfect_information, birthday or first_card",
	ScoreRuleCodeNotUnique:              "Score rule code is not unique",
	ScoreHistoryNotFound:                "Score history is not found",
	UndeletableScoreRule:                "Default score rule can not be delete",
	InvalidMongoId:                      "Invalid MongoId",
	DeleteFailed:                        "Delete failed",
	InvalidEmail:                        "Invalid email",
	MissingPhoneEmailOpenIdUnionId:      "Missing param phone, email, openId or unionId",
	PhoneBlocked:                        "Phone is blocked",
	UnknownProperty:                     "Unknown property",
	OriginOpenIdExist:                   "OpenId in origin has existed",
	MemberPropertyRequired:              "Member property is required",
	PhoneNotUnique:                      "Phone is not unique",
	EmailNotUnique:                      "Email is not unique",
	MissingOpenId:                       "OpenId can not be empty",
	OpenIdExist:                         "OpenId has existed",
	UnionIdExist:                        "UnionId has existed",
	CardNumberExist:                     "Card number has existed",
	CreateBlacklistFail:                 "Create blacklist fail",
	BlacklistHasExisted:                 "Blacklist has existed",
	DeleteBlacklistFail:                 "Delete blacklist fail",
	GetBlacklistFail:                    "Get blacklist fail",
	MemberPropertyNotFound:              "Member property is not found",
	UpdateMemberPropertyFail:            "Update member property fail",
	MemberHasBound:                      "Member has been Bound the channel",
	MemberPropertyExceed100Count:        "Exceed max count (100) of member properties",
	MemberPropertyTypeNotExist:          "Type does not exist",
	MemberPropertyShouldBeUnique:        "Should be a unique property",
	MemberPropertyIsDefaultShouldBeTrue: "Default property's isDefault value should be true",
	MemberPropertyAlreadyExist:          "Propery already exists and is unique",
	MemberPropertyPropertyIdInvalid:     "PropertyId is invalid - incorrect pattern",
	DefaultTagGroupNotFound:             "The default tag group is not found",
	ResourceNotExist:                    "Resource does not exist",
	TagGroupNotFound:                    "The tag group is not found",
	TagNameExist:                        "The tag name has existed",
	TagNotFound:                         "The tag is not found",
	TagNameRequired:                     "The tag name is required",
	TagGroupHasExisted:                  "Tag group has been existed",
	CreateTagGroupFail:                  "Create tag group fail",
	UpdateTagGroupFail:                  "Update tag group fail",
	DeleteTagGroupFail:                  "Delete tag group fail",
	CommonMissingRequiredFields:         "Missing required fields",
	SuspiciousRuleNotFound:              "Get suspiciousRule fail by invalid id",
	SetSuspiciousRuleInRedisFail:        "Set suspiciousRule in redis fail",
	UpdateSuspiciousRuleFail:            "Update suspiciousRule fail",
	UpsertSuspiciousRuleFail:            "Upsert suspiciousRule fail",
	BlockedStatusRemarkRequired:         "Remark is required",
	InvalidScoreResetType:               "Score reset type is invalid",
	MemberInfoLogNotFound:               "Member info log not found",
	SocialFilterRequired:                "You should set at least one social filter",
	MergeMemberNotFound:                 "Invalid IDs for merged members",
	PortalOriginRepeated:                "Repeat set portal origin",
	GeneratePropertyIdFail:              "Generate property id fail",
	MissingChannelId:                    "ChannelId can not be empty",
	InvalidPhone:                        "Invalid phone",
	MemberPropertyValueInvalid:          "Invalid member property value",
	InputNotUnique:                      "Input property is not unique",
	CannotModifyUnknownStage:            "Can not modify 'unknown' stage",
	MemberStageNotFound:                 "Member stage not found",
	UpdateMemberStageFail:               "Update memberStage fail",
	InvalidPropertyId:                   "Member property id is invalid",
	InvalidInformationRuleValue:         "Information rule's value is invalid",
	InvalidInformationRuleOperator:      "Information rule's operator is invalid",
	InvalidMembershipCard:               "Membership card is invalid",
	UpdateMemberDisabledStatusFail:      "Update member disabled status fail",
	MemberAddressNotFound:               "Member address not found",
	InvalidEventProperties:              "Invalid event properties",
	MissingMemberId:                     "Missing params memberId or id",
	GetPhoneFromMiniProgramFail:         "Get phone from mini program fail",
	MissingOrigin:                       "Member must have an origin",
	MissingOpenIdOrUnionId:              "Missing openId or unionId",
	MainMemberExistsInSubMember:         "Main member's id exists in merged members' id list",
	FailedToBindAnonymousToMember:       "Failed to bind anonymous to member",
	MemberDisabled:                      "Member is disabled, can not provide member card",
	MemberIsProcessing:                  "Member is processing",
	WaitingForPreviousJob:               "Previous job is still processing",
	FailedToGetMemberScoreSyncRecords:   "Failed to get member score sync records",
	CountMemberForTagFailed:             "Failed to count member for tag",
	FailedToGetTaobaoPointChangeMsgs:    "Failed to get taobao point change messages",
	MustContainsSocialInfo:              "Socials must contains origin and channel",
	TagNotMoreThanOne:                   "Can only count one tag at a time",
	InvalidPropertyInfo:                 "Invalid property info",
	FailedToCreateMember:                "Failed to create member",
	TooManyUnionIds:                     "Params unionIds exceeds limit (100)",
	InvalidPageSize:                     "Invalid page size",
	FailedToGetMemberEventLog:           "Failed to get member event log",
	InvalidCursor:                       "Invalid cursor",
	ScoreRuleRequiredSpent:              "Consumption should set spent",
	InvalidCardId:                       "Invalid card id",
	InvalidMemberDay:                    "Invalid member day",
	ScoreRuleRewardHistoryNotFound:      "Score rule reward history not found",
	MissingUnionId:                      "UnionId can not be empty",
	MemberDayNeedBonus:                  "Member day need set bonus rate",
	MemberGroupNotFound:                 "Member group not found",
}

func NewError(code codes.Code) error {
	return errors.NewRPCError(code, codeText[code])
}

func NewErrorWithExtra(code codes.Code, extra map[string]interface{}) error {
	return errors.NewRPCErrorWithExtra(code, codeText[code], extra)
}

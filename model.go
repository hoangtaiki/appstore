package appstore

import (
	"github.com/golang-jwt/jwt/v5"
)

// OrderLookupResponse https://developer.apple.com/documentation/appstoreserverapi/orderlookupresponse
type OrderLookupResponse struct {
	Status             int      `json:"status"`
	SignedTransactions []string `json:"signedTransactions"`
}

type Environment string

// Environment https://developer.apple.com/documentation/appstoreserverapi/environment
const (
	Sandbox    Environment = "Sandbox"
	Production Environment = "Production"
)

// HistoryResponse https://developer.apple.com/documentation/appstoreserverapi/historyresponse
type HistoryResponse struct {
	AppAppleId         int64       `json:"appAppleId"`
	BundleId           string      `json:"bundleId"`
	Environment        Environment `json:"environment"`
	HasMore            bool        `json:"hasMore"`
	Revision           string      `json:"revision"`
	SignedTransactions []string    `json:"signedTransactions"`
}

// TransactionInfoResponse https://developer.apple.com/documentation/appstoreserverapi/transactioninforesponse
type TransactionInfoResponse struct {
	SignedTransactionInfo string `json:"signedTransactionInfo"`
}

// RefundLookupResponse same as the RefundHistoryResponse https://developer.apple.com/documentation/appstoreserverapi/refundhistoryresponse
type RefundLookupResponse struct {
	HasMore            bool     `json:"hasMore"`
	Revision           string   `json:"revision"`
	SignedTransactions []string `json:"signedTransactions"`
}

// StatusResponse https://developer.apple.com/documentation/appstoreserverapi/get_all_subscription_statuses
type StatusResponse struct {
	Environment Environment                       `json:"environment"`
	AppAppleId  int64                             `json:"appAppleId"`
	BundleId    string                            `json:"bundleId"`
	Data        []SubscriptionGroupIdentifierItem `json:"data"`
}

type SubscriptionGroupIdentifierItem struct {
	SubscriptionGroupIdentifier string                 `json:"subscriptionGroupIdentifier"`
	LastTransactions            []LastTransactionsItem `json:"lastTransactions"`
}

type LastTransactionsItem struct {
	OriginalTransactionId string `json:"originalTransactionId"`
	Status                int32  `json:"status"`
	SignedRenewalInfo     string `json:"signedRenewalInfo"`
	SignedTransactionInfo string `json:"signedTransactionInfo"`
}

// MassExtendRenewalDateRequest https://developer.apple.com/documentation/appstoreserverapi/massextendrenewaldaterequest
type MassExtendRenewalDateRequest struct {
	RequestIdentifier      string   `json:"requestIdentifier"`
	ExtendByDays           int32    `json:"extendByDays"`
	ExtendReasonCode       int32    `json:"extendReasonCode"`
	ProductId              string   `json:"productId"`
	StorefrontCountryCodes []string `json:"storefrontCountryCodes"`
}

type DeliveryStatus string

// DeliveryStatus https://developer.apple.com/documentation/appstoreserverapi/deliverystatus
const (
	DELIVERED                 DeliveryStatus = "DELIVERED"
	UNDELIVERED_QUALITY_ISSUE DeliveryStatus = "UNDELIVERED_QUALITY_ISSUE"
	UNDELIVERED_WRONG_ITEM    DeliveryStatus = "UNDELIVERED_WRONG_ITEM"
	UNDELIVERED_SERVER_OUTAGE DeliveryStatus = "UNDELIVERED_SERVER_OUTAGE"
	UNDELIVERED_OTHER         DeliveryStatus = "UNDELIVERED_OTHER"
)

type RefundPreference string

// RefundPreference https://developer.apple.com/documentation/appstoreserverapi/refundpreference
const (
	DECLINE        RefundPreference = "DECLINE"
	GRANT_FULL     RefundPreference = "GRANT_FULL"
	GRANT_PRORATED RefundPreference = "GRANT_PRORATED"
)

// ConsumptionRequestBody https://developer.apple.com/documentation/appstoreserverapi/consumptionrequest
type ConsumptionRequest struct {
	CustomerConsented     bool             `json:"customerConsented"`
	ConsumptionPercentage int32            `json:"consumptionPercentage"`
	DeliveryStatus        DeliveryStatus   `json:"deliveryStatus"`
	RefundPreference      RefundPreference `json:"refundPreference"`
	SampleContentProvided bool             `json:"sampleContentProvided"`
}

// ConsumptionRequestBody https://developer.apple.com/documentation/appstoreserverapi/consumptionrequest
type ConsumptionRequestBody struct {
	AccountTenure            int32  `json:"accountTenure"`
	AppAccountToken          string `json:"appAccountToken"`
	ConsumptionStatus        int32  `json:"consumptionStatus"`
	CustomerConsented        bool   `json:"customerConsented"`
	DeliveryStatus           int32  `json:"deliveryStatus"`
	LifetimeDollarsPurchased int32  `json:"lifetimeDollarsPurchased"`
	LifetimeDollarsRefunded  int32  `json:"lifetimeDollarsRefunded"`
	Platform                 int32  `json:"platform"`
	PlayTime                 int32  `json:"playTime"`
	SampleContentProvided    bool   `json:"sampleContentProvided"`
	UserStatus               int32  `json:"userStatus"`
	RefundPreference         int32  `json:"refundPreference"`
}

type AdvancedCommerceDescriptors struct {
	Description string `json:"description"`
	DisplayName string `json:"displayName"`
}

type AdvancedCommercePriceIncreaseInfo struct {
	DependentSKUs []string `json:"dependentSKUs"`
	Price         int64    `json:"price"`
	Status        string   `json:"status"`
}

type AdvancedCommerceOffer struct {
	Period      string `json:"period"`
	PeriodCount int32  `json:"periodCount"`
	Price       int64  `json:"price"`
	Reason      string `json:"reason"`
}

type AdvancedCommerceRenewalItems struct {
	SKU               string                            `json:"SKU"`
	Description       string                            `json:"description"`
	DisplayName       string                            `json:"displayName"`
	Offer             AdvancedCommerceOffer             `json:"offer"`
	Price             int64                             `json:"price"`
	PriceIncreaseInfo AdvancedCommercePriceIncreaseInfo `json:"priceIncreaseInfo"`
}

// AdvancedCommerceRenewalInfo https://developer.apple.com/documentation/appstoreserverapi/advancedcommercerenewalinfo
type AdvancedCommerceRenewalInfo struct {
	ConsistencyToken   string                         `json:"consistencyToken"`
	Descriptors        AdvancedCommerceDescriptors    `json:"descriptors"`
	Items              []AdvancedCommerceRenewalItems `json:"items"`
	Period             string                         `json:"period"`
	RequestReferenceId string                         `json:"requestReferenceId"`
	TaxCode            string                         `json:"taxCode"`
}

type RenewalCommitmentInfo struct {
	CommitmentAutoRenewProductId     string          `json:"commitmentAutoRenewProductId"`
	CommitmentAutoRenewStatus        int32           `json:"commitmentAutoRenewStatus"`
	CommitmentRenewalBillingPlanType BillingPlanType `json:"commitmentRenewalBillingPlanType"`
	CommitmentRenewalDate            int64           `json:"commitmentRenewalDate"`
	CommitmentRenewalPrice           int64           `json:"commitmentRenewalPrice"`
}

// JWSRenewalInfoDecodedPayload https://developer.apple.com/documentation/appstoreserverapi/jwsrenewalinfodecodedpayload
type JWSRenewalInfoDecodedPayload struct {
	AppAccountToken             string                      `json:"appAccountToken,omitempty"`
	AppTransactionId            string                      `json:"appTransactionId,omitempty"`
	AutoRenewProductId          string                      `json:"autoRenewProductId"`
	AutoRenewStatus             int32                       `json:"autoRenewStatus"`
	Environment                 Environment                 `json:"environment"`
	ExpirationIntent            int32                       `json:"expirationIntent"`
	GracePeriodExpiresDate      int64                       `json:"gracePeriodExpiresDate"`
	IsInBillingRetryPeriod      *bool                       `json:"isInBillingRetryPeriod"`
	OfferIdentifier             string                      `json:"offerIdentifier"`
	OfferType                   int32                       `json:"offerType"`
	OfferPeriod                 string                      `json:"offerPeriod"`
	OriginalTransactionId       string                      `json:"originalTransactionId"`
	PriceIncreaseStatus         *int32                      `json:"priceIncreaseStatus"`
	ProductId                   string                      `json:"productId"`
	RecentSubscriptionStartDate int64                       `json:"recentSubscriptionStartDate"`
	RenewalDate                 int64                       `json:"renewalDate"`
	SignedDate                  int64                       `json:"signedDate"`
	RenewalPrice                int64                       `json:"renewalPrice,omitempty"`
	Currency                    string                      `json:"currency,omitempty"`
	OfferDiscountType           OfferDiscountType           `json:"offerDiscountType,omitempty"`
	EligibleWinBackOfferIds     []string                    `json:"eligibleWinBackOfferIds,omitempty"`
	AdvancedCommerceInfo        AdvancedCommerceRenewalInfo `json:"advancedCommerceInfo,omitempty"`
	CommitmentInfo              RenewalCommitmentInfo       `json:"commitmentInfo,omitempty"`
	RenewalBillingPlanType      BillingPlanType             `json:"renewalBillingPlanType,omitempty"`
}

func (J JWSRenewalInfoDecodedPayload) Valid() error {
	return nil
}

func (J JWSRenewalInfoDecodedPayload) GetAudience() (jwt.ClaimStrings, error) {
	return nil, nil
}

func (J JWSRenewalInfoDecodedPayload) GetExpirationTime() (*jwt.NumericDate, error) {
	return nil, nil
}

func (J JWSRenewalInfoDecodedPayload) GetIssuedAt() (*jwt.NumericDate, error) {
	return nil, nil
}

func (J JWSRenewalInfoDecodedPayload) GetIssuer() (string, error) {
	return "", nil
}

func (J JWSRenewalInfoDecodedPayload) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil
}

func (J JWSRenewalInfoDecodedPayload) GetSubject() (string, error) {
	return "", nil
}

// JWSDecodedHeader https://developer.apple.com/documentation/appstoreserverapi/jwsdecodedheader
type JWSDecodedHeader struct {
	Alg string   `json:"alg,omitempty"`
	Kid string   `json:"kid,omitempty"`
	X5C []string `json:"x5c,omitempty"`
}

// TransactionReason indicates the cause of a purchase transaction,
// https://developer.apple.com/documentation/appstoreservernotifications/transactionreason
type TransactionReason string

const (
	TransactionReasonPurchase = "PURCHASE"
	TransactionReasonRenewal  = "RENEWAL"
)

// IAPType https://developer.apple.com/documentation/appstoreserverapi/type
type IAPType string

const (
	AutoRenewable IAPType = "Auto-Renewable Subscription"
	NonConsumable IAPType = "Non-Consumable"
	Consumable    IAPType = "Consumable"
	NonRenewable  IAPType = "Non-Renewing Subscription"
)

type OfferDiscountType string

const (
	OfferDiscountTypeFreeTrial  OfferDiscountType = "FREE_TRIAL"
	OfferDiscountTypePayAsYouGo OfferDiscountType = "PAY_AS_YOU_GO"
	OfferDiscountTypePayUpFront OfferDiscountType = "PAY_UP_FRONT"
)

type RevocationType string

const (
	REFUND_FULL     RevocationType = "REFUND_FULL"
	REFUND_PRORATED RevocationType = "REFUND_PRORATED"
	FAMILY_REVOKE   RevocationType = "FAMILY_REVOKE"
)

type AdvancedCommerceRefundReason string

const (
	AdvancedCommerceRefundUNINTENDED_PURCHASE       AdvancedCommerceRefundReason = "UNINTENDED_PURCHASE"
	AdvancedCommerceRefundFULFILLMENT_ISSUE         AdvancedCommerceRefundReason = "FULFILLMENT_ISSUE"
	AdvancedCommerceRefundUNSATISFIED_WITH_PURCHASE AdvancedCommerceRefundReason = "UNSATISFIED_WITH_PURCHASE"
	AdvancedCommerceRefundLEGAL                     AdvancedCommerceRefundReason = "LEGAL"
	AdvancedCommerceRefundOTHER                     AdvancedCommerceRefundReason = "OTHER"
	AdvancedCommerceRefundMODIFY_ITEMS_REFUND       AdvancedCommerceRefundReason = "MODIFY_ITEMS_REFUND"
	AdvancedCommerceRefundSIMULATE_REFUND_DECLINE   AdvancedCommerceRefundReason = "SIMULATE_REFUND_DECLINE"
)

type AdvancedCommerceRefundType string

const (
	AdvancedCommerceRefundTypeFULL     AdvancedCommerceRefundType = "FULL"
	AdvancedCommerceRefundTypePRORATED AdvancedCommerceRefundType = "PRORATED"
	AdvancedCommerceRefundTypeCUSTOM   AdvancedCommerceRefundType = "CUSTOM"
)

type AdvancedCommerceRefund struct {
	RefundAmount int64                        `json:"refundAmount"`
	RefundDate   int64                        `json:"refundDate"`
	RefundReason AdvancedCommerceRefundReason `json:"refundReason"`
	RefundType   AdvancedCommerceRefundType   `json:"refundType"`
}

type AdvancedCommerceTransactionItem struct {
	SKU            string                   `json:"SKU,omitempty"`
	Description    string                   `json:"description,omitempty"`
	DisplayName    string                   `json:"displayName,omitempty"`
	Offer          AdvancedCommerceOffer    `json:"offer,omitempty"`
	Price          int64                    `json:"price,omitempty"`
	Refunds        []AdvancedCommerceRefund `json:"refunds,omitempty"`
	RevocationDate int64                    `json:"revocationDate,omitempty"`
}

type AdvancedCommercePeriod string

const (
	AdvancedCommercePeriodP1W AdvancedCommercePeriod = "P1W"
	AdvancedCommercePeriodP1M AdvancedCommercePeriod = "P1M"
	AdvancedCommercePeriodP2M AdvancedCommercePeriod = "P2M"
	AdvancedCommercePeriodP3M AdvancedCommercePeriod = "P3M"
	AdvancedCommercePeriodP6M AdvancedCommercePeriod = "P6M"
	AdvancedCommercePeriodP1Y AdvancedCommercePeriod = "P1Y"
)

// advancedCommerceTransactionInfo https://developer.apple.com/documentation/appstoreserverapi/advancedcommercetransactioninfo
type AdvancedCommerceTransactionInfo struct {
	Descriptors        AdvancedCommerceDescriptors       `json:"descriptors"`
	EstimatedTax       int64                             `json:"estimatedTax"`
	Items              []AdvancedCommerceTransactionItem `json:"items"`
	Period             AdvancedCommercePeriod            `json:"period"`
	RequestReferenceId string                            `json:"requestReferenceId,omitempty"`
	TaxCode            string                            `json:"taxCode,omitempty"`
	TaxExclusivePrice  int64                             `json:"taxExclusivePrice,omitempty"`
	TaxRate            string                            `json:"taxRate,omitempty"`
}

type BillingPlanType string

const (
	BillingPlanTypeBILLED_UPFRONT BillingPlanType = "BILLED_UPFRONT"
	BillingPlanTypeMONTHLY        BillingPlanType = "MONTHLY"
)

// TransactionCommitmentInfo https://developer.apple.com/documentation/appstoreserverapi/transactioncommitmentinfo
type TransactionCommitmentInfo struct {
	BillingPeriodNumber   int32 `json:"billingPeriodNumber"`
	CommitmentExpiresDate int64 `json:"commitmentExpiresDate"`
	CommitmentPrice       int64 `json:"commitmentPrice"`
	TotalBillingPeriods   int32 `json:"totalBillingPeriods"`
}

// JWSTransaction https://developer.apple.com/documentation/appstoreserverapi/jwstransaction
type JWSTransaction struct {
	AppTransactionId            string                          `json:"appTransactionId,omitempty"`
	TransactionID               string                          `json:"transactionId,omitempty"`
	OriginalTransactionId       string                          `json:"originalTransactionId,omitempty"`
	WebOrderLineItemId          string                          `json:"webOrderLineItemId,omitempty"`
	BundleID                    string                          `json:"bundleId,omitempty"`
	ProductID                   string                          `json:"productId,omitempty"`
	SubscriptionGroupIdentifier string                          `json:"subscriptionGroupIdentifier,omitempty"`
	PurchaseDate                int64                           `json:"purchaseDate,omitempty"`
	OriginalPurchaseDate        int64                           `json:"originalPurchaseDate,omitempty"`
	ExpiresDate                 int64                           `json:"expiresDate,omitempty"`
	Quantity                    int32                           `json:"quantity,omitempty"`
	Type                        IAPType                         `json:"type,omitempty"`
	AppAccountToken             string                          `json:"appAccountToken,omitempty"`
	InAppOwnershipType          string                          `json:"inAppOwnershipType,omitempty"`
	SignedDate                  int64                           `json:"signedDate,omitempty"`
	OfferType                   int32                           `json:"offerType,omitempty"`
	OfferPeriod                 string                          `json:"offerPeriod,omitempty"`
	OfferIdentifier             string                          `json:"offerIdentifier,omitempty"`
	RevocationDate              int64                           `json:"revocationDate,omitempty"`
	RevocationReason            *int32                          `json:"revocationReason,omitempty"`
	RevocationType              RevocationType                  `json:"revocationType,omitempty"`
	RevocationPercentage        int32                           `json:"revocationPercentage,omitempty"`
	IsUpgraded                  bool                            `json:"isUpgraded,omitempty"`
	Storefront                  string                          `json:"storefront,omitempty"`
	StorefrontId                string                          `json:"storefrontId,omitempty"`
	TransactionReason           TransactionReason               `json:"transactionReason,omitempty"`
	Environment                 Environment                     `json:"environment,omitempty"`
	Price                       int64                           `json:"price,omitempty"`
	Currency                    string                          `json:"currency,omitempty"`
	OfferDiscountType           OfferDiscountType               `json:"offerDiscountType,omitempty"`
	AdvancedCommerceInfo        AdvancedCommerceTransactionInfo `json:"advancedCommerceInfo,omitempty"`
	BillingPlanType             BillingPlanType                 `json:"billingPlanType,omitempty"`
	CommitmentInfo              TransactionCommitmentInfo       `json:"commitmentInfo,omitempty"`
}

func (J JWSTransaction) Valid() error {
	return nil
}
func (J JWSTransaction) GetAudience() (jwt.ClaimStrings, error) {
	return nil, nil
}

func (J JWSTransaction) GetExpirationTime() (*jwt.NumericDate, error) {
	return nil, nil
}

func (J JWSTransaction) GetIssuedAt() (*jwt.NumericDate, error) {
	return nil, nil
}

func (J JWSTransaction) GetIssuer() (string, error) {
	return "", nil
}

func (J JWSTransaction) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil
}

func (J JWSTransaction) GetSubject() (string, error) {
	return "", nil
}

// https://developer.apple.com/documentation/appstoreserverapi/extendreasoncode
type ExtendReasonCode int

const (
	UndeclaredExtendReasonCode = iota
	CustomerSatisfaction
	OtherReasons
	ServiceIssueOrOutage
)

// ExtendRenewalDateRequest https://developer.apple.com/documentation/appstoreserverapi/extendrenewaldaterequest
type ExtendRenewalDateRequest struct {
	ExtendByDays      int32            `json:"extendByDays"`
	ExtendReasonCode  ExtendReasonCode `json:"extendReasonCode"`
	RequestIdentifier string           `json:"requestIdentifier"`
}

// MassExtendRenewalDateStatusResponse https://developer.apple.com/documentation/appstoreserverapi/massextendrenewaldatestatusresponse
type MassExtendRenewalDateStatusResponse struct {
	RequestIdentifier string `json:"requestIdentifier"`
	Complete          bool   `json:"complete"`
	CompleteDate      int64  `json:"completeDate,omitempty"`
	FailedCount       int64  `json:"failedCount,omitempty"`
	SucceededCount    int64  `json:"succeededCount,omitempty"`
}

// NotificationHistoryRequest https://developer.apple.com/documentation/appstoreserverapi/notificationhistoryrequest
type NotificationHistoryRequest struct {
	StartDate             int64              `json:"startDate"`
	EndDate               int64              `json:"endDate"`
	OriginalTransactionId string             `json:"originalTransactionId,omitempty"`
	NotificationType      NotificationTypeV2 `json:"notificationType,omitempty"`
	NotificationSubtype   SubtypeV2          `json:"notificationSubtype,omitempty"`
	OnlyFailures          bool               `json:"onlyFailures,omitempty"`
	TransactionId         string             `json:"transactionId,omitempty"`
}

type NotificationTypeV2 string

// list of notificationType
// https://developer.apple.com/documentation/appstoreservernotifications/notificationtype
const (
	NotificationTypeV2ConsumptionRequest     NotificationTypeV2 = "CONSUMPTION_REQUEST"
	NotificationTypeV2DidChangeRenewalPref   NotificationTypeV2 = "DID_CHANGE_RENEWAL_PREF"
	NotificationTypeV2DidChangeRenewalStatus NotificationTypeV2 = "DID_CHANGE_RENEWAL_STATUS"
	NotificationTypeV2DidFailToRenew         NotificationTypeV2 = "DID_FAIL_TO_RENEW"
	NotificationTypeV2DidRenew               NotificationTypeV2 = "DID_RENEW"
	NotificationTypeV2Expired                NotificationTypeV2 = "EXPIRED"
	NotificationTypeV2GracePeriodExpired     NotificationTypeV2 = "GRACE_PERIOD_EXPIRED"
	NotificationTypeV2OfferRedeemed          NotificationTypeV2 = "OFFER_REDEEMED"
	NotificationTypeV2PriceIncrease          NotificationTypeV2 = "PRICE_INCREASE"
	NotificationTypeV2Refund                 NotificationTypeV2 = "REFUND"
	NotificationTypeV2RefundDeclined         NotificationTypeV2 = "REFUND_DECLINED"
	NotificationTypeV2RenewalExtended        NotificationTypeV2 = "RENEWAL_EXTENDED"
	NotificationTypeV2Revoke                 NotificationTypeV2 = "REVOKE"
	NotificationTypeV2Subscribed             NotificationTypeV2 = "SUBSCRIBED"
	NotificationTypeV2OneTimeCharge          NotificationTypeV2 = "ONE_TIME_CHARGE"
	NotificationTypeV2RefundReversed         NotificationTypeV2 = "REFUND_REVERSED"
	NotificationTypeV2ExternalPurchaseToken  NotificationTypeV2 = "EXTERNAL_PURCHASE_TOKEN"
	NotificationTypeV2RenewalExtension       NotificationTypeV2 = "RENEWAL_EXTENSION"
	NotificationTypeV2Test                   NotificationTypeV2 = "TEST"
)

// SubtypeV2 is type
type SubtypeV2 string

// list of subtypes
// https://developer.apple.com/documentation/appstoreservernotifications/subtype
const (
	SubTypeV2InitialBuy        SubtypeV2 = "INITIAL_BUY"
	SubTypeV2Resubscribe       SubtypeV2 = "RESUBSCRIBE"
	SubTypeV2Downgrade         SubtypeV2 = "DOWNGRADE"
	SubTypeV2Upgrade           SubtypeV2 = "UPGRADE"
	SubTypeV2AutoRenewEnabled  SubtypeV2 = "AUTO_RENEW_ENABLED"
	SubTypeV2AutoRenewDisabled SubtypeV2 = "AUTO_RENEW_DISABLED"
	SubTypeV2Voluntary         SubtypeV2 = "VOLUNTARY"
	SubTypeV2BillingRetry      SubtypeV2 = "BILLING_RETRY"
	SubTypeV2PriceIncrease     SubtypeV2 = "PRICE_INCREASE"
	SubTypeV2GracePeriod       SubtypeV2 = "GRACE_PERIOD"
	SubTypeV2BillingRecovery   SubtypeV2 = "BILLING_RECOVERY"
	SubTypeV2Pending           SubtypeV2 = "PENDING"
	SubTypeV2Accepted          SubtypeV2 = "ACCEPTED"
)

// NotificationHistoryResponses https://developer.apple.com/documentation/appstoreserverapi/notificationhistoryresponse
type NotificationHistoryResponses struct {
	HasMore             bool                              `json:"hasMore"`
	PaginationToken     string                            `json:"paginationToken"`
	NotificationHistory []NotificationHistoryResponseItem `json:"notificationHistory"`
}

// NotificationHistoryResponseItem https://developer.apple.com/documentation/appstoreserverapi/notificationhistoryresponseitem
type NotificationHistoryResponseItem struct {
	SignedPayload          string                 `json:"signedPayload"`
	FirstSendAttemptResult FirstSendAttemptResult `json:"firstSendAttemptResult"`
	SendAttempts           []SendAttemptItem      `json:"sendAttempts"`
}

// SendAttemptItem https://developer.apple.com/documentation/appstoreserverapi/sendattemptitem
type SendAttemptItem struct {
	AttemptDate       int64                  `json:"attemptDate"`
	SendAttemptResult FirstSendAttemptResult `json:"sendAttemptResult"`
}

// https://developer.apple.com/documentation/appstoreserverapi/firstsendattemptresult
type FirstSendAttemptResult string

const (
	FirstSendAttemptResultSuccess            FirstSendAttemptResult = "SUCCESS"
	FirstSendAttemptResultCircularRedirect   FirstSendAttemptResult = "CIRCULAR_REDIRECT"
	FirstSendAttemptResultInvalidResponse    FirstSendAttemptResult = "INVALID_RESPONSE"
	FirstSendAttemptResultNoResponse         FirstSendAttemptResult = "NO_RESPONSE"
	FirstSendAttemptResultOther              FirstSendAttemptResult = "OTHER"
	FirstSendAttemptResultPrematureClose     FirstSendAttemptResult = "PREMATURE_CLOSE"
	FirstSendAttemptResultSocketIssue        FirstSendAttemptResult = "SOCKET_ISSUE"
	FirstSendAttemptResultTimedOut           FirstSendAttemptResult = "TIMED_OUT"
	FirstSendAttemptResultTlsIssue           FirstSendAttemptResult = "TLS_ISSUE"
	FirstSendAttemptResultUnsupportedCharset FirstSendAttemptResult = "UNSUPPORTED_CHARSET"
)

// SendTestNotificationResponse https://developer.apple.com/documentation/appstoreserverapi/sendtestnotificationresponse
type SendTestNotificationResponse struct {
	TestNotificationToken string `json:"testNotificationToken"`
}

// Notification body https://developer.apple.com/documentation/appstoreservernotifications/responsebodyv2
type NotificationV2 struct {
	SignedPayload string `json:"signedPayload"`
}

// Notification signed payload
type NotificationPayload struct {
	jwt.RegisteredClaims
	NotificationType    string           `json:"notificationType"`
	Subtype             string           `json:"subtype"`
	NotificationUUID    string           `json:"notificationUUID"`
	NotificationVersion string           `json:"notificationVersion"`
	Data                NotificationData `json:"data"`
}

// Notification Data
type NotificationData struct {
	jwt.RegisteredClaims
	AppAppleID               int    `json:"appAppleId"`
	BundleID                 string `json:"bundleId"`
	BundleVersion            string `json:"bundleVersion"`
	ConsumptionRequestReason string `json:"consumptionRequestReason"`
	Environment              string `json:"environment"`
	SignedRenewalInfo        string `json:"signedRenewalInfo"`
	SignedTransactionInfo    string `json:"signedTransactionInfo"`
	Status                   int    `json:"status"`
}

// Notification Transaction Info
type TransactionInfo struct {
	jwt.RegisteredClaims
	TransactionId               string `json:"transactionId"`
	OriginalTransactionID       string `json:"originalTransactionId"`
	WebOrderLineItemID          string `json:"webOrderLineItemId"`
	BundleID                    string `json:"bundleId"`
	ProductID                   string `json:"productId"`
	SubscriptionGroupIdentifier string `json:"subscriptionGroupIdentifier"`
	PurchaseDate                int    `json:"purchaseDate"`
	OriginalPurchaseDate        int    `json:"originalPurchaseDate"`
	ExpiresDate                 int    `json:"expiresDate"`
	Type                        string `json:"type"`
	InAppOwnershipType          string `json:"inAppOwnershipType"`
	SignedDate                  int    `json:"signedDate"`
	Environment                 string `json:"environment"`
}

// Notification Renewal Info
type RenewalInfo struct {
	jwt.RegisteredClaims
	OriginalTransactionID  string `json:"originalTransactionId"`
	ExpirationIntent       int    `json:"expirationIntent"`
	AutoRenewProductId     string `json:"autoRenewProductId"`
	ProductID              string `json:"productId"`
	AutoRenewStatus        int    `json:"autoRenewStatus"`
	IsInBillingRetryPeriod *bool  `json:"isInBillingRetryPeriod"`
	SignedDate             int    `json:"signedDate"`
	Environment            string `json:"environment"`
}

type UpdateAppAccountTokenRequest struct {
	AppAccountToken string `json:"appAccountToken"`
}

type AppTransactionInfoResponse struct {
	SignedAppTransactionInfo string `json:"signedAppTransactionInfo"`
}

type JWSAppTransactionDecodedPayload struct {
	AppAppleId                 int64       `json:"appAppleId"`
	AppTransactionId           string      `json:"appTransactionId"`
	BundleId                   string      `json:"bundleId"`
	OriginalApplicationVersion string      `json:"originalApplicationVersion"`
	OriginalPlatform           string      `json:"originalPlatform"`
	OriginalPurchaseDate       int64       `json:"originalPurchaseDate"`
	PreorderDate               int64       `json:"preorderDate,omitempty"`
	ReceiptCreationDate        int64       `json:"receiptCreationDate"`
	ReceiptType                Environment `json:"receiptType"`
}

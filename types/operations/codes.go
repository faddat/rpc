package operations

var opTypes = []OpType{
	TypeVote,
	TypeComment,
	TypeTransfer,
	TypeTransferToVesting,
	TypeWithdrawVesting,
	TypeLimitOrderCreate,
	TypeLimitOrderCancel,
	TypeFeedPublish,
	TypeConvert,
	TypeAccountCreate,
	TypeAccountUpdate,
	TypeWitnessUpdate,
	TypeAccountWitnessVote,
	TypeAccountWitnessProxy,
	TypePOW,
	TypeCustom,
	TypeReportOverProduction,
	TypeDeleteComment,
	TypeCustomJSON,
	TypeCommentOptions,
	TypeSetWithdrawVestingRoute,
	TypeLimitOrderCreate2,
	TypeChallengeAuthority,
	TypeProveAuthority,
	TypeRequestAccountRecoverty,
	TypeRecoverAccount,
	TypeChangeRecoveryAccount,
	TypeEscrowTransfer,
	TypeEscrowDispute,
	TypeEscrowRelease,
	TypePOW2,
}

var opCodes map[OpType]uint64

func init() {
	opCodes = make(map[OpType]uint64, len(opTypes))
	for i, opType := range opTypes {
		opCodes[opType] = i
	}
}

// Code returns the operation code for the given operation type.
func Code(kind OpType) uint64 {
	return opCodes[kind]
}

package staking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Wrapper struct
type Hooks struct {
	k stakingkeeper.Keeper
}

var _ types.StakingHooks = Hooks{}

// Create new distribution hooks
func NewHooks(k stakingkeeper.Keeper) Hooks {
	return Hooks{k}
}

func (h Hooks) AfterValidatorBonded(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) error {
	return nil
}

func (h Hooks) AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, _ sdk.ValAddress) error {
	return nil
}

func (h Hooks) AfterValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) error {
	return nil
}

func (h Hooks) AfterValidatorBeginUnbonding(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress) error {
	return nil
}

func (h Hooks) BeforeValidatorModified(_ sdk.Context, _ sdk.ValAddress) error {
	return nil
}

func (h Hooks) BeforeDelegationCreated(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) error {
	return nil
}

func (h Hooks) BeforeDelegationSharesModified(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) error {
	return nil
}

func (h Hooks) BeforeDelegationRemoved(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) error {
	return nil
}

func (h Hooks) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	validator, found := h.k.GetValidator(ctx, valAddr)

	if !found {
		return sdkerrors.Wrap(sdkerrors.ErrNotFound, "validator does not exist")
	}

	// Get the last Total Power of the validator set
	lastPower := h.k.GetLastTotalPower(ctx)

	// Get the power of the current validator power
	validatorLastPower := sdk.TokensToConsensusPower(validator.Tokens, h.k.PowerReduction(ctx))

	// Compute what the new Validator voting power would be in relation to the new total power
	// validatorIncreasedDelegationPercent := float32(validatorNewPower) / float32(newTotalPower)
	validatorIncreasedDelegationPercent := sdk.NewDec(validatorLastPower).QuoInt(lastPower)

	// If Delegations are allowed, and the Delegation would have increased the Validator to over 20% of the staking power, do not allow the Delegation to proceed
	if validatorIncreasedDelegationPercent.GT(sdk.NewDecWithPrec(20, 2)) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "validator power is over the allowed limit")
	}

	return nil
}

func (h Hooks) BeforeValidatorSlashed(_ sdk.Context, _ sdk.ValAddress, _ sdk.Dec) error {
	return nil
}

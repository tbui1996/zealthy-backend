package request

type PatchFlagRequest struct {
	Name      *string
	IsEnabled *bool
	FlagId    int
}

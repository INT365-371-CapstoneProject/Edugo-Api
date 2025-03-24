package request

type CreateBookmarkRequest struct {
	Announce_ID uint `form:"annouce_id" validate:"required"`
	// Account_ID    uint           `form:"account_id" validate:"required"`
}

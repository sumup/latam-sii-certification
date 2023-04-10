package batches

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/sumup/sii-certification/internal/entities"
	"github.com/sumup/sii-certification/internal/utils"
)

type Handler struct {
	domain IDomain
}

func NewHTTPHandler(domain IDomain) *Handler {
	return &Handler{
		domain,
	}
}

type URLInput interface {
	postURLInput | batchesURLInput
}

type QueryParamsInput interface {
	ByDayParams
}

func validateURLInput[V URLInput](ctx *gin.Context, input V) (V, error) {
	err := ctx.ShouldBindUri(&input)
	if err != nil {
		return input, &InvalidInputError{errDetail: err.Error(), Err: err}
	}

	return input, nil
}

func validateQueryParamsInput[V QueryParamsInput](ctx *gin.Context, input V) (V, error) {
	err := ctx.ShouldBindQuery(&input)
	if err != nil {
		return input, &InvalidInputError{errDetail: err.Error(), Err: err}
	}

	return input, nil
}

type postURLInput struct {
	CountryCode string `uri:"country-code" binding:"required,oneof=CL"`
	Year        string `uri:"year" binding:"required,oneof=2021 2022 2023 2024 2025 2026 2027 2028"`
	Month       string `uri:"month" binding:"required,oneof=01 02 03 04 05 06 07 08 09 10 11 12"`
	Day         string `uri:"day" binding:"required,oneof=01 02 03 04 05 06 07 08 09 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31"`
}

func (h *Handler) Post(ctx *gin.Context) {
	urlInput, err := validateURLInput(ctx, postURLInput{})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	go func() {
		detachedCtx := context.Background()

		span := trace.SpanFromContext(ctx)
		defer span.End()

		span.SetAttributes(attribute.String("day", urlInput.Year+"-"+urlInput.Month+"-"+urlInput.Day))

		startDate, endDate := utils.LimitTimesForDay(urlInput.CountryCode, urlInput.Year, urlInput.Month, urlInput.Day)

		batches, err := h.domain.Generate(detachedCtx, startDate, endDate)

		if errors.Is(err, errBatchesServiceNoTransactions) {
			log.Err(err).Msg("There are no transactions to generate batches")

			return
		}

		if errors.Is(err, errBatchesServiceSanityCheckFailed) {
			log.Err(err).Msg("Sanity check failed")

			return
		}

		if err != nil {
			log.Err(err).Msg("Unknown error in generating batches")

			return
		}

		log.Info().
			Int("batches", len(batches)).
			Msg("created batches successfully")
	}()

	ctx.JSON(http.StatusAccepted, gin.H{
		"status": gin.H{
			"code":        http.StatusAccepted,
			"description": "Accepted",
		},
	})
}

type batchesURLInput struct {
	CountryCode string `uri:"country-code" binding:"required,oneof=CL"`
}

func (h *Handler) Send(ctx *gin.Context) {
	_, err := validateURLInput(ctx, batchesURLInput{})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	go func() {
		detachedCtx := context.Background()

		batches, err := h.domain.Send(detachedCtx)

		if errors.Is(err, errBatchesServiceNoBatchesToSend) {
			log.Info().Msg("There are no batches to send")

			return
		}

		if errors.Is(err, errBatchesServiceCantGetBatches) {
			log.Err(err).Msg("Can't get batches to send")

			return
		}

		if err != nil {
			log.Err(err).Msg("Unknown error in sending batches")

			return
		}

		log.Info().
			Int("batches", len(batches)).
			Msg("batches sent successfully")
	}()

	ctx.JSON(http.StatusAccepted, gin.H{
		"status": gin.H{
			"code":        http.StatusAccepted,
			"description": "Accepted",
		},
	})
}

func (h *Handler) GetTestAuth(ctx *gin.Context) {
	token, tokenResponse, seed, seedResponse, err := h.domain.Auth(context.Background())

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":          err.Error(),
			"token_response": tokenResponse,
			"seed_response":  seedResponse,
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token":          token,
		"token_response": tokenResponse,
		"seed":           seed,
		"seed_response":  seedResponse,
	})
}

type ByDayParams struct {
	Day      string `form:"day" binding:"required,datetime=2006-01-02"`
	Page     int    `form:"page" binding:"required,numeric,gte=0"`
	PageSize int    `form:"pageSize" binding:"required,numeric,gte=0,lte=100"`
}

func (h *Handler) GetBatchesByDay(ctx *gin.Context) {
	_, err := validateURLInput(ctx, batchesURLInput{})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	queryParams, err := validateQueryParamsInput(ctx, ByDayParams{})

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	batches, batchesCount, err := h.domain.GetBatchesByDay(ctx, queryParams.Day, queryParams.Page, queryParams.PageSize)

	if err != nil {
		log.Err(err).Msg("Unknown error getting batches by day")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	response := buildBatchesByDayResponse(batches, queryParams, batchesCount)

	ctx.JSON(http.StatusOK, response)
}

func buildBatchesByDayResponse(batches []entities.Batch, queryParams ByDayParams, batchesCount int64) gin.H {
	totalPages := batchesCount / int64(queryParams.PageSize)
	remainder := batchesCount % int64(queryParams.PageSize)

	if remainder > 0 {
		totalPages++
	}

	batchesResponse := make([]gin.H, len(batches))

	for i, batch := range batches {
		batchesResponse[i] = gin.H{
			"id":                batch.ID,
			"merchant_code":     batch.MerchantCode,
			"day":               batch.Day,
			"amount":            batch.Amount,
			"has_taxes":         batch.HasTaxes,
			"is_cnp":            batch.IsCNP,
			"status":            batch.Status,
			"last_response":     batch.LastResponse,
			"external_track_id": batch.ExternalTrackID,
			"n_transactions":    len(batch.Transactions),
		}
	}

	response := gin.H{
		"day":         queryParams.Day,
		"page":        queryParams.Page,
		"page_size":   queryParams.PageSize,
		"total_pages": totalPages,
		"batches":     batchesResponse,
	}

	return response
}

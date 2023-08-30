package del

import (
	resp "avito-internship/internal/lib/api/response"
	"avito-internship/internal/lib/logger/slogger"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
)

type Request struct {
	SegmentName string `json:"segment_name" validate:"required"`
}

type Response struct {
	resp.Response
	SegmentID int64 `json:"id,omitempty"`
}

type DeleteSegment interface {
	DeleteSegment(segmentName string) (int64, error)
}

func DelSeg(log *slog.Logger, segmentDelete DeleteSegment) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.segments.delete.DelSeg"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slogger.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", slogger.Err(err))

			render.JSON(w, r, resp.Error("invalid request"))
			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		segmentID, err := segmentDelete.DeleteSegment(req.SegmentName)
		if err != nil {
			log.Error("failed to delete segment", slogger.Err(err))

			render.JSON(w, r, resp.Error("failed to delete segment"))

			return
		}

		log.Info("segment deleted", slog.Int64("id", segmentID))

		render.JSON(w, r, Response{
			Response:  resp.OK(),
			SegmentID: segmentID,
		})
	}
}

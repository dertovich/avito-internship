package save

import (
	resp "avito-internship/internal/lib/api/response"
	"avito-internship/internal/lib/logger/slogger"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
)

type Request struct {
	SegmentName string `json:"name" validate:"required,name"`
}

type Response struct {
	resp.Response
	SegmentName string `json:"name,omitempty"`
}

type SegmentCreator interface {
	CreateSegment(segment string) (int64, error)
}

func New(log *slog.Logger, segmentCreator SegmentCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.segments.save.New"

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

		id, err := segmentCreator.CreateSegment(req.SegmentName)
		if err != nil {
			log.Info("segment name already exists", slog.String("segment", req.SegmentName))

			render.JSON(w, r, resp.Error("segment already exists"))

			return
		}
		if err != nil {
			log.Error("failed to add segment", slogger.Err(err))

			render.JSON(w, r, resp.Error("failed to add segment"))

			return
		}

		log.Info("segment added", slog.Int64("id", id))

		render.JSON(w, r, Response{
			Response: resp.OK(),
		})
	}
}

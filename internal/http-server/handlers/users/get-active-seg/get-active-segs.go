package getactiveseg

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
	UserID int64 `json:"id" validate:"required"`
}

type Response struct {
	resp.Response
	Segments []string `json:"segments"`
}

type UserSegments interface {
	ShowActiveSegmentUser(user_id int64) ([]string, error)
}

func GetActiveSegmentsForUser(log *slog.Logger, userSegments UserSegments) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.users.get.active.segments"

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

		segments, err := userSegments.ShowActiveSegmentUser(req.UserID)
		if err != nil {
			log.Error("failed to get active segments for user", slogger.Err(err))

			render.JSON(w, r, resp.Error("failed to get active segments for user"))

			return
		}

		log.Info("active segments for user retrieved", slog.Int64("user_id", req.UserID), slog.Any("segments", segments))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Segments: segments,
		})
	}
}

package saveseguser

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
	UserID   int64    `json:"userId" validate:"required"`
	Segments []string `json:"segments" validate:"required"`
}

type Response struct {
	resp.Response
	UserId int `json:"userId"`
}

type AddUserToSegment interface {
	AddUserToSegment(userID int64, segmentSlugs []string) error
}

func AddUserToSegments(log *slog.Logger, addUserToSegment AddUserToSegment) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.users.save.AddUserToSegments"

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

		err = addUserToSegment.AddUserToSegment(req.UserID, req.Segments)
		if err != nil {
			log.Error("failed to add segments to user", slogger.Err(err))

			render.JSON(w, r, resp.Error("failed to add segments to user"))

			return
		}

		log.Info("segments added to user", slog.Int64("user_id", req.UserID))

		render.JSON(w, r, Response{
			Response: resp.OK(),
		})
	}
}

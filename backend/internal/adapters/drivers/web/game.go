package web

import (
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	game "github.com/taldoflemis/brain.test/internal/core/domain/game_aggregate"
	"github.com/taldoflemis/brain.test/internal/core/services"
	"github.com/taldoflemis/brain.test/internal/ports"
)

type QuestionKind string

const (
	QuizQuestionKind      = "quiz"
	TrueFalseQuestionKind = "true_false"
)

var (
	ErrUnknownQuestionKind = errors.New("Unknown question kind")
)

type Question interface {
	LoadFromMap(data map[string]any) error
	ToQuestion() game.Question
}

// CreateQuizQuestionRequest
//
//	@Description	Request to create a Quiz Question
type CreateQuizQuestionRequest struct {
	Title        string `json:"title"        validate:"required"`
	Points       int    `json:"points"       validate:"required"`
	TimeLimit    int    `json:"time_limit"   validate:"required"`
	Alternatives []struct {
		Data      string `json:"data" validate:"required"`
		IsCorrect bool   `json:"is_correct"`
	} `json:"alternatives" validate:"required,dive,required"`
}

func (r *CreateQuizQuestionRequest) LoadFromMap(data map[string]any) error {
	jsonString, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonString, &r)
	return err
}

func (r *CreateQuizQuestionRequest) ToQuestion() game.Question {
	alternatives := make([]game.Alternative, len(r.Alternatives))
	for i, a := range r.Alternatives {
		alternatives[i] = game.Alternative{
			Data:      a.Data,
			IsCorrect: a.IsCorrect,
		}
	}

	return &game.QuizQuestion{
		Title:        r.Title,
		Points:       game.Points(r.Points),
		TimeLimit:    game.TimeLimit(r.TimeLimit),
		Alternatives: alternatives,
	}
}

// CreateTrueFalseQuestionRequest
//
//	@Description	Request to create a True False Question
type CreateTrueFalseQuestionRequest struct {
	Title            string `json:"title"             validate:"required"`
	Points           int    `json:"points"            validate:"required"`
	TimeLimit        int    `json:"time_limit"        validate:"required"`
	TrueAlternative  string `json:"true_alternative"  validate:"required"`
	FalseAlternative string `json:"false_alternative" validate:"required"`
}

func (r *CreateTrueFalseQuestionRequest) LoadFromMap(data map[string]any) error {
	jsonString, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonString, &r)
	return err
}

func (r *CreateTrueFalseQuestionRequest) ToQuestion() game.Question {
	return &game.TrueFalseQuestion{
		Title:            r.Title,
		Points:           game.Points(r.Points),
		TimeLimit:        game.TimeLimit(r.TimeLimit),
		TrueAlternative:  r.TrueAlternative,
		FalseAlternative: r.FalseAlternative,
	}
}

// CreateQuestionRequest
//
//	@Description	a way to create questions dynamically
type CreateQuestionRequest struct {
	// kind
	Kind QuestionKind `json:"kind" validate:"required"`
	// data
	Data map[string]any `json:"data" validate:"required"`
}

// CreateGameRequest
//
//	@Description	Request to create a Game
type CreateGameRequest struct {
	// the title of a game
	Title string `json:"title"       validate:"required"`
	// the description of a game
	Description string `json:"description" validate:"omitempty"`
	// questions of the game
	Questions []CreateQuestionRequest `json:"questions"   validate:"required,dive,required"`
}

type gameHandler struct {
	jwtMiddleware     fiber.Handler
	validationService *services.ValidationService
	gameService       *services.GameService
}

func NewGameHandler(
	jwtMiddleware fiber.Handler,
	validationService *services.ValidationService,
	gameService *services.GameService,
) *gameHandler {
	return &gameHandler{
		jwtMiddleware:     jwtMiddleware,
		validationService: validationService,
		gameService:       gameService,
	}
}

func (h *gameHandler) RegisterRoutes(router fiber.Router) {
	gameApi := router.Group("/game")

	gameApi.Use(h.jwtMiddleware)
	gameApi.Post("/", h.CreateGame)
	gameApi.Get("/", h.GetGamesByUserId)
	gameApi.Get("/:gameId", h.GetGamesById)
}

//	GetGameByUserId godoc
//
// @Summary	Get games by user id
// @Tags		Game
// @Accept		json
// @Success	200
// @Failure	400		{string}	string
// @Failure	401		{string}	string
// @Failure	422		{object}	ValidationErrorResponse
// @Router		/game/	[get]
func (h *gameHandler) GetGamesByUserId(c *fiber.Ctx) error {
	userId := extractTokenFromContext(c)

	games, err := h.gameService.GetGamesByUserId(c.Context(), userId)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"games": games,
	})
}

//	GetGameById godoc
//
// @Summary	Get games by id
// @Tags		Game
// @Accept		json
// @Success	200
// @Failure	401			{string}	string
// @Failure	404			{string}	string
// @Router		/game/:id	[get]
func (h *gameHandler) GetGamesById(c *fiber.Ctx) error {
	gameId, err := uuid.Parse(c.Params("gameId"))

	if err != nil {
		return err
	}

	game, err := h.gameService.GetGameById(c.Context(), gameId)

	if err != nil {
		if errors.Is(err, ports.ErrGameNotFound) {
			return c.Status(fiber.StatusNotFound).SendString(err.Error())
		}
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"game": game,
	})
}

// CreateGame godoc
//
//	@Summary	Create a Game
//	@Tags		Game
//	@Accept		json
//	@Param		req	body	CreateGameRequest	true	"Create Game Request"
//	@Success	201
//	@Failure	400		{string}	string
//	@Failure	401		{string}	string
//	@Failure	422		{object}	ValidationErrorResponse
//	@Router		/game/	[post]
func (h *gameHandler) CreateGame(c *fiber.Ctx) error {
	userId := extractTokenFromContext(c)

	req := new(CreateGameRequest)
	err := c.BodyParser(req)
	if err != nil {
		return err
	}

	err = h.validationService.Validate(req)
	if err != nil {
		return err
	}

	questions, err := h.parseQuestions(req.Questions)
	if err != nil {
		if errors.Is(err, ErrUnknownQuestionKind) {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		return err
	}

	err = h.gameService.CreateNewGame(c.Context(), userId, &game.Game{
		Title:       req.Title,
		Description: req.Description,
		Questions:   questions,
	})
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *gameHandler) parseQuestions(qs []CreateQuestionRequest) ([]game.Question, error) {
	questions := make([]game.Question, 0)

	for _, q := range qs {
		var candidate game.Question

		var temp Question

		switch q.Kind {
		case QuizQuestionKind:
			temp = new(CreateQuizQuestionRequest)
		case TrueFalseQuestionKind:
			temp = new(CreateTrueFalseQuestionRequest)
		default:
			return nil, ErrUnknownQuestionKind
		}

		err := temp.LoadFromMap(q.Data)
		if err != nil {
			return nil, ErrUnknownQuestionKind
		}
		err = h.validationService.Validate(temp)
		if err != nil {
			return nil, err
		}
		candidate = temp.ToQuestion()

		questions = append(questions, candidate)
	}

	return questions, nil
}

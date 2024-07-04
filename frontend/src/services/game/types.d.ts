type Points = number;
type TimeLimit = number;

type TrueFalseQuestion = {
  id: string;
  title: string;
  points: Points;
  time_limit: TimeLimit;
  true_alternative: string;
  false_alternative: string;
};

type Alternative = {
  data: string;
  is_correct: boolean;
};

type QuizQuestion = {
  id: string;
  title: string;
  points: Points;
  time_limit: TimeLimit;
  alternatives: Alternative[];
};

type Question = TrueFalseQuestion | QuizQuestion;

export type Game = {
  id: string;
  title: string;
  description: string;
  owner_id: string;
  questions: Question[];
};

export type GetGamesByUserResponse = {
  games: Game[];
};

export type GetGameByIdResponse = {
  game: Game;
};

export type QuestionKindRequest = "quiz" | "true_false";

export type TrueFalseQuestionRequest = Omit<TrueFalseQuestion, "id">;

type AlternativeRequest = Alternative;

type QuizQuestionRequest = Omit<QuizQuestion, "id">;

export type QuestionRequest = TrueFalseQuestionRequest | QuizQuestionRequest;

export type CreateQuestionRequest = {
  kind: QuestionKindRequest;
  data: QuestionRequest;
};

export type CreateGameRequest = {
  title: string;
  description?: string;
  questions: CreateQuestionRequest[];
};

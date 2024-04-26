CREATE TABLE games(
	id UUID PRIMARY KEY,
	owner_id TEXT NOT NULL,
	title TEXT NOT NULL,
	description TEXT NOT NULL
);

CREATE TABLE questions(
	id UUID PRIMARY KEY,
	game_id UUID NOT NULL,
	"order" INT NOT NULL,
	title TEXT NOT NULL,
	kind TEXT NOT NULL,
	time_limit INT NOT NULL,
	points INT NOT NULL,
	CONSTRAINT fk_game_id FOREIGN KEY(game_id) REFERENCES games(id)
);

CREATE INDEX questions_game_id ON questions(game_id);

CREATE TABLE true_false_questions(
	question_id UUID PRIMARY KEY REFERENCES questions(id) DEFERRABLE INITIALLY DEFERRED,
	true_alternative TEXT NOT NULL,
	false_alternative TEXT NOT NULL
);

CREATE TABLE quiz_questions (
	id UUID PRIMARY KEY,
	question_id UUID NOT NULL,
	"order" INT NOT NULL,
	data TEXT NOT NULL,
	correct BOOLEAN NOT NULL,
	CONSTRAINT fk_question_id FOREIGN KEY(question_id) REFERENCES questions(id) DEFERRABLE INITIALLY DEFERRED
);

CREATE INDEX quiz_questions_question_id ON quiz_questions(question_id);

# Notes

## Create a new game

### Types

```ts
type CreateGameRequest = {
  title: string;
  description: string;
  questions: Question[];
};

type Question = {
  kind: "quiz" | "true_false";
  data: QuizQuestion | TrueFalseQuestion;
};

type QuizQuestion = {
  title: string;
  points: number;
  time_limit: number;
  alternatives: {
    data: string;
    is_correct: boolean;
  };
};

type TrueFalseQuestion = {
  title: string;
  points: number;
  time_limit: number;
  true_alternative: string;
  false_alternative: string;
};
```

### Endpoint

- POST /game/

"use server";

import gameService from "@/services/game/service";
import { CreateGameFormSchema } from "./create-quiz-form";
import {
  CreateGameRequest,
  CreateQuestionRequest,
} from "@/services/game/types";
import { QuestionSchema } from "./questions";
import { auth } from "@/auth";
import { redirect } from "next/navigation";

const processQuestionSchema = (
  question: QuestionSchema,
): CreateQuestionRequest => {
  let request: CreateQuestionRequest;
  switch (question.kind) {
    case "quiz": {
      request = {
        kind: "quiz",
        data: {
          title: question.title,
          points: 1,
          time_limit: parseInt(question.timeout),
          alternatives: question.alternatives,
        },
      };
      break;
    }
    case "true_false": {
      request = {
        kind: "true_false",
        data: {
          title: question.title,
          points: 1,
          time_limit: parseInt(question.timeout),
          true_alternative: question.trueAlternative,
          false_alternative: question.falseAlternative,
        },
      };
      break;
    }
  }
  return request;
};

const createGame = async (values: CreateGameFormSchema) => {
  const token = await auth();

  if (!token) {
    redirect("/sign-in");
  }

  const access_token = token.access_token;

  const request: CreateGameRequest = {
    title: values.general.title,
    questions: [],
  };

  if (values.general.description && values.general.description != "") {
    request.description = values.general.description;
  }

  for (let i = 0; i < values.questions.length; i++) {
    let question = processQuestionSchema(values.questions[i]);
    request.questions.push(question);
  }

  await gameService.createGame(request, access_token);
  redirect("/dashboard");
};

export { createGame };

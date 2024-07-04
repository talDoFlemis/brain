import { type QuestionSchema } from "./questions";
import TrueFalseQuestion, {
  defaultTrueFalseQuestion,
} from "./true-false-question";
import QuizQuestion, { defaultQuizQuestion } from "./quiz-question";
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import {
  UseFieldArrayRemove,
  UseFieldArrayUpdate,
  UseFormReturn,
  useWatch,
} from "react-hook-form";
import { CreateGameFormSchema } from "./create-quiz-form";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { IoMdTrash } from "react-icons/io";

type QuestionHeaderProps = {
  form: UseFormReturn<CreateGameFormSchema>;
  index: number;
  question: QuestionSchema;
  remove: UseFieldArrayRemove;
  update: UseFieldArrayUpdate<CreateGameFormSchema, "questions">;
};

function QuestionHeader({
  form,
  index,
  remove,
  update,
  question,
}: QuestionHeaderProps) {
  const updateKind = (value: string) => {
    // TODO: optimize to not update if match same value
    switch (value) {
      case "true_false": {
        update(index, {
          ...defaultTrueFalseQuestion,
          kind: "true_false",
          timeout: question.timeout,
          title: question.title,
        });
        break;
      }
      case "quiz": {
        update(index, {
          ...defaultQuizQuestion,
          kind: "quiz",
          timeout: question.timeout,
          title: question.title,
        });
        break;
      }
      default:
        throw new Error("Should not come here");
    }
  };

  return (
    <div className="flex justify-between items-center">
      <FormField
        control={form.control}
        name={`questions.${index}.title`}
        render={({ field }) => (
          <FormItem className="w-2/3">
            <FormLabel
              htmlFor={`question-${index}-title`}
              className="text-card-foreground"
            >
              Título
            </FormLabel>
            <FormControl>
              <Input
                className="text-card-foreground bg-foreground border-accent"
                type="text"
                id={`question-${index}-title`}
                aria-invalid={!!form.formState.errors.questions}
                aria-errormessage={`question-${index}-title-error`}
                {...field}
              />
            </FormControl>
            <FormMessage id={`question-${index}-title-error`} />
          </FormItem>
        )}
      />
      <FormField
        control={form.control}
        name={`questions.${index}.timeout`}
        render={({ field }) => (
          <FormItem>
            <FormLabel
              htmlFor={`question-${index}-timeout`}
              className="text-card-foreground"
            >
              Tempo
            </FormLabel>
            <Select onValueChange={field.onChange} defaultValue={field.value}>
              <SelectTrigger className="text-card-foreground bg-foreground border-accent">
                <FormControl>
                  <SelectValue placeholder="Select a verified email to display" />
                </FormControl>
              </SelectTrigger>
              <SelectContent className="bg-background">
                <SelectItem value="15" className="focus:text-secondary">
                  15
                </SelectItem>
                <SelectItem value="30" className="focus:text-secondary">
                  30
                </SelectItem>
                <SelectItem value="60" className="focus:text-secondary">
                  60
                </SelectItem>
              </SelectContent>
            </Select>
            <FormMessage id={`question-${index}-timeout-error`} />
          </FormItem>
        )}
      />
      <FormField
        control={form.control}
        name={`questions.${index}.kind`}
        render={({ field }) => (
          <FormItem>
            <FormLabel
              htmlFor={`question-${index}-kind`}
              className="text-card-foreground"
            >
              Tipo
            </FormLabel>
            <Select
              onValueChange={(value) => {
                updateKind(value);
              }}
              defaultValue={field.value}
            >
              <SelectTrigger className="text-card-foreground bg-foreground border-accent">
                <FormControl>
                  <SelectValue placeholder="Select a verified email to display" />
                </FormControl>
              </SelectTrigger>
              <SelectContent className="bg-background">
                <SelectItem value="true_false" className="focus:text-secondary">
                  Verdadeiro ou Falso
                </SelectItem>
                <SelectItem value="quiz" className="focus:text-secondary">
                  Multiplas Escolhas
                </SelectItem>
              </SelectContent>
            </Select>
            <FormMessage id={`question-${index}-kind-error`} />
          </FormItem>
        )}
      />
      <IoMdTrash
        className="text-destructive h-6 w-6 cursor-pointer hover:scale-125 transition-all"
        onClick={() => remove(index)}
      />
    </div>
  );
}

type QuestionFactoryProps = {
  form: UseFormReturn<CreateGameFormSchema>;
  index: number;
  remove: UseFieldArrayRemove;
  update: UseFieldArrayUpdate<CreateGameFormSchema, "questions">;
};

function QuestionFactory({
  form,
  remove,
  index,
  update,
}: QuestionFactoryProps) {
  const question = useWatch({
    control: form.control,
    name: `questions.${index}`,
  });
  const factory = (
    form: UseFormReturn<CreateGameFormSchema>,
    index: number,
  ) => {
    switch (question.kind) {
      case "quiz":
        return <QuizQuestion form={form} index={index} />;

      case "true_false":
        return <TrueFalseQuestion form={form} index={index} />;
    }
  };

  return (
    <li className="flex flex-col gap-4 bg-foreground rounded-lg p-4">
      <QuestionHeader
        question={question}
        form={form}
        index={index}
        remove={remove}
        update={update}
      />
      {factory(form, index)}
    </li>
  );
}

export default QuestionFactory;

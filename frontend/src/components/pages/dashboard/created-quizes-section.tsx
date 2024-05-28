import { FaEdit } from "react-icons/fa";
import { FaPlay } from "react-icons/fa";

type QuizCardProps = {
  title: string;
  body: string;
};

function QuizCard({ title, body }: QuizCardProps) {
  return (
    <div className="flex flex-col bg-foreground text-secondary-foreground h-80 w-96 gap-4 p-4 rounded-lg overflow-y-auto">
      <h3 className="text-2xl font-bold">{title}</h3>
      <p className="font-light leading-6">{body}</p>
      <div className="mt-auto flex grow justify-start gap-8 py-8">
        <FaEdit className="w-8 h-8 cursor-pointer" />
        <FaPlay className="w-8 h-8 text-accent cursor-pointer" />
      </div>
    </div>
  );
}

function CreateQuizesSection() {
  return (
    <section className="flex flex-col gap-8 h-screen w-full px-4 overflow-y-hidden">
      <h1 className="text-4xl text-secondary-foreground">Quizzes criados</h1>
      <section className="flex flex-wrap grow w-full gap-4 overflow-y-auto">
        <QuizCard
          title="Test quiz 01"
          body="Lorem ipsum dolor sit amet, officia excepteur ex fugiat reprehenderit enim labore culpa sint ad nisi Lorem pariatur mollit ex esse exercitation amet. Nisi anim cupidatat excepteur officia."
        />
        <QuizCard
          title="Test quiz 01"
          body="Lorem ipsum dolor sit amet, officia excepteur ex fugiat reprehenderit enim labore culpa sint ad nisi Lorem pariatur mollit ex esse exercitation amet. Nisi anim cupidatat excepteur officia."
        />
        <QuizCard
          title="Test quiz 01"
          body="Lorem ipsum dolor sit amet, officia excepteur ex fugiat reprehenderit enim labore culpa sint ad nisi Lorem pariatur mollit ex esse exercitation amet. Nisi anim cupidatat excepteur officia."
        />
        <QuizCard
          title="Test quiz 01"
          body="Lorem ipsum dolor sit amet, officia excepteur ex fugiat reprehenderit enim labore culpa sint ad nisi Lorem pariatur mollit ex esse exercitation amet. Nisi anim cupidatat excepteur officia."
        />
        <QuizCard
          title="Test quiz 01"
          body="Lorem ipsum dolor sit amet, officia excepteur ex fugiat reprehenderit enim labore culpa sint ad nisi Lorem pariatur mollit ex esse exercitation amet. Nisi anim cupidatat excepteur officia."
        />
        <QuizCard
          title="Test quiz 01"
          body="Lorem ipsum dolor sit amet, officia excepteur ex fugiat reprehenderit enim labore culpa sint ad nisi Lorem pariatur mollit ex esse exercitation amet. Nisi anim cupidatat excepteur officia."
        />
      </section>
    </section>
  );
}

export default CreateQuizesSection;

import { FaEdit } from "react-icons/fa";
import { FaPlay } from "react-icons/fa";

export type QuizCardProps = {
  title: string;
  body: string;
};

export function QuizCard({ title, body }: QuizCardProps) {
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

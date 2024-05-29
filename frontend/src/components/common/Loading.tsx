import { CgSpinner } from "react-icons/cg";

function Loading() {
  return (
    <div className="h-fit flex flex-col items-center justify-center">
      <div className="flex items-center gap-4">
        <CgSpinner className="h-10 w-10 animate-spin text-secondary" />
        <h2 className="text-3xl font-bold text-secondary">Carregando...</h2>
      </div>
    </div>
  );
}

export default Loading;

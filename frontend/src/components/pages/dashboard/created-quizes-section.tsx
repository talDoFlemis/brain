"use client";

import { useEffect, useState } from "react";
import { QuizCard } from "./quiz-card";
import useSessionWithRefresh from "@/hooks/useSessionWithRefresh";
import Loading from "@/components/common/Loading";
import { Game } from "@/services/game/types";
import gameService from "@/services/game/service";

function CreateQuizesSection() {
  const [games, setGames] = useState<Game[]>([]);
  const [loading, setLoading] = useState<boolean>(true);

  const { data } = useSessionWithRefresh();

  useEffect(() => {
    if (!data) return;
    gameService
      .getGamesByUser(data.access_token)
      .then((response) => setGames(response.games))
      .catch((error) => {
        console.log(error);
        setGames([]);
      })
      .finally(() => {
        setLoading(false);
      });
  }, [data]);

  const content =
    games.length === 0 ? (
      <h2 className="text-white/80 text-2xl">Nenhum quiz criado no momento</h2>
    ) : (
      games.map((game) => {
        return (
          <QuizCard key={game.id} title={game.title} body={game.description} />
        );
      })
    );

  return (
    <section className="flex flex-col gap-8 h-screen w-full px-4 overflow-y-hidden">
      <h1 className="text-4xl text-secondary-foreground">Quizzes criados</h1>
      <section className="flex flex-wrap grow w-full gap-4 overflow-y-auto">
        {loading ? <Loading /> : content}
      </section>
    </section>
  );
}

export default CreateQuizesSection;

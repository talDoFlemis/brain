function PlayedGamesSection() {
  return (
    <section className="flex flex-col gap-8 h-screen w-full px-4 overflow-y-hidden">
      <h1 className="text-4xl text-secondary-foreground">Quizzes Jogados</h1>
      <section className="flex flex-wrap grow w-full gap-4 overflow-y-auto"></section>
    </section>
  );
}

export default PlayedGamesSection;

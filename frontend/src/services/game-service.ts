import apiProvider from "@/providers/api-provider";

export type Game = {
  id: string;
  title: string;
  description: string;
  owner_id: string;
  // TODO: Type each kind of question
  questions: any[];
};

export type GetGamesByUserResponse = {
  games: Game[];
};

export type GetGameByIdResponse = {
  game: Game;
};

const getGamesByUser = async (
  access_token: string,
): Promise<GetGamesByUserResponse> => {
  try {
    const response = await apiProvider.useGet<GetGamesByUserResponse>(
      "/game/",
      {
        headers: {
          Authorization: `Bearer ${access_token}`,
        },
      },
    );

    return response;
  } catch (error) {
    throw error;
  }
};

const getGameById = async (
  gameId: string,
  access_token: string,
): Promise<GetGameByIdResponse> => {
  try {
    const response = await apiProvider.useGet<GetGameByIdResponse>(
      `/game/${gameId}`,
      {
        headers: {
          Authorization: `Bearer ${access_token}`,
        },
      },
    );

    return response;
  } catch (error) {
    throw error;
  }
};

const gameService = {
  getGamesByUser,
  getGameById,
};

export default gameService;

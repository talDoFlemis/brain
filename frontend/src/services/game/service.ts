import apiProvider from "@/providers/api-provider";
import {
  CreateGameRequest,
  GetGameByIdResponse,
  GetGamesByUserResponse,
} from "./types";

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

const createGame = async (
  request: CreateGameRequest,
  access_token: string,
): Promise<boolean> => {
  try {
    await apiProvider.usePost<string, CreateGameRequest>(`/game/`, request, {
      headers: {
        Authorization: `Bearer ${access_token}`,
      },
    });
    return true;
  } catch (error) {
    throw error;
  }
};

const gameService = {
  getGamesByUser,
  getGameById,
  createGame,
};

export default gameService;

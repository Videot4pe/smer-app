import { atom } from "jotai";

import AuthApi from "../api/auth-api";
import type { ShortUser } from "../models/user";

const token = atom(localStorage.getItem("jwt-token") ?? undefined);
const refreshToken = atom(localStorage.getItem("jwt-token") ?? undefined);

export const jwtToken = atom(
  (get) => get(token),
  (get, set, newToken: string | undefined) => {
    set(token, newToken);
    if (!newToken) {
      localStorage.removeItem("jwt-token");
    } else {
      localStorage.setItem("jwt-token", newToken);
    }
  }
);

export const refreshJwtToken = atom(
  (get) => get(token),
  (get, set, newToken: string | undefined) => {
    set(refreshToken, newToken);
    if (!newToken) {
      localStorage.removeItem("refresh-token");
    } else {
      localStorage.setItem("refresh-token", newToken);
    }
  }
);

export const signinAtom = atom(
  (get) => get(jwtToken),
  async (_get, set, user: ShortUser) => {
    try {
      const newToken = await AuthApi.signin(user);
      set(jwtToken, newToken);
    } catch (error) {
      // eslint-disable-next-line no-console
      console.error(error);
    }
  }
);

export const isLoggedIn = atom<boolean>((get) => !!get(jwtToken));

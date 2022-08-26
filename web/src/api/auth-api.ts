import type { ShortUser, UserInfo } from "../models/user";

import { ApiClient } from "./client/api-client";

const client = new ApiClient("auth");

export interface JwtPayload {
  token: string;
  refreshToken: string;
}

export default {
  signin: (user: ShortUser) => client.post<JwtPayload>("/signin", user),
  signup: (user: ShortUser) => client.post<number>("/signup", user),
  refresh: (token: string) => client.post<JwtPayload>("/refresh", { token }),
  reset: (email: string) => client.post<number>("/password-reset", email),
  changePassword: (hash: string, password: string) => client.post<void>("/change-password", { hash, password }),
  info: () => client.get<UserInfo>("/info"),
};

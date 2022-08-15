import type { User } from "../models/user";

import { ApiClient } from "./client/api-client";

const client = new ApiClient("users");

export default {
  update: (user: User) => client.patch<number>("", user),
  remove: () => client.delete<number>(""),
  view: () => client.get<User>(""),
};

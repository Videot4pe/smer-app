export interface User {
  id?: number;
  avatar?: string | null;
  avatarId?: number | null;
  username: string;
  name: string;
  surname: string;
  patronymic?: string;
  email: string;

  createdAt?: string;
  updatedAt?: string;
}

export interface UserInfo {
  id: number;
  username: string;
  email: string;
}

export interface ShortUser {
  email: string;
  password: string;
}

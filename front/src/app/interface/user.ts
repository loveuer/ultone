import {Enum} from "./enum";

export interface User {
  id: number,
  created_at: number,
  updated_at: number,
  deleted_at: number,
  username: string,
  status: Enum<0 | 1>,
  nickname: string,
  comment: string,
  role: Enum<100 | 254 | 256>,
  privileges: Enum<number>[],
  created_by_id: number,
  created_by_name: string,
  active_at: number,
  deadline: number,
  login_at: number,
}

export interface NewUser {
  username: string,
  nickname: string,
  status: 0 | 1,
  password: string,
  privileges: number[],
  deadline: number,
  role: 100,
  comment: string,
}

import {Enum} from "./enum";

export interface Log {
  id: number,
  created_at: number,
  updated_at: number,
  deleted_at: number,
  user_id: number,
  username: string,
  type: Enum<number>,
  content: { [key: string]: any },
  html: string,
}

import {Injectable, signal, WritableSignal} from '@angular/core';
import { HttpClient } from "@angular/common/http";
import {Router} from "@angular/router";
import {User, NewUser} from "../interface/user";
import {Response} from "../interface/response";
import {tap} from "rxjs";
import {MsgService} from "./msg.service";

@Injectable({
  providedIn: 'root'
})
export class UserService {

  readonly auth_login_url = '/api/user/auth/login'
  readonly auth_logout_url = '/api/user/auth/logout'
  readonly manage_user_create_url = '/api/user/manage/user/create'
  readonly manage_user_update_url= '/api/user/manage/user/update'
  readonly manage_user_delete_url = '/api/user/manage/user/delete'
  readonly manage_user_list_url = '/api/user/manage/user/list'
  readonly init_user: User = {id: 0, username: ""} as User

  readonly user = signal(this.init_user)
  readonly user_list: WritableSignal<{ total: number, page: number, size: number, list: User[] }> = signal({
    total: 0,
    page: 0,
    size: 20,
    list: [] as User[],
  })

  constructor(
    private msg_srv: MsgService,
    private http: HttpClient,
    private router: Router,
  ) {
    this.auth_verify()
  }

  auth_login(username: string, password: string) {
    this.http.post<Response<{ token: string, user: User }>>(this.auth_login_url, {
      username: username,
      password: password
    }).subscribe(val => {
      if (val.status === 200) {
        localStorage.setItem("ult-token", val.data.token)
        this.user.set(val.data.user)
        this.router.navigate([''])
      }
    })
  }

  auth_verify() {
    this.http.get<Response<{ token: string, user: User }>>(this.auth_login_url).subscribe(rs => {
      if (rs.status === 200) {
        this.user.set(rs.data.user)
        console.log("[D] auth verify user=", this.user())
      }
    })
  }

  auth_logout() {
    this.http.post<Response<null>>(this.auth_logout_url, {}).subscribe(rs => {
      this.router.navigate(['login']).finally(() => {
        localStorage.removeItem("ult-token")
        this.user.set(this.init_user)
      })
    })
  }

  manage_user_create(new_user: NewUser) {
    return this.http.post<Response<User>>(this.manage_user_create_url, {...new_user}).pipe(
      tap({
        next: (rs) => {
          if (rs.status === 200) {
            this.manage_user_list()
          }
        }
      })
    )
  }

  manage_user_delete(target: User) {
    return this.http.post<Response<User>>(this.manage_user_delete_url, {id: target.id}).pipe(
      tap({
        next: (rs) => {
          if (rs.status === 200) {
            this.msg_srv.success('删除用户成功')
            this.manage_user_list()
          } else {
            this.msg_srv.error(rs.msg)
          }
        }
      })
    ).subscribe()
  }

  manage_user_update(body: Object) {
   return this.http.post<Response<any>>(this.manage_user_update_url, body).pipe(
     tap({next: (rs) => {
       if (rs.status === 200) {
         this.msg_srv.success("更新用户成功")
         this.manage_user_list()
       } else {
         this.msg_srv.error(rs.msg)
       }
       }})
   )
  }

  manage_user_list() {
    this.http.get<Response<{ total: number, list: User[] }>>(this.manage_user_list_url, {params: {}}).subscribe(rs => {
      if (rs.status === 200) {
        this.user_list.set({...this.user_list(), total: rs.data.total, list: rs.data.list})
        // this.user_list = {...this.user_list, total: rs.data.total, list: rs.data.list}
      } else {
        this.user_list.set({...this.user_list(), list: []})
      }
    })
  }
}

import {Injectable, signal} from '@angular/core';
import {Log} from "../interface/log";
import { HttpClient } from "@angular/common/http";
import {tap} from "rxjs";

@Injectable({
  providedIn: 'root'
})
export class LogService {

  readonly log_list_url = "/api/log/content/list"

  readonly logs = signal({
    list: [] as Log[],
    total: 0,
    page: 0,
    size: 0,
  })

  constructor(
    private http: HttpClient,
  ) {
  }

  get_logs() {
    return this.http.get<{status:number, msg:string, data: {list:Log[], total: number}}>(this.log_list_url).pipe(

    ).subscribe(rs => {
      if (rs.status === 200) {
        this.logs.set({...this.logs(), total: rs.data.total, list: rs.data.list})
      } else {
        this.logs.set({...this.logs(), list: []})
      }
    })}
}

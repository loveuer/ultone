import { HttpEvent, HttpHandler, HttpInterceptor, HttpRequest, HttpResponse } from '@angular/common/http';
import {map, Observable} from "rxjs";
import {Response} from "../interface/response";
import {MsgService} from "../service/msg.service";
import {Router} from "@angular/router";
import {Injectable} from "@angular/core";
import {UserService} from "../service/user.service";

// import {MsgService} from "../service/msg.service";


@Injectable({
  providedIn: 'root'
})
export class alerterInterceptor implements HttpInterceptor {
  constructor(
    private msg_srv: MsgService,
    private router: Router,
  ) {
  }

  intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    // const nr = req.clone().headers.set("Authorization", localStorage.getItem("ult-token"))
    req = req.clone({headers: req.headers.append("Authorization", localStorage.getItem("ult-token") ?? "")})
    return next.handle(req).pipe(
      map((event) => {
        if (event instanceof HttpResponse) {
          try {
            const rs = (event.body) as Response<any>
            console.log(`[D] ${req.method} - ${req.url} =>`, rs)
            if (rs.status > 200) {
              this.msg_srv.error(rs.msg)
              // alert(rs.msg)
            }

            if (rs.status === 401) {
              this.router.navigate(['login'])
            }
          } catch (e) {
           console.warn('[E] http err=',e)
            this.msg_srv.error('无法连接服务器')
          }
        }
        return event
      })
    );
  }
}

import {Routes} from '@angular/router';
import {HomeComponent} from "./page/home/home.component";
import {LoginComponent} from "./page/login/login.component";
import {LogComponent} from "./page/log/log.component";
import {UserComponent} from "./page/user/user.component";

export const routes: Routes = [
  {path: "", component: HomeComponent},
  {path: "login", component: LoginComponent},
  {path: "log", component: LogComponent},
  {path: "user", component: UserComponent},
];

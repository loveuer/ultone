import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import {ActivatedRoute, NavigationEnd, Router, RouterOutlet} from '@angular/router';
import {MatSidenav, MatSidenavContainer, MatSidenavContent} from "@angular/material/sidenav";
import {MatIcon} from "@angular/material/icon";
import {MatToolbar} from "@angular/material/toolbar";
import {MatButton, MatIconButton} from "@angular/material/button";
import {UserService} from "./service/user.service";
import {MatMenu, MatMenuItem, MatMenuTrigger} from "@angular/material/menu";

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, RouterOutlet, MatSidenavContent, MatSidenav, MatSidenavContainer, MatIcon, MatToolbar, MatIconButton, MatButton, MatMenuTrigger, MatMenu, MatMenuItem],
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss',
})
export class AppComponent {
  title = 'front';
  url = '';

  constructor(
    public router: Router,
    public user_srv: UserService,
  ) {
    this.router.events.subscribe(val => {
      if (val instanceof  NavigationEnd) {
        let _url = val.url
        if (_url.startsWith("/")) {
          _url = _url.slice(1)
        }
        _url = _url.split("/")[0]
        console.log('[D] router val=', _url)
        this.url = _url
      }
    })
  }
}

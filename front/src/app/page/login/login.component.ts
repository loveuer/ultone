import {Component, signal} from '@angular/core';
import {CommonModule} from '@angular/common';
import {FormsModule} from "@angular/forms";
import {UserService} from "../../service/user.service";
import {MatButton} from "@angular/material/button";
import {MatFormField, MatInput, MatInputModule} from "@angular/material/input";
import {MatIcon, MatIconModule} from "@angular/material/icon";
import {MatFormFieldModule} from "@angular/material/form-field";

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MatButton,
    MatInput,
    MatFormField,
    MatFormFieldModule,
    MatInputModule,
    MatIconModule,
    MatIcon,
  ],
  templateUrl: './login.component.html',
  styleUrl: './login.component.scss'
})
export class LoginComponent {

  username: string = ''
  password: string = ''

  constructor(
    public user_srv: UserService
  ) {
  }

  login() {
    this.user_srv.auth_login(this.username, this.password)
  }

  enter(event: KeyboardEvent) {
    if (event.key === 'Enter' && this.password && this.username) {
      this.login()
    }
  }
}
